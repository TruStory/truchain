package challenge

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// import (
// 	"testing"

// 	sdk "github.com/cosmos/cosmos-sdk/types"
// 	"github.com/stretchr/testify/assert"
// )

func TestMarshaling(t *testing.T) {
	ctx, k := mockDB()

	challenge := Challenge{
		ID:      k.GetNextID(ctx, k.challengeKey),
		StoryID: int64(5),
	}

	bz := k.marshal(challenge)
	assert.NotNil(t, bz)

	value := k.unmarshal(bz)
	assert.IsType(t, Challenge{}, value, "should be right type")
	assert.Equal(t, challenge.StoryID, value.StoryID, "should be equal")
}

// func TestAddGetStory(t *testing.T) {
// 	ctx, sk, ck := mockDB()

// 	// test getting a non-existant story
// 	_, err := sk.GetStory(ctx, int64(5))
// 	assert.NotNil(t, err)

// 	storyID := createFakeStory(ctx, sk, ck)

// 	// test getting an existing story
// 	savedStory, err := sk.GetStory(ctx, storyID)
// 	assert.Nil(t, err)

// 	story := Story{
// 		ID:           storyID,
// 		Body:         "Body of story.",
// 		CategoryID:   int64(1),
// 		CreatedBlock: int64(0),
// 		Creator:      sdk.AccAddress([]byte{1, 2}),
// 		State:        Created,
// 		Kind:         Default,
// 	}

// 	assert.Equal(t, story, savedStory, "Story received from store does not match expected value")

// 	// test incrementing id by adding another story
// 	body := "Body of story 2."
// 	creator := sdk.AccAddress([]byte{3, 4})
// 	kind := Default

// 	storyID, _ = sk.NewStory(ctx, body, int64(1), creator, kind)
// 	assert.Equal(t, int64(2), storyID, "Story ID did not increment properly")
// }

// func TestGetStoriesWithCategory(t *testing.T) {
// 	ctx, sk, ck := mockDB()

// 	numStories := 100

// 	for i := 0; i < numStories; i++ {
// 		createFakeStory(ctx, sk, ck)
// 	}

// 	stories, _ := sk.GetStoriesWithCategory(ctx, 1)
// 	assert.Equal(t, numStories, len(stories))
// }
