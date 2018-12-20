package story

import (
	"net/url"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

func TestAddGetStory(t *testing.T) {
	ctx, sk, ck := mockDB()

	// test getting a non-existent story
	_, err := sk.Story(ctx, int64(5))
	assert.NotNil(t, err)

	storyID := createFakeStory(ctx, sk, ck)

	// test getting an existing story
	savedStory, err := sk.Story(ctx, storyID)
	assert.Nil(t, err)

	story := Story{
		ID:         storyID,
		Body:       "Body of story.",
		CategoryID: int64(1),
		Creator:    sdk.AccAddress([]byte{1, 2}),
		State:      Unconfirmed,
		Type:       Default,
	}

	assert.Equal(t, story.ID, savedStory.ID, "Story received from store does not match expected value")

	// test incrementing id by adding another story
	body := "Body of story 2."
	creator := sdk.AccAddress([]byte{3, 4})
	kind := Default
	source := url.URL{}
	evidence := []Evidence{}
	argument := []Argument{}

	storyID, _ = sk.Create(ctx, argument, body, int64(1), creator, evidence, source, kind)
	assert.Equal(t, int64(2), storyID, "Story ID did not increment properly")

	coinName, _ := sk.CategoryDenom(ctx, storyID)
	assert.Equal(t, "trudex", coinName)
}

func TestChallenge(t *testing.T) {
	ctx, sk, ck := mockDB()

	storyID := createFakeStory(ctx, sk, ck)
	story, _ := sk.Story(ctx, storyID)
	assert.Equal(t, Unconfirmed, story.State, "state should match")

	sk.StartGame(ctx, storyID)
	story, _ = sk.Story(ctx, storyID)
	assert.Equal(t, Challenged, story.State, "state should match")
}

func TestUpdateStory(t *testing.T) {
	ctx, sk, ck := mockDB()

	storyID := createFakeStory(ctx, sk, ck)
	story, _ := sk.Story(ctx, storyID)

	story.State = Challenged
	story.Body = "akjdsfhadskf"

	sk.UpdateStory(ctx, story)
	updatedStory, _ := sk.Story(ctx, storyID)

	assert.Equal(t, story.Body, updatedStory.Body, "should match")
	assert.Equal(t, story.State, updatedStory.State, "should match")
}

func TestGetStoriesWithCategory(t *testing.T) {
	ctx, sk, ck := mockDB()

	numStories := 100

	for i := 0; i < numStories; i++ {
		createFakeStory(ctx, sk, ck)
	}

	stories, _ := sk.StoriesByCategoryID(ctx, 1)
	assert.Equal(t, numStories, len(stories))
}

func TestGetChallengedStoriesWithCategory(t *testing.T) {
	ctx, sk, ck := mockDB()

	numStories := 5
	for i := 0; i < numStories; i++ {
		createFakeStory(ctx, sk, ck)
	}

	sk.StartGame(ctx, 2)
	sk.StartGame(ctx, 3)

	stories, _ := sk.ChallengedStoriesWithCategory(ctx, 1)
	assert.Equal(t, 2, len(stories))
}

func TestFeedWithCategory(t *testing.T) {
	ctx, sk, ck := mockDB()

	numStories := 5
	for i := 0; i < numStories; i++ {
		createFakeStory(ctx, sk, ck)
	}

	sk.StartGame(ctx, 2)
	sk.StartGame(ctx, 4)

	stories, _ := sk.FeedByCategoryID(ctx, 1)

	assert.Equal(t, 5, len(stories))
	assert.Equal(t, Challenged, stories[0].State)
	assert.Equal(t, Challenged, stories[1].State)
	assert.Equal(t, Unconfirmed, stories[2].State)
}

func TestFeedTrending(t *testing.T) {
	ctx, sk, ck := mockDB()

	numStories := 5
	for i := 0; i < numStories; i++ {
		createFakeStory(ctx, sk, ck)
	}

	stories := sk.Stories(ctx)

	assert.Equal(t, 5, len(stories))
	assert.Equal(t, int64(5), stories[0].ID)
	assert.Equal(t, int64(1), stories[4].ID)
}
