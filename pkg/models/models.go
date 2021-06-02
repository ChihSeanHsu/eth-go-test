package models

type Block struct {
	BlockHash    string
	BlockNum     uint
	BlockTime    uint
	ParentHash   string
	Transactions []string
}
