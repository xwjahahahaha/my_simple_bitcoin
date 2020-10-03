package main

import (
	"bytes"
	"sort"
)

//为了UTXOs的排序实现sort.Interface接口的三个方法
type UTXOs []TxOutput

//实现计算长度
func (utxos UTXOs) Len() int {
	return len(utxos)
}

//实现交换函数
func (utxos UTXOs) Swap(i, j int) {
	utxos[i], utxos[j] = utxos[j], utxos[i]
}

//实现排序规则
func (utxos UTXOs) Less(i, j int) bool {
	return utxos[i].Value < utxos[j].Value		//按余额从小到大排序
}


// 查找账户的所有UTXO(未消费的支出)
func (bc *BlockChain)FindUTXOs(PubKeyHash []byte) []TxOutput {
	//var utxos []TxOutput
	utxos := UTXOs{} //放置所有未花费的交易
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
				//比较两个byte数组,找到公钥的hash相同的，验证身份
				if bytes.Equal(output.PubKeyHash, PubKeyHash) {
					// 是的话就加入到UTXO中
					utxos = append(utxos, output)
				}
			}
			//TODO 加密解密的过程要完善
			//4. 遍历交易中的输入组，如果同地址的话就找其对应交易对应编号的输出进行标记
			//如果当前交易是挖矿交易就不遍历标记
			if tx.isCoinbaseTx() {
				continue
			}
			for _, input := range tx.Vin{
				//比较输入的公钥的hash与查找人的公钥hash
				//先对输入中的公钥进行hash运算
				if bytes.Equal(HashPubKey(input.PubKey), PubKeyHash){
					//加入到标记数组
					//生成交易对应的输出编号数组
					SpentOutputMap[string(input.TxHash)] = append(SpentOutputMap[string(input.TxHash)], input.OutputIndex)
				}
			}
		}

		//结束条件
		if len(block.PrevBlockHash) == 0{
			//fmt.Println("查询余额结束！")
			break
		}
	}
	return utxos
}


// 查找转账账户所有最适合的UTXO，用来转账
// 返回两个值，一个map集合， 一个最合理的金额
func (bc *BlockChain)FindSuitableUTXO(senderPubKeyHash []byte, amount float64) (map[string][]int64, float64) {
	//1. 先查询账户的所有UTXO  (因为每个output都有其编号和所属交易的字段=>可以定位到.所以不用再次遍历区块、交易)
	utxos := UTXOs{}
	utxos = bc.FindUTXOs(senderPubKeyHash)
	var bestUtxosMap = make(map[string][]int64)  	//存储找到的集合
	var cumulativePrice float64						//累计的钱
	//2 .排序找到合适的部分金额
	//排序 //TODO 合适的UTXO拼接方式可以优化
	sort.Sort(utxos)
	//遍历累加金额
	for _, output := range utxos{
		if cumulativePrice < amount {
			//3. 添加到map集合中
			//如果当前output的钱加上累计的钱还比需要的钱小，就加入
			bestUtxosMap[string(output.TxHash)] = append(bestUtxosMap[string(output.TxHash)], output.Index)
			cumulativePrice += output.Value		//累计计算目前的总额钱
			if cumulativePrice >= amount{
				return bestUtxosMap, cumulativePrice 		//判断钱够了，就直接返回
			}
		}
	}
	//4. 返回map和金额
	return bestUtxosMap, cumulativePrice
}

