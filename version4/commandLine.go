package main

import "fmt"

//实现cli特定命令

//添加区块
func (cil *CLI)AddBlock(data string)  {
	cil.bc.AddBlock(data)
}

//打印整个区块链
func (cil *CLI)PrintChain()  {
	blockChain := cil.bc	//获取当前的区块链
	//创建迭代器
	iterator := blockChain.NewIterator()
	for {
		//调用迭代器的Next
		block := iterator.Next()
		fmt.Printf("===============\n前区块hash值：%x\n", block.PrevBlockHash)
		fmt.Printf("本区块hash值：%x\n", block.ThisBlockHash)
		fmt.Printf("本区块时间：%v\n", block.TimeStamp)
		fmt.Printf("区块数据：%s\n", block.Data)

		//结束条件
		if len(block.PrevBlockHash) == 0{
			fmt.Println("区块链遍历结束！")
			break
		}
	}
}