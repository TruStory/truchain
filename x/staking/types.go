package staking

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/TruStory/truchain/x/bank"
)

// Defines staking module constants
const (
	StoreKey          = ModuleName
	RouterKey         = ModuleName
	QuerierRoute      = ModuleName
	DefaultParamspace = ModuleName
)

type StakeType byte

func (t StakeType) String() string {
	if int(t) >= len(StakeTypeName) {
		return "Unknown"
	}
	return StakeTypeName[t]
}

const (
	StakeBacking StakeType = iota
	StakeChallenge
	StakeUpvote
)

var StakeTypeName = []string{
	StakeBacking:   "StakeBacking",
	StakeChallenge: "StakeChallenge",
	StakeUpvote:    "StakeUpvote",
}

var bankTransactionMappings = []TransactionType{
	StakeBacking:   TransactionBacking,
	StakeChallenge: TransactionChallenge,
	StakeUpvote:    TransactionUpvote,
}

func (t StakeType) BankTransactionType() bank.TransactionType {
	if int(t) >= len(bankTransactionMappings) {
		panic("invalid stake type")
	}
	return bankTransactionMappings[t]
}

func (t StakeType) ValidForArgument() bool {
	return t.oneOf([]StakeType{StakeBacking, StakeChallenge})
}

func (t StakeType) ValidForUpvote() bool {
	return t.oneOf([]StakeType{StakeBacking, StakeChallenge})
}

func (t StakeType) Valid() bool {
	return t.oneOf([]StakeType{StakeBacking, StakeChallenge, StakeUpvote})
}

func (t StakeType) oneOf(types []StakeType) bool {
	for _, tType := range types {
		if tType == t {
			return true
		}
	}
	return false
}

type Stake struct {
	ID          uint64         `json:"id"`
	ArgumentID  uint64         `json:"argument_id"`
	Type        StakeType      `json:"type"`
	Amount      sdk.Coin       `json:"amount"`
	Creator     sdk.AccAddress `json:"creator"`
	CreatedTime time.Time      `json:"created_time"`
	EndTime     time.Time      `json:"end_time"`
	Expired     bool           `json:"expired"`
}

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
