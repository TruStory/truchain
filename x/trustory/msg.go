package trustory

import (
	"encoding/json"
	"net/url"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ============================================================================

// PlaceBondMsg defines a message to bond to a story
type PlaceBondMsg struct {
	StoryID int64          `json:"story_id"`
	Amount  sdk.Coin       `json:"amount"`
	Creator sdk.AccAddress `json:"creator"`
	Period  time.Duration  `json:"period"`
}

// NewPlaceBondMsg creates a message to place a new bond
func NewPlaceBondMsg(
	storyID int64,
	amount sdk.Coin,
	creator sdk.AccAddress,
	period time.Duration) PlaceBondMsg {
	return PlaceBondMsg{
		StoryID: storyID,
		Amount:  amount,
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
	return getSignBytes(msg)
}

// ValidateBasic implements Msg
func (msg PlaceBondMsg) ValidateBasic() sdk.Error {
	if msg.StoryID <= 0 {
		return ErrInvalidStoryID("StoryID cannot be negative")
	}
	if len(msg.Creator) == 0 {
		return sdk.ErrInvalidAddress("Invalid address: " + msg.Creator.String())
	}
	if msg.Amount.IsZero() == true {
		return ErrInvalidAmount("Invalid bond amount: " + msg.Amount.String())
	}
	if msg.Period == 0 {
		return ErrInvalidBondPeriod("Invalid bond period: " + msg.Period.String())
	}
	return nil
}

// GetSigners implements Msg
func (msg PlaceBondMsg) GetSigners() []sdk.AccAddress {
	return getSigners(msg.Creator)
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

// ============================================================================

// StoryCategory is a type that defines a story category
type StoryCategory int

// List of accepted categories
const (
	Unknown StoryCategory = iota
	Bitcoin
	Consensus
	DEX
	Ethereum
	StableCoins
)

// IsValid returns true if the value is listed in the enum defintion, false otherwise.
func (i StoryCategory) IsValid() bool {
	switch i {
	case Unknown, Bitcoin, Consensus, DEX, Ethereum, StableCoins:
		return true
	}
	return false
}

// Slug is the short name for a category
func (i StoryCategory) Slug() string {
	return [...]string{"unknown", "btc", "consensus", "dex", "eth", "stablecoins"}[i]
}

func (i StoryCategory) String() string {
	return [...]string{"Unknown", "Bitcoin", "Consensus", "Decentralized Exchanges", "Ethereum", "Stable Coins"}[i]
}

// SubmitStoryMsg defines a message to submit a story
type SubmitStoryMsg struct {
	Body      string         `json:"body"`
	Category  StoryCategory  `json:"category"`
	Creator   sdk.AccAddress `json:"creator"`
	StoryType string         `json:"story_type"`
}

// NewSubmitStoryMsg creates a new message to submit a story
func NewSubmitStoryMsg(body string, category StoryCategory, creator sdk.AccAddress, storyType string) SubmitStoryMsg {
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
	return getSignBytes(msg)
}

// ValidateBasic implements Msg
func (msg SubmitStoryMsg) ValidateBasic() sdk.Error {
	if len(msg.Body) == 0 {
		return ErrInvalidBody("Invalid body: " + msg.Body)
	}
	if msg.Category.IsValid() == false {
		return ErrInvalidCategory("Invalid category: " + msg.Category.String())
	}
	if len(msg.Creator) == 0 {
		return sdk.ErrInvalidAddress("Invalid address: " + msg.Creator.String())
	}
	if len(msg.StoryType) == 0 {
		return ErrInvalidStoryType("Invalid story type: " + msg.StoryType)
	}
	return nil
}

// GetSigners implements Msg
func (msg SubmitStoryMsg) GetSigners() []sdk.AccAddress {
	return getSigners(msg.Creator)
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
	return "Vote"
}

// GetSignBytes implements Msg
func (msg VoteMsg) GetSignBytes() []byte {
	return getSignBytes(msg)
}

// ValidateBasic implements Msg
func (msg VoteMsg) ValidateBasic() sdk.Error {
	if msg.StoryID <= 0 {
		return ErrInvalidStoryID("StoryID cannot be negative")
	}
	if len(msg.Creator) == 0 {
		return sdk.ErrInvalidAddress("Invalid address: " + msg.Creator.String())
	}
	if msg.Stake.IsZero() == true {
		return ErrInvalidAmount("Invalid stake amount: " + msg.Stake.String())
	}
	return nil
}

// GetSigners implements Msg
func (msg VoteMsg) GetSigners() []sdk.AccAddress {
	return getSigners(msg.Creator)
}

// ============================================================================

func getSignBytes(msg sdk.Msg) []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return b
}

func getSigners(addr sdk.AccAddress) []sdk.AccAddress {
	return []sdk.AccAddress{addr}
}
