package models

import (
	"context"
	"errors"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
	"time"
)

type DB struct {
	*gorm.DB
}

func InitDB(pool int, retry int) *DB {
	dsn := os.Getenv("DB_CONN_STR")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil && retry <= 3 {
		log.Println(err)
		// waiting for return
		waitSec := 10 * retry
		retry++
		log.Println("wait for reconnect...")
		time.Sleep(time.Duration(waitSec) * time.Second)
		return InitDB(pool, retry)
	} else if err != nil {
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
			log.Printf("Get chainID err: %s", err)
		}

		receipt, err := ethClient.TransactionReceipt(ctx, tx.Hash())
		if err != nil {
			log.Printf("Get receipt err: %s", err)
		}
		for _, l := range receipt.Logs {
			logs = append(logs, Log{
				Index: l.Index,
				Data:  l.Data,
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

func (db *DB) GetBlocks(ctx context.Context, limit int) ([]Block, error) {
	var blocks []Block
	err := db.Order("block_num desc").Limit(limit).Find(&blocks).Error
	if len(blocks) == 0 {
		return blocks, ErrNotFound
	}
	return blocks, err
}

func (db *DB) GetBlockByID(ctx context.Context, id uint64) (Block, error) {
	var block Block
	err := db.Where("blocks.block_num = ?", id).Preload("Transactions").First(&block).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return block, ErrNotFound
	}
	return block, err
}

func (db *DB) GetTxByHash(ctx context.Context, hash string) (Transaction, error) {
	var tx Transaction
	err := db.Where("transactions.transaction_hash = ?", hash).Preload("Logs").
		First(&tx).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return tx, ErrNotFound
	}
	return tx, err
}
