package main

import (
	"fmt"
	)


func main() {
	bc := NewBlockChain()
	bc.AddBlock("第二个比特比")
	bc.AddBlock("第三个比特比")
	for i, block := range bc.blocks {
		fmt.Printf("=====区块高度%d======\n", i)
		fmt.Printf("前区块哈希值:%x\n", block.PrevHash)
		fmt.Printf("当前区块哈希值:%x\n", block.Hash)
		fmt.Printf("区块数据%s\n", block.Data)
	}
}
