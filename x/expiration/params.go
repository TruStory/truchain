package expiration

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
)

// store keys for expiration params
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

// ParamSetPairs implements params.ParamSet
func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		{Key: KeyAmountWeight, Value: &p.AmountWeight},
		{Key: KeyPeriodWeight, Value: &p.PeriodWeight},
		{Key: KeyMinInterestRate, Value: &p.MinInterestRate},
		{Key: KeyMaxInterestRate, Value: &p.MaxInterestRate},
	}
}

// ParamKeyTable for story module
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&Params{})
}

// GetParams gets the genesis params for the type
func (k Keeper) GetParams(ctx sdk.Context) Params {
	var paramSet Params
	k.paramStore.GetParamSet(ctx, &paramSet)
	return paramSet
}

// SetParams sets the params for the expiration module
func (k Keeper) SetParams(ctx sdk.Context, params Params) {
	logger := ctx.Logger().With("module", StoreKey)
	k.paramStore.SetParamSet(ctx, &params)
	logger.Info(fmt.Sprintf("Loaded expiration module params: %+v", params))
}
