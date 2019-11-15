package v0_3

import (
	"github.com/cosmos/cosmos-sdk/x/auth"
	authexported "github.com/cosmos/cosmos-sdk/x/auth/exported"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	"github.com/cosmos/cosmos-sdk/x/supply"

	"github.com/TruStory/truchain/x/account"
	trubank "github.com/TruStory/truchain/x/bank"
	trubankexported "github.com/TruStory/truchain/x/bank/exported"
	"github.com/TruStory/truchain/x/claim"
	truslashing "github.com/TruStory/truchain/x/slashing"
	trustaking "github.com/TruStory/truchain/x/staking"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	"github.com/cosmos/cosmos-sdk/x/gov"
	"github.com/cosmos/cosmos-sdk/x/mint"
	"github.com/cosmos/cosmos-sdk/x/staking"
)

const (
	targetDenom       = "utru"
	currentDenom      = "tru"
	oldBankModuleName = "trubank2"
)

func convert(coin sdk.Coin) (sdk.Coin, error) {
	if coin.Amount.IsZero() {
		return sdk.NewInt64Coin(targetDenom, 0), nil
	}
	return sdk.ConvertCoin(coin, targetDenom)
}

func trackEarnedCoins(tx trubank.Transaction, coins sdk.Coins) sdk.Coins {
	coins = coins.Sort()
	earned := tx.Type.OneOf(trubankexported.AllowedTransactionsForEarning)
	lost := tx.Type.OneOf(trubankexported.AllowedTransactionsForEarningDeduction)
	if earned {
		earnedCoin := sdk.NewCoin(tx.CommunityID, tx.Amount.Amount)
		return coins.Add(sdk.NewCoins(earnedCoin))
	}
	if lost {
		lostCoin := sdk.NewCoin(tx.CommunityID, tx.Amount.Amount)
		return coins.Sub(sdk.NewCoins(lostCoin))
	}
	return coins
}

