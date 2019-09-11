package bank

import (
	"testing"

	app "github.com/TruStory/truchain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

func TestHandle_MsgPayReward(t *testing.T) {
	ctx, keeper, ak := mockDB()
	handler := NewHandler(keeper)

	amount := sdk.NewCoin(app.StakeDenom, sdk.NewInt(app.Shanev*100))

	brokerAddress := createFakeFundedAccount(ctx, ak, sdk.Coins{amount})
	assert.NotNil(t, handler)
	recipientAddr := createFakeFundedAccount(ctx, ak, sdk.Coins{})
	reward := sdk.NewCoin(app.StakeDenom, sdk.NewInt(app.Shanev*15))
	msg := NewMsgPayReward(brokerAddress, recipientAddr, reward, 1)

	assert.Equal(t, msg.Route(), RouterKey)
	assert.Equal(t, msg.Type(), TypeMsgPayReward)
	res := handler(ctx, msg)

	assert.Equal(t, ErrorCodeInvalidRewardBrokerAddress, res.Code)
	assert.Equal(t, DefaultCodespace, res.Codespace)
	p := Params{RewardBrokerAddress: brokerAddress}
	keeper.SetParams(ctx, p)
	res = handler(ctx, msg)
	assert.True(t, res.IsOK())

	recipientCoins := keeper.bankKeeper.GetCoins(ctx, recipientAddr)
	assert.True(t, recipientCoins.AmountOf(app.StakeDenom).Equal(sdk.NewInt(app.Shanev*15)))
}

func TestMsgPayReward_Invalid(t *testing.T) {
	ctx, keeper, _ := mockDB()
	_, _, validAddr := keyPubAddr()

	msg := NewMsgPayReward(sdk.AccAddress{}, sdk.AccAddress{}, sdk.Coin{}, 1)
	handler := NewHandler(keeper)

	res := handler(ctx, msg)
	assert.False(t, res.IsOK())
	assert.Equal(t, sdk.CodeInvalidAddress, res.Code)
	assert.Equal(t, sdk.CodespaceRoot, res.Codespace)

	msg = NewMsgPayReward(validAddr, sdk.AccAddress{}, sdk.Coin{}, 1)
	handler = NewHandler(keeper)

	res = handler(ctx, msg)
	assert.False(t, res.IsOK())
	assert.Equal(t, sdk.CodeInvalidAddress, res.Code)
	assert.Equal(t, sdk.CodespaceRoot, res.Codespace)

}

func TestByzantineMsg(t *testing.T) {
	ctx, keeper, _ := mockDB()

	handler := NewHandler(keeper)
	assert.NotNil(t, handler)

	res := handler(ctx, nil)
	assert.Equal(t, sdk.CodeUnknownRequest, res.Code)
	assert.Equal(t, sdk.CodespaceRoot, res.Codespace)
}
