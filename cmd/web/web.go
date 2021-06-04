package main

import (
	"context"
	"errors"
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
	var result gin.H
	var status int

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
	blocks, err := DB.GetBlocks(ctx, int(limit))
	if errors.Is(err, models.ErrNotFound) {
		status = http.StatusNotFound
		result = gin.H{
			"err": err,
		}

	} else if err != nil {
		status = http.StatusInternalServerError
		result = gin.H{
			"err": err,
		}
	} else {
		var array []gin.H
		for _, block := range blocks {
			array = append(array, gin.H{
				"block_num":   block.BlockNum,
				"block_hash":  block.BlockHash,
				"block_time":  block.BlockTime,
				"parent_hash": block.ParentHash,
			})
		}
		status = http.StatusOK
		result = gin.H{
			"block": array,
		}
	}
	c.JSON(status, result)
}

func getBlockByID(c *gin.Context) {
	var txs []string
	var result gin.H
	var status int

	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	ctx := context.Background()
	// TODO: error handling
	block, err := DB.GetBlockByID(ctx, id)
	if errors.Is(err, models.ErrNotFound) {
		result = gin.H{"err": err.Error()}
		status = http.StatusNotFound
	} else if err != nil {
		result = gin.H{"err": err.Error()}
		status = http.StatusInternalServerError
	} else {
		for _, tx := range block.Transactions {
			txs = append(txs, tx.TransactionHash)
		}
		result = gin.H{
			"block_num":    block.BlockNum,
			"block_hash":   block.BlockHash,
			"block_time":   block.BlockTime,
			"parent_hash":  block.ParentHash,
			"transactions": txs,
		}
		status = http.StatusOK
	}
	c.JSON(status, result)
}

func getTxByHash(c *gin.Context) {
	var logs []gin.H
	var result gin.H
	var status int

	hash := c.Param("txHash")
	ctx := context.Background()
	// TODO: add error handling
	tx, err := DB.GetTxByHash(ctx, hash)
	if errors.Is(err, models.ErrNotFound) {
		result = gin.H{"err": err.Error()}
		status = http.StatusNotFound
	} else if err != nil {
		result = gin.H{"err": err.Error()}
		status = http.StatusInternalServerError
	} else {
		for _, log := range tx.Logs {
			logs = append(logs, gin.H{
				"index": log.Index,
				"data":  log.Data,
			})
		}
		result = gin.H{
			"tx_hash": tx.TransactionHash,
			"to":      tx.To,
			"from":    tx.From,
			"nonce":   tx.Nonce,
			"data":    tx.Data,
			"value":   tx.Value,
			"logs":    logs,
		}
		status = http.StatusOK

	}
	c.JSON(status, result)
}

func main() {
	DB = models.InitDB(20, 1)
	r := gin.Default()
	r.GET("/blocks", getBlocks)
	r.GET("/blocks/:id", getBlockByID)
	r.GET("/transaction/:txHash", getTxByHash)
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
