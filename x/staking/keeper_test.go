package staking

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"

	app "github.com/TruStory/truchain/types"
	"github.com/TruStory/truchain/x/bank"
)

func TestKeeper_SubmitArgumentMaxLimit(t *testing.T) {
	ctx, k, mdb := mockDB()
	addr := createFakeFundedAccount(ctx, mdb.authAccKeeper, sdk.Coins{sdk.NewInt64Coin(app.StakeDenom, app.Shanev*300)})

	// max number of arguments
	arg1, err := k.SubmitArgument(ctx, "arg1", "summary1", addr, 1, StakeChallenge)
	assert.NoError(t, err)
	arg2, err := k.SubmitArgument(ctx, "arg2", "summary2", addr, 1, StakeBacking)
	assert.NoError(t, err)
	arg3, err := k.SubmitArgument(ctx, "arg3", "summary3", addr, 1, StakeChallenge)
	assert.NoError(t, err)
	arg4, err := k.SubmitArgument(ctx, "arg4", "summary4", addr, 1, StakeBacking)
	assert.NoError(t, err)
	arg5, err := k.SubmitArgument(ctx, "arg5", "summary5", addr, 1, StakeChallenge)
	assert.NoError(t, err)
	_, err = k.SubmitArgument(ctx, "arg6", "summary6", addr, 1, StakeBacking)
	assert.Error(t, err)
	assert.Equal(t, ErrorCodeMaxNumOfArgumentsReached, err.Code())
	userArguments := k.UserArguments(ctx, addr)
	assert.Equal(t, []Argument{arg1, arg2, arg3, arg4, arg5}, userArguments)
}

func TestKeeper_SubmitArgument(t *testing.T) {
	ctx, k, mdb := mockDB()
	ctx.WithBlockTime(time.Now())
	mockedAccountKeeper := mdb.accountKeeper.(*mockedAccountKeeper)
	addr := createFakeFundedAccount(ctx, mdb.authAccKeeper, sdk.Coins{sdk.NewInt64Coin(app.StakeDenom, app.Shanev*300)})
	addr2 := createFakeFundedAccount(ctx, mdb.authAccKeeper, sdk.Coins{sdk.NewInt64Coin(app.StakeDenom, app.Shanev*300)})
	mockedAccountKeeper.jail(addr)

	_, err := k.SubmitArgument(ctx, "body", "summary", addr, 1, StakeUpvote)
	assert.Error(t, err)
	assert.Equal(t, ErrorCodeInvalidStakeType, err.Code())

	_, err = k.SubmitArgument(ctx, "body", "summary", addr, 1, 127)
	assert.Error(t, err)
	assert.Equal(t, ErrorCodeInvalidStakeType, err.Code())

	_, err = k.SubmitArgument(ctx, "body", "summary", addr, 1, StakeBacking)
	assert.Error(t, err)
	assert.Equal(t, ErrorCodeAccountJailed, err.Code())
	_ = mdb.accountKeeper.UnJail(ctx, addr)

	argument, err := k.SubmitArgument(ctx, "body", "summary", addr, 1, StakeBacking)
	assert.NoError(t, err)
	expectedArgument := Argument{
		ID:           1,
		Creator:      addr,
		ClaimID:      1,
		Summary:      "summary",
		Body:         "body",
		StakeType:    StakeBacking,
		CreatedTime:  ctx.BlockHeader().Time,
		UpdatedTime:  ctx.BlockHeader().Time,
		UpvotedCount: 0,
		UpvotedStake: sdk.NewInt64Coin(app.StakeDenom, 0),
		TotalStake:   sdk.NewInt64Coin(app.StakeDenom, app.Shanev*50),
	}
	assert.Equal(t, expectedArgument, argument)
	argument, ok := k.getArgument(ctx, expectedArgument.ID)
	assert.True(t, ok)
	assert.Equal(t, expectedArgument, argument)

	expectedStake := Stake{
		ID:          1,
		ArgumentID:  1,
		Type:        StakeBacking,
		Amount:      sdk.NewInt64Coin(app.StakeDenom, app.Shanev*50),
		Creator:     addr,
		CreatedTime: ctx.BlockHeader().Time,
		EndTime:     ctx.BlockHeader().Time.Add(time.Hour * 24 * 7),
	}
	s, _ := k.getStake(ctx, 1)
	assert.Equal(t, expectedStake, s)
	argument2, err := k.SubmitArgument(ctx, "body2", "summary2", addr2, 1, StakeChallenge)
	expectedArgument2 := Argument{
		ID:           2,
		Creator:      addr2,
		ClaimID:      1,
		Summary:      "summary2",
		Body:         "body2",
		StakeType:    StakeChallenge,
		CreatedTime:  ctx.BlockHeader().Time,
		UpdatedTime:  ctx.BlockHeader().Time,
		UpvotedStake: sdk.NewInt64Coin(app.StakeDenom, 0),
		TotalStake:   sdk.NewInt64Coin(app.StakeDenom, app.Shanev*50),
	}
	expectedStake2 := Stake{
		ID:          2,
		ArgumentID:  2,
		Type:        StakeChallenge,
		Amount:      sdk.NewInt64Coin(app.StakeDenom, app.Shanev*50),
		Creator:     addr2,
		CreatedTime: ctx.BlockHeader().Time,
		EndTime:     ctx.BlockHeader().Time.Add(time.Hour * 24 * 7),
	}
	assert.NoError(t, err)
	assert.Equal(t, expectedArgument2, argument2)
	s, ok = k.getStake(ctx, 2)
	assert.True(t, ok)
	assert.Equal(t, expectedStake2, s)
	associatedArguments := k.ClaimArguments(ctx, 1)
	assert.Len(t, associatedArguments, 2)
	assert.Equal(t, expectedArgument, associatedArguments[0])
	assert.Equal(t, expectedArgument2, associatedArguments[1])

	associatedStakes := k.ArgumentStakes(ctx, expectedArgument.ID)
	assert.Len(t, associatedStakes, 1)
	assert.Equal(t, associatedStakes[0], expectedStake)

	// user <-> argument associations
	user1Arguments := k.UserArguments(ctx, addr)
	user2Arguments := k.UserArguments(ctx, addr2)

	assert.Len(t, user1Arguments, 1)
	assert.Len(t, user2Arguments, 1)

	assert.Equal(t, user1Arguments[0], expectedArgument)
	assert.Equal(t, user2Arguments[0], expectedArgument2)

	// user <-> stakes

	user1Stakes := k.UserStakes(ctx, addr)
	user2Stakes := k.UserStakes(ctx, addr2)

	assert.Len(t, user1Stakes, 1)
	assert.Len(t, user2Stakes, 1)

	assert.Equal(t, user1Stakes[0], expectedStake)
	assert.Equal(t, user2Stakes[0], expectedStake2)

	expiringStakes := make([]Stake, 0)

	k.IterateActiveStakeQueue(ctx, ctx.BlockHeader().Time, func(stake Stake) bool {
		expiringStakes = append(expiringStakes, stake)
		return false
	})
	// shouldn't have any expiring stake
	assert.Len(t, expiringStakes, 0)

	period := k.GetParams(ctx).Period
	k.IterateActiveStakeQueue(ctx, ctx.BlockHeader().Time.Add(period), func(stake Stake) bool {
		expiringStakes = append(expiringStakes, stake)
		return false
	})

	assert.Len(t, expiringStakes, 2)
	assert.Equal(t, []Stake{expectedStake, expectedStake2}, expiringStakes)
}

