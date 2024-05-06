package main

import (
	"fmt"
	"os"
)

type CLI struct {
	bc *BlockChain
}

const Usage = `
	addBlock --data DATA	"add data to blockchain"
	printChain				"print all blockchain data"
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
	default:
		fmt.Printf("无效的命令，请检查！\n")
		fmt.Printf(Usage)
	}
}
