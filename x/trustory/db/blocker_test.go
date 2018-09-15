package db

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/stretchr/testify/assert"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
)

func TestNotMeetVoteMinNewResponseEndBlock(t *testing.T) {
	ms, storyKey, voteKey := setupMultiStore()
	cdc := makeCodec()
	k := NewTruKeeper(storyKey, voteKey, bank.Keeper{}, cdc)

	// create fake context with fake block time in header
	time := time.Date(2018, time.September, 14, 23, 0, 0, 0, time.UTC)
	header := abci.Header{Time: time}
	ctx := sdk.NewContext(ms, header, false, log.NewNopLogger())

	// create fake story with vote end after block time
	_ = createFakeStory(ms, k)

	r := k.NewResponseEndBlock(ctx)
	assert.NotNil(t, r)
}

func TestMeetVoteMinNewResponseEndBlock(t *testing.T) {
	ms, storyKey, voteKey := setupMultiStore()
	cdc := makeCodec()
	k := NewTruKeeper(storyKey, voteKey, bank.Keeper{}, cdc)

	// create fake context with fake block time in header
	time := time.Date(2018, time.September, 14, 23, 0, 0, 0, time.UTC)
	header := abci.Header{Time: time}
	ctx := sdk.NewContext(ms, header, false, log.NewNopLogger())

	// create fake story with vote end after block time
	storyID := createFakeStory(ms, k)

	fakeAddr := sdk.AccAddress([]byte{1, 2})
	fakeCoins, _ := sdk.ParseCoins("10memecoin")

	_, _ = k.VoteStory(ctx, storyID, fakeAddr, true, fakeCoins)

	r := k.NewResponseEndBlock(ctx)
	assert.NotNil(t, r)
}