func TestKeeper_AfterTimeStakesIterator(t *testing.T) {
	ctx, k, mdb := mockDB()
	ctx = ctx.WithBlockTime(mustParseTime("2019-01-15"))
	addr, _ := sdk.AccAddressFromBech32("cosmos18pkfm85y3v65rrmn8f2y2z2ytenhq0943q5unm")
	addr2, _ := sdk.AccAddressFromBech32("cosmos16nycdfk7293jrj42rke9dg9ffz5qmj3kzrcddc")
	setCoins(ctx, mdb.authAccKeeper, sdk.Coins{sdk.NewInt64Coin(app.StakeDenom, app.Shanev*300)}, addr)
	setCoins(ctx, mdb.authAccKeeper, sdk.Coins{sdk.NewInt64Coin(app.StakeDenom, app.Shanev*300)}, addr2)

	_, err := k.SubmitArgument(ctx, "body", "summary", addr, 1, StakeBacking)
	assert.NoError(t, err)
	ctx = ctx.WithBlockTime(mustParseTime("2019-01-17"))
	_, err = k.SubmitArgument(ctx, "body", "summary", addr2, 2, StakeChallenge)
	assert.NoError(t, err)
	_, err = k.SubmitArgument(ctx, "body", "summary", addr2, 2, StakeBacking)
	assert.NoError(t, err)
	ctx = ctx.WithBlockTime(mustParseTime("2019-01-18"))
	_, err = k.SubmitArgument(ctx, "body", "summary", addr, 3, StakeBacking)
	assert.NoError(t, err)
	ctx = ctx.WithBlockTime(mustParseTime("2019-01-19"))
	_, err = k.SubmitArgument(ctx, "body", "summary", addr2, 4, StakeChallenge)
	assert.NoError(t, err)
	ctx = ctx.WithBlockTime(mustParseTime("2019-01-20"))
	_, err = k.SubmitArgument(ctx, "body", "summary", addr2, 5, StakeChallenge)
	assert.NoError(t, err)
	_, err = k.SubmitArgument(ctx, "body", "summary", addr2, 6, StakeBacking)
	assert.NoError(t, err)
	ctx = ctx.WithBlockTime(mustParseTime("2019-01-21"))
	_, err = k.SubmitArgument(ctx, "body", "summary", addr, 7, StakeBacking)
	assert.NoError(t, err)

	stakes := make([]Stake, 0)
	ctx = ctx.WithBlockTime(mustParseTime("2019-01-22"))
	k.IterateAfterCreatedTimeUserStakes(ctx, addr,
		mustParseTime("2019-01-17"), func(stake Stake) bool {
			stakes = append(stakes, stake)
			return false
		},
	)
	assert.Len(t, stakes, 2)
}

