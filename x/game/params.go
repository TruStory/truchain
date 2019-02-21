package game

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
)

// store keys for expiration params
var (
	KeyMinQuorum    = []byte("minQuorum")
	KeyVotingPeriod = []byte("votingPeriod")
)

// Params holds parameters for a game
type Params struct {
	MinQuorum    int           `json:"min_quorum"`
	VotingPeriod time.Duration `json:"voting_period"`
}

// DefaultParams is the story params for testing
func DefaultParams() Params {
	return Params{
		MinQuorum:    3,
		VotingPeriod: 1 * 24 * time.Hour,
	}
}

// KeyValuePairs implements params.ParamSet
func (p *Params) KeyValuePairs() params.KeyValuePairs {
	return params.KeyValuePairs{
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

// SetParams sets the params for the expiration module
func (k Keeper) SetParams(ctx sdk.Context, params Params) {
	logger := ctx.Logger().With("module", "game")
	k.paramStore.SetParamSet(ctx, &params)
	logger.Info(fmt.Sprintf("Loaded game module params: %+v", params))
}
