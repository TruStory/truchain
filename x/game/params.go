package game

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
)

// DefaultParamspace defines the default game module parameter subspace
const DefaultParamspace = "games"

// store keys for expiration params
var (
	KeyChallengeToBackingRatio = []byte("challengeToBackingRatio")
	KeyMinChallengeThreshold   = []byte("minChallengeThreshold")
	KeyMinQuorum               = []byte("minQuorum")
	KeyVotingPeriod            = []byte("votingPeriod")
)

// Params holds parameters for a game
type Params struct {
	ChallengeToBackingRatio sdk.Dec       `json:"challenge_to_backing_ratio"`
	MinChallengeThreshold   sdk.Int       `json:"min_challenge_threshold"`
	MinQuorum               int           `json:"min_quorum"`
	VotingPeriod            time.Duration `json:"voting_period"`
}

// DefaultParams is the story params for testing
func DefaultParams() Params {
	return Params{
		ChallengeToBackingRatio: sdk.NewDecWithPrec(100, 2), // 100%
		MinChallengeThreshold:   sdk.NewInt(1),              // 1 preethi
		MinQuorum:               3,
		VotingPeriod:            1 * 24 * time.Hour,
	}
}

// KeyValuePairs implements params.ParamSet
func (p *Params) KeyValuePairs() params.KeyValuePairs {
	return params.KeyValuePairs{
		{Key: KeyChallengeToBackingRatio, Value: &p.ChallengeToBackingRatio},
		{Key: KeyMinChallengeThreshold, Value: &p.MinChallengeThreshold},
		{Key: KeyMinQuorum, Value: &p.MinQuorum},
		{Key: KeyVotingPeriod, Value: &p.VotingPeriod},
	}
}

// ParamTypeTable for story module
func ParamTypeTable() params.TypeTable {
	return params.NewTypeTable().RegisterParamSet(&Params{})
}

func (k Keeper) minQuorum(ctx sdk.Context) (res int) {
	k.paramStore.Get(ctx, KeyMinQuorum, &res)
	return
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
	logger := ctx.Logger().With("module", "game")
	k.paramStore.SetParamSet(ctx, &params)
	logger.Info(fmt.Sprintf("Loaded game module params: %+v", params))
}
