package voting

import (
	"testing"
	"time"

	app "github.com/TruStory/truchain/types"
	"github.com/TruStory/truchain/x/backing"
	"github.com/TruStory/truchain/x/challenge"
	tokenVote "github.com/TruStory/truchain/x/vote"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	abci "github.com/tendermint/tendermint/abci/types"
)

// [CONFIRMED STORY] ==================================================
func TestDistributeRewards(t *testing.T) {
	ctx, votes, k := fakeConfirmedGame()

	pool := sdk.NewCoin(app.StakeDenom, sdk.NewInt(35))
	categoryID := int64(1)

	err := k.distributeRewards(ctx, pool, votes, true, categoryID)
	assert.Nil(t, err)
}

func TestDistributeRewardsNoStakers(t *testing.T) {
	ctx, votes, k := fakeConfirmedGameNoStakers()

	pool := sdk.NewCoin(app.StakeDenom, sdk.NewInt(35))
	categoryID := int64(1)

	err := k.distributeRewards(ctx, pool, votes, true, categoryID)
	assert.Nil(t, err)
}

func TestDistributeRewardsConfirmed(t *testing.T) {
	ctx, votes, k := fakeConfirmedGame()
	categoryID := int64(1)
	cred := "crypto"

	// fake future block time
	interestStopTime := ctx.BlockHeader().Time.Add(24 * time.Hour)
	ctx = ctx.WithBlockHeader(abci.Header{Time: interestStopTime})

	pool, _ := k.rewardPool(ctx, votes, true, categoryID)
	err := k.distributeRewards(ctx, pool, votes, true, categoryID)
	assert.Nil(t, err)

	coins := sdk.Coins{}

	winningBacker1 := votes.trueVotes[0].(backing.Backing)
	coins = k.bankKeeper.GetCoins(ctx, winningBacker1.Creator())
	assert.Equal(t, "2000000000000", coins.AmountOf(app.StakeDenom).String())
	assert.Equal(t, "1054295762500", coins.AmountOf(cred).String())

	winningBacker2 := votes.trueVotes[1].(backing.Backing)
	coins = k.bankKeeper.GetCoins(ctx, winningBacker2.Creator())
	assert.Equal(t, "2000000000000", coins.AmountOf(app.StakeDenom).String())
	assert.Equal(t, "1054295762500", coins.AmountOf(cred).String())

	winningBacker3 := votes.trueVotes[2].(backing.Backing)
	coins = k.bankKeeper.GetCoins(ctx, winningBacker3.Creator())
	assert.Equal(t, "2000000000000", coins.AmountOf(app.StakeDenom).String())
	assert.Equal(t, "1054295762500", coins.AmountOf(cred).String())

	winningBacker4 := votes.trueVotes[3].(backing.Backing)
	coins = k.bankKeeper.GetCoins(ctx, winningBacker4.Creator())
	assert.Equal(t, "2000000000000", coins.AmountOf(app.StakeDenom).String())
	assert.Equal(t, "1054295762500", coins.AmountOf(cred).String())

	winningVoter1 := votes.trueVotes[4].(tokenVote.TokenVote)
	coins = k.bankKeeper.GetCoins(ctx, winningVoter1.Creator())
	assert.Equal(t, "2000000000000", coins.AmountOf(app.StakeDenom).String())
	assert.Equal(t, "438916650000", coins.AmountOf(cred).String())

	winningVoter2 := votes.trueVotes[5].(tokenVote.TokenVote)
	coins = k.bankKeeper.GetCoins(ctx, winningVoter2.Creator())
	assert.Equal(t, "2000000000000", coins.AmountOf(app.StakeDenom).String())
	assert.Equal(t, "438916650000", coins.AmountOf(cred).String())

	winningVoter3 := votes.trueVotes[6].(tokenVote.TokenVote)
	coins = k.bankKeeper.GetCoins(ctx, winningVoter3.Creator())
	assert.Equal(t, "2000000000000", coins.AmountOf(app.StakeDenom).String())
	assert.Equal(t, "438916650000", coins.AmountOf(cred).String())

	losingChallenger1 := votes.falseVotes[0].(challenge.Challenge)
	coins = k.bankKeeper.GetCoins(ctx, losingChallenger1.Creator())
	assert.Equal(t, "1000000000000", coins.AmountOf(app.StakeDenom).String())
	assert.Equal(t, "0", coins.AmountOf(cred).String())

	losingChallenger2 := votes.falseVotes[1].(challenge.Challenge)
	coins = k.bankKeeper.GetCoins(ctx, losingChallenger2.Creator())
	assert.Equal(t, "1000000000000", coins.AmountOf(app.StakeDenom).String())
	assert.Equal(t, "0", coins.AmountOf(cred).String())

	losingChallenger3 := votes.falseVotes[2].(challenge.Challenge)
	coins = k.bankKeeper.GetCoins(ctx, losingChallenger3.Creator())
	assert.Equal(t, "0", coins.AmountOf(app.StakeDenom).String())
	assert.Equal(t, "0", coins.AmountOf(cred).String())

	losingVoter1 := votes.falseVotes[3].(tokenVote.TokenVote)
	coins = k.bankKeeper.GetCoins(ctx, losingVoter1.Creator())
	assert.Equal(t, "1000000000000", coins.AmountOf(app.StakeDenom).String())
	assert.Equal(t, "0", coins.AmountOf(cred).String())
}

func TestCount(t *testing.T) {
	_, votes, _ := fakeConfirmedGame()

	cAmount, cCount, vCount, _ := winnerInfo(votes.falseVotes)
	assert.Equal(t, int64(4000000000000), cAmount.Int64())
	assert.Equal(t, int64(3), cCount)
	assert.Equal(t, int64(1), vCount)
}

// [REJECTED STORY] ==================================================

func TestDistributeRewardsRejected(t *testing.T) {
	ctx, votes, k := fakeRejectedGame()
	categoryID := int64(1)
	cred := "crypto"

	// fake future block time
	interestStopTime := ctx.BlockHeader().Time.Add(24 * time.Hour)
	ctx = ctx.WithBlockHeader(abci.Header{Time: interestStopTime})

	pool, _ := k.rewardPool(ctx, votes, false, categoryID)
	err := k.distributeRewards(ctx, pool, votes, false, categoryID)
	assert.Nil(t, err)

	coins := sdk.Coins{}

	winningChallenger1 := votes.falseVotes[0].(challenge.Challenge)
	coins = k.bankKeeper.GetCoins(ctx, winningChallenger1.Creator())
	assert.Equal(t, "2000000000000", coins.AmountOf(app.StakeDenom).String())
	assert.Equal(t, "66733300000", coins.AmountOf(cred).String())
}
