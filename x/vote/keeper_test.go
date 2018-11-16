package vote

import (
	"net/url"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

func TestCreateGetVote(t *testing.T) {
	ctx, k, sk, ck, challengeKeeper, bankKeeper, _, _ := mockDB()

	storyID := createFakeStory(ctx, sk, ck)
	amount := sdk.NewCoin("trudex", sdk.NewInt(15))
	comment := "test comment is long enough"
	creator := sdk.AccAddress([]byte{1, 2})
	cnn, _ := url.Parse("http://www.cnn.com")
	evidence := []url.URL{*cnn}

	// give user some funds
	bankKeeper.AddCoins(ctx, creator, sdk.Coins{amount.Plus(amount)})

	argument := "test argument"
	_, err := challengeKeeper.Create(ctx, storyID, amount, argument, creator, evidence)
	assert.Nil(t, err)

	voteID, err := k.Create(ctx, storyID, amount, true, comment, creator, evidence)
	assert.Nil(t, err)

	vote, _ := k.Vote(ctx, voteID)
	assert.Equal(t, voteID, vote.ID)
}

func TestGetVotesByGame(t *testing.T) {
	ctx, k, sk, ck, challengeKeeper, bankKeeper, _, _ := mockDB()

	storyID := createFakeStory(ctx, sk, ck)
	amount := sdk.NewCoin("trudex", sdk.NewInt(15))
	comment := "test comment is long enough"
	creator := sdk.AccAddress([]byte{1, 2})
	creator2 := sdk.AccAddress([]byte{3, 4})
	cnn, _ := url.Parse("http://www.cnn.com")
	evidence := []url.URL{*cnn}

	// give user some funds
	bankKeeper.AddCoins(ctx, creator, sdk.Coins{amount.Plus(amount)})
	bankKeeper.AddCoins(ctx, creator2, sdk.Coins{amount})

	argument := "test argument"
	_, err := challengeKeeper.Create(ctx, storyID, amount, argument, creator, evidence)
	assert.Nil(t, err)

	_, err = k.Create(ctx, storyID, amount, true, comment, creator, evidence)
	assert.Nil(t, err)

	_, err = k.Create(ctx, storyID, amount, true, comment, creator2, evidence)
	assert.Nil(t, err)

	story, _ := sk.GetStory(ctx, storyID)

	votes, _ := k.VotesByGame(ctx, story.GameID)
	assert.Equal(t, 2, len(votes))
}

func TestCreateVote_ErrGameNotStarted(t *testing.T) {
	ctx, k, sk, ck, _, _, _, _ := mockDB()

	storyID := createFakeStory(ctx, sk, ck)
	amount := sdk.NewCoin("trudex", sdk.NewInt(15))
	comment := "test comment is long enough"
	creator := sdk.AccAddress([]byte{1, 2})
	cnn, _ := url.Parse("http://www.cnn.com")
	evidence := []url.URL{*cnn}

	vote := true

	_, err := k.Create(ctx, storyID, amount, vote, comment, creator, evidence)
	assert.NotNil(t, err)
	assert.Equal(t, ErrGameNotStarted(storyID).Code(), err.Code())
}
