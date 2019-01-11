package params

import (
	"time"

	"github.com/TruStory/truchain/x/backing"
	"github.com/TruStory/truchain/x/game"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Params defines defaults for a story
type Params struct {
	BackingMinDuration time.Duration `json:"backing_min_duration,omitempty"`
	BackingMaxDuration time.Duration `json:"backing_max_duration,omitempty"`
	ChallengeMinStake  sdk.Int       `json:"challenge_min_stake,omitempty"`
}

// DefaultParams creates the default params
func DefaultParams() Params {

	return Params{
		BackingMinDuration: backing.DefaultMsgParams().MinPeriod,
		BackingMaxDuration: backing.DefaultMsgParams().MaxPeriod,
		ChallengeMinStake:  game.DefaultParams().MinChallengeStake,
	}
}
