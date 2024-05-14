package main

import (
	"blockChain/bolt"
	"fmt"
	"log"
)

// BlockChain 4.引入区块链
// blockchain结构重写 使用数据库代替数组
type BlockChain struct {
	//定义一个区块链数组
	//blocks []*Block
	db *bolt.DB

	tail []byte //存储最后一个区块的哈希
}

const blockChainDb = "blockChain.db"
const blockBucket = "blockBucket"

// NewBlockChain 5.定义一个区块链
func NewBlockChain(address string) *BlockChain {

	/*return &BlockChain{
		blocks: []*Block{genesisBlock},
	}*/
	//最后一个区块的哈希，从数据库中读出来的
	var lastHash []byte

	//1.打开数据库
	db, err := bolt.Open(blockChainDb, 0600, nil)
	//defer db.Close()
	if err != nil {
		log.Panic("打开数据库失败！")
	}
	//操作数据库(改写)
	db.Update(func(tx *bolt.Tx) error {
		//2.找到抽屉bucket(如果没有就创建)
		bucket := tx.Bucket([]byte(blockBucket))
		if bucket == nil {
			//没有抽屉，我们需要创建
			bucket, err = tx.CreateBucket([]byte(blockBucket))
			if err != nil {
				log.Panic("创建blockBucket失败")
			}
			//创建一个创世块,并作为第一个区块添加到区块链中
			genesisBlock := GenesisBlock(address)
			//写数据 hash作为key，block字节流作为value
			bucket.Put(genesisBlock.Hash, genesisBlock.Serialize())
			bucket.Put([]byte("LastHashKey"), genesisBlock.Hash)
			lastHash = genesisBlock.Hash
		} else {
			lastHash = bucket.Get([]byte("LastHashKey"))
		}

		return nil
	})
	return &BlockChain{db, lastHash}
}

// GenesisBlock 定义一个创世块
func GenesisBlock(address string) *Block {
	coinbase := NewCoinbaseTX(address, "Go一期创世块,老牛逼了！")
	return NewBlock([]*Transaction{coinbase}, []byte{})
}

// AddBlock 添加区块
func (bc *BlockChain) AddBlock(txs []*Transaction) {
	//完成数据添加
	db := bc.db         //区块链数据库
	lastHash := bc.tail //最后一个区块的哈希

	db.Update(func(tx *bolt.Tx) error {

		//完成数据添加
		bucket := tx.Bucket([]byte(blockBucket))
		if bucket == nil {
			log.Panic("bucket 不应该为空，请检查！")
		}
		block := NewBlock(txs, lastHash)

		//hash作为key, block的字节流作为value,尚未发现
		bucket.Put(block.Hash, block.Serialize())
		bucket.Put([]byte("LastHashKey"), block.Hash)

		//c.跟新一下内存中的区块链，指的是把最后的小尾巴tail更新一下
		bc.tail = block.Hash

		return nil
	})

}

// FindUTXOs 找到指定地址的所有的utxo
func (bc *BlockChain) FindUTXOs(address string) []TXOutput {
	// UTXO 未消费输出
	var UTXO []TXOutput

	txs := bc.FindUTXOTransactions(address)

	for _, tx := range txs {
		for _, output := range tx.TXOutputs {
			if address == output.PubKeyHash {
				UTXO = append(UTXO, output)
			}
		}
	}
	return UTXO
}

func (bc *BlockChain) FindNeedUTXOs(from string, amount float64) (map[string][]uint64, float64) {
	//找到的合理的utxos集合
	utxos := make(map[string][]uint64)
	//找到的utxos里面包含的钱的总数
	var calc float64

	txs := bc.FindUTXOTransactions(from)

	for _, tx := range txs {
		for i, output := range tx.TXOutputs {
			if from == output.PubKeyHash {
				//我们需要实现的逻辑就在这里，找到自己需要的最少的utxo
				//比较一下是否满足转账需求
				//a.如果满足的话，直接返回 utxos, calc
				//b.不满足继续统计
				if calc < amount {
					//1.把utxo加进来
					utxos[string(tx.TXID)] = append(utxos[string(tx.TXID)], uint64(i))
					//2.统计一下当前utxo的总额
					calc += output.Value

					//加完之后满足条件了，return
					if calc >= amount {
						fmt.Printf("找到了满足的金额: %f\n", calc)
						return utxos, calc
					} else {
						fmt.Printf("不满足转账金额，当前总额: %f,目标金额:%f\n", calc, amount)
					}
				}
			}
		}
	}
	return utxos, calc
}

func (bc *BlockChain) FindUTXOTransactions(address string) []*Transaction {
	// UTXO 未消费输出
	//var UTXO []TXOutput
	var txs []*Transaction //存储所有包含utxo交易集合
	//我们定义一个map来保存消费过的output,key是这个output的交易id，value是这个交易中索引的数组
	spentOutputs := make(map[string][]int64)

	//创建迭代器
	it := bc.NewIterator()

	for {
		//1.遍历区块
		block := it.Next()

		//2.遍历交易
		for _, tx := range block.Transactions {
			//fmt.Printf("current txid : %x\n", tx.TXID)

			//打标签
		OUTPUT:
			//3.遍历output，找到和自己相关的utxo(在添加output前检查一下是否已经消耗过)
			for i, output := range tx.TXOutputs {
				//fmt.Printf("current index : %d\n", i)
				//在这里做一个过滤，将所有消耗过的outputs和当前所有即将添加的output对比一下
				//如果相同，则跳过，否则添加
				//如果当前的交易id存在于我们已经表示的map，那么说明这个交易里面有消耗过的output
				//map[2222] = []int64{0} map[3333] = []int64{0,1}
				if spentOutputs[string(tx.TXID)] != nil {
					for _, j := range spentOutputs[string(tx.TXID)] {
						//[]int64{0,1} j : 0,1
						if int64(i) == j {
							//当前准备添加的output已经消耗过了，不要再加了
							continue OUTPUT //跳出到标签
						}
					}
				}

				//这个output和我们目标的地址相同，满足条件，加到返回UTXO数组中
				if output.PubKeyHash == address {
					//UTXO = append(UTXO, output)
					//重点！！！返回所有包含我的utxo的交易的集合
					txs = append(txs, tx)
				}
			}
			//如果当前交易时挖矿交易，那么不做遍历，直接跳过
			if !tx.IsCoinbase() {
				//4. 遍历input，找到自己花费过的utxo的集合(把自己消耗过的表示出来)
				for _, input := range tx.TXInputs {
					//判断一下当前这个input和目标(李四)是否一致，如果相同，说明这个李四消耗过
					if input.Sig == address {
						//indexArray := spentOutputs[string(input.TXid)]
						//indexArray = append(indexArray, input.Index)
						spentOutputs[string(input.TXid)] = append(spentOutputs[string(input.TXid)], input.Index)
					}
				}
			} else {
				//fmt.Println("这是挖矿交易，不做遍历！")
			}
		}
		if len(block.PrevHash) == 0 {
			break
			fmt.Printf("区块链遍历完成退出！")
		}
	}

	return txs
}
