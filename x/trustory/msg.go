package trustory

import (
	"encoding/json"
	"net/url"
	"time"

	"github.com/TruStory/cosmos-sdk/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ============================================================================

// PlaceBondMsg defines a message to bond to a story
type PlaceBondMsg struct {
	StoryID int64          `json:"story_id"`
	Stake   sdk.Coin       `json:"stake"`
	Creator sdk.AccAddress `json:"creator"`
	Period  time.Time      `json:"period"`
}

// NewPlaceBondMsg creates a message to place a new bond
func NewPlaceBondMsg(
	storyID int64,
	stake sdk.Coin,
	creator sdk.AccAddress,
	period time.Time) PlaceBondMsg {
	return PlaceBondMsg{
		StoryID: storyID,
		Stake:   stake,
		Creator: creator,
		Period:  period,
	}
}

// Type implements Msg
func (msg PlaceBondMsg) Type() string {
	return "truStory"
}

// GetSignBytes implements Msg
func (msg PlaceBondMsg) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return b
}

// ValidateBasic implements Msg
func (msg PlaceBondMsg) ValidateBasic() types.Error {
	if msg.StoryID <= 0 {
		return ErrInvalidStoryID("StoryID cannot be negative")
	}
	if len(msg.Creator) == 0 {
		return sdk.ErrInvalidAddress("Invalid address: " + msg.Creator.String())
	}
	if msg.Stake.IsValid == false {
		return sdk.ErrInvalidBondAmount("Invalid bond amount: " + msg.Stake.String())
	}
	if msg.Period.IsZero == true {
		return sdk.ErrInvalidBondPeriod("Invalid bond period: " + msg.Period.String())
	}
	return nil
}

// GetSigners implements Msg
func (msg PlaceBondMsg) GetSigners() []types.Address {
	return []sdk.AccAddress{msg.Creator}
}

// ============================================================================

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
func (msg AddCommentMsg) Type() string {
	return "truStory"
}

// GetSignBytes implements Msg
func (msg AddCommentMsg) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return b
}

// ValidateBasic implements Msg
func (msg AddCommentMsg) ValidateBasic() types.Error {
	if msg.StoryID <= 0 {
		return ErrInvalidStoryID("StoryID cannot be negative")
	}
	if len(msg.Creator) == 0 {
		return sdk.ErrInvalidAddress("Invalid address: " + msg.Creator.String())
	}
	if len(msg.Body) == 0 {
		return sdk.ErrInvalidBody("Invalid comment body: " + msg.Body.String())
	}
	return nil
}

// GetSigners implements Msg
func (msg AddCommentMsg) GetSigners() []types.Address {
	return []sdk.AccAddress{msg.Creator}
}

// ============================================================================

// SubmitEvidenceMsg defines a message to submit evidence for a story
type SubmitEvidenceMsg struct {
	StoryID int64          `json:"story_id"`
	Creator sdk.AccAddress `json:"creator"`
	URL     url.URL        `json:"url"`
}

// NewSubmitEvidenceMsg creates a new message to submit evidence for a story
func NewSubmitEvidenceMsg(storyID int64, creator sdk.AccAddress, url url.URL) SubmitEvidenceMsg {
	return SubmitEvidenceMsg{
		StoryID: storyID,
		Creator: creator,
		URL:     url,
	}
}
