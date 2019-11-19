package slashing

import (
	"fmt"
	"strings"
	"testing"

	app "github.com/TruStory/truchain/types"
	"github.com/TruStory/truchain/x/staking"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/stretchr/testify/assert"
	abci "github.com/tendermint/tendermint/abci/types"
)

func TestQuerySlash_Success(t *testing.T) {
	ctx, keeper := mockDB()

	stakeID := uint64(1)
	creator := keeper.GetParams(ctx).SlashAdmins[0]
	createdSlash, _, err := keeper.CreateSlash(ctx, stakeID, SlashTypeUnhelpful, SlashReasonPlagiarism, "", creator)
	assert.NoError(t, err)

	params := keeper.codec.MustMarshalJSON(QuerySlashParams{
		ID: 1,
	})

	query := abci.RequestQuery{
		Path: strings.Join([]string{"custom", QuerierRoute, QuerySlash}, "/"),
		Data: params,
	}

	result, sdkErr := querySlash(ctx, query, keeper)
	assert.NoError(t, sdkErr)

	var returnedSlash Slash
	jsonErr := keeper.codec.UnmarshalJSON(result, &returnedSlash)
	assert.NoError(t, jsonErr)
	assert.Equal(t, returnedSlash, createdSlash)
}

func TestQuerySlash_ErrNotFound(t *testing.T) {
	ctx, keeper := mockDB()

	params := keeper.codec.MustMarshalJSON(QuerySlashParams{
		ID: 1,
	})
	query := abci.RequestQuery{
		Path: strings.Join([]string{"custom", QuerierRoute, QuerySlash}, "/"),
		Data: params,
	}

	_, sdkErr := querySlash(ctx, query, keeper)
	assert.NotNil(t, sdkErr)
	assert.Equal(t, ErrSlashNotFound(1).Code(), sdkErr.Code())
}

func TestQuerySlashes_Success(t *testing.T) {
	ctx, keeper := mockDB()
	_, _, addr1, _ := getFakeAppAccountParams()
	_, _, addr2, _ := getFakeAppAccountParams()
	earned := sdk.NewCoins(sdk.NewInt64Coin("general", 70*app.Shanev))
	usersEarnings := []staking.UserEarnedCoins{
		staking.UserEarnedCoins{Address: addr1, Coins: earned},
		staking.UserEarnedCoins{Address: addr2, Coins: earned},
	}
	genesis := staking.DefaultGenesisState()
	genesis.UsersEarnings = usersEarnings
	staking.InitGenesis(ctx, keeper.stakingKeeper, genesis)

	p := keeper.GetParams(ctx)
	p.MinSlashCount = 2
	keeper.SetParams(ctx, p)

	staker := keeper.GetParams(ctx).SlashAdmins[1]
	_, err := keeper.stakingKeeper.SubmitArgument(ctx, "arg1", "summary1", staker, 1, staking.StakeBacking)
	assert.NoError(t, err)

	stakeID := uint64(1)

	first, _, err := keeper.CreateSlash(ctx, stakeID, SlashTypeUnhelpful, SlashReasonPlagiarism, "", addr1)
	assert.Nil(t, err)

	another, _, err := keeper.CreateSlash(ctx, stakeID, SlashTypeUnhelpful, SlashReasonPlagiarism, "", addr2)
	assert.Nil(t, err)

	result, sdkErr := querySlashes(ctx, keeper)
	assert.Nil(t, sdkErr)

	var all []Slash
	jsonErr := keeper.codec.UnmarshalJSON(result, &all)
	assert.NoError(t, jsonErr)
	assert.Len(t, all, 2)
	assert.Equal(t, all[0], first)
	assert.Equal(t, all[1], another)
}

func TestQueryParams_Success(t *testing.T) {
	ctx, keeper := mockDB()

	onChainParams := keeper.GetParams(ctx)

	query := abci.RequestQuery{
		Path: fmt.Sprintf("/custom/%s/%s", ModuleName, QueryParams),
	}

	querier := NewQuerier(keeper)
	resBytes, err := querier(ctx, []string{QueryParams}, query)
	assert.Nil(t, err)

	var returnedParams Params
	sdkErr := ModuleCodec.UnmarshalJSON(resBytes, &returnedParams)
	assert.Nil(t, sdkErr)
	assert.Equal(t, returnedParams, onChainParams)
}
