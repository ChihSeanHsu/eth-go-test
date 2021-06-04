package models

type Log struct {
	TransactionHash string `gorm:"primaryKey"`
	Index           uint   `gorm:"primaryKey;type:numeric"`
	Data            []byte `gorm:"type:bytea"`
}
