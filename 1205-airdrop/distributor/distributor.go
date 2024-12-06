package distributor

import (
	"context"
	"fmt"
	"math"
	"math/big"
	"time"

	"github.com/multiversx/mx-chain-crypto-go/signing"
	"github.com/multiversx/mx-chain-crypto-go/signing/ed25519"
	"github.com/multiversx/mx-sdk-go/blockchain"
	"github.com/multiversx/mx-sdk-go/blockchain/cryptoProvider"
	"github.com/multiversx/mx-sdk-go/builders"
	"github.com/multiversx/mx-sdk-go/core"
	"github.com/multiversx/mx-sdk-go/interactors"

	"emr/mvx-winter-challenge/airdrop/account"
	"emr/mvx-winter-challenge/airdrop/api"
)

// Config holds the distributor configuration
type Config struct {
	ProxyURL    string
	ApiURL      string
	TokenName   string
	TokenAmount int64
	GasLimit    uint64
	GasPrice    uint64
	BatchSize   int
}

// Distributor handles token distribution operations
type Distributor struct {
	config    Config
	apiClient *api.Client
}

// DistributionResult holds the operation result information
type DistributionResult struct {
	AccountAddress string
	Tokens         []api.TokenData
	TxHashes       map[string][]string
	Error          error
}

// New creates a new Distributor instance
func NewDistributor(config Config) *Distributor {
	return &Distributor{
		config:    config,
		apiClient: api.NewClient(config.ApiURL),
	}
}

// Distribute distributes tokens from the account to the addresses provided
func (d *Distributor) Distribute(ctx context.Context, acc account.Account, addresses []string) DistributionResult {
	result := DistributionResult{
		AccountAddress: acc.AddressBench32,
		TxHashes:       make(map[string][]string),
	}

	tokens, err := d.apiClient.GetAccountTokens(acc.AddressBench32, d.config.TokenName)
	if err != nil {
		result.Error = fmt.Errorf("failed to get tokens: %w", err)
		return result
	}

	result.Tokens = tokens
	for _, token := range tokens {
		if token.Owner != acc.AddressBench32 {
			continue
		}

		requiredAmount := big.NewInt(d.config.TokenAmount)
		requiredAmount.Mul(requiredAmount, big.NewInt(int64(math.Pow10(token.Decimals))))
		requiredAmount.Mul(requiredAmount, big.NewInt(int64(len(addresses))))

		if token.Balance.Cmp(requiredAmount) < 0 {
			continue
		}

		txHashes, err := d.sendTransactions(ctx, acc, addresses, token)
		if err != nil {
			result.Error = fmt.Errorf("failed to send transactions: %w", err)
			return result
		}
		result.TxHashes[token.Identifier] = txHashes
	}

	return result
}

func (d *Distributor) sendTransactions(ctx context.Context, acc account.Account, addresses []string, token api.TokenData) ([]string, error) {
	suite := ed25519.NewEd25519()
	keyGen := signing.NewKeyGenerator(suite)
	signer := cryptoProvider.NewSigner()

	proxyArgs := blockchain.ArgsProxy{
		ProxyURL:            d.config.ProxyURL,
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

	networkConfigs, err := proxy.GetNetworkConfig(ctx)
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

	holder, err := cryptoProvider.NewCryptoComponentsHolder(keyGen, acc.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create crypto components holder: %w", err)
	}

	transferTx, _, err := proxy.GetDefaultTransactionArguments(ctx, acc.Address, networkConfigs)
	if err != nil {
		return nil, fmt.Errorf("failed to get default transaction arguments: %w", err)
	}
	transferTx.GasLimit = d.config.GasLimit
	transferTx.Value = "0"

	for _, addr := range addresses {
		tx := transferTx
		tx.Receiver = addr

		data, err := d.buildTransferData(token)
		if err != nil {
			return nil, fmt.Errorf("failed to build transfer data: %w", err)
		}
		tx.Data = data

		err = txi.ApplyUserSignature(holder, &tx)
		if err != nil {
			return nil, fmt.Errorf("failed to sign transaction: %w", err)
		}
		txi.AddTransaction(&tx)
		transferTx.Nonce++
	}

	txHashes, err := txi.SendTransactionsAsBunch(ctx, d.config.BatchSize)
	if err != nil {
		return nil, fmt.Errorf("failed to send transaction: %w", err)
	}

	return txHashes, nil
}

func (d *Distributor) buildTransferData(token api.TokenData) ([]byte, error) {
	txData := builders.NewTxDataBuilder()
	txData.Function("ESDTTransfer")
	txData.ArgBytes([]byte(token.Identifier))

	tokenAmountWithDecimals := big.NewInt(int64(math.Pow10(token.Decimals)))
	tokenAmountWithDecimals.Mul(tokenAmountWithDecimals, big.NewInt(d.config.TokenAmount))
	txData.ArgInt64(tokenAmountWithDecimals.Int64())

	return txData.ToDataBytes()
}
