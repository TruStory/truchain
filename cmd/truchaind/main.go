package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	gaiaInit "github.com/cosmos/cosmos-sdk/cmd/gaia/init"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"

	"github.com/TruStory/truchain/app"
	"github.com/TruStory/truchain/x/category"
	"github.com/TruStory/truchain/x/story"
	"github.com/cosmos/cosmos-sdk/client"

	"github.com/cosmos/cosmos-sdk/server"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/cli"
	"github.com/tendermint/tendermint/libs/common"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
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

	server.AddCommands(
		ctx,
		cdc,
		rootCmd,
		newApp,
		exportAppStateAndTMValidators)

	// prepare and add flags
	rootDir := os.ExpandEnv("$HOME/.truchaind")
	executor := cli.PrepareBaseCmd(rootCmd, "BC", rootDir)

	err := executor.Execute()
	if err != nil {
		// Note: Handle with #870
		panic(err)
	}
}

// InitCmd initializes all files for tendermint and application
func InitCmd(ctx *server.Context, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize genesis config, priv-validator file, and p2p-node file",
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {

			config := ctx.Config
			config.SetRoot(viper.GetString(cli.HomeFlag))

			chainID := viper.GetString(client.FlagChainID)
			if chainID == "" {
				chainID = fmt.Sprintf("test-chain-%v", common.RandStr(6))
			}

			_, pk, err := gaiaInit.InitializeNodeValidatorFiles(config)
			if err != nil {
				return err
			}

			var appState json.RawMessage
			genFile := config.GenesisFile()

			if !viper.GetBool(flagOverwrite) && common.FileExists(genFile) {
				return fmt.Errorf("genesis.json file already exists: %v", genFile)
			}

			genesis := app.GenesisState{
				AuthData:   auth.DefaultGenesisState(),
				BankData:   bank.DefaultGenesisState(),
				Categories: category.DefaultCategories(),
				StoryData:  story.DefaultGenesisState(),
			}

			appState, err = codec.MarshalJSONIndent(cdc, genesis)
			if err != nil {
				return err
			}

			_, _, validator, err := server.SimpleAppGenTx(cdc, pk)
			if err != nil {
				return err
			}

			if err = gaiaInit.ExportGenesisFile(genFile, chainID, []tmtypes.GenesisValidator{validator}, appState); err != nil {
				return err
			}

			fmt.Printf("Initialized truchaind configuration and bootstrapping files in %s...\n", viper.GetString(cli.HomeFlag))
			return nil
		},
	}

	cmd.Flags().String(cli.HomeFlag, app.DefaultNodeHome, "node's home directory")
	cmd.Flags().String(client.FlagChainID, "", "genesis file chain-id, if left blank will be randomly created")
	cmd.Flags().BoolP(flagOverwrite, "o", false, "overwrite the genesis.json file")

	return cmd
}

func newApp(logger log.Logger, db dbm.DB, traceStore io.Writer) abci.Application {
	return app.NewTruChain(logger, db)
}

func exportAppStateAndTMValidators(logger log.Logger, db dbm.DB, _ io.Writer, _ int64, _ bool) (
	json.RawMessage, []tmtypes.GenesisValidator, error) {
	bapp := app.NewTruChain(logger, db)
	return bapp.ExportAppStateAndValidators()
}
