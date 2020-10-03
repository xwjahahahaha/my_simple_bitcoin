package main

import (
	"itcast_Go/bolt"
	"log"
)

//区块链迭代器
type BlockChainIterator struct {
	db *bolt.DB
	moveHashPoint []byte
}

//创建迭代器
func (bc BlockChain)NewIterator() *BlockChainIterator{
	return &BlockChainIterator{
		db:            bc.db,
		moveHashPoint: bc.tailHash,
	}
}

//迭代函数
func (bci *BlockChainIterator)Next() *Block {
	//1. 返回当前区块
	var block Block
	bci.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BucketName))
		if bucket != nil{
			bytes := bucket.Get(bci.moveHashPoint)
			//反序列化
			block = Deserialize(bytes)
			//hash左移
			bci.moveHashPoint = block.PrevBlockHash
		}else{
			log.Panic("没有此Bucket！")
		}
		return nil
	})
	//2. 前移一个位置
	return &block
}
