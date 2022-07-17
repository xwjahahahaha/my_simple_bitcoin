# 请查看最新版本Version9文件夹

# simpleBitCoin
简单版本的比特币系统
# 功能概要
* 区块的创建
* 交易的打包
* 用户余额的查询
* pow挖矿算法
* 钱包功能，增加账户、管理账户等
* 转账功能
* 公私钥签名以及验证
# 文件介绍
![](http://xwjpics.gumptlu.work/qiniu_picGo/20201004092202.png)

注意：
1. refresh重置文件不会重置wallet.dat即钱包存储文件
2. 如果更新了版本名称则需要修改refresh中对应内容
3. 项目所依赖的包文件不在此仓库中，需自行安装配置
4. **详细的开发细节文档可见:**

# 使用
![](http://xwjpics.gumptlu.work/qiniu_picGo/20201004090620.png)

## 1. printChain 输出整条区块链
simple:
![](http://xwjpics.gumptlu.work/qiniu_picGo/20201004091517.png)

## 2. getBalance ADDRESS 查询账户余额
参数： ADDRESS-账户地址
simple:
![](http://xwjpics.gumptlu.work/qiniu_picGo/20201004091724.png)

## 3. send FROM TO AMOUNT MINER DATA  由FROM转AMOUNT钱给TO，由MINER挖矿，同时写入DATA
参数: FROM-转出人 TO-转入人 AMOUNT-转账金额 MINER-挖矿人 DATA-铸币交易可以自添加的数据
simple:
![](http://xwjpics.gumptlu.work/qiniu_picGo/20201004092752.png)

## 4. getTransaction TXHASH 查询交易信息
参数 TXHASH-交易hash
simple:
![](http://xwjpics.gumptlu.work/qiniu_picGo/20201004092924.png)

## 5. newWallet 创建一个新的钱包（公私钥对）
simple:
![](http://xwjpics.gumptlu.work/qiniu_picGo/20201004093103.png)

## 6. listAddress 列举所有的钱包地址
simple:
![](http://xwjpics.gumptlu.work/qiniu_picGo/20201004093157.png)

# 待完善
* 每个区块都是只能打包一个交易就直接发布了，没有区块链网络体系去获取交易
* 分布式网络共识协议没有实现
* 梅克尔树root的计算，目前项目只是简单的拼接字节
* 签名机制待完善，项目使用的签名验证方式是P2PKH，还有P2SH、多重签名等可以完善
* 远程访问rpc调用，类似于geth的远程访问
* 客户端的构建
