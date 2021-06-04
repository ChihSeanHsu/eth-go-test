package models

type Log struct {
	TransactionHash string `gorm:"primaryKey"`
	Index           uint   `gorm:"primaryKey"`
	Data            []byte `gorm:"type:bytea"`
}
