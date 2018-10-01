package types

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ReserveTokenDenom is the coin  denom for Trustory's native reserve token.
const ReserveTokenDenom = "trustake"

// BackStoryMsg defines a message to bond to a story.
// It implements the Cosmos `Msg` interface which is required
// for transactions on Cosmos blockchains.
type BackStoryMsg struct {
	StoryID  int64          `json:"story_id"`
	Amount   sdk.Coin       `json:"amount"`
	Creator  sdk.AccAddress `json:"creator"`
	Duration time.Duration  `json:"duration"`
}

// NewBackStoryMsg creates a message to place a new bond
func NewBackStoryMsg(
	storyID int64,
	amount sdk.Coin,
	creator sdk.AccAddress,
	duration time.Duration) BackStoryMsg {
	return BackStoryMsg{
		StoryID:  storyID,
		Amount:   amount,
		Creator:  creator,
		Duration: duration,
	}
}

// Type implements Msg
func (msg BackStoryMsg) Type() string { return "BackStory" }

// Name implements Msg
func (msg BackStoryMsg) Name() string { return msg.Type() }

// GetSignBytes implements Msg
func (msg BackStoryMsg) GetSignBytes() []byte {
	return getSignBytes(msg)
}

// ValidateBasic implements Msg
func (msg BackStoryMsg) ValidateBasic() sdk.Error {

	params := NewBackingParams()

	if msg.StoryID <= 0 {
		return ErrInvalidStoryID("StoryID cannot be negative")
	}
	if len(msg.Creator) == 0 {
		return sdk.ErrInvalidAddress("Invalid address: " + msg.Creator.String())
	}
	if msg.Amount.IsZero() == true {
		return ErrInvalidAmount("Invalid backing amount: " + msg.Amount.String())
	}
	if msg.Duration < params.MinPeriod || msg.Duration > params.MaxPeriod {
		return ErrInvalidBackingPeriod("Invalid backing duration: " + msg.Duration.String())
	}
	return nil
}

// GetSigners implements Msg
func (msg BackStoryMsg) GetSigners() []sdk.AccAddress {
	return getSigners(msg.Creator)
}
