package main

import (
	"fmt"
	"os"
)

//接收命令行参数，处理命令行。操作区块链
type CLI struct {
	bc BlockChain	//要操作的区块链
}

const Usage  = `
	addBlock --data DATA	'add data to BlockChain'
	printChain	'print all data of BlockChain'
`

//接收参数的处理函数
func (cli *CLI)Run()  {
	//1. 接收所有命令
	args := os.Args
	//2. 分析所有命令
	//去除无效输入
	if len(args) < 2{
		fmt.Println("命令参数数量出错！")
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
		case "printChain":
			//打印整条区块链
			if len(args) == 2 {
				cli.PrintChain()
			}else {
				fmt.Println("printChain命令参数出错！")
			}
		default:
			//不是命令则输入提示信息
			fmt.Println("无效的选项！")
			fmt.Printf(Usage)
	}
}
