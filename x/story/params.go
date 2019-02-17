package story

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
)

// DefaultParamspace defines the default auth module parameter subspace
const DefaultParamspace = "story"

// KeyExpireDuration is store's key for expire duration
var (
	KeyExpireDuration    = []byte("expireDuration")
	KeyMinStoryLength    = []byte("minStoryLength")
	keyMaxStoryLength    = []byte("maxStoryLength")
	keyMinArgumentLength = []byte("minArgumentLength")
	keyMaxArgumentLength = []byte("maxArugmentLength")
)

// Params holds parameters for a story
type Params struct {
	ExpireDuration    time.Duration `json:"expire_duration"`
	MinStoryLength    int           `json:"min_story_length"`
	MaxStoryLength    int           `json:"max_story_length"`
	MinArgumentLength int           `json:"min_argument_length"`
	MaxArgumentLength int           `json:"max_argument_length"`
}

// DefaultParams is the story params for testing
func DefaultParams() Params {
	return Params{
		ExpireDuration: 1 * 24 * time.Hour,
	}
}

// KeyValuePairs implements params.ParamSet
func (p *Params) KeyValuePairs() params.KeyValuePairs {
	return params.KeyValuePairs{
		{KeyExpireDuration, &p.ExpireDuration},
		{KeyMinStoryLength, &p.MinStoryLength},
		{keyMaxStoryLength, &p.MaxStoryLength},
		{keyMinArgumentLength, &p.MinArgumentLength},
		{keyMaxArgumentLength, &p.MaxArgumentLength},
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
	k.paramStore.SetParamSet(ctx, &params)
}
