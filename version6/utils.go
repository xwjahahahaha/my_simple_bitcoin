package main

import (
	"bytes"
	"encoding/binary"
	"log"
)

//Uint64 转[]byte数组
func Uint64ToByte(num uint64) []byte {
	var buffer bytes.Buffer
	err := binary.Write(&buffer, binary.BigEndian, num)
	if err != nil{
		log.Panic(err)
	}
	return buffer.Bytes()
}
