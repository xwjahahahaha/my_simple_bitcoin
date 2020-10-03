package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"errors"
	"fmt"
	"log"
	"math/big"
	"strings"
	"time"
)

const reward = 50		//目前的挖矿奖励

//1. 定义交易结构
type Transaction struct {
	TxHash []byte		//交易hash
	//一个交易中有多个输入和多个输出
	Vin []TxInput		//交易输入
	Vout []TxOutput		//交易输出
	TimeStamp uint64	//产生时间
}

//定义交易输入结构
type TxInput struct {
	TxHash []byte		//上一个输出所在交易的hash
	OutputIndex int64	//上一个输出的编号 就是其在VOut中的序号
	//此版本使用的是P2PKH的脚本校验方式
	ScriptSig []byte	//私钥签名
	PubKey []byte 		//公钥
}
//定义交易输出结构
type TxOutput struct {
	Value float64		//输出的金额
	Index int64			//输出的编号(在此交易中)
	TxHash []byte		//该输出所属于的交易hash
	PubKeyHash []byte	//公钥的Hash
}

//地址转其公钥的hash函数
//（地址是由公钥计算过来的, 可以逆推回去到公钥的hash，但是无法逆推到原公钥，原公钥无法逆推到私钥，因为hash函数不可逆）
func (Output *TxOutput)Lock(address string)  {
	// 1.base58函数的解码
	// 2.去除尾部添加的4byte校验码和首部添加的1byte版本号
	// 3.赋值给Output
	Output.PubKeyHash = adsToPubKeyHash(address)
}

//TxOutput的创建函数
//因为其中要填充PubKeyHash要调用Lock函数，所以要一个函数取实现创建,以此实现内部调用
func NewTxOutput(value float64, index int64, txHash[]byte, address string) *TxOutput {
	output := TxOutput{
		Value:      value,
		Index:      index,
		TxHash:     txHash,
		PubKeyHash: nil,
	}
	//调用lock得出公钥的hash
	output.Lock(address)
	return &output
}


//设置交易hash
func (tx * Transaction)SetTxHash() {
	//采用序列化
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(tx)
	if err != nil{
		log.Panic("设置交易hash失败！", err)
	}
	data := buffer.Bytes()
	//sha256
	hash := sha256.Sum256(data)
	//赋值
	tx.TxHash = hash[:]
}


//2. 提供创建交易方法
//创建普通交易
func NewTransaction(from, to string, amount float64, bc *BlockChain) *Transaction {
	//创建交易之后要进行数字签名，签名需要私钥，这里新的交易的输入也需要公钥，所以需要打开存储的钱包文件获取公私钥
	// 1. 打开钱包，获取钱包
	ws := OpenWallets(WalletsFileName)
	// 2. 使用地址匹配到钱包中对应的公私钥对 (注意，这里的地址只能是本地管理的地址，因为只有这样才有公私钥)
	wallet := ws.WalletsMap[from]	//from地址不一定是本地管理的地址，所以需要判断返回的指针
	if wallet == nil{
		fmt.Println("非本地钱包管理地址，交易创建失败！")
		return nil
	}
	pubKey := wallet.PublicKey		// 这个公钥是创建交易的发送方的公钥

	//1. 找到合理的UTXO集合，用于转账. utxo是一个map，key为交易哈希，val为交易编号数组
	// 将公钥做hash去和output中的PubKeyHash作比对
	pubKeyHash := HashPubKey(pubKey)
	utxosMap, resAmount := bc.FindSuitableUTXO(pubKeyHash, amount)
	//判断最合适的UTXO是否满足需求，不满足直接返回
	if resAmount < amount{
		fmt.Printf("需转账账户[%s]当前的最高余额是:%.4f不足以支付\n", from, resAmount)
		return nil
	}
	//2. 把这些交易转换为input输入
	var inputArray []TxInput 
	var outputArray []TxOutput
	// 遍历map
	for txHash, indexArray := range utxosMap{
		// 遍历编号数组
		for _, index := range indexArray{
			//转换
			input := TxInput{
				TxHash:      []byte(txHash),
				OutputIndex: index,
				// 私钥签名
				ScriptSig:   []byte(from),		//TODO 这里要解锁验证
				//公钥
				PubKey:      pubKey,			// 注意这里的原始公钥不作为签名的直接来源，不然签名的验证将无任何意义
			}
			//加入到输入组中:
			inputArray = append(inputArray, input)
		}
	}
	//3. 再创建对应的输出
	//amount这的钱是转账需要的钱, 这里的index肯定是第一个，因为是第一个输出
	//TxHash先写空，等交易id算出来了在添加
	//把输出者的地址传入即可，NewTxOutput内部自动会回推公钥hash
	output := NewTxOutput(amount, 0, nil, to)
	outputArray = append(outputArray, *output)
	//判断是否要找零
	if resAmount > amount {
		//找零, 添加一个输出
		//resAmount-amount找零的钱, 这是第二个输出index就是1, 交易先写空，等交易id算出来了在添加
		//找零的output转给自己，所以是from
		backOutput := NewTxOutput(resAmount-amount, 1, nil, from)
		//添加进输出组
		outputArray = append(outputArray, *backOutput)
	}
	//4. 生成交易
	newTx := Transaction{
		TxHash:    nil,
		Vin:       inputArray,
		Vout:      outputArray,
		TimeStamp: uint64(time.Now().Unix()),
	}
	newTx.SetTxHash()	//设置交易的hash
	newTxHash := newTx.TxHash
	// 注意Go语言和python一样循环修改数据的值必须使用索引！！！！
	for i := range outputArray{
		outputArray[i].TxHash = newTxHash
	}
	newTx.Vout = outputArray 	//重新再塞进去，不要忘记！

	//对当前交易进行签名
	// 1. 获取私钥
	privateKey := wallet.Private
	// 2. 对每个input生成签名
	newTx.Signature(privateKey, bc)
	return &newTx
}


