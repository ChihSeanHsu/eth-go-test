package models

type Block struct {
	BlockHash    string `gorm:"primaryKey;type:varchar(256)"`
	BlockNum     uint64 `gorm:"type:bigint"`
	BlockTime    uint64
	ParentHash   string        `gorm:"type:varchar(256)"`
	Transactions []Transaction `gorm:"foreignKey:BlockHash;references:BlockHash"`
}
