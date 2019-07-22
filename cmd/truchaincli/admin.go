package main

import (
	"fmt"

	"github.com/TruStory/truchain/x/staking"

	"github.com/TruStory/truchain/x/claim"

	"github.com/TruStory/truchain/x/community"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/client/utils"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtxb "github.com/cosmos/cosmos-sdk/x/auth/client/txbuilder"

	"github.com/spf13/cobra"
)

// GetAdminCmd will create a new community
func GetAdminCmd(cdc *codec.Codec) *cobra.Command {
	adminCmd := &cobra.Command{
		Use:                        "admin",
		Short:                      "Add or remove an admin to a module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
	}

	adminCmd.PersistentFlags().String("action", "", "Choose either 'add' or 'remove'.")
	adminCmd.PersistentFlags().String("auth", "", "The cosmos address that is authorised to perform this action.")
	err := adminCmd.MarkPersistentFlagRequired("action")
	if err != nil {
		panic(err)
	}
	err = adminCmd.MarkPersistentFlagRequired("auth")
	if err != nil {
		panic(err)
	}

	adminCmd.AddCommand(CommunityAdminCmd(cdc))
	adminCmd.AddCommand(ClaimAdminCmd(cdc))
	adminCmd.AddCommand(StakingAdminCmd(cdc))
	adminCmd.AddCommand(SlashingAdminCmd(cdc))

	return adminCmd
}

// CommunityAdminCmd commands exposes the commands to interact with community admins
func CommunityAdminCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "community [admin]",
		Short: "Add/remove an admin to/from the community module",
		Args:  cobra.ExactArgs(1),

		RunE: func(cmd *cobra.Command, args []string) error {
			auth := cmd.Flag("auth").Value.String()
			action := cmd.Flag("action").Value.String()
			// build and sign the transaction, then broadcast to Tendermint
			authAdmin, err := sdk.AccAddressFromBech32(auth)
			if err != nil {
				panic(err)
			}
			newAdmin, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				panic(err)
			}
			var msg sdk.Msg
			if action == "add" {
				msg = community.NewMsgAddAdmin(newAdmin, authAdmin)
			} else if action == "remove" {
				msg = community.NewMsgRemoveAdmin(newAdmin, authAdmin)
			}

			return executeMsg(cmd, args, cdc, msg)
		},
	}

	cmd = client.PostCommands(cmd)[0]

	return cmd
}

// ClaimAdminCmd commands exposes the commands to interact with claim admins
func ClaimAdminCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "claim [admin]",
		Short: "Add/remove an admin to/from the claim module",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			auth := cmd.Flag("auth").Value.String()
			action := cmd.Flag("action").Value.String()
			// build and sign the transaction, then broadcast to Tendermint
			authAdmin, err := sdk.AccAddressFromBech32(auth)
			if err != nil {
				panic(err)
			}
			newAdmin, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				panic(err)
			}
			var msg sdk.Msg
			if action == "add" {
				msg = claim.NewMsgAddAdmin(newAdmin, authAdmin)
			} else if action == "remove" {
				msg = claim.NewMsgRemoveAdmin(newAdmin, authAdmin)
			}

			return executeMsg(cmd, args, cdc, msg)
		},
	}

	cmd = client.PostCommands(cmd)[0]

	return cmd
}

// StakingAdminCmd commands exposes the commands to interact with staking admins
func StakingAdminCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "staking [admin]",
		Short: "Add/remove an admin to/from the staking module",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			auth := cmd.Flag("auth").Value.String()
			action := cmd.Flag("action").Value.String()
			// build and sign the transaction, then broadcast to Tendermint
			authAdmin, err := sdk.AccAddressFromBech32(auth)
			if err != nil {
				panic(err)
			}
			newAdmin, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				panic(err)
			}
			var msg sdk.Msg
			if action == "add" {
				msg = staking.NewMsgAddAdmin(newAdmin, authAdmin)
			} else if action == "remove" {
				msg = staking.NewMsgRemoveAdmin(newAdmin, authAdmin)
			}

			return executeMsg(cmd, args, cdc, msg)
		},
	}

	cmd = client.PostCommands(cmd)[0]

	return cmd
}

// SlashingAdminCmd commands exposes the commands to interact with slashing admins
func SlashingAdminCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "slashing [admin]",
		Short: "Add/remove an admin to/from the slashing module",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			auth := cmd.Flag("auth").Value.String()
			action := cmd.Flag("action").Value.String()
			// build and sign the transaction, then broadcast to Tendermint
			authAdmin, err := sdk.AccAddressFromBech32(auth)
			if err != nil {
				panic(err)
			}
			newAdmin, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				panic(err)
			}
			var msg sdk.Msg
			if action == "add" {
				msg = claim.NewMsgAddAdmin(newAdmin, authAdmin)
			} else if action == "remove" {
				msg = claim.NewMsgRemoveAdmin(newAdmin, authAdmin)
			}

			return executeMsg(cmd, args, cdc, msg)
		},
	}

	cmd = client.PostCommands(cmd)[0]

	return cmd
}

func executeMsg(cmd *cobra.Command, args []string, cdc *codec.Codec, msg sdk.Msg) error {
	auth := cmd.Flag("auth").Value.String()
	cliCtx := context.NewCLIContextWithFrom(auth).
		WithCodec(cdc).
		WithAccountDecoder(cdc)
	seq, _ := cliCtx.GetAccountSequence(cliCtx.FromAddress)
	txBldr := authtxb.NewTxBuilderFromCLI().WithSequence(seq).WithTxEncoder(utils.GetTxEncoder(cliCtx.Codec))

	// build and sign the transaction, then broadcast to Tendermint
	fromName := cliCtx.GetFromName()
	passphrase, err := keys.GetPassphrase(fromName)
	if err != nil {
		return err
	}

	txBytes, err := txBldr.BuildAndSign(fromName, passphrase, []sdk.Msg{msg})
	if err != nil {
		return err
	}

	// broadcast to a Tendermint node
	res, err := cliCtx.WithBroadcastMode(client.BroadcastBlock).BroadcastTx(txBytes)
	if err != nil {
		return err
	}
	fmt.Println(res)
	return nil
}
