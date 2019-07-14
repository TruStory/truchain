package slashing

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Defines slashing module constants
const (
	StoreKey          = ModuleName
	RouterKey         = ModuleName
	QuerierRoute      = ModuleName
	DefaultParamspace = ModuleName
)

// Slash stores data about a slashing
type Slash struct {
	ID             uint64
	StakeID        uint64
	Type           SlashType
	Reason         SlashReason
	DetailedReason string
	Creator        sdk.AccAddress
	CreatedTime    time.Time
}

// Slashes is an array of slashes
type Slashes []Slash

func (s Slash) String() string {
	return fmt.Sprintf(`Slash %d:
  StakeID: %d
  Creator: %s
  CreatedTime: %s`,
		s.ID, s.StakeID, s.Creator.String(), s.CreatedTime.String())
}

// SlashType enum
type SlashType int

const (
	// SlashTypeUnhelpful represents the unhelpful slashing type
	SlashTypeUnhelpful SlashType = iota // 0
)

// SlashReason enum
type SlashReason int

const (
	// SlashReasonLogicOrEvidenceAbsent represents the reason when no clear logic or evidence is present
	SlashReasonLogicOrEvidenceAbsent SlashReason = iota
	// SlashReasonIssueNotAddressed represents the reason when the issue at hand is not addressed
	SlashReasonIssueNotAddressed
	// SlashReasonFocusedOnPerson represents the reason when the argument is focused on the person, not the idea
	SlashReasonFocusedOnPerson
	// SlashReasonPlagiarism represents the reason when the argument is plagiarised
	SlashReasonPlagiarism
	// SlashReasonOthers represents the reason that is any other than the above
	SlashReasonOthers
)
