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
