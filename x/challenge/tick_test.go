package challenge

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

func TestNewResponseEndBlock(t *testing.T) {
	ctx, k := fakePendingGameQueue()

	tags := k.NewResponseEndBlock(ctx)
	assert.Equal(t, sdk.Tags{}, tags)
}

func Test_pendingGameQueue(t *testing.T) {
	ctx, k := fakePendingGameQueue()

	q := k.pendingGameQueue(ctx)
	assert.Equal(t, uint64(2), q.List.Len())

	err := k.checkPendingQueue(ctx, q)
	assert.Nil(t, err)
}
