package category

import sdk "github.com/cosmos/cosmos-sdk/types"

// DefaultCategories for tests and chain init
func DefaultCategories() []Category {
	return []Category{
		{Title: "Cryptocurrency", Slug: "crypto"},
		{Title: "Memes", Slug: "meme"},
	}
}

// InitGenesis loads initial categories from the genesis file
func InitGenesis(ctx sdk.Context, categoryKeeper WriteKeeper, categories []Category) {
	for _, cat := range categories {
		categoryKeeper.Create(ctx, cat.Title, cat.Slug, cat.Description)
	}
}
