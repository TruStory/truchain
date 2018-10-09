package db

import (
	"testing"

	ts "github.com/TruStory/truchain/x/truchain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"

	"github.com/stretchr/testify/assert"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/tendermint/tendermint/libs/log"
)

func TestAddGetStory(t *testing.T) {
	ms, _, storyKey, catKey, backingKey := setupMultiStore()
	cdc := makeCodec()
	keeper := NewTruKeeper(storyKey, catKey, backingKey, bank.BaseKeeper{}, cdc)
	ctx := sdk.NewContext(ms, abci.Header{}, false, log.NewNopLogger())
	storyID := CreateFakeStory(ms, keeper)

	// test getting a non-existant story
	_, err := keeper.GetStory(ctx, int64(5))
	assert.NotNil(t, err)

	// test getting an existing story
	savedStory, err := keeper.GetStory(ctx, storyID)
	assert.Nil(t, err)

	story := ts.Story{
		ID:           storyID,
		Body:         "Body of story.",
		CategoryID:   int64(1),
		CreatedBlock: int64(0),
		Creator:      sdk.AccAddress([]byte{1, 2}),
		State:        ts.Created,
		StoryType:    ts.Default,
	}

	assert.Equal(t, savedStory, story, "Story received from store does not match expected value")

	// test incrementing id by adding another story
	body := "Body of story 2."
	// category := fakeCategory(ctx, k)
	creator := sdk.AccAddress([]byte{3, 4})
	storyType := ts.Default

	storyID, _ = keeper.NewStory(ctx, body, int64(1), creator, storyType)
	assert.Equal(t, int64(2), storyID, "Story ID did not increment properly")
}
