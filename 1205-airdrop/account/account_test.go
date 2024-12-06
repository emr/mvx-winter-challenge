package account

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoadAccounts(t *testing.T) {
	accounts, err := LoadAccounts(".")
	require.NoError(t, err)
	require.Len(t, accounts, 2)

	// First account
	require.NotEmpty(t, accounts[0].PrivateKey)
	require.NotNil(t, accounts[0].Address)
	require.Equal(t, accounts[0].AddressBench32, "erd17k8mese63jmj7ztfzv34dltudzgz7dksaz8pg0saj5crlppfu4uqe5x4ey")

	// Second account
	require.NotEmpty(t, accounts[1].PrivateKey)
	require.NotNil(t, accounts[1].Address)
	require.Equal(t, accounts[1].AddressBench32, "erd120m8zm9tw7q3h9eqxkyfe6j2ypv9lx08zpy3cw85ecnyjgs2gqxqumnedl")
}
