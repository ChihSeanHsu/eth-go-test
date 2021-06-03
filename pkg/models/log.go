package models

import (
	"gorm.io/gorm"
)

type Log struct {
	gorm.Model
	TransactionHash string `gorm:"index:transaction_hash__idx"`
	Index           uint   `gorm:"index:transaction_hash__idx"`
	Data            []byte `gorm:"type:bytea"`
}
