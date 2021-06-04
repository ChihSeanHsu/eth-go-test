package main

import (
	"github.com/eth-go-test/pkg/models"
)

func main() {
	db := models.InitDB(1, 1)
	db.AutoMigrate(&models.Block{}, &models.Transaction{}, &models.Log{})
}
