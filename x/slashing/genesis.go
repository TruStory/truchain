package slashing

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

// GenesisState defines genesis data for the module
type GenesisState struct {
	Slashes      []Slash         `json:"slashes"`
	Params       Params          `json:"params"`
	AdminPubKeys []crypto.PubKey `json:"admin_public_keys"`
}

// NewGenesisState creates a new genesis state.
func NewGenesisState() GenesisState {
	// make a new random account for admin
	admin := secp256k1.GenPrivKey()
	admin2 := secp256k1.GenPrivKey()

	return GenesisState{
		Slashes:      []Slash{},
		Params:       DefaultParams(),
		AdminPubKeys: []crypto.PubKey{admin.PubKey(), admin2.PubKey()},
	}
}

// DefaultGenesisState returns a default genesis state
func DefaultGenesisState() GenesisState { return NewGenesisState() }

// InitGenesis initializes story state from genesis file
func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) {
	for _, slash := range data.Slashes {
		keeper.setSlash(ctx, slash)
		keeper.setCreatorSlash(ctx, slash.Creator, slash.ID)
		keeper.setStakeSlash(ctx, slash.StakeID, slash.ID)
	}
	keeper.setSlashID(ctx, uint64(len(data.Slashes)+1))

	for _, admin := range data.AdminPubKeys {
		data.Params.SlashAdmins = append(data.Params.SlashAdmins, sdk.AccAddress(admin.Address()))
	}
	keeper.SetParams(ctx, data.Params)
}

// ExportGenesis exports the genesis state
func ExportGenesis(ctx sdk.Context, keeper Keeper) GenesisState {
	return GenesisState{
		Slashes: keeper.Slashes(ctx),
		Params:  keeper.GetParams(ctx),
	}
}

// ValidateGenesis validates the genesis state data
func ValidateGenesis(data GenesisState) error {
	if data.Params.MaxStakeSlashCount < 1 {
		return fmt.Errorf("Param: MaxStakeSlashCount, must have a positive value")
	}

	if data.Params.SlashMagnitude < 1 {
		return fmt.Errorf("Param: SlashMagnitude, must have a positive value")
	}

	if data.Params.SlashMinStake.IsNegative() {
		return fmt.Errorf("Param: SlashMinStake, cannot be a negative value")
	}

	if len(data.Params.SlashAdmins) < 1 {
		return fmt.Errorf("Param: SlashAdmins, must have atleast one admin")
	}

	if data.Params.JailTime.Seconds() < 1 {
		return fmt.Errorf("Param: JailTime, must have a positive value")
	}

	if data.Params.CuratorShare.IsNegative() {
		return fmt.Errorf("Param: CuratorShare, cannot be a negative value")
	}

	return nil
}
