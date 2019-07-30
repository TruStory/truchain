package staking

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"

	app "github.com/TruStory/truchain/types"
	"github.com/TruStory/truchain/x/bank"
	"github.com/TruStory/truchain/x/claim"
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

	_, err = k.SubmitArgument(ctx, "body", "summary", addr, 1, StakeType(0xFF))
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
		CommunityID:  "testunit",
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
	argument, ok := k.Argument(ctx, expectedArgument.ID)
	assert.True(t, ok)
	assert.Equal(t, expectedArgument, argument)

	expectedStake := Stake{
		ID:          1,
		ArgumentID:  1,
		CommunityID: "testunit",
		Type:        StakeBacking,
		Amount:      sdk.NewInt64Coin(app.StakeDenom, app.Shanev*50),
		Creator:     addr,
		CreatedTime: ctx.BlockHeader().Time,
		EndTime:     ctx.BlockHeader().Time.Add(time.Hour * 24 * 7),
	}
	s, _ := k.Stake(ctx, 1)
	assert.Equal(t, expectedStake, s)
	argument2, err := k.SubmitArgument(ctx, "body2", "summary2", addr2, 1, StakeChallenge)
	expectedArgument2 := Argument{
		ID:           2,
		Creator:      addr2,
		ClaimID:      1,
		CommunityID:  "testunit",
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
		CommunityID: "testunit",
		Type:        StakeChallenge,
		Amount:      sdk.NewInt64Coin(app.StakeDenom, app.Shanev*50),
		Creator:     addr2,
		CreatedTime: ctx.BlockHeader().Time,
		EndTime:     ctx.BlockHeader().Time.Add(time.Hour * 24 * 7),
	}
	assert.NoError(t, err)
	assert.Equal(t, expectedArgument2, argument2)
	s, ok = k.Stake(ctx, 2)
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
		CommunityID: "testunit",
		Type:        StakeBacking,
		Amount:      sdk.NewInt64Coin(app.StakeDenom, app.Shanev*50),
		Creator:     addr,
		CreatedTime: ctx.BlockHeader().Time,
		EndTime:     ctx.BlockHeader().Time.Add(time.Hour * 24 * 7),
	}
	s, ok := k.Stake(ctx, 1)
	assert.True(t, ok)
	assert.Equal(t, expectedStake, s)
	_, err = k.SubmitUpvote(ctx, argument.ID, addr2)
	assert.NoError(t, err)
	expectedStake2 := Stake{
		ID:          2,
		ArgumentID:  1,
		CommunityID: "testunit",
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

	argument, ok = k.Argument(ctx, argument.ID)
	assert.True(t, ok)
	assert.Equal(t, 2, argument.UpvotedCount)
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

func TestKeeper_StakeLimitTiers(t *testing.T) {

	type tierTest struct {
		name        string
		balance     int64
		nArguments  int
		earnedCoins int64
	}

	var tierTests = []tierTest{
		{"500 limit", 700, 11, 12},
		{"1000 limit", 1200, 21, 20},

		{"1500 limit", 1800, 31, 35},

		{"2000 limit", 2100, 41, 49},

		{"2500 limit", 5000, 51, 51},
	}
	assert.Len(t, tierTests, len(tierLimitsEarnedCoins))
	for _, tt := range tierTests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, k, mdb := mockDB()
			addr := createFakeFundedAccount(ctx, mdb.authAccKeeper, sdk.Coins{sdk.NewInt64Coin(app.StakeDenom, app.Shanev*tt.balance)})
			k.setEarnedCoins(ctx, addr, sdk.NewCoins(sdk.NewInt64Coin("crypto", app.Shanev*tt.earnedCoins)))
			argumentsToBeCreated := tt.nArguments
			for i := 1; i < tt.nArguments; i++ {
				_, err := k.SubmitArgument(ctx, "arg1", "summary1", addr, uint64(i), StakeChallenge)
				assert.NoError(t, err)
				argumentsToBeCreated--
			}
			_, err := k.SubmitArgument(ctx, "arg1", "summary1", addr, uint64(tt.nArguments), StakeChallenge)
			argumentsToBeCreated--
			assert.Error(t, err)
			assert.Equal(t, ErrorCodeMaxAmountStakingReached, err.Code())
			assert.Zero(t, argumentsToBeCreated)

		})
	}

}
func TestKeeper_StakeLimitDefaultTier(t *testing.T) {
	ctx, k, mdb := mockDB()
	addr := createFakeFundedAccount(ctx, mdb.authAccKeeper, sdk.Coins{sdk.NewInt64Coin(app.StakeDenom, app.Shanev*500)})

	_, err := k.SubmitArgument(ctx, "arg1", "summary1", addr, 1, StakeChallenge)
	assert.NoError(t, err)

	_, err = k.SubmitArgument(ctx, "arg1", "summary1", addr, 2, StakeChallenge)
	assert.NoError(t, err)

	_, err = k.SubmitArgument(ctx, "arg1", "summary1", addr, 3, StakeChallenge)
	assert.NoError(t, err)

	_, err = k.SubmitArgument(ctx, "arg1", "summary1", addr, 4, StakeChallenge)
	assert.NoError(t, err)

	_, err = k.SubmitArgument(ctx, "arg1", "summary1", addr, 5, StakeChallenge)
	assert.NoError(t, err)

	_, err = k.SubmitArgument(ctx, "arg1", "summary1", addr, 6, StakeChallenge)
	assert.NoError(t, err)

	_, err = k.SubmitArgument(ctx, "arg1", "summary1", addr, 7, StakeChallenge)
	assert.Error(t, err)
	assert.Equal(t, ErrorCodeMaxAmountStakingReached, err.Code())
}

