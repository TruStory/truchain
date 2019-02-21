package game

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
)

// DefaultParamspace defines the default game module parameter subspace
const DefaultParamspace = "game"

// store keys for expiration params
var (
	KeyQuorum = []byte("quorum")
)

// Params holds parameters for a story
type Params struct {
	Quorum int `json:"quorum"`
}

// DefaultParams is the story params for testing
func DefaultParams() Params {
	return Params{
		Quorum: 3,
	}
}

// KeyValuePairs implements params.ParamSet
func (p *Params) KeyValuePairs() params.KeyValuePairs {
	return params.KeyValuePairs{
		{Key: KeyQuorum, Value: &p.Quorum},
	}
}

// ParamTypeTable for story module
func ParamTypeTable() params.TypeTable {
	return params.NewTypeTable().RegisterParamSet(&Params{})
}

func (k Keeper) quorum(ctx sdk.Context) (res int) {
	k.paramStore.Get(ctx, KeyQuorum, &res)
	return
}

// SetParams sets the params for the expiration module
func (k Keeper) SetParams(ctx sdk.Context, params Params) {
	logger := ctx.Logger().With("module", "expiration")
	k.paramStore.SetParamSet(ctx, &params)
	logger.Info(fmt.Sprintf("Loaded expiration module params: %+v", params))
}
