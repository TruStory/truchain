package category

import sdk "github.com/cosmos/cosmos-sdk/types"

// ExportGenesis ...
func ExportGenesis(ctx sdk.Context, ck WriteKeeper) []Category {

	categories, _ := ck.GetAllCategories(ctx)

	return categories
}
