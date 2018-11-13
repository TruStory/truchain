package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/tendermint/tendermint/crypto"
)

var _ auth.Account = (*AppAccount)(nil)

// RegistrarAccAddress is the UTF8-encoded address of the account which is responsible for registering other accounts
const RegistrarAccAddress = "truchainaccregistrar"

// AppAccount is a custom extension for this application. It is an example of
// extending auth.BaseAccount with custom fields. It is compatible with the
// stock auth.AccountStore, since auth.AccountStore uses the flexible go-amino
// library.
type AppAccount struct {
	auth.BaseAccount
}

// GetAccountNumber returns the account number
func (acc AppAccount) GetAccountNumber() int64 {
	return acc.AccountNumber
}

// GetCoins returns the coins
func (acc AppAccount) GetCoins() sdk.Coins {
	return acc.Coins
}

// GetSequence returns the sequence
func (acc AppAccount) GetSequence() int64 {
	return acc.Sequence
}

// SetAccountNumber sets the account number
func (acc AppAccount) SetAccountNumber(accNumber int64) error {
	return acc.BaseAccount.SetAccountNumber(accNumber)
}

// SetAddress sets the address
func (acc AppAccount) SetAddress(address sdk.AccAddress) error {
	return acc.BaseAccount.SetAddress(address)
}

// SetSequence sets the sequence
func (acc AppAccount) SetSequence(seq int64) error {
	return acc.BaseAccount.SetSequence(seq)
}

// SetCoins sets the coins
func (acc AppAccount) SetCoins(coins sdk.Coins) error {
	return acc.BaseAccount.SetCoins(coins)
}

// SetPubKey sets the public key
func (acc AppAccount) SetPubKey(pubkey crypto.PubKey) error {
	return acc.BaseAccount.SetPubKey(pubkey)
}

// NewAppAccount returns a reference to a new AppAccount given a name and an
// auth.BaseAccount.
func NewAppAccount(baseAcct auth.BaseAccount) *AppAccount {
	return &AppAccount{BaseAccount: baseAcct}
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
	Address sdk.AccAddress `json:"address"`
	Coins   sdk.Coins      `json:"coins"`
}

// NewGenesisAccount returns a reference to a new GenesisAccount given an
// AppAccount.
func NewGenesisAccount(aa *AppAccount) *GenesisAccount {
	return &GenesisAccount{
		Address: aa.Address,
		Coins:   aa.Coins.Sort(),
	}
}

// ToAppAccount converts a GenesisAccount to an AppAccount.
func (ga *GenesisAccount) ToAppAccount() (acc *AppAccount, err error) {
	return &AppAccount{
		BaseAccount: auth.BaseAccount{
			Address: ga.Address,
			Coins:   ga.Coins.Sort(),
		},
	}, nil
}
