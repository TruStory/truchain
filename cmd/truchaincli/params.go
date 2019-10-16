package main

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/TruStory/truchain/x/account"
	"github.com/TruStory/truchain/x/bank"
	"github.com/TruStory/truchain/x/claim"
	"github.com/TruStory/truchain/x/community"
	"github.com/TruStory/truchain/x/slashing"
	"github.com/TruStory/truchain/x/staking"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cobra"
)

// GetParamsCmd will create a new community
func GetParamsCmd(cdc *codec.Codec) *cobra.Command {
	paramsCmd := &cobra.Command{
		Use:                        "params",
		Short:                      "Update the params of various modules",
		SuggestionsMinimumDistance: 2,
	}

	paramsCmd.AddCommand(AccountParamsCmd(cdc))
	paramsCmd.AddCommand(BankParamsCmd(cdc))
	paramsCmd.AddCommand(CommunityParamsCmd(cdc))
	paramsCmd.AddCommand(ClaimParamsCmd(cdc))
	paramsCmd.AddCommand(StakingParamsCmd(cdc))
	paramsCmd.AddCommand(SlashingParamsCmd(cdc))

	return paramsCmd
}

// AccountParamsCmd commands exposes the commands to interact with account params
func AccountParamsCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "account [auth-address]",
		Short: "Update the params of the account module",
		Args:  cobra.ExactArgs(1),

		RunE: func(cmd *cobra.Command, args []string) error {
			mapInput := make(map[string]string)
			updates := account.Params{}
			updatedFields := make([]string, 0)
			mapParams(updates, func(param string, _ int, field reflect.StructField) {
				input := cmd.Flag(param).Value.String()
				if input != "" {
					if field.Type.PkgPath() == "github.com/cosmos/cosmos-sdk/types" {
						// if cosmos type, we'll make the cosmos object
						reflect.ValueOf(&updates).Elem().FieldByName(field.Name).Set(
							makeCosmosObject(field.Type.String(), cmd.Flag(param).Value.String()),
						)
					} else {
						mapInput[param] = input
					}

					updatedFields = append(updatedFields, param)
				}
			})

			msConfig := &mapstructure.DecoderConfig{
				TagName:          "json",
				WeaklyTypedInput: true,
				Result:           &updates,
			}
			decoder, err := mapstructure.NewDecoder(msConfig)
			if err != nil {
				panic(err)
			}
			err = decoder.Decode(mapInput)
			if err != nil {
				panic(err)
			}

			cliCtx := context.NewCLIContextWithFrom(args[0]).WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(auth.DefaultTxEncoder(cdc))

			msg := account.NewMsgUpdateParams(updates, updatedFields, cliCtx.GetFromAddress())
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

	// Adding the available flags
	mapParams(account.Params{}, func(param string, index int, field reflect.StructField) {
		cmd.Flags().String(param, "", "Updates the param: "+param)
	})

	cmd = client.PostCommands(cmd)[0]

	return cmd
}

// BankParamsCmd commands exposes the commands to interact with bank params
func BankParamsCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bank [auth-address]",
		Short: "Update the params of the bank module",
		Args:  cobra.ExactArgs(1),

		RunE: func(cmd *cobra.Command, args []string) error {
			mapInput := make(map[string]string)
			updates := bank.Params{}
			updatedFields := make([]string, 0)
			mapParams(updates, func(param string, _ int, field reflect.StructField) {
				input := cmd.Flag(param).Value.String()
				if input != "" {
					if field.Type.PkgPath() == "github.com/cosmos/cosmos-sdk/types" {
						// if cosmos type, we'll make the cosmos object
						reflect.ValueOf(&updates).Elem().FieldByName(field.Name).Set(
							makeCosmosObject(field.Type.String(), cmd.Flag(param).Value.String()),
						)
					} else {
						mapInput[param] = input
					}

					updatedFields = append(updatedFields, param)
				}
			})

			msConfig := &mapstructure.DecoderConfig{
				TagName:          "json",
				WeaklyTypedInput: true,
				Result:           &updates,
			}
			decoder, err := mapstructure.NewDecoder(msConfig)
			if err != nil {
				panic(err)
			}
			err = decoder.Decode(mapInput)
			if err != nil {
				panic(err)
			}

			cliCtx := context.NewCLIContextWithFrom(args[0]).WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(auth.DefaultTxEncoder(cdc))

			msg := bank.NewMsgUpdateParams(updates, updatedFields, cliCtx.GetFromAddress())
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

	// Adding the available flags
	mapParams(bank.Params{}, func(param string, index int, field reflect.StructField) {
		cmd.Flags().String(param, "", "Updates the param: "+param)
	})

	cmd = client.PostCommands(cmd)[0]

	return cmd
}

