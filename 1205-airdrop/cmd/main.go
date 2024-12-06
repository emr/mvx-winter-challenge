package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"

	"emr/mvx-winter-challenge/airdrop/account"
	"emr/mvx-winter-challenge/airdrop/address"
	"emr/mvx-winter-challenge/airdrop/distributor"
)

const (
	numHolders      = 1000
	tokenAmount     = 10000
	proxyURL        = "https://testnet-gateway.multiversx.com"
	apiURL          = "https://testnet-api.multiversx.com"
	tokenNameFilter = "WinterToken"
	gasLimit        = 500000
	gasPrice        = 1000000000
	batchSize       = 10000
	accountsDir     = "accounts"
)

func main() {
	accounts, err := account.LoadAccounts(accountsDir)
	if err != nil {
		fmt.Printf("Failed to load accounts: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Loaded %d accounts\n", len(accounts))

	config := distributor.Config{
		ProxyURL:    proxyURL,
		ApiURL:      apiURL,
		TokenName:   tokenNameFilter,
		TokenAmount: tokenAmount,
		GasLimit:    gasLimit,
		GasPrice:    gasPrice,
		BatchSize:   batchSize,
	}

	dist := distributor.NewDistributor(config)

	var wg sync.WaitGroup
	for _, acc := range accounts {
		ctx := context.Background()
		wg.Add(1)
		go func(acc account.Account) {
			defer wg.Done()
			addresses, err := address.GenerateRandom(numHolders)
			if err != nil {
				fmt.Printf("Error generating addresses: %v\n", err)
				return
			}
			result := dist.Distribute(ctx, acc, addresses)
			if result.Error != nil {
				fmt.Printf("Error processing account %s: %v\n", acc.AddressBench32, result.Error)
				return
			}
			printResult(acc, addresses, result)
		}(acc)
	}
	wg.Wait()
}

func printResult(acc account.Account, addresses []string, result distributor.DistributionResult) {
	tokenIds := make([]string, len(result.Tokens))

	for i, t := range result.Tokens {
		tokenIds[i] = t.Identifier
	}

	if len(tokenIds) == 0 {
		fmt.Printf("Account %s doesn't own any tokens\n", acc.AddressBench32)
		return
	}

	totalTxs := 0
	for _, hashes := range result.TxHashes {
		totalTxs += len(hashes)
	}

	if len(tokenIds) == 1 {
		fmt.Printf("Account %s distributed the token %v to %d addresses in %d transactions\n",
			acc.AddressBench32,
			tokenIds[0],
			len(addresses),
			totalTxs)
		return
	}

	fmt.Printf("Account %s distributed %d tokens (%v) to %d addresses in %d transactions\n",
		acc.AddressBench32,
		len(tokenIds),
		strings.Join(tokenIds, ", "),
		len(addresses),
		totalTxs)
}
