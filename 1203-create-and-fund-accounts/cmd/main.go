package main

import (
	"fmt"
	"log"
	"time"

	"emr/mvx-winter-challenge/create-and-fund-accounts/internal/browser"
	"emr/mvx-winter-challenge/create-and-fund-accounts/internal/generator"
)

func main() {
	// Generate accounts
	gen, err := generator.NewGenerator("accounts")
	if err != nil {
		log.Fatalf("Failed to create generator: %v", err)
	}

	accounts, err := gen.GenerateAccounts(3)
	if err != nil {
		log.Fatalf("Failed to generate accounts: %v", err)
	}

	// Print generated accounts
	for _, acc := range accounts {
		fmt.Printf("Generated account %s in shard %d: %s\n", acc.Id, acc.ShardID, acc.Address)
	}

	// Initialize faucet requester
	browser := browser.NewBrowser()
	defer browser.Close()

	// Request tokens for each account
	for _, acc := range accounts {
		fmt.Printf("Requesting tokens for account %s (%s)...\n", acc.Id, acc.Address)

		err := browser.RequestTokens(acc.PEMFile)
		if err != nil {
			log.Printf("Failed to request tokens for account %s: %v", acc.Id, err)
			continue
		}

		fmt.Printf("Successfully requested tokens for account %s\n", acc.Id)
		time.Sleep(2 * time.Second) // Add delay between requests
	}
}
