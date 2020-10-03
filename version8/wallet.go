package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"github.com/btcsuite/btcutil/base58"
	"log"
)

//这里的钱包是一种结构，每一个钱包都保存了公钥私钥对

type Wallet struct {
	//私钥
	Private *ecdsa.PrivateKey
	//公钥，这里的PublicKey不存储原始的公钥，而是x和y的拼接字符串，在校验端重新拆分
	PublicKey []byte
}

//创建钱包
func NewWallet() *Wallet {
	// 创建椭圆曲线
	curve := elliptic.P256()
	//生成私钥
	privateKey, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil{
		log.Panic()
	}
	//得到原始公钥
	publicKeyOrg := privateKey.PublicKey
	//拼接字符串
	publicKey := append(publicKeyOrg.X.Bytes(), publicKeyOrg.Y.Bytes()...)

	//返回钱包
	return &Wallet{
		Private:   privateKey,
		PublicKey: publicKey,
	}
}

//根据公钥生成地址
func (w Wallet) NewAddress() string {
	//获取公钥
	pubKey := w.PublicKey

	//第一步: 对公钥做hash
	rip160HashValue := HashPubKey(pubKey)

	//第二步：第一步结果与版本号1bytes拼接，形成21Bytes的结果
	version := byte(00)
	Bytes21Data := append([]byte{version}, rip160HashValue...)

	//第三步：分支拷贝加密
	//两次hash256,截取最后4字节
	Bytes4Date := CheckSum(Bytes21Data)

	//第四步:拼接
	Bytes25Data := append(Bytes21Data, Bytes4Date...)

	//最后一步：base58
	address := base58.Encode(Bytes25Data)

	return address
}
