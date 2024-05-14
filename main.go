package main

func main() {
	//block := NewBlock("中本聪转我50枚比特币！", []byte{})

	bc := NewBlockChain("张三")
	cli := CLI{bc}
	cli.Run()
	/*bc.AddBlock("班长向班花转了50枚比特币！")
	bc.AddBlock("班长又向班花转了50枚比特币！")*/

}
