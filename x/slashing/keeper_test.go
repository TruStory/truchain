package slashing

import (
	"testing"

	"github.com/TruStory/truchain/x/staking"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

func TestNewSlash_Success(t *testing.T) {
	ctx, keeper := mockDB()

	staker := keeper.GetParams(ctx).SlashAdmins[1]
	arg, err := keeper.stakingKeeper.SubmitArgument(ctx, "arg1", "summary1", staker, 1, staking.StakeChallenge)
	assert.NoError(t, err)

	stakeID := uint64(1)
	creator := keeper.GetParams(ctx).SlashAdmins[1]
	slash, _, err := keeper.CreateSlash(ctx, stakeID, SlashTypeUnhelpful, SlashReasonPlagiarism, "", creator)
	assert.NoError(t, err)

	assert.NotZero(t, slash.ID)
	assert.Equal(t, uint64(2), arg.ID)
	assert.Equal(t, slash.Creator, creator)
}

func TestNewSlash_InvalidArgument(t *testing.T) {
	ctx, keeper := mockDB()

	invalidArgumentID := uint64(404)
	creator := keeper.GetParams(ctx).SlashAdmins[0]
	_, _, err := keeper.CreateSlash(ctx, invalidArgumentID, SlashTypeUnhelpful, SlashReasonPlagiarism, "", creator)

	assert.NotNil(t, err)
	assert.Equal(t, ErrInvalidArgument(invalidArgumentID).Code(), err.Code())
}

func TestNewSlash_InvalidDetailedReason(t *testing.T) {
	ctx, keeper := mockDB()

	stakeID := uint64(1)
	creator := keeper.GetParams(ctx).SlashAdmins[0]
	longDetailedReason := "This is a very very very descriptive reason to slash an argument. I am writing it in this detail to make the validation fail. I hope it works!"
	_, _, err := keeper.CreateSlash(ctx, stakeID, SlashTypeUnhelpful, SlashReasonOther, longDetailedReason, creator)

	assert.NotNil(t, err)
	assert.Equal(t, ErrInvalidSlashReason("").Code(), err.Code())
}

func TestNewSlash_ErrNotEnoughEarnedStake(t *testing.T) {
	ctx, keeper := mockDB()

	stakeID := uint64(1)
	invalidCreator := sdk.AccAddress([]byte{1, 2})
	_, _, err := keeper.CreateSlash(ctx, stakeID, SlashTypeUnhelpful, SlashReasonPlagiarism, "", invalidCreator)

	assert.NotNil(t, err)
	assert.Equal(t, ErrNotEnoughEarnedStake(invalidCreator).Code(), err.Code())
}

func TestSlash_Success(t *testing.T) {
	ctx, keeper := mockDB()

	stakeID := uint64(1)
	creator := keeper.GetParams(ctx).SlashAdmins[0]
	createdSlash, _, err := keeper.CreateSlash(ctx, stakeID, SlashTypeUnhelpful, SlashReasonPlagiarism, "", creator)
	assert.Nil(t, err)

	returnedSlash, err := keeper.Slash(ctx, createdSlash.ID)
	assert.NoError(t, err)
	assert.Equal(t, createdSlash, returnedSlash)
}

func TestSlash_ErrNotFound(t *testing.T) {
	ctx, keeper := mockDB()

	stakeID := uint64(1)
	creator := keeper.GetParams(ctx).SlashAdmins[0]
	_, _, err := keeper.CreateSlash(ctx, stakeID, SlashTypeUnhelpful, SlashReasonPlagiarism, "", creator)
	assert.Nil(t, err)

	_, err = keeper.Slash(ctx, uint64(404))

	assert.NotNil(t, err)
	assert.Equal(t, ErrSlashNotFound(uint64(404)).Code(), err.Code())
}

func TestSlashes_Success(t *testing.T) {
	ctx, keeper := mockDB()

	staker := keeper.GetParams(ctx).SlashAdmins[1]
	_, err := keeper.stakingKeeper.SubmitArgument(ctx, "arg1", "summary1", staker, 1, staking.StakeBacking)
	assert.NoError(t, err)

	stakeID := uint64(1)
	creator := keeper.GetParams(ctx).SlashAdmins[0]
	first, _, err := keeper.CreateSlash(ctx, stakeID, SlashTypeUnhelpful, SlashReasonPlagiarism, "", creator)
	assert.NoError(t, err)

	creator2 := keeper.GetParams(ctx).SlashAdmins[1]
	another, _, err := keeper.CreateSlash(ctx, stakeID, SlashTypeUnhelpful, SlashReasonPlagiarism, "", creator2)
	assert.NoError(t, err)

	all := keeper.Slashes(ctx)
	assert.Len(t, all, 2)
	assert.Equal(t, all[0], first)
	assert.Equal(t, all[1], another)
}

func Test_punishment(t *testing.T) {
	ctx, keeper := mockDB()

	staker := keeper.GetParams(ctx).SlashAdmins[1]
	slasher := keeper.GetParams(ctx).SlashAdmins[2]
	slashMagnitude := keeper.GetParams(ctx).SlashMagnitude
	stakerStartingBalance := keeper.bankKeeper.GetCoins(ctx, staker)
	slasherStartingBalance := keeper.bankKeeper.GetCoins(ctx, slasher)
	assert.Equal(t, "300000000000tru", stakerStartingBalance.String())
	assert.Equal(t, "300000000000tru", slasherStartingBalance.String())

	claim, _ := keeper.claimKeeper.Claim(ctx, 1)
	assert.Equal(t, "0tru", claim.TotalChallenged.String())

	argument, err := keeper.stakingKeeper.SubmitArgument(ctx, "arg2", "summary2", staker, claim.ID, staking.StakeChallenge)
	assert.NoError(t, err)

	stake, _ := keeper.stakingKeeper.Stake(ctx, 2)
	assert.Equal(t, argument.ID, stake.ArgumentID)
	assert.Equal(t, "50000000000tru", stake.Amount.String())

	claim, _ = keeper.claimKeeper.Claim(ctx, 1)
	assert.Equal(t, stake.Amount.String(), claim.TotalChallenged.String())

	// staker should have = starting balance - stake amount
	stakerEndingBalance := keeper.bankKeeper.GetCoins(ctx, staker)
	expectedBalance := stakerStartingBalance.Sub(sdk.Coins{stake.Amount})
	assert.Equal(t, expectedBalance.String(), stakerEndingBalance.String())

	// this also does a punish because slasher is an admin
	_, _, err = keeper.CreateSlash(ctx, argument.ID, SlashTypeUnhelpful, SlashReasonPlagiarism, "", slasher)
	assert.NoError(t, err)

	// staker should have = starting balance - (stake amount * slashMagnitude)
	slashPenalty := sdk.NewCoin(stake.Amount.Denom, stake.Amount.Amount.MulRaw(int64(slashMagnitude)))
	stakerEndingBalance = keeper.bankKeeper.GetCoins(ctx, staker)
	expectedBalance = stakerStartingBalance.Sub(sdk.Coins{slashPenalty})
	assert.Equal(t, expectedBalance.String(), stakerEndingBalance.String())

	// slasher should have = starting balance + reward (25% stake)
	slasherEndingBalance := keeper.bankKeeper.GetCoins(ctx, slasher)
	reward := stake.Amount.Amount.ToDec().Mul(sdk.NewDecWithPrec(25, 2)).TruncateInt()
	rewardCoin := sdk.NewCoin(stake.Amount.Denom, reward)
	expectedBalance = slasherStartingBalance.Add(sdk.Coins{rewardCoin})
	assert.Equal(t, expectedBalance.String(), slasherEndingBalance.String())

	claim, _ = keeper.claimKeeper.Claim(ctx, 1)
	assert.Equal(t, "0tru", claim.TotalChallenged.String())
}

func TestAddAdmin_Success(t *testing.T) {
	ctx, keeper := mockDB()

	creator := keeper.GetParams(ctx).SlashAdmins[0]
	_, _, newAdmin, _ := getFakeAppAccountParams()

	err := keeper.AddAdmin(ctx, newAdmin, creator)
	assert.Nil(t, err)

	newAdmins := keeper.GetParams(ctx).SlashAdmins
	assert.Subset(t, newAdmins, []sdk.AccAddress{newAdmin})
}

func TestAddAdmin_CreatorNotAuthorised(t *testing.T) {
	ctx, keeper := mockDB()

	invalidCreator := sdk.AccAddress([]byte{1, 2})
	_, _, newAdmin, _ := getFakeAppAccountParams()

	err := keeper.AddAdmin(ctx, newAdmin, invalidCreator)
	assert.NotNil(t, err)
	assert.Equal(t, ErrAddressNotAuthorised().Code(), err.Code())
}

func TestRemoveAdmin_Success(t *testing.T) {
	ctx, keeper := mockDB()

	currentAdmins := keeper.GetParams(ctx).SlashAdmins
	adminToRemove := currentAdmins[0]

	err := keeper.RemoveAdmin(ctx, adminToRemove, adminToRemove) // removing self
	assert.Nil(t, err)
	newAdmins := keeper.GetParams(ctx).SlashAdmins
	assert.Equal(t, len(currentAdmins)-1, len(newAdmins))
}

func TestRemoveAdmin_RemoverNotAuthorised(t *testing.T) {
	ctx, keeper := mockDB()

	invalidRemover := sdk.AccAddress([]byte{1, 2})
	currentAdmins := keeper.GetParams(ctx).SlashAdmins
	adminToRemove := currentAdmins[0]

	err := keeper.AddAdmin(ctx, adminToRemove, invalidRemover)
	assert.NotNil(t, err)
	assert.Equal(t, ErrAddressNotAuthorised().Code(), err.Code())
}
