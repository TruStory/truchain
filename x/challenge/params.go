package challenge

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
)

// store keys for challenge params
var (
	KeyMinChallengeStake = []byte("minChallengeStake")
)

// Params holds parameters for a challenge
type Params struct {
	MinChallengeStake sdk.Int `json:"min_challenge_stake"`
}

// DefaultParams is the story params for testing
func DefaultParams() Params {
	return Params{
		MinChallengeStake: sdk.NewInt(1), // 1 preethi
	}
}

// ParamSetPairs implements params.ParamSet
func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		{Key: KeyMinChallengeStake, Value: &p.MinChallengeStake},
	}
}

// ParamKeyTable for story module
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&Params{})
}

// GetParams gets the genesis params for the type
func (k Keeper) GetParams(ctx sdk.Context) Params {
	var paramSet Params
	k.paramStore.GetParamSet(ctx, &paramSet)
	return paramSet
}

// SetParams sets the params for the expiration module
func (k Keeper) SetParams(ctx sdk.Context, params Params) {
	logger := ctx.Logger().With("module", "challenge")
	k.paramStore.SetParamSet(ctx, &params)
	logger.Info(fmt.Sprintf("Loaded challenge module params: %+v", params))
}

func (k Keeper) minChallengeStake(ctx sdk.Context) (res sdk.Int) {
	k.paramStore.Get(ctx, KeyMinChallengeStake, &res)
	return
}
