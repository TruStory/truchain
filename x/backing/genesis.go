package backing

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GenesisState - all slashing state that must be provided at genesis
type GenesisState struct {
	Backings  []Backing `json:"backing"`
	Unexpired []Backing `json:"unexpired"`
}

// ExportGenesis ...
func ExportGenesis(ctx sdk.Context, bk WriteKeeper) (data GenesisState) {

	// Get array of Backing struct
	backings := bk.Backings(ctx)

	// Get list of unexpired backings
	unexpiredBackings := bk.BackList(ctx)
	// fmt.Printf("%+v\n", unexpiredBackings)

	// Get backings <-> story mappings

	return GenesisState{
		Backings:  backings,
		Unexpired: unexpiredBackings,
	}
}
