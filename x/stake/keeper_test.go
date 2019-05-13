package stake

import (
	"testing"
	"time"

	"github.com/TruStory/truchain/x/trubank"
	abci "github.com/tendermint/tendermint/abci/types"

	app "github.com/TruStory/truchain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

type FakeStaker struct {
	*Vote `json:"vote"`
}

// ID implements `Voter.ID`
func (s FakeStaker) ID() int64 {
	return s.Vote.ID
}

// StoryID implements `Voter.StoryID`
func (s FakeStaker) StoryID() int64 {
	return s.Vote.StoryID
}

// Amount implements `Voter.Amount`
func (s FakeStaker) Amount() sdk.Coin {
	return s.Vote.Amount
}

// Creator implements `Voter.Creator`
func (s FakeStaker) Creator() sdk.AccAddress {
	return s.Vote.Creator
}

// VoteChoice implements `Voter.VoteChoice`
func (s FakeStaker) VoteChoice() bool {
	return s.Vote.Vote
}

// Timestamp implements `Voter.Timestamp`
func (s FakeStaker) Timestamp() app.Timestamp {
	return s.Vote.Timestamp
}

func TestRedistributeStakeTrueWins(t *testing.T) {
	ctx, k := mockDB()

	amount := sdk.NewCoin("trusteak", sdk.NewInt(10*app.Shanev))
	creator1 := sdk.AccAddress([]byte{1, 2})
	creator2 := sdk.AccAddress([]byte{2, 4})
	trueVote := Vote{
		ID:         1,
		StoryID:    1,
		Amount:     amount,
		ArgumentID: 0,
		Creator:    creator1,
		Vote:       true,
		Timestamp:  app.NewTimestamp(ctx.BlockHeader()),
	}
	backingVote := FakeStaker{&trueVote}

	falseVote := Vote{
		ID:         2,
		StoryID:    1,
		Amount:     amount,
		ArgumentID: 0,
		Creator:    creator2,
		Vote:       false,
		Timestamp:  app.NewTimestamp(ctx.BlockHeader()),
	}
	challengeVote := FakeStaker{&falseVote}

	votes := []Voter{backingVote, backingVote, backingVote, challengeVote}

	_, err := k.RedistributeStake(ctx, votes)
	assert.NoError(t, err)

	// check stake amounts
	transactions, err := k.truBankKeeper.TransactionsByCreator(ctx, creator1)
	assert.NoError(t, err)
	// 3 backers should receive stake * 2 transaction each (refund + reward)
	assert.Len(t, transactions, 6)
	assert.Equal(t, trubank.BackingReturned, transactions[0].TransactionType)
	assert.Equal(t, trubank.RewardPool, transactions[1].TransactionType)
	assert.Equal(t, trubank.BackingReturned, transactions[2].TransactionType)
	assert.Equal(t, trubank.RewardPool, transactions[3].TransactionType)
	assert.Equal(t, trubank.BackingReturned, transactions[4].TransactionType)
	assert.Equal(t, trubank.RewardPool, transactions[5].TransactionType)
	assert.Equal(t, "10000000000trusteak", transactions[0].Amount.String())
	assert.Equal(t, "3333333333trusteak", transactions[1].Amount.String())
	assert.Equal(t, "10000000000trusteak", transactions[2].Amount.String())
	assert.Equal(t, "3333333333trusteak", transactions[3].Amount.String())
	assert.Equal(t, "10000000000trusteak", transactions[4].Amount.String())
	assert.Equal(t, "3333333333trusteak", transactions[5].Amount.String())

	// false voters get nothing back
	transactions, err = k.truBankKeeper.TransactionsByCreator(ctx, creator2)
	assert.NoError(t, err)
	assert.Len(t, transactions, 0)

}

func TestRedistributeStakeOneStaker(t *testing.T) {
	ctx, k := mockDB()

	amount := sdk.NewCoin("trusteak", sdk.NewInt(10*app.Shanev))
	creator1 := sdk.AccAddress([]byte{1, 2})
	trueVote := Vote{
		ID:         1,
		StoryID:    1,
		Amount:     amount,
		ArgumentID: 0,
		Creator:    creator1,
		Vote:       true,
		Timestamp:  app.NewTimestamp(ctx.BlockHeader()),
	}
	backingVote := FakeStaker{&trueVote}

	votes := []Voter{backingVote}

	_, err := k.RedistributeStake(ctx, votes)
	assert.NoError(t, err)

	// check stake amounts
	transactions, err := k.truBankKeeper.TransactionsByCreator(ctx, creator1)
	assert.NoError(t, err)
	// Stake should be returned
	assert.Len(t, transactions, 1)
	assert.Equal(t, trubank.BackingReturned, transactions[0].TransactionType)
	assert.Equal(t, "10000000000trusteak", transactions[0].Amount.String())

}

func TestRedistributeStakeFalseWins(t *testing.T) {
	ctx, k := mockDB()

	amount := sdk.NewCoin("trusteak", sdk.NewInt(10*app.Shanev))
	creator1 := sdk.AccAddress([]byte{1, 2})
	creator2 := sdk.AccAddress([]byte{2, 4})
	trueVote := Vote{
		ID:         1,
		StoryID:    1,
		Amount:     amount,
		ArgumentID: 0,
		Creator:    creator1,
		Vote:       true,
		Timestamp:  app.NewTimestamp(ctx.BlockHeader()),
	}
	backingVote := FakeStaker{&trueVote}

	falseVote := Vote{
		ID:         2,
		StoryID:    1,
		Amount:     amount,
		ArgumentID: 0,
		Creator:    creator2,
		Vote:       false,
		Timestamp:  app.NewTimestamp(ctx.BlockHeader()),
	}
	challengeVote := FakeStaker{&falseVote}

	votes := []Voter{backingVote, challengeVote, challengeVote, challengeVote}

	_, err := k.RedistributeStake(ctx, votes)
	assert.NoError(t, err)

	transactions, err := k.truBankKeeper.TransactionsByCreator(ctx, creator1)
	assert.NoError(t, err)
	// true voters don't get anything back
	assert.Len(t, transactions, 0)

	transactions, err = k.truBankKeeper.TransactionsByCreator(ctx, creator2)
	assert.NoError(t, err)
	// 3 challengers should receive stake * 2 transaction each (refund + reward)
	assert.Len(t, transactions, 6)
	assert.Equal(t, trubank.ChallengeReturned, transactions[0].TransactionType)
	assert.Equal(t, trubank.RewardPool, transactions[1].TransactionType)
	assert.Equal(t, trubank.ChallengeReturned, transactions[2].TransactionType)
	assert.Equal(t, trubank.RewardPool, transactions[3].TransactionType)
	assert.Equal(t, trubank.ChallengeReturned, transactions[4].TransactionType)
	assert.Equal(t, trubank.RewardPool, transactions[5].TransactionType)
	assert.Equal(t, "10000000000trusteak", transactions[0].Amount.String())
	assert.Equal(t, "3333333333trusteak", transactions[1].Amount.String())
	assert.Equal(t, "10000000000trusteak", transactions[2].Amount.String())
	assert.Equal(t, "3333333333trusteak", transactions[3].Amount.String())
	assert.Equal(t, "10000000000trusteak", transactions[4].Amount.String())
	assert.Equal(t, "3333333333trusteak", transactions[5].Amount.String())
}

func TestRedistributeStakeNoMajority(t *testing.T) {
	ctx, k := mockDB()

	amount := sdk.NewCoin("trusteak", sdk.NewInt(10*app.Shanev))
	creator1 := sdk.AccAddress([]byte{1, 2})
	creator2 := sdk.AccAddress([]byte{2, 4})
	trueVote := Vote{
		ID:         1,
		StoryID:    1,
		Amount:     amount,
		ArgumentID: 0,
		Creator:    creator1,
		Vote:       true,
		Timestamp:  app.NewTimestamp(ctx.BlockHeader()),
	}
	backingVote := FakeStaker{&trueVote}

	falseVote := Vote{
		ID:         2,
		StoryID:    1,
		Amount:     amount,
		ArgumentID: 0,
		Creator:    creator2,
		Vote:       false,
		Timestamp:  app.NewTimestamp(ctx.BlockHeader()),
	}
	challengeVote := FakeStaker{&falseVote}

	votes := []Voter{backingVote, challengeVote}

	_, err := k.RedistributeStake(ctx, votes)
	assert.NoError(t, err)

	transactions, err := k.truBankKeeper.TransactionsByCreator(ctx, creator1)
	assert.NoError(t, err)
	assert.Len(t, transactions, 1)
	assert.Equal(t, trubank.BackingReturned, transactions[0].TransactionType)
	assert.Equal(t, "10000000000trusteak", transactions[0].Amount.String())

	transactions, err = k.truBankKeeper.TransactionsByCreator(ctx, creator2)
	assert.NoError(t, err)
	assert.Len(t, transactions, 1)
	assert.Equal(t, trubank.ChallengeReturned, transactions[0].TransactionType)
	assert.Equal(t, "10000000000trusteak", transactions[0].Amount.String())
}

func TestDistributeInterest(t *testing.T) {
	ctx, k := mockDB()

	blockTime := time.Now()
	ctx = ctx.WithBlockHeader(abci.Header{Time: blockTime})

	// 10,000,000,000 amount
	//     20,547,945 interest ✅
	amount := sdk.NewCoin("trusteak", sdk.NewInt(10*app.Shanev))
	creator1 := sdk.AccAddress([]byte{1, 2})
	creator2 := sdk.AccAddress([]byte{2, 4})
	voteEndTime := blockTime.Add(-time.Hour * 24 * 3)

	trueVote := Vote{
		ID:         1,
		StoryID:    1,
		Amount:     amount,
		ArgumentID: 0,
		Creator:    creator1,
		Vote:       true,
		Timestamp:  app.NewTimestamp(abci.Header{Time: voteEndTime}),
	}
	backingVote := FakeStaker{&trueVote}

	falseVote := Vote{
		ID:         2,
		StoryID:    1,
		Amount:     amount,
		ArgumentID: 0,
		Creator:    creator2,
		Vote:       false,
		Timestamp:  app.NewTimestamp(abci.Header{Time: voteEndTime}),
	}
	challengeVote := FakeStaker{&falseVote}

	votes := []Voter{backingVote, backingVote, backingVote, challengeVote}

	_, err := k.DistributeInterest(ctx, votes)
	assert.NoError(t, err)

	transactions, err := k.truBankKeeper.TransactionsByCreator(ctx, creator1)
	assert.NoError(t, err)
	assert.Len(t, transactions, 3)
	assert.Equal(t, trubank.Interest, transactions[0].TransactionType)
	assert.Equal(t, trubank.Interest, transactions[1].TransactionType)
	assert.Equal(t, trubank.Interest, transactions[2].TransactionType)
	assert.Equal(t, "20547945trusteak", transactions[0].Amount.String())
	assert.Equal(t, "20547945trusteak", transactions[1].Amount.String())
	assert.Equal(t, "20547945trusteak", transactions[2].Amount.String())

	transactions, err = k.truBankKeeper.TransactionsByCreator(ctx, creator2)
	assert.NoError(t, err)
	assert.Len(t, transactions, 1)
	assert.Equal(t, trubank.Interest, transactions[0].TransactionType)
	assert.Equal(t, "20547945trusteak", transactions[0].Amount.String())
}

func Test_interest_3days(t *testing.T) {
	ctx, k := mockDB()

	amount := sdk.NewCoin("crypto", sdk.NewInt(500000000000000))
	period := time.Hour * 24 * 3

	interest := k.interest(ctx, amount, period)

	// 500,000,000,000,000 amount
	//   1,027,397,260,274 interest
	assert.Equal(t, sdk.NewInt(1027397260274).String(), interest.String())
}

func Test_interest_0days(t *testing.T) {
	ctx, k := mockDB()

	amount := sdk.NewCoin("crypto", sdk.NewInt(500000000000000))
	period := time.Hour * 24 * 0

	interest := k.interest(ctx, amount, period)

	// 500,000,000,000,000 amount
	//                   0 interest
	assert.Equal(t, sdk.NewInt(0).String(), interest.String())
}
