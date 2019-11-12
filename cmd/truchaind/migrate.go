package main

import (
	"fmt"
	"io/ioutil"
	"sort"
	"strings"
	"time"

	"github.com/cosmos/cosmos-sdk/x/auth"
	authexported "github.com/cosmos/cosmos-sdk/x/auth/exported"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/tendermint/tendermint/types"

	"github.com/TruStory/truchain/x/account"
	trubank "github.com/TruStory/truchain/x/bank"
	trubankexported "github.com/TruStory/truchain/x/bank/exported"
	"github.com/TruStory/truchain/x/claim"
	truslashing "github.com/TruStory/truchain/x/slashing"
	trustaking "github.com/TruStory/truchain/x/staking"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	extypes "github.com/cosmos/cosmos-sdk/x/genutil"
	"github.com/cosmos/cosmos-sdk/x/gov"
	"github.com/cosmos/cosmos-sdk/x/mint"
	"github.com/cosmos/cosmos-sdk/x/staking"
)

const (
	flagGenesisTime = "genesis-time"
	flagChainID     = "chain-id"
)

// Copied from truchain/truapi until truapi is moved into Octopus
func humanReadable(coin sdk.Coin, prec int64) string {
	// empty struct
	if (sdk.Coin{}) == coin {
		return "0"
	}
	shanevs := sdk.NewDecFromIntWithPrec(coin.Amount, prec).String()
	parts := strings.Split(shanevs, ".")
	number := parts[0]
	decimal := parts[1]
	// If greater than 1.0 => show two decimal digits, truncate trailing zeros
	displayDecimalPlaces := 2
	if number == "0" {
		// If less than 1.0 => show four decimal digits, truncate trailing zeros
		displayDecimalPlaces = 4
	}
	decimal = strings.TrimRight(decimal, "0")
	numberOfDecimalPlaces := len(decimal)
	if numberOfDecimalPlaces > displayDecimalPlaces {
		numberOfDecimalPlaces = displayDecimalPlaces
	}
	decimal = decimal[0:numberOfDecimalPlaces]
	decimal = strings.TrimRight(decimal, "0")
	if decimal == "" {
		return number
	}
	return fmt.Sprintf("%s%s%s", number, ".", decimal)
}

var denom = "utru"

