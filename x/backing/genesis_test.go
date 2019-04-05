package backing

import (
	"testing"
	"time"

	app "github.com/TruStory/truchain/types"
	"github.com/TruStory/truchain/x/stake"
	"github.com/cosmos/cosmos-sdk/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

func TestExportGenesis(t *testing.T) {
	ctx, keeper, sk, ck, bankKeeper, _ := mockDB()
	storyID := createFakeStory(ctx, sk, ck)
	amount := sdk.NewCoin(app.StakeDenom, sdk.NewInt(5000000))
	argument := "cool story brew.."
	creator := sdk.AccAddress([]byte{1, 2})
	bankKeeper.AddCoins(ctx, creator, sdk.Coins{amount})
	_, err := keeper.Create(ctx, storyID, amount, 0, argument, creator)
	assert.NoError(t, err)

	genesisState := ExportGenesis(ctx, keeper)
	assert.Equal(t, 1, len(genesisState.Backings))
}

func TestImportGenesis(t *testing.T) {
	backing := Backing{
		Vote: &stake.Vote{
			ID:      0,
			StoryID: 0,
			Amount: types.Coin{
				Denom:  "",
				Amount: types.Int{},
			},
			ArgumentID: 0,
			Creator:    nil,
			Vote:       false,
			Timestamp: app.Timestamp{
				CreatedBlock: 0,
				CreatedTime:  time.Time{},
				UpdatedBlock: 0,
				UpdatedTime:  time.Time{},
			},
		},
	}
	genesisState := GenesisState{
		Backings: []Backing{backing},
	}

	ctx, keeper, _, _, _, _ := mockDB()
	InitGenesis(ctx, keeper, genesisState)

	assert.Equal(t, 1, len(genesisState.Backings))
}
