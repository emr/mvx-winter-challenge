package account

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/multiversx/mx-sdk-go/core"
	"github.com/multiversx/mx-sdk-go/interactors"
)

// Account represents a MultiversX account with its credentials
type Account struct {
	PrivateKey     []byte
	Address        core.AddressHandler
	AddressBench32 string
}

// LoadAccounts loads all PEM accounts from the specified directory
func LoadAccounts(dir string) ([]Account, error) {
	var accounts []Account
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}

	w := interactors.NewWallet()

	for _, file := range files {
		if !strings.HasSuffix(file.Name(), ".pem") {
			continue
		}

		privateKey, err := w.LoadPrivateKeyFromPemFile(filepath.Join(dir, file.Name()))
		if err != nil {
			return nil, fmt.Errorf("failed to load private key from %s: %w", file.Name(), err)
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
			PrivateKey:     privateKey,
			Address:        address,
			AddressBench32: addressBench32,
		})
	}

	return accounts, nil
}
