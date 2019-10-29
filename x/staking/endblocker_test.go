package staking

import (
	"testing"
	"time"

	app "github.com/TruStory/truchain/types"
	"github.com/TruStory/truchain/x/claim"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

func TestKeeper_TestEarnedCoins(t *testing.T) {
	ctx, k, mdb := mockDB()
	mockedClaimKeeper := mdb.claimKeeper.(*mockClaimKeeper)
	claims := make(map[uint64]claim.Claim)
	addr := createFakeFundedAccount(ctx, mdb.authAccKeeper, sdk.Coins{sdk.NewInt64Coin(app.StakeDenom, app.Shanev*250)})
	addr2 := createFakeFundedAccount(ctx, mdb.authAccKeeper, sdk.Coins{sdk.NewInt64Coin(app.StakeDenom, app.Shanev*250)})
	claims[1] = claim.Claim{
		ID:          1,
		CommunityID: "crypto",
		Body:        "body",
		Creator:     addr,
	}
	claims[2] = claim.Claim{
		ID:          2,
		CommunityID: "random",
		Body:        "body",
		Creator:     addr,
	}
	mockedClaimKeeper.SetClaims(claims)

	_, err := k.SubmitArgument(ctx.WithBlockTime(mustParseTime("2019-01-01")),
		"arg1", "summary1", addr, 1, StakeChallenge)
	assert.NoError(t, err)
	arg2, err := k.SubmitArgument(ctx.WithBlockTime(mustParseTime("2019-01-02")),
		"arg2", "summary2", addr2, 2, StakeBacking)
	assert.NoError(t, err)
	_, err = k.SubmitUpvote(ctx.WithBlockTime(mustParseTime("2019-01-03")), arg2.ID, addr)
	assert.NoError(t, err)

	EndBlocker(ctx.WithBlockTime(mustParseTime("2019-01-13")), k)
	usersEarnings := k.UsersEarnings(ctx)
	assert.Len(t, usersEarnings, 2)
	earnings := make(map[string]UserEarnedCoins)
	earnings[usersEarnings[0].Address.String()] = usersEarnings[0]
	earnings[usersEarnings[1].Address.String()] = usersEarnings[1]
	argumentInterest := k.interest(ctx, sdk.NewInt64Coin(app.StakeDenom, app.Shanev*50), time.Hour*24*7).RoundInt()
	upvoteInterest := k.interest(ctx, sdk.NewInt64Coin(app.StakeDenom, app.Shanev*10), time.Hour*24*7)
	upvoteAfterSplitInterest := upvoteInterest.Mul(sdk.NewDecWithPrec(50, 2)).RoundInt()

	assert.Equal(t, argumentInterest.String(), earnings[addr.String()].Coins.AmountOf("crypto").String())
	assert.Equal(t, upvoteAfterSplitInterest.String(), earnings[addr.String()].Coins.AmountOf("random").String())
	t.Log(argumentInterest.String())
	t.Log(upvoteInterest.String())
	t.Log(upvoteAfterSplitInterest.String())

	assert.Equal(t, sdk.NewInt(0), earnings[addr2.String()].Coins.AmountOf("crypto"))
	argumentAndUpvoteReceived := argumentInterest.Add(upvoteAfterSplitInterest)
	assert.Equal(t, argumentAndUpvoteReceived.String(), earnings[addr2.String()].Coins.AmountOf("random").String())
}

func TestKeeper_TestRefundStake(t *testing.T) {
	ctx, k, mdb := mockDB()
	mockedClaimKeeper := mdb.claimKeeper.(*mockClaimKeeper)
	claims := make(map[uint64]claim.Claim)
	addr := createFakeFundedAccount(ctx, mdb.authAccKeeper, sdk.Coins{sdk.NewInt64Coin(app.StakeDenom, app.Shanev*250)})
	addr2 := createFakeFundedAccount(ctx, mdb.authAccKeeper, sdk.Coins{sdk.NewInt64Coin(app.StakeDenom, app.Shanev*250)})
	claims[1] = claim.Claim{
		ID:          1,
		CommunityID: "crypto",
		Body:        "body",
		Creator:     addr,
	}
	claims[2] = claim.Claim{
		ID:          2,
		CommunityID: "random",
		Body:        "body",
		Creator:     addr,
	}
	mockedClaimKeeper.SetClaims(claims)

	_, err := k.SubmitArgument(ctx.WithBlockTime(mustParseTime("2019-01-01")),
		"arg1", "summary1", addr, 1, StakeChallenge)
	assert.NoError(t, err)
	arg2, err := k.SubmitArgument(ctx.WithBlockTime(mustParseTime("2019-01-03")),
		"arg2", "summary2", addr2, 2, StakeBacking)
	assert.NoError(t, err)
	_, err = k.SubmitUpvote(ctx.WithBlockTime(mustParseTime("2019-01-05")), arg2.ID, addr)
	assert.NoError(t, err)
	EndBlocker(ctx.WithBlockTime(mustParseTime("2019-01-08")), k)
	_, err = k.SubmitArgument(ctx.WithBlockTime(mustParseTime("2019-01-11")),
		"arg2", "summary2", addr, 2, StakeBacking)
	assert.NoError(t, err)
	EndBlocker(ctx.WithBlockTime(mustParseTime("2019-01-13")), k)
	addr1Txs := k.bankKeeper.TransactionsByAddress(ctx, addr)
	addr2Txs := k.bankKeeper.TransactionsByAddress(ctx, addr2)

	// 3 stakes + 2 interest + 2 refund
	assert.Len(t, addr1Txs, 7)
	txTypes := make([]TransactionType, 0)
	txCommunities := make([]string, 0)
	for _, tx := range addr1Txs {
		txTypes = append(txTypes, tx.Type)
		txCommunities = append(txCommunities, tx.CommunityID)
	}
	expected := []TransactionType{
		// first interactions
		TransactionChallenge, TransactionUpvote,
		// first end block
		TransactionChallengeReturned, TransactionInterestArgumentCreation,
		// second interactions
		TransactionBacking,
		// second end block
		TransactionUpvoteReturned, TransactionInterestUpvoteGiven}
	expectedTxCommunities := []string{"crypto", "random", "crypto", "crypto", "random", "random", "random"}
	assert.Equal(t, expected, txTypes)
	assert.Equal(t, expectedTxCommunities, txCommunities)

	// 1 stakes + 2 interest + 1 refund
	assert.Len(t, addr2Txs, 4)
	txTypes = make([]TransactionType, 0)

	for _, tx := range addr2Txs {

		txTypes = append(txTypes, tx.Type)
	}
	expected = []TransactionType{TransactionBacking, TransactionBackingReturned,
		TransactionInterestArgumentCreation, TransactionInterestUpvoteReceived}
	assert.Equal(t, expected, txTypes)
}

func TestKeeper_TestStakeRewardResult(t *testing.T) {
	ctx, k, mdb := mockDB()
	mockedClaimKeeper := mdb.claimKeeper.(*mockClaimKeeper)
	claims := make(map[uint64]claim.Claim)
	addr := createFakeFundedAccount(ctx, mdb.authAccKeeper, sdk.Coins{sdk.NewInt64Coin(app.StakeDenom, app.Shanev*250)})
	addr2 := createFakeFundedAccount(ctx, mdb.authAccKeeper, sdk.Coins{sdk.NewInt64Coin(app.StakeDenom, app.Shanev*250)})

	claims[1] = claim.Claim{
		ID:          1,
		CommunityID: "crypto",
		Body:        "body",
		Creator:     addr,
	}
	claims[2] = claim.Claim{
		ID:          2,
		CommunityID: "random",
		Body:        "body",
		Creator:     addr,
	}
	mockedClaimKeeper.SetClaims(claims)

	_, err := k.SubmitArgument(ctx.WithBlockTime(mustParseTime("2019-01-01")),
		"arg1", "summary1", addr, 1, StakeChallenge)
	assert.NoError(t, err)
	// check  the stake is added to the stake pool
	c := mdb.supplyKeeper.GetModuleAccount(ctx, UserStakesPoolName).GetCoins()
	assert.Equal(t, c, sdk.Coins{k.GetParams(ctx).ArgumentCreationStake})

	EndBlocker(ctx.WithBlockTime(mustParseTime("2019-01-08")), k)
	// check  the stake is removed from the pool
	c = mdb.supplyKeeper.GetModuleAccount(ctx, UserStakesPoolName).GetCoins()
	assert.Equal(t, c.AmountOf(app.StakeDenom).String(), "0")

	arg2, err := k.SubmitArgument(ctx.WithBlockTime(mustParseTime("2019-01-03")),
		"arg2", "summary2", addr2, 2, StakeBacking)
	assert.NoError(t, err)
	EndBlocker(ctx.WithBlockTime(mustParseTime("2019-01-11")), k)

	_, err = k.SubmitUpvote(ctx.WithBlockTime(mustParseTime("2019-01-05")), arg2.ID, addr)
	assert.NoError(t, err)
	// check  the stake is added to the stake pool
	c = mdb.supplyKeeper.GetModuleAccount(ctx, UserStakesPoolName).GetCoins()
	assert.Equal(t, c, sdk.Coins{k.GetParams(ctx).UpvoteStake})
	EndBlocker(ctx.WithBlockTime(mustParseTime("2019-01-13")), k)
	// check  the stake is removed from the pool
	c = mdb.supplyKeeper.GetModuleAccount(ctx, UserStakesPoolName).GetCoins()
	assert.Equal(t, c.AmountOf(app.StakeDenom).String(), "0")

	stakes := k.UserStakes(ctx, addr)

	assert.Len(t, stakes, 2)
	assert.NotNil(t, stakes[0].Result)
	assert.NotNil(t, stakes[1].Result)
	stakesUser2 := k.UserStakes(ctx, addr2)
	assert.Len(t, stakesUser2, 1)
	assert.NotNil(t, stakesUser2[0].Result)

	assert.Equal(t, RewardResultArgumentCreation, stakes[0].Result.Type)
	assert.Equal(t, RewardResultUpvoteSplit, stakes[1].Result.Type)
	assert.Equal(t, RewardResultArgumentCreation, stakesUser2[0].Result.Type)

	argumentInterest := k.interest(ctx, sdk.NewInt64Coin(app.StakeDenom, app.Shanev*50), time.Hour*24*7).RoundInt()
	upvoteInterest := k.interest(ctx, sdk.NewInt64Coin(app.StakeDenom, app.Shanev*10), time.Hour*24*7)
	upvoteAfterSplitInterest := upvoteInterest.Mul(sdk.NewDecWithPrec(50, 2)).RoundInt()

	assert.Equal(t, argumentInterest.String(), stakes[0].Result.ArgumentCreatorReward.Amount.String())
	assert.Equal(t, upvoteAfterSplitInterest.String(), stakes[1].Result.ArgumentCreatorReward.Amount.String())
	assert.Equal(t, upvoteAfterSplitInterest.String(), stakes[1].Result.StakeCreatorReward.Amount.String())
	assert.Equal(t, addr, stakes[1].Result.StakeCreator)
	assert.Equal(t, addr2, stakes[1].Result.ArgumentCreator)
	totalRewards := stakes[0].Result.ArgumentCreatorReward
	totalRewards = totalRewards.Add(stakes[1].Result.ArgumentCreatorReward)
	totalRewards = totalRewards.Add(stakes[1].Result.StakeCreatorReward)
	totalRewards = totalRewards.Add(stakesUser2[0].Result.ArgumentCreatorReward)
	coins, _ := sdk.ParseCoins("1000000000000tru")
	result := coins.Sub(sdk.Coins{totalRewards})
	c = mdb.supplyKeeper.GetModuleAccount(ctx, UserRewardPoolName).GetCoins()
	assert.Equal(t, c.AmountOf(app.StakeDenom).String(), result.AmountOf(app.StakeDenom).String())

}
