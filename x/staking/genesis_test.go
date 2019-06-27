package staking

import (
	"testing"
	"time"

	app "github.com/TruStory/truchain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

func TestDefaultGenesisState(t *testing.T) {
	state := DefaultGenesisState()
	assert.Len(t, state.Arguments, 0)
	assert.Len(t, state.Stakes, 0)
	assert.Equal(t, state.Params, DefaultParams())
}

func TestInitGenesis(t *testing.T) {
	ctx, k, _, _ := mockDB()

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
		Body:         "body",
		StakeType:    StakeBacking,
		CreatedTime:  ctx.BlockHeader().Time,
		UpdatedTime:  ctx.BlockHeader().Time,
		UpvotedCount: 1,
		UpvotedStake: sdk.NewInt64Coin(app.StakeDenom, app.Shanev*10),
		TotalStake:   sdk.NewInt64Coin(app.StakeDenom, app.Shanev*60),
	}
	arguments := []Argument{argument1}
	genesisState := NewGenesisState(arguments, stakes, DefaultParams())
	InitGenesis(ctx, k, genesisState)
	actualGenesis := ExportGenesis(ctx, k)
	assert.Equal(t, genesisState, actualGenesis)

	// test association list are imported

	claimArguments := k.ClaimArguments(ctx, 1)
	assert.Equal(t, arguments, claimArguments)

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
