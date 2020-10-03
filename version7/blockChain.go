package main

import (
	"itcast_Go/bolt"
	"log"
)
//4. 引入区块链
type BlockChain struct{
	//使用blot数据库代替数组
	db *bolt.DB	//数据库
	tailHash []byte	//最后入链区块的hash
}

const BlockChainDBName  = "blockChain.db"
const BucketName  = "BlockBucket"

//5. 定义一个区块链
func NewBlockChain(miner string) *BlockChain {
	var lastBlockHash []byte
	//第一个参数是名字，第二个参数是权限6代表允许读写
	db, err := bolt.Open(BlockChainDBName, 0600, nil)
	//defer db.Close() 先不关
	if err != nil{
		log.Panic("打开数据库失败！" , err)
		return nil
	}
	//操作数据库
	db.Update(func(tx *bolt.Tx) error {
		//2. 打开抽屉(没有就创建)
		var bucketName []byte = []byte(BucketName)
		bucket := tx.Bucket(bucketName)
		if bucket == nil{
			//没有就创建
			bucket, err = tx.CreateBucket(bucketName)
			if err != nil{
				log.Panic("打开桶", BucketName,"失败！" ,err)
				return err
			}
			//3. 写数据
			//操作抽屉中的数据，添加数据
			genesisBlock := GenesisBlock(miner)				//创建创世区块，要传入挖矿人地址
			lastBlockHash = genesisBlock.ThisBlockHash			//创世区块的hash就是最后一个hash
			bucket.Put([]byte("lastHashKey"), lastBlockHash)	//修改数据库中最后一个常量键值对（保存最后的hash）
			//区块的hash作为key， 区块的字节流作为value
			bucket.Put(lastBlockHash, genesisBlock.Serialize())
		}else {
			//已经有桶了就读取最后的常量key: lastHashKey的value值（保存的是最后一个块的hash）
			//lastHashKey -> 最后一个块的hash
			lastBlockHash = bucket.Get([]byte("lastHashKey"))
		}
		return nil
	})

	return &BlockChain{
		db:       db,
		tailHash: lastBlockHash,
	}
}

//创世区块
func GenesisBlock(miner string) *Block {
	//创建出块交易coinBase,第二个参数是出块可以自由发挥的内容
	coinbasetx := createCoinBaseTx(miner, "我是创世块")
	// 把创世交易打包到第一个区块中
	return NewBlock(miner, []*Transaction{coinbasetx}, []byte{})
}


//添加区块
func (blockChain *BlockChain) AddBlock(miner string, txs []*Transaction)  {
	//加入到区块链数组
	blockChain.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BucketName))
		if bucket != nil{
			//使用数据库的区块链tailHash字段保存了最后一个hash
			//创建区块
			newBlock := NewBlock(miner, txs, blockChain.tailHash)
			//区块的hash
			newBlockHash := newBlock.ThisBlockHash
			bucket.Put(newBlockHash, newBlock.Serialize())
			//把两边（chain和db）的标记最后一个hash的字段修改
			blockChain.tailHash = newBlockHash				//内存中的区块链blockchain的末尾字段也要改
			bucket.Put([]byte("lastHashKey"), newBlockHash)	//修改数据库中最后一个常量键值对（保存最后的hash）
		}else {
			log.Panic("没有此bucket！")
		}
		return nil
	})
}