// ClaimParamsCmd commands exposes the commands to interact with claim params
func ClaimParamsCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "claim [auth-address]",
		Short: "Update the params of the claim module",
		Args:  cobra.ExactArgs(1),

		RunE: func(cmd *cobra.Command, args []string) error {
			mapInput := make(map[string]string)
			updates := claim.Params{}
			updatedFields := make([]string, 0)
			mapParams(updates, func(param string, _ int, field reflect.StructField) {
				input := cmd.Flag(param).Value.String()
				if input != "" {
					if field.Type.PkgPath() == "github.com/cosmos/cosmos-sdk/types" {
						// if cosmos type, we'll make the cosmos object
						reflect.ValueOf(&updates).Elem().FieldByName(field.Name).Set(
							makeCosmosObject(field.Type.String(), cmd.Flag(param).Value.String()),
						)
					} else {
						mapInput[param] = input
					}

					updatedFields = append(updatedFields, param)
				}
			})

			msConfig := &mapstructure.DecoderConfig{
				TagName:          "json",
				WeaklyTypedInput: true,
				Result:           &updates,
			}
			decoder, err := mapstructure.NewDecoder(msConfig)
			if err != nil {
				panic(err)
			}
			err = decoder.Decode(mapInput)
			if err != nil {
				panic(err)
			}

			cliCtx := context.NewCLIContextWithFrom(args[0]).WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(auth.DefaultTxEncoder(cdc))

			msg := claim.NewMsgUpdateParams(updates, updatedFields, cliCtx.GetFromAddress())
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

	// Adding the available flags
	mapParams(claim.Params{}, func(param string, index int, field reflect.StructField) {
		cmd.Flags().String(param, "", "Updates the param: "+param)
	})

	cmd = client.PostCommands(cmd)[0]

	return cmd
}

// CommunityParamsCmd commands exposes the commands to interact with community params
func CommunityParamsCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "community [auth-address]",
		Short: "Update the params of the community module",
		Args:  cobra.ExactArgs(1),

		RunE: func(cmd *cobra.Command, args []string) error {
			mapInput := make(map[string]string)
			updates := community.Params{}
			updatedFields := make([]string, 0)
			mapParams(updates, func(param string, _ int, field reflect.StructField) {
				input := cmd.Flag(param).Value.String()
				if input != "" {
					if field.Type.PkgPath() == "github.com/cosmos/cosmos-sdk/types" {
						// if cosmos type, we'll make the cosmos object
						reflect.ValueOf(&updates).Elem().FieldByName(field.Name).Set(
							makeCosmosObject(field.Type.String(), cmd.Flag(param).Value.String()),
						)
					} else {
						mapInput[param] = input
					}

					updatedFields = append(updatedFields, param)
				}
			})

			msConfig := &mapstructure.DecoderConfig{
				TagName:          "json",
				WeaklyTypedInput: true,
				Result:           &updates,
			}
			decoder, err := mapstructure.NewDecoder(msConfig)
			if err != nil {
				panic(err)
			}
			err = decoder.Decode(mapInput)
			if err != nil {
				panic(err)
			}

			cliCtx := context.NewCLIContextWithFrom(args[0]).WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(auth.DefaultTxEncoder(cdc))

			msg := community.NewMsgUpdateParams(updates, updatedFields, cliCtx.GetFromAddress())
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

	// Adding the available flags
	mapParams(community.Params{}, func(param string, index int, field reflect.StructField) {
		cmd.Flags().String(param, "", "Updates the param: "+param)
	})

	cmd = client.PostCommands(cmd)[0]

	return cmd
}

