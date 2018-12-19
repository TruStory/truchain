package vote

import (
	"crypto/rand"
	"net/url"
	"testing"
	"time"

	"github.com/TruStory/truchain/x/backing"
	"github.com/TruStory/truchain/x/challenge"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/stretchr/testify/assert"
)

func fakeFundedCreator(ctx sdk.Context, k bank.Keeper) sdk.AccAddress {
	bz := make([]byte, 4)
	rand.Read(bz)
	creator := sdk.AccAddress(bz)

	// give user some category coins
	amount := sdk.NewCoin("trudex", sdk.NewInt(2000))
	k.AddCoins(ctx, creator, sdk.Coins{amount})

	return creator
}

func fakeValidationGame() (ctx sdk.Context, votes poll, k Keeper) {

	ctx, k, ck := mockDB()

	storyID := createFakeStory(ctx, k.storyKeeper, ck)
	amount := sdk.NewCoin("trudex", sdk.NewInt(1000))
	trustake := sdk.NewCoin("trusteak", sdk.NewInt(1000))
	argument := "test argument"
	cnn, _ := url.Parse("http://www.cnn.com")
	evidence := []url.URL{*cnn}

	creator1 := fakeFundedCreator(ctx, k.bankKeeper)
	// remove cat coins to simulate backing conversion from trusteak
	k.bankKeeper.SubtractCoins(ctx, creator1, sdk.Coins{amount})
	// add trustake
	k.bankKeeper.AddCoins(ctx, creator1, sdk.Coins{trustake})

	creator2 := fakeFundedCreator(ctx, k.bankKeeper)
	creator3 := fakeFundedCreator(ctx, k.bankKeeper)
	creator4 := fakeFundedCreator(ctx, k.bankKeeper)
	creator5 := fakeFundedCreator(ctx, k.bankKeeper)
	creator6 := fakeFundedCreator(ctx, k.bankKeeper)
	creator7 := fakeFundedCreator(ctx, k.bankKeeper)
	creator8 := fakeFundedCreator(ctx, k.bankKeeper)
	creator9 := fakeFundedCreator(ctx, k.bankKeeper)

	// fake backings
	duration := 1 * time.Hour
	b1id, _ := k.backingKeeper.Create(ctx, storyID, trustake, argument, creator1, duration)
	b2id, _ := k.backingKeeper.Create(ctx, storyID, amount, argument, creator2, duration)
	b3id, _ := k.backingKeeper.Create(ctx, storyID, amount, argument, creator3, duration)
	b4id, _ := k.backingKeeper.Create(ctx, storyID, amount, argument, creator4, duration)

	// fake challenges
	c1id, _ := k.challengeKeeper.Create(ctx, storyID, amount, argument, creator5, evidence)
	c2id, _ := k.challengeKeeper.Create(ctx, storyID, amount, argument, creator6, evidence)

	// fake votes
	v1id, _ := k.Create(ctx, storyID, amount, true, argument, creator7, evidence)
	v2id, _ := k.Create(ctx, storyID, amount, true, argument, creator8, evidence)
	v3id, _ := k.Create(ctx, storyID, amount, false, argument, creator9, evidence)

	b1, _ := k.backingKeeper.Backing(ctx, b1id)
	// fake an interest
	b1.Interest = sdk.NewCoin("trudex", sdk.NewInt(500))
	k.backingKeeper.Update(ctx, b1)

	b2, _ := k.backingKeeper.Backing(ctx, b2id)
	b2.Interest = sdk.NewCoin("trudex", sdk.NewInt(500))
	k.backingKeeper.Update(ctx, b2)

	b3, _ := k.backingKeeper.Backing(ctx, b3id)
	b3.Interest = sdk.NewCoin("trudex", sdk.NewInt(
		500))
	k.backingKeeper.Update(ctx, b3)

	b4, _ := k.backingKeeper.Backing(ctx, b4id)
	b4.Interest = sdk.NewCoin("trudex", sdk.NewInt(500))
	k.backingKeeper.Update(ctx, b4)
	// change last backing vote to FALSE
	k.backingKeeper.ToggleVote(ctx, b4.ID())

	c1, _ := k.challengeKeeper.Challenge(ctx, c1id)
	c2, _ := k.challengeKeeper.Challenge(ctx, c2id)

	v1, _ := k.TokenVote(ctx, v1id)
	v2, _ := k.TokenVote(ctx, v2id)
	v3, _ := k.TokenVote(ctx, v3id)

	votes.trueVotes = append(votes.trueVotes, b1, b2, b3, v1, v2)
	votes.falseVotes = append(votes.falseVotes, b4, c1, c2, v3)

	return
}

