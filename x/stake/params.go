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
	KeyMajorityPercent  = []byte("majorityPercent")
	KeyStakeToCredRatio = []byte("stakeToCredRatio")
	KeyInterestRate     = []byte("interestRate")
)

// Params holds parameters for voting
type Params struct {
	MaxAmount        sdk.Coin `json:"max_amount"`
	MajorityPercent  sdk.Dec  `json:"majority_percent"`
	StakeToCredRatio sdk.Int  `json:"stake_to_cred_ratio"`
	InterestRate     sdk.Dec  `json:"interest_rate"`
}

// DefaultParams is the default parameters for voting
func DefaultParams() Params {
	return Params{
		MaxAmount:        sdk.NewCoin(app.StakeDenom, sdk.NewInt(100*app.Shanev)),
		MajorityPercent:  sdk.NewDecWithPrec(51, 2), // 51%
		StakeToCredRatio: sdk.NewInt(10),            // 10:1 ratio
		InterestRate:     sdk.NewDecWithPrec(25, 2), // 25% for Year 1
	}
}

// KeyValuePairs implements params.ParamSet
func (p *Params) KeyValuePairs() params.KeyValuePairs {
	return params.KeyValuePairs{
		{Key: KeyMaxAmount, Value: &p.MaxAmount},
		{Key: KeyInterestRate, Value: &p.InterestRate},
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
