package slashing

import (
	"net/url"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// StakeType enum
type StakeType int

// Stake is the palceholder here until Staking module is done
type Stake struct {
	ID          uint64
	ArgumentID  uint64
	Type        StakeType
	Amount      sdk.Coin
	Creator     sdk.AccAddress
	CreatedTime time.Time
	EndTime     time.Time
	Expired     bool
}

// Argument ..
type Argument struct {
	ID             uint64         `json:"id"`
	Creator        sdk.AccAddress `json:"creator"`
	ClaimID        uint64         `json:"claim_id"`
	Summary        string         `json:"summary"`
	Body           string         `json:"body"`
	StakeType      StakeType      `json:"stake_type"`
	UpvotedCount   uint64         `json:"upvoted_count"`
	UpvotedStake   sdk.Coin       `json:"upvoted_stake"`
	TotalStake     sdk.Coin       `json:"total_stake"`
	UnhelpfulCount uint64         `json:"unhelpful_count"`
	IsUnhelpful    bool           `json:"is_unhelpful"`
	CreatedTime    time.Time      `json:"created_time"`
	UpdatedTime    time.Time      `json:"updated_time"`
}

// StakingKeeper is the expected Staking keeper interface for this module
type StakingKeeper interface {
	Stake(ctx sdk.Context, id uint64) (Stake, sdk.Error)
	Argument(ctx sdk.Context, argumentID uint64) (Argument, bool)
}

// Claim is the placeholder here until claim module is done
type Claim struct {
	ID                  uint64
	CommunityID         string
	Body                string
	Creator             sdk.AccAddress
	Source              url.URL
	TotalParticipants   int64
	TotalBackingStake   sdk.Coin
	TotalChallengeStake sdk.Coin
	CreatedTime         time.Time
}

// ClaimKeeper is the expected Claim keeper interface for this module
type ClaimKeeper interface {
	Claim(ctx sdk.Context, id uint64) (Claim, bool)
}
