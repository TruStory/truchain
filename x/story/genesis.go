package story

import "time"
import sdk "github.com/cosmos/cosmos-sdk/types"

// GenesisState - all story state that must be provided at genesis
type GenesisState struct {
	ExpireDuration time.Duration `json:"expire_duration"`
}

// InitGenesis initializes story state from genesis file
func InitGenesis(ctx sdk.Context, storyKeeper WriteKeeper) {

}
