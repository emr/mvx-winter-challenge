package shard

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestComputeShardID(t *testing.T) {
	tests := []struct {
		pubKey []byte
		want   uint32
	}{
		{[]byte{0x00}, 0},
		{[]byte{0x01}, 1},
		{[]byte{0x02}, 2},
	}

	for _, tt := range tests {
		got := ComputeShardID(tt.pubKey, 3)
		assert.Equal(t, tt.want, got, "ComputeShardID(%v) = %v, want %v", tt.pubKey, got, tt.want)
	}
}
