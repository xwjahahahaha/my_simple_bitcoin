package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
)

//创建“人”结构体
type Person struct {
	Name string
	Age uint
}

func main()  {
	//定义一个“人”结构
	var xiaoming Person
	xiaoming.Name = "小明"
	xiaoming.Age = 18
	//编码的数据放进buffer
	var buffle bytes.Buffer

	//使用gob序列化得到字节流
	//定义一个编码器encoder
	encoder := gob.NewEncoder(&buffle)
	//编码结构体
	err := encoder.Encode(&xiaoming)
	if err != nil{
		log.Panic(err)
	}
	fmt.Printf("小明编码后的结果为：%v\n", buffle.Bytes())

	//使用gob反序列化得到结构体
	//创建byte读input流,然后创建解码器
	decoder := gob.NewDecoder(bytes.NewReader(buffle.Bytes()))
	var daming Person
	//解码
	err = decoder.Decode(&daming)
	if err != nil{
		log.Panic(err)
	}

	fmt.Printf("解码后的小明: %v\n", daming)

}