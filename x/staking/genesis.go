package staking

import (
	"fmt"

	"github.com/TruStory/truchain/x/account"
	bankexported "github.com/TruStory/truchain/x/bank/exported"

	app "github.com/TruStory/truchain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
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
		k.setArgument(ctx, a)
		k.setClaimArgument(ctx, a.ClaimID, a.ID)
		k.setUserArgument(ctx, a.Creator, a.ID)
	}
	for _, s := range data.Stakes {
		k.setStake(ctx, s)
		if !s.Expired {
			k.InsertActiveStakeQueue(ctx, s.ID, s.EndTime)
			err := k.supplyKeeper.MintCoins(ctx, UserStakesPoolName, sdk.NewCoins(s.Amount))
			if err != nil {
				panic(err)
			}
		}
		k.setArgumentStake(ctx, s.ArgumentID, s.ID)
		k.setUserStake(ctx, s.Creator, s.CreatedTime, s.ID)

		arg, ok := k.Argument(ctx, s.ArgumentID)
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

	err := initUserRewardsPool(ctx, k)
	if err != nil {
		panic(err)
	}
}

func initUserRewardsPool(ctx sdk.Context, keeper Keeper) sdk.Error {
	userGrowthAcc := keeper.supplyKeeper.GetModuleAccount(ctx, UserRewardPoolName)
	if userGrowthAcc.GetCoins().Empty() {
		amount := app.NewShanevCoin(5000000)
		err := keeper.supplyKeeper.MintCoins(ctx, UserRewardPoolName, sdk.NewCoins(amount))
		if err != nil {
			return err
		}

		keeper.accountKeeper.IterateAppAccounts(ctx, func(acc account.AppAccount) (stop bool) {
			addr := acc.PrimaryAddress()
			//fmt.Println(addr.String())
			keeper.bankKeeper.IterateUserTransactions(ctx, addr, false, func(tx bankexported.Transaction) bool {
				switch tx.Type {
				case TransactionInterestArgumentCreation, TransactionInterestUpvoteGiven,
					TransactionInterestUpvoteReceived:
					//fmt.Println("processing " + tx.Type.String() + " " + tx.Amount.String())
					err := keeper.supplyKeeper.BurnCoins(ctx, UserRewardPoolName, sdk.NewCoins(tx.Amount))
					if err != nil {
						panic(err)
					}
				}
				return false
			})
			return false
		})
	}
	return nil
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
