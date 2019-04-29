package db

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_GetPoolSize(t *testing.T) {
	assert.Equal(t, getPoolSize(), 25)
	os.Setenv("PG_POOL_SIZE", "100")
	assert.Equal(t, getPoolSize(), 100)
	os.Setenv("PG_POOL_SIZE", "abcdefg")
	assert.Equal(t, getPoolSize(), 25)
	os.Setenv("PG_POOL_SIZE", "0")
	assert.Equal(t, getPoolSize(), 25)
	os.Setenv("PG_POOL_SIZE", "-25")
	assert.Equal(t, getPoolSize(), 25)

}
