package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SubmitStoryMsg defines a message to submit a story
type SubmitStoryMsg struct {
	Body       string         `json:"body"`
	CategoryID int64          `json:"category_id"`
	Creator    sdk.AccAddress `json:"creator"`
	StoryType  StoryType      `json:"story_type"`
}

// NewSubmitStoryMsg creates a new message to submit a story
func NewSubmitStoryMsg(body string, categoryID int64, creator sdk.AccAddress, storyType StoryType) SubmitStoryMsg {
	return SubmitStoryMsg{
		Body:       body,
		CategoryID: categoryID,
		Creator:    creator,
		StoryType:  storyType,
	}
}

// Type implements Msg
func (msg SubmitStoryMsg) Type() string { return MsgType }

// Name implements Msg
func (msg SubmitStoryMsg) Name() string { return "submit_story" }

// GetSignBytes implements Msg
func (msg SubmitStoryMsg) GetSignBytes() []byte {
	return getSignBytes(msg)
}

// ValidateBasic implements Msg
func (msg SubmitStoryMsg) ValidateBasic() sdk.Error {
	if len(msg.Body) == 0 {
		return ErrInvalidBody("Invalid body: " + msg.Body)
	}
	if msg.CategoryID == 0 {
		return ErrInvalidCategory(msg.CategoryID)
	}
	if len(msg.Creator) == 0 {
		return sdk.ErrInvalidAddress("Invalid address: " + msg.Creator.String())
	}
	if msg.StoryType.IsValid() == false {
		return ErrInvalidStoryType("Invalid story type: " + msg.StoryType.String())
	}
	return nil
}

// GetSigners implements Msg
func (msg SubmitStoryMsg) GetSigners() []sdk.AccAddress {
	return getSigners(msg.Creator)
}
