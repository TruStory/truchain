package argument

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
)

// store keys for argument params
var (
	KeyMinArgumentLength = []byte("minArgumentLength")
	KeyMaxArgumentLength = []byte("maxArgumentLength")
)

// Params holds parameters for voting
type Params struct {
	MinArgumentLength int `json:"min_argument_length"`
	MaxArgumentLength int `json:"max_argument_length"`
}

// DefaultParams is the default parameters for voting
func DefaultParams() Params {
	return Params{
		MinArgumentLength: 10,
		MaxArgumentLength: 1000,
	}
}

// KeyValuePairs implements params.ParamSet
func (p *Params) KeyValuePairs() params.KeyValuePairs {
	return params.KeyValuePairs{
		{Key: KeyMinArgumentLength, Value: &p.MinArgumentLength},
		{Key: KeyMaxArgumentLength, Value: &p.MaxArgumentLength},
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
	logger.Info(fmt.Sprintf("Loaded argument module params: %+v", params))
}
