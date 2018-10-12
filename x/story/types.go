package story

import sdk "github.com/cosmos/cosmos-sdk/types"

// Evidence for a story
type Evidence struct {
	ID           int64          `json:"id"`
	StoryID      int64          `json:"story_id"`
	CreatedBlock int64          `json:"created_block"`
	Creator      sdk.AccAddress `json:"creator"`
	URI          string         `json:"uri"`
}

// ============================================================================

// State is a type that defines a story state
type State int

// List of acceptable story states
const (
	Created State = iota
	Validated
	Rejected
	Unverifiable
	Challenged
	Revoked
)

// IsValid returns true if the value is listed in the enum defintion, false otherwise.
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
	ID           int64            `json:"id"`
	BackIDs      []int64          `json:"back_ids,omitempty"`
	EvidenceIDs  []int64          `json:"evidence_i_ds,omitempty"`
	Thread       []int64          `json:"thread,omitempty"`
	Body         string           `json:"body"`
	CategoryID   int64            `json:"category_id"`
	CreatedBlock int64            `json:"created_block"`
	Creator      sdk.AccAddress   `json:"creator"`
	Round        int64            `json:"round"`
	State        State            `json:"state"`
	Kind         Kind             `json:"kind"`
	UpdatedBlock int64            `json:"updated_block"`
	Users        []sdk.AccAddress `json:"users"`
}

// NewStory creates a new story
func NewStory(
	id int64,
	backIDs []int64,
	evidenceIDs []int64,
	thread []int64,
	body string,
	categoryID int64,
	createdBlock int64,
	creator sdk.AccAddress,
	round int64,
	state State,
	kind Kind,
	updatedBlock int64,
	users []sdk.AccAddress) Story {

	return Story{
		ID:           id,
		BackIDs:      backIDs,
		EvidenceIDs:  evidenceIDs,
		Thread:       thread,
		Body:         body,
		CategoryID:   categoryID,
		CreatedBlock: createdBlock,
		Creator:      creator,
		Round:        round,
		State:        Created,
		Kind:         kind,
		UpdatedBlock: updatedBlock,
		Users:        users,
	}
}

// List defines a list of story IDs
type List []int64

// IsEmpty checks if the story list is empty
func (sl List) IsEmpty() bool {
	if len(sl) == 0 {
		return true
	}
	return false
}
