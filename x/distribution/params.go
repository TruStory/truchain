package distribution

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
)

// DefaultParamspace defines the default distribution module parameter subspace
const DefaultParamspace = "distribution"

// store keys for distribution params
var (
	KeyAmountWeight    = []byte("amountWeight")
	KeyPeriodWeight    = []byte("periodWeight")
	KeyMinInterestRate = []byte("minInterestRate")
	KeyMaxInterestRate = []byte("maxInterestRate")
)

// Params holds parameters for a story
type Params struct {
	AmountWeight    sdk.Dec `json:"amount_weight"`
	PeriodWeight    sdk.Dec `json:"period_weight"`
	MinInterestRate sdk.Dec `json:"min_interest_rate"`
	MaxInterestRate sdk.Dec `json:"max_interest_rate"`
}

// DefaultParams is the story params for testing
func DefaultParams() Params {
	return Params{
		AmountWeight:    sdk.NewDecWithPrec(333, 3), // 33.3%
		PeriodWeight:    sdk.NewDecWithPrec(667, 3), // 66.7%
		MinInterestRate: sdk.ZeroDec(),              // 0%
		MaxInterestRate: sdk.NewDecWithPrec(10, 2),  // 10%
	}
}

// KeyValuePairs implements params.ParamSet
func (p *Params) KeyValuePairs() params.KeyValuePairs {
	return params.KeyValuePairs{
		{Key: KeyAmountWeight, Value: &p.AmountWeight},
		{Key: KeyPeriodWeight, Value: &p.PeriodWeight},
		{Key: KeyMinInterestRate, Value: &p.MinInterestRate},
		{Key: KeyMaxInterestRate, Value: &p.MaxInterestRate},
	}
}

// ParamTypeTable for story module
func ParamTypeTable() params.TypeTable {
	return params.NewTypeTable().RegisterParamSet(&Params{})
}

// AmountWeight for the backing/challenge interest
func (k Keeper) AmountWeight(ctx sdk.Context) (res sdk.Dec) {
	k.paramStore.Get(ctx, KeyAmountWeight, &res)
	return
}

// SetParams sets the params for the distribution module
func (k Keeper) SetParams(ctx sdk.Context, params Params) {
	logger := ctx.Logger().With("module", "distribution")
	k.paramStore.SetParamSet(ctx, &params)
	logger.Info(fmt.Sprintf("Loaded distribution module params: %+v", params))
}