func TestProcessGame(t *testing.T) {
	ctx, _, k := fakeValidationGame()

	gameID := int64(1)
	game, _ := k.gameKeeper.Game(ctx, gameID)

	err := processGame(ctx, k, game)
	assert.Nil(t, err)
}

func TestTally(t *testing.T) {
	ctx, _, k := fakeValidationGame()

	gameID := int64(1)
	game, _ := k.gameKeeper.Game(ctx, gameID)

	votes, _ := tally(ctx, k, game)

	assert.Equal(t, 5, len(votes.trueVotes))
	assert.Equal(t, 4, len(votes.falseVotes))
}

func TestRewardPool(t *testing.T) {
	ctx, votes, _ := fakeValidationGame()

	expectedPool := sdk.NewCoin("trudex", sdk.NewInt(3500))

	pool, _ := rewardPool(ctx, votes, true)
	assert.Equal(t, expectedPool.String(), pool.String())
}

func TestDistributeRewards(t *testing.T) {
	ctx, votes, k := fakeValidationGame()

	pool := sdk.NewCoin("trudex", sdk.NewInt(35))

	err := distributeRewards(ctx, k.bankKeeper, pool, votes, true)
	assert.Nil(t, err)
}

func TestConfirmStory(t *testing.T) {
	ctx, votes, k := fakeValidationGame()

	confirmed, _ := confirmStory(ctx, k.accountKeeper, votes)
	assert.True(t, confirmed)
}

func TestWeightedVote(t *testing.T) {
	ctx, votes, k := fakeValidationGame()

	trueWeights, _ := weightedVote(ctx, k.accountKeeper, votes.trueVotes)
	falseWeights, _ := weightedVote(ctx, k.accountKeeper, votes.falseVotes)

	assert.Equal(t, "5000", trueWeights.String())
	assert.Equal(t, "4000", falseWeights.String())
}

func TestConfirmedStoryRewardPool(t *testing.T) {
	ctx, votes, _ := fakeValidationGame()

	pool := sdk.NewCoin("trudex", sdk.ZeroInt())

	confirmedPool(ctx, votes.falseVotes, &pool)
	assert.Equal(t, "3500trudex", pool.String())
}

func TestDistributeRewardsConfirmed(t *testing.T) {
	ctx, votes, k := fakeValidationGame()
	pool := sdk.NewCoin("trudex", sdk.ZeroInt())
	confirmedPool(ctx, votes.falseVotes, &pool)

	err := distributeRewardsConfirmed(ctx, k.bankKeeper, votes, pool)
	assert.Nil(t, err)

	coins := sdk.Coins{}

	winningBacker1 := votes.trueVotes[0].(backing.Backing)
	coins = k.bankKeeper.GetCoins(ctx, winningBacker1.Creator())
	assert.Equal(t, "1000", coins.AmountOf("trudex").String())

	winningBacker2 := votes.trueVotes[1].(backing.Backing)
	coins = k.bankKeeper.GetCoins(ctx, winningBacker2.Creator())
	assert.Equal(t, "1000", coins.AmountOf("trudex").String())

	winningBacker3 := votes.trueVotes[2].(backing.Backing)
	coins = k.bankKeeper.GetCoins(ctx, winningBacker3.Creator())
	assert.Equal(t, "1000", coins.AmountOf("trudex").String())

	winningVoter1 := votes.trueVotes[3].(TokenVote)
	coins = k.bankKeeper.GetCoins(ctx, winningVoter1.Creator())
	assert.Equal(t, "3750", coins.AmountOf("trudex").String())

	winningVoter2 := votes.trueVotes[4].(TokenVote)
	coins = k.bankKeeper.GetCoins(ctx, winningVoter2.Creator())
	assert.Equal(t, "3750", coins.AmountOf("trudex").String())

	losingBacker1 := votes.falseVotes[0].(backing.Backing)
	coins = k.bankKeeper.GetCoins(ctx, losingBacker1.Creator())
	assert.Equal(t, "2000", coins.AmountOf("trudex").String())

	losingChallenger1 := votes.falseVotes[1].(challenge.Challenge)
	coins = k.bankKeeper.GetCoins(ctx, losingChallenger1.Creator())
	assert.Equal(t, "1000", coins.AmountOf("trudex").String())

	losingChallenger2 := votes.falseVotes[2].(challenge.Challenge)
	coins = k.bankKeeper.GetCoins(ctx, losingChallenger2.Creator())
	assert.Equal(t, "1000", coins.AmountOf("trudex").String())

	losingVoter1 := votes.falseVotes[3].(TokenVote)
	coins = k.bankKeeper.GetCoins(ctx, losingVoter1.Creator())
	assert.Equal(t, "1000", coins.AmountOf("trudex").String())
}

func TestRejectedStoryRewardPool(t *testing.T) {
	ctx, votes, _ := fakeValidationGame()

	pool := sdk.NewCoin("trudex", sdk.ZeroInt())

	rejectedPool(ctx, votes, &pool)
	assert.Equal(t, "7000trudex", pool.String())
}

