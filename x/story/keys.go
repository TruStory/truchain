package story

import (
	"fmt"

	app "github.com/TruStory/truchain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// create subspace prefix "categories:id:[int64]:stories:time:[time.Time]"
func storyIDsByCategoryKey(
	k Keeper, catID int64, storyTimestamp app.Timestamp, challenged bool) []byte {

	return append(
		storyIDsByCategorySubspaceKey(k, catID, challenged),
		sdk.FormatTimeBytes(storyTimestamp.CreatedTime)...)
}

// create subspace prefix "categories:id:[int64]:stories:time:"
func storyIDsByCategorySubspaceKey(k Keeper, catID int64, challenged bool) []byte {
	format := "%s:id:%d:%s:time:"
	if challenged {
		format = "%s:id:%d:%s:challenged:time:"
	}

	key := fmt.Sprintf(
		format,
		k.categoryKeeper.GetStoreKey().Name(),
		catID,
		k.GetStoreKey().Name())

	return []byte(key)
}
