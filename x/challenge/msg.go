package challenge

import (
	"github.com/TruStory/truchain/x/story"

	app "github.com/TruStory/truchain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// StartChallengeMsg defines a message to challenge a story
type StartChallengeMsg struct {
	StoryID  int64          `json:"story_id"`
	Amount   sdk.Coin       `json:"amount"`
	Argument string         `json:"argument,omitempty"`
	Creator  sdk.AccAddress `json:"creator"`
	Evidence string         `json:"evidence,omitempty"`
}

// NewStartChallengeMsg creates a message to challenge a story
func NewStartChallengeMsg(
	storyID int64,
	amount sdk.Coin,
	argument string,
	creator sdk.AccAddress,
	evidence string) StartChallengeMsg {
	return StartChallengeMsg{
		StoryID:  storyID,
		Amount:   amount,
		Argument: argument,
		Creator:  creator,
		Evidence: evidence,
	}
}

// Type implements Msg
func (msg StartChallengeMsg) Type() string { return app.GetType(msg) }

// Name implements Msg
func (msg StartChallengeMsg) Name() string { return app.GetName(msg) }

// GetSignBytes implements Msg. Story creator should sign this message.
// Serializes Msg into JSON bytes for transport.
func (msg StartChallengeMsg) GetSignBytes() []byte {
	return app.MustGetSignBytes(msg)
}

// ValidateBasic implements Msg
func (msg StartChallengeMsg) ValidateBasic() sdk.Error {
	// if len(msg.Body) == 0 {
	// 	return ErrInvalidStoryBody(msg.Body)
	// }
	if msg.StoryID == 0 {
		return story.ErrInvalidStoryID(msg.StoryID)
	}
	if len(msg.Creator) == 0 {
		return sdk.ErrInvalidAddress("Invalid address: " + msg.Creator.String())
	}
	// if msg.Kind.IsValid() == false {
	// 	return ErrInvalidStoryKind(msg.Kind.String())
	// }
	return nil
}

// GetSigners implements Msg. Story creator is the only signer of this message.
func (msg StartChallengeMsg) GetSigners() []sdk.AccAddress {
	return app.GetSigners(msg.Creator)
}
