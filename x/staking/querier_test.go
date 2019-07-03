package staking

import (
	"strings"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	abci "github.com/tendermint/tendermint/abci/types"

	app "github.com/TruStory/truchain/types"
)

func TestQuerier_EmptyTopArgument(t *testing.T) {
	ctx, k, _ := mockDB()

	querier := NewQuerier(k)
	queryParams := QueryClaimTopArgumentParams{
		ClaimID: 1,
	}

	query := abci.RequestQuery{
		Path: strings.Join([]string{"custom", QuerierRoute, QueryClaimTopArgument}, "/"),
		Data: []byte{},
	}

	query.Data = k.codec.MustMarshalJSON(&queryParams)
	bz, err := querier(ctx, []string{QueryClaimTopArgument}, query)
	assert.NoError(t, err)
	argument := Argument{}
	jsonErr := k.codec.UnmarshalJSON(bz, &argument)
	assert.NoError(t, jsonErr)
	assert.Equal(t, uint64(0), argument.ID)

}

func TestQuerier_EarnedCoins(t *testing.T) {
	ctx, k, _ := mockDB()
	_, _, address := keyPubAddr()
	usersEarnings := make([]UserEarnedCoins, 0)
	coins := sdk.NewCoins(sdk.NewInt64Coin("crypto", app.Shanev*10),
		sdk.NewInt64Coin("random", app.Shanev*30))
	userEarnings := UserEarnedCoins{
		Address: address,
		Coins:   coins,
	}
	usersEarnings = append(usersEarnings, userEarnings)
	genesisState := NewGenesisState(nil, nil, usersEarnings, DefaultParams())
	InitGenesis(ctx, k, genesisState)

	querier := NewQuerier(k)
	queryParams := QueryEarnedCoinsParams{
		Address: address,
	}

	query := abci.RequestQuery{
		Path: strings.Join([]string{"custom", QuerierRoute, QueryEarnedCoins}, "/"),
		Data: []byte{},
	}

	query.Data = k.codec.MustMarshalJSON(&queryParams)
	bz, err := querier(ctx, []string{QueryEarnedCoins}, query)
	assert.NoError(t, err)
	earnedCoins := sdk.Coins{}
	jsonErr := k.codec.UnmarshalJSON(bz, &earnedCoins)
	assert.NoError(t, jsonErr)
	assert.True(t, coins.IsEqual(earnedCoins))
	// total

	queryTotalParams := QueryTotalEarnedCoinsParams{
		Address: address,
	}

	query = abci.RequestQuery{
		Path: strings.Join([]string{"custom", QuerierRoute, QueryTotalEarnedCoins}, "/"),
		Data: []byte{},
	}

	query.Data = k.codec.MustMarshalJSON(&queryTotalParams)
	bz, err = querier(ctx, []string{QueryTotalEarnedCoins}, query)
	assert.NoError(t, err)
	totalEarned := sdk.Coin{}
	jsonErr = k.codec.UnmarshalJSON(bz, &totalEarned)
	assert.NoError(t, jsonErr)
	assert.Equal(t, sdk.NewInt64Coin(app.StakeDenom, app.Shanev*40), totalEarned)

}
