package main

import (
	"bytes"
	"crypto/elliptic"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"log"
)

// wallets保存了所有的钱包wallet结构，其中包含了钱包地址和其公私钥对的映射

type Wallets struct {
	WalletsMap map[string]*Wallet  //所有钱包地址和其公私钥对的映射
}

//打开钱包
//读取钱包文件，加载其中钱包结构
func OpenWallets(walletsName string) *Wallets {
	//创建钱包
	var ws Wallets
	//读取文件
	ret := ws.loadWalletsFile(walletsName)
	if !ret {
		//没创建就创建
		fmt.Println("还未创建钱包, 已创建...")
		ws.WalletsMap = make(map[string]*Wallet)
		ws.createWallet(walletsName)
	}
	//已创建就直接返回
	return &ws
}

//添加钱包内容（添加一个地址与公私钥的映射），并保存到文件中
func (ws *Wallets)createWallet(saveName string) string {
	// 先创建Wallet
	wallet := NewWallet()
	// 再生成地址
	address := wallet.NewAddress()
	//加入到钱包中
	ws.WalletsMap[address] = wallet
	//保存整个钱包
	ws.saveWalletsFile(saveName)
	//返回地址
	return address
}

//保存钱包
func (ws *Wallets)saveWalletsFile(saveName string) {
	var buffer bytes.Buffer
	//注册接口
	gob.Register(elliptic.P256())
	//创建编码器
	encoder := gob.NewEncoder(&buffer)
	//编码
	err := encoder.Encode(ws)
	if err != nil{
		log.Panic(err)
	}
	//保存到文件  第三个参数是权限
	err = ioutil.WriteFile(saveName, buffer.Bytes(), 0600)
	if err != nil{
		log.Panic(err)
	}
}

//读取钱包文件
func (ws *Wallets)loadWalletsFile(walletsName string) bool {
	//读文件
	content, err := ioutil.ReadFile(walletsName)
	if err != nil{
		//如果没有钱包，那么就创建
		return false
		log.Panic(err)
	}
	//解码
	gob.Register(elliptic.P256())	//注册接口
	var load_ws Wallets
	decoder := gob.NewDecoder(bytes.NewReader(content))
	err = decoder.Decode(&load_ws)
	if err != nil{
		log.Panic(err)
	}
	//赋值
	ws.WalletsMap = load_ws.WalletsMap
	return true
}

//获取当前钱包的所有地址
func (ws *Wallets) GetAllAddress() []string {
	var addresses []string
	for address := range ws.WalletsMap{
		addresses = append(addresses, address)
	}
	return addresses
}
