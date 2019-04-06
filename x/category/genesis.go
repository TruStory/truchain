package category

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GenesisState for categories
type GenesisState struct {
	Categories []Category `json:"categories"`
}

// DefaultGenesisState for tests
func DefaultGenesisState() GenesisState {
	categories := []Category{
		{ID: 1, Title: "Cryptocurrency", Slug: "crypto", TotalCred: sdk.NewCoin("crypto", sdk.ZeroInt())},
		{ID: 2, Title: "Memes", Slug: "meme", TotalCred: sdk.NewCoin("meme", sdk.ZeroInt())},
	}

	return GenesisState{
		Categories: categories,
	}
}

// InitGenesis loads initial categories from the genesis file
func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) {
	for _, cat := range data.Categories {
		keeper.setCategory(ctx, cat)
	}
	keeper.SetLen(ctx, int64(len(data.Categories)))
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
