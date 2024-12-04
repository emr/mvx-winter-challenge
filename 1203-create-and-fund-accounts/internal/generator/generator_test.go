package generator

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateAccounts(t *testing.T) {
	gen, err := NewGenerator("test")
	require.NoError(t, err)

	accounts, err := gen.GenerateAccounts(3)
	require.NoError(t, err)
	assert.Len(t, accounts, 9)

	// Verify we have exactly 3 accounts per shard
	shardCounts := make(map[uint32]int)
	for _, acc := range accounts {
		shardCounts[acc.ShardID]++
	}

	for shardID := uint32(0); shardID < 3; shardID++ {
		assert.Equal(t, 3, shardCounts[shardID], "Expected exactly 3 accounts in shard %d", shardID)
	}
}
