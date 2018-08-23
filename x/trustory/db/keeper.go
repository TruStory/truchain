package db

import (
	"encoding/binary"

	ts "github.com/TruStory/trucoin/x/trustory/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	amino "github.com/tendermint/go-amino"
)

// ============================================================================

// TruKeeper data type
type TruKeeper struct {
	StoryKey sdk.StoreKey
	VoteKey  sdk.StoreKey
	Cdc      *amino.Codec
}

// NewTruKeeper creates a new keeper with write and read access
func NewTruKeeper(storyKey sdk.StoreKey, voteKey sdk.StoreKey, cdc *amino.Codec) TruKeeper {
	return TruKeeper{
		StoryKey: storyKey,
		VoteKey:  voteKey,
		Cdc:      cdc,
	}
}

// GetStory gets the story with the given id from the key-value store
func (k TruKeeper) GetStory(ctx sdk.Context, storyID int64) (ts.Story, sdk.Error) {
	store := ctx.KVStore(k.StoryKey)
	key := generateKey(k.StoryKey.String(), storyID)
	val := store.Get(key)
	if val == nil {
		return ts.Story{}, ts.ErrStoryNotFound(storyID)
	}
	story := &ts.Story{}
	err := k.Cdc.UnmarshalBinary(val, story)
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
	storyType ts.StoryType) (int64, sdk.Error) {

	store := ctx.KVStore(k.StoryKey)

	story := ts.Story{
		ID:           k.newID(ctx, k.StoryKey),
		Body:         body,
		Category:     category,
		CreatedBlock: ctx.BlockHeight(),
		Creator:      creator,
		State:        ts.Created,
		StoryType:    storyType,
	}

	val, err := k.Cdc.MarshalBinary(story)
	if err != nil {
		panic(err)
	}

	key := generateKey(k.StoryKey.String(), story.ID)
	store.Set(key, val)

	return story.ID, nil
}

// VoteStory saves a vote to a story
func (k TruKeeper) VoteStory(ctx sdk.Context, storyID int64, creator sdk.AccAddress, choice bool, stake sdk.Coin) (int64, sdk.Error) {
	// access story DB
	storyStore := ctx.KVStore(k.StoryKey)
	storyKey := generateKey(k.StoryKey.String(), storyID)
	storyVal := storyStore.Get(storyKey)

	if storyVal == nil {
		return -1, ts.ErrStoryNotFound(storyID)
	}

	// get existing story
	story := &ts.Story{}
	err := k.Cdc.UnmarshalBinary(storyVal, story)
	if err != nil {
		panic(err)
	}

	// create new vote struct
	vote := ts.Vote{
		ID:           k.newID(ctx, k.VoteKey),
		StoryID:      story.ID,
		CreatedBlock: ctx.BlockHeight(),
		Creator:      creator,
		Round:        story.Round + 1,
		Stake:        stake,
		Vote:         choice,
	}

	// store vote in vote store
	voteStore := ctx.KVStore(k.VoteKey)
	voteKey := generateKey(k.VoteKey.String(), vote.ID)
	voteVal, err := k.Cdc.MarshalBinary(vote)
	if err != nil {
		panic(err)
	}
	voteStore.Set(voteKey, voteVal)

	// add vote id to story
	story.VoteIDs = append(story.VoteIDs, vote.ID)

	// create new story with vote
	newStory, err := k.Cdc.MarshalBinary(*story)
	if err != nil {
		panic(err)
	}

	// replace old story with new one in story store
	storyStore.Set(storyKey, newStory)

	return vote.ID, nil
}

// ============================================================================

// GetVote gets a vote with the given id from the key-value store
func (k TruKeeper) GetVote(ctx sdk.Context, voteID int64) (ts.Vote, sdk.Error) {
	store := ctx.KVStore(k.VoteKey)
	key := generateKey(k.VoteKey.String(), voteID)
	val := store.Get(key)
	if val == nil {
		return ts.Vote{}, ts.ErrVoteNotFound(voteID)
	}
	vote := &ts.Vote{}
	err := k.Cdc.UnmarshalBinary(val, vote)
	if err != nil {
		panic(err)
	}
	return *vote, nil
}

// ============================================================================

// newID creates a new id for a key by incrementing the last one by 1
func (k TruKeeper) newID(ctx sdk.Context, storeKey sdk.StoreKey) int64 {
	store := ctx.KVStore(storeKey)

	// create key of form "keyName|TotalID", i.e: "stories|TotalID" or "votes|TotalID"
	key := storeKey.Name() + "|TotalID"
	keyVal := []byte(key)
	lastID := store.Get(keyVal)

	// if we don't have an ID yet, set it to zero
	if lastID == nil {
		zero, err := k.Cdc.MarshalBinary(int64(0))
		if err != nil {
			panic(err)
		}
		store.Set(keyVal, zero)

		return 0
	}

	// convert from binary to int64
	ID := new(int64)
	err := k.Cdc.UnmarshalBinary(lastID, ID)
	if err != nil {
		panic(err)
	}

	// increment id by 1 and update value in blockchain
	newID := *ID + 1
	newIDVal, err := k.Cdc.MarshalBinary(newID)
	if err != nil {
		panic(err)
	}
	store.Set(keyVal, newIDVal)

	return newID
}

// generateKey creates a key of the form "keyName"|{id}
func generateKey(keyName string, id int64) []byte {
	var key []byte
	idBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(idBytes, uint64(id))
	key = []byte(keyName)
	key = append(key, idBytes...)
	return key
}
