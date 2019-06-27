package staking

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"

	app "github.com/TruStory/truchain/types"
)

func TestKeeper_IterateAfterCreatedTimeUserStakes(t *testing.T) {
	ctx, k, accKeeper, _ := mockDB()
	addr := createFakeFundedAccount(ctx, accKeeper, sdk.Coins{sdk.NewInt64Coin(app.StakeDenom, app.Shanev*300)})

	_, err := k.SubmitArgument(ctx.WithBlockTime(mustParseTime("2019-01-01")),
		"arg1", "summary1", addr, 1, StakeChallenge)
	_, err = k.SubmitArgument(ctx.WithBlockTime(mustParseTime("2019-01-03")),
		"arg2", "summary2", addr, 1, StakeBacking)
	assert.NoError(t, err)
	_, err = k.SubmitArgument(ctx.WithBlockTime(mustParseTime("2019-01-05")),
		"arg3", "summary3", addr, 1, StakeChallenge)
	assert.NoError(t, err)
	_, err = k.SubmitArgument(ctx.WithBlockTime(mustParseTime("2019-01-07")),
		"arg4", "summary4", addr, 1, StakeBacking)
	assert.NoError(t, err)
	_, err = k.SubmitArgument(ctx.WithBlockTime(mustParseTime("2019-01-10")),
		"arg5", "summary5", addr, 1, StakeChallenge)
	assert.NoError(t, err)

	stakes := afterCreatedTimeStakes(ctx, k, addr, mustParseTime("2019-01-01"))
	assert.Len(t, stakes, 5)
	stakes = afterCreatedTimeStakes(ctx, k, addr, mustParseTime("2019-01-02"))
	assert.Len(t, stakes, 4)

	stakes = afterCreatedTimeStakes(ctx, k, addr, mustParseTime("2019-01-07"))
	assert.Len(t, stakes, 2)
}
