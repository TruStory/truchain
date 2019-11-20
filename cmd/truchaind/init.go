package main

import (
	"github.com/TruStory/truchain/app"
	truchain "github.com/TruStory/truchain/types"
	"github.com/TruStory/truchain/x/account"
	trubank "github.com/TruStory/truchain/x/bank"
	"github.com/TruStory/truchain/x/claim"
	"github.com/TruStory/truchain/x/community"
	truslashing "github.com/TruStory/truchain/x/slashing"
	trustaking "github.com/TruStory/truchain/x/staking"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutilcli "github.com/cosmos/cosmos-sdk/x/genutil/client/cli"
	"github.com/cosmos/cosmos-sdk/x/gov"
	"github.com/cosmos/cosmos-sdk/x/mint"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/libs/cli"
	"github.com/tendermint/tendermint/types"
	"os"
)

func InitCmd(ctx *server.Context, cdc *codec.Codec, mbm module.BasicManager, defaultNodeHome string) *cobra.Command {
	init := genutilcli.InitCmd(ctx, cdc, app.ModuleBasics, app.DefaultNodeHome)
	init.Args = cobra.ExactArgs(2)
	init.PostRunE = func(cmd *cobra.Command, args []string) error {
		config := ctx.Config
		config.SetRoot(viper.GetString(cli.HomeFlag))
		genFile := config.GenesisFile()
		genDoc := &types.GenesisDoc{}
		addr, e := sdk.AccAddressFromBech32(args[1])
		if e != nil {
			panic(e)
		}

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
		// migrate gov state
		if appState[gov.ModuleName] != nil {
			var govGenState gov.GenesisState
			cdc.MustUnmarshalJSON(appState[gov.ModuleName], &govGenState)
			minDeposit := sdk.NewInt64Coin(truchain.StakeDenom, 10_000_000)
			govGenState.DepositParams.MinDeposit = sdk.NewCoins(minDeposit)
			appState[gov.ModuleName] = cdc.MustMarshalJSON(govGenState)
		}
		// migrate mint state
		if appState[mint.ModuleName] != nil {
			var mintGenState mint.GenesisState
			cdc.MustUnmarshalJSON(appState[mint.ModuleName], &mintGenState)
			mintGenState.Params.MintDenom = truchain.StakeDenom
			appState[mint.ModuleName] = cdc.MustMarshalJSON(mintGenState)
		}
		// migrate crisis state
		if appState[crisis.ModuleName] != nil {
			var crisisGenState crisis.GenesisState
			cdc.MustUnmarshalJSON(appState[crisis.ModuleName], &crisisGenState)
			crisisGenState.ConstantFee.Denom = truchain.StakeDenom
			appState[crisis.ModuleName] = cdc.MustMarshalJSON(crisisGenState)
		}
		// migrate account state
		if appState[account.ModuleName] != nil {
			var accountGenState account.GenesisState
			cdc.MustUnmarshalJSON(appState[account.ModuleName], &accountGenState)
			accountGenState.Params.Registrar = addr
			appState[account.ModuleName] = cdc.MustMarshalJSON(accountGenState)
		}
		// migrate community state
		if appState[community.ModuleName] != nil {
			var communityGenState community.GenesisState
			cdc.MustUnmarshalJSON(appState[community.ModuleName], &communityGenState)
			communityGenState.Params.CommunityAdmins = []sdk.AccAddress{addr}
			appState[community.ModuleName] = cdc.MustMarshalJSON(communityGenState)
		}
		// migrate claim state
		if appState[claim.ModuleName] != nil {
			var genState claim.GenesisState
			cdc.MustUnmarshalJSON(appState[claim.ModuleName], &genState)
			genState.Params.ClaimAdmins = []sdk.AccAddress{addr}
			appState[claim.ModuleName] = cdc.MustMarshalJSON(genState)
		}
		// migrate staking state
		if appState[trustaking.ModuleName] != nil {
			var genState trustaking.GenesisState
			cdc.MustUnmarshalJSON(appState[trustaking.ModuleName], &genState)
			genState.Params.StakingAdmins = []sdk.AccAddress{addr}
			appState[trustaking.ModuleName] = cdc.MustMarshalJSON(genState)
		}
		// migrate slashing state
		if appState[truslashing.ModuleName] != nil {
			var genState truslashing.GenesisState
			cdc.MustUnmarshalJSON(appState[truslashing.ModuleName], &genState)
			genState.Params.SlashAdmins = []sdk.AccAddress{addr}
			appState[truslashing.ModuleName] = cdc.MustMarshalJSON(genState)
		}
		// migrate trubank state
		if appState[trubank.ModuleName] != nil {
			var genState trubank.GenesisState
			cdc.MustUnmarshalJSON(appState[trubank.ModuleName], &genState)
			genState.Params.RewardBrokerAddress = addr
			appState[trubank.ModuleName] = cdc.MustMarshalJSON(genState)
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
