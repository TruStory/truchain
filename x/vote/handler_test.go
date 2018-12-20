package vote

import (
	"encoding/json"
	"net/url"
	"testing"

	"github.com/TruStory/truchain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

func TestCreateVoteMsg(t *testing.T) {
	ctx, k, ck := mockDB()

	h := NewHandler(k)
	assert.NotNil(t, h)

	storyID := createFakeStory(ctx, k.storyKeeper, ck)
	amount := sdk.NewCoin("trudex", sdk.NewInt(15))
	creator := sdk.AccAddress([]byte{1, 2})
	testURL, _ := url.Parse("http://www.trustory.io")
	evidence := []url.URL{*testURL}

	// give user some funds
	k.bankKeeper.AddCoins(ctx, creator, sdk.Coins{amount.Plus(amount)})

	argument := "test argument"
	_, err := k.challengeKeeper.Create(ctx, storyID, amount, argument, creator, evidence)
	assert.Nil(t, err)

	evidence1 := []string{"http://www.trustory.io"}
	msg := NewCreateVoteMsg(storyID, amount, "valid comment", creator, evidence1, true)
	assert.NotNil(t, msg)

	res := h(ctx, msg)
	idres := new(types.IDResult)
	_ = json.Unmarshal(res.Data, &idres)

	assert.Equal(t, int64(1), idres.ID)
}