//3. 创建挖矿交易
//挖矿交易
func createCoinBaseTx(address string, commitment string)  *Transaction{
	//创建coinbase
	input := TxInput{
		TxHash:      nil,					//铸币交易没有输入
		OutputIndex: -1,					//没有上一个输出
		ScriptSig:   nil,					//coinbase的input不需要验证也就不需要签名
		PubKey:		[]byte(commitment),		//本来是写上一个output的公钥，但这里没有上一个，所以是coinBase域，可以写想保存的内容
	}
	//reward铸币交易, Index是第一个, TxHash先按空值去计算交易的hash，算出来了在填充回去
	// 直接传字符串地址函数内部创建会进行地址=>公钥hash的转换
	Output := NewTxOutput(reward, 0, nil, address)
	tx :=  &Transaction{
		TxHash:    nil,
		Vin:       []TxInput{input},		//把铸币交易放进区块中
		Vout:      []TxOutput{*Output},
		TimeStamp: uint64(time.Now().Unix()),
	}
	//设置交易的Hash值
	tx.SetTxHash()
	//设置output的TxHash, 再塞进Vout中！
	Output.TxHash = tx.TxHash
	tx.Vout = []TxOutput{*Output}		//再塞进Vout中
	return tx
}

//判断一个交易是否是挖矿交易
func (tx *Transaction)isCoinbaseTx() bool {
	//条件
	//1. 交易只有一个输入
	//2. 输入中上一个交易hash为空
	//3. 输入中的上一个输出编号为-1
	if len(tx.Vin) == 1 && bytes.Equal(tx.Vin[0].TxHash , []byte{}) && tx.Vin[0].OutputIndex == -1{
		return true
	}
	return false
}



//签名的实现
//参数：账户的私钥
func (tx *Transaction) Signature(privateKey *ecdsa.PrivateKey, bc *BlockChain)  {
	if tx.isCoinbaseTx(){
		return		//铸币交易不需要签名
	}
	// 1.复制一份input，获取其output的pubKeyHash赋值到input的pubKey中(只是为了计算签名)
	trimmedCopyTx := tx.TrimmedCopy(bc)
	//对于每个交易中的input其自己生成的output就在同交易中并且其中自带pubKeyHash和转账金额
	// 2.对这个整体交易做hash，赋值到input的签名中
	for i, input := range trimmedCopyTx.Vin{
		//2.1 找到每个input关联的上一个output的公钥hash，并添加到当前的input的pubKey中
		//获取上一个交易
		preTx, err := FindTxByTxHash(input.TxHash, bc)
		if err != nil{
			log.Panic("签名时查找相关输出交易出错！", err)
		}
		//获取该交易中的output中的公钥hash
		prePubKeyHash := preTx.Vout[input.OutputIndex].PubKeyHash
		//注意！在这里直接对input赋值是无效的！！！
		trimmedCopyTx.Vin[i].PubKey = prePubKeyHash
		//2.2 签名需要的数据都具备了,做hash处理
		trimmedCopyTx.SetTxHash()	//交易的hash就是需要的签名数据
		signDataHash := trimmedCopyTx.TxHash
		//2.3 重要的一步！把当前交易中的这个input的pubKey还原为空，保证不影响其他input的签名
		trimmedCopyTx.Vin[i].PubKey = nil
		//2.4 执行签名动作得到r，s字节流
		r, s, err := ecdsa.Sign(rand.Reader, privateKey, signDataHash)
		if err != nil{
			log.Panic(err)
		}
		//2.5 把签名放到原本交易的ScriptSig中
		signnature := append(r.Bytes(), s.Bytes()...)
		tx.Vin[i].ScriptSig = signnature
	}
}

