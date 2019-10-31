package distribution

import (
	"fmt"
	"reflect"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
)

// Keys for params
var (
	KeyUserGrowthAllocation  = []byte("userGrowthAllocation")
	KeyUserRewardAllocation  = []byte("userRewardAllocation")
	KeyStakeholderAllocation = []byte("stakeholderAllocation")
)

// Params holds parameters for Auth
type Params struct {
	UserGrowthAllocation  sdk.Dec `json:"user_growth_allocation"`
	UserRewardAllocation  sdk.Dec `json:"user_reward_allocation"`
	StakeholderAllocation sdk.Dec `json:"stakeholder_allocation"`
}

// DefaultParams is the auth params for testing
func DefaultParams() Params {
	return Params{
		UserGrowthAllocation:  sdk.NewDecWithPrec(25, 2),
		UserRewardAllocation:  sdk.NewDecWithPrec(25, 2),
		StakeholderAllocation: sdk.NewDecWithPrec(25, 2),
	}
}

// ParamSetPairs implements params.ParamSet
func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		{Key: KeyUserGrowthAllocation, Value: &p.UserGrowthAllocation},
		{Key: KeyUserRewardAllocation, Value: &p.UserRewardAllocation},
		{Key: KeyStakeholderAllocation, Value: &p.StakeholderAllocation},
	}
}

// ParamKeyTable for auth module
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&Params{})
}

// GetParams gets the genesis params for the auth
func (k Keeper) GetParams(ctx sdk.Context) Params {
	var paramSet Params
	k.paramStore.GetParamSet(ctx, &paramSet)
	return paramSet
}

// SetParams sets the params for the auth
func (k Keeper) SetParams(ctx sdk.Context, params Params) {
	logger := ctx.Logger().With("module", ModuleName)
	k.paramStore.SetParamSet(ctx, &params)
	logger.Info(fmt.Sprintf("Loaded account params: %+v", params))
}

// UpdateParams updates the required params
func (k Keeper) UpdateParams(ctx sdk.Context, updates Params, updatedFields []string) sdk.Error {
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
