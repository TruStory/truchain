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

	AttributeKeyMinSlashCountKey = "min-slash-count"
	AttributeKeySlashResults     = "slash-results"
)

// Slash stores data about a slashing
type Slash struct {
	ID             uint64
	ArgumentID     uint64
	Type           SlashType
	Reason         SlashReason
	DetailedReason string
	Creator        sdk.AccAddress
	CreatedTime    time.Time
}

type PunishmentResultType int

const (
	PunishmentInterestSlashed PunishmentResultType = iota
	PunishmentStakeSlashed
	PunishmentCuratorRewarded
	PunishmentJailed
)

type PunishmentResult struct {
	Type          PunishmentResultType `json:"type"`
	AppAccAddress sdk.AccAddress       `json:"address"`
	Coin          sdk.Coin             `json:"coin"`
}

// Slashes is an array of slashes
type Slashes []Slash

func (s Slash) String() string {
	return fmt.Sprintf(`Slash %d:
  ArgumentID: %d
  Creator: %s
  Reason: %d
  CreatedTime: %s`,
		s.ID, s.ArgumentID, s.Creator.String(), s.Reason, s.CreatedTime.String())
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
	// SlashNoOriginalThought represents the reasons when a thought isn't original;
	SlashNoOriginalThought
	// SlashReasonPlagiarism represents the reason when the argument is plagiarised
	SlashReasonPlagiarism
	// SlashReasonOther represents the reason that is any other than the above
	SlashReasonOther
	// SlashReasonHarassment ...
	SlashReasonHarassment
	// SlashReasonSpam ...
	SlashReasonSpam
	// SlashReasonOffensiveContent ...
	SlashReasonOffensiveContent
)

func (r SlashReason) String() string {
	if int(r) >= len(SlashReasonName) {
		return "Unknown"
	}
	return SlashReasonName[r]
}

// SlashReasonName is the reason for the slash
var SlashReasonName = []string{
	SlashReasonLogicOrEvidenceAbsent: "No clear logic or evidence",
	SlashReasonIssueNotAddressed:     "Doesn't address the issue",
	SlashReasonFocusedOnPerson:       "Focuses on the person",
	SlashNoOriginalThought:           "No original thought",
	SlashReasonPlagiarism:            "Plagiarism",
	SlashReasonOther:                 "Other",
	SlashReasonHarassment:            "Harassment",
	SlashReasonSpam:                  "Spam",
	SlashReasonOffensiveContent:      "Offensive Content",
}
