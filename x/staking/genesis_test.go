package staking

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"

	app "github.com/TruStory/truchain/types"
)

func TestDefaultGenesisState(t *testing.T) {
	state := DefaultGenesisState()
	assert.Len(t, state.Arguments, 0)
	assert.Len(t, state.Stakes, 0)
	assert.Equal(t, state.Params, DefaultParams())
}

func TestInitGenesis(t *testing.T) {
	ctx, k, _ := mockDB()

	_, _, addr1 := keyPubAddr()
	_, _, addr2 := keyPubAddr()
	ctx = ctx.WithBlockTime(mustParseTime("2019-06-01"))
	stake1 := Stake{
		ID:          1,
		ArgumentID:  1,
		Type:        StakeBacking,
		Amount:      sdk.NewInt64Coin(app.StakeDenom, app.Shanev*50),
		Creator:     addr1,
		CreatedTime: ctx.BlockHeader().Time,
		EndTime:     ctx.BlockHeader().Time.Add(time.Hour * 24 * 7),
	}

	stake2 := Stake{
		ID:          2,
		ArgumentID:  1,
		Type:        StakeUpvote,
		Amount:      sdk.NewInt64Coin(app.StakeDenom, app.Shanev*10),
		Creator:     addr2,
		CreatedTime: ctx.BlockHeader().Time,
		EndTime:     ctx.BlockHeader().Time.Add(time.Hour * 24 * 7),
	}

	stakes := []Stake{stake1, stake2}

	argument1 := Argument{
		ID:           1,
		Creator:      addr1,
		ClaimID:      1,
		Summary:      "summary",
		Body:         "summary with *markdown* [trustory](http://trustory.io). and body, testing cuttoff on a [URL](http://somereally.long.url.even.longer.to.get.140.chars)",
		StakeType:    StakeBacking,
		CreatedTime:  ctx.BlockHeader().Time,
		UpdatedTime:  ctx.BlockHeader().Time,
		UpvotedCount: 1,
		UpvotedStake: sdk.NewInt64Coin(app.StakeDenom, app.Shanev*10),
		TotalStake:   sdk.NewInt64Coin(app.StakeDenom, app.Shanev*60),
	}

	expectedSummary := "summary with markdown trustory. and body, testing cuttoff on a URL"
	arguments := []Argument{argument1}
	arguments[0].Summary = expectedSummary

	usersEarnings := make([]UserEarnedCoins, 0)
	genesisState := NewGenesisState(arguments, stakes, usersEarnings, DefaultParams())
	InitGenesis(ctx, k, genesisState)
	actualGenesis := ExportGenesis(ctx, k)
	assert.Equal(t, genesisState, actualGenesis)

	// test association list are imported

	claimArguments := k.ClaimArguments(ctx, 1)
	assert.Equal(t, arguments, claimArguments)

	assert.Equal(t, expectedSummary, claimArguments[0].Summary)

	argumentStakes := k.ArgumentStakes(ctx, 1)
	assert.Equal(t, stakes, argumentStakes)

	assert.Equal(t, arguments, k.UserArguments(ctx, addr1))
	assert.Equal(t, []Argument{}, k.UserArguments(ctx, addr2))

	assert.Equal(t, []Stake{stake1}, k.UserStakes(ctx, addr1))
	assert.Equal(t, []Stake{stake2}, k.UserStakes(ctx, addr2))

	expiringStakes := make([]Stake, 0)

	k.IterateActiveStakeQueue(ctx, mustParseTime("2019-06-08"), func(stake Stake) bool {
		expiringStakes = append(expiringStakes, stake)
		return false
	})

	assert.Equal(t, stakes, expiringStakes)

}

func TestValidateGenesis(t *testing.T) {
	genesisState := NewGenesisState(nil, nil, nil, DefaultParams())
	genesisState.Params.ArgumentCreationStake.Denom = "my-denom"
	err := ValidateGenesis(genesisState)
	assert.Error(t, err)
	assert.Equal(t, ErrInvalidArgumentStakeDenom, err)
	genesisState.Params.ArgumentCreationStake.Denom = app.StakeDenom
	genesisState.Params.UpvoteStake.Denom = "my-denom"
	err = ValidateGenesis(genesisState)
	assert.Error(t, err)
	assert.Equal(t, ErrInvalidUpvoteStakeDenom, err)
}
