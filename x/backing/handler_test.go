package backing

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	params "github.com/TruStory/truchain/parameters"
	"github.com/TruStory/truchain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

func TestBackStoryMsg_FailBasicValidation(t *testing.T) {
	ctx, bk, _, _, _, _ := mockDB()

	h := NewHandler(bk)
	assert.NotNil(t, h)

	storyID := int64(1)
	amount, _ := sdk.ParseCoin("5trushane")
	argument := "cool story brew"
	creator := sdk.AccAddress([]byte{1, 2})
	// duration := 5 * time.Hour
	duration := 15 * time.Minute
	msg := NewBackStoryMsg(storyID, amount, argument, creator, duration)
	assert.NotNil(t, msg)

	res := h(ctx, msg)
	hasInvalidBackingPeriod := strings.Contains(res.Log, "901")
	assert.True(t, hasInvalidBackingPeriod, "should return err code")
}

func TestBackStoryMsg_FailInsufficientFunds(t *testing.T) {
	ctx, bk, _, _, _, _ := mockDB()

	h := NewHandler(bk)
	assert.NotNil(t, h)

	storyID := int64(1)
	amount := sdk.NewCoin(params.StakeDenom, sdk.NewInt(5000000))
	argument := "cool story brew"
	creator := sdk.AccAddress([]byte{1, 2})
	// duration := 99 * time.Hour
	duration := 24 * time.Hour
	msg := NewBackStoryMsg(storyID, amount, argument, creator, duration)
	assert.NotNil(t, msg)

	res := h(ctx, msg)
	hasInsufficientFunds := strings.Contains(res.Log, "65541")
	assert.True(t, hasInsufficientFunds, "should return err code")
}

func TestBackStoryMsg(t *testing.T) {
	ctx, bk, sk, ck, _, am := mockDB()

	h := NewHandler(bk)
	assert.NotNil(t, h)

	storyID := createFakeStory(ctx, sk, ck)
	amount := sdk.NewCoin(params.StakeDenom, sdk.NewInt(5000000))
	argument := "cool story brew"
	creator := createFakeFundedAccount(ctx, am, sdk.Coins{amount})
	// duration := 99 * time.Hour
	duration := 24 * time.Hour
	msg := NewBackStoryMsg(storyID, amount, argument, creator, duration)
	assert.NotNil(t, msg)

	res := h(ctx, msg)
	idres := new(types.IDResult)
	_ = json.Unmarshal(res.Data, &idres)

	assert.Equal(t, int64(1), idres.ID, "incorrect result backing id")
}
