package address

import (
	"testing"

	"github.com/multiversx/mx-sdk-go/data"
	"github.com/stretchr/testify/assert"
)

func TestGenerateRandom(t *testing.T) {
	addresses, err := GenerateRandom(10000)
	if err != nil {
		t.Errorf("GenerateRandom() error = %v", err)
	}

	assert.Equal(t, 10000, len(addresses))

	for _, addr := range addresses {
		_, err := data.NewAddressFromBech32String(addr)
		assert.Nil(t, err, "Address should be valid bech32 format")
	}
}
