package challenge

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
	ctx, keeper, sk, _, bankKeeper := mockDB()

	storyID := createFakeStory(ctx, sk)
	amount := sdk.NewCoin(app.StakeDenom, sdk.NewInt(15000000000))
	argument := "test argument is long enough"
	creator := sdk.AccAddress([]byte{1, 2})
	bankKeeper.AddCoins(ctx, creator, sdk.Coins{amount})
	_, err := keeper.Create(ctx, storyID, amount, 0, argument, creator)
	assert.NoError(t, err)

	genesisState := ExportGenesis(ctx, keeper)
	assert.Equal(t, 1, len(genesisState.Challenges))
	assert.Equal(t, "1", genesisState.Params.MinChallengeStake.String())
}

func TestImportGenesis(t *testing.T) {
	challenge := Challenge{
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
		Challenges: []Challenge{challenge},
		Params:     DefaultParams(),
	}

	ctx, keeper, _, _, _ := mockDB()
	InitGenesis(ctx, keeper, genesisState)

	assert.Equal(t, 1, len(genesisState.Challenges))
	assert.Equal(t, "1", genesisState.Params.MinChallengeStake.String())
}
