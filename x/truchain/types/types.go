package types

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ============================================================================

// BackingQueue is a list of all backings
type BackingQueue []int64

// IsEmpty checks if the queue is empty
func (asq BackingQueue) IsEmpty() bool {
	if len(asq) == 0 {
		return true
	}
	return false
}

// ============================================================================

// Backing type
type Backing struct {
	ID      int64          `json:"id"`
	StoryID int64          `json:"story_id"`
	Amount  sdk.Coin       `json:"amount"`
	Expires time.Time      `json:"expires"`
	User    sdk.AccAddress `json:"user"`
}

// NewBacking creates a new backing type
func NewBacking(
	id int64,
	storyID int64,
	amount sdk.Coin,
	expires time.Time,
	creator sdk.AccAddress) Backing {

	return Backing{
		ID:      id,
		StoryID: storyID,
		Amount:  amount,
		Expires: expires,
		User:    creator,
	}
}

// ============================================================================

// Comment for a story
type Comment struct {
	ID           int64          `json:"id"`
	StoryID      int64          `json:"story_id"`
	Body         string         `json:"body"`
	CreatedBlock int64          `json:"created_block"`
	Creator      sdk.AccAddress `json:"creator"`
}

// ============================================================================

// Evidence for a story
type Evidence struct {
	ID           int64          `json:"id"`
	StoryID      int64          `json:"story_id"`
	CreatedBlock int64          `json:"created_block"`
	Creator      sdk.AccAddress `json:"creator"`
	URI          string         `json:"uri"`
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

// CoinDenom is the coin denomination for the category (trubtc, trustablecoins, etc)
func (i StoryCategory) CoinDenom() string {
	return "tru" + i.Slug()
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
	BackIDs      []int64          `json:"back_ids,omitempty"`
	CommentIDs   []int64          `json:"comment_i_ds,omitempty"`
	EvidenceIDs  []int64          `json:"evidence_i_ds,omitempty"`
	Thread       []int64          `json:"thread,omitempty"`
	VoteIDs      []int64          `json:"vote_i_ds"`
	Body         string           `json:"body"`
	Category     StoryCategory    `json:"category"`
	CreatedBlock int64            `json:"created_block"`
	Creator      sdk.AccAddress   `json:"creator"`
	Escrow       sdk.AccAddress   `json:"escrow"`
	Round        int64            `json:"round"`
	State        StoryState       `json:"state"`
	StoryType    StoryType        `json:"type"`
	UpdatedBlock int64            `json:"updated_block"`
	Users        []sdk.AccAddress `json:"users"`
	VoteMaxNum   int64            `json:"vote_max_num"`
	VoteStart    time.Time        `json:"vote_start"`
	VoteEnd      time.Time        `json:"vote_end"`
}

// NewStory creates a new story
func NewStory(
	id int64,
	backIDs []int64,
	commentIDs []int64,
	evidenceIDs []int64,
	thread []int64,
	voteIDs []int64,
	body string,
	category StoryCategory,
	createdBlock int64,
	creator sdk.AccAddress,
	escrow sdk.AccAddress,
	round int64,
	state StoryState,
	storyType StoryType,
	updatedBlock int64,
	users []sdk.AccAddress,
	voteMax int64,
	voteStart time.Time,
	voteEnd time.Time) Story {

	return Story{
		ID:           id,
		BackIDs:      backIDs,
		CommentIDs:   commentIDs,
		EvidenceIDs:  evidenceIDs,
		Thread:       thread,
		VoteIDs:      voteIDs,
		Body:         body,
		Category:     category,
		CreatedBlock: createdBlock,
		Creator:      creator,
		Escrow:       escrow,
		Round:        round,
		State:        Created,
		StoryType:    storyType,
		UpdatedBlock: updatedBlock,
		Users:        users,
		VoteMaxNum:   voteMax,
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

// NewVote creates a new Vote instance
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
