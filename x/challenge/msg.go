package challenge

import (
	"net/url"

	"github.com/TruStory/truchain/x/story"

	app "github.com/TruStory/truchain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SubmitChallengeMsg defines a message to challenge a story
type SubmitChallengeMsg struct {
	StoryID  int64          `json:"story_id"`
	Amount   sdk.Coin       `json:"amount"`
	Argument string         `json:"argument,omitempty"`
	Creator  sdk.AccAddress `json:"creator"`
	Evidence []url.URL      `json:"evidence,omitempty"`
}

// NewSubmitChallengeMsg creates a message to challenge a story
func NewSubmitChallengeMsg(
	storyID int64,
	amount sdk.Coin,
	argument string,
	creator sdk.AccAddress,
	evidence []url.URL) SubmitChallengeMsg {
	return SubmitChallengeMsg{
		StoryID:  storyID,
		Amount:   amount,
		Argument: argument,
		Creator:  creator,
		Evidence: evidence,
	}
}

// Type implements Msg
func (msg SubmitChallengeMsg) Type() string { return app.GetType(msg) }

// Route implements Msg
func (msg SubmitChallengeMsg) Route() string { return app.GetName(msg) }

// GetSignBytes implements Msg. Story creator should sign this message.
// Serializes Msg into JSON bytes for transport.
func (msg SubmitChallengeMsg) GetSignBytes() []byte {
	return app.MustGetSignBytes(msg)
}

// ValidateBasic implements Msg
func (msg SubmitChallengeMsg) ValidateBasic() sdk.Error {
	params := DefaultParams()

	if msg.StoryID == 0 {
		return story.ErrInvalidStoryID(msg.StoryID)
	}
	if msg.Amount.IsZero() == true {
		return sdk.ErrInsufficientFunds("Invalid challenge amount" + msg.Amount.String())
	}
	if len := len(msg.Argument); len < params.MinArgumentLength || len > params.MaxArgumentLength {
		return ErrInvalidMsg(msg.Argument)
	}
	if len(msg.Creator) == 0 {
		return sdk.ErrInvalidAddress("Invalid address: " + msg.Creator.String())
	}
	if len := len(msg.Evidence); len < params.MinEvidenceCount || len > params.MaxEvidenceCount {
		return ErrInvalidMsg(msg.Evidence)
	}
	return nil
}

// GetSigners implements Msg. Story creator is the only signer of this message.
func (msg SubmitChallengeMsg) GetSigners() []sdk.AccAddress {
	return app.GetSigners(msg.Creator)
}
