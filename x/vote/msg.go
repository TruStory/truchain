package vote

import (
	app "github.com/TruStory/truchain/types"
	"github.com/TruStory/truchain/x/story"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// CreateVoteMsg defines a message to create a vote
type CreateVoteMsg struct {
	app.VoteStoryMsg

	// explicit vote
	Vote bool `json:"vote"`
}

// NewCreateVoteMsg creates a message to vote
func NewCreateVoteMsg(
	storyID int64,
	amount sdk.Coin,
	argument string,
	creator sdk.AccAddress,
	evidence []string,
	vote bool) CreateVoteMsg {

	// populate embedded vote msg struct
	voteMsg := app.VoteStoryMsg{
		StoryID:  storyID,
		Amount:   amount,
		Argument: argument,
		Creator:  creator,
		Evidence: []string{},
	}

	return CreateVoteMsg{voteMsg, vote}
}

// Route implements Msg
func (msg CreateVoteMsg) Route() string { return app.GetRoute(msg) }

// Type implements Msg
func (msg CreateVoteMsg) Type() string { return app.GetType(msg) }

// GetSignBytes implements Msg. Story creator should sign this message.
// Serializes Msg into JSON bytes for transport.
func (msg CreateVoteMsg) GetSignBytes() []byte {
	return app.MustGetSignBytes(msg)
}

// ValidateBasic implements Msg
func (msg CreateVoteMsg) ValidateBasic() sdk.Error {
	params := DefaultMsgParams()

	if msg.StoryID == 0 {
		return story.ErrInvalidStoryID(msg.StoryID)
	}
	if msg.Amount.IsZero() == true {
		return sdk.ErrInsufficientFunds("Invalid vote amount" + msg.Amount.String())
	}
	if len := len(msg.Argument); len < params.MinArgumentLength || len > params.MaxArgumentLength {
		return app.ErrInvalidArgumentMsg()
	}
	if len(msg.Creator) == 0 {
		return sdk.ErrInvalidAddress("Invalid address: " + msg.Creator.String())
	}
	if len := len(msg.Evidence); len < params.MinEvidenceCount || len > params.MaxEvidenceCount {
		return app.ErrInvalidEvidenceMsg()
	}
	return nil
}

// GetSigners implements Msg. Story creator is the only signer of this message.
func (msg CreateVoteMsg) GetSigners() []sdk.AccAddress {
	return app.GetSigners(msg.Creator)
}
