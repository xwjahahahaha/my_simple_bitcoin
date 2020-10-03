package main

import "fmt"

//计算比特币总量2100万
func  main()  {
	bitTotal := 0.0; //总量
	blockInterval := 21.0 //单位是万
	currentReward := 50.0 //出块奖励

	for currentReward > 0 {
		//每个区间（4年）的量
		amount := blockInterval * currentReward;
		currentReward *= 0.5 //出块奖励每4年减半 乘法效率比除法高
		bitTotal += amount
	}
	fmt.Println(bitTotal , "万")

}