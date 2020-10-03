package main

//第二个版本：完整字段
func main()  {
	blockChain01 := NewBlockChain()
	cli := CLI{*blockChain01}
	cli.Run();
}  