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
	New State = iota
	Voting
	Confirmed
	Rejected
	Expired
)

// IsValid returns true if the value is listed in the enum definition, false otherwise.
func (s State) IsValid() bool {
	switch s {
	case New, Voting, Confirmed, Rejected, Expired:
		return true
	}
	return false
}

func (s State) String() string {
	return [...]string{"New", "Voting", "Confirmed", "Rejected", "Expired"}[s]
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
	ID              int64          `json:"id"`
	Body            string         `json:"body"`
	CategoryID      int64          `json:"category_id"`
	Creator         sdk.AccAddress `json:"creator"`
	ExpireTime      time.Time      `json:"expire_time"`
	Flagged         bool           `json:"flagged,omitempty"`
	Source          url.URL        `json:"source,omitempty"`
	State           State          `json:"state"`
	Type            Type           `json:"type"`
	VotingStartTime time.Time      `json:"voting_start_time,omitempty"`
	VotingEndTime   time.Time      `json:"voting_end_time,omitempty"`
	Timestamp       app.Timestamp  `json:"timestamp"`
}

func (s Story) String() string {
	return fmt.Sprintf(
		"Story <%d %s %s %d %s>",
		s.ID, s.Body, s.ExpireTime, s.State, s.VotingEndTime)
}
