package vote

import (
	"crypto/rand"
	"net/url"
	"testing"
	"time"

	app "github.com/TruStory/truchain/types"
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

	// give user some funds
	amount := sdk.NewCoin("trudex", sdk.NewInt(20))
	k.AddCoins(ctx, creator, sdk.Coins{amount})

	return creator
}

func fakeValidationGame() (
	ctx sdk.Context, trueVotes []interface{}, falseVotes []interface{}, k Keeper) {

	ctx, k, ck := mockDB()

	storyID := createFakeStory(ctx, k.storyKeeper, ck)
	amount := sdk.NewCoin("trudex", sdk.NewInt(10))
	argument := "test argument"
	cnn, _ := url.Parse("http://www.cnn.com")
	evidence := []url.URL{*cnn}

	creator1 := fakeFundedCreator(ctx, k.bankKeeper)
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
	b1id, _ := k.backingKeeper.Create(ctx, storyID, amount, creator1, duration)
	b2id, _ := k.backingKeeper.Create(ctx, storyID, amount, creator2, duration)
	b3id, _ := k.backingKeeper.Create(ctx, storyID, amount, creator3, duration)
	b4id, _ := k.backingKeeper.Create(ctx, storyID, amount, creator4, duration)

	// fake challenges
	c1id, _ := k.challengeKeeper.Create(ctx, storyID, amount, argument, creator5, evidence)
	c2id, _ := k.challengeKeeper.Create(ctx, storyID, amount, argument, creator6, evidence)

	// fake votes
	v1id, _ := k.Create(ctx, storyID, amount, true, argument, creator7, evidence)
	v2id, _ := k.Create(ctx, storyID, amount, true, argument, creator8, evidence)
	v3id, _ := k.Create(ctx, storyID, amount, false, argument, creator9, evidence)

	b1, _ := k.backingKeeper.Backing(ctx, b1id)
	// fake an interest
	b1.Interest = sdk.NewCoin("trudex", sdk.NewInt(5))
	k.backingKeeper.Update(ctx, b1)

	b2, _ := k.backingKeeper.Backing(ctx, b2id)
	b2.Interest = sdk.NewCoin("trudex", sdk.NewInt(5))
	k.backingKeeper.Update(ctx, b2)

	b3, _ := k.backingKeeper.Backing(ctx, b3id)
	b3.Interest = sdk.NewCoin("trudex", sdk.NewInt(5))
	k.backingKeeper.Update(ctx, b3)

	b4, _ := k.backingKeeper.Backing(ctx, b4id)
	// change last backing vote to FALSE
	b4.Vote.Vote = false
	b4.Interest = sdk.NewCoin("trudex", sdk.NewInt(5))
	k.backingKeeper.Update(ctx, b4)

	c1, _ := k.challengeKeeper.Challenge(ctx, c1id)
	c2, _ := k.challengeKeeper.Challenge(ctx, c2id)

	v1, _ := k.Vote(ctx, v1id)
	v2, _ := k.Vote(ctx, v2id)
	v3, _ := k.Vote(ctx, v3id)

	trueVotes = append(trueVotes, b1, b2, b3, v1, v2)
	falseVotes = append(falseVotes, b4, c1, c2, v3)

	return
}

func TestProcessGame(t *testing.T) {
	ctx, _, _, k := fakeValidationGame()

	gameID := int64(1)
	game, _ := k.gameKeeper.Get(ctx, gameID)

	err := processGame(ctx, k, game)
	assert.Nil(t, err)
}

func TestTally(t *testing.T) {
	ctx, _, _, k := fakeValidationGame()

	gameID := int64(1)
	game, _ := k.gameKeeper.Get(ctx, gameID)

	trueVotes, falseVotes, _ := tally(ctx, k, game)

	assert.Equal(t, 5, len(trueVotes))
	assert.Equal(t, 4, len(falseVotes))
}

func TestRewardPool(t *testing.T) {
	ctx, trueVotes, falseVotes, _ := fakeValidationGame()

	expectedPool := sdk.NewCoin("trudex", sdk.NewInt(35))

	pool, _ := rewardPool(ctx, trueVotes, falseVotes, true)
	assert.Equal(t, expectedPool, pool)
}

