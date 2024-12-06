package main

import (
	"fmt"
	"sync"
	"time"

	"emr/mvx-winter-challenge/1206-query-transactions/pkg/csv"
	"emr/mvx-winter-challenge/1206-query-transactions/pkg/transactions"
)

const (
	apiEndpoint     = "https://testnet-api.multiversx.com"
	resultsDir      = "results"
	maxRetries      = 10
	retryDelay      = 5 * time.Second
	fetcherPageSize = 200
)

var addresses = []string{
	"erd1dn8cdeljw9e5mkwnw8wnanwwl7h0g07kur04jva8eewdf8xggc2qcykpxj",
	"erd14azdv5gzrr5udqlls4vs0n69ewpj7wkpxvctr44aaj2t7pgrstzqclpq4s",
	"erd156ym2y90ts0hpk0nvvd2xm0u7mtc6y3p2mqxagf4qmlj9672fgxq48ew6j",
	"erd1zhwcnrp3q3l5vn7k65wqmqk84lqqlyj4l6jmke2qx8e0pxfg5glsut63ct",
	"erd10x6mj2a02cd33tx2nt9l5qfz750ened8rh5uuwsdfy4x8r6qjt3sp34der",
	"erd13fdhz46q4m7tvy43sg6mpcwmpwxu7jseq2hgd353wed4g89ppy2stn26c3",
	"erd1sng37zsvu2u6s6sfst5zd3e0l6nt26pv6u9amkzx0snzkwe7cl0q0zdhmq",
	"erd10fknux996mp7nauw8wklu2rs7sskmtwdj85aqw4qd7w2wrksf68qnvfldr",
	"erd1ylv7su0lcce8xw97aujqguqhfeckd3mgh7zhkrpytu282e26kc8qszzdpf",
}

func main() {
	txFetcher := transactions.NewFetcher(apiEndpoint, maxRetries, retryDelay, fetcherPageSize)

	var wg sync.WaitGroup
	for _, addr := range addresses {
		wg.Add(1)
		go func(address string) {
			defer wg.Done()
			if err := processAddress(txFetcher, address); err != nil {
				fmt.Printf("Error processing address %s: %v\n", address, err)
			}
		}(addr)
	}
	wg.Wait()
}

func processAddress(fetcher *transactions.Fetcher, address string) error {
	writer, err := csv.OpenFile(resultsDir, address)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer writer.CloseFile()

	total := 0
	for {
		transactions, err := fetcher.FetchTransactions(address, total)
		if err != nil {
			return fmt.Errorf("failed to fetch transactions: %w", err)
		}

		if err := writer.WriteTransactionsPage(transactions); err != nil {
			return fmt.Errorf("failed to write transactions page: %w", err)
		}

		if len(transactions) < fetcherPageSize {
			break
		}
		total += len(transactions)
	}

	fmt.Printf("%d transactions written to %s\n", total, writer.File.Name())

	return nil
}
