package main

import (
	"fmt"
	"itcast_Go/bolt"
	"log"
)

func main()  {
	//1. 打开数据库
	//第一个参数是名字，第二个参数是权限6代表允许读写
	db, err := bolt.Open("test.db", 0600, nil)
	defer db.Close()
	if err != nil{
		log.Panic("打开数据库失败！" , err)
	}
	//操作数据库
	db.Update(func(tx *bolt.Tx) error {
		//2. 打开抽屉(没有就创建)
		var bucketName []byte = []byte("b1")
		bucket := tx.Bucket(bucketName)
		if bucket == nil{
			//没有就创建
			bucket, err = tx.CreateBucket(bucketName)
			if err != nil{
				log.Panic(err)
			}
		}
		//操作抽屉中的数据，添加数据
		//3. 写数据
		bucket.Put([]byte("1111"), []byte("hello"))
		bucket.Put([]byte("2222"), []byte("world"))
		return nil
	})

	//4. 读数据
	db.View(func(tx *bolt.Tx) error {
		//找到抽屉
		bucket := tx.Bucket([]byte("b1"))
		if bucket != nil{
			//如果存在就读取
			v1 := bucket.Get([]byte("1111"))
			v2 := bucket.Get([]byte("2222"))
			//输出
			fmt.Printf("'1111'-> %s\n", v1)
			fmt.Printf("'2222'-> %s\n", v2)
		}
		return nil
	})
}
