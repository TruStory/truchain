package challenge

import (
	"github.com/TruStory/truchain/x/story"

	app "github.com/TruStory/truchain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// CreateChallengeMsg defines a message to challenge a story
type CreateChallengeMsg struct {
	app.VoteStoryMsg
}

// NewCreateChallengeMsg creates a message to challenge a story
func NewCreateChallengeMsg(
	storyID int64,
	amount sdk.Coin,
	argument string,
	creator sdk.AccAddress) CreateChallengeMsg {

	// populate embedded vote msg struct
	voteMsg := app.VoteStoryMsg{
		StoryID:  storyID,
		Amount:   amount,
		Argument: argument,
		Creator:  creator,
	}

	return CreateChallengeMsg{voteMsg}
}

// Route implements Msg
func (msg CreateChallengeMsg) Route() string { return app.GetRoute(msg) }

// Type implements Msg
func (msg CreateChallengeMsg) Type() string { return app.GetType(msg) }

// GetSignBytes implements Msg. Story creator should sign this message.
// Serializes Msg into JSON bytes for transport.
func (msg CreateChallengeMsg) GetSignBytes() []byte {
	return app.MustGetSignBytes(msg)
}

// ValidateBasic implements Msg
func (msg CreateChallengeMsg) ValidateBasic() sdk.Error {
	params := DefaultMsgParams()

	if msg.StoryID == 0 {
		return story.ErrInvalidStoryID(msg.StoryID)
	}
	if msg.Amount.IsZero() == true {
		return sdk.ErrInsufficientFunds("Invalid challenge amount" + msg.Amount.String())
	}
	if len := len([]rune(msg.Argument)); len > 0 && (len < params.MinArgumentLength || len > params.MaxArgumentLength) {
		return app.ErrInvalidArgumentMsg()
	}
	if len(msg.Creator) == 0 {
		return sdk.ErrInvalidAddress("Invalid address: " + msg.Creator.String())
	}

	return nil
}

// GetSigners implements Msg. Story creator is the only signer of this message.
func (msg CreateChallengeMsg) GetSigners() []sdk.AccAddress {
	return app.GetSigners(msg.Creator)
}