func TestKeeper_StakeMinBalance(t *testing.T) {
	ctx, k, mdb := mockDB()
	addr := createFakeFundedAccount(ctx, mdb.authAccKeeper, sdk.Coins{sdk.NewInt64Coin(app.StakeDenom, app.Shanev*300)})

	_, err := k.SubmitArgument(ctx, "arg1", "summary1", addr, 1, StakeChallenge)
	assert.NoError(t, err)

	_, err = k.SubmitArgument(ctx, "arg1", "summary1", addr, 2, StakeChallenge)
	assert.NoError(t, err)

	_, err = k.SubmitArgument(ctx, "arg1", "summary1", addr, 3, StakeChallenge)
	assert.NoError(t, err)

	_, err = k.SubmitArgument(ctx, "arg1", "summary1", addr, 4, StakeChallenge)
	assert.NoError(t, err)

	_, err = k.SubmitArgument(ctx, "arg1", "summary1", addr, 5, StakeChallenge)
	assert.NoError(t, err)

	_, err = k.SubmitArgument(ctx, "arg1", "summary1", addr, 6, StakeChallenge)
	assert.Error(t, err)
	assert.Equal(t, ErrorCodeMinBalance, err.Code())

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

func TestKeeper_ClaimTotalsAdded(t *testing.T) {
	ctx, k, mdb := mockDB()
	mockedClaimKeeper := mdb.claimKeeper.(*mockClaimKeeper)
	mockedClaimKeeper.enableTrackStake = true
	claims := make(map[uint64]claim.Claim)
	addr := createFakeFundedAccount(ctx, mdb.authAccKeeper, sdk.Coins{sdk.NewInt64Coin(app.StakeDenom, app.Shanev*250)})
	addr2 := createFakeFundedAccount(ctx, mdb.authAccKeeper, sdk.Coins{sdk.NewInt64Coin(app.StakeDenom, app.Shanev*250)})
	addr3 := createFakeFundedAccount(ctx, mdb.authAccKeeper, sdk.Coins{sdk.NewInt64Coin(app.StakeDenom, app.Shanev*250)})
	claims[1] = claim.Claim{
		ID:              1,
		CommunityID:     "crypto",
		Body:            "body",
		Creator:         addr,
		TotalBacked:     sdk.NewInt64Coin(app.StakeDenom, 0),
		TotalChallenged: sdk.NewInt64Coin(app.StakeDenom, 0),
	}
	claims[2] = claim.Claim{
		ID:              2,
		CommunityID:     "random",
		Body:            "body",
		Creator:         addr,
		TotalBacked:     sdk.NewInt64Coin(app.StakeDenom, 0),
		TotalChallenged: sdk.NewInt64Coin(app.StakeDenom, 0),
	}
	mockedClaimKeeper.SetClaims(claims)

	_, err := k.SubmitArgument(ctx, "arg1", "summary1", addr, 1, StakeChallenge)
	assert.NoError(t, err)

	arg2, err := k.SubmitArgument(ctx, "arg2", "summary2", addr2, 2, StakeBacking)
	assert.NoError(t, err)
	_, err = k.SubmitUpvote(ctx, arg2.ID, addr)
	_, err = k.SubmitUpvote(ctx, arg2.ID, addr3)
	assert.NoError(t, err)

	arg3, err := k.SubmitArgument(ctx,
		"arg3", "summary3", addr, 2, StakeChallenge)
	_, err = k.SubmitUpvote(ctx, arg3.ID, addr3)
	claim1, ok := mockedClaimKeeper.Claim(ctx, 1)
	assert.True(t, ok)
	assert.Equal(t, claim1.TotalChallenged, sdk.NewInt64Coin(app.StakeDenom, app.Shanev*50))

	claim2, ok := mockedClaimKeeper.Claim(ctx, 2)
	assert.Equal(t, sdk.NewInt64Coin(app.StakeDenom, app.Shanev*60).String(), claim2.TotalChallenged.String())
	assert.Equal(t, sdk.NewInt64Coin(app.StakeDenom, app.Shanev*70).String(), claim2.TotalBacked.String())
	assert.True(t, ok)
	assert.NoError(t, err)
}

func TestAddAdmin_Success(t *testing.T) {
	ctx, keeper, _ := mockDB()

	creator := keeper.GetParams(ctx).StakingAdmins[0]
	_, _, newAdmin := keyPubAddr()

	err := keeper.AddAdmin(ctx, newAdmin, creator)
	assert.Nil(t, err)

	newAdmins := keeper.GetParams(ctx).StakingAdmins
	assert.Subset(t, newAdmins, []sdk.AccAddress{newAdmin})
}

func TestAddAdmin_CreatorNotAuthorised(t *testing.T) {
	ctx, keeper, _ := mockDB()

	invalidCreator := sdk.AccAddress([]byte{1, 2})
	_, _, newAdmin := keyPubAddr()

	err := keeper.AddAdmin(ctx, newAdmin, invalidCreator)
	assert.NotNil(t, err)
	assert.Equal(t, ErrAddressNotAuthorised().Code(), err.Code())
}

func TestRemoveAdmin_Success(t *testing.T) {
	ctx, keeper, _ := mockDB()

	currentAdmins := keeper.GetParams(ctx).StakingAdmins
	adminToRemove := currentAdmins[0]

	err := keeper.RemoveAdmin(ctx, adminToRemove, adminToRemove) // removing self
	assert.Nil(t, err)
	newAdmins := keeper.GetParams(ctx).StakingAdmins
	assert.Equal(t, len(currentAdmins)-1, len(newAdmins))
}

func TestRemoveAdmin_RemoverNotAuthorised(t *testing.T) {
	ctx, keeper, _ := mockDB()

	invalidRemover := sdk.AccAddress([]byte{1, 2})
	currentAdmins := keeper.GetParams(ctx).StakingAdmins
	adminToRemove := currentAdmins[0]

	err := keeper.AddAdmin(ctx, adminToRemove, invalidRemover)
	assert.NotNil(t, err)
	assert.Equal(t, ErrAddressNotAuthorised().Code(), err.Code())
}
