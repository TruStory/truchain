package trustory

import (
	"encoding/binary"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	amino "github.com/tendermint/go-amino"
)

// StoryKeeper data type
type StoryKeeper struct {
	StoryKey sdk.StoreKey
	// CoinKey  bank.Keeper
	Cdc *wire.Codec
}

// NewStoryKeeper creates a new keeper with write and read access
func NewStoryKeeper(storyKey sdk.StoreKey, cdc *amino.Codec) StoryKeeper {
	return StoryKeeper{
		StoryKey: storyKey,
		// CoinKey:  coinKey,
		Cdc: cdc,
	}
}

// GetStory gets the story with the given id from the key-value store
func (sk StoryKeeper) GetStory(ctx sdk.Context, storyID int64) (Story, sdk.Error) {
	store := ctx.KVStore(sk.StoryKey)
	key := generateStoryKey(storyID)
	val := store.Get(key)
	if val == nil {
		return Story{}, ErrStoryNotFound(storyID)
	}
	story := &Story{}
	err := sk.Cdc.UnmarshalBinary(val, story)
	if err != nil {
		panic(err)
	}
	return *story, nil
}

// AddStory adds a story to the key-value store
func (sk StoryKeeper) AddStory(ctx sdk.Context, body string, creator sdk.AccAddress) (int64, sdk.Error) {
	store := ctx.KVStore(sk.StoryKey)

	story := Story{
		ID:   sk.newStoryID(store),
		Body: body,
		// Category: category,
		Creator: creator,
	}

	val, err := sk.Cdc.MarshalBinary(story)
	if err != nil {
		panic(err)
	}

	key := generateStoryKey(story.ID)
	store.Set(key, val)

	return story.ID, nil
}

// ============================================================================

// newStoryID creates a new id for a story by incrementing the last story id by 1
func (sk StoryKeeper) newStoryID(store sdk.KVStore) int64 {
	lastStoryID := store.Get([]byte("StoryID"))
	if lastStoryID == nil {
		return 0
	}

	storyID := new(int64)
	err := sk.Cdc.UnmarshalBinary(lastStoryID, storyID)
	if err != nil {
		panic(err)
	}

	return (*storyID + 1)
}

// generateStoryKey creates a key of the form "stories"|{storyID}
func generateStoryKey(storyID int64) []byte {
	var key []byte
	storyIDBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(storyIDBytes, uint64(storyID))

	key = []byte("stories")
	key = append(key, storyIDBytes...)
	return key
}
