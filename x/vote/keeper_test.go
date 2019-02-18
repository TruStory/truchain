package vote

import (
	"testing"

	params "github.com/TruStory/truchain/parameters"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

func TestCreateGetVote(t *testing.T) {
	ctx, k, ck := mockDB()

	storyID := createFakeStory(ctx, k.storyKeeper, ck)
	amount := sdk.NewCoin(params.StakeDenom, sdk.NewInt(15000000000))
	comment := "test comment is long enough"
	creator := sdk.AccAddress([]byte{1, 2})

	// give user some funds
	k.bankKeeper.AddCoins(ctx, creator, sdk.Coins{amount.Plus(amount)})

	argument := "test argument"
	_, err := k.challengeKeeper.Create(ctx, storyID, amount, argument, creator)
	assert.Nil(t, err)

	voteID, err := k.Create(ctx, storyID, amount, true, comment, creator)
	assert.Nil(t, err)

	vote, _ := k.TokenVote(ctx, voteID)
	assert.Equal(t, voteID, vote.ID())
}

func TestGetVotesByGameID(t *testing.T) {
	ctx, k, ck := mockDB()

	storyID := createFakeStory(ctx, k.storyKeeper, ck)
	amount := sdk.NewCoin(params.StakeDenom, sdk.NewInt(15000000000))
	comment := "test comment is long enough"
	creator := sdk.AccAddress([]byte{1, 2})
	creator2 := sdk.AccAddress([]byte{3, 4})

	// give user some funds
	k.bankKeeper.AddCoins(ctx, creator, sdk.Coins{amount.Plus(amount)})
	k.bankKeeper.AddCoins(ctx, creator2, sdk.Coins{amount})

	argument := "test argument"
	_, err := k.challengeKeeper.Create(ctx, storyID, amount, argument, creator)
	assert.Nil(t, err)

	_, err = k.Create(ctx, storyID, amount, true, comment, creator)
	assert.Nil(t, err)

	_, err = k.Create(ctx, storyID, amount, true, comment, creator2)
	assert.Nil(t, err)

	story, _ := k.storyKeeper.Story(ctx, storyID)

	votes, _ := k.TokenVotesByGameID(ctx, story.ID)
	assert.Equal(t, 2, len(votes))
}

func TestGetVotesByStoryIDAndCreator(t *testing.T) {
	ctx, k, ck := mockDB()

	storyID := createFakeStory(ctx, k.storyKeeper, ck)
	amount := sdk.NewCoin(params.StakeDenom, sdk.NewInt(15000000000))
	comment := "test comment is long enough"
	creator := sdk.AccAddress([]byte{1, 2})

	// give user some funds
	k.bankKeeper.AddCoins(ctx, creator, sdk.Coins{amount.Plus(amount)})

	argument := "test argument"
	_, err := k.challengeKeeper.Create(ctx, storyID, amount, argument, creator)
	assert.Nil(t, err)

	_, err = k.Create(ctx, storyID, amount, true, comment, creator)
	assert.Nil(t, err)

	vote, _ := k.TokenVotesByStoryIDAndCreator(ctx, storyID, creator)
	assert.Equal(t, int64(1), vote.ID())
}

func TestTotalVoteAmountByGameID(t *testing.T) {
	ctx, k, ck := mockDB()

	storyID := createFakeStory(ctx, k.storyKeeper, ck)
	amount := sdk.NewCoin(params.StakeDenom, sdk.NewInt(15000000000))
	comment := "test comment is long enough"
	creator := sdk.AccAddress([]byte{1, 2})
	creator1 := sdk.AccAddress([]byte{2, 3})

	// give user some funds
	k.bankKeeper.AddCoins(ctx, creator, sdk.Coins{amount.Plus(amount)})
	k.bankKeeper.AddCoins(ctx, creator1, sdk.Coins{amount.Plus(amount)})

	argument := "test argument"
	_, err := k.challengeKeeper.Create(ctx, storyID, amount, argument, creator)
	assert.Nil(t, err)

	// create votes
	_, err = k.Create(ctx, storyID, amount, true, comment, creator)
	assert.Nil(t, err)
	_, err = k.Create(ctx, storyID, amount, true, comment, creator1)
	assert.Nil(t, err)

	story, _ := k.storyKeeper.Story(ctx, storyID)

	totalAmount, _ := k.TotalVoteAmountByGameID(ctx, story.ID)
	assert.Equal(t, "30000000000trusteak", totalAmount.String())
}

// func TestCreateVote_ErrGameNotStarted(t *testing.T) {
// 	ctx, k, ck := mockDB()

// 	storyID := createFakeStory(ctx, k.storyKeeper, ck)
// 	amount := sdk.NewCoin(params.StakeDenom, sdk.NewInt(5000000))
// 	comment := "test comment is long enough"
// 	creator := sdk.AccAddress([]byte{1, 2})

// 	vote := true

// 	_, err := k.Create(ctx, storyID, amount, vote, comment, creator)
// 	assert.NotNil(t, err)
// 	assert.Equal(t, ErrGameNotStarted(storyID).Code(), err.Code())
// }
