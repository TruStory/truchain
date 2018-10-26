package main

import (
	"encoding/json"
	"io"
	"os"

	gaiaInit "github.com/cosmos/cosmos-sdk/cmd/gaia/init"

	"github.com/TruStory/truchain/app"

	"github.com/cosmos/cosmos-sdk/server"
	"github.com/spf13/cobra"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/cli"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
	tmtypes "github.com/tendermint/tendermint/types"
)

func main() {
	cdc := app.MakeCodec()
	ctx := server.NewDefaultContext()

	rootCmd := &cobra.Command{
		Use:               "truchaind",
		Short:             "TruChain Daemon (server)",
		PersistentPreRunE: server.PersistentPreRunEFn(ctx),
	}

	appInit := server.DefaultAppInit
	rootCmd.AddCommand(gaiaInit.InitCmd(ctx, cdc, appInit))
	rootCmd.AddCommand(gaiaInit.TestnetFilesCmd(ctx, cdc, appInit))

	server.AddCommands(
		ctx,
		cdc,
		rootCmd,
		appInit,
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

func newApp(logger log.Logger, db dbm.DB, _ io.Writer) abci.Application {
	return app.NewTruChain(logger, db)
}

func exportAppStateAndTMValidators(logger log.Logger, db dbm.DB, _ io.Writer) (json.RawMessage, []tmtypes.GenesisValidator, error) {
	bapp := app.NewTruChain(logger, db)
	return bapp.ExportAppStateAndValidators()
}
