package story

import (
	"net/url"

	app "github.com/TruStory/truchain/types"
	"github.com/TruStory/truchain/x/category"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SubmitStoryMsg defines a message to submit a story
type SubmitStoryMsg struct {
	Body       string         `json:"body"`
	CategoryID int64          `json:"category_id"`
	Creator    sdk.AccAddress `json:"creator"`
	Kind       Kind           `json:"kind"`
}

// NewSubmitStoryMsg creates a new message to submit a story
func NewSubmitStoryMsg(body string, categoryID int64, creator sdk.AccAddress, kind Kind) SubmitStoryMsg {
	return SubmitStoryMsg{
		Body:       body,
		CategoryID: categoryID,
		Creator:    creator,
		Kind:       kind,
	}
}

// Route implements Msg
func (msg SubmitStoryMsg) Route() string { return app.GetType(msg) }

// Type implements Msg
func (msg SubmitStoryMsg) Type() string { return app.GetName(msg) }

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
	if msg.Kind.IsValid() == false {
		return ErrInvalidStoryKind(msg.Kind.String())
	}
	return nil
}

// GetSigners implements Msg. Story creator is the only signer of this message.
func (msg SubmitStoryMsg) GetSigners() []sdk.AccAddress {
	return app.GetSigners(msg.Creator)
}

// ============================================================================

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

// Route implements Msg
func (msg SubmitEvidenceMsg) Route() string { return app.GetName(msg) }

// Type implements Msg
func (msg SubmitEvidenceMsg) Type() string { return app.GetType(msg) }

// GetSignBytes implements Msg
func (msg SubmitEvidenceMsg) GetSignBytes() []byte {
	return app.MustGetSignBytes(msg)
}

// ValidateBasic implements Msg
func (msg SubmitEvidenceMsg) ValidateBasic() sdk.Error {
	if msg.StoryID <= 0 {
		return ErrInvalidStoryID(msg.StoryID)
	}
	if len(msg.Creator) == 0 {
		return sdk.ErrInvalidAddress("Invalid address: " + msg.Creator.String())
	}
	_, err := url.ParseRequestURI(msg.URI)
	if err != nil {
		return ErrInvalidEvidenceURL(msg.URI)
	}
	return nil
}

// GetSigners implements Msg
func (msg SubmitEvidenceMsg) GetSigners() []sdk.AccAddress {
	return app.GetSigners(msg.Creator)
}