// Migrate utru denom
func Migrate(appState genutil.AppMap) genutil.AppMap {
	cdc := codec.New()
	codec.RegisterCrypto(cdc)
	err := sdk.RegisterDenom(targetDenom, sdk.NewDecWithPrec(1, 6))
	if err != nil {
		panic(err)
	}
	err = sdk.RegisterDenom(currentDenom, sdk.NewDecWithPrec(1, 9))
	if err != nil {
		panic(err)
	}

	auth.RegisterCodec(cdc)
	supply.RegisterCodec(cdc)
	earned := make(map[string]sdk.Coins)
	getEarnedCoins := func(address string) sdk.Coins {
		e, ok := earned[address]
		if ok {
			return e
		}
		return sdk.NewCoins()
	}

	// cosmos modules
	if appState[crisis.ModuleName] != nil {
		var crisisGenState crisis.GenesisState
		cdc.MustUnmarshalJSON(appState[crisis.ModuleName], &crisisGenState)
		crisisGenState.ConstantFee = sdk.NewInt64Coin(targetDenom, 1000)
		appState[crisis.ModuleName] = cdc.MustMarshalJSON(&crisisGenState)
	}

	if appState[supply.ModuleName] != nil {
		var supplyGenState supply.GenesisState
		cdc.MustUnmarshalJSON(appState[supply.ModuleName], &supplyGenState)
		supplyGenState.Supply = sdk.NewCoins()
		appState[supply.ModuleName] = cdc.MustMarshalJSON(&supplyGenState)
	}

	if appState[staking.ModuleName] != nil {
		var stakingGenState staking.GenesisState
		cdc.MustUnmarshalJSON(appState[staking.ModuleName], &stakingGenState)
		stakingGenState.Params.BondDenom = targetDenom
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
		govGenState.DepositParams.MinDeposit = sdk.NewCoins(sdk.NewInt64Coin(targetDenom, 1000))
		appState[gov.ModuleName] = cdc.MustMarshalJSON(&govGenState)
	}
	if appState[mint.ModuleName] != nil {
		var mintGenState mint.GenesisState
		cdc.MustUnmarshalJSON(appState[mint.ModuleName], &mintGenState)
		mintGenState.Params.MintDenom = targetDenom
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
			var skipConvertion bool
			for _, ma := range moduleAccounts {
				if supply.NewModuleAddress(ma).Equals(acc.GetAddress()) {
					skipConvertion = true
					err := acc.SetCoins(sdk.NewCoins())
					if err != nil {
						panic(err)
					}
				}
			}
			if !skipConvertion {
				utru, err := convert(sdk.NewCoin(currentDenom, acc.GetCoins().AmountOf(currentDenom)))
				if err != nil {
					panic(err)
				}
				err = acc.SetCoins(sdk.Coins{utru})
				if err != nil {
					panic(err)
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
	trubankModuleName := trubank.ModuleName
	if appState[oldBankModuleName] != nil {
		trubankModuleName = oldBankModuleName
	}
	// trustory modules
	if appState[trubankModuleName] != nil && appState[account.ModuleName] != nil {
		var bankGenState trubank.GenesisState
		cdc.MustUnmarshalJSON(appState[trubankModuleName], &bankGenState)
		transactions := make([]trubankexported.Transaction, 0)
		for _, tx := range bankGenState.Transactions {

			txAmount, err := sdk.ConvertCoin(tx.Amount, targetDenom)
			if err != nil {
				panic(err)
			}
			tx.Amount = txAmount
			address := tx.AppAccountAddress.String()
			earned[address] = trackEarnedCoins(tx, getEarnedCoins(address))
			transactions = append(transactions, tx)
		}
		bankGenState.Transactions = transactions
		delete(appState, trubankModuleName)
		appState[trubank.ModuleName] = cdc.MustMarshalJSON(&bankGenState)
	}

	if appState[trustaking.ModuleName] != nil {
		var stakingGenState trustaking.GenesisState
		cdc.MustUnmarshalJSON(appState[trustaking.ModuleName], &stakingGenState)

		argCreationStake, err := sdk.ConvertCoin(stakingGenState.Params.ArgumentCreationStake, targetDenom)
		if err != nil {
			panic(err)
		}

		upvoteStake, err := sdk.ConvertCoin(stakingGenState.Params.UpvoteStake, targetDenom)
		if err != nil {
			panic(err)
		}

		stakingGenState.Params.ArgumentCreationStake = argCreationStake
		stakingGenState.Params.UpvoteStake = upvoteStake

		arguments := make([]trustaking.Argument, 0, len(stakingGenState.Arguments))
		for _, arg := range stakingGenState.Arguments {
			totalStake, err := sdk.ConvertCoin(arg.TotalStake, targetDenom)
			if err != nil {
				panic(err)
			}
			upvotedStake, err := sdk.ConvertCoin(arg.UpvotedStake, targetDenom)
			if err != nil {
				panic(err)
			}
			arg.TotalStake = totalStake
			arg.UpvotedStake = upvotedStake
			arguments = append(arguments, arg)
		}
		stakingGenState.Arguments = arguments
		stakes := make([]trustaking.Stake, 0, len(stakingGenState.Stakes))
		for _, s := range stakingGenState.Stakes {
			utru, err := sdk.ConvertCoin(s.Amount, targetDenom)
			if err != nil {
				panic(err)
			}
			s.Amount = utru
			if s.Result != nil {
				creatorReward, err := convert(s.Result.ArgumentCreatorReward)
				if err != nil {
					panic(err)
				}
				stakerReward, err := convert(s.Result.StakeCreatorReward)
				if err != nil {
					panic(err)
				}
				s.Result.ArgumentCreatorReward = creatorReward
				s.Result.StakeCreatorReward = stakerReward
			}
			stakes = append(stakes, s)
		}
		userEarnedCoins := make([]trustaking.UserEarnedCoins, 0, len(stakingGenState.UsersEarnings))
		for a, c := range earned {
			address, err := sdk.AccAddressFromBech32(a)
			if err != nil {
				panic(err)
			}
			userEarnedCoins = append(userEarnedCoins, trustaking.UserEarnedCoins{Address: address, Coins: c})
		}
		stakingGenState.UsersEarnings = userEarnedCoins
		stakingGenState.Stakes = stakes
		appState[trustaking.ModuleName] = cdc.MustMarshalJSON(&stakingGenState)
	}

	if appState[claim.ModuleName] != nil {
		var claimGenState claim.GenesisState
		cdc.MustUnmarshalJSON(appState[claim.ModuleName], &claimGenState)
		claims := make([]claim.Claim, 0, len(claimGenState.Claims))
		for _, c := range claimGenState.Claims {
			totalBacked, err := convert(c.TotalBacked)
			if err != nil {
				panic(err)
			}

			totalChallenged, err := convert(c.TotalChallenged)
			if err != nil {
				panic(err)
			}
			c.TotalBacked = totalBacked
			c.TotalChallenged = totalChallenged
			claims = append(claims, c)

		}
		claimGenState.Claims = claims
		appState[claim.ModuleName] = cdc.MustMarshalJSON(claimGenState)
	}

	if appState[truslashing.ModuleName] != nil {
		var slashingGenState truslashing.GenesisState
		cdc.MustUnmarshalJSON(appState[truslashing.ModuleName], &slashingGenState)
		slashMinStake, err := sdk.ConvertCoin(slashingGenState.Params.SlashMinStake, targetDenom)
		if err != nil {
			panic(err)
		}
		slashingGenState.Params.SlashMinStake = slashMinStake
		appState[truslashing.ModuleName] = cdc.MustMarshalJSON(slashingGenState)
	}
	return appState
}
