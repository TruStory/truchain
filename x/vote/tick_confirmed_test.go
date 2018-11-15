package vote

import (
	"crypto/rand"
	"net/url"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/davecgh/go-spew/spew"
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
	ctx sdk.Context, falseVotes []interface{}, trueVotes []interface{}) {

	ctx, k, sk, ck, challengeKeeper, bankKeeper, backingKeeper := mockDB()

	storyID := createFakeStory(ctx, sk, ck)
	amount := sdk.NewCoin("trudex", sdk.NewInt(15))
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
	b1id, _ := backingKeeper.NewBacking(ctx, storyID, amount, creator1, duration)
	b2id, _ := backingKeeper.NewBacking(ctx, storyID, amount, creator2, duration)
	b3id, _ := backingKeeper.NewBacking(ctx, storyID, amount, creator3, duration)
	b4id, _ := backingKeeper.NewBacking(ctx, storyID, amount, creator4, duration)

	// fake challenges
	c1id, _ := challengeKeeper.Create(ctx, storyID, amount, argument, creator5, evidence)
	c2id, _ := challengeKeeper.Create(ctx, storyID, amount, argument, creator6, evidence)

	// fake votes
	v1id, _ := k.Create(ctx, storyID, amount, true, argument, creator7, evidence)
	v2id, _ := k.Create(ctx, storyID, amount, true, argument, creator8, evidence)
	v3id, _ := k.Create(ctx, storyID, amount, false, argument, creator9, evidence)

	b1, _ := backingKeeper.Backing(ctx, b1id)
	b2, _ := backingKeeper.Backing(ctx, b2id)
	b3, _ := backingKeeper.Backing(ctx, b3id)
	b4, _ := backingKeeper.Backing(ctx, b4id)

	c1, _ := challengeKeeper.Challenge(ctx, c1id)
	c2, _ := challengeKeeper.Challenge(ctx, c2id)

	v1, _ := k.Get(ctx, v1id)
	v2, _ := k.Get(ctx, v2id)
	v3, _ := k.Get(ctx, v3id)

	trueVotes = append(trueVotes, b1, b2, b3, v1, v2)
	falseVotes = append(falseVotes, b4, c1, c2, v3)

	return
}

func TestConfirmedStoryRewardPool(t *testing.T) {
	_, _, _, _, _, bankKeeper, _ := mockDB()

	ctx, _, falseVotes := createFakeConfirmedStory()

	pool, err := confirmedStoryRewardPool(ctx, bankKeeper, falseVotes)
	spew.Dump(pool)
	spew.Dump(err)
	assert.NotNil(t, pool)
}
