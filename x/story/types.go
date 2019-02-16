package story

import (
	"fmt"
	"net/url"
	"time"

	app "github.com/TruStory/truchain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ============================================================================

// State is a type that defines a story state
type State int8

// List of acceptable story states
const (
	Unconfirmed = State(iota)
	Challenged
	Confirmed
	Rejected
)

// IsValid returns true if the value is listed in the enum definition, false otherwise.
func (i State) IsValid() bool {
	switch i {
	case Unconfirmed, Challenged, Confirmed, Rejected:
		return true
	}
	return false
}

func (i State) String() string {
	return [...]string{"Unconfirmed", "Challenged", "Confirmed", "Rejected"}[i]
}

// Type is a type that defines a story type
type Type int

// List of acceptable story types
const (
	Default Type = iota
	Identity
	Recovery
)

// IsValid returns true if a story type is valid, false otherwise.
func (i Type) IsValid() bool {
	switch i {
	case Default, Identity, Recovery:
		return true
	}
	return false
}

func (i Type) String() string {
	return [...]string{"Default", "Identity", "Recovery"}[i]
}

// Story type
type Story struct {
	ID         int64          `json:"id"`
	Argument   string         `json:"arguments,omitempty"`
	Body       string         `json:"body"`
	CategoryID int64          `json:"category_id"`
	Creator    sdk.AccAddress `json:"creator"`
	ExpireTime time.Time      `json:"expire_time"`
	Flagged    bool           `json:"flagged,omitempty"`
	GameID     int64          `json:"game_id,omitempty"`
	Source     url.URL        `json:"source,omitempty"`
	State      State          `json:"state"`
	Type       Type           `json:"type"`
	Timestamp  app.Timestamp  `json:"timestamp"`
}

func (s Story) String() string {
	return fmt.Sprintf(
		"Story <%d %s %s %d>", s.ID, s.Body, s.ExpireTime, s.GameID)
}

// Params holds parameters for a story
type Params struct {
	ExpireDuration time.Duration
}

// DefaultParams is the default parameters for voting
func DefaultParams() Params {
	return Params{
		ExpireDuration: 1 * 24 * time.Hour,
	}
}

// MsgParams holds default parameters for a story
type MsgParams struct {
	MinStoryLength    int // min number of chars for story body
	MaxStoryLength    int // max number of chars for story body
	MinArgumentLength int // min number of chars for argument
	MaxArgumentLength int // max number of chars for argument
}

// DefaultMsgParams creates a new MsgParams type with defaults
func DefaultMsgParams() MsgParams {
	return MsgParams{
		MinStoryLength:    25,
		MaxStoryLength:    350,
		MinArgumentLength: 10,
		MaxArgumentLength: 3000,
	}
}
