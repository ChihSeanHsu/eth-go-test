package models

import (
	"gorm.io/gorm"
)

type Transaction struct {
	gorm.Model
	ID              int64
	BlockHash       string `gorm:"type:varchar(256)"`
	TransactionHash string `gorm:"primaryKey;type:varchar(256)"`
	From            string `gorm:"type:varchar(256)"`
	To              string `gorm:"type:varchar(256)"`
	Nonce           uint64
	Data            []byte `gorm:"type:bytea"`
	Value           uint64 `gorm:"type:bigint"`
	Logs            []Log  `gorm:"foreignKey:TransactionHash;references:TransactionHash"`
}
