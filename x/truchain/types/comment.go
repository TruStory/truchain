package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// AddCommentMsg defines a message to add a comment to a story
type AddCommentMsg struct {
	StoryID int64          `json:"story_id"`
	Body    string         `json:"body"`
	Creator sdk.AccAddress `json:"creator"`
}

// NewAddCommentMsg creates a message to add a new comment to a story
func NewAddCommentMsg(storyID int64, body string, creator sdk.AccAddress) AddCommentMsg {
	return AddCommentMsg{
		StoryID: storyID,
		Body:    body,
		Creator: creator,
	}
}

// Type implements Msg
func (msg AddCommentMsg) Type() string { return "AddComment" }

// Name implements Msg
func (msg AddCommentMsg) Name() string { return msg.Type() }

// GetSignBytes implements Msg
func (msg AddCommentMsg) GetSignBytes() []byte {
	return getSignBytes(msg)
}

// ValidateBasic implements Msg
func (msg AddCommentMsg) ValidateBasic() sdk.Error {
	if msg.StoryID <= 0 {
		return ErrInvalidStoryID("StoryID cannot be negative")
	}
	if len(msg.Creator) == 0 {
		return sdk.ErrInvalidAddress("Invalid address: " + msg.Creator.String())
	}
	if len(msg.Body) == 0 {
		return ErrInvalidBody("Invalid comment body: " + msg.Body)
	}
	return nil
}

// GetSigners implements Msg
func (msg AddCommentMsg) GetSigners() []sdk.AccAddress {
	return getSigners(msg.Creator)
}
