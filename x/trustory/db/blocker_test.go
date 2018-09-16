package db

import (
	"testing"
	"time"

	ts "github.com/TruStory/trucoin/x/trustory/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/stretchr/testify/assert"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tendermint/tendermint/libs/log"
)

func TestNotMeetVoteMinNewResponseEndBlock(t *testing.T) {
	ms, accKey, storyKey, voteKey := setupMultiStore()
	cdc := makeCodec()
	am := auth.NewAccountMapper(cdc, accKey, auth.ProtoBaseAccount)
	ck := bank.NewKeeper(am)
	k := NewTruKeeper(storyKey, voteKey, ck, cdc)

	// create fake context with fake block time in header
	time := time.Date(2018, time.September, 14, 23, 0, 0, 0, time.UTC)
	header := abci.Header{Time: time}
	ctx := sdk.NewContext(ms, header, false, log.NewNopLogger())

	// create fake story with vote end after block time
	_ = createFakeStory(ms, k)

	r := k.NewResponseEndBlock(ctx)
	assert.NotNil(t, r)
}

func TestMeetVoteMinNewResponseEndBlock(t *testing.T) {
	ms, accKey, storyKey, voteKey := setupMultiStore()
	cdc := makeCodec()
	am := auth.NewAccountMapper(cdc, accKey, auth.ProtoBaseAccount)
	ck := bank.NewKeeper(am)
	k := NewTruKeeper(storyKey, voteKey, ck, cdc)

	// create fake context with fake block time in header
	time := time.Date(2018, time.September, 14, 23, 0, 0, 0, time.UTC)
	header := abci.Header{Time: time}
	ctx := sdk.NewContext(ms, header, false, log.NewNopLogger())

	// create fake story with vote end after block time
	storyID := createFakeStoryWithEscrow(ctx, am, ms, k)

	// fund voter account
	_, _, addr := keyPubAddr()
	baseAcct := auth.NewBaseAccountWithAddress(addr)
	coins, _ := sdk.ParseCoins("5memecoin")
	_ = baseAcct.SetCoins(coins)
	am.SetAccount(ctx, &baseAcct)

	// fake 10 votes
	for i := 0; i < 10; i++ {
		_, _ = k.VoteStory(ctx, storyID, addr, true, coins)
	}

	r := k.NewResponseEndBlock(ctx)
	assert.NotNil(t, r)
}

func createFakeStoryWithEscrow(ctx sdk.Context, am auth.AccountMapper, ms sdk.MultiStore, k TruKeeper) int64 {
	body := "Body of story."
	category := ts.DEX
	creator := sdk.AccAddress([]byte{1, 2})

	coins, _ := sdk.ParseCoins("10memecoin")
	_, _, escrowAddr := keyPubAddr()
	baseAcct := auth.NewBaseAccountWithAddress(escrowAddr)
	_ = baseAcct.SetCoins(coins)
	am.SetAccount(ctx, &baseAcct)

	storyType := ts.Default
	t := time.Date(2018, time.September, 13, 23, 0, 0, 0, time.UTC)

	storyID, _ := k.AddStory(ctx, body, category, creator, escrowAddr, storyType, t, t)
	return storyID
}

func keyPubAddr() (crypto.PrivKey, crypto.PubKey, sdk.AccAddress) {
	key := ed25519.GenPrivKey()
	pub := key.PubKey()
	addr := sdk.AccAddress(pub.Address())
	return key, pub, addr
}