func TestChallengerPool(t *testing.T) {
	ctx, votes, _ := fakeValidationGame()
	pool := sdk.NewCoin("trudex", sdk.ZeroInt())
	rejectedPool(ctx, votes, &pool)

	coin := challengerPool(pool, DefaultParams())
	assert.Equal(t, "5250trudex", coin.String())
}

func TestVoterPool(t *testing.T) {
	ctx, votes, _ := fakeValidationGame()
	pool := sdk.NewCoin("trudex", sdk.ZeroInt())
	rejectedPool(ctx, votes, &pool)

	coin := voterPool(pool, DefaultParams())
	assert.Equal(t, "1750trudex", coin.String())
}

func TestCount(t *testing.T) {
	ctx, votes, _ := fakeValidationGame()
	pool := sdk.NewCoin("trudex", sdk.ZeroInt())
	rejectedPool(ctx, votes, &pool)

	cAmount, vCount, _ := winnerInfo(votes.falseVotes)
	assert.Equal(t, int64(2000), cAmount.Int64())
	assert.Equal(t, int64(1), vCount)
}

func TestChallengerRewardAmount(t *testing.T) {
	coin := challengerRewardAmount(
		sdk.NewCoin("trudex", sdk.NewInt(1000)),
		sdk.NewInt(2000),
		sdk.NewCoin("trudex", sdk.NewInt(5250)))

	assert.Equal(t, "2625", coin.String())
}

func TestChallengerRewardAmount2(t *testing.T) {
	coin := challengerRewardAmount(
		sdk.NewCoin("trudex", sdk.NewInt(1500)),
		sdk.NewInt(2000),
		sdk.NewCoin("trudex", sdk.NewInt(5250)))

	assert.Equal(t, "3938", coin.String())
}

func TestChallengerRewardAmount3(t *testing.T) {
	coin := challengerRewardAmount(
		sdk.NewCoin("trudex", sdk.NewInt(500)),
		sdk.NewInt(2000),
		sdk.NewCoin("trudex", sdk.NewInt(5250)))

	assert.Equal(t, "1312", coin.String())
}

func TestDistributeRewardsRejected(t *testing.T) {
	ctx, votes, k := fakeValidationGame()
	pool := sdk.NewCoin("trudex", sdk.ZeroInt())
	rejectedPool(ctx, votes, &pool)

	err := distributeRewardsRejected(ctx, k.bankKeeper, votes.falseVotes, pool)
	assert.Nil(t, err)

	coins := sdk.Coins{}

	winningBacker1 := votes.falseVotes[0].(backing.Backing)
	coins = k.bankKeeper.GetCoins(ctx, winningBacker1.Creator())
	assert.Equal(t, "2000", coins.AmountOf("trudex").String())

	winningChallenger1 := votes.falseVotes[1].(challenge.Challenge)
	coins = k.bankKeeper.GetCoins(ctx, winningChallenger1.Creator())
	assert.Equal(t, "4625", coins.AmountOf("trudex").String())

	winningChallenger2 := votes.falseVotes[2].(challenge.Challenge)
	coins = k.bankKeeper.GetCoins(ctx, winningChallenger2.Creator())
	assert.Equal(t, "4625", coins.AmountOf("trudex").String())

	winningVoter1 := votes.falseVotes[3].(TokenVote)
	coins = k.bankKeeper.GetCoins(ctx, winningVoter1.Creator())
	assert.Equal(t, "3750", coins.AmountOf("trudex").String())

	losingBacker1 := votes.trueVotes[0].(backing.Backing)
	coins = k.bankKeeper.GetCoins(ctx, losingBacker1.Creator())
	assert.Equal(t, "1000", coins.AmountOf("trudex").String())

	losingBacker2 := votes.trueVotes[1].(backing.Backing)
	coins = k.bankKeeper.GetCoins(ctx, losingBacker2.Creator())
	assert.Equal(t, "1000", coins.AmountOf("trudex").String())

	losingBacker3 := votes.trueVotes[2].(backing.Backing)
	coins = k.bankKeeper.GetCoins(ctx, losingBacker3.Creator())
	assert.Equal(t, "1000", coins.AmountOf("trudex").String())

	losingVoter1 := votes.trueVotes[3].(TokenVote)
	coins = k.bankKeeper.GetCoins(ctx, losingVoter1.Creator())
	assert.Equal(t, "1000", coins.AmountOf("trudex").String())

	losingVoter2 := votes.trueVotes[4].(TokenVote)
	coins = k.bankKeeper.GetCoins(ctx, losingVoter2.Creator())
	assert.Equal(t, "1000", coins.AmountOf("trudex").String())
}
