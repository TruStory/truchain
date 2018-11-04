package story

import "fmt"

// create subspace prefix "categories:id:[CategoryID]:stories:challenged:id:[StoryID]"
func challengedStoriesByCategoryKey(
	k Keeper, catID int64, storyID int64) []byte {

	key := fmt.Sprintf(
		"%s:id:%d:%s:challenged:id:%d",
		k.categoryKeeper.GetStoreKey().Name(),
		catID,
		k.GetStoreKey().Name(),
		storyID)

	return []byte(key)
}

// create subspace prefix "categories:id:[CategoryID]:stories:id:[StoryID]"
func storiesByCategoryKey(k Keeper, catID int64, storyID int64) []byte {
	key := fmt.Sprintf(
		"%s:id:%d:%s:id:%d",
		k.categoryKeeper.GetStoreKey().Name(),
		catID,
		k.GetStoreKey().Name(),
		storyID)

	return []byte(key)
}

// create subspace prefix "categories:id:[CategoryID]:stories:challenged:id:"
func challengedStoryIDsByCategorySubspaceKey(k Keeper, catID int64) []byte {
	prefix := fmt.Sprintf(
		"%s:id:%d:%s:challenged:id:",
		k.categoryKeeper.GetStoreKey().Name(),
		catID,
		k.GetStoreKey().Name())

	return []byte(prefix)
}

// create subspace prefix "categories:id:[CategoryID]:stories:id:"
func storyIDsByCategorySubspaceKey(k Keeper, catID int64) []byte {
	prefix := fmt.Sprintf(
		"%s:id:%d:%s:id:",
		k.categoryKeeper.GetStoreKey().Name(),
		catID,
		k.GetStoreKey().Name())

	return []byte(prefix)
}
