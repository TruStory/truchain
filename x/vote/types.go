package vote

import sdk "github.com/cosmos/cosmos-sdk/types"

// NOTE: the main type for vote lives in the top-level `types` package.
// This is because it is embedded in multiple other types like
// `Backing` and `Challenge` in other modules. So it should exist
// at a higher level than this module.

// MsgParams holds default parameters for a vote
type MsgParams struct {
	MinCommentLength int // min number of chars for argument
	MaxCommentLength int // max number of chars for argument
	MinEvidenceCount int // min number of evidence URLs
	MaxEvidenceCount int // max number of evidence URLs
}

// DefaultMsgParams creates a new MsgParams type with defaults
func DefaultMsgParams() MsgParams {
	return MsgParams{
		MinCommentLength: 10,
		MaxCommentLength: 340,
		MinEvidenceCount: 0,
		MaxEvidenceCount: 10,
	}
}

// Params holds parameters for voting
type Params struct {
	ChallengerRewardPoolShare sdk.Dec
	SupermajorityPercent      sdk.Dec
}

// DefaultParams is the default parameters for voting
func DefaultParams() Params {
	return Params{
		ChallengerRewardPoolShare: sdk.NewDecWithPrec(75, 2), // 75%
		SupermajorityPercent:      sdk.NewDecWithPrec(51, 2), // 51%
	}
}
