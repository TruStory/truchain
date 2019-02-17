package category

import sdk "github.com/cosmos/cosmos-sdk/types"

// InitGenesis loads initial categories from the genesis file
func InitGenesis(ctx sdk.Context, categoryKeeper WriteKeeper, categories []Category) {
	for _, cat := range categories {
		categoryKeeper.Create(ctx, cat.Title, cat.Slug, cat.Description)
	}
}
