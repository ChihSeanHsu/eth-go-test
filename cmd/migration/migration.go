package main

import (
	"github.com/eth-go-test/pkg/models"
)

func main() {
	dsn := "host=db user=postgres password=example dbname=db port=5432 sslmode=disable TimeZone=Asia/Taipei"
	db := models.InitDB(dsn, 1, 1)
	db.AutoMigrate(&models.Block{}, &models.Transaction{}, &models.Log{})
}
