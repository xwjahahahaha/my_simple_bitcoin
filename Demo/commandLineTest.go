package main

//go命令行测试

import (
"fmt"
"os"
)

//go命令行练习
func main()  {
	len1 := len(os.Args)
	fmt.Printf("命令长度为：%d\n", len1)
	for i, cmd := range os.Args{
		fmt.Printf("arg[%d]: %s\n", i, cmd)
	}

}