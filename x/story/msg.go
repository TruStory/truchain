package story

import (
	"net/url"

	c "github.com/TruStory/truchain/x/category"

	t "github.com/TruStory/truchain/types"
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

// Type implements Msg
func (msg SubmitStoryMsg) Type() string { return t.GetType(msg) }

// Name implements Msg
func (msg SubmitStoryMsg) Name() string { return t.GetName(msg) }

// GetSignBytes implements Msg. Story creator should sign this message.
// Serializes Msg into JSON bytes for transport.
func (msg SubmitStoryMsg) GetSignBytes() []byte {
	return t.MustGetSignBytes(msg)
}

// ValidateBasic implements Msg
func (msg SubmitStoryMsg) ValidateBasic() sdk.Error {
	if len(msg.Body) == 0 {
		return ErrInvalidStoryBody(msg.Body)
	}
	if msg.CategoryID == 0 {
		return c.ErrInvalidCategory(msg.CategoryID)
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
	return t.GetSigners(msg.Creator)
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

// Type implements Msg
func (msg SubmitEvidenceMsg) Type() string { return t.GetType(msg) }

// Name implements Msg
func (msg SubmitEvidenceMsg) Name() string { return t.GetName(msg) }

// GetSignBytes implements Msg
func (msg SubmitEvidenceMsg) GetSignBytes() []byte {
	return t.MustGetSignBytes(msg)
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
	return t.GetSigners(msg.Creator)
}
