package params

import (
	"github.com/TruStory/truchain/x/account"
	"github.com/TruStory/truchain/x/bank"
	"github.com/TruStory/truchain/x/claim"
	"github.com/TruStory/truchain/x/community"
	"github.com/TruStory/truchain/x/slashing"
	"github.com/TruStory/truchain/x/staking"
)

// TODO [shanev]: these will be added by https://github.com/TruStory/truchain/issues/399
// MinBackingAmount : '1000000000',
// MaxBackingAmount:  '100000000000',

// Defines module constants
const (
	QuerierRoute = ModuleName
)

// Params defines defaults for a story
type Params struct {
	AccountParams   account.Params
	CommunityParams community.Params
	ClaimParams     claim.Params
	BankParams      bank.Params
	StakingParams   staking.Params
	SlashingParams  slashing.Params
}

// DefaultParams creates the default params
func DefaultParams() Params {
	return Params{
		AccountParams:   account.DefaultParams(),
		CommunityParams: community.DefaultParams(),
		ClaimParams:     claim.DefaultParams(),
		BankParams:      bank.DefaultParams(),
		StakingParams:   staking.DefaultParams(),
		SlashingParams:  slashing.DefaultParams(),
	}
}