func TestKeeper_SubmitUpvote(t *testing.T) {
	ctx, k, mdb := mockDB()
	ctx.WithBlockTime(time.Now())
	addr := createFakeFundedAccount(ctx, mdb.authAccKeeper, sdk.Coins{sdk.NewInt64Coin(app.StakeDenom, app.Shanev*300)})
	addr2 := createFakeFundedAccount(ctx, mdb.authAccKeeper, sdk.Coins{sdk.NewInt64Coin(app.StakeDenom, app.Shanev*300)})
	addr3 := createFakeFundedAccount(ctx, mdb.authAccKeeper, sdk.Coins{sdk.NewInt64Coin(app.StakeDenom, app.Shanev*300)})
	argument, err := k.SubmitArgument(ctx, "body", "summary", addr, 1, StakeBacking)
	assert.NoError(t, err)
	expectedStake := Stake{
		ID:          1,
		ArgumentID:  1,
		Type:        StakeBacking,
		Amount:      sdk.NewInt64Coin(app.StakeDenom, app.Shanev*50),
		Creator:     addr,
		CreatedTime: ctx.BlockHeader().Time,
		EndTime:     ctx.BlockHeader().Time.Add(time.Hour * 24 * 7),
	}
	s, ok := k.getStake(ctx, 1)
	assert.True(t, ok)
	assert.Equal(t, expectedStake, s)
	_, err = k.SubmitUpvote(ctx, argument.ID, addr2)
	assert.NoError(t, err)
	expectedStake2 := Stake{
		ID:          2,
		ArgumentID:  1,
		Type:        StakeUpvote,
		Amount:      sdk.NewInt64Coin(app.StakeDenom, app.Shanev*10),
		Creator:     addr2,
		CreatedTime: ctx.BlockHeader().Time,
		EndTime:     ctx.BlockHeader().Time.Add(time.Hour * 24 * 7),
	}
	// fail if argument doesn't exist
	_, err = k.SubmitUpvote(ctx, 9999, addr)
	assert.Error(t, err)
	assert.Equal(t, ErrorCodeUnknownArgument, err.Code())
	// don't let stake twice
	_, err = k.SubmitUpvote(ctx, argument.ID, addr)
	assert.Error(t, err)
	assert.Equal(t, ErrorCodeDuplicateStake, err.Code())
	_, err = k.SubmitUpvote(ctx, argument.ID, addr2)
	assert.Error(t, err)
	assert.Equal(t, ErrorCodeDuplicateStake, err.Code())

	// user <-> stakes
	user1Stakes := k.UserStakes(ctx, addr)
	user2Stakes := k.UserStakes(ctx, addr2)

	assert.Len(t, user1Stakes, 1)
	assert.Len(t, user2Stakes, 1)

	assert.Equal(t, user1Stakes[0], expectedStake)
	assert.Equal(t, user2Stakes[0], expectedStake2)

	expiringStakes := make([]Stake, 0)

	k.IterateActiveStakeQueue(ctx, ctx.BlockHeader().Time, func(stake Stake) bool {
		expiringStakes = append(expiringStakes, stake)
		return false
	})
	// shouldn't have any expiring stake
	assert.Len(t, expiringStakes, 0)

	period := k.GetParams(ctx).Period
	k.IterateActiveStakeQueue(ctx, ctx.BlockHeader().Time.Add(period), func(stake Stake) bool {
		expiringStakes = append(expiringStakes, stake)
		return false
	})

	assert.Len(t, expiringStakes, 2)
	assert.Equal(t, []Stake{expectedStake, expectedStake2}, expiringStakes)

	userTxs := k.bankKeeper.TransactionsByAddress(ctx, addr)
	user2Txs := k.bankKeeper.TransactionsByAddress(ctx, addr2)
	assert.Len(t, userTxs, 1)
	assert.Len(t, user2Txs, 1)
	assert.Equal(t, bank.TransactionBacking, userTxs[0].Type)
	assert.Equal(t, bank.TransactionUpvote, user2Txs[0].Type)
	assert.Equal(t, sdk.NewInt64Coin(app.StakeDenom, app.Shanev*50), userTxs[0].Amount)
	assert.Equal(t, sdk.NewInt64Coin(app.StakeDenom, app.Shanev*10), user2Txs[0].Amount)

	_, err = k.SubmitUpvote(ctx, argument.ID, addr3)
	assert.NoError(t, err)

	argument, ok = k.getArgument(ctx, argument.ID)
	assert.True(t, ok)
	assert.Equal(t, uint64(2), argument.UpvotedCount)
	assert.Equal(t, sdk.NewInt64Coin(app.StakeDenom, app.Shanev*20), argument.UpvotedStake)
	assert.Equal(t, sdk.NewInt64Coin(app.StakeDenom, app.Shanev*70), argument.TotalStake)
}

