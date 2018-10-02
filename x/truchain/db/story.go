package db

import (
	ts "github.com/TruStory/truchain/x/truchain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewStory adds a story to the key-value store
func (k TruKeeper) NewStory(
	ctx sdk.Context,
	body string,
	category ts.StoryCategory,
	creator sdk.AccAddress,
	storyType ts.StoryType) (int64, sdk.Error) {

	store := ctx.KVStore(k.storyKey)

	story := ts.Story{
		ID:           k.id(ctx, k.storyKey),
		Body:         body,
		Category:     category,
		CreatedBlock: ctx.BlockHeight(),
		Creator:      creator,
		State:        ts.Created,
		StoryType:    storyType,
	}

	key := key(k.storyKey.String(), story.ID)
	val := k.cdc.MustMarshalBinary(story)
	store.Set(key, val)

	return story.ID, nil
}

// GetStory gets the story with the given id from the key-value store
func (k TruKeeper) GetStory(ctx sdk.Context, storyID int64) (ts.Story, sdk.Error) {
	store := ctx.KVStore(k.storyKey)
	key := key(k.storyKey.String(), storyID)
	val := store.Get(key)
	if val == nil {
		return ts.Story{}, ts.ErrStoryNotFound(storyID)
	}
	story := &ts.Story{}
	k.cdc.MustUnmarshalBinary(val, story)

	return *story, nil
}

// UpdateStory updates an existing story in the store
func (k TruKeeper) UpdateStory(ctx sdk.Context, story ts.Story) {
	newStory := ts.NewStory(
		story.ID,
		story.BackIDs,
		story.CommentIDs,
		story.EvidenceIDs,
		story.Thread,
		story.Body,
		story.Category,
		story.CreatedBlock,
		story.Creator,
		story.Round,
		story.State,
		story.StoryType,
		ctx.BlockHeight(),
		story.Users)

	store := ctx.KVStore(k.storyKey)
	key := key(k.storyKey.String(), story.ID)
	val := k.cdc.MustMarshalBinary(newStory)
	store.Set(key, val)
}
