package vote

import (
	"crypto/rand"
	"net/url"
	"testing"
	"time"

	"github.com/TruStory/truchain/x/backing"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/stretchr/testify/assert"
)

func fakeFundedCreator(ctx sdk.Context, k bank.Keeper) sdk.AccAddress {
	bz := make([]byte, 4)
	rand.Read(bz)
	creator := sdk.AccAddress(bz)

	// give user some funds
	amount := sdk.NewCoin("trudex", sdk.NewInt(15))
	k.AddCoins(ctx, creator, sdk.Coins{amount})

	return creator
}

func createFakeConfirmedStory() (
	ctx sdk.Context, trueVotes []interface{}, falseVotes []interface{},
	bankKeeper bank.Keeper) {

	ctx, k, sk, ck, challengeKeeper, bankKeeper, backingKeeper := mockDB()

	storyID := createFakeStory(ctx, sk, ck)
	amount := sdk.NewCoin("trudex", sdk.NewInt(10))
	argument := "test argument"
	cnn, _ := url.Parse("http://www.cnn.com")
	evidence := []url.URL{*cnn}

	creator1 := fakeFundedCreator(ctx, bankKeeper)
	creator2 := fakeFundedCreator(ctx, bankKeeper)
	creator3 := fakeFundedCreator(ctx, bankKeeper)
	creator4 := fakeFundedCreator(ctx, bankKeeper)
	creator5 := fakeFundedCreator(ctx, bankKeeper)
	creator6 := fakeFundedCreator(ctx, bankKeeper)
	creator7 := fakeFundedCreator(ctx, bankKeeper)
	creator8 := fakeFundedCreator(ctx, bankKeeper)
	creator9 := fakeFundedCreator(ctx, bankKeeper)

	// fake backings
	duration := 1 * time.Hour
	b1id, _ := backingKeeper.Create(ctx, storyID, amount, creator1, duration)
	b2id, _ := backingKeeper.Create(ctx, storyID, amount, creator2, duration)
	b3id, _ := backingKeeper.Create(ctx, storyID, amount, creator3, duration)
	b4id, _ := backingKeeper.Create(ctx, storyID, amount, creator4, duration)

	// fake challenges
	c1id, _ := challengeKeeper.Create(ctx, storyID, amount, argument, creator5, evidence)
	c2id, _ := challengeKeeper.Create(ctx, storyID, amount, argument, creator6, evidence)

	// fake votes
	v1id, _ := k.Create(ctx, storyID, amount, true, argument, creator7, evidence)
	v2id, _ := k.Create(ctx, storyID, amount, true, argument, creator8, evidence)
	v3id, _ := k.Create(ctx, storyID, amount, false, argument, creator9, evidence)

	b1, _ := backingKeeper.Backing(ctx, b1id)
	// fake an interest
	b1.Interest = sdk.NewCoin("trudex", sdk.NewInt(5))
	backingKeeper.Update(ctx, b1)

	b2, _ := backingKeeper.Backing(ctx, b2id)
	b2.Interest = sdk.NewCoin("trudex", sdk.NewInt(5))
	backingKeeper.Update(ctx, b2)

	b3, _ := backingKeeper.Backing(ctx, b3id)
	b3.Interest = sdk.NewCoin("trudex", sdk.NewInt(5))
	backingKeeper.Update(ctx, b3)

	b4, _ := backingKeeper.Backing(ctx, b4id)
	// change last backing vote to FALSE
	b4.Vote.Vote = false
	b4.Interest = sdk.NewCoin("trudex", sdk.NewInt(5))
	backingKeeper.Update(ctx, b4)

	c1, _ := challengeKeeper.Challenge(ctx, c1id)
	c2, _ := challengeKeeper.Challenge(ctx, c2id)

	v1, _ := k.Vote(ctx, v1id)
	v2, _ := k.Vote(ctx, v2id)
	v3, _ := k.Vote(ctx, v3id)

	trueVotes = append(trueVotes, b1, b2, b3, v1, v2)
	falseVotes = append(falseVotes, b4, c1, c2, v3)

	return
}

func TestConfirmedStoryRewardPool(t *testing.T) {
	ctx, _, falseVotes, _ := createFakeConfirmedStory()

	pool := sdk.NewCoin("trudex", sdk.ZeroInt())

	confirmedPool(ctx, falseVotes, &pool)
	assert.Equal(t, "35trudex", pool.String())
}

func TestRejectedStoryRewardPool(t *testing.T) {
	ctx, trueVotes, falseVotes, _ := createFakeConfirmedStory()

	pool := sdk.NewCoin("trudex", sdk.ZeroInt())

	rejectedPool(ctx, trueVotes, falseVotes, &pool)
	assert.Equal(t, "70trudex", pool.String())
}

func TestDistributeRewardsConfirmed(t *testing.T) {
	ctx, trueVotes, falseVotes, bankKeeper := createFakeConfirmedStory()
	pool := sdk.NewCoin("trudex", sdk.ZeroInt())
	confirmedPool(ctx, falseVotes, &pool)

	err := distributeRewardsConfirmed(ctx, bankKeeper, trueVotes, falseVotes, pool)
	assert.Nil(t, err)

	coins := sdk.Coins{}

	winningBacker1 := trueVotes[0].(backing.Backing)
	coins = bankKeeper.GetCoins(ctx, winningBacker1.Creator)
	assert.Equal(t, "0", coins.AmountOf("trudex").String())

	// winningBacker2 := trueVotes[1].(backing.Backing)
	// coins = bankKeeper.GetCoins(ctx, winningBacker2.Creator)
	// assert.Equal(t, "0", coins.AmountOf("trudex").String())

	// winningBacker3 := trueVotes[2].(backing.Backing)
	// coins = bankKeeper.GetCoins(ctx, winningBacker3.Creator)
	// assert.Equal(t, "0", coins.AmountOf("trudex").String())

	// winningVoter1 := trueVotes[3].(app.Vote)
	// coins := bankKeeper.GetCoins(ctx, winningVoter1.Creator)
	// assert.Equal(t, "27.5", coins.AmountOf("trudex").String())

	// winningVoter2 := trueVotes[4].(app.Vote)
	// coins = bankKeeper.GetCoins(ctx, winningVoter2.Creator)
	// assert.Equal(t, "27.5", coins.AmountOf("trudex").String())

	// losingBacker1 := falseVotes[0].(backing.Backing)
	// coins = bankKeeper.GetCoins(ctx, losingBacker1.Creator)
	// assert.Equal(t, "10", coins.AmountOf("trudex").String())

	// losingChallenger1 := falseVotes[1].(challenge.Challenge)
	// coins = bankKeeper.GetCoins(ctx, losingChallenger1.Creator)
	// assert.Equal(t, "0", coins.AmountOf("trudex").String())

	// losingChallenger2 := falseVotes[2].(challenge.Challenge)
	// coins = bankKeeper.GetCoins(ctx, losingChallenger2.Creator)
	// assert.Equal(t, "0", coins.AmountOf("trudex").String())

	// losingVoter1 := falseVotes[3].(app.Vote)
	// coins = bankKeeper.GetCoins(ctx, losingVoter1.Creator)
	// assert.Equal(t, "0", coins.AmountOf("trudex").String())
}
