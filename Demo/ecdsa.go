package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"log"
	"math/big"
)

func main()  {
	//创建椭圆曲线
	curve := elliptic.P256()
	//生成私钥
	privateKey, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil{
		log.Panic()
	}
	//得到公钥
	publicKey := privateKey.PublicKey

	//数据取hash
	data := "hahahah"
	hash := sha256.Sum256([]byte(data))

	//签名
	r, s, err := ecdsa.Sign(rand.Reader, privateKey, hash[:])
	if err != nil{
		log.Panic()
	}

	//把r和s拼接到一起传输
	signature := append(r.Bytes(), s.Bytes()...)
	//传输....

	//拆开签名，得到r和s
	r1 := big.Int{}
	s1 := big.Int{}
	//r是前半部分，s是后半部分
	r1.SetBytes(signature[:len(signature)/2])
	s1.SetBytes(signature[len(signature)/2:])

	fmt.Printf("r：%v\n", r1)
	fmt.Printf("s：%v\n", s1)
	fmt.Printf("签名：%v\n", signature)


	//校验
	//校验需要三个部分：数据hash， 签名， 公钥
	ret := ecdsa.Verify(&publicKey, hash[:], r, s)
	fmt.Printf("校验结果%v\n", ret)



}