package reset

import (
	"github.com/cosmos/cosmos-sdk/x/auth"
	authexported "github.com/cosmos/cosmos-sdk/x/auth/exported"
	"github.com/cosmos/cosmos-sdk/x/supply"

	app "github.com/TruStory/truchain/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	"github.com/cosmos/cosmos-sdk/x/gov"
	"github.com/cosmos/cosmos-sdk/x/mint"
	"github.com/cosmos/cosmos-sdk/x/staking"
)

// Migrate utru denom
func Migrate(appState genutil.AppMap) genutil.AppMap {
	cdc := codec.New()
	codec.RegisterCrypto(cdc)
	auth.RegisterCodec(cdc)
	supply.RegisterCodec(cdc)
	if appState[supply.ModuleName] != nil {
		var supplyGenState supply.GenesisState
		cdc.MustUnmarshalJSON(appState[supply.ModuleName], &supplyGenState)
		supplyGenState.Supply = sdk.NewCoins()
		appState[supply.ModuleName] = cdc.MustMarshalJSON(&supplyGenState)
	}

	if appState[staking.ModuleName] != nil {
		var stakingGenState staking.GenesisState
		cdc.MustUnmarshalJSON(appState[staking.ModuleName], &stakingGenState)
		stakingGenState.Params.BondDenom = app.StakeDenom
		stakingGenState.Validators = nil
		stakingGenState.LastTotalPower = sdk.NewInt(0)
		stakingGenState.LastValidatorPowers = nil
		stakingGenState.Delegations = nil
		stakingGenState.UnbondingDelegations = nil
		stakingGenState.Redelegations = nil
		stakingGenState.Exported = false
		appState[staking.ModuleName] = cdc.MustMarshalJSON(&stakingGenState)
	}

	if appState[gov.ModuleName] != nil {
		var govGenState gov.GenesisState
		cdc.MustUnmarshalJSON(appState[gov.ModuleName], &govGenState)
		govGenState.DepositParams.MinDeposit = sdk.NewCoins(sdk.NewInt64Coin(app.StakeDenom, 1000))
		appState[gov.ModuleName] = cdc.MustMarshalJSON(&govGenState)
	}
	if appState[mint.ModuleName] != nil {
		var mintGenState mint.GenesisState
		cdc.MustUnmarshalJSON(appState[mint.ModuleName], &mintGenState)
		mintGenState.Minter.AnnualProvisions = sdk.NewDec(0)
		mintGenState.Params.MintDenom = app.StakeDenom
		appState[mint.ModuleName] = cdc.MustMarshalJSON(&mintGenState)
	}
	if appState[distribution.ModuleName] != nil {
		var distrGenState distribution.GenesisState
		cdc.MustUnmarshalJSON(appState[distribution.ModuleName], &distrGenState)
		distrGenState.FeePool.CommunityPool = sdk.NewDecCoins(sdk.NewCoins())
		distrGenState.PreviousProposer = sdk.ConsAddress{}
		distrGenState.DelegatorWithdrawInfos = make([]distribution.DelegatorWithdrawInfo, 0)
		distrGenState.DelegatorStartingInfos = make([]distribution.DelegatorStartingInfoRecord, 0)
		distrGenState.OutstandingRewards = make([]distribution.ValidatorOutstandingRewardsRecord, 0)
		distrGenState.ValidatorAccumulatedCommissions = make([]distribution.ValidatorAccumulatedCommissionRecord, 0)
		distrGenState.ValidatorHistoricalRewards = make([]distribution.ValidatorHistoricalRewardsRecord, 0)
		distrGenState.ValidatorCurrentRewards = make([]distribution.ValidatorCurrentRewardsRecord, 0)
		distrGenState.ValidatorSlashEvents = make([]distribution.ValidatorSlashEventRecord, 0)
		appState[distribution.ModuleName] = cdc.MustMarshalJSON(&distrGenState)
	}

	if appState[auth.ModuleName] != nil {
		var authGenState auth.GenesisState
		cdc.MustUnmarshalJSON(appState[auth.ModuleName], &authGenState)
		accounts := make([]authexported.GenesisAccount, 0)
		moduleAccounts := []string{"user_stakes_tokens_pool", "user_growth_tokens_pool",
			"bonded_tokens_pool", "user_reward_tokens_pool", "not_bonded_tokens_pool", "distribution"}

		for _, acc := range authGenState.Accounts {
			for _, ma := range moduleAccounts {
				if supply.NewModuleAddress(ma).Equals(acc.GetAddress()) {
					err := acc.SetCoins(sdk.NewCoins())
					if err != nil {
						panic(err)
					}
				}
			}
			accounts = append(accounts, acc)
		}
		authGenState.Accounts = accounts
		appState[auth.ModuleName] = cdc.MustMarshalJSON(authGenState)
	}
	if appState[genutil.ModuleName] == nil {
		genState := genutil.NewGenesisState(nil)
		appState[genutil.ModuleName] = cdc.MustMarshalJSON(&genState)
	}

	return appState
}
