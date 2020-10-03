package main

//第二个版本：完整字段
func main()  {
	blockChain01 := NewBlockChain("文杰") //TODO  这里挖矿人的地址随便写的
	cli := CLI{*blockChain01}
	cli.Run();
}  