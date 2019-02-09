package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"

	sdk "github.com/cosmos/cosmos-sdk/types"
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
	AuthData auth.GenesisState   `json:"auth"`
	BankData bank.GenesisState   `json:"bank"`
	Accounts []*auth.BaseAccount `json:"accounts"`
}
