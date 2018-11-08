package story

import (
	app "github.com/TruStory/truchain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Evidence for a story
type Evidence struct {
	ID        int64          `json:"id"`
	StoryID   int64          `json:"story_id"`
	Creator   sdk.AccAddress `json:"creator"`
	URI       string         `json:"uri"`
	Timestamp app.Timestamp  `json:"timestamp"`
}

// ============================================================================

// State is a type that defines a story state
type State int8

// List of acceptable story states
const (
	Created = State(iota)
	Validated
	Rejected
	Unverifiable
	Challenged
	Revoked
)

// IsValid returns true if the value is listed in the enum definition, false otherwise.
func (i State) IsValid() bool {
	switch i {
	case Created, Validated, Rejected, Unverifiable, Challenged, Revoked:
		return true
	}
	return false
}

func (i State) String() string {
	return [...]string{"Created", "Validated", "Rejected", "Unverifiable", "Challenged", "Revoked"}[i]
}

// Kind is a type that defines a story type
type Kind int

// List of acceptable story types
const (
	Default Kind = iota
	Identity
	Recovery
)

// IsValid returns true if a story type is valid, false otherwise.
func (i Kind) IsValid() bool {
	switch i {
	case Default, Identity, Recovery:
		return true
	}
	return false
}

func (i Kind) String() string {
	return [...]string{"Default", "Identity", "Recovery"}[i]
}

// Story type
type Story struct {
	ID         int64          `json:"id"`
	Body       string         `json:"body"`
	CategoryID int64          `json:"category_id"`
	Creator    sdk.AccAddress `json:"creator"`
	GameID     int64          `json:"game_id"`
	State      State          `json:"state"`
	Kind       Kind           `json:"kind"`
	Timestamp  app.Timestamp  `json:"timestamp"`
}
