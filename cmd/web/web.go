package main

import (
	"context"
	"fmt"
	"github.com/eth-go-test/pkg/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

var (
	DB *models.DB
)

func getBlocks(c *gin.Context) {
	var limit, defaultLimit uint64
	var result []gin.H

	defaultLimit = 10
	limitString := c.Query("limit")
	if limitString != "" {
		l, err := strconv.ParseUint(limitString, 10, 64)
		if err != nil {
			fmt.Printf("Parse limit err: %s", err)
			c.String(http.StatusBadRequest, fmt.Sprintf("error: %s", err))
		}
		limit = l
	} else {
		limit = defaultLimit
	}
	ctx := context.Background()
	// TODO: error handling
	blocks := DB.GetBlocks(ctx, int(limit))

	for _, block := range blocks {
		result = append(result, gin.H{
			"block_num":   block.BlockNum,
			"block_hash":  block.BlockHash,
			"block_time":  block.BlockTime,
			"parent_hash": block.ParentHash,
		})
	}
	c.JSON(200, gin.H{
		"blocks": result,
	})
}

func getBlockByID(c *gin.Context) {
	var txs []string
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	ctx := context.Background()
	// TODO: error handling
	block := DB.GetBlockByID(ctx, id)
	for _, tx := range block.Transactions {
		txs = append(txs, tx.TransactionHash)
	}
	c.JSON(200, gin.H{
		"block_num":    block.BlockNum,
		"block_hash":   block.BlockHash,
		"block_time":   block.BlockTime,
		"parent_hash":  block.ParentHash,
		"transactions": txs,
	})
}

func getTxByHash(c *gin.Context) {
	var logs []gin.H
	hash := c.Param("txHash")
	ctx := context.Background()
	// TODO: add error handling
	tx := DB.GetTxByHash(ctx, hash)
	for _, log := range tx.Logs {
		logs = append(logs, gin.H{
			"index": log.Index,
			"data":  log.Data,
		})
	}

	c.JSON(200, gin.H{
		"tx_hash": tx.TransactionHash,
		"to":      tx.To,
		"from":    tx.From,
		"nonce":   tx.Nonce,
		"data":    tx.Data,
		"value":   tx.Value,
		"logs":    logs,
	})
}

func main() {
	dsn := "host=db user=postgres password=example dbname=db port=5432 sslmode=disable TimeZone=Asia/Taipei"
	DB = models.InitDB(dsn, 20)
	r := gin.Default()
	r.GET("/blocks", getBlocks)
	r.GET("/blocks/:id", getBlockByID)
	r.GET("/transaction/:txHash", getTxByHash)
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
