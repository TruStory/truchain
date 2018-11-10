package vote

import (
	"net/url"

	app "github.com/TruStory/truchain/types"
	"github.com/TruStory/truchain/x/story"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// CreateVoteMsg defines a message to create a vote
type CreateVoteMsg struct {
	StoryID  int64          `json:"story_id"`
	Amount   sdk.Coin       `json:"amount"`
	Comment  string         `json:"comment,omitempty"`
	Creator  sdk.AccAddress `json:"creator"`
	Evidence []url.URL      `json:"evidence,omitempty"`
	Vote     bool           `json:"vote"`
}

// NewCreateVoteMsg creates a message to vote
func NewCreateVoteMsg(
	storyID int64,
	amount sdk.Coin,
	comment string,
	creator sdk.AccAddress,
	evidence []url.URL,
	vote bool) CreateVoteMsg {

	return CreateVoteMsg{
		StoryID:  storyID,
		Amount:   amount,
		Comment:  comment,
		Creator:  creator,
		Evidence: evidence,
		Vote:     vote,
	}
}

// Type implements Msg
func (msg CreateVoteMsg) Type() string { return app.GetType(msg) }

// Route implements Msg
func (msg CreateVoteMsg) Route() string { return app.GetName(msg) }

// GetSignBytes implements Msg. Story creator should sign this message.
// Serializes Msg into JSON bytes for transport.
func (msg CreateVoteMsg) GetSignBytes() []byte {
	return app.MustGetSignBytes(msg)
}

// ValidateBasic implements Msg
func (msg CreateVoteMsg) ValidateBasic() sdk.Error {
	params := DefaultParams()

	if msg.StoryID == 0 {
		return story.ErrInvalidStoryID(msg.StoryID)
	}
	if msg.Amount.IsZero() == true {
		return sdk.ErrInsufficientFunds("Invalid vote amount" + msg.Amount.String())
	}
	if len := len(msg.Comment); len < params.MinCommentLength || len > params.MaxCommentLength {
		return app.ErrInvalidCommentMsg()
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
