package main

import (
	"fmt"
	"os"
	"strconv"
)

//接收命令行参数，处理命令行。操作区块链
type CLI struct {
	bc BlockChain	//要操作的区块链
}

const Usage  = `
	printChain				'打印区块数据'
	getBalance ADDRESS			'查询账户余额'
	getTransaction TXHASH 				'查询交易信息'
	send FROM TO AMOUNT MINER DATA 		'由FROM转AMOUNT钱给TO，由MINER挖矿，同时写入DATA'
	newWallet				"创建一个新的钱包（公私钥对）"
	listAddress				"列举所有的钱包地址"
	
`

//接收参数的处理函数
func (cli *CLI)Run()  {
	//1. 接收所有命令
	args := os.Args
	//2. 分析所有命令
	if len(args)<2  {
		fmt.Printf(Usage)
		return
	}
	command := args[1]
	switch command {
		case "printChain":
			fmt.Println("打印整条区块链")
			//反向打印整条区块链
			if len(args) == 2 {
				cli.PrintChainR()
			}else {
				fmt.Println("printChainR命令参数出错！")
			}
		case "getBalance":
			fmt.Println("===========================================================")
			fmt.Println("查询账户余额...")
			if len(args) == 3{
				//执行查询余额部分
				cli.getBalance(args[2])
			}else{
				fmt.Println("getBalance命令参数出错！")
			}
		case "getTransaction":
			fmt.Println("===========================================================")
			fmt.Println("查询交易信息...")
			if len(args) == 3{
				cli.getTransaction(args[2])
			}else {
				fmt.Println("getTransaction命令参数出错！")
			}
		case "send":
			fmt.Println("===========================================================")
			fmt.Println("转账开始...")
			if len(args) == 7{
				//拿参数
				from := args[2]
				to := args[3]
				amount, _ := strconv.ParseFloat(os.Args[4], 64) //string类型的参数转换为float64
				miner := args[5]
				commitment := args[6]
				//执行发送交易
				cli.Send(from, to, amount, miner, commitment)
			}else{
				fmt.Println("send命令参数出错！")
			}
		case "newWallet":
			fmt.Println("===========================================================")
			fmt.Println("开始打开钱包内容...")
			cli.newWallet()
		case "listAddress":
			fmt.Println("===========================================================")
			fmt.Println("读取钱包中的所有地址...")
			cli.listAddress()
		default:
			//不是命令则输入提示信息
			fmt.Println("无效的选项！")
			fmt.Printf(Usage)
	}
}

