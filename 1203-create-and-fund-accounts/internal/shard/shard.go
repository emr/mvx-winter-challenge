package shard

import "math"

// ComputeShardID calculates the shard ID from a public key
func ComputeShardID(pubKey []byte, shardCount uint32) uint32 {
	// Use the last byte of the public key
	lastByte := pubKey[len(pubKey)-1]
	addr := uint32(lastByte)

	// Calculate number of bits needed to represent number of shards
	n := uint32(math.Ceil(math.Log2(float64(shardCount))))
	maskHigh := uint32((1 << n) - 1)
	maskLow := uint32((1 << (n - 1)) - 1)

	shard := addr & maskHigh
	if shard > 2 {
		shard = addr & maskLow
	}

	return shard
}
