package main

import (
	"encoding/hex"
	"fmt"
	"log"
	"time"
)

//实现cli特定命令


//反向打印整个区块链
func (cli *CLI)PrintChainR()  {
	blockChain := cli.bc	//获取当前的区块链
	//创建迭代器
	iterator := blockChain.NewIterator()
	for {
		//调用迭代器的Next
		block := iterator.Next()
		fmt.Printf("================================================================================================\n")
		fmt.Printf("版本号：%d\n", block.Version)
		fmt.Printf("前区块hash值：%x\n", block.PrevBlockHash)
		fmt.Printf("本区块hash值：%x\n", block.ThisBlockHash)
		fmt.Printf("梅克尔根：%x\n", block.MerkelRoot)
		//转换下时间
		timeFormat := time.Unix(int64(block.TimeStamp), 0).Format("2006-01-02 15:04:05")
		fmt.Printf("本区块时间：%s\n", timeFormat)
		fmt.Printf("当前难度：%d\n", block.Difficulty)
		fmt.Printf("本区块随机数Nonce：%d\n", block.Nonce)
		fmt.Printf("区块的创世交易写的内容：%s\n", block.Transactions[0].Vin[0].PubKey)
		fmt.Printf("挖矿人：%s\n", block.Miner)
		//输出区块中的所有交易
		fmt.Println("区块中的交易:")
		for _, tx := range block.Transactions{
			fmt.Println(tx)			//格式化打印交易
		}
		//结束条件
		if len(block.PrevBlockHash) == 0{
			fmt.Println("区块链遍历结束！")
			break
		}
	}
}

//查询账户余额
func (cli *CLI)getBalance(address string) {
	// 校验地址
	if !checkAddress(address){
		fmt.Printf("[%s]地址无效!\n", address)
		return
	}

	//获取UTXOS(所有的未支付输出)
	//地址类型转换
	utxos := cli.bc.FindUTXOs(adsToPubKeyHash(address))
	//总和余额
	total := 0.0
	//计算余额
	for _,utxo := range utxos {
		total += utxo.Value
	}
	fmt.Printf("账户[%s]总余额是%f\n", address, total)
	fmt.Println("余额查询结束...")
}

//查询交易信息
func (cli * CLI) getTransaction(TxHash string)  {
	//1. 找到此交易
	TxHashBytes, _ := hex.DecodeString(TxHash)	//把16进制字符串转换为[]byte
	tx, err := FindTxByTxHash(TxHashBytes, &cli.bc)
	if err != nil{
		log.Panic("找不到此交易！", err)
	}
	//2. 打印
	fmt.Println(tx)
}


//转账命令
func (cli *CLI) Send(from, to string, amount float64, miner string, commitment string){
	//0. 校验
	if !checkAddress(from){
		fmt.Printf("from : [%s]地址无效!\n", from)
		return
	}
	if !checkAddress(to){
		fmt.Printf("to : [%s]地址无效!\n", to)
		return
	}
	if !checkAddress(miner){
		fmt.Printf("miner : [%s]地址无效!\n", miner)
		return
	}
	if amount < 0 {
		fmt.Println("转账金额无效！")
		return
	}
	//1. 创建挖矿交易
	coinbaseTx := createCoinBaseTx(miner, commitment)
	//2. 发起新交易
	// 把地址转换成为
	newTx := NewTransaction(from, to, amount, &cli.bc)
	if newTx == nil{
		//如果新交易为空，那么是余额不足，直接终止
		return
	}
	//3. 打包到区块
	cli.bc.AddBlock(miner, []*Transaction{coinbaseTx, newTx})
	fmt.Printf("生成的交易Hash:%x\n", newTx.TxHash)
	fmt.Println("转账结束...")
	fmt.Println("===========================================================")
}

//读取钱包的文件名
const WalletsFileName = "wallets.dat"

//添加新钱包内容 (公私钥对)
func (cli *CLI) newWallet() {
	//打开钱包
	ws := OpenWallets(WalletsFileName)
	//创建钱包结构
	address := ws.createWallet(WalletsFileName)
	fmt.Printf("新账户地址: %s\n", address)
	//列举下目前的所有地址
	cli.listAddress()
}

//列举钱包中的所有地址
func (cli *CLI) listAddress()  {
	//打开钱包
	ws := OpenWallets(WalletsFileName)
	//得到地址数组
	//遍历钱包中的地址
	fmt.Println("当前钱包所有地址:")
	for _, address := range ws.GetAllAddress(){
		fmt.Printf("地址: %s\n", address)
	}
}