package types

import (
	"net/url"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SubmitEvidenceMsg defines a message to submit evidence for a story
type SubmitEvidenceMsg struct {
	StoryID int64          `json:"story_id"`
	Creator sdk.AccAddress `json:"creator"`
	URI     string         `json:"url"`
}

// NewSubmitEvidenceMsg creates a new message to submit evidence for a story
func NewSubmitEvidenceMsg(storyID int64, creator sdk.AccAddress, uri string) SubmitEvidenceMsg {
	return SubmitEvidenceMsg{
		StoryID: storyID,
		Creator: creator,
		URI:     uri,
	}
}

// Type implements Msg
func (msg SubmitEvidenceMsg) Type() string { return "SubmitEvidence" }

// Name implements Msg
func (msg SubmitEvidenceMsg) Name() string { return msg.Type() }

// GetSignBytes implements Msg
func (msg SubmitEvidenceMsg) GetSignBytes() []byte {
	return getSignBytes(msg)
}

// ValidateBasic implements Msg
func (msg SubmitEvidenceMsg) ValidateBasic() sdk.Error {
	if msg.StoryID <= 0 {
		return ErrInvalidStoryID("StoryID cannot be negative")
	}
	if len(msg.Creator) == 0 {
		return sdk.ErrInvalidAddress("Invalid address: " + msg.Creator.String())
	}
	_, err := url.ParseRequestURI(msg.URI)
	if err != nil {
		return ErrInvalidURL("Invalid URL: " + msg.URI)
	}
	return nil
}

// GetSigners implements Msg
func (msg SubmitEvidenceMsg) GetSigners() []sdk.AccAddress {
	return getSigners(msg.Creator)
}
