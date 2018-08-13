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
	return "PlaceBond"
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
		return sdk.ErrInvalidAmount("Invalid bond amount: " + msg.Stake.String())
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
	return "AddComment"
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
func (msg SubmitEvidenceMsg) Type() string {
	return "SubmitEvidence"
}

// GetSignBytes implements Msg
func (msg SubmitEvidenceMsg) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return b
}

// ValidateBasic implements Msg
func (msg SubmitEvidenceMsg) ValidateBasic() types.Error {
	if msg.StoryID <= 0 {
		return ErrInvalidStoryID("StoryID cannot be negative")
	}
	if len(msg.Creator) == 0 {
		return sdk.ErrInvalidAddress("Invalid address: " + msg.Creator.String())
	}
	err := url.ParseRequestURI(msg.URI)
	if err != nil {
		return sdk.ErrInvalidURL("Invalid URL: " + msg.URI.String())
	}
	return nil
}

// GetSigners implements Msg
func (msg SubmitEvidenceMsg) GetSigners() []types.Address {
	return []sdk.AccAddress{msg.Creator}
}

// ============================================================================

// SubmitStoryMsg defines a message to submit a story
type SubmitStoryMsg struct {
	Body      string         `json:"body"`
	Category  string         `json:"category"`
	Creator   sdk.AccAddress `json:"creator"`
	StoryType string         `json:"story_type"`
}

// NewSubmitStoryMsg creates a new message to submit a story
func NewSubmitStoryMsg(body string, category string, creator sdk.AccAddress, storyType string) SubmitStoryMsg {
	return SubmitStoryMsg{
		Body:      body,
		Category:  category,
		Creator:   creator,
		StoryType: storyType,
	}
}

// Type implements Msg
func (msg SubmitStoryMsg) Type() string {
	return "SubmitStory"
}

// GetSignBytes implements Msg
func (msg SubmitStoryMsg) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return b
}

// ValidateBasic implements Msg
func (msg SubmitStoryMsg) ValidateBasic() types.Error {
	if len(msg.Body) == 0 {
		return sdk.ErrInvalidBody("Invalid body: " + msg.Body.String())
	}
	if len(msg.Category) == 0 {
		return sdk.ErrInvalidCategory("Invalid category: " + msg.Category.String())
	}
	if len(msg.Creator) == 0 {
		return sdk.ErrInvalidAddress("Invalid address: " + msg.Creator.String())
	}
	if len(msg.StoryType) == 0 {
		return sdk.ErrInvalidStoryType("Invalid story type: " + msg.StoryType.String())
	}
	return nil
}

// GetSigners implements Msg
func (msg SubmitStoryMsg) GetSigners() []types.Address {
	return []sdk.AccAddress{msg.Creator}
}

// ============================================================================

// VoteMsg defines a message to vote on a story
type VoteMsg struct {
	StoryID int64
	Creator sdk.AccAddress
	Stake   sdk.Coin
	Vote    bool
}

// NewVoteMsg creates a new message to vote on a story
func NewVoteMsg(storyID int64, creator sdk.AccAddress, stake sdk.Coin, vote bool) VoteMsg {
	return VoteMsg{
		StoryID: storyID,
		Creator: creator,
		Stake:   stake,
		Vote:    vote,
	}
}

// Type implements Msg
func (msg VoteMsg) Type() string {
	return "VoteMsg"
}

// GetSignBytes implements Msg
func (msg VoteMsg) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return b
}

// ValidateBasic implements Msg
func (msg VoteMsg) ValidateBasic() types.Error {
	if msg.StoryID <= 0 {
		return ErrInvalidStoryID("StoryID cannot be negative")
	}
	if len(msg.Creator) == 0 {
		return sdk.ErrInvalidAddress("Invalid address: " + msg.Creator.String())
	}
	if msg.Stake.IsValid == false {
		return sdk.ErrInvalidAmount("Invalid stake amount: " + msg.Stake.String())
	}
	return nil
}

// GetSigners implements Msg
func (msg VoteMsg) GetSigners() []types.Address {
	return []sdk.AccAddress{msg.Creator}
}
