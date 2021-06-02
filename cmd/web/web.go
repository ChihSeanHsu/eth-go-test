package main

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
	"math/big"
	"net/http"
	"strconv"
)
import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/ethclient"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"log"
)

var EthClient *ethclient.Client

const (
	EtheUrl = "https://data-seed-prebsc-2-s3.binance.org:8545/"
)


func addBlock(block *ethTypes.Block, result *[]gin.H)  {
	*result = append(*result, gin.H{
		"block_num": block.Number(),
		"block_hash": block.Hash(),
		"block_time": block.Time(),
		"parent_hash": block.ParentHash(),
	})
}

func getBlocks(c *gin.Context) {
	var limit, i, defaultLimit uint32
	var result []gin.H

	defaultLimit = 10
	limitString := c.Query("limit")
	if limitString != "" {
		l, err := strconv.ParseUint(limitString, 10, 32)
		if err != nil {
			fmt.Printf("Parse limit err: %s", err)
			c.String(http.StatusBadRequest, fmt.Sprintf("error: %s", err))
		}
		limit = uint32(l)
	} else {
		limit = defaultLimit
	}
	ctx := context.Background()

	header, err := EthClient.HeaderByNumber(ctx, nil)
	if err != nil {
		fmt.Printf("Get header err: %s", err)
		c.String(http.StatusInternalServerError, fmt.Sprintf("error: %s", err))
	}

	block, err := EthClient.BlockByNumber(ctx, header.Number)
	if err != nil {
		fmt.Printf("Get block err: %s", err)
		c.String(http.StatusInternalServerError, fmt.Sprintf("error: %s", err))
	}

	addBlock(block, &result)
	for i = 1; i < limit; i++ {
		block, err := EthClient.BlockByHash(ctx, block.ParentHash())
		if err != nil {
			fmt.Printf("Get block err: %s", err)
			c.String(http.StatusInternalServerError, fmt.Sprintf("error: %s", err))
		}
		addBlock(block, &result)
	}
	c.JSON(200, gin.H{
		"blocks": result,
	})
}

func getBlockByID(c *gin.Context) {
	var transactions []common.Hash
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	ctx := context.Background()
	block, err := EthClient.BlockByNumber(ctx, big.NewInt(id))
	if err != nil {
		fmt.Printf("Get block err: %s", err)
		c.String(http.StatusInternalServerError, fmt.Sprintf("error: %s", err))
	}
	for _, tx := range block.Transactions() {
		transactions = append(transactions, tx.Hash())
	}
	c.JSON(200, gin.H{
		"block_num": block.Number(),
		"block_hash": block.Hash(),
		"block_time": block.Time(),
		"parent_hash": block.ParentHash(),
		"transactions": transactions,
	})
}

func getTxByHash(c *gin.Context) {
	var from common.Address
	hash := c.Param("txHash")
	ctx := context.Background()
	// TODO: enhance error handling, and pending
	tx, _, err := EthClient.TransactionByHash(ctx, common.HexToHash(hash))
	if err != nil {
		fmt.Printf("Get block err: %s", err)
		c.String(http.StatusInternalServerError, fmt.Sprintf("error: %s", err))
	}
	// TODO: enhance error handling
	chainID, err := EthClient.NetworkID(ctx)
	if err != nil {
		fmt.Printf("Get chainID err: %s", err)
		c.String(http.StatusInternalServerError, fmt.Sprintf("error: %s", err))
	}
	// TODO: enhance error handling
	if msg, err := tx.AsMessage(ethTypes.NewEIP155Signer(chainID), tx.FeeCap()); err == nil {
		from = msg.From()
	}
	// TODO: enhance error handling
	receipt, err := EthClient.TransactionReceipt(ctx, tx.Hash())
	if err != nil {
		fmt.Printf("Get receipt err: %s", err)
		c.String(http.StatusInternalServerError, fmt.Sprintf("error: %s", err))
	}
	c.JSON(200, gin.H{
		"tx_hash": tx.Hash(),
		"to": tx.To(),
		"from": from,
		"nonce": tx.Nonce(),
		"data": tx.Data(),
		"value": tx.Value(),
		"logs": receipt.Logs,
	})
}

func initEthClient() {
	client, err := ethclient.Dial(EtheUrl)
	if err != nil {
		log.Fatal("connection fail", err)
	}
	EthClient = client
}

func main() {
	initEthClient()
	r := gin.Default()
	r.GET("/blocks", getBlocks)
	r.GET("/blocks/:id", getBlockByID)
	r.GET("/transaction/:txHash", getTxByHash)
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
