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
	StoryID  int64          `json:"story_id"`
	Amount   sdk.Coin       `json:"amount"`
	Creator  sdk.AccAddress `json:"creator"`
	Duration int64          `json:"duration"`
}

// NewBackStoryMsg creates a message to back a story
func NewBackStoryMsg(
	storyID int64,
	amount sdk.Coin,
	creator sdk.AccAddress,
	duration time.Duration) BackStoryMsg {
	return BackStoryMsg{
		StoryID:  storyID,
		Amount:   amount,
		Creator:  creator,
		Duration: int64(duration),
	}
}

// Type implements Msg
func (msg BackStoryMsg) Type() string { return app.GetType(msg) }

// Name implements Msg
func (msg BackStoryMsg) Name() string { return app.GetName(msg) }

// GetSignBytes implements Msg
func (msg BackStoryMsg) GetSignBytes() []byte {
	return app.MustGetSignBytes(msg)
}

func (msg BackStoryMsg) GetDuration() time.Duration {
	return time.Duration(time.Duration(msg.Duration) * time.Second)
}

// ValidateBasic implements Msg
func (msg BackStoryMsg) ValidateBasic() sdk.Error {

	params := DefaultParams()

	if msg.StoryID <= 0 {
		return story.ErrInvalidStoryID(msg.StoryID)
	}
	if len(msg.Creator) == 0 {
		return sdk.ErrInvalidAddress("Invalid address: " + msg.Creator.String())
	}
	if msg.Amount.IsZero() == true {
		return sdk.ErrInsufficientFunds("Invalid backing amount" + msg.Amount.String())
	}

	duration := time.Duration(time.Duration(msg.Duration) * time.Second)

	if duration < params.MinPeriod || duration > params.MaxPeriod {
		return ErrInvalidPeriod(duration)
	}

	return nil
}

// GetSigners implements Msg
func (msg BackStoryMsg) GetSigners() []sdk.AccAddress {
	return app.GetSigners(msg.Creator)
}
