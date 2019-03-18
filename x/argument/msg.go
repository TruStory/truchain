package argument

import (
	app "github.com/TruStory/truchain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// LikeArgumentMsg sends a message to like an argument
type LikeArgumentMsg struct {
	ArgumentID int64          `json:"argument_id"`
	Creator    sdk.AccAddress `json:"creator"`
	Amount     sdk.Coin       `json:"amount"`
}

// NewLikeArgumentMsg constructs a new like argument message
func NewLikeArgumentMsg(
	argumentID int64,
	creator sdk.AccAddress,
	amount sdk.Coin) LikeArgumentMsg {

	return LikeArgumentMsg{
		ArgumentID: argumentID,
		Creator:    creator,
		Amount:     amount,
	}
}

// GetSignBytes implements Msg.GetSignBytes()
func (msg LikeArgumentMsg) GetSignBytes() []byte {
	return app.MustGetSignBytes(msg)
}

// GetSigners implements Msg.GetSigners()
func (msg LikeArgumentMsg) GetSigners() []sdk.AccAddress {
	return app.GetSigners(msg.Creator)
}

// Route implements Msg.Route()
func (msg LikeArgumentMsg) Route() string { return app.GetRoute(msg) }

// Type implements Msg.Type()
func (msg LikeArgumentMsg) Type() string { return app.GetType(msg) }

// ValidateBasic implements Msg.ValidateBasic()
func (msg LikeArgumentMsg) ValidateBasic() sdk.Error {
	if msg.ArgumentID == 0 {
		return ErrInvalidArgumentID()
	}

	if len(msg.Creator) == 0 {
		return sdk.ErrInvalidAddress("Invalid address: " + msg.Creator.String())
	}

	if !msg.Amount.IsPositive() {
		return sdk.ErrInsufficientFunds("Invalid staking amount")
	}

	return nil
}
