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

// BackingParams holds data for backing interest calculations
type BackingParams struct {
	AmountWeight    sdk.Dec
	PeriodWeight    sdk.Dec
	MinPeriod       time.Duration
	MaxPeriod       time.Duration
	MinInterestRate sdk.Dec
	MaxInterestRate sdk.Dec
}

// NewBackingParams creates a new BackingParams type with defaults
func NewBackingParams() BackingParams {
	return BackingParams{
		AmountWeight:    sdk.NewDecWithPrec(333, 3), // 33.3%
		PeriodWeight:    sdk.NewDecWithPrec(667, 3), // 66.7%
		MinPeriod:       3 * 24 * time.Hour,         // 3 days
		MaxPeriod:       90 * 24 * time.Hour,        // 90 days
		MinInterestRate: sdk.ZeroDec(),              // 0%
		MaxInterestRate: sdk.NewDecWithPrec(10, 2),  // 10%
	}
}

// Backing type
type Backing struct {
	ID        int64          `json:"id"`
	StoryID   int64          `json:"story_id"`
	Principal sdk.Coin       `json:"principal"`
	Interest  sdk.Coin       `json:"interest"`
	Expires   time.Time      `json:"expires"`
	Params    BackingParams  `json:"params"`
	Period    time.Duration  `json:"period"`
	User      sdk.AccAddress `json:"user"`
}

// NewBacking creates a new backing type
func NewBacking(
	id int64,
	storyID int64,
	principal sdk.Coin,
	interest sdk.Coin,
	expires time.Time,
	params BackingParams,
	period time.Duration,
	creator sdk.AccAddress) Backing {

	return Backing{
		ID:        id,
		StoryID:   storyID,
		Principal: principal,
		Interest:  interest,
		Expires:   expires,
		Params:    params,
		Period:    period,
		User:      creator,
	}
}

// ============================================================================

// Category is a type that defines the category for a story
type Category struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description,omitempty"`
}

// CoinName returns the name of the coin, alias for slug
func (c Category) CoinName() string {
	return c.Slug
}

// NewCategory creates a new story category type
func NewCategory(id int64, name string, slug string, description string) Category {
	return Category{
		ID:          id,
		Name:        name,
		Slug:        slug,
		Description: description,
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
	Body         string           `json:"body"`
	CategoryID   int64            `json:"category_id"`
	CreatedBlock int64            `json:"created_block"`
	Creator      sdk.AccAddress   `json:"creator"`
	Round        int64            `json:"round"`
	State        StoryState       `json:"state"`
	StoryType    StoryType        `json:"type"`
	UpdatedBlock int64            `json:"updated_block"`
	Users        []sdk.AccAddress `json:"users"`
}

// NewStory creates a new story
func NewStory(
	id int64,
	backIDs []int64,
	commentIDs []int64,
	evidenceIDs []int64,
	thread []int64,
	body string,
	categoryID int64,
	createdBlock int64,
	creator sdk.AccAddress,
	round int64,
	state StoryState,
	storyType StoryType,
	updatedBlock int64,
	users []sdk.AccAddress) Story {

	return Story{
		ID:           id,
		BackIDs:      backIDs,
		CommentIDs:   commentIDs,
		EvidenceIDs:  evidenceIDs,
		Thread:       thread,
		Body:         body,
		CategoryID:   categoryID,
		CreatedBlock: createdBlock,
		Creator:      creator,
		Round:        round,
		State:        Created,
		StoryType:    storyType,
		UpdatedBlock: updatedBlock,
		Users:        users,
	}
}
