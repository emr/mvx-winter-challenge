package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"math"
	"math/big"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/multiversx/mx-chain-crypto-go/signing"
	"github.com/multiversx/mx-chain-crypto-go/signing/ed25519"
	"github.com/multiversx/mx-sdk-go/blockchain"
	"github.com/multiversx/mx-sdk-go/blockchain/cryptoProvider"
	"github.com/multiversx/mx-sdk-go/builders"
	"github.com/multiversx/mx-sdk-go/core"
	"github.com/multiversx/mx-sdk-go/interactors"
)

const (
	gasPrice             = 1000000000
	gasLimit             = 60000000
	networkConfig        = "testnet"
	proxyURL             = "https://testnet-gateway.multiversx.com"
	tokenName            = "WinterToken"
	tokenTicker          = "WINTER"
	supply               = 100000000
	decimals             = 8
	esdtIssueAddr        = "erd1qqqqqqqqqqqqqqqpqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqzllls8a5w6u"
	esdtIssueCost        = "50000000000000000" // 0.05 EGLD
	tokenCountPerAccount = 3
)

type Account struct {
	Index          int
	PrivateKey     []byte
	Address        core.AddressHandler
	AddressBench32 string
}

func main() {
	accounts, err := loadAccounts("accounts")
	if err != nil {
		fmt.Printf("Failed to load accounts: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Loaded %d accounts\n", len(accounts))

	for _, acc := range accounts {
		txHashes, err := issueTokensForAccount(acc)
		if err != nil {
			fmt.Printf("Failed to issue tokens for account %s: %v\n", acc.AddressBench32, err)
			continue
		}
		fmt.Printf("Issued tokens for account %s\n", acc.AddressBench32)
		for _, txHash := range txHashes {
			fmt.Printf("  - https://testnet-explorer.multiversx.com/transactions/%s\n", txHash)
		}
	}
}

func loadAccounts(dir string) ([]Account, error) {
	var accounts []Account
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}

	for i, file := range files {
		if !strings.HasSuffix(file.Name(), ".pem") {
			continue
		}

		w := interactors.NewWallet()

		privateKey, err := w.LoadPrivateKeyFromPemFile(filepath.Join(dir, file.Name()))
		if err != nil {
			return nil, fmt.Errorf("failed to load private key: %w", err)
		}

		address, err := w.GetAddressFromPrivateKey(privateKey)
		if err != nil {
			return nil, fmt.Errorf("failed to get address from private key: %w", err)
		}

		addressBench32, err := address.AddressAsBech32String()
		if err != nil {
			return nil, fmt.Errorf("failed to get address as bech32 string: %w", err)
		}

		accounts = append(accounts, Account{
			Index:          i,
			PrivateKey:     privateKey,
			Address:        address,
			AddressBench32: addressBench32,
		})
	}

	return accounts, nil
}

func issueTokensForAccount(account Account) ([]string, error) {
	suite := ed25519.NewEd25519()
	keyGen := signing.NewKeyGenerator(suite)
	signer := cryptoProvider.NewSigner()
	proxyArgs := blockchain.ArgsProxy{
		ProxyURL:            proxyURL,
		Client:              nil,
		SameScState:         false,
		ShouldBeSynced:      false,
		FinalityCheck:       false,
		CacheExpirationTime: time.Minute,
		EntityType:          core.Proxy,
	}

	proxy, err := blockchain.NewProxy(proxyArgs)
	if err != nil {
		return nil, fmt.Errorf("failed to create proxy: %w", err)
	}

	networkConfigs, err := proxy.GetNetworkConfig(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to get network config: %w", err)
	}

	txBuilder, err := builders.NewTxBuilder(signer)
	if err != nil {
		return nil, fmt.Errorf("failed to create tx builder: %w", err)
	}

	txi, err := interactors.NewTransactionInteractor(proxy, txBuilder)
	if err != nil {
		return nil, fmt.Errorf("failed to create transaction interactor: %w", err)
	}

	holder, err := cryptoProvider.NewCryptoComponentsHolder(keyGen, account.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create crypto components holder: %w", err)
	}

	issueTx, _, err := proxy.GetDefaultTransactionArguments(context.Background(), account.Address, networkConfigs)
	if err != nil {
		return nil, fmt.Errorf("failed to get default transaction arguments: %w", err)
	}

	supplyWithDecimals := big.NewInt(int64(math.Pow10(decimals)))
	supplyWithDecimals.Mul(supplyWithDecimals, big.NewInt(supply))

	args := []string{
		"issue",
		hex.EncodeToString([]byte(tokenName)),
		hex.EncodeToString([]byte(tokenTicker)),
		fmt.Sprintf("%016x", supplyWithDecimals),
		fmt.Sprintf("%02x", decimals),
	}

	issueTx.Data = []byte(strings.Join(args, "@"))
	issueTx.Value = esdtIssueCost
	issueTx.Receiver = esdtIssueAddr
	issueTx.GasLimit = uint64(gasLimit)

	for i := 0; i < 3; i++ {
		tx := issueTx
		err = txi.ApplyUserSignature(holder, &tx)
		if err != nil {
			return nil, fmt.Errorf("failed to sign transaction: %w", err)
		}
		txi.AddTransaction(&tx)
		issueTx.Nonce++
	}

	txHashes, err := txi.SendTransactionsAsBunch(context.Background(), 3)
	if err != nil {
		return nil, fmt.Errorf("failed to send transaction: %w", err)
	}

	return txHashes, nil
}
