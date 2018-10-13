package story

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

func TestAddGetStory(t *testing.T) {
	ctx, sk, ck := mockDB()

	// test getting a non-existant story
	_, err := sk.GetStory(ctx, int64(5))
	assert.NotNil(t, err)

	storyID := createFakeStory(ctx, sk, ck)

	// test getting an existing story
	savedStory, err := sk.GetStory(ctx, storyID)
	assert.Nil(t, err)

	story := Story{
		ID:           storyID,
		Body:         "Body of story.",
		CategoryID:   int64(1),
		CreatedBlock: int64(0),
		Creator:      sdk.AccAddress([]byte{1, 2}),
		State:        Created,
		Kind:         Default,
	}

	assert.Equal(t, story, savedStory, "Story received from store does not match expected value")

	// test incrementing id by adding another story
	body := "Body of story 2."
	creator := sdk.AccAddress([]byte{3, 4})
	kind := Default

	storyID, _ = sk.NewStory(ctx, body, int64(1), creator, kind)
	assert.Equal(t, int64(2), storyID, "Story ID did not increment properly")
}

func TestGetStoriesWithCategory(t *testing.T) {
	ctx, sk, ck := mockDB()

	numStories := 100

	for i := 0; i < numStories; i++ {
		createFakeStory(ctx, sk, ck)
	}

	stories, _ := sk.GetStoriesWithCategory(ctx, 1)
	assert.Equal(t, numStories, len(stories))
}

// all = [1,2,3,4,5]
// challenged = [2,4]

// id = 1
// cid = 2
// unchallenged = [1]

// id = 2
// cid = 2
// break

// id = 3
// cid = 2
// cid = 4
// unchallenged = [1,3]

// id = 4
// cid = 2
// cid = 4
// break

// id = 5
// cid = 2
// cid = 4
// unchallenged = [1,3,5]

// unchallenged = [1,3,5]
// feed = [2,4,1,3,5]
