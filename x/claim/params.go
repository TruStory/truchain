package claim

import (
	"fmt"
	"reflect"

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
func (k Keeper) UpdateParams(ctx sdk.Context, updater sdk.AccAddress, updates Params, updatedFields []string) sdk.Error {
	if !k.isAdmin(ctx, updater) {
		err := ErrAddressNotAuthorised()
		return err
	}

	current := k.GetParams(ctx)
	updated := k.getUpdatedParams(current, updates, updatedFields)
	k.SetParams(ctx, updated)

	return nil
}

func (k Keeper) getUpdatedParams(current Params, updates Params, updatedFields []string) Params {
	updated := current
	mapParams(updates, func(param string, index int, field reflect.StructField) {
		if isIn(param, updatedFields) {
			reflect.ValueOf(&updated).Elem().FieldByName(field.Name).Set(
				reflect.ValueOf(
					reflect.ValueOf(updates).FieldByName(field.Name).Interface(),
				),
			)
		}
	})

	return updated
}

func isIn(needle string, haystack []string) bool {
	for _, value := range haystack {
		if needle == value {
			return true
		}
	}

	return false
}

// mapParams walks over each param, and ignores the *_admins param because they are out of scope for this CLI command
func mapParams(params interface{}, fn func(param string, index int, field reflect.StructField)) {
	rParams := reflect.TypeOf(params)
	for i := 0; i < rParams.NumField(); i++ {
		field := rParams.Field(i)
		param := field.Tag.Get("json")
		fn(param, i, field)
	}
}
