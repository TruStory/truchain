package init

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	cfg "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/libs/cli"
	"github.com/tendermint/tendermint/libs/common"
	"github.com/tendermint/tendermint/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/cmd/gaia/app"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
)

// NOTE: Much of this file is copied from it's gaia equivalent. It initializes a genesis file with defaults.
// This should go away after it becomes part of a module in the Cosmos SDK (in the works).

const (
	flagOverwrite = "overwrite"
	flagMoniker   = "moniker"
)

type printInfo struct {
	Moniker    string          `json:"moniker"`
	ChainID    string          `json:"chain_id"`
	NodeID     string          `json:"node_id"`
	GenTxsDir  string          `json:"gentxs_dir"`
	AppMessage json.RawMessage `json:"app_message"`
}

// nolint: errcheck
func DisplayInfo(cdc *codec.Codec, info printInfo) error {
	out, err := codec.MarshalJSONIndent(cdc, info)
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "%s\n", string(out))
	return nil
}

// InitCmd initializes all files for truchain and application
func InitCmd(ctx *server.Context, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize private validator, p2p, genesis, and application configuration files",
		Long:  `Initialize validators's and node's configuration files.`,
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			config := ctx.Config
			config.SetRoot(viper.GetString(cli.HomeFlag))

			chainID := viper.GetString(client.FlagChainID)
			if chainID == "" {
				chainID = fmt.Sprintf("test-chain-%v", common.RandStr(6))
			}

			nodeID, _, err := InitializeNodeValidatorFiles(config)
			if err != nil {
				return err
			}

			config.Moniker = viper.GetString(flagMoniker)

			var appState json.RawMessage
			genFile := config.GenesisFile()

			if appState, err = initializeEmptyGenesis(cdc, genFile, chainID,
				viper.GetBool(flagOverwrite)); err != nil {
				return err
			}

			genDoc := &types.GenesisDoc{}
			if _, err := os.Stat(genFile); err != nil {
				if !os.IsNotExist(err) {
					return err
				}
			} else {
				genDoc, err = types.GenesisDocFromFile(genFile)
				if err != nil {
					return err
				}
			}

			genDoc.ChainID = chainID
			genDoc.Validators = nil
			genDoc.AppState = appState
			if err = ExportGenesisFile(genDoc, genFile); err != nil {
				return err
			}

			toPrint := NewPrintInfo(config.Moniker, chainID, nodeID, "", appState)

			cfg.WriteConfigFile(filepath.Join(config.RootDir, "config", "config.toml"), config)

			return DisplayInfo(cdc, toPrint)
		},
	}

	cmd.Flags().String(cli.HomeFlag, app.DefaultNodeHome, "node's home directory")
	cmd.Flags().BoolP(flagOverwrite, "o", false, "overwrite the genesis.json file")
	cmd.Flags().String(client.FlagChainID, "", "genesis file chain-id, if left blank will be randomly created")
	cmd.Flags().String(flagMoniker, "", "set the validator's moniker")
	err := cmd.MarkFlagRequired(flagMoniker)
	if err != nil {
		panic(err)
	}

	return cmd
}
