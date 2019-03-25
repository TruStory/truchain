package argument

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	ctx, k, _, _ := mockDB()

	stakeID := int64(0)
	storyID := int64(0)
	creator := sdk.AccAddress([]byte{1, 2})

	argumentID, err := k.Create(ctx, stakeID, storyID, 0, "argument body", creator)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), argumentID)
}

func TestRegisterLike(t *testing.T) {
	ctx, k, _, _ := mockDB()

	stakeID := int64(0)
	storyID := int64(0)
	creator := sdk.AccAddress([]byte{1, 2})

	argumentID, err := k.Create(ctx, stakeID, storyID, 0, "argument body", creator)
	assert.NoError(t, err)

	err = k.RegisterLike(ctx, argumentID, creator)
	assert.NoError(t, err)
}

func TestArgument(t *testing.T) {
	ctx, k, _, _ := mockDB()

	stakeID := int64(0)
	storyID := int64(0)
	creator := sdk.AccAddress([]byte{1, 2})

	argumentID, err := k.Create(ctx, stakeID, storyID, 0, "argument body", creator)
	assert.NoError(t, err)

	argument, err := k.Argument(ctx, argumentID)
	assert.NoError(t, err)
	assert.Equal(t, argumentID, argument.ID)
}

func TestArgumentValidation(t *testing.T) {
	ctx, k, _, _ := mockDB()

	stakeID := int64(0)
	storyID := int64(0)
	creator := sdk.AccAddress([]byte{1, 2})

	_, err := k.Create(ctx, stakeID, storyID, int64(0), "", creator)
	assert.Error(t, err)

	_, err = k.Create(ctx, stakeID, storyID, int64(5), "argument body", creator)
	assert.Error(t, err)
}

func TestLikesByArgumentID(t *testing.T) {
	ctx, k, _, _ := mockDB()

	stakeID := int64(0)
	storyID := int64(0)
	creator1 := sdk.AccAddress([]byte{1, 2})
	creator2 := sdk.AccAddress([]byte{2, 4})

	argumentID, err := k.Create(ctx, stakeID, storyID, 0, "argument body", creator1)
	assert.NoError(t, err)

	err = k.RegisterLike(ctx, argumentID, creator1)
	assert.NoError(t, err)

	err = k.RegisterLike(ctx, argumentID, creator2)
	assert.NoError(t, err)

	likes, err := k.LikesByArgumentID(ctx, argumentID)
	assert.NoError(t, err)
	assert.Len(t, likes, 2)
}
