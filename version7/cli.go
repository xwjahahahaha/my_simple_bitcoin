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
	addBlock --data DATA			'添加区块数据'
	printChainR				'反向打印区块数据'
	printChain				'正向打印区块链数据'
	getBalance ADDRESS			'查询账户余额'
	send FROM TO AMOUNT MINER DATA 		'由FROM转AMOUNT钱给TO，由MINER挖矿，同时写入DATA'
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
		//3. 操作区块链  这里都是逻辑控制，具体过程在commandLine中实现
		case "addBlock" :
			//添加区块
			if len(args) == 4 && args[2]=="--data"{
				cli.AddBlock(args[3])
			}else{
				fmt.Println("addBlock命令参数出错！")
			}
		case "printChainR":
			fmt.Println("反向打印整条区块链")
			//反向打印整条区块链
			if len(args) == 2 {
				cli.PrintChainR()
			}else {
				fmt.Println("printChainR命令参数出错！")
			}
		case "printChain":
			fmt.Println("正向打印整条区块链")
			//正向打印整条区块链
			if len(args) == 2 {
				cli.PrintChain()
			}else {
				fmt.Println("printChain命令参数出错！")
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
		default:
			//不是命令则输入提示信息
			fmt.Println("无效的选项！")
			fmt.Printf(Usage)
	}
}
