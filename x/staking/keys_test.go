package staking

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStakeKey(t *testing.T) {
	var stakeID uint64 = 0x1A2B3C4D
	key := stakeKey(stakeID)
	assert.Equal(t, key, []byte{0x00, 0x0, 0x0, 0x0, 0x0, 0x1A, 0x2B, 0x3C, 0x4D})
}

func TestArgumentKey(t *testing.T) {
	var argumentID uint64 = 0xD4C3B2A11A2B3C4D
	key := argumentKey(argumentID)
	assert.Equal(t, key, []byte{0x01, 0xD4, 0xC3, 0xB2, 0xA1, 0x1A, 0x2B, 0x3C, 0x4D})
}
