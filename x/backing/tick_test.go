package backing

import (
	"testing"
	"time"

	app "github.com/TruStory/truchain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

func TestNewResponseEndBlock(t *testing.T) {
	ctx, bk, _, _, _, _ := mockDB()

	tags := bk.NewResponseEndBlock(ctx)
	assert.Nil(t, tags)
}

func Test_processMaturedBackings_QueueEmpty(t *testing.T) {
	ctx, bk, _, _, _, _ := mockDB()

	err := bk.processMaturedBackings(ctx)
	assert.Nil(t, err)
}

func Test_processMaturedBackings_UnmaturedBackings(t *testing.T) {
	ctx, bk, sk, ck, bankKeeper, _ := mockDB()
	storyID := createFakeStory(ctx, sk, ck)
	amount, _ := sdk.ParseCoin("5trudex")
	argument := "cool story brew"
	creator := sdk.AccAddress([]byte{1, 2})
	duration := 99 * time.Hour
	bankKeeper.AddCoins(ctx, creator, sdk.Coins{amount})
	bk.Create(ctx, storyID, amount, argument, creator, duration)
	bk.Create(ctx, storyID, amount, argument, creator, duration)
	bk.Create(ctx, storyID, amount, argument, creator, duration)

	err := bk.processMaturedBackings(ctx)
	assert.Nil(t, err)
}

func Test_processMaturedBackings_MaturedBackings(t *testing.T) {
	ctx, bk, sk, ck, bankKeeper, _ := mockDB()
	storyID := createFakeStory(ctx, sk, ck)
	amount, _ := sdk.ParseCoin("5trudex")
	argument := "cool story brew"
	creator := sdk.AccAddress([]byte{1, 2})
	duration := 4 * time.Hour
	bankKeeper.AddCoins(ctx, creator, sdk.Coins{amount})
	bk.Create(ctx, storyID, amount, argument, creator, duration)
	bk.Create(ctx, storyID, amount, argument, creator, duration)
	bk.Create(ctx, storyID, amount, argument, creator, duration)
	bk.Create(ctx, storyID, amount, argument, creator, duration)
	bk.Create(ctx, storyID, amount, argument, creator, duration)

	err := bk.processMaturedBackings(ctx)
	assert.Nil(t, err)
}

func Test_processMaturedBackings_GameInSession(t *testing.T) {
	ctx, bk, sk, ck, bankKeeper, _ := mockDB()
	storyID := createFakeStory(ctx, sk, ck)
	amount, _ := sdk.ParseCoin("5trudex")
	argument := "cool story brew"
	creator := sdk.AccAddress([]byte{1, 2})
	duration := 99 * time.Hour
	bankKeeper.AddCoins(ctx, creator, sdk.Coins{amount})
	bk.Create(ctx, storyID, amount, argument, creator, duration)
	bk.Create(ctx, storyID, amount, argument, creator, duration)
	bk.Create(ctx, storyID, amount, argument, creator, duration)
	bk.Create(ctx, storyID, amount, argument, creator, duration)
	bk.Create(ctx, storyID, amount, argument, creator, duration)

	story, _ := sk.Story(ctx, storyID)
	story.GameID = 5
	sk.UpdateStory(ctx, story)

	gameList := bk.pendingGameList(ctx)
	var testID int64 = 5
	gameList.Push(testID)

	err := bk.processMaturedBackings(ctx)
	assert.Nil(t, err)
}

func Test_distributeEarnings(t *testing.T) {
	ctx, bk, _, _, _, _ := mockDB()

	principal, _ := sdk.ParseCoin("5trudex")
	interest, _ := sdk.ParseCoin("2trudex")
	matures := time.Now().Add(24 * time.Hour)
	params := DefaultParams()
	duration := 24 * time.Hour
	creator := sdk.AccAddress([]byte{1, 2})

	// create new vote type
	vote := app.Vote{
		ID:        int64(1),
		Amount:    principal,
		Argument:  "",
		Creator:   creator,
		Vote:      true,
		Timestamp: app.NewTimestamp(ctx.BlockHeader()),
	}

	// create new backing type with embedded vote
	backing := Backing{
		vote,
		int64(5),
		interest,
		matures,
		params,
		duration,
	}

	err := bk.distributeEarnings(ctx, backing)
	assert.Nil(t, err)
}

func Test_isGameInList(t *testing.T) {
	ctx, k, _, _, _, _ := mockDB()
	gameList := k.pendingGameList(ctx)
	var testID int64 = 5
	gameList.Push(testID)

	found := k.isGameInList(gameList, 5)
	assert.True(t, found)

	found = k.isGameInList(gameList, 4)
	assert.False(t, found)
}