func convert(coin sdk.Coin) (sdk.Coin, error) {
	if coin.Amount.IsZero() {
		return sdk.NewInt64Coin(denom, 0), nil
	}
	return sdk.ConvertCoin(coin, denom)
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

func Migrate(appState genutil.AppMap) genutil.AppMap {
	cdc := codec.New()
	codec.RegisterCrypto(cdc)
	sdk.RegisterDenom("utru", sdk.NewDecWithPrec(1, 6))
	sdk.RegisterDenom("tru", sdk.NewDecWithPrec(1, 9))
	auth.RegisterCodec(cdc)
	supply.RegisterCodec(cdc)
	earned := make(map[string]sdk.Coins, 0)
	getEarnedCoins := func(address string) sdk.Coins {
		e, ok := earned[address]
		if ok {
			return e
		}
		return sdk.NewCoins()
	}

	if appState[trubank.ModuleName] != nil && appState[account.ModuleName] != nil {
		var bankGenState trubank.GenesisState
		cdc.MustUnmarshalJSON(appState[trubank.ModuleName], &bankGenState)
		transactions := make([]trubankexported.Transaction, 0)
		for _, tx := range bankGenState.Transactions {

			txAmount, err := sdk.ConvertCoin(tx.Amount, denom)
			if err != nil {
				panic(err)
			}
			tx.Amount = txAmount
			address := tx.AppAccountAddress.String()
			earned[address] = trackEarnedCoins(tx, getEarnedCoins(address))
			transactions = append(transactions, tx)
		}
		bankGenState.Transactions = transactions
		appState[trubank.ModuleName] = cdc.MustMarshalJSON(&bankGenState)
	}

	if appState[crisis.ModuleName] != nil {
		var crisisGenState crisis.GenesisState
		cdc.MustUnmarshalJSON(appState[trubank.ModuleName], &crisisGenState)
		crisisGenState.ConstantFee = sdk.NewInt64Coin(denom, 1000)
		appState[crisis.ModuleName] = cdc.MustMarshalJSON(&crisisGenState)
	}

	if appState[supply.ModuleName] != nil {
		var supplyGenState supply.GenesisState
		cdc.MustUnmarshalJSON(appState[supply.ModuleName], &supplyGenState)
		supplyGenState.Supply = sdk.NewCoins()
		// sdk.NewPermissionsForAddress("user_growth_tokens_pool", []string{supply.Burner, supply.Staking})
		appState[supply.ModuleName] = cdc.MustMarshalJSON(&supplyGenState)
	}

	if appState[staking.ModuleName] != nil {
		var stakingGenState staking.GenesisState
		cdc.MustUnmarshalJSON(appState[staking.ModuleName], &stakingGenState)
		stakingGenState.Params.BondDenom = denom
		appState[staking.ModuleName] = cdc.MustMarshalJSON(&stakingGenState)
	}

	if appState[gov.ModuleName] != nil {
		var govGenState gov.GenesisState
		cdc.MustUnmarshalJSON(appState[gov.ModuleName], &govGenState)
		govGenState.DepositParams.MinDeposit = sdk.NewCoins(sdk.NewInt64Coin(denom, 1000))
		appState[gov.ModuleName] = cdc.MustMarshalJSON(&govGenState)
	}
	if appState[mint.ModuleName] != nil {
		var mintGenState mint.GenesisState
		cdc.MustUnmarshalJSON(appState[mint.ModuleName], &mintGenState)
		mintGenState.Params.MintDenom = denom
		appState[mint.ModuleName] = cdc.MustMarshalJSON(&mintGenState)
	}
	if appState[distribution.ModuleName] != nil {
		var distrGenState distribution.GenesisState
		cdc.MustUnmarshalJSON(appState[distribution.ModuleName], &distrGenState)
		distrGenState.FeePool.CommunityPool = sdk.NewDecCoins(sdk.NewCoins())
		appState[distribution.ModuleName] = cdc.MustMarshalJSON(&distrGenState)
	}

	// migrate all stakes
	if appState[trustaking.ModuleName] != nil {
		var stakingGenState trustaking.GenesisState
		cdc.MustUnmarshalJSON(appState[trustaking.ModuleName], &stakingGenState)

		argCreationStake, err := sdk.ConvertCoin(stakingGenState.Params.ArgumentCreationStake, denom)
		if err != nil {
			panic(err)
		}

		upvoteStake, err := sdk.ConvertCoin(stakingGenState.Params.UpvoteStake, denom)
		if err != nil {
			panic(err)
		}

		stakingGenState.Params.ArgumentCreationStake = argCreationStake
		stakingGenState.Params.UpvoteStake = upvoteStake

		arguments := make([]trustaking.Argument, 0, len(stakingGenState.Arguments))
		for _, arg := range stakingGenState.Arguments {
			totalStake, err := sdk.ConvertCoin(arg.TotalStake, denom)
			if err != nil {
				panic(err)
			}
			upvotedStake, err := sdk.ConvertCoin(arg.UpvotedStake, denom)
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
			utru, err := sdk.ConvertCoin(s.Amount, denom)
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

	if appState[auth.ModuleName] != nil {
		var authGenState auth.GenesisState
		cdc.MustUnmarshalJSON(appState[auth.ModuleName], &authGenState)
		accounts := make([]authexported.GenesisAccount, 0)
		for _, acc := range authGenState.Accounts {
			if supply.NewModuleAddress("user_stakes_tokens_pool").Equals(acc.GetAddress()) {
				fmt.Println("user_stakes_tokens_pool", acc.GetCoins())
			}
			if supply.NewModuleAddress("user_growth_tokens_pool").Equals(acc.GetAddress()) {
				fmt.Println("user_growth_tokens_pool", acc.GetCoins())
			}
			if supply.NewModuleAddress("bonded_tokens_pool").Equals(acc.GetAddress()) {
				fmt.Println("bonded_tokens_pool", acc.GetCoins())
			}
			if supply.NewModuleAddress("user_reward_tokens_pool").Equals(acc.GetAddress()) {
				fmt.Println("user_reward_tokens_pool", acc.GetCoins())
			}

			utru, err := convert(sdk.NewCoin("tru", acc.GetCoins().AmountOf("tru")))
			if err != nil {
				panic(err)
			}
			err = acc.SetCoins(sdk.Coins{utru})
			if err != nil {
				panic(err)
			}
			accounts = append(accounts, acc)
		}
		authGenState.Accounts = accounts
		appState[auth.ModuleName] = cdc.MustMarshalJSON(authGenState)
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
		slashMinStake, err := sdk.ConvertCoin(slashingGenState.Params.SlashMinStake, denom)
		if err != nil {
			panic(err)
		}
		slashingGenState.Params.SlashMinStake = slashMinStake
		appState[truslashing.ModuleName] = cdc.MustMarshalJSON(slashingGenState)
	}
	for _, v := range earned["cosmos1xqc5gwzpfp8ygkzdfdpnq4j3xd8y6djy5z8gfn"] {
		fmt.Println("community", v.Denom, humanReadable(v, 6))
	}
	fmt.Println("earned", earned["cosmos1xqc5gwzpfp8ygkzdfdpnq4j3xd8y6djy5z8gfn"].String())
	return appState
}

// Allow applications to extend and modify the migration process.
//
// Ref: https://github.com/cosmos/cosmos-sdk/issues/5041
var migrationMap = extypes.MigrationMap{
	"v0.1.29": Migrate,
}

// GetMigrationCallback returns a MigrationCallback for a given version.
func GetMigrationCallback(version string) extypes.MigrationCallback {
	return migrationMap[version]
}

// GetMigrationVersions get all migration version in a sorted slice.
func GetMigrationVersions() []string {
	versions := make([]string, len(migrationMap))

	var i int
	for version := range migrationMap {
		versions[i] = version
		i++
	}

	sort.Strings(versions)
	return versions
}

// MigrateGenesisCmd returns a command to execute genesis state migration.
// nolint: funlen
func MigrateGenesisCmd(_ *server.Context, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tru_migrate [target-version] [genesis-file]",
		Short: "Migrate genesis to a specified target version",
		Long: fmt.Sprintf(`Migrate the source genesis into the target version and print to STDOUT.
Example:
$ %s migrate v0.36 /path/to/genesis.json --chain-id=cosmoshub-3 --genesis-time=2019-04-22T17:00:00Z
`, version.ServerName),
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error

			target := args[0]
			importGenesis := args[1]

			genDoc, err := types.GenesisDocFromFile(importGenesis)
			if err != nil {
				return errors.Wrapf(err, "failed to read genesis document from file %s", importGenesis)
			}

			var initialState extypes.AppMap
			if err := cdc.UnmarshalJSON(genDoc.AppState, &initialState); err != nil {
				return errors.Wrap(err, "failed to JSON unmarshal initial genesis state")
			}

			migrationFunc := GetMigrationCallback(target)
			if migrationFunc == nil {
				return fmt.Errorf("unknown migration function for version: %s", target)
			}

			// TODO: handler error from migrationFunc call
			newGenState := migrationFunc(initialState)

			genDoc.AppState, err = cdc.MarshalJSON(newGenState)
			if err != nil {
				return errors.Wrap(err, "failed to JSON marshal migrated genesis state")
			}

			genesisTime := cmd.Flag(flagGenesisTime).Value.String()
			if genesisTime != "" {
				var t time.Time

				err := t.UnmarshalText([]byte(genesisTime))
				if err != nil {
					return errors.Wrap(err, "failed to unmarshal genesis time")
				}

				genDoc.GenesisTime = t
			}

			chainID := cmd.Flag(flagChainID).Value.String()
			if chainID != "" {
				genDoc.ChainID = chainID
			}

			bz, err := cdc.MarshalJSONIndent(genDoc, "", "\t")
			if err != nil {
				return errors.Wrap(err, "failed to marshal genesis doc")
			}

			sortedBz, err := sdk.SortJSON(bz)
			if err != nil {
				return errors.Wrap(err, "failed to sort JSON genesis doc")
			}
			// os.New
			fmt.Println(len(sortedBz))
			err = ioutil.WriteFile("migrated.json", sortedBz, 0644)
			// fmt.Println(string(sortedBz))
			return err
		},
	}

	cmd.Flags().String(flagGenesisTime, "", "override genesis_time with this flag")
	cmd.Flags().String(flagChainID, "", "override chain_id with this flag")

	return cmd
}
