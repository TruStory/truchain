package main

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/x/auth"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/spf13/cobra"

	"github.com/TruStory/truchain/x/bank"
)

// SendGiftCmd will create a send tx and sign it with the given key.
func SendGiftCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "send_gift [from_key_or_address] [to_address] [amount]",
		Short: "Create and sign a send gift tx",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(auth.DefaultTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			to, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				return err
			}

			// parse coin trying to be sent
			coin, err := sdk.ParseCoin(args[2])
			if err != nil {
				return err
			}
			// build and sign the transaction, then broadcast to Tendermint
			msg := bank.NewMsgSendGift(cliCtx.GetFromAddress(), to, coin)
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
