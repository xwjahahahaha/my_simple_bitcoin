package main

//4. 引入区块链
type BlockChain struct{
	blocks []*Block
}

//5. 定义一个区块链
func NewBlockChain() *BlockChain {
	genesisBlock := GenesisBlock()
	return  &BlockChain{
		blocks : []*Block{genesisBlock},  //创世块的加入
	}

}

//创世区块
func GenesisBlock() *Block {
	return NewBlock("我是创世块", []byte{0x0})
}


//6. 添加区块
func (blockChain *BlockChain) AddBlock(data string)  {
	//获取链中上一区块的hash->区块链数组中的最后一个块
	lastBlock := blockChain.blocks[len(blockChain.blocks) - 1]
	//创建区块
	newBlock := NewBlock(data, lastBlock.ThisBlockHash)
	//加入到区块链数组
	blockChain.blocks = append(blockChain.blocks, newBlock)	
}