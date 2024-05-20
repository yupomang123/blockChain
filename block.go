package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/gob"
	"log"
	"time"
)

// Block 1.定义结构
type Block struct {
	// 1.版本号
	Version uint64
	// 2.前区块哈希
	PrevHash []byte
	// 3.Merkel根(一个哈希值)
	MerkelRoot []byte
	// 4.时间戳
	TimeStamp uint64
	// 5,难度值
	Difficulty uint64
	// 6.随机数，也就是挖矿要找的数据
	Nonce uint64

	// a.当前区块哈希
	Hash []byte
	// b.数据
	//Data []byte
	//真实的交易数组
	Transactions []*Transaction
}

// Uint64ToByte 实现一个辅助函数，功能是将uint64转成[]byte
func Uint64ToByte(num uint64) []byte {
	var buffer bytes.Buffer

	err := binary.Write(&buffer, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}
	return buffer.Bytes()
}

// NewBlock 2.创建区块
func NewBlock(txs []*Transaction, prevBlockHash []byte) *Block {
	block := Block{
		Version:    00,
		PrevHash:   prevBlockHash,
		MerkelRoot: []byte{},
		TimeStamp:  uint64(time.Now().Unix()),
		Difficulty: 0,
		Nonce:      0,
		Hash:       []byte{},
		//Data:       []byte(data),
		Transactions: txs,
	}

	block.MerkelRoot = block.MakeMerkelRoot()

	//block.SetHash()
	//创建一个pow对象
	pow := NewProofOfWork(&block)
	//查找随机数,不停的进行哈希运算
	hash, nonce := pow.Run()

	//根据挖矿结果对区块数据进行更新
	block.Hash = hash
	block.Nonce = nonce

	return &block
}

// 序列化
func (block *Block) Serialize() []byte {
	//编码的数据放到buffer
	var buffer bytes.Buffer

	//使用gob进行序列化(编码)得到字节流
	//1.定义一个编码器 使用编码器进行编码
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(&block)
	if err != nil {
		log.Panic("编码出错")
	}
	//fmt.Printf("编码后的小明: %v\n", buffer.Bytes())
	return buffer.Bytes()
}

// 反序列化
func DeSerialize(data []byte) Block {
	//使用gob进行反序列化(解码) 得到person结构
	//1.定义一个解码器 使用解码器进行解码
	decoder := gob.NewDecoder(bytes.NewReader(data))

	var block Block

	err := decoder.Decode(&block)
	if err != nil {
		log.Panic("解码出错")
	}
	return block
}

// SetHash 3.生成哈希
/*func (block *Block) SetHash() {
var blockInfo []byte
//1.拼装数据
/*blockInfo = append(blockInfo, Uint64ToByte(block.Version)...)
blockInfo = append(blockInfo, block.PrevHash...)
blockInfo = append(blockInfo, block.MerkelRoot...)
blockInfo = append(blockInfo, Uint64ToByte(block.TimeStamp)...)
blockInfo = append(blockInfo, Uint64ToByte(block.Difficulty)...)
blockInfo = append(blockInfo, Uint64ToByte(block.Nonce)...)
blockInfo = append(blockInfo, block.Data...)*/
/*tmp := [][]byte{
		Uint64ToByte(block.Version),
		block.PrevHash,
		block.MerkelRoot,
		Uint64ToByte(block.TimeStamp),
		Uint64ToByte(block.Difficulty),
		Uint64ToByte(block.Nonce),
		block.Data,
	}
	//将二维的切片数组链接起来，返回一个一维的切片
	blockInfo = bytes.Join(tmp, []byte{})
	//2.sha256
	hash := sha256.Sum256(blockInfo)
	block.Hash = hash[:]
}*/

// 模拟梅克尔根，只是对交易的数据做简单的拼接，而不做二叉树处理
func (block *Block) MakeMerkelRoot() []byte {

	var info []byte
	for _, tx := range block.Transactions {
		//将交易的哈希值拼接起来，再整体做哈希处理
		info = append(info, tx.TXID...)
	}
	hash := sha256.Sum256(info)
	return hash[:]
}
