package db

import (
	"time"

	ts "github.com/TruStory/trucoin/x/trustory/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetStory gets the story with the given id from the key-value store
func (k TruKeeper) GetStory(ctx sdk.Context, storyID int64) (ts.Story, sdk.Error) {
	store := ctx.KVStore(k.storyKey)
	key := generateKey(k.storyKey.String(), storyID)
	val := store.Get(key)
	if val == nil {
		return ts.Story{}, ts.ErrStoryNotFound(storyID)
	}
	story := &ts.Story{}
	k.cdc.MustUnmarshalBinary(val, story)

	return *story, nil
}

// AddStory adds a story to the key-value store
func (k TruKeeper) AddStory(
	ctx sdk.Context,
	body string,
	category ts.StoryCategory,
	creator sdk.AccAddress,
	escrow sdk.AccAddress,
	storyType ts.StoryType,
	voteStart time.Time,
	voteEnd time.Time) (int64, sdk.Error) {

	store := ctx.KVStore(k.storyKey)

	story := ts.Story{
		ID:           k.newID(ctx, k.storyKey),
		Body:         body,
		Category:     category,
		CreatedBlock: ctx.BlockHeight(),
		Creator:      creator,
		Escrow:       escrow,
		State:        ts.Created,
		StoryType:    storyType,
		VoteStart:    voteStart,
		VoteEnd:      voteEnd,
	}

	key := generateKey(k.storyKey.String(), story.ID)
	val := k.cdc.MustMarshalBinary(story)
	store.Set(key, val)

	// add story to the active story queue (for in-progress stories)
	err := k.ActiveStoryQueuePush(ctx, story.ID)
	if err != nil {
		return -1, err
	}

	return story.ID, nil
}

// UpdateStory updates an existing story in the store
func (k TruKeeper) UpdateStory(ctx sdk.Context, story ts.Story) sdk.Error {
	newStory := ts.NewStory(
		story.ID,
		story.BondIDs,
		story.CommentIDs,
		story.EvidenceIDs,
		story.Thread,
		story.VoteIDs,
		story.Body,
		story.Category,
		story.CreatedBlock,
		story.Creator,
		story.Escrow,
		story.Round,
		story.State,
		story.SubmitBlock,
		story.StoryType,
		ctx.BlockHeight(),
		story.Users,
		story.VoteStart,
		story.VoteEnd)

	store := ctx.KVStore(k.storyKey)
	key := generateKey(k.storyKey.String(), story.ID)
	val := k.cdc.MustMarshalBinary(newStory)
	store.Set(key, val)

	return nil
}
