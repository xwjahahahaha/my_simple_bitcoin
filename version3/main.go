package main

import "fmt"
//第二个版本：完整字段
func main()  {
	blockChain01 := NewBlockChain()
	//加点数据
	blockChain01.AddBlock("文杰")
	blockChain01.AddBlock("小阳")
	blockChain01.AddBlock("剂量")
	for i, block := range blockChain01.blocks{
		fmt.Printf("================\n当前区块的高度：%d\n", i)
		fmt.Printf("前区块hash值：%x\n", block.PrevBlockHash)
		fmt.Printf("本区块hash值：%x\n", block.ThisBlockHash)
		//fmt.Printf("本区块时间：%t\n", block.TimeStamp)
		fmt.Printf("区块数据：%\n", block.Data)
	}
}  