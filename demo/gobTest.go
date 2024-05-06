package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
)

type Person struct {
	Name string
	Age  uint
}

func main() {
	//定义一个结构Person
	var xiaoMing Person
	xiaoMing.Name = "小明"
	xiaoMing.Age = 20

	//编码的数据放到buffer
	var buffer bytes.Buffer

	//使用gob进行序列化(编码)得到字节流
	//1.定义一个编码器 使用编码器进行编码
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(&xiaoMing)
	if err != nil {
		log.Panic("编码出错，小明不知去向")
	}
	fmt.Printf("编码后的小明: %v\n", buffer.Bytes())

	//使用gob进行反序列化(解码) 得到person结构
	//1.定义一个解码器 使用解码器进行解码
	decoder := gob.NewDecoder(bytes.NewReader(buffer.Bytes()))

	var Daming Person

	err = decoder.Decode(&Daming)
	if err != nil {
		log.Panic("解码出错")
	}

	fmt.Printf("解码后的小明: %v\n", &Daming)
}
