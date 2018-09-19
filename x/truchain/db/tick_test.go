package db

import (
	"testing"
	"time"

	ts "github.com/TruStory/truchain/x/truchain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/stretchr/testify/assert"
)

func TestNotMeetVoteMinNewResponseEndBlock(t *testing.T) {
	ctx, ms, _, k := mockDB()

	// create fake story with vote end after block time
	_ = createFakeStory(ms, k)

	r := k.NewResponseEndBlock(ctx)
	assert.NotNil(t, r)
}

func TestMeetVoteMinNewResponseEndBlock(t *testing.T) {
	ctx, ms, am, k := mockDB()

	// create fake story with vote end after block time
	storyID := createFakeStoryWithEscrow(ctx, am, ms, k)

	// fund voter account
	coins, _ := sdk.ParseCoins("5memecoin")
	addr := createFundedAccount(ctx, am, coins)

	// fake 10 votes
	for i := 0; i < 10; i++ {
		_, _ = k.VoteStory(ctx, storyID, addr, true, coins)
	}

	r := k.NewResponseEndBlock(ctx)
	assert.NotNil(t, r)
}

// ============================================================================

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

	storyID, _ := k.AddStory(ctx, body, category, creator, escrowAddr, storyType, 10, t, t)
	return storyID
}
