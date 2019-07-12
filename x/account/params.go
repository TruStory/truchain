package account

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
)

// Keys for params
var (
	KeyRegistrar     = []byte("registrar")
	KeyMaxSlashCount = []byte("maxSlashCount")
	KeyJailDuration  = []byte("jailTime")
)

// Params holds parameters for Auth
type Params struct {
	Registrar     sdk.AccAddress `json:"registrar"`
	MaxSlashCount int            `json:"max_slash_count"`
	JailDuration  time.Duration  `json:"jail_duration"`
}

// DefaultParams is the auth params for testing
func DefaultParams() Params {
	return Params{
		Registrar:     nil,
		MaxSlashCount: 3,
		JailDuration:  24 * time.Hour * 7,
	}
}

// ParamSetPairs implements params.ParamSet
func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		{Key: KeyRegistrar, Value: &p.Registrar},
		{Key: KeyMaxSlashCount, Value: &p.MaxSlashCount},
		{Key: KeyJailDuration, Value: &p.JailDuration},
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
