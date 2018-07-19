package trustory

import (
	"encoding/binary"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/stake"
)

// GenerateStoryKey creates a key of the form "story"|{storyID}
func GenerateStoryKey(storyID int64) []byte {
	var key []byte
	storyIDBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(storyIDBytes, uint64(storyID))

	key = []byte("stories")
	key = append(key, storyIDBytes...)
	return key
}

// Keeper data type
type Keeper struct {
	TruStory  sdk.StoreKey
	codespace sdk.CodespaceType
	cdc       *wire.Codec

	ck bank.Keeper
	sm stake.Keeper
}

// NewKeeper creates a new keeper with write and read access
func NewKeeper(TruStory sdk.StoreKey, ck bank.Keeper, sm stake.Keeper, codespace sdk.CodespaceType) Keeper {
	cdc := wire.NewCodec()

	return Keeper{
		TruStory:  TruStory,
		cdc:       cdc,
		ck:        ck,
		sm:        sm,
		codespace: codespace,
	}
}

// NewStoryID created a new id for a story: not sure how this works
func (k Keeper) NewStoryID(ctx sdk.Context) int64 {
	store := ctx.KVStore(k.TruStory)
	bid := store.Get([]byte("TotalID"))
	if bid == nil {
		return 0
	}

	totalID := new(int64)
	err := k.cdc.UnmarshalBinary(bid, totalID)
	if err != nil {
		panic(err)
	}

	return (*totalID + 1)
}

// GetStory gets the story with the given id from the context
func (k Keeper) GetStory(ctx sdk.Context, storyID int64) (Story, sdk.Error) {
	store := ctx.KVStore(k.TruStory)

	key := GenerateStoryKey(storyID)
	bp := store.Get(key)
	if bp == nil {
		return Story{}, ErrStoryNotFound(storyID)
	}

	story := Story{}

	err := k.cdc.UnmarshalBinary(bp, story)
	if err != nil {
		panic(err)
	}

	return story, nil
}

// SetStory sets a story to the context
func (k Keeper) SetStory(ctx sdk.Context, storyID int64, story Story) sdk.Error {
	store := ctx.KVStore(k.TruStory)

	bp, err := k.cdc.MarshalBinary(story)
	if err != nil {
		panic(err)
	}

	key := GenerateStoryKey(storyID)

	store.Set(key, bp)
	return nil
}
