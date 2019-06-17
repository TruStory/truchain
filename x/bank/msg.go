package bank

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	TypeMsgPayReward = "pay_reward"
)

var (
	_ sdk.Msg = &MsgPayReward{}
)

type MsgPayReward struct {
	Creator   sdk.AccAddress
	Recipient sdk.AccAddress
	Reward    sdk.Coin
	InviteID  int64
}

func NewMsgPayReward(creator, recipient sdk.AccAddress, reward sdk.Coin, inviteID int64) MsgPayReward {
	return MsgPayReward{
		Creator:   creator,
		Recipient: recipient,
		Reward:    reward,
		InviteID:  inviteID,
	}
}
func (msg MsgPayReward) Route() string { return RouterKey }

func (msg MsgPayReward) Type() string {
	return TypeMsgPayReward
}

func (msg MsgPayReward) ValidateBasic() sdk.Error {
	if len(msg.Creator) == 0 {
		return sdk.ErrInvalidAddress("invalid creator address")
	}
	if len(msg.Recipient) == 0 {
		return sdk.ErrInvalidAddress("invalid recipient address")
	}

	if msg.Reward.IsNegative() || msg.Reward.IsZero() {
		return sdk.ErrInvalidCoins("invalid coins")
	}
	return nil
}

func (msg MsgPayReward) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Creator}
}

func (msg MsgPayReward) GetSignBytes() []byte {
	bz := moduleCodec.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}
