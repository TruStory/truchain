package main

import (
	"fmt"
	"io/ioutil"
	"time"

	"github.com/TruStory/truchain/cmd/truchaind/migration/v0_3"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/types"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	extypes "github.com/cosmos/cosmos-sdk/x/genutil"
)

const (
	flagGenesisTime     = "genesis-time"
	flagChainID         = "chain-id"
	flagResetValidators = "reset-validators"
)

var migrationMap = extypes.MigrationMap{
	"v0.3.1": v0_3.Migrate,
}

// GetMigrationCallback returns a MigrationCallback for a given version.
func GetMigrationCallback(version string) extypes.MigrationCallback {
	return migrationMap[version]
}

// MigrateGenesisCmd returns a command to execute genesis state migration.
// nolint: funlen
func MigrateGenesisCmd(_ *server.Context, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tru_migrate [target-version] [genesis-file] [output-file]",
		Short: "Migrate genesis to a specified target version",
		Long: fmt.Sprintf(`Migrate the source genesis into the target version and print to STDOUT.
Example:
$ %s tru_migrate v0.3.1 /path/to/genesis.json /path/to/output-genesis.json --chain-id=betanet-2 --genesis-time=2019-04-22T17:00:00Z
`, version.ServerName),
		Args: cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error

			target := args[0]
			importGenesis := args[1]
			outputGenesis := args[2]
			if outputGenesis == "" {
				return errors.New("must provide a valid path for output file")
			}

			genDoc, err := types.GenesisDocFromFile(importGenesis)
			if err != nil {
				return errors.Wrapf(err, "failed to read genesis document from file %s", importGenesis)
			}

			var initialState extypes.AppMap
			if err := cdc.UnmarshalJSON(genDoc.AppState, &initialState); err != nil {
				return errors.Wrap(err, "failed to JSON unmarshal initial genesis state")
			}

			migrationFunc := GetMigrationCallback(target)
			if migrationFunc == nil {
				return fmt.Errorf("unknown migration function for version: %s", target)
			}

			// TODO: handler error from migrationFunc call
			newGenState := migrationFunc(initialState)
			resetValidators := viper.GetBool(flagResetValidators)
			if resetValidators {
				genDoc.Validators = nil
			}
			genDoc.AppState, err = cdc.MarshalJSON(newGenState)
			if err != nil {
				return errors.Wrap(err, "failed to JSON marshal migrated genesis state")
			}

			genesisTime := cmd.Flag(flagGenesisTime).Value.String()
			if genesisTime != "" {
				var t time.Time

				err := t.UnmarshalText([]byte(genesisTime))
				if err != nil {
					return errors.Wrap(err, "failed to unmarshal genesis time")
				}

				genDoc.GenesisTime = t
			}

			chainID := cmd.Flag(flagChainID).Value.String()
			if chainID != "" {
				genDoc.ChainID = chainID
			}

			bz, err := cdc.MarshalJSONIndent(genDoc, "", "\t")
			if err != nil {
				return errors.Wrap(err, "failed to marshal genesis doc")
			}

			sortedBz, err := sdk.SortJSON(bz)
			if err != nil {
				return errors.Wrap(err, "failed to sort JSON genesis doc")
			}
			err = ioutil.WriteFile(outputGenesis, sortedBz, 0644)
			return err
		},
	}

	cmd.Flags().String(flagGenesisTime, "", "override genesis_time with this flag")
	cmd.Flags().String(flagChainID, "", "override chain_id with this flag")
	cmd.Flags().Bool(flagResetValidators, false, "remove validators set with this flag")
	return cmd
}
