package main

import (
	"reflect"
	"strings"

	"github.com/mitchellh/mapstructure"

	"github.com/TruStory/truchain/x/community"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/spf13/cobra"
)

// GetParamsCmd will create a new community
func GetParamsCmd(cdc *codec.Codec) *cobra.Command {
	paramsCmd := &cobra.Command{
		Use:                        "params",
		Short:                      "Update the params of various modules",
		SuggestionsMinimumDistance: 2,
	}

	paramsCmd.AddCommand(CommunityParamsCmd(cdc))

	return paramsCmd
}

// CommunityParamsCmd commands exposes the commands to interact with community params
func CommunityParamsCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "community",
		Short: "Update the params of the community module",
		Args:  cobra.ExactArgs(0),

		RunE: func(cmd *cobra.Command, args []string) error {
			var params community.Params
			mapInput := make(map[string]string)
			mapParams(params, func(param string, _ int) {
				input := cmd.Flag(param).Value.String()
				if input != "" {
					mapInput[param] = input
				}
			})
			msConfig := &mapstructure.DecoderConfig{
				TagName:          "json",
				WeaklyTypedInput: true,
				Result:           &params,
			}
			decoder, err := mapstructure.NewDecoder(msConfig)
			if err != nil {
				panic(err)
			}
			err = decoder.Decode(mapInput)
			if err != nil {
				panic(err)
			}

			return nil
		},
	}

	// Adding the available flags
	params := community.Params{}
	mapParams(params, func(param string, index int) {
		cmd.Flags().String(param, "", "Updates the param: "+param)
	})

	cmd = client.PostCommands(cmd)[0]

	return cmd
}

// mapParams walks over each param, and ignores the *_admins param because they are out of scope for this CLI command
func mapParams(params interface{}, fn func(param string, index int)) {
	rParams := reflect.TypeOf(params)
	for i := 0; i < rParams.NumField(); i++ {
		param := rParams.Field(i).Tag.Get("json")
		if !strings.HasSuffix(param, "_admins") {
			fn(param, i)
		}
	}
}
