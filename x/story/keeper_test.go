package story

import (
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
)

func TestValidKeys(t *testing.T) {
	ctx, sk, ck := mockDB()

	storyID := createFakeStory(ctx, sk, ck)
	story, _ := sk.GetStory(ctx, storyID)

	key := getChallengedStoriesKey(sk, story.CategoryID)
	assert.Equal(t, "challenges:categories:id:1:stories", fmt.Sprintf("%s", key), "should be equal")
}

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

func TestChallenge(t *testing.T) {
	ctx, sk, ck := mockDB()

	storyID := createFakeStory(ctx, sk, ck)
	story, _ := sk.GetStory(ctx, storyID)
	assert.Equal(t, Created, story.State, "state should match")

	sk.StartChallenge(ctx, storyID)
	story, _ = sk.GetStory(ctx, storyID)
	spew.Dump(story)
	assert.Equal(t, Challenged, story.State, "state should match")
}

func TestUpdateStory(t *testing.T) {
	ctx, sk, ck := mockDB()

	storyID := createFakeStory(ctx, sk, ck)
	story, _ := sk.GetStory(ctx, storyID)

	story.State = Challenged
	story.Body = "akjdsfhadskf"

	sk.UpdateStory(ctx, story)
	updatedStory, _ := sk.GetStory(ctx, storyID)

	assert.Equal(t, story.Body, updatedStory.Body, "should match")
	assert.Equal(t, story.State, updatedStory.State, "should match")
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

func TestGetChallengedStoriesWithCategory(t *testing.T) {
	ctx, sk, ck := mockDB()

	numStories := 5
	for i := 0; i < numStories; i++ {
		createFakeStory(ctx, sk, ck)
	}

	sk.StartChallenge(ctx, 2)
	sk.StartChallenge(ctx, 3)

	stories, _ := sk.GetChallengedStoriesWithCategory(ctx, 1)
	assert.Equal(t, 2, len(stories))
}

func TestFeedWithCategory(t *testing.T) {
	ctx, sk, ck := mockDB()

	numStories := 5
	for i := 0; i < numStories; i++ {
		createFakeStory(ctx, sk, ck)
	}

	sk.StartChallenge(ctx, 2)
	sk.StartChallenge(ctx, 4)

	stories, _ := sk.GetFeedWithCategory(ctx, 1)

	assert.Equal(t, 5, len(stories))
	assert.Equal(t, Challenged, stories[0].State)
	assert.Equal(t, Challenged, stories[1].State)
	assert.Equal(t, Created, stories[2].State)
}
