package community

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCommunity_Success(t *testing.T) {
	ctx, keeper := mockDB()

	id, name, description := getFakeCommunityParams()

	community, err := keeper.NewCommunity(ctx, id, name, description)
	assert.Nil(t, err)

	assert.Equal(t, community.Name, name)
	assert.Equal(t, community.ID, id)
	assert.Equal(t, community.Description, description)
}

func TestNewCommunity_InvalidName(t *testing.T) {
	ctx, keeper := mockDB()

	id, _, description := getFakeCommunityParams()
	invalidName := "Some really really really long name for a community"

	_, err := keeper.NewCommunity(ctx, id, invalidName, description)
	assert.NotNil(t, err)
	assert.Equal(t, ErrInvalidCommunityMsg("").Code(), err.Code())
}

func TestNewCommunity_InvalidID(t *testing.T) {
	ctx, keeper := mockDB()

	_, name, description := getFakeCommunityParams()
	invalidID := "some-really-really-really-long-name-for-a-community"

	_, err := keeper.NewCommunity(ctx, invalidID, name, description)
	assert.NotNil(t, err)
	assert.Equal(t, ErrInvalidCommunityMsg("").Code(), err.Code())
}

func TestNewCommunity_InvalidDescription(t *testing.T) {
	ctx, keeper := mockDB()

	id, name, _ := getFakeCommunityParams()
	invalidDescription := "If I could ever think of a really silly day of my life, I would choose the day when I tried fitting in more than 140 chars in a tweet. How silly it was!"

	_, err := keeper.NewCommunity(ctx, id, name, invalidDescription)
	assert.NotNil(t, err)
	assert.Equal(t, ErrInvalidCommunityMsg("").Code(), err.Code())
}

func TestCommunity_Success(t *testing.T) {
	ctx, keeper := mockDB()

	id, name, description := getFakeCommunityParams()

	createdCommunity, err := keeper.NewCommunity(ctx, id, name, description)
	assert.Nil(t, err)

	returnedCommunity, err := keeper.Community(ctx, id)
	assert.Nil(t, err)
	assert.Equal(t, createdCommunity.Name, returnedCommunity.Name)
	assert.Equal(t, createdCommunity.ID, returnedCommunity.ID)
	assert.Equal(t, createdCommunity.Description, returnedCommunity.Description)
}

func TestCommunity_ErrCommunityNotFound(t *testing.T) {
	ctx, keeper := mockDB()
	id := "slug"

	_, err := keeper.Community(ctx, id)
	assert.NotNil(t, err)
	assert.Equal(t, ErrCommunityNotFound(id).Code(), err.Code())
}

func TestCommunities_Success(t *testing.T) {
	ctx, keeper := mockDB()

	id, name, description := getFakeCommunityParams()
	first, err := keeper.NewCommunity(ctx, id, name, description)
	assert.Nil(t, err)

	id2, name2, description2 := getAnotherFakeCommunityParams()
	another, err := keeper.NewCommunity(ctx, id2, name2, description2)
	assert.Nil(t, err)

	all := keeper.Communities(ctx)

	assert.Len(t, all, 4)
	assert.Equal(t, all[0].ID, "crypto")
	assert.Equal(t, all[1].ID, "meme")
	assert.Equal(t, all[2].ID, first.ID)
	assert.Equal(t, all[3].ID, another.ID)
}
