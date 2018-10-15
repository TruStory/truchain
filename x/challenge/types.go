package challenge

import (
	"net/url"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Reason is a type that defines a reason for a challenge
type Reason int8

// List of acceptable challenge reasons
const (
	False = Reason(iota)
	NotEasilyFalsifiable
	Spam
	TooObvious
	TooSubjective
)

// IsValid return true if a reason is valid
func (i Reason) IsValid() bool {
	switch i {
	case False, NotEasilyFalsifiable, Spam, TooObvious, TooSubjective:
		return true
	}
	return false
}

// Challenge defines a challenge on a story
type Challenge struct {
	ID              int64            `json:"id"`
	StoryID         int64            `json:"story_id"`
	Arugment        string           `json:"arugment"`
	Challengers     []sdk.AccAddress `json:"challengers"`
	Creator         sdk.AccAddress   `json:"creator"`
	Evidence        []url.URL        `json:"evidence,omitempty"`
	Pool            sdk.Coin         `json:"pool"`
	Reason          Reason           `json:"reason"`
	ThresholdAmount sdk.Int          `json:"threshold_amount"`
	CreatedBlock    int64            `json:"created_block"`
	CreatedTime     time.Time        `json:"created_time"`
	UpdatedBlock    int64            `json:"updated_block"`
	UpdatedTime     time.Time        `json:"updated_time"`
}

// NewChallenge creates a new `Challenge` type with defaults
func NewChallenge(
	id int64,
	storyID int64,
	amount sdk.Coin,
	argument string,
	challengers []sdk.AccAddress,
	creator sdk.AccAddress,
	evidence []url.URL,
	pool sdk.Coin,
	reason Reason,
	thresholdAmount sdk.Int,
	createdBlock int64,
	createdTime time.Time) Challenge {

	return Challenge{
		ID:              id,
		StoryID:         storyID,
		Arugment:        argument,
		Challengers:     challengers,
		Creator:         creator,
		Evidence:        evidence,
		Pool:            pool,
		Reason:          reason,
		ThresholdAmount: thresholdAmount,
		CreatedBlock:    createdBlock,
		CreatedTime:     time.Now(),
	}
}

// Params holds default parameters for a challenge
type Params struct {
	MinArgumentLength int           // min number of chars for argument
	MaxArgumentLength int           // max number of chars for argument
	MinEvidenceCount  int           // min number of evidence URLs
	MaxEvidenceCount  int           // max number of evidence URLs
	MinChallengeStake sdk.Int       // min amount required to challenge
	Expires           time.Duration // time to expire if threshold not met
	Period            time.Duration // length of challenge game / voting period
	Threshold         int64         // amount at which game begins
}

// NewParams creates a new Params type with defaults
func NewParams() Params {
	return Params{
		MinArgumentLength: 10,
		MaxArgumentLength: 340,
		MinEvidenceCount:  0,
		MaxEvidenceCount:  10,
		MinChallengeStake: sdk.NewInt(10),
		Expires:           10 * 24 * time.Hour,
		Period:            30 * 24 * time.Hour,
		Threshold:         10,
	}
}
