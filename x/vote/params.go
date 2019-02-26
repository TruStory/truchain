package vote

import (
	"fmt"

	app "github.com/TruStory/truchain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
)

// store keys for vote params
var (
	KeyStakeAmount = []byte("stakeAmount")
)

// Params holds parameters for voting
type Params struct {
	StakeAmount sdk.Coin `json:"stake_amount"`
}

// DefaultParams is the default parameters for voting
func DefaultParams() Params {
	return Params{
		StakeAmount: sdk.NewCoin(app.StakeDenom, sdk.NewInt(10*app.Shanev)),
	}
}

// KeyValuePairs implements params.ParamSet
func (p *Params) KeyValuePairs() params.KeyValuePairs {
	return params.KeyValuePairs{
		{Key: KeyStakeAmount, Value: &p.StakeAmount},
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
	logger := ctx.Logger().With("module", "vote")
	k.paramStore.SetParamSet(ctx, &params)
	logger.Info(fmt.Sprintf("Loaded vote module params: %+v", params))
}
