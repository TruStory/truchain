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

// NewCommunityCmd will create a new community
func NewCommunityCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "new-community [id] [name] [description] [creator]",
		Short: "Create a new community with the given details",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {

			cliCtx := context.NewCLIContextWithFrom(args[3]).
				WithCodec(cdc).
				WithAccountDecoder(cdc)
			seq, _ := cliCtx.GetAccountSequence(cliCtx.FromAddress)
			txBldr := authtxb.NewTxBuilderFromCLI().WithSequence(seq).WithTxEncoder(utils.GetTxEncoder(cliCtx.Codec))

			// build and sign the transaction, then broadcast to Tendermint
			msg := community.NewMsgNewCommunity(args[0], args[1], args[2], cliCtx.GetFromAddress())
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
