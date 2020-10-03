package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math/big"
)


//1. 创建pow数据结构
type ProofOfWork struct {
	//a. block
	block *Block
	//b. 挖矿难度值
	target big.Int
}

//2. 提供创建pow的函数
func NewProofOfWork(block *Block) *ProofOfWork {
	//传的块写入pow中
	pow := ProofOfWork{
		block : block,
	}
	//加入手动写的难度值str
	targetStr := "0001000000000000000000000000000000000000000000000000000000000000"
	//转换str到big.int
	tmpInt := big.Int{}
	tmpInt.SetString(targetStr, 16)
	pow.target = tmpInt
	return &pow
}


//3. 提供不断计算hash的函数
//返回值是当前区块的hash和目标随机数
func (pow *ProofOfWork)Run() ([]byte, uint64) {
	var nonce uint64
	var hash [32]byte
	fmt.Printf("黄金矿工:[%s]开始挖矿...\n", pow.block.Miner)
	for {
		//1. 拼接区块 "头" 信息，包括nonce
		//!!!!注意是区块头而不是整个区块，因为使用了梅克尔树的树根包含了所有交易生成的root
		//二维切片放置每一个一维[]byte数组,以后需要加参数在此添加即可
		tmp := [][]byte{
			Uint64ToByte(pow.block.Version),
			pow.block.PrevBlockHash,
			pow.block.MerkelRoot,	//区块体交易通过root影响区块头
			Uint64ToByte(pow.block.TimeStamp),
			Uint64ToByte(pow.block.Difficulty),
			Uint64ToByte(nonce), //使用自己的随机数
		}
		//将二维切片数组连接成一维的切片[]byte，二维数组元素中间的连接内容为空（第二个参数）
		blockInfo := bytes.Join(tmp, []byte{})
		//2. 求hash
		hash = sha256.Sum256(blockInfo)
		//3. 比较难度值target
		//把[32]byte 数组类型转换成big.Int
		tmpInt := big.Int{}
		tmpInt.SetBytes(hash[:])
		//big.int间比较
		//   -1 if x <  y
		//    0 if x == y
		//   +1 if x >  y
		if tmpInt.Cmp(&pow.target) == -1{
			//a. 小于，达标返回
			fmt.Printf("新块的nonce：%d 新区块hash：%x\n挖矿成功！\n", nonce, hash)
			return hash[:], nonce
		}else {
			//b. 大于，修改nonce继续计算
			nonce ++
		}
	}
}

