package models

import (
	"context"
	"fmt"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"log"
	"time"
)

type DB struct {
	*gorm.DB
}

func InitDB(dsn string, pool int) *DB {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(pool)
	sqlDB.SetConnMaxIdleTime(time.Hour)

	return &DB{db}
}

func (db *DB) CreateBlock(ctx context.Context, ethClient *ethclient.Client, block *ethTypes.Block) {
	var transactions []Transaction
	for _, tx := range block.Transactions() {
		var logs []Log
		chainID, err := ethClient.NetworkID(ctx)
		if err != nil {
			fmt.Printf("Get chainID err: %s", err)
		}

		receipt, err := ethClient.TransactionReceipt(ctx, tx.Hash())
		if err != nil {
			fmt.Printf("Get receipt err: %s", err)
		}
		for _, log := range receipt.Logs {
			logs = append(logs, Log{
				Index: log.Index,
				Data:  log.Data,
			})
		}
		txModel := Transaction{
			TransactionHash: tx.Hash().Hex(),
			Nonce:           tx.Nonce(),
			Data:            tx.Data(),
			Value:           tx.Value().Uint64(),
			Logs:            logs,
		}
		if msg, err := tx.AsMessage(ethTypes.NewEIP155Signer(chainID)); err == nil {
			txModel.From = msg.From().Hex()
		}
		if tx.To() != nil {
			txModel.To = tx.To().Hex()
		}
		transactions = append(transactions, txModel)
	}
	db.Create(&Block{
		BlockHash:    block.Hash().Hex(),
		BlockNum:     block.Number().Uint64(),
		BlockTime:    block.Time(),
		ParentHash:   block.ParentHash().Hex(),
		Transactions: transactions,
	})
}

func (db *DB) GetBlocks(ctx context.Context, limit int) []Block {
	var blocks []Block
	db.Order("block_num desc").Limit(limit).Find(&blocks)
	return blocks
}

func (db *DB) GetBlockByID(ctx context.Context, id uint64) Block {
	var block Block
	db.Preload("Transactions.TransactionHash").Preload(clause.Associations).
		First(&block, "blocks.block_num = ?", id)
	return block
}

func (db *DB) GetTxByHash(ctx context.Context, hash string) Transaction {
	var tx Transaction
	db.Preload("Logs.Index Logs.Data").Preload(clause.Associations).
		First(&tx, "transactions.transaction_hash = ?", hash)
	return tx
}
