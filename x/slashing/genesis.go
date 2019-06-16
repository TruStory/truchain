package slashing

import (
	"fmt"

	"github.com/TruStory/truchain/x/auth"
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

	return GenesisState{
		Slashes:      []Slash{},
		Params:       DefaultParams(),
		AdminPubKeys: []crypto.PubKey{admin.PubKey()},
	}
}

// DefaultGenesisState returns a default genesis state
func DefaultGenesisState() GenesisState { return NewGenesisState() }

// InitGenesis initializes story state from genesis file
func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) {
	for _, slash := range data.Slashes {
		keeper.Set(ctx, slash.ID, slash)
	}

	for _, admin := range data.AdminPubKeys {
		appAccount := keeper.appAccountKeeper.NewAppAccount(ctx, sdk.AccAddress(admin.Address()), sdk.Coins{}, admin, 0, 0)
		err := keeper.appAccountKeeper.AddToEarnedStake(ctx, appAccount.Address, auth.EarnedCoin{Coin: sdk.NewCoin("default", sdk.NewInt(1)), CommunityID: 1})
		if err != nil {
			panic(err)
		}
		data.Params.SlashAdmins = append(data.Params.SlashAdmins, appAccount.GetAddress())
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

	if !data.Params.SlashMagnitude.IsPositive() {
		return fmt.Errorf("Param: SlashMagnitude, must have a positive value")
	}

	for _, coin := range data.Params.SlashMinStake {
		if coin.IsNegative() {
			return fmt.Errorf("Param: SlashMinStake, cannot be a negative value")
		}
	}

	if len(data.Params.SlashAdmins) < 1 {
		return fmt.Errorf("Param: SlashAdmins, must have atleast one admin")
	}

	if data.Params.JailTime.Seconds() < 1 {
		return fmt.Errorf("Param: JailTime, must have a positive value")
	}

	return nil
}
