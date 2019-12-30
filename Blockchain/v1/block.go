package main

import "crypto/sha256"

// 定义结构
type Block struct {
	// 前区块哈希
	PrevHash []byte
	// 当前区块哈希
	Hash []byte
	// 数据
	Data []byte
}

// 创建区块
func NewBlock(data string, prevBlockHash []byte) *Block {
	block := Block{
		PrevHash: prevBlockHash,
		Hash:     []byte{}, // 后面再补回来
		Data:     []byte(data),
	}

	block.SetHash()
	return &block
}

// 生成哈希
func (block *Block) SetHash() {
	// 1.拼装数据
	blockinfo := append(block.PrevHash, block.Data...)
	// 2.sha256
	hash := sha256.Sum256(blockinfo)
	block.Hash = hash[:]
}
