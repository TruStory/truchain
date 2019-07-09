package slashing

import (
	"fmt"

	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/secp256k1"

	app "github.com/TruStory/truchain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
)

// Keys for params
var (
	KeyMinSlashCount  = []byte("minStakeSlashCount")
	KeySlashMagnitude = []byte("slashMagnitude")
	KeySlashMinStake  = []byte("slashMinStake")
	KeySlashAdmins    = []byte("slashAdmins")
	KeyCuratorShare   = []byte("curatorShare")
)

// Params holds parameters for Slashing
type Params struct {
	MinSlashCount  int              `json:"min_slash_count"`
	SlashMagnitude int              `json:"slash_magnitude"`
	SlashMinStake  sdk.Coin         `json:"slash_min_stake"`
	SlashAdmins    []sdk.AccAddress `json:"slash_admins"`
	CuratorShare   sdk.Dec          `json:"curator_share"`
}

// DefaultParams is the Slashing params for testing
func DefaultParams() Params {
	_, _, adminAddr1 := getFakeKeyPubAddr()
	_, _, adminAddr2 := getFakeKeyPubAddr()
	return Params{
		MinSlashCount:  5,
		SlashMagnitude: 3,
		SlashMinStake:  sdk.NewCoin(app.StakeDenom, sdk.NewInt(10*app.Shanev)),
		SlashAdmins:    []sdk.AccAddress{adminAddr1, adminAddr2},
		CuratorShare:   sdk.NewDecWithPrec(25, 2),
	}
}

// ParamSetPairs implements params.ParamSet
func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		{Key: KeyMinSlashCount, Value: &p.MinSlashCount},
		{Key: KeySlashMagnitude, Value: &p.SlashMagnitude},
		{Key: KeySlashMinStake, Value: &p.SlashMinStake},
		{Key: KeySlashAdmins, Value: &p.SlashAdmins},
		{Key: KeyCuratorShare, Value: &p.CuratorShare},
	}
}

// ParamKeyTable for slashing module
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&Params{})
}

// GetParams gets the genesis params for the slashing
func (k Keeper) GetParams(ctx sdk.Context) Params {
	var paramSet Params
	k.paramStore.GetParamSet(ctx, &paramSet)
	return paramSet
}

// SetParams sets the params for the slashing
func (k Keeper) SetParams(ctx sdk.Context, params Params) {
	logger := ctx.Logger().With("module", ModuleName)
	k.paramStore.SetParamSet(ctx, &params)
	logger.Info(fmt.Sprintf("Loaded slashing params: %+v", params))
}

// unexported and used for testing...
func getFakeKeyPubAddr() (crypto.PrivKey, crypto.PubKey, sdk.AccAddress) {
	key := secp256k1.GenPrivKey()
	pub := key.PubKey()
	addr := sdk.AccAddress(pub.Address())
	return key, pub, addr
}
