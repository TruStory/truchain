package backing

import (
	"time"

	app "github.com/TruStory/truchain/types"
	"github.com/TruStory/truchain/x/story"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// BackStoryMsg defines a message to back a story. It implements the
// `Msg` interface which is required for transactions on Cosmos blockchains.
type BackStoryMsg struct {
	app.VoteStoryMsg

	Duration time.Duration `json:"duration"`
}

// NewBackStoryMsg creates a message to back a story
func NewBackStoryMsg(
	storyID int64,
	amount sdk.Coin,
	creator sdk.AccAddress,
	duration time.Duration) BackStoryMsg {

	// populate embedded vote msg struct
	voteMsg := app.VoteStoryMsg{
		StoryID:  storyID,
		Amount:   amount,
		Argument: "",
		Creator:  creator,
		Evidence: []string{},
	}

	return BackStoryMsg{voteMsg, duration}
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

	params := DefaultMsgParams()

	if msg.StoryID <= 0 {
		return story.ErrInvalidStoryID(msg.StoryID)
	}
	if len(msg.Creator) == 0 {
		return sdk.ErrInvalidAddress("Invalid address: " + msg.Creator.String())
	}
	if msg.Amount.IsZero() == true {
		return sdk.ErrInsufficientFunds("Invalid backing amount" + msg.Amount.String())
	}

	if msg.Duration < params.MinPeriod || msg.Duration > params.MaxPeriod {
		return ErrInvalidPeriod(msg.Duration)
	}

	return nil
}

// GetSigners implements Msg
func (msg BackStoryMsg) GetSigners() []sdk.AccAddress {
	return app.GetSigners(msg.Creator)
}
