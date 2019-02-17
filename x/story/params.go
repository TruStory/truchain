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
var KeyExpireDuration = []byte("expireDuration")

// Params holds parameters for a story
type Params struct {
	ExpireDuration time.Duration `json:"expire_duration"`
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
	}
}

func (p Params) String() string {
	return fmt.Sprintf("Params <%s>", p.ExpireDuration)
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