func Test_interest(t *testing.T) {
	ctx, k, _ := mockDB()
	amount := sdk.NewInt64Coin(app.StakeDenom, 50000000000)
	now := time.Now()
	p := k.GetParams(ctx)
	after7days := now.Add(p.Period)
	interest := k.interest(ctx, amount, after7days.Sub(now))
	assert.Equal(t, sdk.NewInt(1006849315), interest.RoundInt())
}

func Test_splitReward(t *testing.T) {
	ctx, k, _ := mockDB()
	amount := sdk.NewInt64Coin(app.StakeDenom, 500000000000000)
	now := time.Now()
	p := k.GetParams(ctx)
	after7days := now.Add(p.Period)
	interest := k.interest(ctx, amount, after7days.Sub(now))
	t.Log("interest: " + interest.String())

	creatorReward, stakerReward := k.splitReward(ctx, interest)
	expectedCreatorReward := sdk.NewDecFromInt(sdk.NewInt(10068493150685)).
		Mul(sdk.NewDecWithPrec(50, 2))
	t.Log("expected creator reward: " + expectedCreatorReward.String())

	assert.True(t, amount.Amount.GT(interest.RoundInt()))
	assert.True(t, interest.RoundInt().GT(creatorReward))
	assert.True(t, interest.RoundInt().GT(stakerReward))
	assert.True(t, creatorReward.Equal(stakerReward))
	assert.Equal(t,
		expectedCreatorReward.RoundInt().String(),
		creatorReward.String(),
	)
	t.Log("actual creator reward: " + creatorReward.String())

	assert.Equal(t,
		interest.Sub(expectedCreatorReward).RoundInt().String(),
		stakerReward.String(),
	)
	t.Log("actual staker reward: " + stakerReward.String())
}

func TestKeeper_StakePeriodAmountLimit(t *testing.T) {
	ctx, k, mdb := mockDB()
	addr := createFakeFundedAccount(ctx, mdb.authAccKeeper, sdk.Coins{sdk.NewInt64Coin(app.StakeDenom, app.Shanev*215)})

	_, err := k.SubmitArgument(ctx.WithBlockTime(mustParseTime("2019-01-01")),
		"arg1", "summary1", addr, 1, StakeChallenge)
	_, err = k.SubmitArgument(ctx.WithBlockTime(mustParseTime("2019-01-05")),
		"arg2", "summary2", addr, 2, StakeBacking)
	assert.NoError(t, err)
	_, err = k.SubmitArgument(ctx.WithBlockTime(mustParseTime("2019-01-07")),
		"arg3", "summary3", addr, 3, StakeChallenge)
	assert.NoError(t, err)

	// should mark first stake as expired and refund stake
	EndBlocker(ctx.WithBlockTime(mustParseTime("2019-01-08")), k)
	_, err = k.SubmitArgument(ctx.WithBlockTime(mustParseTime("2019-01-08")),
		"arg4", "summary4", addr, 4, StakeBacking)

	_, err = k.SubmitArgument(ctx.WithBlockTime(mustParseTime("2019-01-10")),
		"arg5", "summary5", addr, 5, StakeChallenge)
	assert.Error(t, err)
	assert.Equal(t, ErrorCodeMaxAmountStakingReached, err.Code())

}

func TestKeeper_InsufficientCoins(t *testing.T) {
	ctx, k, mdb := mockDB()
	_, _, unfundedAddress := keyPubAddr()
	addr := createFakeFundedAccount(ctx, mdb.authAccKeeper, sdk.Coins{sdk.NewInt64Coin(app.StakeDenom, app.Shanev*300)})
	_, err := k.SubmitArgument(ctx, "body", "summary", unfundedAddress, 1, StakeBacking)
	assert.Error(t, err)
	assert.Equal(t, sdk.CodeInsufficientFunds, err.Code())

	argument, err := k.SubmitArgument(ctx, "body", "summary", addr, 1, StakeBacking)
	assert.NoError(t, err)

	_, err = k.SubmitUpvote(ctx, argument.ID, unfundedAddress)
	assert.Error(t, err)
	assert.Equal(t, sdk.CodeInsufficientFunds, err.Code())
}
