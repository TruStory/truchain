package db

import (
	"testing"

	ts "github.com/TruStory/trucoin/x/trustory/types"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"

	"github.com/stretchr/testify/assert"
	amino "github.com/tendermint/go-amino"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
)

func TestAddGetStory(t *testing.T) {
	ms, storyKey, voteKey := setupMultiStore()
	cdc := makeCodec()
	keeper := NewTruKeeper(storyKey, voteKey, auth.AccountMapper{}, bank.Keeper{}, cdc)
	ctx := sdk.NewContext(ms, abci.Header{}, false, log.NewNopLogger())
	storyID := createFakeStory(ms, keeper)

	// test getting a non-existant story
	_, err := keeper.GetStory(ctx, int64(5))
	assert.NotNil(t, err)

	// test getting an existing story
	savedStory, err := keeper.GetStory(ctx, storyID)
	assert.Nil(t, err)

	story := ts.Story{
		ID:           storyID,
		Body:         "Body of story.",
		Category:     ts.DEX,
		CreatedBlock: int64(0),
		Creator:      sdk.AccAddress([]byte{1, 2}),
		State:        ts.Created,
		StoryType:    ts.Default,
	}

	assert.Equal(t, savedStory, story, "Story received from store does not match expected value")

	// test incrementing id by adding another story
	body := "Body of story 2."
	category := ts.Bitcoin
	creator := sdk.AccAddress([]byte{3, 4})
	storyType := ts.Default

	storyID, _ = keeper.AddStory(ctx, body, category, creator, storyType)
	assert.Equal(t, int64(1), storyID, "Story ID did not increment properly")
}

func TestVoteStory(t *testing.T) {
	ms, storyKey, voteKey := setupMultiStore()
	cdc := makeCodec()
	keeper := NewTruKeeper(storyKey, voteKey, auth.AccountMapper{}, bank.Keeper{}, cdc)
	ctx := sdk.NewContext(ms, abci.Header{}, false, log.NewNopLogger())

	storyID := createFakeStory(ms, keeper)

	creator := sdk.AccAddress([]byte{3, 4})
	vote := true
	stake, _ := sdk.ParseCoins("10trustake")

	// test voting on a non-existant story
	_, err := keeper.VoteStory(ctx, int64(5), creator, vote, stake)
	assert.NotNil(t, err)

	// test voting on a story
	voteID, err := keeper.VoteStory(ctx, storyID, creator, vote, stake)
	assert.Nil(t, err)
	assert.Equal(t, voteID, int64(0), "Vote ID does not match")

	// test getting a non-existant vote
	_, err = keeper.GetVote(ctx, int64(5))
	assert.NotNil(t, err)

	// test getting vote and comparing fields
	savedVote, err := keeper.GetVote(ctx, voteID)
	assert.Nil(t, err)
	assert.Equal(t, savedVote.Vote, true, "Vote choice  does not match")

	assert.Equal(t, savedVote.Amount.AmountOf("trustake"), sdk.NewInt(10), "Vote amount does not match")
}

// ============================================================================

func createFakeStory(ms sdk.MultiStore, k TruKeeper) int64 {
	ctx := sdk.NewContext(ms, abci.Header{}, false, log.NewNopLogger())
	body := "Body of story."
	category := ts.DEX
	creator := sdk.AccAddress([]byte{1, 2})
	storyType := ts.Default

	storyID, _ := k.AddStory(ctx, body, category, creator, storyType)
	return storyID
}

func setupMultiStore() (sdk.MultiStore, *sdk.KVStoreKey, *sdk.KVStoreKey) {
	db := dbm.NewMemDB()
	storyKey := sdk.NewKVStoreKey("stories")
	voteKey := sdk.NewKVStoreKey("votes")
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(storyKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(voteKey, sdk.StoreTypeIAVL, db)
	ms.LoadLatestVersion()
	return ms, storyKey, voteKey
}

func makeCodec() *amino.Codec {
	cdc := amino.NewCodec()
	ts.RegisterAmino(cdc)
	crypto.RegisterAmino(cdc)
	cdc.RegisterInterface((*auth.Account)(nil), nil)
	cdc.RegisterConcrete(&auth.BaseAccount{}, "cosmos-sdk/BaseAccount", nil)
	return cdc
}
