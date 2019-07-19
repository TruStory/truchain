package main

import (
	"fmt"

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

	// adminCmd = client.PostCommands(
	// 	CommunityAdminCmd(cdc),
	// )[0]

	adminCmd.AddCommand(CommunityAdminCmd(cdc))

	return adminCmd
}

// CommunityAdminCmd commands exposes the commands to interact with community admins
func CommunityAdminCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "community-add [auth-admin] [new-admin]",
		Short: "Add a new admin to the community admins",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {

			cliCtx := context.NewCLIContextWithFrom(args[0]).
				WithCodec(cdc).
				WithAccountDecoder(cdc)
			seq, _ := cliCtx.GetAccountSequence(cliCtx.FromAddress)
			txBldr := authtxb.NewTxBuilderFromCLI().WithSequence(seq).WithTxEncoder(utils.GetTxEncoder(cliCtx.Codec))

			// build and sign the transaction, then broadcast to Tendermint
			newAdmin, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				panic(err)
			}
			msg := community.NewMsgAddAdmin(newAdmin, cliCtx.GetFromAddress())
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
		},
	}

	cmd = client.PostCommands(cmd)[0]

	return cmd
}
