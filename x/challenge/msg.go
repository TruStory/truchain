package challenge

import (
	"net/url"

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
	Evidence []url.URL      `json:"evidence,omitempty"`
	Reason   Reason         `json:"reason"`
}

// NewStartChallengeMsg creates a message to challenge a story
func NewStartChallengeMsg(
	storyID int64,
	amount sdk.Coin,
	argument string,
	creator sdk.AccAddress,
	evidence []url.URL,
	reason Reason) StartChallengeMsg {
	return StartChallengeMsg{
		StoryID:  storyID,
		Amount:   amount,
		Argument: argument,
		Creator:  creator,
		Evidence: evidence,
		Reason:   reason,
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
	params := NewParams()

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
	if !msg.Reason.IsValid() {
		return ErrInvalidMsg(msg.Reason)
	}
	return nil
}

// GetSigners implements Msg. Story creator is the only signer of this message.
func (msg StartChallengeMsg) GetSigners() []sdk.AccAddress {
	return app.GetSigners(msg.Creator)
}

// ============================================================================

// UpdateChallengeMsg defines a message to challenge a story
type UpdateChallengeMsg struct {
	ChallengeID int64          `json:"challenge_id"`
	Amount      sdk.Coin       `json:"amount"`
	Creator     sdk.AccAddress `json:"creator"`
}

// NewUpdateChallengeMsg creates a message to challenge a story
func NewUpdateChallengeMsg(
	challengeID int64,
	amount sdk.Coin,
	creator sdk.AccAddress) UpdateChallengeMsg {
	return UpdateChallengeMsg{
		ChallengeID: challengeID,
		Amount:      amount,
		Creator:     creator,
	}
}

// Type implements Msg
func (msg UpdateChallengeMsg) Type() string { return app.GetType(msg) }

// Name implements Msg
func (msg UpdateChallengeMsg) Name() string { return app.GetName(msg) }

// GetSignBytes implements Msg. Story creator should sign this message.
// Serializes Msg into JSON bytes for transport.
func (msg UpdateChallengeMsg) GetSignBytes() []byte {
	return app.MustGetSignBytes(msg)
}

// ValidateBasic implements Msg
func (msg UpdateChallengeMsg) ValidateBasic() sdk.Error {
	if msg.ChallengeID == 0 {
		return story.ErrInvalidStoryID(msg.ChallengeID)
	}
	if msg.Amount.IsZero() == true {
		return sdk.ErrInvalidCoins("Invalid challenge amount" + msg.Amount.String())
	}
	if len(msg.Creator) == 0 {
		return sdk.ErrInvalidAddress("Invalid address: " + msg.Creator.String())
	}
	return nil
}

// GetSigners implements Msg. Story creator is the only signer of this message.
func (msg UpdateChallengeMsg) GetSigners() []sdk.AccAddress {
	return app.GetSigners(msg.Creator)
}
