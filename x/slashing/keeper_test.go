package slashing

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

func TestNewSlash_Success(t *testing.T) {
	ctx, keeper := mockDB()

	stakeID := uint64(1)
	creator := DefaultParams().SlashAdmins[0]
	slash, err := keeper.NewSlash(ctx, stakeID, creator)
	assert.Nil(t, err)

	assert.NotZero(t, slash.ID)
	assert.Equal(t, slash.StakeID, stakeID)
	assert.Equal(t, slash.Creator, creator)
}

func TestNewSlash_InvalidStake(t *testing.T) {
	ctx, keeper := mockDB()

	invalidStakeID := uint64(404)
	creator := DefaultParams().SlashAdmins[0]
	_, err := keeper.NewSlash(ctx, invalidStakeID, creator)

	assert.NotNil(t, err)
	assert.Equal(t, ErrInvalidStake(invalidStakeID).Code(), err.Code())
}

func TestNewSlash_InvalidCreator(t *testing.T) {
	ctx, keeper := mockDB()

	stakeID := uint64(1)
	invalidCreator := sdk.AccAddress([]byte{1, 2})
	_, err := keeper.NewSlash(ctx, stakeID, invalidCreator)

	assert.NotNil(t, err)
	assert.Equal(t, ErrInvalidCreator(invalidCreator).Code(), err.Code())
}

func TestSlash_Success(t *testing.T) {
	ctx, keeper := mockDB()

	stakeID := uint64(1)
	creator := DefaultParams().SlashAdmins[0]
	createdSlash, err := keeper.NewSlash(ctx, stakeID, creator)
	assert.Nil(t, err)

	returnedSlash, err := keeper.Slash(ctx, createdSlash.ID)

	assert.Nil(t, err)
	assert.Equal(t, createdSlash, returnedSlash)
}

func TestSlash_ErrNotFound(t *testing.T) {
	ctx, keeper := mockDB()

	stakeID := uint64(1)
	creator := DefaultParams().SlashAdmins[0]
	_, err := keeper.NewSlash(ctx, stakeID, creator)
	assert.Nil(t, err)

	_, err = keeper.Slash(ctx, uint64(404))

	assert.NotNil(t, err)
	assert.Equal(t, ErrSlashNotFound(uint64(404)).Code(), err.Code())
}

func TestSlashes_Success(t *testing.T) {
	ctx, keeper := mockDB()

	stakeID := uint64(1)
	creator := DefaultParams().SlashAdmins[0]
	first, err := keeper.NewSlash(ctx, stakeID, creator)
	assert.Nil(t, err)

	stakeID2 := uint64(2)
	another, err := keeper.NewSlash(ctx, stakeID2, creator)
	assert.Nil(t, err)

	all := keeper.Slashes(ctx)
	assert.Len(t, all, 2)
	assert.Equal(t, all[0], first)
	assert.Equal(t, all[1], another)
}
