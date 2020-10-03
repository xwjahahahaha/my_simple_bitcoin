package main
import "time"


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

//3. 生成Hash
//func (block *Block) SetHash()  {
//	//二维切片放置每一个一维[]byte数组,以后需要加参数在此添加即可
//	tmp := [][]byte{
//		Uint64ToByte(block.Version),
//		block.PrevBlockHash,
//		block.MerkelRoot,
//		Uint64ToByte(block.TimeStamp),
//		Uint64ToByte(block.Difficulty),
//		Uint64ToByte(block.Nonce),
//		block.Data,
//	}
//	//将二维切片数组连接成一维的切片[]byte，二维数组元素中间的连接内容为空（第二个参数）
//	blockInfo := bytes.Join(tmp, []byte{})
//	//2.sha256
//	hash := sha256.Sum256(blockInfo)
//	block.ThisBlockHash = hash[:]
//}