// 根据交易hash查找交易
func FindTxByTxHash(txHash []byte, bc *BlockChain) (*Transaction, error) {
	//遍历区块链
	iterator := bc.NewIterator()
	for{
		block := iterator.Next()
		//遍历区块中的交易数组
		for _, tx := range block.Transactions{
			if bytes.Equal(tx.TxHash, txHash){
				//如果两个hash相同就找到了返回
				return tx, nil
			}
		}
		//结束条件
		if len(block.PrevBlockHash) == 0{
			break
		}
	}
	return nil, errors.New("查询的交易不存在！")
}

//特殊的copy，为了计算签名
func (tx *Transaction) TrimmedCopy(bc *BlockChain) Transaction {
	var copyInputs []TxInput
	var copyOutputs []TxOutput
	// 把此交易中每个input中的签名都设置为空，pubKey设置为其关联的上一个输出的公钥hash
	for _, input := range tx.Vin{
		//添加，设置签名和公钥为空
		copyInputs = append(copyInputs, TxInput{
			TxHash:      input.TxHash,
			OutputIndex: input.OutputIndex,
			ScriptSig:   nil,
			PubKey:      nil,  		//这两个都需要设置为空值
		})
	}
	for _, output := range tx.Vout{ //保留 不变
		copyOutputs = append(copyOutputs, output)
	}
	return Transaction{
		TxHash:    tx.TxHash, 		//保留
		Vin:       copyInputs,
		Vout:      copyOutputs,
		TimeStamp: tx.TimeStamp,	//保留
	}
}


//验证交易
func (tx *Transaction) Verify(bc *BlockChain) bool {
	if tx.isCoinbaseTx(){
		return true			//铸币交易无需验证
	}
	//1. 获取验证所需要的数据
	// 1.1 Data
	trimmedCopy := tx.TrimmedCopy(bc)
	for i, input := range tx.Vin{  //注意，遍历的是原本的交易
		preTx, err := FindTxByTxHash(input.TxHash, bc)
		if err != nil{
			log.Panic(err)
		}
		trimmedCopy.Vin[i].PubKey = preTx.Vout[input.OutputIndex].PubKeyHash
		//计算hash
		trimmedCopy.SetTxHash()
		//a. Data得到
		dataHash := trimmedCopy.TxHash
		// 还原
		trimmedCopy.Vin[i].PubKey = nil
		//b. 签名得到
		signature := input.ScriptSig
		//c. 公钥
		//拆解PubKey， X, Y得到原生公钥
		PubKey := input.PubKey

		//拆开签名，得到r和s
		r1 := big.Int{}
		s1 := big.Int{}
		//r是前半部分，s是后半部分
		r1.SetBytes(signature[:len(signature)/2])
		s1.SetBytes(signature[len(signature)/2:])

		//拆开公钥，得到x和y
		x := big.Int{}
		y := big.Int{}
		//r是前半部分，s是后半部分
		x.SetBytes(PubKey[:len(PubKey)/2])
		y.SetBytes(PubKey[len(PubKey)/2:])
		//得到公钥原型
		pubKeyOrigin := ecdsa.PublicKey{elliptic.P256(), &x, &y}

		//1.2 verify
		if !ecdsa.Verify(&pubKeyOrigin, dataHash, &r1, &s1){
			return false			//一旦有一个input验证错误就失败
		}
	}
	return true
}

func (tx *Transaction) String() string {
	var lines []string
	lines = append(lines, fmt.Sprintf("--- Transaction:  【%x】", tx.TxHash))
	for i, input := range tx.Vin{
		lines = append(lines, fmt.Sprintf("  -Input    :   %d", i))
		lines = append(lines, fmt.Sprintf("    TxHash  :   %x", input.TxHash))
		lines = append(lines, fmt.Sprintf("    Out     :   %d", input.OutputIndex))
		lines = append(lines, fmt.Sprintf("    Sig     :   %x", input.ScriptSig))
		lines = append(lines, fmt.Sprintf("    PubKey  :   %x", input.PubKey))
	}

	for i, output := range tx.Vout{
		lines = append(lines, fmt.Sprintf("  -Output   :    %d", i))
		lines = append(lines, fmt.Sprintf("    Value   :    %f", output.Value))
		lines = append(lines, fmt.Sprintf("    Script  :    %x", output.PubKeyHash))
		lines = append(lines, fmt.Sprintf("    TxHash  :    %x", output.TxHash))
		lines = append(lines, fmt.Sprintf("    Index   :    %d", output.Index))
	}
	return strings.Join(lines, "\n")
}