package db

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_key(t *testing.T) {
	bz1 := key("stories", int64(5))
	bz2 := key("stories", int64(math.MaxInt64))

	assert.Equal(t, "stories:5", string(bz1), "should generate valid key")
	assert.Equal(t, "stories:9223372036854775807", string(bz2), "should generate valid key")
}
