package backing

import (
	app "github.com/TruStory/truchain/types"
	"github.com/TruStory/truchain/x/stake"
	"github.com/TruStory/truchain/x/story"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// BackStoryMsg defines a message to back a story. It implements the
// `Msg` interface which is required for transactions on Cosmos blockchains.
type BackStoryMsg struct {
	stake.Msg
}

// NewBackStoryMsg creates a message to back a story
func NewBackStoryMsg(
	storyID int64,
	amount sdk.Coin,
	argument string,
	creator sdk.AccAddress) BackStoryMsg {

	// populate embedded vote msg struct
	stakeMsg := stake.Msg{
		StoryID:  storyID,
		Amount:   amount,
		Argument: argument,
		Creator:  creator,
	}

	return BackStoryMsg{stakeMsg}
}

// Route implements Msg
func (msg BackStoryMsg) Route() string { return app.GetRoute(msg) }

// Type implements Msg
func (msg BackStoryMsg) Type() string { return app.GetType(msg) }

// GetSignBytes implements Msg
func (msg BackStoryMsg) GetSignBytes() []byte {
	return app.MustGetSignBytes(msg)
}

// ValidateBasic implements Msg
func (msg BackStoryMsg) ValidateBasic() sdk.Error {
	if msg.StoryID <= 0 {
		return story.ErrInvalidStoryID(msg.StoryID)
	}
	if len(msg.Creator) == 0 {
		return sdk.ErrInvalidAddress("Invalid address: " + msg.Creator.String())
	}
	if msg.Amount.IsZero() == true {
		return sdk.ErrInsufficientFunds("Invalid backing amount: " + msg.Amount.String())
	}

	return nil
}

// GetSigners implements Msg
func (msg BackStoryMsg) GetSigners() []sdk.AccAddress {
	return app.GetSigners(msg.Creator)
}

// LikeBackingArgumentMsg represents a like on a backing message
type LikeBackingArgumentMsg struct {
	ArgumentID int64          `json:"argument_id"`
	Creator    sdk.AccAddress `json:"creator"`
	Amount     sdk.Coin       `json:"amount"`
}

// NewLikeBackingArgumentMsg constructs a new like argument message
func NewLikeBackingArgumentMsg(
	argumentID int64,
	creator sdk.AccAddress,
	amount sdk.Coin) LikeBackingArgumentMsg {

	return LikeBackingArgumentMsg{
		ArgumentID: argumentID,
		Creator:    creator,
		Amount:     amount,
	}
}

// GetSignBytes implements Msg.GetSignBytes()
func (msg LikeBackingArgumentMsg) GetSignBytes() []byte {
	return app.MustGetSignBytes(msg)
}

// GetSigners implements Msg.GetSigners()
func (msg LikeBackingArgumentMsg) GetSigners() []sdk.AccAddress {
	return app.GetSigners(msg.Creator)
}

// Route implements Msg.Route()
func (msg LikeBackingArgumentMsg) Route() string { return app.GetRoute(msg) }

// Type implements Msg.Type()
func (msg LikeBackingArgumentMsg) Type() string { return app.GetType(msg) }

// ValidateBasic implements Msg.ValidateBasic()
func (msg LikeBackingArgumentMsg) ValidateBasic() sdk.Error {
	if msg.ArgumentID == 0 {
		// return ErrInvalidArgumentID()
		return sdk.ErrInternal("Invalid argument id")
	}

	if len(msg.Creator) == 0 {
		return sdk.ErrInvalidAddress("Invalid address: " + msg.Creator.String())
	}

	if !msg.Amount.IsPositive() {
		return sdk.ErrInsufficientFunds("Invalid staking amount")
	}

	return nil
}
