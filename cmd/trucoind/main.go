package main

import (
	"encoding/json"
	"io"
	"os"

	"github.com/TruStory/trucoin/app"
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
		Use:               "trucoind",
		Short:             "TruStoryApp Daemon (server)",
		PersistentPreRunE: server.PersistentPreRunEFn(ctx),
	}

	server.AddCommands(ctx, cdc, rootCmd, server.DefaultAppInit,
		server.ConstructAppCreator(newApp, "trucoin"),
		server.ConstructAppExporter(exportAppStateAndTMValidators, "trucoin"))

	// prepare and add flags
	rootDir := os.ExpandEnv("$HOME/.trucoind")
	executor := cli.PrepareBaseCmd(rootCmd, "BC", rootDir)

	err := executor.Execute()
	if err != nil {
		// Note: Handle with #870
		panic(err)
	}
}

func newApp(logger log.Logger, db dbm.DB, _ io.Writer) abci.Application {
	return app.NewTruStoryApp(logger, db)
}

func exportAppStateAndTMValidators(logger log.Logger, db dbm.DB, _ io.Writer) (json.RawMessage, []tmtypes.GenesisValidator, error) {
	bapp := app.NewTruStoryApp(logger, db)
	return bapp.ExportAppStateAndValidators()
}
