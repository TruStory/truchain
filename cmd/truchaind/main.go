package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/TruStory/truchain/app"
	truchainInit "github.com/TruStory/truchain/cmd/truchaind/init"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	abci "github.com/tendermint/tendermint/abci/types"
	cfg "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/libs/cli"
	"github.com/tendermint/tendermint/libs/common"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tendermint/p2p"
	"github.com/tendermint/tendermint/privval"
	"github.com/tendermint/tendermint/types"
	tmtypes "github.com/tendermint/tendermint/types"
)

const (
	flagClientHome = "home-client"
	flagOverwrite  = "overwrite"
)

func main() {
	cobra.EnableCommandSorting = false

	cdc := app.MakeCodec()
	ctx := server.NewDefaultContext()

	rootCmd := &cobra.Command{
		Use:               "truchaind",
		Short:             "TruChain Daemon (server)",
		PersistentPreRunE: server.PersistentPreRunEFn(ctx),
	}

	rootCmd.AddCommand(InitCmd(ctx, cdc))
	rootCmd.AddCommand(AddGenesisAccountCmd(ctx, cdc))
	rootCmd.AddCommand(truchainInit.TestnetFilesCmd(ctx, cdc))

	server.AddCommands(ctx, cdc, rootCmd, newApp, exportAppStateAndTMValidators)

	// prepare and add flags
	rootDir := os.ExpandEnv("$HOME/.truchaind")
	executor := cli.PrepareBaseCmd(rootCmd, "BC", rootDir)

	err := executor.Execute()
	if err != nil {
		// Note: Handle with #870
		panic(err)
	}
}

