package main

import (
	"encoding/json"
	"io"
	"os"

	"github.com/TruStory/truchain/app"
	truchainInit "github.com/TruStory/truchain/cmd/truchaind/init"
	"github.com/TruStory/truchain/x/argument"
	"github.com/TruStory/truchain/x/category"
	"github.com/TruStory/truchain/x/challenge"
	"github.com/TruStory/truchain/x/expiration"
	"github.com/TruStory/truchain/x/stake"
	"github.com/TruStory/truchain/x/story"
	"github.com/cosmos/cosmos-sdk/client"
	gaiaInit "github.com/cosmos/cosmos-sdk/cmd/gaia/init"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/spf13/cobra"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/cli"
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

<<<<<<< HEAD
	rootCmd.AddCommand(InitCmd(ctx, cdc))
	rootCmd.AddCommand(truchainInit.TestnetFilesCmd(ctx, cdc))
	// rootCmd.AddCommand(InitCmd(ctx, cdc))
>>>>>>> commented out init code
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
// func InitCmd(ctx *server.Context, cdc *codec.Codec) *cobra.Command {
// 	cmd := &cobra.Command{
// 		Use:   "init",
// 		Short: "Initialize genesis config, priv-validator file, and p2p-node file",
// 		Args:  cobra.NoArgs,
// 		RunE: func(_ *cobra.Command, _ []string) error {

// 			config := ctx.Config
// 			config.SetRoot(viper.GetString(cli.HomeFlag))

// 			chainID := viper.GetString(client.FlagChainID)
// 			if chainID == "" {
// 				chainID = fmt.Sprintf("test-chain-%v", common.RandStr(6))
// 			}

// 			_, pk, err := gaiaInit.InitializeNodeValidatorFiles(config)
// 			if err != nil {
// 				return err
// 			}

// 			var appState json.RawMessage
// 			genFile := config.GenesisFile()

// 			if !viper.GetBool(flagOverwrite) && common.FileExists(genFile) {
// 				return fmt.Errorf("genesis.json file already exists: %v", genFile)
// 			}
// 			genesis := app.GenesisState{
// 				ArgumentData:   argument.DefaultGenesisState(),
// 				AuthData:       auth.DefaultGenesisState(),
// 				BankData:       bank.DefaultGenesisState(),
// 				CategoryData:   category.DefaultGenesisState(),
// 				ChallengeData:  challenge.DefaultGenesisState(),
// 				ExpirationData: expiration.DefaultGenesisState(),
// 				StakeData:      stake.DefaultGenesisState(),
// 				StoryData:      story.DefaultGenesisState(),
// 			}

// 			appState, err = codec.MarshalJSONIndent(cdc, genesis)
// 			if err != nil {
// 				return err
// 			}

// 			_, _, validator, err := server.SimpleAppGenTx(cdc, pk)
// 			if err != nil {
// 				return err
// 			}

// 			if err = gaiaInit.ExportGenesisFile(genFile, chainID, []tmtypes.GenesisValidator{validator}, appState); err != nil {
// 				return err
// 			}

// 			fmt.Printf("Initialized truchaind configuration and bootstrapping files in %s...\n", viper.GetString(cli.HomeFlag))
// 			return nil
// 		},
// 	}

// 	cmd.Flags().String(cli.HomeFlag, app.DefaultNodeHome, "node's home directory")
// 	cmd.Flags().String(client.FlagChainID, "", "genesis file chain-id, if left blank will be randomly created")
// 	cmd.Flags().BoolP(flagOverwrite, "o", false, "overwrite the genesis.json file")

// 	return cmd
// }

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
