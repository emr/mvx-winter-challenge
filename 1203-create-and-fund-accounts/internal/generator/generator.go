package generator

import (
	"fmt"
	"os"
	"path/filepath"

	"emr/mvx-winter-challenge/create-and-fund-accounts/internal/shard"

	"github.com/multiversx/mx-sdk-go/interactors"
)

type Account struct {
	Address    string
	Id         string
	Password   string
	Mnemonic   string
	ShardID    uint32
	PrivateKey []byte
	PEMFile    string
}

const shardCount = 3

type generator struct {
	outputDir string
}

func NewGenerator(dirName string) (*generator, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("error getting current directory: %w", err)
	}

	outputDir := filepath.Join(currentDir, dirName)

	// Create output directory if it doesn't exist
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return nil, fmt.Errorf("error creating output directory: %w", err)
	}
	return &generator{
		outputDir: outputDir,
	}, nil
}

func (g *generator) GenerateAccounts(accountsPerShard uint32) ([]Account, error) {
	accounts := make([]Account, 0, shardCount*accountsPerShard)

	for shardID := uint32(0); shardID < shardCount; shardID++ {
		for i := uint32(0); i < accountsPerShard; i++ {

			account, err := g.generateAccountInShard(shardID, i)
			if err != nil {
				return nil, fmt.Errorf("error generating account for shard %d: %w", shardID, err)
			}
			accounts = append(accounts, account)
		}
	}
	return accounts, nil
}

// generateAccountInShard keeps generating accounts until one in the desired shard is found
func (g *generator) generateAccountInShard(targetShard uint32, index uint32) (Account, error) {
	w := interactors.NewWallet()

	for {
		mnemonic, err := w.GenerateMnemonic()
		if err != nil {
			return Account{}, fmt.Errorf("error generating mnemonic: %w", err)
		}

		privateKey := w.GetPrivateKeyFromMnemonic(mnemonic, 0, 0)
		address, err := w.GetAddressFromPrivateKey(privateKey)
		if err != nil {
			return Account{}, fmt.Errorf("error getting address from private key: %w", err)
		}

		pubKey := address.AddressBytes()
		computedShard := shard.ComputeShardID(pubKey, shardCount)

		if computedShard != targetShard {
			// Continue generating if shard doesn't match
			continue
		}

		addressBech32, err := address.AddressAsBech32String()
		if err != nil {
			return Account{}, fmt.Errorf("error getting address as bech32 string: %w", err)
		}

		pemFilePath := fmt.Sprintf("%s/shard%d_acc%d.pem", g.outputDir, targetShard, index)
		err = w.SavePrivateKeyToPemFile(privateKey, pemFilePath)
		if err != nil {
			return Account{}, fmt.Errorf("error saving private key to pem file: %w", err)
		}

		id := fmt.Sprintf("Shard%d#Acc%d", targetShard, index)

		return Account{
			Address:    addressBech32,
			Id:         id,
			Password:   id,
			Mnemonic:   string(mnemonic),
			ShardID:    targetShard,
			PrivateKey: privateKey,
			PEMFile:    pemFilePath,
		}, nil
	}
}
