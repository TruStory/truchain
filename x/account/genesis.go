package account

import (
	"fmt"

	app "github.com/TruStory/truchain/types"
	bankexported "github.com/TruStory/truchain/x/bank/exported"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authexported "github.com/cosmos/cosmos-sdk/x/auth/exported"
)

// GenesisState defines genesis data for the module
type GenesisState struct {
	AppAccounts []AppAccount `json:"app_accounts"`
	Params      Params       `json:"params"`
}

// NewGenesisState creates a new genesis state.
func NewGenesisState() GenesisState {
	return GenesisState{
		AppAccounts: nil,
		Params:      DefaultParams(),
	}
}

// DefaultGenesisState returns a default genesis state
func DefaultGenesisState() GenesisState { return NewGenesisState() }

// InitGenesis initializes account state from genesis file
func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) {
	for _, acc := range data.AppAccounts {
		keeper.setAppAccount(ctx, acc)
		if acc.IsJailed {
			keeper.setJailEndTimeAccount(ctx, acc.JailEndTime, acc.PrimaryAddress())
		}
	}
	keeper.SetParams(ctx, data.Params)

	err := initUserGrowthPool(ctx, keeper)
	if err != nil {
		panic(err)
	}
}

func initUserGrowthPool(ctx sdk.Context, keeper Keeper) sdk.Error {
	userGrowthAcc := keeper.supplyKeeper.GetModuleAccount(ctx, UserGrowthPoolName)
	if userGrowthAcc.GetCoins().Empty() {
		fmt.Println("Setting user growth pool")
		amount := app.NewShanevCoin(5_000_000)
		err := keeper.supplyKeeper.MintCoins(ctx, UserGrowthPoolName, sdk.NewCoins(amount))
		if err != nil {
			return err
		}

		userBalanceTotal := sdk.NewCoin(app.StakeDenom, sdk.ZeroInt())
		keeper.accountKeeper.IterateAccounts(ctx, func(acc authexported.Account) (stop bool) {
			addr := acc.GetAddress()
			amt := acc.GetCoins().AmountOf(app.StakeDenom)
			userBalanceTotal = userBalanceTotal.Add(sdk.NewCoin(app.StakeDenom, amt))
			keeper.bankKeeper.IterateUserTransactions(ctx, addr, false, func(tx bankexported.Transaction) bool {
				if tx.Type == bankexported.TransactionGift {
					err := keeper.supplyKeeper.BurnCoins(ctx, UserGrowthPoolName, sdk.NewCoins(tx.Amount))
					if err != nil {
						panic(err)
					}
				}
				return false
			})
			return false
		})

		fmt.Println("Completed: Setting user growth pool")

	}
	return nil
}

// ExportGenesis exports the genesis state
func ExportGenesis(ctx sdk.Context, keeper Keeper) GenesisState {
	return GenesisState{
		AppAccounts: keeper.AppAccounts(ctx),
		Params:      keeper.GetParams(ctx),
	}
}

// ValidateGenesis validates the genesis state data
func ValidateGenesis(data GenesisState) error {
	if len(data.Params.Registrar) == 0 {
		return fmt.Errorf("Param: Registrar, must be a valid address")
	}

	if data.Params.MaxSlashCount < 1 {
		return fmt.Errorf("Param: MaxSlashCount, must have a positive value")
	}

	if data.Params.JailDuration.Seconds() < 1 {
		return fmt.Errorf("Param: JailTime, must have a positive value")
	}

	return nil
}
