package db

import (
	"testing"

	ts "github.com/TruStory/trucoin/x/trustory/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/stretchr/testify/assert"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
)

func TestActiveStoryQueue(t *testing.T) {
	ms, storyKey, voteKey := setupMultiStore()
	cdc := makeCodec()
	k := NewTruKeeper(storyKey, voteKey, bank.Keeper{}, cdc)
	ctx := sdk.NewContext(ms, abci.Header{}, false, log.NewNopLogger())
	storyID := createFakeStory(ms, k)

	_, err := k.ActiveStoryQueueHead(ctx)
	assert.NotNil(t, err)
	assert.Equal(t, sdk.CodeType(711), err.Code(), err.Error())

	_, err = k.ActiveStoryQueuePop(ctx)
	assert.NotNil(t, err)
	assert.Equal(t, sdk.CodeType(711), err.Code(), err.Error())

	// create an empty story queue
	k.setActiveStoryQueue(ctx, ts.ActiveStoryQueue{})

	_, err = k.ActiveStoryQueueHead(ctx)
	assert.NotNil(t, err)
	assert.Equal(t, sdk.CodeType(712), err.Code(), err.Error())

	_, err = k.ActiveStoryQueuePop(ctx)
	assert.NotNil(t, err)
	assert.Equal(t, sdk.CodeType(712), err.Code(), err.Error())

	err = k.ActiveStoryQueuePush(ctx, storyID)
	assert.Nil(t, err)

	story, _ := k.ActiveStoryQueueHead(ctx)
	assert.NotNil(t, story)

	story, _ = k.ActiveStoryQueuePop(ctx)
	assert.NotNil(t, story)
}
