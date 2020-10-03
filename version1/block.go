package main
import "crypto/sha256"


//1.定义结构
type Block struct{
	//1.前一个区块的hash
	PrevBlockHash []byte
	//2.本区块的hash 
	ThisBlockHash []byte
	//3.数据
	Data []byte
}

//2.创建区块
func NewBlock(data string, prevBlockHash []byte)  *Block{
	block := Block{
		PrevBlockHash : prevBlockHash,
		ThisBlockHash : []byte{}, 
		Data : []byte(data),
	}
	block.SetHash()
	return &block 
}

//3. 生成Hash
func (block *Block) SetHash()  {
	//1.拼装数据
	blockInfo := append(block.PrevBlockHash, block.Data...)
	//2.sha256
	hash := sha256.Sum256(blockInfo)
	block.ThisBlockHash = hash[:]
}