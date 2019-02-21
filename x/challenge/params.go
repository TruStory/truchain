package challenge

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
)

// store keys for challenge params
var (
	KeyChallengeToBackingRatio = []byte("challengeToBackingRatio")
	KeyMinChallengeThreshold   = []byte("minChallengeThreshold")
)

// Params holds parameters for a challenge
type Params struct {
	ChallengeToBackingRatio sdk.Dec `json:"challenge_to_backing_ratio"`
	MinChallengeThreshold   sdk.Int `json:"min_challenge_threshold"`
}

// DefaultParams is the story params for testing
func DefaultParams() Params {
	return Params{
		ChallengeToBackingRatio: sdk.NewDecWithPrec(100, 2), // 100%
		MinChallengeThreshold:   sdk.NewInt(1),              // 1 preethi
	}
}

// KeyValuePairs implements params.ParamSet
func (p *Params) KeyValuePairs() params.KeyValuePairs {
	return params.KeyValuePairs{
		{Key: KeyChallengeToBackingRatio, Value: &p.ChallengeToBackingRatio},
		{Key: KeyMinChallengeThreshold, Value: &p.MinChallengeThreshold},
	}
}

// ParamTypeTable for story module
func ParamTypeTable() params.TypeTable {
	return params.NewTypeTable().RegisterParamSet(&Params{})
}

func (k Keeper) challengeToBackingRatio(ctx sdk.Context) (res sdk.Dec) {
	k.paramStore.Get(ctx, KeyChallengeToBackingRatio, &res)
	return
}

func (k Keeper) minChallengeThreshold(ctx sdk.Context) (res sdk.Int) {
	k.paramStore.Get(ctx, KeyMinChallengeThreshold, &res)
	return
}

// SetParams sets the params for the expiration module
func (k Keeper) SetParams(ctx sdk.Context, params Params) {
	logger := ctx.Logger().With("module", "challenge")
	k.paramStore.SetParamSet(ctx, &params)
	logger.Info(fmt.Sprintf("Loaded challenge module params: %+v", params))
}
