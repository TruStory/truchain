package types

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ============================================================================

// Bond placed on a story
type Bond struct {
	ID           int64          `json:"id"`
	StoryID      int64          `json:"story_id"`
	Amount       sdk.Coin       `json:"amount"`
	CreatedBlock int64          `json:"created_block"`
	Creator      sdk.AccAddress `json:"creator"`
	Period       time.Duration  `json:"period"`
}

// NewBond creates a new bond
func NewBond(
	id int64,
	storyID int64,
	amount sdk.Coin,
	createdBlock int64,
	creator sdk.AccAddress,
	period time.Duration) Bond {
	return Bond{
		ID:           id,
		StoryID:      storyID,
		Amount:       amount,
		CreatedBlock: createdBlock,
		Creator:      creator,
		Period:       period,
	}
}

// ============================================================================

// Comment for a story
type Comment struct {
	ID      int64          `json:"id"`
	StoryID int64          `json:"story_id"`
	Body    string         `json:"body"`
	Creator sdk.AccAddress `json:"creator"`
}

// NewComment creates a new comment for a given story
func NewComment(id int64, storyID int64, body string, creator sdk.AccAddress) Comment {
	return Comment{
		ID:      id,
		StoryID: storyID,
		Body:    body,
		Creator: creator,
	}
}

// ============================================================================

// Evidence for a story
type Evidence struct {
	ID      int64          `json:"id"`
	StoryID int64          `json:"story_id"`
	Creator sdk.AccAddress `json:"creator"`
	URI     string         `json:"uri"`
}

// NewEvidence creates new evidence for a story
func NewEvidence(id int64, storyID int64, creator sdk.AccAddress, uri string) Evidence {
	return Evidence{
		ID:      id,
		StoryID: storyID,
		Creator: creator,
		URI:     uri,
	}
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

// StoryState is a type that defines a story state
type StoryState int

// List of acceptable story states
const (
	Created StoryState = iota
	Validated
	Rejected
	Unverifiable
	Challenged
	Revoked
)

// IsValid returns true if the value is listed in the enum defintion, false otherwise.
func (i StoryState) IsValid() bool {
	switch i {
	case Created, Validated, Rejected, Unverifiable, Challenged, Revoked:
		return true
	}
	return false
}

func (i StoryState) String() string {
	return [...]string{"Created", "Validated", "Rejected", "Unverifiable", "Challenged", "Revoked"}[i]
}

// StoryType is a type that defines a story type
type StoryType int

// List of acceptable story types
const (
	Default StoryType = iota
	Identity
	Recovery
)

// IsValid returns true if a story type is valid, false otherwise.
func (i StoryType) IsValid() bool {
	switch i {
	case Default, Identity, Recovery:
		return true
	}
	return false
}

func (i StoryType) String() string {
	return [...]string{"Default", "Identity", "Recovery"}[i]
}

// Story type
type Story struct {
	ID           int64            `json:"id"`
	BondIDs      []int64          `json:"bond_i_ds,omitempty"`
	CommentIDs   []int64          `json:"comment_i_ds,omitempty"`
	EvidenceIDs  []int64          `json:"evidence_i_ds,omitempty"`
	Thread       []int64          `json:"thread,omitempty"`
	VoteIDs      []int64          `json:"vote_i_ds"`
	Body         string           `json:"body"`
	Category     StoryCategory    `json:"category"`
	CreatedBlock int64            `json:"created_block"`
	Creator      sdk.AccAddress   `json:"creator"`
	Escrow       sdk.AccAddress   `json:"escrow"`
	Expiration   time.Time        `json:"expiration,omitempty"`
	Round        int64            `json:"round"`
	State        StoryState       `json:"state"`
	SubmitBlock  int64            `json:"submit_block"`
	StoryType    StoryType        `json:"type"`
	UpdatedBlock int64            `json:"updated_block"`
	Users        []sdk.AccAddress `json:"users"`
	VoteStart    time.Time        `json:"vote_start"`
	VoteEnd      time.Time        `json:"vote_end"`
}

// NewStory creates a new story
func NewStory(
	id int64,
	bondIDs []int64,
	commentIDs []int64,
	evidenceIDs []int64,
	thread []int64,
	voteIDs []int64,
	body string,
	category StoryCategory,
	createdBlock int64,
	creator sdk.AccAddress,
	escrow sdk.AccAddress,
	expiration time.Time,
	round int64,
	state StoryState,
	submitBlock int64,
	storyType StoryType,
	updatedBlock int64,
	users []sdk.AccAddress,
	voteStart time.Time,
	voteEnd time.Time) Story {
	return Story{
		ID:           id,
		BondIDs:      bondIDs,
		CommentIDs:   commentIDs,
		EvidenceIDs:  evidenceIDs,
		Thread:       thread,
		VoteIDs:      voteIDs,
		Body:         body,
		Category:     category,
		CreatedBlock: createdBlock,
		Creator:      creator,
		Escrow:       escrow,
		Expiration:   expiration,
		Round:        round,
		State:        Created,
		SubmitBlock:  submitBlock,
		StoryType:    storyType,
		UpdatedBlock: updatedBlock,
		Users:        users,
		VoteStart:    voteStart,
		VoteEnd:      voteEnd,
	}
}

// ============================================================================

// Vote for a story
type Vote struct {
	ID           int64          `json:"id"`
	StoryID      int64          `json:"story_id"`
	Amount       sdk.Coins      `json:"amount"`
	CreatedBlock int64          `json:"created_block"`
	Creator      sdk.AccAddress `json:"creator"`
	Round        int64          `json:"round"`
	Vote         bool           `json:"vote"`
}

// NewVote creates a new vote for a story
func NewVote(
	id int64,
	storyID int64,
	amount sdk.Coins,
	createdBlock int64,
	creator sdk.AccAddress,
	round int64,
	vote bool) Vote {
	return Vote{
		ID:           id,
		StoryID:      storyID,
		Amount:       amount,
		CreatedBlock: createdBlock,
		Creator:      creator,
		Round:        round,
		Vote:         vote,
	}
}

// ============================================================================

// ActiveStoryQueue is a queue of in-progress stories -- `Created` and `Challenged`
type ActiveStoryQueue []int64

// IsEmpty checks if the queue is empty
func (asq ActiveStoryQueue) IsEmpty() bool {
	if len(asq) == 0 {
		return true
	}
	return false
}