func TestDistributeRewards(t *testing.T) {
	ctx, trueVotes, falseVotes, k := fakeValidationGame()

	pool := sdk.NewCoin("trudex", sdk.NewInt(35))

	err := distributeRewards(ctx, k.bankKeeper, pool, trueVotes, falseVotes, true)
	assert.Nil(t, err)
}

func TestConfirmStory(t *testing.T) {
	ctx, trueVotes, falseVotes, k := fakeValidationGame()

	confirmed, _ := confirmStory(ctx, k.accountKeeper, trueVotes, falseVotes)
	assert.True(t, confirmed)
}

func TestWeightedVote(t *testing.T) {
	ctx, trueVotes, falseVotes, k := fakeValidationGame()

	trueWeights, _ := weightedVote(ctx, k.accountKeeper, trueVotes)
	falseWeights, _ := weightedVote(ctx, k.accountKeeper, falseVotes)

	assert.Equal(t, sdk.NewInt(50), trueWeights)
	assert.Equal(t, sdk.NewInt(40), falseWeights)
}

func TestConfirmedStoryRewardPool(t *testing.T) {
	ctx, _, falseVotes, _ := fakeValidationGame()

	pool := sdk.NewCoin("trudex", sdk.ZeroInt())

	confirmedPool(ctx, falseVotes, &pool)
	assert.Equal(t, "35trudex", pool.String())
}

func TestDistributeRewardsConfirmed(t *testing.T) {
	ctx, trueVotes, falseVotes, k := fakeValidationGame()
	pool := sdk.NewCoin("trudex", sdk.ZeroInt())
	confirmedPool(ctx, falseVotes, &pool)

	err := distributeRewardsConfirmed(ctx, k.bankKeeper, trueVotes, falseVotes, pool)
	assert.Nil(t, err)

	coins := sdk.Coins{}

	winningBacker1 := trueVotes[0].(backing.Backing)
	coins = k.bankKeeper.GetCoins(ctx, winningBacker1.Creator)
	assert.Equal(t, "10", coins.AmountOf("trudex").String())

	winningBacker2 := trueVotes[1].(backing.Backing)
	coins = k.bankKeeper.GetCoins(ctx, winningBacker2.Creator)
	assert.Equal(t, "10", coins.AmountOf("trudex").String())

	winningBacker3 := trueVotes[2].(backing.Backing)
	coins = k.bankKeeper.GetCoins(ctx, winningBacker3.Creator)
	assert.Equal(t, "10", coins.AmountOf("trudex").String())

	winningVoter1 := trueVotes[3].(app.Vote)
	coins = k.bankKeeper.GetCoins(ctx, winningVoter1.Creator)
	assert.Equal(t, "37", coins.AmountOf("trudex").String())

	winningVoter2 := trueVotes[4].(app.Vote)
	coins = k.bankKeeper.GetCoins(ctx, winningVoter2.Creator)
	assert.Equal(t, "37", coins.AmountOf("trudex").String())

	losingBacker1 := falseVotes[0].(backing.Backing)
	coins = k.bankKeeper.GetCoins(ctx, losingBacker1.Creator)
	assert.Equal(t, "20", coins.AmountOf("trudex").String())

	losingChallenger1 := falseVotes[1].(challenge.Challenge)
	coins = k.bankKeeper.GetCoins(ctx, losingChallenger1.Creator)
	assert.Equal(t, "10", coins.AmountOf("trudex").String())

	losingChallenger2 := falseVotes[2].(challenge.Challenge)
	coins = k.bankKeeper.GetCoins(ctx, losingChallenger2.Creator)
	assert.Equal(t, "10", coins.AmountOf("trudex").String())

	losingVoter1 := falseVotes[3].(app.Vote)
	coins = k.bankKeeper.GetCoins(ctx, losingVoter1.Creator)
	assert.Equal(t, "10", coins.AmountOf("trudex").String())
}

func TestRejectedStoryRewardPool(t *testing.T) {
	ctx, trueVotes, falseVotes, _ := fakeValidationGame()

	pool := sdk.NewCoin("trudex", sdk.ZeroInt())

	rejectedPool(ctx, trueVotes, falseVotes, &pool)
	assert.Equal(t, "70trudex", pool.String())
}

