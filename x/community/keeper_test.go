package community

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateCommunity_Success(t *testing.T) {
	ctx, keeper := mockDB()

	name, slug, description := getFakeCommunityParams()

	community := keeper.Create(ctx, name, slug, description)

	assert.NotZero(t, community.ID)
	assert.Equal(t, community.Name, name)
	assert.Equal(t, community.Slug, slug)
	assert.Equal(t, community.Description, description)
}

func TestGetCommunity_Success(t *testing.T) {
	ctx, keeper := mockDB()

	name, slug, description := getFakeCommunityParams()

	createdCommunity := keeper.Create(ctx, name, slug, description)

	returnedCommunity, err := keeper.Get(ctx, createdCommunity.ID)

	assert.Nil(t, err)
	assert.Equal(t, createdCommunity.ID, returnedCommunity.ID)
	assert.Equal(t, createdCommunity.Name, returnedCommunity.Name)
	assert.Equal(t, createdCommunity.Slug, returnedCommunity.Slug)
	assert.Equal(t, createdCommunity.Description, returnedCommunity.Description)
}

func TestGetCommunity_ErrCategoryNotFound(t *testing.T) {
	ctx, keeper := mockDB()
	id := int64(314) // any random number, what better than a pie 🥧

	_, err := keeper.Get(ctx, id)
	assert.NotNil(t, err)
	assert.Equal(t, ErrCommunityNotFound(id).Code(), err.Code(), "should get error")
}
