package staking

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// names used as root for pool module accounts:
const (
	UserRewardPoolName = "user_reward_tokens_pool"
)

// Pool - tracking user reward funds
type Pool struct {
	UserRewardTokens sdk.Int `json:"user_reward_tokens" yaml:"user_reward_tokens"`
}

// NewPool creates a new Pool instance used for queries
func NewPool(userReward sdk.Int) Pool {
	return Pool{
		UserRewardTokens: userReward,
	}
}

// String returns a human readable string representation of a pool.
func (p Pool) String() string {
	return fmt.Sprintf(`Pool:	
  User Reward Tokens:  %s`, p.UserRewardTokens)
}
