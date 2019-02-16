package category

import sdk "github.com/cosmos/cosmos-sdk/types"

// InitCategories creates the initial set of categories
// func (k Keeper) InitCategories(
// 	ctx sdk.Context, creator sdk.AccAddress, categories map[string]string) (err sdk.Error) {

// 	// sort keys to make category creation deterministic
// 	var keys []string
// 	for key := range categories {
// 		keys = append(keys, key)
// 	}
// 	sort.Strings(keys)

// 	for _, key := range keys {
// 		title := categories[key]
// 		slug := key
// 		_, err = k.NewCategory(ctx, title, creator, slug, "")
// 		if err != nil {
// 			return err
// 		}
// 	}

// 	return
// }

// InitGenesis loads initial categories from the genesis file
func InitGenesis(ctx sdk.Context, categoryKeeper WriteKeeper, categories []Category) {
	for _, cat := range categories {
		categoryKeeper.NewCategory(ctx, cat.Title, cat.Slug, cat.Description)
	}
}
