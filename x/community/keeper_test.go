package community

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCommunity_Success(t *testing.T) {
	ctx, keeper := mockDB()

	name, slug, description := getFakeCommunityParams()

	community, err := keeper.NewCommunity(ctx, name, slug, description)
	assert.Nil(t, err)

	assert.NotZero(t, community.ID)
	assert.Equal(t, community.Name, name)
	assert.Equal(t, community.Slug, slug)
	assert.Equal(t, community.Description, description)
}

func TestNewCommunity_InvalidName(t *testing.T) {
	ctx, keeper := mockDB()

	_, slug, description := getFakeCommunityParams()
	invalidName := "Some really really really long name for a community"

	_, err := keeper.NewCommunity(ctx, invalidName, slug, description)
	assert.NotNil(t, err)
	assert.Equal(t, ErrInvalidCommunityMsg("").Code(), err.Code())
}

func TestNewCommunity_InvalidSlug(t *testing.T) {
	ctx, keeper := mockDB()

	name, _, description := getFakeCommunityParams()
	invalidSlug := "some-really-really-really-long-name-for-a-community"

	_, err := keeper.NewCommunity(ctx, name, invalidSlug, description)
	assert.NotNil(t, err)
	assert.Equal(t, ErrInvalidCommunityMsg("").Code(), err.Code())
}

func TestNewCommunity_InvalidDescription(t *testing.T) {
	ctx, keeper := mockDB()

	name, slug, _ := getFakeCommunityParams()
	invalidDescription := "If I could ever think of a really silly day of my life, I would choose the day when I tried fitting in more than 140 chars in a tweet. How silly it was!"

	_, err := keeper.NewCommunity(ctx, name, slug, invalidDescription)
	assert.NotNil(t, err)
	assert.Equal(t, ErrInvalidCommunityMsg("").Code(), err.Code())
}

func TestCommunity_Success(t *testing.T) {
	ctx, keeper := mockDB()

	name, slug, description := getFakeCommunityParams()

	createdCommunity, err := keeper.NewCommunity(ctx, name, slug, description)
	assert.Nil(t, err)

	returnedCommunity, err := keeper.Community(ctx, createdCommunity.ID)
	assert.Nil(t, err)
	assert.Equal(t, createdCommunity.ID, returnedCommunity.ID)
	assert.Equal(t, createdCommunity.Name, returnedCommunity.Name)
	assert.Equal(t, createdCommunity.Slug, returnedCommunity.Slug)
	assert.Equal(t, createdCommunity.Description, returnedCommunity.Description)
}

func TestCommunity_ErrCommunityNotFound(t *testing.T) {
	ctx, keeper := mockDB()
	id := uint64(314) // any random number, what better than a pie 🥧

	_, err := keeper.Community(ctx, id)
	assert.NotNil(t, err)
	assert.Equal(t, ErrCommunityNotFound(id).Code(), err.Code())
}

func TestCommunities_Success(t *testing.T) {
	ctx, keeper := mockDB()

	name, slug, description := getFakeCommunityParams()
	first, err := keeper.NewCommunity(ctx, name, slug, description)
	assert.Nil(t, err)

	name2, slug2, description2 := getAnotherFakeCommunityParams()
	another, err := keeper.NewCommunity(ctx, name2, slug2, description2)
	assert.Nil(t, err)

	all := keeper.Communities(ctx)

	assert.Len(t, all, 2)
	assert.Equal(t, all[0].ID, first.ID)
	assert.Equal(t, all[1].ID, another.ID)
}
