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
	Expired
)

// IsValid returns true if the value is listed in the enum definition, false otherwise.
func (i State) IsValid() bool {
	switch i {
	case Unconfirmed, Challenged, Confirmed, Rejected, Expired:
		return true
	}
	return false
}

func (i State) String() string {
	return [...]string{"Unconfirmed", "Challenged", "Confirmed", "Rejected", "Expired"}[i]
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
		"Story <%d %s %s %d %s>", s.ID, s.Body, s.ExpireTime, s.GameID, s.State)
}