func TestChallengerPool(t *testing.T) {
	ctx, trueVotes, falseVotes, _ := fakeValidationGame()
	pool := sdk.NewCoin("trudex", sdk.ZeroInt())
	rejectedPool(ctx, trueVotes, falseVotes, &pool)

	coin := challengerPool(pool, DefaultParams())
	assert.Equal(t, "52trudex", coin.String())
}

func TestVoterPool(t *testing.T) {
	ctx, trueVotes, falseVotes, _ := fakeValidationGame()
	pool := sdk.NewCoin("trudex", sdk.ZeroInt())
	rejectedPool(ctx, trueVotes, falseVotes, &pool)

	coin := voterPool(pool, DefaultParams())
	assert.Equal(t, "18trudex", coin.String())
}

func TestCount(t *testing.T) {
	ctx, trueVotes, falseVotes, _ := fakeValidationGame()
	pool := sdk.NewCoin("trudex", sdk.ZeroInt())
	rejectedPool(ctx, trueVotes, falseVotes, &pool)

	cCount, vCount, _ := count(falseVotes)
	assert.Equal(t, int64(2), cCount)
	assert.Equal(t, int64(1), vCount)
}

func TestChallengerRewardAmount(t *testing.T) {
	coin := challengerRewardAmount(
		sdk.NewCoin("trudex", sdk.NewInt(10)),
		int64(2),
		sdk.NewCoin("trudex", sdk.NewInt(52)))

	assert.Equal(t, "26", coin.String())
}

func TestDistributeRewardsRejected(t *testing.T) {
	ctx, trueVotes, falseVotes, k := fakeValidationGame()
	pool := sdk.NewCoin("trudex", sdk.ZeroInt())
	rejectedPool(ctx, trueVotes, falseVotes, &pool)

	err := distributeRewardsRejected(ctx, k.bankKeeper, falseVotes, pool)
	assert.Nil(t, err)

	coins := sdk.Coins{}

	winningBacker1 := falseVotes[0].(backing.Backing)
	coins = k.bankKeeper.GetCoins(ctx, winningBacker1.Creator)
	assert.Equal(t, "20", coins.AmountOf("trudex").String())

	winningChallenger1 := falseVotes[1].(challenge.Challenge)
	coins = k.bankKeeper.GetCoins(ctx, winningChallenger1.Creator)
	assert.Equal(t, "46", coins.AmountOf("trudex").String())

	winningChallenger2 := falseVotes[2].(challenge.Challenge)
	coins = k.bankKeeper.GetCoins(ctx, winningChallenger2.Creator)
	assert.Equal(t, "46", coins.AmountOf("trudex").String())

	winningVoter1 := falseVotes[3].(app.Vote)
	coins = k.bankKeeper.GetCoins(ctx, winningVoter1.Creator)
	assert.Equal(t, "38", coins.AmountOf("trudex").String())

	losingBacker1 := trueVotes[0].(backing.Backing)
	coins = k.bankKeeper.GetCoins(ctx, losingBacker1.Creator)
	assert.Equal(t, "10", coins.AmountOf("trudex").String())

	losingBacker2 := trueVotes[1].(backing.Backing)
	coins = k.bankKeeper.GetCoins(ctx, losingBacker2.Creator)
	assert.Equal(t, "10", coins.AmountOf("trudex").String())

	losingBacker3 := trueVotes[2].(backing.Backing)
	coins = k.bankKeeper.GetCoins(ctx, losingBacker3.Creator)
	assert.Equal(t, "10", coins.AmountOf("trudex").String())

	losingVoter1 := trueVotes[3].(app.Vote)
	coins = k.bankKeeper.GetCoins(ctx, losingVoter1.Creator)
	assert.Equal(t, "10", coins.AmountOf("trudex").String())

	losingVoter2 := trueVotes[4].(app.Vote)
	coins = k.bankKeeper.GetCoins(ctx, losingVoter2.Creator)
	assert.Equal(t, "10", coins.AmountOf("trudex").String())
}
