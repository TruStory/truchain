package claim

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
)

// Keys for params
var (
	KeyMinClaimLength = []byte("minClaimLength")
	KeyMaxClaimLength = []byte("maxClaimLength")
)

// Params holds parameters for a Claim
type Params struct {
	MinClaimLength int `json:"min_claim_length"`
	MaxClaimLength int `json:"max_claim_length"`
}

// DefaultParams is the Claim params for testing
func DefaultParams() Params {
	return Params{
		MinClaimLength: 25,
		MaxClaimLength: 350,
	}
}

// ParamSetPairs implements params.ParamSet
func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		{Key: KeyMinClaimLength, Value: &p.MinClaimLength},
		{Key: KeyMaxClaimLength, Value: &p.MaxClaimLength},
	}
}

// ParamKeyTable for claim module
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&Params{})
}

// GetParams gets the genesis params for the claim
func (k Keeper) GetParams(ctx sdk.Context) Params {
	var paramSet Params
	k.paramStore.GetParamSet(ctx, &paramSet)
	return paramSet
}

// SetParams sets the params for the claim
func (k Keeper) SetParams(ctx sdk.Context, params Params) {
	k.paramStore.SetParamSet(ctx, &params)
	logger(ctx).Info(fmt.Sprintf("Loaded claim params: %+v", params))
}
