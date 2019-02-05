package backing

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GenesisState - all slashing state that must be provided at genesis
type GenesisState struct {
	Backings     []Backing `json:"backings"`
	BackingsList []int64   `json:"backing_list"`
}

// ExportGenesis ...
func ExportGenesis(ctx sdk.Context, bk WriteKeeper) (data GenesisState) {

	// Get array of Backing struct
	backings, _ := bk.GetAllBackings(ctx)

	// Get list of unexpired backings
	backingsList := bk.GetAllBackingList(ctx)

	return GenesisState{
		Backings:     backings,
		BackingsList: backingsList,
	}
}
