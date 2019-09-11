package account

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// names used as root for pool module accounts:
const (
	UserGrowthPoolName  = "user_growth_tokens_pool"
	StakeholderPoolName = "stakeholder_tokens_pool"
)

// Pool - tracking user growth and stakeholder funds
type Pool struct {
	UserGrowthTokens  sdk.Int `json:"user_growth_tokens" yaml:"user_growth_tokens"`
	StakeholderTokens sdk.Int `json:"stakeholder_tokens" yaml:"stakeholder_tokens"`
}

// NewPool creates a new Pool instance used for queries
func NewPool(userGrowth, stakeholder sdk.Int) Pool {
	return Pool{
		UserGrowthTokens:  userGrowth,
		StakeholderTokens: stakeholder,
	}
}

// String returns a human readable string representation of a pool.
func (p Pool) String() string {
	return fmt.Sprintf(`Pool:	
  User Growth Tokens:  %s	
  Stakeholder Tokens:      %s`, p.UserGrowthTokens,
		p.StakeholderTokens)
}
