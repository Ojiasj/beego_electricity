package main

import (
	"crypto/sha256"
	"time"
	"bytes"
	"encoding/binary"
	"log"
)

// 定义结构
type Block struct {
	//版本号
	Version uint64
	// 前区块哈希
	PrevHash []byte
	// Merkel根
	MerkelRoot []byte
	// 时间戳
	TimeStamp uint64
	// 难度值
	Difficulty uint64
	// 随机数
	Nonce uint64
	// 当前区块哈希
	Hash []byte
	// 数据
	Data []byte
}

// 实现辅助函数，将uint64转[]byte
func Uint64ToByte(num uint64) []byte {
	var buffer bytes.Buffer
	err := binary.Write(&buffer, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}

	return buffer.Bytes()
}

// 创建区块
func NewBlock(data string, prevBlockHash []byte) *Block {
	block := Block{
		Version:    00,
		PrevHash:   prevBlockHash,
		MerkelRoot: []byte{},
		TimeStamp:  uint64(time.Now().Unix()),
		Difficulty: 0,
		Nonce:      0,
		Hash:       []byte{}, // 后面再补回来
		Data:       []byte(data),
	}

	block.SetHash()
	return &block
}

// 生成哈希
func (block *Block) SetHash() {
	//var blockInfo []byte
	//// 1.拼装数据
	//blockInfo = append(blockInfo, Uint64ToByte(block.Version)...)
	//blockInfo = append(blockInfo, block.PrevHash...)
	//blockInfo = append(blockInfo, block.MerkelRoot...)
	//blockInfo = append(blockInfo, Uint64ToByte(block.TimeStamp)...)
	//blockInfo = append(blockInfo, Uint64ToByte(block.Difficulty)...)
	//blockInfo = append(blockInfo, Uint64ToByte(block.Nonce)...)
	//blockInfo = append(blockInfo, block.Data...)
	// 2.sha256

	tmp := [][]byte{
		Uint64ToByte(block.Version),
		block.PrevHash,
		block.MerkelRoot,
		Uint64ToByte(block.TimeStamp),
		Uint64ToByte(block.Difficulty),
		Uint64ToByte(block.Nonce),
		block.Data,
	}
	blockInfo := bytes.Join(tmp, []byte{})

	//将而为的切片叔祖联结起来

	hash := sha256.Sum256(blockInfo)
	block.Hash = hash[:]
}
