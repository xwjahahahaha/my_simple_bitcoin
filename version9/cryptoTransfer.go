package main

import (
	"bytes"
	"crypto/sha256"
	"github.com/btcsuite/btcutil/base58"
	"golang.org/x/crypto/ripemd160"
	"log"
)

//4字节校验码生成函数
//进行两次SHA256，再截断最后4字节
func CheckSum(data []byte) []byte {
	//两次hash256
	hash1 := sha256.Sum256(data[:])
	hash2 := sha256.Sum256(hash1[:])
	//提取前4个字节
	Bytes4Date := hash2[:4]
	return Bytes4Date
}

//公钥 => 公钥hash
//公钥做hash   先SHA256再ripemd160
func HashPubKey(pubKey []byte) []byte {
	// 第一步：RIPEMD160(sha256(pk))
	//做一次sha256
	pk_hash_1 := sha256.Sum256(pubKey)
	//ripemd160加密函数加密
	//创建编码器
	rip160hasher := ripemd160.New()
	//写入公钥hash
	_, err := rip160hasher.Write(pk_hash_1[:])
	if err != nil{
		log.Panic(err)
	}
	//返回Hash结果
	rip160HashValue := rip160hasher.Sum(nil)
	return rip160HashValue
}

//地址 => 公钥的Hash
func adsToPubKeyHash(address string) []byte {
	// 1.base58函数的解码
	bytes25Data := base58.Decode(address)
	// 2.去除尾部添加的4byte校验码和首部添加的1byte版本号
	addressHash := bytes25Data[1: len(bytes25Data)-4]
	return addressHash
}

//校验地址
//校验流程思路:  （先反向走在正向走岔路回来）
func checkAddress(address string) bool {
	//1. 根据地址反推出25byte的数据。截断后4字节得到21字节的数据
	bytes25Data := base58.Decode(address)
	if len(bytes25Data) < 4 {		//地址长度不够直接返回
		return false
	}
	this4Bytes := bytes25Data[len(bytes25Data)-4 :]		//取最后四个
	//2. 将这21byte数据进行两次SHA256，再截取4字节的校验码
	org4Bytes := CheckSum(bytes25Data[: len(bytes25Data)-4])
	//3. 通过校验码与原本的25字节数据后四位比对
	return bytes.Equal(this4Bytes, org4Bytes)
}