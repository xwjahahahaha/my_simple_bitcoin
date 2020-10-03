package main

import (
	"fmt"
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
func NewBlockChain(address string) *BlockChain {
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
			genesisBlock := GenesisBlock(address)				//创建创世区块，要传入挖矿人地址
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
func GenesisBlock(address string) *Block {
	//创建出块交易coinBase,第二个参数是出块可以自由发挥的内容
	coinbasetx := createCoinBaseTx(address, "我是创世块")
	// 把创世交易打包到第一个区块中
	return NewBlock([]*Transaction{coinbasetx}, []byte{})
}


//添加区块
func (blockChain *BlockChain) AddBlock(txs []*Transaction)  {
	//加入到区块链数组
	blockChain.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BucketName))
		if bucket != nil{
			//使用数据库的区块链tailHash字段保存了最后一个hash
			//创建区块
			newBlock := NewBlock(txs, blockChain.tailHash)
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


// 查找账户的所有UTXO(未消费的支出)
func (bc *BlockChain)FindUTXOs(address string) []TxOutput {
	var UTXO []TxOutput		//放置所有未花费的交易
	//创建标记Map，key为交易编号，value是要标记的输出编号数组
	var SpentOutputMap = make(map[string][]int64)	//标记已消费的输出的Map

 	// 遍历交易中的输入组（标记输入）
	// 1. 遍历区块
 	//创建迭代器
 	iterator := bc.NewIterator()
 	for {
 		block := iterator.Next() 	//返回当前区块并向前一步
		// 2.遍历交易
		for _, tx := range block.Transactions{
			// 3.遍历交易中的输出组(检查是否被标记，如果没被标记那么就是用过了)
			OUTPUT:
			for i, output := range tx.Vout{
				//判断是否被标记, 过滤掉消耗
				if SpentOutputMap[string(tx.TxHash)] != nil{
					//先判断当前交易是否存在与已消费的map记录中
					//为空则不需要检查编号标记，非空则检查编号是否存在
					for _, index := range SpentOutputMap[string(tx.TxHash)] {
						if index == int64(i){
							//当前交易的当前编号存在，表明已消费过，那么就不加入到UTXO中
							//跳转不在执行下面的步骤
							continue OUTPUT  	//注意continue跳转到外循环
						}
					}
				}

				//判断地址是否是当前查询人
				if string(output.ScriptPubKey) == address {
					// 是的话就加入到UTXO中
					UTXO = append(UTXO, output)
				}
			}
			//TODO 加密解密的过程要完善
			//4. 遍历交易中的输入组，如果同地址的话就找其对应交易对应编号的输出进行标记
			//如果当前交易是挖矿交易就不遍历标记
			if tx.isCoinbaseTx() {
				continue
			}
			for _, input := range tx.Vin{
				if string(input.ScriptSig) == address{
					//加入到标记数组
					//生成交易对应的输出编号数组
					SpentOutputMap[string(input.TxHash)] = append(SpentOutputMap[string(input.TxHash)], input.OutputIndex)
				}
			}


		}

		//结束条件
		if len(block.PrevBlockHash) == 0{
			fmt.Println("查询余额结束！")
			break
		}
	}

	return UTXO
}


// 查找转账账户所有最适合的UTXO，用来转账
// 返回两个值，一个map集合， 一个最合理的金额
func (bc *BlockChain)FindSuitableUTXO(address string, amount float64) (map[string][]int64, float64) {
	//项目中并没有优化计算最合适的将零钱拼装，而是简单的遍历逐步统计，满足要求了就转账
	//TODO 最合理的计算（优化）
	var bestUtxosMap = make(map[string][]int64)  	//存储找到的集合
	var SpentOutputMap = make(map[string][]int64)	//标记已消费的输出的Map
	var cumulativePrice float64						//累计的钱
	// 1. 遍历区块
	//创建迭代器
	iterator := bc.NewIterator()
	for {
		block := iterator.Next() 	//返回当前区块并向前一步
		// 2.遍历交易
		for _, tx := range block.Transactions{
			// 3.遍历交易中的输出组(检查是否被标记，如果没被标记那么就是用过了)
		OUTPUT:
			for i, output := range tx.Vout{
				//判断是否被标记, 过滤掉消耗
				if SpentOutputMap[string(tx.TxHash)] != nil{
					//先判断当前交易是否存在与已消费的map记录中
					//为空则不需要检查编号标记，非空则检查编号是否存在
					for _, index := range SpentOutputMap[string(tx.TxHash)] {
						if index == int64(i){
							//当前交易的当前编号存在，表明已消费过，那么就不加入到UTXO中
							//跳转不在执行下面的步骤
							continue OUTPUT  	//注意continue跳转到外循环
						}
					}
				}

				//判断地址是否是当前查询人
				if string(output.ScriptPubKey) == address {
					// 判断是否足够，不足就累加
					if cumulativePrice < amount {
						//如果当前output的钱加上累计的钱还比需要的钱小，就加入
						bestUtxosMap[string(tx.TxHash)] = append(bestUtxosMap[string(tx.TxHash)], int64(i))
						cumulativePrice += output.Value		//累计计算目前的总额钱
						if cumulativePrice >= amount{
							return bestUtxosMap, cumulativePrice 		//判断钱够了，就直接返回
						}
					}
				}
			}
			//TODO 加密解密的过程要完善
			//4. 遍历交易中的输入组，如果同地址的话就找其对应交易对应编号的输出进行标记
			//如果当前交易是挖矿交易就不遍历标记
			if tx.isCoinbaseTx() {
				continue
			}
			for _, input := range tx.Vin{
				if string(input.ScriptSig) == address{
					//加入到标记数组
					//生成交易对应的输出编号数组
					SpentOutputMap[string(input.TxHash)] = append(SpentOutputMap[string(input.TxHash)], input.OutputIndex)
				}
			}

		}
		//结束条件
		if len(block.PrevBlockHash) == 0{
			fmt.Println("余额不足！！")
			break
		}
	}

	return bestUtxosMap, cumulativePrice
}


