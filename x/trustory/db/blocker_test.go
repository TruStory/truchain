package db

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/stretchr/testify/assert"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
)

func TestNewResponseEndBlock(t *testing.T) {
	ms, storyKey, voteKey := setupMultiStore()
	cdc := makeCodec()
	k := NewTruKeeper(storyKey, voteKey, bank.Keeper{}, cdc)
	ctx := sdk.NewContext(ms, abci.Header{}, false, log.NewNopLogger())
	_ = createFakeStory(ms, k)

	r := k.NewResponseEndBlock(ctx)
	assert.NotNil(t, r)
	assert.Equal(t, abci.ResponseEndBlock{}, r, "Invalid response end block")
}
