package trubank

import (
	app "github.com/TruStory/truchain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// PayRewardMsg defines a message to pay a reward
type PayRewardMsg struct {
	Creator   sdk.AccAddress `json:"creator"`
	Recipient sdk.AccAddress `json:"recipient"`
	Reward    sdk.Coin       `json:"reward"`
	InviteID  int64          `json:"invite_id"`
}

// NewPayRewardMsg creates a new message to pay a reward
func NewPayRewardMsg(
	creator sdk.AccAddress,
	recipient sdk.AccAddress,
	reward sdk.Coin,
	inviteID int64) PayRewardMsg {

	return PayRewardMsg{
		Creator:   creator,
		Recipient: recipient,
		Reward:    reward,
		InviteID:  inviteID,
	}
}

// Route implements Msg
func (msg PayRewardMsg) Route() string { return app.GetRoute(msg) }

// Type implements Msg
func (msg PayRewardMsg) Type() string { return app.GetType(msg) }

// GetSignBytes implements Msg. Reward creator should sign this message.
// Serializes Msg into JSON bytes for transport.
func (msg PayRewardMsg) GetSignBytes() []byte {
	return app.MustGetSignBytes(msg)
}

// ValidateBasic implements Msg
func (msg PayRewardMsg) ValidateBasic() sdk.Error {
	if len(msg.Creator) == 0 {
		return sdk.ErrInvalidAddress("Invalid creator address: " + msg.Creator.String())
	}
	if len(msg.Recipient) == 0 {
		return sdk.ErrInvalidAddress("Invalid recipient address: " + msg.Recipient.String())
	}
	if !msg.Reward.IsPositive() {
		return sdk.ErrInsufficientFunds("Reward amount must be greater than zero")
	}
	if msg.Reward.Denom != app.StakeDenom {
		return sdk.ErrInvalidCoins("Rewards can only be paid in " + app.StakeDenom + ", not " + msg.Reward.Denom)
	}
	if msg.InviteID == 0 {
		return ErrTransferringCoinsToUser(msg.Recipient)
	}

	return nil
}

// GetSigners implements Msg. Reward creator is the only signer of this message.
func (msg PayRewardMsg) GetSigners() []sdk.AccAddress {
	return app.GetSigners(msg.Creator)
}
