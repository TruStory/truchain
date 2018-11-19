package vote

import (
	app "github.com/TruStory/truchain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// TokenVote defines a simple token vote on a story
type TokenVote struct {
	app.Vote
}

// Amount implements `Voter`
func (v TokenVote) Amount() sdk.Coin {
	return v.Vote.Amount
}

// Creator implements `Voter`
func (v TokenVote) Creator() sdk.AccAddress {
	return v.Vote.Creator
}

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
	MajorityPercent           sdk.Dec
}

// DefaultParams is the default parameters for voting
func DefaultParams() Params {
	return Params{
		ChallengerRewardPoolShare: sdk.NewDecWithPrec(75, 2), // 75%
		MajorityPercent:           sdk.NewDecWithPrec(51, 2), // 51%
	}
}
