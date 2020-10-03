package main

import (
	"bytes"
	"encoding/gob"
	"log"
	"time"
)


//1.定义结构
type Block struct{
	//1. 版本号
	Version uint64
	//2.前一个区块的hash
	PrevBlockHash []byte
	//3. 梅克尔根
	MerkelRoot []byte
	//4. 时间戳
	TimeStamp uint64
	//5. 难度值
	Difficulty uint64
	//6. nonce随机数(挖矿要找的数)
	Nonce uint64

	//a.本区块的hash   => bitcoin系统实际没有此字段，这里为了简化操作
	ThisBlockHash []byte
	//b.打包的所有交易->数组
	Transactions []*Transaction
}

//2.创建区块
func NewBlock(txs []*Transaction, prevBlockHash []byte)  *Block{
	block := Block{
		Version : 00,
		PrevBlockHash : prevBlockHash,
		MerkelRoot: []byte{},
		TimeStamp: uint64(time.Now().Unix()),
		Difficulty: 0,
		Nonce: 0,
		ThisBlockHash : []byte{}, 
		Transactions : txs,
	}
	//计算区块中包含所有交易的MerkelRoot
	block.MerkelRoot = block.CalculMerkelRoot()
	//需要不断做pow获取本区块hash和随机值
	pow := NewProofOfWork(&block)
	//不断做Hash运算
	thisBlockHash, nonce := pow.Run()
	//算出符合难度的nonce后计算hash
	//给区块hash和nonce赋值
	block.ThisBlockHash = thisBlockHash
	block.Nonce = nonce


	return &block 
}

//将区块的数据转换成字节流（序列化/编码）
func (block *Block)Serialize() []byte {
	//创建字节流
	var buffer bytes.Buffer
	//创建编码器
	encoder := gob.NewEncoder(&buffer)
	//编码
	err := encoder.Encode(&block)
	if err != nil{
		log.Panic("编码错误！", err)
		return nil
	}
	return buffer.Bytes()
}

//将区块的数据转换成字节流（反序列化/解码）
func Deserialize(bytesArry []byte) Block {
	//创建block
	var block Block
	//创建Reader流,创建解码器
	decoder := gob.NewDecoder(bytes.NewReader(bytesArry))
	err := decoder.Decode(&block)
	if err != nil{
		log.Panic("解码错误!", err)
	}
	return block
}

//计算梅克尔树的树根，这里先只是简单的累计计算		//TODO
//func (block *Block)SetMerkelRoot() {
//	for i, tx := range block.Transactions {
//		//统计串上次交易中所有的交易输入,得到交易输入[]byte
//		var inputInfo []byte
//		var outputInfo []byte
//		for i, input := range tx.Vin {
//			//组合一个交易中的所有输入
//			inputTmp := [][]byte{
//				input.TxHash,
//				Uint64ToByte(uint64(input.OutputIndex)),
//				input.ScriptSig,
//			}
//			//加到此交易组的总byte中
//			inputInfo = bytes.Join(inputTmp, []byte{})
//		}
//
//		for i, output := range tx.Vout {
//			//组合一个交易中的所有交易输出
//			outputTmp := [][]byte{
//				Uint64ToByte(uint64(output.Value)),			//TODO
//				Uint64ToByte(output.Index),
//				output.scriptPubKey,
//			}
//			//加到此交易组的总byte中
//			outputInfo = bytes.Join(outputTmp, []byte{})
//		}
//
//		//二维切片放置每一个一维[]byte数组,以后需要加参数在此添加即可
//		tmp := [][]byte{
//			Uint64ToByte(tx.TxId),
//			tx.TxHash,
//			inputInfo, 	//交易输入组bytes
//			outputInfo,	//交易输出组bytes
//			Uint64ToByte(tx.TimeStamp),
//		}
//		//将二维切片数组连接成一维的切片[]byte，二维数组元素中间的连接内容为空（第二个参数）
//		txInfo := bytes.Join(tmp, []byte{})
//		//2.sha256
//		hash := sha256.Sum256(txInfo)
//		//设置树根的hash
//		block.MerkelRoot = hash[:]
//	}
//}

//计算梅克尔树的树根，这里先只是简单的累计计算
func (block *Block) CalculMerkelRoot() []byte{
	//TODO
	return []byte{}
}