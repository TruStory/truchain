package story

import (
	app "github.com/TruStory/truchain/types"
	"github.com/TruStory/truchain/x/category"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SubmitStoryMsg defines a message to submit a story
type SubmitStoryMsg struct {
	Argument   string         `json:"argument,omitempty"`
	Body       string         `json:"body"`
	CategoryID int64          `json:"category_id"`
	Creator    sdk.AccAddress `json:"creator"`
	Source     string         `json:"source,omitempty"`
	StoryType  Type           `json:"story_type"`
}

// NewSubmitStoryMsg creates a new message to submit a story
func NewSubmitStoryMsg(
	argument string,
	body string,
	categoryID int64,
	creator sdk.AccAddress,
	source string,
	storyType Type) SubmitStoryMsg {

	return SubmitStoryMsg{
		Argument:   argument,
		Body:       body,
		CategoryID: categoryID,
		Creator:    creator,
		Source:     source,
		StoryType:  storyType,
	}
}

// Route implements Msg
func (msg SubmitStoryMsg) Route() string { return app.GetRoute(msg) }

// Type implements Msg
func (msg SubmitStoryMsg) Type() string { return app.GetType(msg) }

// GetSignBytes implements Msg. Story creator should sign this message.
// Serializes Msg into JSON bytes for transport.
func (msg SubmitStoryMsg) GetSignBytes() []byte {
	return app.MustGetSignBytes(msg)
}

// ValidateBasic implements Msg
func (msg SubmitStoryMsg) ValidateBasic() sdk.Error {
	if len(msg.Body) == 0 {
		return ErrInvalidStoryBody(msg.Body)
	}
	if msg.CategoryID == 0 {
		return category.ErrInvalidCategory(msg.CategoryID)
	}
	if len(msg.Creator) == 0 {
		return sdk.ErrInvalidAddress("Invalid address: " + msg.Creator.String())
	}
	if msg.StoryType.IsValid() == false {
		return ErrInvalidStoryType(msg.StoryType.String())
	}
	if len(msg.Argument) == 0 {
		return ErrInvalidStoryArgument(msg.Argument)
	}
	return nil
}

// GetSigners implements Msg. Story creator is the only signer of this message.
func (msg SubmitStoryMsg) GetSigners() []sdk.AccAddress {
	return app.GetSigners(msg.Creator)
}

// ============================================================================

// FlagStoryMsg defines a message to flag a story
type FlagStoryMsg struct {
	StoryID int64          `json:"story_id"`
	Creator sdk.AccAddress `json:"creator"`
}

// NewFlagStoryMsg creates a new message to flag a story
func NewFlagStoryMsg(storyID int64, creator sdk.AccAddress) FlagStoryMsg {
	return FlagStoryMsg{
		StoryID: storyID,
		Creator: creator,
	}
}

// Route implements Msg.Route
func (msg FlagStoryMsg) Route() string { return app.GetRoute(msg) }

// Type implements Msg.Type
func (msg FlagStoryMsg) Type() string { return app.GetType(msg) }

// GetSignBytes implements Msg.GetSignBytes
func (msg FlagStoryMsg) GetSignBytes() []byte {
	return app.MustGetSignBytes(msg)
}

// ValidateBasic implements Msg.ValidateBasic
func (msg FlagStoryMsg) ValidateBasic() sdk.Error {
	if msg.StoryID <= 0 {
		return ErrInvalidStoryID(msg.StoryID)
	}
	if len(msg.Creator) == 0 {
		return sdk.ErrInvalidAddress("Invalid address: " + msg.Creator.String())
	}

	return nil
}

// GetSigners implements Msg.GetSigners
func (msg FlagStoryMsg) GetSigners() []sdk.AccAddress {
	return app.GetSigners(msg.Creator)
}
