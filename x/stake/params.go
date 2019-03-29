package stake

import (
	"fmt"

	app "github.com/TruStory/truchain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
)

// store keys for voting params
var (
	KeyMaxAmount        = []byte("maxAmount")
	KeyMinInterestRate  = []byte("minInterestRate")
	KeyMaxInterestRate  = []byte("maxInterestRate")
	KeyMajorityPercent  = []byte("majorityPercent")
	KeyAmountWeight     = []byte("amountWeight")
	KeyPeriodWeight     = []byte("periodWeight")
	KeyStakeToCredRatio = []byte("stakeToCredRatio")
)

// Params holds parameters for voting
type Params struct {
	MaxAmount        sdk.Coin `json:"max_amount"`
	MinInterestRate  sdk.Dec  `json:"min_interest_rate"`
	MaxInterestRate  sdk.Dec  `json:"max_interest_rate"`
	MajorityPercent  sdk.Dec  `json:"majority_percent"`
	AmountWeight     sdk.Dec  `json:"amount_weight"`
	PeriodWeight     sdk.Dec  `json:"period_weight"`
	StakeToCredRatio sdk.Int  `json:"stake_to_cred_ratio"`
}

// DefaultParams is the default parameters for voting
func DefaultParams() Params {
	return Params{
		MaxAmount:        sdk.NewCoin(app.StakeDenom, sdk.NewInt(100*app.Shanev)),
		AmountWeight:     sdk.NewDecWithPrec(333, 3), // 33.3%
		PeriodWeight:     sdk.NewDecWithPrec(667, 3), // 66.7%
		MinInterestRate:  sdk.ZeroDec(),              // 0%
		MaxInterestRate:  sdk.NewDecWithPrec(10, 2),  // 10%
		MajorityPercent:  sdk.NewDecWithPrec(51, 2),  // 51%
		StakeToCredRatio: sdk.NewInt(10),             // 10:1 ratio
	}
}

// KeyValuePairs implements params.ParamSet
func (p *Params) KeyValuePairs() params.KeyValuePairs {
	return params.KeyValuePairs{
		{Key: KeyMaxAmount, Value: &p.MaxAmount},
		{Key: KeyMinInterestRate, Value: &p.MinInterestRate},
		{Key: KeyMaxInterestRate, Value: &p.MaxInterestRate},
		{Key: KeyAmountWeight, Value: &p.AmountWeight},
		{Key: KeyPeriodWeight, Value: &p.PeriodWeight},
		{Key: KeyMajorityPercent, Value: &p.MajorityPercent},
		{Key: KeyStakeToCredRatio, Value: &p.StakeToCredRatio},
	}
}

// ParamTypeTable for story module
func ParamTypeTable() params.TypeTable {
	return params.NewTypeTable().RegisterParamSet(&Params{})
}

// GetParams gets the genesis params for the type
func (k Keeper) GetParams(ctx sdk.Context) Params {
	var paramSet Params
	k.paramStore.GetParamSet(ctx, &paramSet)
	return paramSet
}

// SetParams sets the params for the module
func (k Keeper) SetParams(ctx sdk.Context, params Params) {
	logger := ctx.Logger().With("module", StoreKey)
	k.paramStore.SetParamSet(ctx, &params)
	logger.Info(fmt.Sprintf("Loaded stake module params: %+v", params))
}

func (k Keeper) majorityPercent(ctx sdk.Context) sdk.Dec {
	return k.GetParams(ctx).MajorityPercent
}
