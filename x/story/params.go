package story

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
)

// DefaultParamspace defines the default auth module parameter subspace
const DefaultParamspace = "story"

// KeyExpireDuration is store's key for expire duration
var (
	KeyExpireDuration = []byte("expireDuration")
	KeyMinStoryLength = []byte("minStoryLength")
	KeyMaxStoryLength = []byte("maxStoryLength")
)

// Params holds parameters for a story
type Params struct {
	ExpireDuration time.Duration `json:"expire_duration"`
	MinStoryLength int           `json:"min_story_length"`
	MaxStoryLength int           `json:"max_story_length"`
}

// DefaultParams is the story params for testing
func DefaultParams() Params {
	return Params{
		ExpireDuration: 1 * 24 * time.Hour,
		MinStoryLength: 25,
		MaxStoryLength: 350,
	}
}

// KeyValuePairs implements params.ParamSet
func (p *Params) KeyValuePairs() params.KeyValuePairs {
	return params.KeyValuePairs{
		{Key: KeyExpireDuration, Value: &p.ExpireDuration},
		{Key: KeyMinStoryLength, Value: &p.MinStoryLength},
		{Key: KeyMaxStoryLength, Value: &p.MaxStoryLength},
	}
}

// ParamTypeTable for story module
func ParamTypeTable() params.TypeTable {
	return params.NewTypeTable().RegisterParamSet(&Params{})
}

// ExpireDuration for the story
func (k Keeper) ExpireDuration(ctx sdk.Context) (res time.Duration) {
	k.paramStore.Get(ctx, KeyExpireDuration, &res)
	return
}

// SetParams sets the params for the story
func (k Keeper) SetParams(ctx sdk.Context, params Params) {
	logger := ctx.Logger().With("module", "story")
	k.paramStore.SetParamSet(ctx, &params)
	logger.Info(fmt.Sprintf("Loaded story params: %+v", params))
}
