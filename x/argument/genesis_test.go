package argument

import (
	"testing"
	"time"

	"github.com/TruStory/truchain/types"
	"github.com/davecgh/go-spew/spew"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

func TestExportGenesis(t *testing.T) {
	ctx, k, _, _ := mockDB()

	stakeID := int64(0)
	storyID := int64(0)
	creator := sdk.AccAddress([]byte{1, 2})
	liker1 := sdk.AccAddress([]byte{3, 4})
	liker2 := sdk.AccAddress([]byte{4, 5})

	argumentID, err := k.Create(ctx, stakeID, storyID, 0, "argument body", creator)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), argumentID)
	_, err = k.Create(ctx, stakeID, storyID, 0, "argument body", creator)
	assert.NoError(t, err)
	_, err = k.Create(ctx, stakeID, storyID, 0, "argument body", creator)
	assert.NoError(t, err)

	err = k.RegisterLike(ctx, argumentID, liker1)
	assert.NoError(t, err)
	err = k.RegisterLike(ctx, argumentID, liker2)
	assert.NoError(t, err)

	genesisState := ExportGenesis(ctx, k)
	spew.Dump(genesisState)
	assert.Equal(t, 3, len(genesisState.Arguments))
	assert.Equal(t, 1000, genesisState.Params.MaxArgumentLength)
	assert.Equal(t, 2, len(genesisState.Likes))
}

func TestImportGenesis(t *testing.T) {
	argument := Argument{
		ID:      1,
		StoryID: 0,
		StakeID: 0,
		Body:    "test argument",
		Creator: sdk.AccAddress([]byte{1, 2}),
		Timestamp: types.Timestamp{
			CreatedBlock: 0,
			CreatedTime:  time.Time{},
			UpdatedBlock: 0,
			UpdatedTime:  time.Time{},
		},
	}
	like := Like{
		ArgumentID: 1,
		Creator:    sdk.AccAddress([]byte{3, 4}),
		Timestamp: types.Timestamp{
			CreatedBlock: 0,
			CreatedTime:  time.Time{},
			UpdatedBlock: 0,
			UpdatedTime:  time.Time{},
		},
	}
	genesisState := GenesisState{
		Arguments: []Argument{argument},
		Likes:     []Like{like},
		Params:    DefaultParams(),
	}

	ctx, keeper, _, _ := mockDB()

	InitGenesis(ctx, keeper, genesisState)

	assert.Equal(t, 1, len(genesisState.Arguments))
	assert.Equal(t, 1000, genesisState.Params.MaxArgumentLength)
	assert.Equal(t, 1, len(genesisState.Likes))
}
