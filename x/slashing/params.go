package slashing

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
)

// Keys for params
var (
	KeyMaxStakeSlashCount = []byte("maxStakeSlashCount")
	KeySlashMagnitude     = []byte("slashMagnitude")
	KeySlashMinStake      = []byte("slashMinStake")
	KeySlashAdmins        = []byte("slashAdmins")
	KeyJailTime           = []byte("jailTime")
)

// Params holds parameters for Slashing
type Params struct {
	MaxStakeSlashCount int              `json:"max_slash_stake_count"`
	SlashMagnitude     sdk.Dec          `json:"slash_magnitude"`
	SlashMinStake      sdk.Coin         `json:"slash_min_stake"`
	SlashAdmins        []sdk.AccAddress `json:"slash_admins"`
	JailTime           time.Duration    `json:"jail_time"`
}

// DefaultParams is the Slashing params for testing
func DefaultParams() Params {
	admin, err := sdk.AccAddressFromBech32("cosmos1xqc5gwzpgdr4wjz8xscnys2jx3f9x4zy223g9w")
	if err != nil {
		panic(err)
	}
	return Params{
		MaxStakeSlashCount: 50,
		SlashMagnitude:     sdk.NewDec(3),
		SlashMinStake:      sdk.NewCoin("trustake", sdk.NewInt(50)),
		JailTime:           time.Duration((7 * 24) * time.Hour),
		SlashAdmins:        []sdk.AccAddress{admin},
	}
}

// ParamSetPairs implements params.ParamSet
func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		{Key: KeyMaxStakeSlashCount, Value: &p.MaxStakeSlashCount},
		{Key: KeySlashMagnitude, Value: &p.SlashMagnitude},
		{Key: KeySlashMinStake, Value: &p.SlashMinStake},
		{Key: KeySlashAdmins, Value: &p.SlashAdmins},
		{Key: KeyJailTime, Value: &p.JailTime},
	}
}

// ParamKeyTable for slashing module
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&Params{})
}

// GetParams gets the genesis params for the slashing
func (k Keeper) GetParams(ctx sdk.Context) Params {
	var paramSet Params
	k.paramStore.GetParamSet(ctx, &paramSet)
	return paramSet
}

// SetParams sets the params for the slashing
func (k Keeper) SetParams(ctx sdk.Context, params Params) {
	logger := ctx.Logger().With("module", ModuleName)
	k.paramStore.SetParamSet(ctx, &params)
	logger.Info(fmt.Sprintf("Loaded slashing params: %+v", params))
}
