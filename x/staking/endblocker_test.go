package staking

import (
	"sort"
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
	EndBlocker(ctx.WithBlockTime(mustParseTime("2019-01-08")), k)
	arg2, err := k.SubmitArgument(ctx.WithBlockTime(mustParseTime("2019-01-03")),
		"arg2", "summary2", addr2, 2, StakeBacking)
	assert.NoError(t, err)
	EndBlocker(ctx.WithBlockTime(mustParseTime("2019-01-11")), k)
	_, err = k.SubmitUpvote(ctx.WithBlockTime(mustParseTime("2019-01-05")), arg2.ID, addr)
	assert.NoError(t, err)
	EndBlocker(ctx.WithBlockTime(mustParseTime("2019-01-13")), k)
	addr1Txs := k.bankKeeper.TransactionsByAddress(ctx, addr)
	addr2Txs := k.bankKeeper.TransactionsByAddress(ctx, addr2)
	sort.Slice(addr1Txs, func(i, j int) bool {
		return addr1Txs[i].ID < addr1Txs[j].ID
	})

	sort.Slice(addr2Txs, func(i, j int) bool {
		return addr2Txs[i].ID < addr2Txs[j].ID
	})
	// 2 stakes + 2 interest + 2 refund
	assert.Len(t, addr1Txs, 6)
	txTypes := make([]TransactionType, 0)

	for _, tx := range addr1Txs {

		txTypes = append(txTypes, tx.Type)
	}
	expected := []TransactionType{TransactionChallenge, TransactionChallengeReturned,
		TransactionInterestArgumentCreation, TransactionUpvote, TransactionUpvoteReturned, TransactionInterestUpvoteGiven}
	assert.Equal(t, expected, txTypes)

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
