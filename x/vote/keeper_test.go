package vote

import (
	"testing"

	"github.com/TruStory/truchain/x/story"

	app "github.com/TruStory/truchain/types"
	"github.com/TruStory/truchain/x/stake"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

func TestCreateGetVote(t *testing.T) {
	ctx, k, ck := mockDB()

	storyID := createFakeStory(ctx, k.storyKeeper, ck, story.Pending)
	amount := sdk.NewCoin(app.StakeDenom, sdk.NewInt(15000000000))
	comment := "test comment is long enough"
	creator := sdk.AccAddress([]byte{1, 2})

	// give user some funds
	k.bankKeeper.AddCoins(ctx, creator, sdk.Coins{amount.Plus(amount)})

	argument := "test argument"
	_, err := k.challengeKeeper.Create(ctx, storyID, amount, argument, creator, false)
	assert.Nil(t, err)

	voteID, err := k.Create(ctx, storyID, amount, true, comment, creator)
	assert.Nil(t, err)

	vote, _ := k.TokenVote(ctx, voteID)
	assert.Equal(t, voteID, vote.ID())
}

func TestInValidCreateVoteMsgArgumentTooShort(t *testing.T) {
	ctx, k, ck := mockDB()

	storyID := createFakeStory(ctx, k.storyKeeper, ck, story.Challenged)
	amount := sdk.NewCoin(app.StakeDenom, sdk.NewInt(15000000000))
	creator := sdk.AccAddress([]byte{1, 2})

	// give user some funds
	k.bankKeeper.AddCoins(ctx, creator, sdk.Coins{amount.Plus(amount)})

	argument := "too short"
	_, err := k.Create(ctx, storyID, amount, true, argument, creator)
	assert.NotNil(t, err)
	assert.Equal(t, stake.ErrArgumentTooShortMsg(argument, len(argument)).Code(), err.Code())
}

func TestInValidCreateVoteMsgArgumentTooLong(t *testing.T) {
	ctx, k, ck := mockDB()

	storyID := createFakeStory(ctx, k.storyKeeper, ck, story.Challenged)
	amount := sdk.NewCoin(app.StakeDenom, sdk.NewInt(15000000000))
	creator := sdk.AccAddress([]byte{1, 2})

	// give user some funds
	k.bankKeeper.AddCoins(ctx, creator, sdk.Coins{amount.Plus(amount)})

	argument := "Escaping predators, digestion and other animal activities—including those of humans—require oxygen. But that essential ingredient is no longer so easy for marine life to obtain, several new studies reveal. In the past decade ocean oxygen levels have taken a dive—an alarming trend that is linked to climate change, says Andreas Oschlies, an oceanographer at the Helmholtz Center for Ocean Research Kiel in Germany, whose team tracks ocean oxygen levels worldwide. “We were surprised by the intensity of the changes we saw, how rapidly oxygen is going down in the ocean and how large the effects on marine ecosystems are,” he says. It is no surprise to scientists that warming oceans are losing oxygen, but the scale of the dip calls for urgent attention, Oschlies says. Oxygen levels in some tropical regions have dropped by a startling 40 percent in the last 50 years, some recent studies reveal. Levels have dropped more subtly elsewhere, with an average loss of 2 percent globally. Ocean animals large and small, however, respond to even slight changes in oxygen by seeking refuge in higher oxygen zones or by adjusting behavior, Oschlies and others in his field have found. These adjustments can expose animals to new predators or force them into food-scarce regions. Climate change already poses serious problems for marine life, such as ocean acidification, but deoxygenation is the most pressing issue facing sea animals today, Oschlies says. After all, he says, “they all have to breathe.”"
	_, err := k.Create(ctx, storyID, amount, true, argument, creator)
	assert.NotNil(t, err)
	assert.Equal(t, stake.ErrArgumentTooLongMsg(len(argument)).Code(), err.Code())
}

func TestGetVotesByGameID(t *testing.T) {
	ctx, k, ck := mockDB()

	storyID := createFakeStory(ctx, k.storyKeeper, ck, story.Pending)
	amount := sdk.NewCoin(app.StakeDenom, sdk.NewInt(15000000000))
	comment := "test comment is long enough"
	creator := sdk.AccAddress([]byte{1, 2})
	creator2 := sdk.AccAddress([]byte{3, 4})

	// give user some funds
	k.bankKeeper.AddCoins(ctx, creator, sdk.Coins{amount.Plus(amount)})
	k.bankKeeper.AddCoins(ctx, creator2, sdk.Coins{amount})

	argument := "test argument"
	_, err := k.challengeKeeper.Create(ctx, storyID, amount, argument, creator, false)
	assert.Nil(t, err)

	_, err = k.Create(ctx, storyID, amount, true, comment, creator)
	assert.Nil(t, err)

	_, err = k.Create(ctx, storyID, amount, true, comment, creator2)
	assert.Nil(t, err)

	story, _ := k.storyKeeper.Story(ctx, storyID)

	votes, _ := k.TokenVotesByStoryID(ctx, story.ID)
	assert.Equal(t, 2, len(votes))
}

func TestGetVotesByStoryIDAndCreator(t *testing.T) {
	ctx, k, ck := mockDB()

	storyID := createFakeStory(ctx, k.storyKeeper, ck, story.Pending)
	amount := sdk.NewCoin(app.StakeDenom, sdk.NewInt(15000000000))
	comment := "test comment is long enough"
	creator := sdk.AccAddress([]byte{1, 2})

	// give user some funds
	k.bankKeeper.AddCoins(ctx, creator, sdk.Coins{amount.Plus(amount)})

	argument := "test argument"
	_, err := k.challengeKeeper.Create(ctx, storyID, amount, argument, creator, false)
	assert.Nil(t, err)

	_, err = k.Create(ctx, storyID, amount, true, comment, creator)
	assert.Nil(t, err)

	vote, _ := k.TokenVotesByStoryIDAndCreator(ctx, storyID, creator)
	assert.Equal(t, int64(1), vote.ID())
}

func TestTotalVoteAmountByGameID(t *testing.T) {
	ctx, k, ck := mockDB()

	storyID := createFakeStory(ctx, k.storyKeeper, ck, story.Pending)
	amount := sdk.NewCoin(app.StakeDenom, sdk.NewInt(15000000000))
	comment := "test comment is long enough"
	creator := sdk.AccAddress([]byte{1, 2})
	creator1 := sdk.AccAddress([]byte{2, 3})

	// give user some funds
	k.bankKeeper.AddCoins(ctx, creator, sdk.Coins{amount.Plus(amount)})
	k.bankKeeper.AddCoins(ctx, creator1, sdk.Coins{amount.Plus(amount)})

	argument := "test argument"
	_, err := k.challengeKeeper.Create(ctx, storyID, amount, argument, creator, false)
	assert.Nil(t, err)

	// create votes
	_, err = k.Create(ctx, storyID, amount, true, comment, creator)
	assert.Nil(t, err)
	_, err = k.Create(ctx, storyID, amount, true, comment, creator1)
	assert.Nil(t, err)

	story, _ := k.storyKeeper.Story(ctx, storyID)

	totalAmount, _ := k.TotalVoteAmountByStoryID(ctx, story.ID)
	assert.Equal(t, "30000000000trusteak", totalAmount.String())
}

func TestCreateVote_ErrGameNotStarted(t *testing.T) {
	ctx, k, ck := mockDB()

	storyID := createFakeStory(ctx, k.storyKeeper, ck, story.Pending)
	amount := sdk.NewCoin(app.StakeDenom, sdk.NewInt(50000000000000))
	comment := "test comment is long enough"
	creator := sdk.AccAddress([]byte{1, 2})

	vote := true

	_, err := k.Create(ctx, storyID, amount, vote, comment, creator)
	assert.NotNil(t, err)
	assert.Equal(t, ErrVotingNotStarted(storyID).Code(), err.Code())
}

func TestCreateVote_ErrBelowMinStake(t *testing.T) {
	ctx, k, ck := mockDB()

	storyID := createFakeStory(ctx, k.storyKeeper, ck, story.Challenged)
	amount := sdk.NewCoin(app.StakeDenom, sdk.NewInt(5000000))
	comment := "test comment is long enough"
	creator := sdk.AccAddress([]byte{1, 2})

	vote := true

	_, err := k.Create(ctx, storyID, amount, vote, comment, creator)
	assert.NotNil(t, err)
	assert.Equal(t, sdk.ErrInsufficientFunds("Below minimum stake.").Code(), err.Code())
}

func TestUpdateVote_AddWeightOnTally(t *testing.T) {
	ctx, k, ck := mockDB()

	storyID := createFakeStory(ctx, k.storyKeeper, ck, 0)
	amount := sdk.NewCoin(app.StakeDenom, sdk.NewInt(15000000000))
	comment := "test comment is long enough"
	creator := sdk.AccAddress([]byte{1, 2})

	// give user some funds
	k.bankKeeper.AddCoins(ctx, creator, sdk.Coins{amount.Plus(amount)})

	argument := "test argument"
	_, err := k.challengeKeeper.Create(ctx, storyID, amount, argument, creator, false)
	assert.Nil(t, err)

	_, err1 := k.Create(ctx, storyID, amount, true, comment, creator)
	assert.Nil(t, err1)

	vote, _ := k.TokenVotesByStoryIDAndCreator(ctx, storyID, creator)
	assert.Equal(t, int64(1), vote.ID())

	vote.UpdateWeight(sdk.NewInt(1000000000))

	k.Update(ctx, vote)
	updatedVote, _ := k.TokenVotesByStoryIDAndCreator(ctx, storyID, creator)
	assert.Equal(t, int64(1), updatedVote.ID())

	assert.Equal(t, updatedVote.Weight().String(), "1000000000")
}

func Test_ToggleVote_ChallengeToBacking(t *testing.T) {
	ctx, k, ck := mockDB()

	storyID := createFakeStory(ctx, k.storyKeeper, ck, story.Pending)
	amount := sdk.NewCoin(app.StakeDenom, sdk.NewInt(15000000000))

	creator := sdk.AccAddress([]byte{1, 2})

	// give user some funds
	k.bankKeeper.AddCoins(ctx, creator, sdk.Coins{amount.Plus(amount)})

	argument := "test argument"
	_, err := k.challengeKeeper.Create(ctx, storyID, amount, argument, creator, false)
	assert.Nil(t, err)

	backingID, err := k.ToggleVote(ctx, storyID, amount, argument, creator)
	assert.NoError(t, err)

	b, err := k.backingKeeper.Backing(ctx, backingID)
	assert.NoError(t, err)
	assert.NotNil(t, b.Vote)
}

func Test_ToggleVote_BackingToChallenge(t *testing.T) {
	ctx, k, ck := mockDB()

	storyID := createFakeStory(ctx, k.storyKeeper, ck, story.Pending)
	amount := sdk.NewCoin(app.StakeDenom, sdk.NewInt(15000000000))

	creator := sdk.AccAddress([]byte{1, 2})
	creator1 := sdk.AccAddress([]byte{2, 3})

	// give user some funds
	k.bankKeeper.AddCoins(ctx, creator, sdk.Coins{amount.Plus(amount)})
	k.bankKeeper.AddCoins(ctx, creator1, sdk.Coins{amount.Plus(amount)})

	argument := "test backing argument"
	_, err := k.backingKeeper.Create(ctx, storyID, amount, argument, creator, false)
	assert.Nil(t, err)

	argument = "test  challenge argument"
	_, err = k.challengeKeeper.Create(ctx, storyID, amount, argument, creator1, false)
	assert.NoError(t, err)

	s, err := k.storyKeeper.Story(ctx, storyID)
	assert.NoError(t, err)

	assert.Equal(t, story.Challenged, s.Status)

	challengeID, err := k.ToggleVote(ctx, storyID, amount, argument, creator)
	assert.NoError(t, err)

	b, err := k.challengeKeeper.Challenge(ctx, challengeID)
	assert.NoError(t, err)
	assert.NotNil(t, b.Vote)
}

func Test_ToggleVote_TokenVote_FalseTrue(t *testing.T) {
	ctx, k, ck := mockDB()

	storyID := createFakeStory(ctx, k.storyKeeper, ck, story.Pending)
	amount := sdk.NewCoin(app.StakeDenom, sdk.NewInt(15000000000))

	creator := sdk.AccAddress([]byte{1, 2})
	creator1 := sdk.AccAddress([]byte{2, 3})

	// give user some funds
	k.bankKeeper.AddCoins(ctx, creator, sdk.Coins{amount.Plus(amount)})
	k.bankKeeper.AddCoins(ctx, creator1, sdk.Coins{amount.Plus(amount)})

	argument := "test  challenge argument"
	_, err := k.challengeKeeper.Create(ctx, storyID, amount, argument, creator, false)
	assert.NoError(t, err)

	// Test story switched to challenged state
	s, err := k.storyKeeper.Story(ctx, storyID)
	assert.NoError(t, err)
	assert.Equal(t, story.Challenged, s.Status)

	// Create Token vote
	_, err = k.Create(ctx, storyID, amount, false, argument, creator1)
	assert.NoError(t, err)

	// Should return the ID of the token vote
	tokenVoteID, err := k.ToggleVote(ctx, storyID, amount, argument, creator1)
	assert.NoError(t, err)

	b, err := k.TokenVote(ctx, tokenVoteID)
	assert.NoError(t, err)
	assert.NotNil(t, b.Vote)
	assert.True(t, b.VoteChoice())
}

func Test_ToggleVote_TokenVote_TrueFalse(t *testing.T) {
	ctx, k, ck := mockDB()

	storyID := createFakeStory(ctx, k.storyKeeper, ck, story.Pending)
	amount := sdk.NewCoin(app.StakeDenom, sdk.NewInt(15000000000))

	creator := sdk.AccAddress([]byte{1, 2})
	creator1 := sdk.AccAddress([]byte{2, 3})

	// give user some funds
	k.bankKeeper.AddCoins(ctx, creator, sdk.Coins{amount.Plus(amount)})
	k.bankKeeper.AddCoins(ctx, creator1, sdk.Coins{amount.Plus(amount)})

	argument := "test  challenge argument"
	_, err := k.challengeKeeper.Create(ctx, storyID, amount, argument, creator, false)
	assert.NoError(t, err)

	// Test story switched to challenged state
	s, err := k.storyKeeper.Story(ctx, storyID)
	assert.NoError(t, err)
	assert.Equal(t, story.Challenged, s.Status)

	// Create Token vote
	_, err = k.Create(ctx, storyID, amount, true, argument, creator1)
	assert.NoError(t, err)

	// Should return the ID of the token vote
	tokenVoteID, err := k.ToggleVote(ctx, storyID, amount, argument, creator1)
	assert.NoError(t, err)

	b, err := k.TokenVote(ctx, tokenVoteID)
	assert.NoError(t, err)
	assert.NotNil(t, b.Vote)
	assert.False(t, b.VoteChoice())
}

func Test_ToggleVote_ChallengeState(t *testing.T) {
	ctx, k, ck := mockDB()

	storyID := createFakeStory(ctx, k.storyKeeper, ck, story.Pending)
	amount := sdk.NewCoin(app.StakeDenom, sdk.NewInt(15000000000))

	creator := sdk.AccAddress([]byte{1, 2})

	// give user some funds
	k.bankKeeper.AddCoins(ctx, creator, sdk.Coins{amount.Plus(amount)})

	argument := "test backing argument"
	_, err := k.backingKeeper.Create(ctx, storyID, amount, argument, creator, false)
	assert.Nil(t, err)

	s, err := k.storyKeeper.Story(ctx, storyID)
	assert.NoError(t, err)

	assert.Equal(t, story.Pending, s.Status)

	_, err = k.ToggleVote(ctx, storyID, amount, argument, creator)
	assert.Error(t, err)
	assert.Equal(t, CodeInvalidStoryState, err.Code())
}

func Test_Toggle_InvalidStoryID(t *testing.T) {
	ctx, k, _ := mockDB()

	amount := sdk.NewCoin(app.StakeDenom, sdk.NewInt(15000000000))

	creator := sdk.AccAddress([]byte{1, 2})

	argument := "test backing argument"
	_, err := k.ToggleVote(ctx, 0, amount, argument, creator)
	assert.Error(t, err)
	assert.Equal(t, story.CodeStoryNotFound, err.Code())
}

func Test_Toggle_VoteNotFound(t *testing.T) {
	ctx, k, ck := mockDB()

	amount := sdk.NewCoin(app.StakeDenom, sdk.NewInt(15000000000))
	storyID := createFakeStory(ctx, k.storyKeeper, ck, story.Challenged)
	creator := sdk.AccAddress([]byte{1, 2})

	argument := "test backing argument"
	_, err := k.ToggleVote(ctx, storyID, amount, argument, creator)
	assert.Error(t, err)
	assert.Equal(t, CodeNotFound, err.Code())

}
