package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math/big"
)

//定义一个工作量证明的结构ProofOfWork

type ProofOfWork struct {
	//a. block
	block *Block
	//b. 目标值 一个非常大的数
	target *big.Int
}

// NewProofOfWork 提供创建POW的函数
func NewProofOfWork(block *Block) *ProofOfWork {
	pow := ProofOfWork{
		block: block,
	}
	//我们指定难度值，现在是一个string类型
	targetStr := "0001000000000000000000000000000000000000000000000000000000000000"
	//引入的辅助变量，目的是将上面的难度值转成big.int
	tmpInt := big.Int{}
	//将难度值赋值给big.int, 指定16进制的格式
	tmpInt.SetString(targetStr, 16)

	pow.target = &tmpInt
	return &pow
}

// Run 3.提供计算不断计算哈希的函数
func (pow *ProofOfWork) Run() ([]byte, uint64) {
	var nonce uint64
	block := pow.block
	var hash [32]byte

	fmt.Println("开始挖矿！！！")
	for {
		//1.拼装数据(区块的数据，还有不断变化的随机数)
		tmp := [][]byte{
			Uint64ToByte(block.Version),
			block.PrevHash,
			block.MerkelRoot,
			Uint64ToByte(block.TimeStamp),
			Uint64ToByte(block.Difficulty),
			Uint64ToByte(nonce),
			//只对区块头做哈希，区块体通过MerkelRoot产生影响
			//block.Data,
		}
		//将二维的切片数组链接起来，返回一个一维的切片
		blockInfo := bytes.Join(tmp, []byte{})

		//2.做哈希运算
		hash = sha256.Sum256(blockInfo)
		//3.与pow中的target进行比较
		tmpInt := big.Int{}
		//将我们得到的hash数组转化为一个big.Int
		tmpInt.SetBytes(hash[:])

		//比较当前的哈希与目标哈希,如果当前哈希值小于目标哈希值,就说明找到了
		//	-1 if x <  y
		//	 0 if x == y
		//	+1 if x >  y
		if tmpInt.Cmp(pow.target) == -1 {
			//a.找到了,推出返回
			fmt.Printf("挖矿成功！hash: %x, nonce: %d\n", hash, nonce)
			//break
			return hash[:], nonce
		} else {
			//b.没找到，继续找，随机数加一
			nonce++
		}
	}
	//return hash[:], nonce
	//return []byte("HelloWorld"), 10
}
