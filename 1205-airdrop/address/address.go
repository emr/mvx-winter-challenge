package address

import (
	"crypto/rand"
	"fmt"

	"github.com/multiversx/mx-sdk-go/data"
)

// GenerateRandom generates the specified number of random valid MultiversX addresses
func GenerateRandom(count int) ([]string, error) {
	addresses := make([]string, count)
	for i := 0; i < count; i++ {
		addrBytes := make([]byte, 32)
		if _, err := rand.Read(addrBytes); err != nil {
			return nil, fmt.Errorf("failed to generate random bytes: %w", err)
		}

		addr := data.NewAddressFromBytes(addrBytes)
		addrBech32, err := addr.AddressAsBech32String()
		if err != nil {
			return nil, fmt.Errorf("failed to get address as bech32 string: %w", err)
		}

		addresses[i] = addrBech32
	}

	return addresses, nil
}
