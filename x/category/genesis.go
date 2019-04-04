package category

import sdk "github.com/cosmos/cosmos-sdk/types"

// GenesisState for categories
type GenesisState struct {
	Categories []Category `json:"categories"`
}

// DefaultGenesisState for tests
func DefaultGenesisState() GenesisState {
	return GenesisState{
		Categories: DefaultCategories(),
	}
}

// DefaultCategories for tests and chain init
func DefaultCategories() []Category {
	return []Category{
		{Title: "Cryptocurrency", Slug: "crypto"},
		{Title: "Memes", Slug: "meme"},
	}
}

// InitGenesis loads initial categories from the genesis file
func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) {
	for _, cat := range data.Categories {
		keeper.setCategory(ctx, cat)
	}
}

// ExportGenesis exports the genesis state
func ExportGenesis(ctx sdk.Context, keeper Keeper) GenesisState {
	categories, err := keeper.GetAllCategories(ctx)
	if err != nil {
		// it is okay to panic here because the chain is not running
		panic(err)
	}

	return GenesisState{
		Categories: categories,
	}
}
