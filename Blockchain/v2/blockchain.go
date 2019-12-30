package main

// 引入区块连
type BlockChain struct {
	// 定义一个区块连数据
	blocks []*Block
}

// 定义一个区块连
func NewBlockChain() *BlockChain {
	genesisBlock := GenesisBlock()
	return &BlockChain{
		blocks: []*Block{genesisBlock},
	}
}

// 创世块
func GenesisBlock() *Block {
	return NewBlock("我是创世神", []byte{})
}

// 添加区块
func (bc *BlockChain) AddBlock(data string) {
	// 获取前区块
	lastBlock := bc.blocks[len(bc.blocks)-1]
	prevHash := lastBlock.Hash
	// 创建新的区块
	block := NewBlock(data, prevHash)
	// 添加区块到区块连中
	bc.blocks = append(bc.blocks, block)

}
