package staking

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/TruStory/truchain/x/bank"
)

// Defines staking module constants
const (
	RouterKey    = ModuleName
	QuerierRoute = ModuleName
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

var bankTransactionMappings = []bank.TransactionType{
	StakeBacking:   bank.TransactionBacking,
	StakeChallenge: bank.TransactionChallenge,
	StakeUpvote:    bank.TransactionUpvote,
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
	ID          uint64
	ArgumentID  uint64
	Type        StakeType
	Amount      sdk.Coin
	Creator     sdk.AccAddress
	CreatedTime time.Time
	EndTime     time.Time
	Expired     bool
}

type Argument struct {
	ID             uint64
	Creator        sdk.AccAddress
	ClaimID        uint64
	Summary        string
	Body           string
	StakeType      StakeType
	UpvotedCount   uint64
	UpvotedStake   sdk.Coin
	TotalStake     sdk.Coin
	UnhelpfulCount uint64
	IsUnhelpful    bool
	CreatedTime    time.Time
	UpdatedTime    time.Time
}
