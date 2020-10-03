package main

import (
	"bytes"
	"fmt"
	"itcast_Go/bolt"
	"log"
)

//实现cli特定命令

//添加区块
func (cli *CLI)AddBlock(data string)  {
	//cil.bc.AddBlock(data) //TODO
}

//反向打印整个区块链
func (cli *CLI)PrintChainR()  {
	blockChain := cli.bc	//获取当前的区块链
	//创建迭代器
	iterator := blockChain.NewIterator()
	for {
		//调用迭代器的Next
		block := iterator.Next()
		fmt.Printf("================================\n")
		fmt.Printf("版本号：%d\n", block.Version)
		fmt.Printf("前区块hash值：%x\n", block.PrevBlockHash)
		fmt.Printf("本区块hash值：%x\n", block.ThisBlockHash)
		fmt.Printf("梅克尔根：%x\n", block.MerkelRoot)
		fmt.Printf("本区块时间：%d\n", block.TimeStamp)
		fmt.Printf("当前难度：%d\n", block.Difficulty)
		fmt.Printf("本区块随机数Nonce：%d\n", block.Nonce)
		fmt.Printf("区块的创世交易写的内容：%s\n", block.Transactions[0].Vin[0].ScriptSig)
		//输出区块中的所有交易
		//TODO

		//结束条件
		if len(block.PrevBlockHash) == 0{
			fmt.Println("区块链遍历结束！")
			break
		}
	}
}

//正向打印整个区块链数据
func (cli *CLI)PrintChain()  {
	blockChain := cli.bc
	db := blockChain.db 	//获得区块链的区块blot数据库
	//直接查询数据库
	blockNumber	:= 0	//区块高度
	db.View(func(tx *bolt.Tx) error {
		//2. 打开抽屉(没有就创建)
		var bucketName = []byte(BucketName)
		bucket := tx.Bucket(bucketName)
		if bucket == nil {
			var err error
			//没有就创建
			bucket, err = tx.CreateBucket(bucketName)
			if err != nil {
				log.Panic("打开桶", BucketName, "失败！", err)
				return err
			}
		}
		//遍历：从第一个key->value进行遍历，到最后一个固定的key直接返回
		bucket.ForEach(func(k, v []byte) error {
			if bytes.Equal(k, []byte("lastHashKey")) {
				fmt.Println("区块链遍历结束！")
				//最后一个固定的key直接返回
				return nil
			}
			//将区块反序列化
			block := Deserialize(v)
			//输出区块内容
			fmt.Printf("=================区块高度:%d===============\n", blockNumber)
			fmt.Printf("版本号：%d\n", block.Version)
			fmt.Printf("前区块hash值：%x\n", block.PrevBlockHash)
			fmt.Printf("本区块hash值：%x\n", block.ThisBlockHash)
			fmt.Printf("梅克尔根：%x\n", block.MerkelRoot)
			fmt.Printf("本区块时间：%d\n", block.TimeStamp)
			fmt.Printf("当前难度：%d\n", block.Difficulty)
			fmt.Printf("本区块随机数Nonce：%d\n", block.Nonce)
			fmt.Printf("区块的创世交易写的内容：%s\n", block.Transactions[0].Vin[0].ScriptSig)
			//输出区块中的所有交易
			//TODO


			blockNumber++
			return nil
		})
		return nil
	})
}

//查询账户余额
func (cli *CLI)getBalance(address string) uint64{
	//获取UTXOS(所有的未支付输出)
	utxos := cli.bc.FindUTXOs(address)
	//总和余额
	total := 0.0
	//计算余额
	for _,utxo := range utxos {
		total += utxo.Value
	}
	fmt.Printf("账户%s总余额是%f\n", address, total)
	return 0
}

//转账命令
func (cli *CLI) Send(from, to string, amount float64, miner string, commitment string){
	//1. 创建挖矿交易
	coinbaseTx := createCoinBaseTx(miner, commitment)
	//2. 发起新交易
	newTx := NewTransaction(from, to, amount, cli.bc)
	//3. 打包到区块
	cli.bc.AddBlock([]*Transaction{coinbaseTx, newTx})
}