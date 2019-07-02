package claim

import (
	"testing"

	"github.com/magiconair/properties/assert"
)

func TestClaimIDKey(t *testing.T) {
	var claimID uint64 = 0x1A2B3C4D
	key := key(claimID)
	assert.Equal(t, key, []byte{0x00, 0x0, 0x0, 0x0, 0x00, 0x1A, 0x2B, 0x3C, 0x4D})
}

