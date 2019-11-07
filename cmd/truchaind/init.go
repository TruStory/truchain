package main

import (
	"os"

	"github.com/TruStory/truchain/app"
	truchain "github.com/TruStory/truchain/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutilcli "github.com/cosmos/cosmos-sdk/x/genutil/client/cli"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/libs/cli"
	"github.com/tendermint/tendermint/types"
)

func InitCmd(ctx *server.Context, cdc *codec.Codec, mbm module.BasicManager, defaultNodeHome string) *cobra.Command {
	init := genutilcli.InitCmd(ctx, cdc, app.ModuleBasics, app.DefaultNodeHome)
	init.PostRunE = func(cmd *cobra.Command, args []string) error {
		config := ctx.Config
		config.SetRoot(viper.GetString(cli.HomeFlag))
		genFile := config.GenesisFile()
		genDoc := &types.GenesisDoc{}

		if _, err := os.Stat(genFile); err != nil {
			if !os.IsNotExist(err) {
				return err
			}
		} else {
			genDoc, err = types.GenesisDocFromFile(genFile)
			if err != nil {
				return errors.Wrap(err, "Failed to read genesis doc from file")
			}
		}
		var appState genutil.AppMap
		if err := cdc.UnmarshalJSON(genDoc.AppState, &appState); err != nil {
			return errors.Wrap(err, "failed to JSON unmarshal initial genesis state")
		}

		if err := genutil.ExportGenesisFile(genDoc, genFile); err != nil {
			return errors.Wrap(err, "Failed to export gensis file")
		}

		cdc := codec.New()
		codec.RegisterCrypto(cdc)
		// migrate staking state
		if appState[staking.ModuleName] != nil {
			var stakingGenState staking.GenesisState
			cdc.MustUnmarshalJSON(appState[staking.ModuleName], &stakingGenState)
			stakingGenState.Params.BondDenom = truchain.StakeDenom
			appState[staking.ModuleName] = cdc.MustMarshalJSON(stakingGenState)
		}
		var err error
		genDoc.AppState, err = cdc.MarshalJSON(appState)
		if err != nil {
			return errors.Wrap(err, "failed to JSON marshal migrated genesis state")
		}
		if err = genutil.ExportGenesisFile(genDoc, genFile); err != nil {
			return errors.Wrap(err, "Failed to export gensis file")
		}
		return nil
	}
	return init
}
