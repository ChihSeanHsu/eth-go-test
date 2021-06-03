package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/eth-go-test/pkg/models"
	"github.com/ethereum/go-ethereum"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"math/big"
	"sync"
	"time"
)

const (
	EtheUrlHttp      = "https://data-seed-prebsc-2-s3.binance.org:8545/"
	EtheUrlWebSocket = "wss://bsc-ws-node.nariox.org:443"
)

var DB *models.DB

type retryCount struct {
	BlockNum *big.Int
	Retry    int
}

func subNewBlock(sub ethereum.Subscription, headers <-chan *ethTypes.Header, jobs chan *big.Int) {
	for {
		select {
		case err := <-sub.Err():
			fmt.Errorf("%s", err)
		case header := <-headers:
			fmt.Println("sub insert ", header.Number)
			jobs <- header.Number
		}
	}
}

func forLoopCheckNewBlock(client *ethclient.Client, jobs chan *big.Int, end big.Int) {
	ctx := context.Background()
	var i *big.Int
	for {
		header, err := client.HeaderByNumber(ctx, nil)
		if err != nil {
			fmt.Println("get header err: ", err)
		}
		if header.Number.Cmp(&end) != 0 {
			one := big.NewInt(1)
			for i = &end; header.Number.Cmp(i) != 0; i.Add(i, one) {
				jobs <- big.NewInt(i.Int64())
				fmt.Println("loop insert ", i)
			}
			end = *header.Number
		}
		time.Sleep(2 * time.Second)
	}
}

func worker(id int, jobs <-chan *big.Int, rescues chan retryCount) {
	ethClient, err := ethclient.Dial(EtheUrlHttp)
	if err != nil {
		log.Fatal("worker", id, "err", err)
	}
	ctx := context.Background()
	for blockNum := range jobs {
		if block, err := ethClient.BlockByNumber(ctx, blockNum); err == nil {
			// TODO: store to db
			DB.CreateBlock(ctx, ethClient, block)
			fmt.Printf("worker%d block: %s, %s\n", id, block.Hash(), block.Number())
		} else if err.Error() == "not found" {
			rescues <- retryCount{blockNum, 1}
		} else {
			fmt.Printf("worker%d, block: %v, %v\n", id, blockNum, err)

		}
	}
}

func rescuer(id int, rescues <-chan retryCount, rescuesQueue chan retryCount) {
	ethClient, err := ethclient.Dial(EtheUrlHttp)
	if err != nil {
		log.Fatal("rescuer", id, "err", err)
	}
	ctx := context.Background()
	for retry := range rescues {
		if block, err := ethClient.BlockByNumber(ctx, retry.BlockNum); err == nil {
			// TODO: store to db
			time.Sleep(time.Duration(retry.Retry) * time.Second)
			DB.CreateBlock(ctx, ethClient, block)
			fmt.Printf("rescuer%d block: %s, %s\n", id, block.Hash(), block.Number())
		} else if err.Error() == "not found" && retry.Retry <= 3 {
			retry.Retry++
			rescuesQueue <- retry
		} else {
			fmt.Printf("rescuer%d, %v, %v\n", id, retry.BlockNum, err)
		}
	}
}

func main() {
	var j, end *big.Int
	var wg sync.WaitGroup
	startFrom := flag.Int64("start", -1, "start from which block index")
	workerNum := flag.Int("worker", 5, "How many workers")
	wsMode := flag.Bool("ws", false, "How many workers")
	flag.Parse()

	dsn := "host=db user=postgres password=example dbname=db port=5432 sslmode=disable TimeZone=Asia/Taipei"
	DB = models.InitDB(dsn, 20)
	jobs := make(chan *big.Int)
	rescues := make(chan retryCount)
	connStr := EtheUrlHttp
	if *wsMode {
		connStr = EtheUrlWebSocket
	}
	client, err := ethclient.Dial(connStr)
	if err != nil {
		log.Fatal("connect err: ", err)
	}

	ctx := context.Background()
	header, err := client.HeaderByNumber(ctx, nil)
	if err != nil {
		log.Fatal("get header err: ", err)
	}
	end = header.Number

	if *wsMode {
		headers := make(chan *ethTypes.Header)
		sub, err := client.SubscribeNewHead(context.Background(), headers)
		if err != nil {
			log.Fatal("sub err: ", err)
		}
		go subNewBlock(sub, headers, jobs)
	} else {
		go forLoopCheckNewBlock(client, jobs, *end)
	}

	for w := 1; w <= *workerNum; w++ {
		go worker(w, jobs, rescues)
		wg.Add(1)
	}
	for w := 1; w <= *workerNum; w++ {
		go rescuer(w, rescues, rescues)
		wg.Add(1)
	}
	if *startFrom == -1 {
		*startFrom = header.Number.Int64()
	}
	one := big.NewInt(1)
	for j = big.NewInt(*startFrom); j.Cmp(end) != 0; j.Add(j, one) {
		fmt.Println("outer insert ", j)
		jobs <- big.NewInt(j.Int64())
	}

	wg.Wait()
}
