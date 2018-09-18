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
	err := k.cdc.UnmarshalBinary(val, story)
	if err != nil {
		panic(err)
	}
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

	val, err := k.cdc.MarshalBinary(story)
	if err != nil {
		panic(err)
	}

	key := generateKey(k.storyKey.String(), story.ID)
	store.Set(key, val)

	// add story to the active story queue (for in-progress stories)
	err = k.ActiveStoryQueuePush(ctx, story.ID)
	if err != nil {
		panic(err)
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
		story.Expiration,
		story.Round,
		story.State,
		story.SubmitBlock,
		story.StoryType,
		ctx.BlockHeight(),
		story.Users,
		story.VoteStart,
		story.VoteEnd)

	val, err := k.cdc.MarshalBinary(newStory)
	if err != nil {
		panic(err)
	}

	store := ctx.KVStore(k.storyKey)
	key := generateKey(k.storyKey.String(), story.ID)
	store.Set(key, val)

	return nil
}