// SlashingParamsCmd commands exposes the commands to interact with slashing params
func SlashingParamsCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "slashing [auth-address]",
		Short: "Update the params of the slashing module",
		Args:  cobra.ExactArgs(1),

		RunE: func(cmd *cobra.Command, args []string) error {
			mapInput := make(map[string]string)
			updates := slashing.Params{}
			updatedFields := make([]string, 0)
			mapParams(updates, func(param string, _ int, field reflect.StructField) {
				input := cmd.Flag(param).Value.String()
				if input != "" {
					if field.Type.PkgPath() == "github.com/cosmos/cosmos-sdk/types" {
						// if cosmos type, we'll make the cosmos object
						reflect.ValueOf(&updates).Elem().FieldByName(field.Name).Set(
							makeCosmosObject(field.Type.String(), cmd.Flag(param).Value.String()),
						)
					} else {
						mapInput[param] = input
					}

					updatedFields = append(updatedFields, param)
				}
			})

			msConfig := &mapstructure.DecoderConfig{
				TagName:          "json",
				WeaklyTypedInput: true,
				Result:           &updates,
			}
			decoder, err := mapstructure.NewDecoder(msConfig)
			if err != nil {
				panic(err)
			}
			err = decoder.Decode(mapInput)
			if err != nil {
				panic(err)
			}

			cliCtx := context.NewCLIContextWithFrom(args[0]).WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(auth.DefaultTxEncoder(cdc))

			msg := slashing.NewMsgUpdateParams(updates, updatedFields, cliCtx.GetFromAddress())
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

	// Adding the available flags
	mapParams(slashing.Params{}, func(param string, index int, field reflect.StructField) {
		cmd.Flags().String(param, "", "Updates the param: "+param)
	})

	cmd = client.PostCommands(cmd)[0]

	return cmd
}

// StakingParamsCmd commands exposes the commands to interact with staking params
func StakingParamsCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "staking [auth-address]",
		Short: "Update the params of the staking module",
		Args:  cobra.ExactArgs(1),

		RunE: func(cmd *cobra.Command, args []string) error {
			mapInput := make(map[string]string)
			updates := staking.Params{}
			updatedFields := make([]string, 0)
			mapParams(updates, func(param string, _ int, field reflect.StructField) {
				input := cmd.Flag(param).Value.String()
				if input != "" {
					if field.Type.PkgPath() == "github.com/cosmos/cosmos-sdk/types" {
						// if cosmos type, we'll make the cosmos object
						reflect.ValueOf(&updates).Elem().FieldByName(field.Name).Set(
							makeCosmosObject(field.Type.String(), cmd.Flag(param).Value.String()),
						)
					} else {
						mapInput[param] = input
					}

					updatedFields = append(updatedFields, param)
				}
			})

			msConfig := &mapstructure.DecoderConfig{
				TagName:          "json",
				WeaklyTypedInput: true,
				Result:           &updates,
			}
			decoder, err := mapstructure.NewDecoder(msConfig)
			if err != nil {
				panic(err)
			}
			err = decoder.Decode(mapInput)
			if err != nil {
				panic(err)
			}

			cliCtx := context.NewCLIContextWithFrom(args[0]).WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(auth.DefaultTxEncoder(cdc))

			msg := staking.NewMsgUpdateParams(updates, updatedFields, cliCtx.GetFromAddress())
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

	// Adding the available flags
	mapParams(staking.Params{}, func(param string, index int, field reflect.StructField) {
		cmd.Flags().String(param, "", "Updates the param: "+param)
	})

	cmd = client.PostCommands(cmd)[0]

	return cmd
}

// mapParams walks over each param, and ignores the *_admins param because they are out of scope for this CLI command
func mapParams(params interface{}, fn func(param string, index int, field reflect.StructField)) {
	rParams := reflect.TypeOf(params)
	for i := 0; i < rParams.NumField(); i++ {
		field := rParams.Field(i)
		param := field.Tag.Get("json")
		if !strings.HasSuffix(param, "_admins") {
			fn(param, i, field)
		}
	}
}

// makeCosmosObject converts the input string into correct cosmos object
func makeCosmosObject(cosmosType string, value string) reflect.Value {
	if cosmosType == "types.Dec" {
		dec := sdk.MustNewDecFromStr(value)
		return reflect.ValueOf(dec)
	}

	if cosmosType == "types.Coin" {
		coin, err := sdk.ParseCoin(value)
		if err != nil {
			panic(err)
		}
		return reflect.ValueOf(coin)
	}

	if cosmosType == "types.AccAddress" {
		address, err := sdk.AccAddressFromBech32(value)
		if err != nil {
			panic(err)
		}
		return reflect.ValueOf(address)
	}

	return reflect.ValueOf(value)
}
