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
	KeyClaimAdmins    = []byte("claimAdmins")
)

// Params holds parameters for a Claim
type Params struct {
	MinClaimLength int              `json:"min_claim_length"`
	MaxClaimLength int              `json:"max_claim_length"`
	ClaimAdmins    []sdk.AccAddress `json:"claim_admins"`
}

// DefaultParams is the Claim params for testing
func DefaultParams() Params {
	return Params{
		MinClaimLength: 25,
		MaxClaimLength: 140,
		ClaimAdmins:    []sdk.AccAddress{},
	}
}

// ParamSetPairs implements params.ParamSet
func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		{Key: KeyMinClaimLength, Value: &p.MinClaimLength},
		{Key: KeyMaxClaimLength, Value: &p.MaxClaimLength},
		{Key: KeyClaimAdmins, Value: &p.ClaimAdmins},
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

// UpdateParams updates the required params
func (k Keeper) UpdateParams(ctx sdk.Context, updatesMap map[string]interface{}) sdk.Error {
	current := k.GetParams(ctx)
	updated := k.getUpdatedParams(current, updatesMap)
	k.SetParams(ctx, updated)

	return nil
}

func (k Keeper) getUpdatedParams(current Params, updatesMap map[string]interface{}) Params {
	// TODO: to be implemented
	return current
}
