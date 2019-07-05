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
	Sender    sdk.AccAddress
	Recipient sdk.AccAddress
	Reward    sdk.Coin
	InviteID  uint64
}

func NewMsgPayReward(sender, recipient sdk.AccAddress, reward sdk.Coin, inviteID uint64) MsgPayReward {
	return MsgPayReward{
		Sender:    sender,
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
	if len(msg.Sender) == 0 {
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
	return []sdk.AccAddress{msg.Sender}
}

func (msg MsgPayReward) GetSignBytes() []byte {
	bz := ModuleCodec.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}
