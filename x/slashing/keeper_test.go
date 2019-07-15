package slashing

import (
	"testing"

	"github.com/TruStory/truchain/x/staking"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

func TestNewSlash_Success(t *testing.T) {
	ctx, keeper := mockDB()

	staker := keeper.GetParams(ctx).SlashAdmins[0]
	_, err := keeper.stakingKeeper.SubmitArgument(ctx, "arg1", "summary1", staker, 1, staking.StakeChallenge)
	assert.NoError(t, err)

	stakeID := uint64(1)
	creator := keeper.GetParams(ctx).SlashAdmins[0]
	slash, err := keeper.CreateSlash(ctx, stakeID, SlashTypeUnhelpful, SlashReasonPlagiarism, "", creator)
	assert.NoError(t, err)

	assert.NotZero(t, slash.ID)
	assert.Equal(t, slash.StakeID, stakeID)
	assert.Equal(t, slash.Creator, creator)
}

func TestNewSlash_InvalidStake(t *testing.T) {
	ctx, keeper := mockDB()

	invalidStakeID := uint64(404)
	creator := keeper.GetParams(ctx).SlashAdmins[0]
	_, err := keeper.CreateSlash(ctx, invalidStakeID, SlashTypeUnhelpful, SlashReasonPlagiarism, "", creator)

	assert.NotNil(t, err)
	assert.Equal(t, ErrInvalidStake(invalidStakeID).Code(), err.Code())
}

func TestNewSlash_InvalidDetailedReason(t *testing.T) {
	ctx, keeper := mockDB()

	stakeID := uint64(1)
	creator := keeper.GetParams(ctx).SlashAdmins[0]
	longDetailedReason := "This is a very very very descriptive reason to slash an argument. I am writing it in this detail to make the validation fail. I hope it works!"
	_, err := keeper.CreateSlash(ctx, stakeID, SlashTypeUnhelpful, SlashReasonOther, longDetailedReason, creator)

	assert.NotNil(t, err)
	assert.Equal(t, ErrInvalidSlashReason("").Code(), err.Code())
}

func TestNewSlash_ErrNotEnoughEarnedStake(t *testing.T) {
	ctx, keeper := mockDB()

	stakeID := uint64(1)
	invalidCreator := sdk.AccAddress([]byte{1, 2})
	_, err := keeper.CreateSlash(ctx, stakeID, SlashTypeUnhelpful, SlashReasonPlagiarism, "", invalidCreator)

	assert.NotNil(t, err)
	assert.Equal(t, ErrNotEnoughEarnedStake(invalidCreator).Code(), err.Code())
}

func TestSlash_Success(t *testing.T) {
	ctx, keeper := mockDB()

	stakeID := uint64(1)
	creator := keeper.GetParams(ctx).SlashAdmins[0]
	createdSlash, err := keeper.CreateSlash(ctx, stakeID, SlashTypeUnhelpful, SlashReasonPlagiarism, "", creator)
	assert.Nil(t, err)

	returnedSlash, err := keeper.Slash(ctx, createdSlash.ID)

	assert.Nil(t, err)
	assert.Equal(t, createdSlash, returnedSlash)
}

func TestSlash_ErrNotFound(t *testing.T) {
	ctx, keeper := mockDB()

	stakeID := uint64(1)
	creator := keeper.GetParams(ctx).SlashAdmins[0]
	_, err := keeper.CreateSlash(ctx, stakeID, SlashTypeUnhelpful, SlashReasonPlagiarism, "", creator)
	assert.Nil(t, err)

	_, err = keeper.Slash(ctx, uint64(404))

	assert.NotNil(t, err)
	assert.Equal(t, ErrSlashNotFound(uint64(404)).Code(), err.Code())
}

func TestSlashes_Success(t *testing.T) {
	ctx, keeper := mockDB()

	staker := keeper.GetParams(ctx).SlashAdmins[0]
	_, err := keeper.stakingKeeper.SubmitArgument(ctx, "arg1", "summary1", staker, 1, staking.StakeBacking)
	assert.NoError(t, err)

	stakeID := uint64(1)
	creator := keeper.GetParams(ctx).SlashAdmins[0]
	first, err := keeper.CreateSlash(ctx, stakeID, SlashTypeUnhelpful, SlashReasonPlagiarism, "", creator)
	assert.NoError(t, err)

	creator2 := keeper.GetParams(ctx).SlashAdmins[1]
	another, err := keeper.CreateSlash(ctx, stakeID, SlashTypeUnhelpful, SlashReasonPlagiarism, "", creator2)
	assert.NoError(t, err)

	all := keeper.Slashes(ctx)
	assert.Len(t, all, 2)
	assert.Equal(t, all[0], first)
	assert.Equal(t, all[1], another)
}

func Test_punishment(t *testing.T) {
	ctx, keeper := mockDB()

	staker := keeper.GetParams(ctx).SlashAdmins[0]
	slasher := keeper.GetParams(ctx).SlashAdmins[1]
	stakerStartingBalance := keeper.bankKeeper.GetCoins(ctx, staker)
	slasherStartingBalance := keeper.bankKeeper.GetCoins(ctx, slasher)
	assert.Equal(t, "100000000000trusteak", stakerStartingBalance.String())
	assert.Equal(t, "100000000000trusteak", slasherStartingBalance.String())

	_, err := keeper.stakingKeeper.SubmitArgument(ctx, "arg1", "summary1", staker, 1, staking.StakeChallenge)
	assert.NoError(t, err)

	stake, _ := keeper.stakingKeeper.Stake(ctx, 1)
	assert.Equal(t, "50000000000trusteak", stake.Amount.String())

	// this also does a punish because slasher is an admin
	_, err = keeper.CreateSlash(ctx, stake.ID, SlashTypeUnhelpful, SlashReasonPlagiarism, "", slasher)
	assert.NoError(t, err)

	//staker should have = starting balance - stake amount
	stakerEndingBalance := keeper.bankKeeper.GetCoins(ctx, staker)
	expectedBalance := stakerStartingBalance.Sub(sdk.Coins{stake.Amount})
	assert.Equal(t, expectedBalance.String(), stakerEndingBalance.String())

	// slasher should have = starting balance + reward (25% stake)
	slasherEndingBalance := keeper.bankKeeper.GetCoins(ctx, slasher)
	reward := stake.Amount.Amount.ToDec().Mul(sdk.NewDecWithPrec(25, 2)).TruncateInt()
	rewardCoin := sdk.NewCoin(stake.Amount.Denom, reward)
	expectedBalance = slasherStartingBalance.Add(sdk.Coins{rewardCoin})
	assert.Equal(t, expectedBalance.String(), slasherEndingBalance.String())
}
