package db

import (
	"encoding/binary"
	"strconv"
	"strings"
	"time"

	ts "github.com/TruStory/trucoin/x/trustory/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	amino "github.com/tendermint/go-amino"
)

// ============================================================================

// TruKeeper data type storing keys to the key-value store
type TruKeeper struct {
	// key to access coin store
	ck bank.Keeper

	// unexposed keys to access store from context
	storyKey sdk.StoreKey
	voteKey  sdk.StoreKey

	// wire codec for binary encoding/decoding
	cdc *amino.Codec
}

// NewTruKeeper creates a new keeper with write and read access
func NewTruKeeper(storyKey sdk.StoreKey, voteKey sdk.StoreKey, am auth.AccountMapper, ck bank.Keeper, cdc *amino.Codec) TruKeeper {
	return TruKeeper{
		ck:       ck,
		storyKey: storyKey,
		voteKey:  voteKey,
		cdc:      cdc,
	}
}

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

// VoteStory saves a vote to a story
func (k TruKeeper) VoteStory(ctx sdk.Context, storyID int64, creator sdk.AccAddress, choice bool, amount sdk.Coins) (int64, sdk.Error) {
	// access story DB
	storyStore := ctx.KVStore(k.storyKey)
	storyKey := generateKey(k.storyKey.String(), storyID)
	storyVal := storyStore.Get(storyKey)

	if storyVal == nil {
		return -1, ts.ErrStoryNotFound(storyID)
	}

	// get existing story
	story := &ts.Story{}
	err := k.cdc.UnmarshalBinary(storyVal, story)
	if err != nil {
		panic(err)
	}

	// create new vote struct
	vote := ts.Vote{
		ID:           k.newID(ctx, k.voteKey),
		StoryID:      story.ID,
		CreatedBlock: ctx.BlockHeight(),
		Creator:      creator,
		Round:        story.Round + 1,
		Amount:       amount,
		Vote:         choice,
	}

	// store vote in vote store
	voteStore := ctx.KVStore(k.voteKey)
	voteKey := generateKey(k.voteKey.String(), vote.ID)
	voteVal, err := k.cdc.MarshalBinary(vote)
	if err != nil {
		panic(err)
	}
	voteStore.Set(voteKey, voteVal)

	// add vote id to story
	story.VoteIDs = append(story.VoteIDs, vote.ID)

	// create new story with vote
	newStory, err := k.cdc.MarshalBinary(*story)
	if err != nil {
		panic(err)
	}

	// replace old story with new one in story store
	storyStore.Set(storyKey, newStory)

	// add vote to vote list
	votes, err := k.GetActiveVotes(ctx, story.ID)
	if err != nil {
		panic(err)
	}
	votes = append(votes, vote.ID)
	err = k.SetActiveVotes(ctx, story.ID, votes)

	return vote.ID, nil
}

// ============================================================================

// GetActiveVotes gets all votes for the current round of a story
func (k TruKeeper) GetActiveVotes(ctx sdk.Context, storyID int64) ([]int64, sdk.Error) {
	store := ctx.KVStore(k.storyKey)
	key := generateVoteListKey(storyID)
	val := store.Get(key)
	if val == nil {
		return []int64{}, nil // FIXME: add error
	}
	votes := &[]int64{}
	err := k.cdc.UnmarshalBinary(val, votes)
	if err != nil {
		panic(err)
	}
	return *votes, nil
}

// SetActiveVotes sets all votes for the current round of a story
func (k TruKeeper) SetActiveVotes(ctx sdk.Context, storyID int64, votes []int64) sdk.Error {
	store := ctx.KVStore(k.storyKey)
	key := generateVoteListKey(storyID)
	value := k.cdc.MustMarshalBinary(votes)
	store.Set(key, value)

	return nil
}

// ============================================================================

// GetVote gets a vote with the given id from the key-value store
func (k TruKeeper) GetVote(ctx sdk.Context, voteID int64) (ts.Vote, sdk.Error) {
	store := ctx.KVStore(k.voteKey)
	key := generateKey(k.voteKey.String(), voteID)
	val := store.Get(key)
	if val == nil {
		return ts.Vote{}, ts.ErrVoteNotFound(voteID)
	}
	vote := &ts.Vote{}
	err := k.cdc.UnmarshalBinary(val, vote)
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
		zero, err := k.cdc.MarshalBinary(int64(0))
		if err != nil {
			panic(err)
		}
		store.Set(keyVal, zero)

		return 0
	}

	// convert from binary to int64
	ID := new(int64)
	err := k.cdc.UnmarshalBinary(lastID, ID)
	if err != nil {
		panic(err)
	}

	// increment id by 1 and update value in blockchain
	newID := *ID + 1
	newIDVal, err := k.cdc.MarshalBinary(newID)
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

// generateVoteListKey creates a key for a vote list of form "stories|ID|votes"
func generateVoteListKey(storyID int64) []byte {
	return []byte(strings.Join([]string{"stories", strconv.Itoa(int(storyID)), "votes"}, ""))
}