// InitCmd returns a command that initializes all files needed for Tendermint
// and the respective application.
func InitCmd(ctx *server.Context, cdc *codec.Codec) *cobra.Command { // nolint: golint
	cmd := &cobra.Command{
		Use:   "init [moniker]",
		Short: "Initialize private validator, p2p, genesis, and application configuration files",
		Long:  `Initialize validators's and node's configuration files.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			config := ctx.Config
			config.SetRoot(viper.GetString(cli.HomeFlag))

			chainID := viper.GetString(client.FlagChainID)
			if chainID == "" {
				chainID = fmt.Sprintf("test-chain-%v", common.RandStr(6))
			}

			nodeID, pk, err := InitializeNodeValidatorFiles(config)
			if err != nil {
				return err
			}

			config.Moniker = args[0]

			var appState json.RawMessage
			genFile := config.GenesisFile()

			if appState, err = initializeEmptyGenesis(cdc, genFile, chainID,
				viper.GetBool(flagOverwrite)); err != nil {
				return err
			}

			genDoc := &types.GenesisDoc{}
			if _, err := os.Stat(genFile); err != nil {
				if !os.IsNotExist(err) {
					return err
				}
			} else {
				genDoc, err = types.GenesisDocFromFile(genFile)
				if err != nil {
					return err
				}
			}

			_, _, validator, err := SimpleAppGenTx(cdc, pk)
			if err != nil {
				return err
			}

			genDoc.ChainID = chainID
			genDoc.Validators = []tmtypes.GenesisValidator{validator}
			genDoc.AppState = appState
			if err = truchainInit.ExportGenesisFile(genDoc, genFile); err != nil {
				return err
			}

			toPrint := truchainInit.NewPrintInfo(config.Moniker, chainID, nodeID, "", appState)

			cfg.WriteConfigFile(filepath.Join(config.RootDir, "config", "config.toml"), config)
			return truchainInit.DisplayInfo(cdc, toPrint)
		},
	}

	cmd.Flags().String(cli.HomeFlag, app.DefaultNodeHome, "node's home directory")
	cmd.Flags().BoolP(flagOverwrite, "o", false, "overwrite the genesis.json file")
	cmd.Flags().String(client.FlagChainID, "", "genesis file chain-id, if left blank will be randomly created")

	return cmd
}

// InitializeNodeValidatorFiles creates private validator and p2p configuration files.
func InitializeNodeValidatorFiles(
	config *cfg.Config) (nodeID string, valPubKey crypto.PubKey, err error,
) {

	nodeKey, err := p2p.LoadOrGenNodeKey(config.NodeKeyFile())
	if err != nil {
		return nodeID, valPubKey, err
	}

	nodeID = string(nodeKey.ID())
	server.UpgradeOldPrivValFile(config)

	pvKeyFile := config.PrivValidatorKeyFile()
	if err := common.EnsureDir(filepath.Dir(pvKeyFile), 0777); err != nil {
		return nodeID, valPubKey, nil
	}

	pvStateFile := config.PrivValidatorStateFile()
	if err := common.EnsureDir(filepath.Dir(pvStateFile), 0777); err != nil {
		return nodeID, valPubKey, nil
	}

	valPubKey = privval.LoadOrGenFilePV(pvKeyFile, pvStateFile).GetPubKey()

	return nodeID, valPubKey, nil
}

func newApp(logger log.Logger, db dbm.DB, traceStore io.Writer) abci.Application {
	return app.NewTruChain(logger, db, true)
}

func exportAppStateAndTMValidators(
	logger log.Logger, db dbm.DB, traceStore io.Writer, height int64, forZeroHeight bool, jailWhiteList []string,
) (json.RawMessage, []tmtypes.GenesisValidator, error) {
	if height != -1 {
		tApp := app.NewTruChain(logger, db, false)
		err := tApp.LoadHeight(height)
		if err != nil {
			return nil, nil, err
		}
		return tApp.ExportAppStateAndValidators()
	}

	tApp := app.NewTruChain(logger, db, true)
	return tApp.ExportAppStateAndValidators()
}

func initializeEmptyGenesis(
	cdc *codec.Codec, genFile, chainID string, overwrite bool,
) (appState json.RawMessage, err error) {

	if !overwrite && common.FileExists(genFile) {
		return nil, fmt.Errorf("genesis.json file already exists: %v", genFile)
	}

	return codec.MarshalJSONIndent(cdc, app.NewDefaultGenesisState())
}

// SimpleAppGenTx returns a simple GenTx command that makes the node a valdiator from the start
func SimpleAppGenTx(cdc *codec.Codec, pk crypto.PubKey) (
	appGenTx, cliPrint json.RawMessage, validator tmtypes.GenesisValidator, err error) {

	addr, secret, err := server.GenerateCoinKey()
	if err != nil {
		return
	}

	bz, err := cdc.MarshalJSON(struct {
		Addr sdk.AccAddress `json:"addr"`
	}{addr})
	if err != nil {
		return
	}

	appGenTx = json.RawMessage(bz)

	bz, err = cdc.MarshalJSON(map[string]string{"secret": secret})
	if err != nil {
		return
	}

	cliPrint = json.RawMessage(bz)

	validator = tmtypes.GenesisValidator{
		PubKey: pk,
		Power:  10,
	}

	return
}

// AddGenesisAccountCmd allows users to add accounts to the genesis file
func AddGenesisAccountCmd(ctx *server.Context, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-genesis-account [address] [coins[,coins]]",
		Short: "Adds an account to the genesis file",
		Args:  cobra.ExactArgs(2),
		Long: strings.TrimSpace(`
Adds accounts to the genesis file so that you can start a chain with coins in the CLI:
$ nsd add-genesis-account cosmos1tse7r2fadvlrrgau3pa0ss7cqh55wrv6y9alwh 1000STAKE,1000nametoken
`),
		RunE: func(_ *cobra.Command, args []string) error {
			addr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}
			coins, err := sdk.ParseCoins(args[1])
			if err != nil {
				return err
			}
			coins.Sort()

			var genDoc tmtypes.GenesisDoc
			config := ctx.Config
			genFile := config.GenesisFile()
			if !common.FileExists(genFile) {
				return fmt.Errorf("%s does not exist, run `gaiad init` first", genFile)
			}
			genContents, err := ioutil.ReadFile(genFile)
			if err != nil {
			}

			if err = cdc.UnmarshalJSON(genContents, &genDoc); err != nil {
				return err
			}

			var appState app.GenesisState
			if err = cdc.UnmarshalJSON(genDoc.AppState, &appState); err != nil {
				return err
			}

			for _, stateAcc := range appState.Accounts {
				if stateAcc.Address.Equals(addr) {
					return fmt.Errorf("the application state already contains account %v", addr)
				}
			}

			acc := app.NewDefaultGenesisAccount(addr)
			acc.Coins = coins
			appState.Accounts = append(appState.Accounts, acc)
			appStateJSON, err := cdc.MarshalJSON(appState)
			if err != nil {
				return err
			}

			genDoc.AppState = appStateJSON
			if err = truchainInit.ExportGenesisFile(&genDoc, genFile); err != nil {
				return err
			}

			return nil
		},
	}
	return cmd
}
