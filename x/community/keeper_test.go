package community

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCommunity_Success(t *testing.T) {
	ctx, keeper := mockDB()

	name, slug, description := getFakeCommunityParams()

	community := keeper.NewCommunity(ctx, name, slug, description)

	assert.NotZero(t, community.ID)
	assert.Equal(t, community.Name, name)
	assert.Equal(t, community.Slug, slug)
	assert.Equal(t, community.Description, description)
}

func TestCommunity_Success(t *testing.T) {
	ctx, keeper := mockDB()

	name, slug, description := getFakeCommunityParams()

	createdCommunity := keeper.NewCommunity(ctx, name, slug, description)

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
	first := keeper.NewCommunity(ctx, name, slug, description)

	name2, slug2, description2 := getAnotherFakeCommunityParams()
	another := keeper.NewCommunity(ctx, name2, slug2, description2)

	all := keeper.Communities(ctx)

	assert.Len(t, all, 2)
	assert.Equal(t, all[0].ID, first.ID)
	assert.Equal(t, all[1].ID, another.ID)
}
