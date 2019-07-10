package staking

import (
	"fmt"

	app "github.com/TruStory/truchain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stripmd "github.com/writeas/go-strip-markdown"
)

type UserEarnedCoins struct {
	Address sdk.AccAddress `json:"address"`
	Coins   sdk.Coins      `json:"coins"`
}

// GenesisState defines genesis data for the module
type GenesisState struct {
	Arguments     []Argument        `json:"arguments"`
	Params        Params            `json:"params"`
	Stakes        []Stake           `json:"stakes"`
	UsersEarnings []UserEarnedCoins `json:"users_earnings"`
}

// NewGenesisState creates a new genesis state.
func NewGenesisState(arguments []Argument, stakes []Stake, userEarnings []UserEarnedCoins, params Params) GenesisState {
	return GenesisState{
		Arguments:     arguments,
		Params:        params,
		Stakes:        stakes,
		UsersEarnings: userEarnings,
	}
}

// DefaultGenesisState returns a default genesis state
func DefaultGenesisState() GenesisState {
	return GenesisState{
		Params:        DefaultParams(),
		Stakes:        make([]Stake, 0),
		Arguments:     make([]Argument, 0),
		UsersEarnings: make([]UserEarnedCoins, 0),
	}
}

// InitGenesis initializes staking state from genesis file
func InitGenesis(ctx sdk.Context, k Keeper, data GenesisState) {
	for _, a := range data.Arguments {
		strippedBody := stripmd.Strip(a.Body)
		strippedBodyLen := len(strippedBody)
		stripLen := 140
		if strippedBodyLen < 140 {
			stripLen = strippedBodyLen
		}
		a.Summary = strippedBody[:stripLen]

		k.setArgument(ctx, a)
		k.setClaimArgument(ctx, a.ClaimID, a.ID)
		k.serUserArgument(ctx, a.Creator, a.ID)
	}
	for _, s := range data.Stakes {
		k.setStake(ctx, s)
		if !s.Expired {
			k.InsertActiveStakeQueue(ctx, s.ID, s.EndTime)
		}
		k.setArgumentStake(ctx, s.ArgumentID, s.ID)
		k.setUserStake(ctx, s.Creator, s.CreatedTime, s.ID)

		arg, ok := k.getArgument(ctx, s.ArgumentID)
		if !ok {
			panic(fmt.Sprintf("failed getting argument %d", s.ArgumentID))
		}
		// NOTE: this InitGenesis must run *after* claim InitGenesis
		claim, ok := k.claimKeeper.Claim(ctx, arg.ClaimID)
		if !ok {
			panic(fmt.Sprintf("failed getting claim %d", arg.ClaimID))
		}
		k.setCommunityStake(ctx, claim.CommunityID, s.ID)
		k.setUserCommunityStake(ctx, s.Creator, claim.CommunityID, s.ID)

	}
	k.setArgumentID(ctx, uint64(len(data.Arguments)+1))
	k.setStakeID(ctx, uint64(len(data.Stakes)+1))

	for _, e := range data.UsersEarnings {
		e.Coins.Sort()
		if !e.Coins.IsValid() {
			panic(fmt.Sprintf("user earnings for account %s are invalid %s", e.Address.String(), e.Coins.String()))
		}
		k.setEarnedCoins(ctx, e.Address, e.Coins.Sort())
	}
	k.SetParams(ctx, data.Params)
}

// ExportGenesis exports the genesis state
func ExportGenesis(ctx sdk.Context, keeper Keeper) GenesisState {
	return GenesisState{
		Params:        keeper.GetParams(ctx),
		Arguments:     keeper.Arguments(ctx),
		Stakes:        keeper.Stakes(ctx),
		UsersEarnings: keeper.UsersEarnings(ctx),
	}
}

// ValidateGenesis validates the genesis state data
func ValidateGenesis(data GenesisState) error {
	if data.Params.ArgumentCreationStake.Denom != app.StakeDenom {
		return ErrInvalidArgumentStakeDenom
	}
	if data.Params.UpvoteStake.Denom != app.StakeDenom {
		return ErrInvalidUpvoteStakeDenom
	}
	return nil
}
