package main

import (
	"bytes"
	"encoding/gob"
	"log"
	"time"
)


//1.定义结构
type Block struct{
	//1. 版本号
	Version uint64
	//2.前一个区块的hash
	PrevBlockHash []byte
	//3. 梅克尔根   => v4版本再完善
	MerkelRoot []byte
	//4. 时间戳
	TimeStamp uint64
	//5. 难度值
	Difficulty uint64
	//6. nonce随机数(挖矿要找的数)
	Nonce uint64

	//a.本区块的hash   => bitcoin系统实际没有此字段，这里为了简化操作
	ThisBlockHash []byte
	//b.数据
	Data []byte
}

//2.创建区块
func NewBlock(data string, prevBlockHash []byte)  *Block{
	block := Block{
		Version : 00,
		PrevBlockHash : prevBlockHash,
		MerkelRoot: []byte{},
		TimeStamp: uint64(time.Now().Unix()),
		Difficulty: 0,
		Nonce: 0,
		ThisBlockHash : []byte{}, 
		Data : []byte(data),
	}
	//block.SetHash()
	//需要不断做pow获取本区块hash和随机值
	pow := NewProofOfWork(&block)
	//不断做Hash运算
	thisBlockHash, nonce := pow.Run()
	block.ThisBlockHash = thisBlockHash
	block.Nonce = nonce
	return &block 
}

//将区块的数据转换成字节流（序列化/编码）
func (block *Block)Serialize() []byte {
	//创建字节流
	var buffer bytes.Buffer
	//创建编码器
	encoder := gob.NewEncoder(&buffer)
	//编码
	err := encoder.Encode(&block)
	if err != nil{
		log.Panic("编码错误！", err)
		return nil
	}
	return buffer.Bytes()
}

//将区块的数据转换成字节流（反序列化/解码）
func Deserialize(bytesArry []byte) Block {
	//创建block
	var block Block
	//创建Reader流,创建解码器
	decoder := gob.NewDecoder(bytes.NewReader(bytesArry))
	err := decoder.Decode(&block)
	if err != nil{
		log.Panic("解码错误!", err)
	}
	return block
}
