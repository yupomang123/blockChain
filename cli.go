package main

import (
	"fmt"
	"os"
	"strconv"
)

type CLI struct {
	bc *BlockChain
}

// "add data to blockchain"
// "print all blockchain data"
const Usage = `
	addBlock --data DATA	"添加区块"		
	printChain				"打印所有区块链"
	getBalance --address ADDRESS "获取指定地址的余额"
	send FROM TO AMOUNT MINER DATA "由FROM转AMOUNT给TO,由MINER挖矿，同时写入DATA"
`

//接收参数的动作，我们放到一个函数中

func (cli *CLI) Run() {
	//./block printChain
	//./block addBlock --data "HelloWorld"
	//1.得到所有的命令
	args := os.Args
	if len(args) < 2 {
		fmt.Printf(Usage)
		return
	}

	//2.分析命令
	cmd := args[1]
	fmt.Println("输入命令", cmd)
	switch cmd {
	case "addBlock":
		//3.执行相应的动作
		fmt.Printf("添加区块\n")

		//确保命令有效
		if len(args) == 4 && args[2] == "--data" {
			//a.获取数据
			data := args[3]
			//b.使用bc添加区块AddBlock
			cli.AddBlock(data)
		} else {
			fmt.Printf("添加区块参数使用不当，请检查")
			fmt.Printf(Usage)
		}
	case "printChain":
		fmt.Printf("打印区块\n")
		cli.PrintBlockChain()
	case "getBalance":
		fmt.Printf("获取余额\n")
		if len(args) == 4 && args[2] == "--address" {
			address := args[3]
			cli.GetBalance(address)
		}
	case "send":
		fmt.Printf("转账开始...\n")
		if len(args) != 7 {
			fmt.Printf("参数个数错误，请检查！")
			fmt.Printf(Usage)
			return
		}
		// .exe send FROM TO AMOUNT MINER DATA "由FROM转AMOUNT给TO,由MINER挖矿，同时写入DATA"
		from := args[2]
		to := args[3]
		amount, _ := strconv.ParseFloat(args[4], 64)
		miner := args[5]
		data := args[6]
		cli.Send(from, to, amount, miner, data)
	default:
		fmt.Printf("无效的命令，请检查！\n")
		fmt.Printf(Usage)
	}
}
