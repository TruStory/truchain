package community

import (
	"fmt"
	"reflect"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
)

// Keys for params
var (
	KeyMinIDLength          = []byte("minIDLength")
	KeyMaxIDLength          = []byte("maxIDLength")
	KeyMinNameLength        = []byte("minNameLength")
	KeyMaxNameLength        = []byte("maxNameLength")
	KeyMaxDescriptionLength = []byte("maxDescriptionLength")
	KeyCommunityAdmins      = []byte("communityAdmins")
)

// Params holds parameters for a Community
type Params struct {
	MinIDLength          int              `json:"min_id_length"`
	MaxIDLength          int              `json:"max_id_length"`
	MinNameLength        int              `json:"min_name_length"`
	MaxNameLength        int              `json:"max_name_length"`
	MaxDescriptionLength int              `json:"max_description_length"`
	CommunityAdmins      []sdk.AccAddress `json:"community_admins"`
}

// DefaultParams is the Community params for testing
func DefaultParams() Params {
	return Params{
		MinNameLength:        5,
		MaxNameLength:        25,
		MinIDLength:          3,
		MaxIDLength:          15,
		MaxDescriptionLength: 140,
		CommunityAdmins:      []sdk.AccAddress{},
	}
}

// ParamSetPairs implements params.ParamSet
func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		{Key: KeyMinNameLength, Value: &p.MinNameLength},
		{Key: KeyMaxNameLength, Value: &p.MaxNameLength},
		{Key: KeyMinIDLength, Value: &p.MinIDLength},
		{Key: KeyMaxIDLength, Value: &p.MaxIDLength},
		{Key: KeyMaxDescriptionLength, Value: &p.MaxDescriptionLength},
		{Key: KeyCommunityAdmins, Value: &p.CommunityAdmins},
	}
}

// ParamKeyTable for community module
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&Params{})
}

// GetParams gets the genesis params for the community
func (k Keeper) GetParams(ctx sdk.Context) Params {
	var paramSet Params
	k.paramStore.GetParamSet(ctx, &paramSet)
	return paramSet
}

// SetParams sets the params for the community
func (k Keeper) SetParams(ctx sdk.Context, params Params) {
	logger := ctx.Logger().With("module", ModuleName)
	k.paramStore.SetParamSet(ctx, &params)
	logger.Info(fmt.Sprintf("Loaded community params: %+v", params))
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
