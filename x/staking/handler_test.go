package staking

import (
	"encoding/json"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"

	app "github.com/TruStory/truchain/types"
)

func TestHandle_SubmitArgument(t *testing.T) {
	ctx, k, mdb := mockDB()
	handler := NewHandler(k)
	addr1 := createFakeFundedAccount(ctx, mdb.authAccKeeper, sdk.Coins{sdk.NewInt64Coin(app.StakeDenom, app.Shanev*300)})

	msg1 := NewMsgSubmitArgument(addr1, 1, "summary 1", "body 1", StakeBacking)

	assert.Equal(t, msg1.Route(), RouterKey)
	assert.Equal(t, msg1.Type(), TypeMsgSubmitArgument)
	res := handler(ctx, msg1)

	assert.True(t, res.IsOK())

}

func TestHandle_SubmitUpvote(t *testing.T) {
	ctx, k, mdb := mockDB()
	handler := NewHandler(k)
	addr1 := createFakeFundedAccount(ctx, mdb.authAccKeeper, sdk.Coins{sdk.NewInt64Coin(app.StakeDenom, app.Shanev*300)})
	addr2 := createFakeFundedAccount(ctx, mdb.authAccKeeper, sdk.Coins{sdk.NewInt64Coin(app.StakeDenom, app.Shanev*300)})

	msg1 := NewMsgSubmitArgument(addr1, 1, "summary 1", "body 1", StakeBacking)

	assert.Equal(t, msg1.Route(), RouterKey)
	assert.Equal(t, msg1.Type(), TypeMsgSubmitArgument)
	res := handler(ctx, msg1)
	assert.True(t, res.IsOK())

	msg2 := NewMsgSubmitUpvote(addr2, 1)
	assert.Equal(t, msg2.Route(), RouterKey)
	assert.Equal(t, msg2.Type(), TypeMsgSubmitUpvote)
	res2 := handler(ctx, msg2)
	assert.True(t, res2.IsOK())

	assert.Equal(t, sdk.NewInt(app.Shanev*250), k.bankKeeper.GetCoins(ctx, addr1).AmountOf(app.StakeDenom))
	assert.Equal(t, sdk.NewInt(app.Shanev*290), k.bankKeeper.GetCoins(ctx, addr2).AmountOf(app.StakeDenom))

}

func TestHandleMsgAddAdmin(t *testing.T) {
	ctx, keeper, _ := mockDB()
	handler := NewHandler(keeper)
	assert.NotNil(t, handler) // assert handler is present

	_, _, admin := keyPubAddr()
	creator := keeper.GetParams(ctx).StakingAdmins[0]
	msg := NewMsgAddAdmin(admin, creator)
	assert.NotNil(t, msg) // assert msgs can be created

	result := handler(ctx, msg)
	var success bool
	err := json.Unmarshal(result.Data, &success)
	assert.NoError(t, err)
	assert.Equal(t, success, true)
}

func TestHandleMsgRemoveAdmin(t *testing.T) {
	ctx, keeper, _ := mockDB()
	handler := NewHandler(keeper)
	assert.NotNil(t, handler) // assert handler is present

	_, _, admin := keyPubAddr()
	creator := keeper.GetParams(ctx).StakingAdmins[0]
	msg := NewMsgRemoveAdmin(admin, creator)
	assert.NotNil(t, msg) // assert msgs can be created

	result := handler(ctx, msg)
	var success bool
	err := json.Unmarshal(result.Data, &success)
	assert.NoError(t, err)
	assert.Equal(t, success, true)
}

func TestByzantineMsg(t *testing.T) {
	ctx, k, _ := mockDB()

	handler := NewHandler(k)
	assert.NotNil(t, handler)

	res := handler(ctx, nil)
	assert.Equal(t, sdk.CodeUnknownRequest, res.Code)
	assert.Equal(t, sdk.CodespaceRoot, res.Codespace)
}
