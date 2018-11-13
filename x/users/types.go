package users

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	tcmn "github.com/tendermint/tendermint/libs/common"
)

// User is the externally-facing account object
type User struct {
	Address       string        `json:"address"`
	AccountNumber int64         `json:"account_number"`
	Coins         sdk.Coins     `json:"coins"`
	Sequence      int64         `json:"sequence"`
	Pubkey        tcmn.HexBytes `json:"pubkey"`
}

// NewUser creates a new User struct from an auth.Account (like AppAccount)
func NewUser(acc auth.Account) User {
	return User{
		Address:       acc.GetAddress().String(),
		AccountNumber: acc.GetAccountNumber(),
		Coins:         acc.GetCoins(),
		Sequence:      acc.GetSequence(),
		Pubkey:        tcmn.HexBytes(acc.GetPubKey().Bytes()),
	}
}

// TwitterProfile is the Twitter profile associated with the account with address `Address`
type TwitterProfile struct {
	ID       string `json:"id"`
	Address  string `json:"address"`
	Username string `json:"username"`
	FullName string `json:"full_name"`
}
