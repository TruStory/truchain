package vote

import (
	app "github.com/TruStory/truchain/types"
	"github.com/TruStory/truchain/x/stake"
	"github.com/TruStory/truchain/x/story"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// CreateVoteMsg defines a message to create a vote
type CreateVoteMsg struct {
	stake.Msg

	// explicit vote
	Vote bool `json:"vote"`
}

// NewCreateVoteMsg creates a message to vote
func NewCreateVoteMsg(
	storyID int64,
	amount sdk.Coin,
	argument string,
	creator sdk.AccAddress,
	vote bool) CreateVoteMsg {

	// populate embedded vote msg struct
	voteMsg := stake.Msg{
		StoryID:  storyID,
		Amount:   amount,
		Argument: argument,
		Creator:  creator,
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
	if len := len([]rune(msg.Argument)); len > 0 && (len < params.MinArgumentLength || len > params.MaxArgumentLength) {
		return app.ErrInvalidArgumentMsg()
	}
	if len(msg.Creator) == 0 {
		return sdk.ErrInvalidAddress("Invalid address: " + msg.Creator.String())
	}
	return nil
}

// GetSigners implements Msg. Story creator is the only signer of this message.
func (msg CreateVoteMsg) GetSigners() []sdk.AccAddress {
	return app.GetSigners(msg.Creator)
}
