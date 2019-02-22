package voting

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
)

// store keys for voting params
var (
	KeyChallengerRewardPoolShare = []byte("challengerRewardPoolShare")
	KeyMajorityPercent           = []byte("majorityPercent")
)

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

// KeyValuePairs implements params.ParamSet
func (p *Params) KeyValuePairs() params.KeyValuePairs {
	return params.KeyValuePairs{
		{Key: KeyChallengerRewardPoolShare, Value: &p.ChallengerRewardPoolShare},
		{Key: KeyMajorityPercent, Value: &p.MajorityPercent},
	}
}

// ParamTypeTable for story module
func ParamTypeTable() params.TypeTable {
	return params.NewTypeTable().RegisterParamSet(&Params{})
}

// func (k Keeper) amountWeight(ctx sdk.Context) (res sdk.Dec) {
// 	k.paramStore.Get(ctx, KeyAmountWeight, &res)
// 	return
// }

// func (k Keeper) periodWeight(ctx sdk.Context) (res sdk.Dec) {
// 	k.paramStore.Get(ctx, KeyPeriodWeight, &res)
// 	return
// }

// SetParams sets the params for the expiration module
func (k Keeper) SetParams(ctx sdk.Context, params Params) {
	logger := ctx.Logger().With("module", "voting")
	k.paramStore.SetParamSet(ctx, &params)
	logger.Info(fmt.Sprintf("Loaded voting module params: %+v", params))
}
