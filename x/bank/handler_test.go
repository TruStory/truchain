package bank

import (
	"fmt"
	"testing"

	app "github.com/TruStory/truchain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

func TestHandle_MsgSendGift(t *testing.T) {
	ctx, keeper, ak := mockDB()
	handler := NewHandler(keeper)

	amount := sdk.NewCoin(app.StakeDenom, sdk.NewInt(app.Shanev*100))

	brokerAddress := createFakeFundedAccount(ctx, ak, sdk.Coins{amount})
	assert.NotNil(t, handler)
	recipientAddr := createFakeFundedAccount(ctx, ak, sdk.Coins{})
	reward := sdk.NewCoin(app.StakeDenom, sdk.NewInt(app.Shanev*15))
	msg := NewMsgSendGift(brokerAddress, recipientAddr, reward)

	assert.Equal(t, msg.Route(), RouterKey)
	assert.Equal(t, msg.Type(), TypeMsgSendGift)
	res := handler(ctx, msg)

	assert.Equal(t, ErrorCodeInvalidRewardBrokerAddress, res.Code)
	assert.Equal(t, DefaultCodespace, res.Codespace)
	p := Params{RewardBrokerAddress: brokerAddress}
	keeper.SetParams(ctx, p)
	res = handler(ctx, msg)
	fmt.Println(res)
	assert.True(t, res.IsOK())

	recipientCoins := keeper.bankKeeper.GetCoins(ctx, recipientAddr)
	assert.True(t, recipientCoins.AmountOf(app.StakeDenom).Equal(sdk.NewInt(app.Shanev*15)))
}

func TestMsgSendGift_Invalid(t *testing.T) {
	ctx, keeper, _ := mockDB()
	_, _, validAddr := keyPubAddr()

	msg := NewMsgSendGift(sdk.AccAddress{}, sdk.AccAddress{}, sdk.Coin{})
	handler := NewHandler(keeper)

	res := handler(ctx, msg)
	assert.False(t, res.IsOK())
	assert.Equal(t, sdk.CodeInvalidAddress, res.Code)
	assert.Equal(t, sdk.CodespaceRoot, res.Codespace)

	msg = NewMsgSendGift(validAddr, sdk.AccAddress{}, sdk.Coin{})
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
