package staking

import (
"testing"

"github.com/magiconair/properties/assert"
)

func TestStakeKey(t *testing.T) {
	var stakeID uint64 = 0x1A2B3C4D
	key := stakeKey(stakeID)
	assert.Equal(t, key, []byte{0x00, 0x4D, 0x3C, 0x2B, 0x1A, 0x0, 0x0, 0x0, 0x0})
}

func TestArgumentKey(t *testing.T) {
	var argumentID uint64 = 0xD4C3B2A11A2B3C4D
	key := argumentKey(argumentID)
	assert.Equal(t, key, []byte{0x01, 0x4D, 0x3C, 0x2B, 0x1A, 0xA1, 0xB2, 0xC3, 0xD4})
}
