package community

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

func TestNewCommunity_Success(t *testing.T) {
	ctx, keeper := mockDB()

	id, name, description := getFakeCommunityParams()
	creator := keeper.GetParams(ctx).CommunityAdmins[0]

	community, err := keeper.NewCommunity(ctx, id, name, description, creator)
	assert.Nil(t, err)

	assert.Equal(t, community.Name, name)
	assert.Equal(t, community.ID, id)
	assert.Equal(t, community.Description, description)
}

func TestNewCommunity_InvalidName(t *testing.T) {
	ctx, keeper := mockDB()

	id, _, description := getFakeCommunityParams()
	invalidName := "Some really really really long name for a community"
	creator := keeper.GetParams(ctx).CommunityAdmins[0]

	_, err := keeper.NewCommunity(ctx, id, invalidName, description, creator)
	assert.NotNil(t, err)
	assert.Equal(t, ErrInvalidCommunityMsg("").Code(), err.Code())
}

func TestNewCommunity_InvalidID(t *testing.T) {
	ctx, keeper := mockDB()

	_, name, description := getFakeCommunityParams()
	invalidID := "some-really-really-really-long-name-for-a-community"
	creator := keeper.GetParams(ctx).CommunityAdmins[0]

	_, err := keeper.NewCommunity(ctx, invalidID, name, description, creator)
	assert.NotNil(t, err)
	assert.Equal(t, ErrInvalidCommunityMsg("").Code(), err.Code())
}

func TestNewCommunity_InvalidDescription(t *testing.T) {
	ctx, keeper := mockDB()

	id, name, _ := getFakeCommunityParams()
	invalidDescription := "If I could ever think of a really silly day of my life, I would choose the day when I tried fitting in more than 140 chars in a tweet. How silly it was!"
	creator := keeper.GetParams(ctx).CommunityAdmins[0]

	_, err := keeper.NewCommunity(ctx, id, name, invalidDescription, creator)
	assert.NotNil(t, err)
	assert.Equal(t, ErrInvalidCommunityMsg("").Code(), err.Code())
}

func TestNewCommunity_CreatorNotAuthorised(t *testing.T) {
	ctx, keeper := mockDB()

	id, name, description := getFakeCommunityParams()
	unauthorisedCreator := sdk.AccAddress([]byte{1, 2})
	_, err := keeper.NewCommunity(ctx, id, name, description, unauthorisedCreator)
	assert.NotNil(t, err)
	assert.Equal(t, ErrAddressNotAuthorised().Code(), err.Code())
}

func TestCommunity_Success(t *testing.T) {
	ctx, keeper := mockDB()

	id, name, description := getFakeCommunityParams()
	creator := keeper.GetParams(ctx).CommunityAdmins[0]

	createdCommunity, err := keeper.NewCommunity(ctx, id, name, description, creator)
	assert.Nil(t, err)

	returnedCommunity, err := keeper.Community(ctx, id)
	assert.Nil(t, err)
	assert.Equal(t, createdCommunity.Name, returnedCommunity.Name)
	assert.Equal(t, createdCommunity.ID, returnedCommunity.ID)
	assert.Equal(t, createdCommunity.Description, returnedCommunity.Description)
}

func TestCommunity_ErrCommunityNotFound(t *testing.T) {
	ctx, keeper := mockDB()
	id := "id"

	_, err := keeper.Community(ctx, id)
	assert.NotNil(t, err)
	assert.Equal(t, ErrCommunityNotFound(id).Code(), err.Code())
}

func TestCommunities_Success(t *testing.T) {
	ctx, keeper := mockDB()

	id, name, description := getFakeCommunityParams()
	creator := keeper.GetParams(ctx).CommunityAdmins[0]
	first, err := keeper.NewCommunity(ctx, id, name, description, creator)
	assert.Nil(t, err)

	id2, name2, description2 := getAnotherFakeCommunityParams()
	another, err := keeper.NewCommunity(ctx, id2, name2, description2, creator)
	assert.Nil(t, err)

	all := keeper.Communities(ctx)

	assert.Len(t, all, 4)
	assert.Equal(t, all[0].ID, "crypto")
	assert.Equal(t, all[1].ID, "meme")
	assert.Equal(t, all[2].ID, first.ID)
	assert.Equal(t, all[3].ID, another.ID)
}

func TestUpdateParams_Success(t *testing.T) {
	ctx, keeper := mockDB()

	current := keeper.GetParams(ctx)
	updater := keeper.GetParams(ctx).CommunityAdmins[0]
	updates := Params{
		MinIDLength: current.MinIDLength + 20,
	}
	updatedFields := []string{"min_id_length"}
	keeper.UpdateParams(ctx, updater, updates, updatedFields)

	updated := keeper.GetParams(ctx)
	assert.Equal(t, current.MinIDLength+20, updated.MinIDLength)
}

func TestAddAdmin_Success(t *testing.T) {
	ctx, keeper := mockDB()

	creator := keeper.GetParams(ctx).CommunityAdmins[0]
	newAdmin := getFakeAdmin()

	err := keeper.AddAdmin(ctx, newAdmin, creator)
	assert.Nil(t, err)

	newAdmins := keeper.GetParams(ctx).CommunityAdmins
	assert.Subset(t, newAdmins, []sdk.AccAddress{newAdmin})
}

func TestAddAdmin_CreatorNotAuthorised(t *testing.T) {
	ctx, keeper := mockDB()

	invalidCreator := sdk.AccAddress([]byte{1, 2})
	newAdmin := getFakeAdmin()

	err := keeper.AddAdmin(ctx, newAdmin, invalidCreator)
	assert.NotNil(t, err)
	assert.Equal(t, ErrAddressNotAuthorised().Code(), err.Code())
}

func TestRemoveAdmin_Success(t *testing.T) {
	ctx, keeper := mockDB()

	currentAdmins := keeper.GetParams(ctx).CommunityAdmins
	adminToRemove := currentAdmins[0]

	err := keeper.RemoveAdmin(ctx, adminToRemove, adminToRemove) // removing self
	assert.Nil(t, err)
	newAdmins := keeper.GetParams(ctx).CommunityAdmins
	assert.Equal(t, len(currentAdmins)-1, len(newAdmins))
}

func TestRemoveAdmin_RemoverNotAuthorised(t *testing.T) {
	ctx, keeper := mockDB()

	invalidRemover := sdk.AccAddress([]byte{1, 2})
	currentAdmins := keeper.GetParams(ctx).CommunityAdmins
	adminToRemove := currentAdmins[0]

	err := keeper.AddAdmin(ctx, adminToRemove, invalidRemover)
	assert.NotNil(t, err)
	assert.Equal(t, ErrAddressNotAuthorised().Code(), err.Code())
}
