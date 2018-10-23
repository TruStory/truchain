package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/tendermint/tendermint/crypto"
)

var _ auth.Account = (*AppAccount)(nil)

const RegistrarAccAddress = "truchainaccregistrar"

// AppAccount is a custom extension for this application. It is an example of
// extending auth.BaseAccount with custom fields. It is compatible with the
// stock auth.AccountStore, since auth.AccountStore uses the flexible go-amino
// library.
type AppAccount struct {
	auth.BaseAccount

	Name string `json:"name"`
}

// nolint
func (acc AppAccount) GetName() string { return acc.Name }

// nolint
func (acc *AppAccount) SetName(name string) { acc.Name = name }

// nolint
func (acc AppAccount) GetAccountNumber() int64 {
	return acc.AccountNumber
}

// nolint
func (acc AppAccount) GetCoins() sdk.Coins {
	return acc.Coins
}

// nolint
func (acc AppAccount) GetSequence() int64 {
	return acc.Sequence
}

// nolint
func (acc AppAccount) SetAccountNumber(accNumber int64) error {
	acc.BaseAccount.SetAccountNumber(accNumber)
	return nil
}

// nolint
func (acc AppAccount) SetAddress(address sdk.AccAddress) error {
	acc.BaseAccount.SetAddress(address)
	return nil
}

// nolint
func (acc AppAccount) SetSequence(seq int64) error {
	acc.BaseAccount.SetSequence(seq)
	return nil
}

// nolint
func (acc AppAccount) SetCoins(coins sdk.Coins) error {
	acc.BaseAccount.SetCoins(coins)
	return nil
}

// nolint
func (acc AppAccount) SetPubKey(pubkey crypto.PubKey) error {
	acc.PubKey = pubkey
	return nil
}

// NewAppAccount returns a reference to a new AppAccount given a name and an
// auth.BaseAccount.
func NewAppAccount(name string, baseAcct auth.BaseAccount) *AppAccount {
	return &AppAccount{BaseAccount: baseAcct, Name: name}
}

// GetAccountDecoder returns the AccountDecoder function for the custom
// AppAccount.
func GetAccountDecoder(cdc *codec.Codec) auth.AccountDecoder {
	return func(accBytes []byte) (auth.Account, error) {
		if len(accBytes) == 0 {
			return nil, sdk.ErrTxDecode("accBytes are empty")
		}

		acct := new(AppAccount)
		err := cdc.UnmarshalBinaryBare(accBytes, &acct)
		if err != nil {
			panic(err)
		}

		return acct, err
	}
}

// GenesisState reflects the genesis state of the application.
type GenesisState struct {
	Accounts []*GenesisAccount `json:"accounts"`
}

// GenesisAccount reflects a genesis account the application expects in it's
// genesis state.
type GenesisAccount struct {
	Name    string         `json:"name"`
	Address sdk.AccAddress `json:"address"`
	Coins   sdk.Coins      `json:"coins"`
}

// NewGenesisAccount returns a reference to a new GenesisAccount given an
// AppAccount.
func NewGenesisAccount(aa *AppAccount) *GenesisAccount {
	return &GenesisAccount{
		Name:    aa.Name,
		Address: aa.Address,
		Coins:   aa.Coins.Sort(),
	}
}

// ToAppAccount converts a GenesisAccount to an AppAccount.
func (ga *GenesisAccount) ToAppAccount() (acc *AppAccount, err error) {
	return &AppAccount{
		Name: ga.Name,
		BaseAccount: auth.BaseAccount{
			Address: ga.Address,
			Coins:   ga.Coins.Sort(),
		},
	}, nil
}
