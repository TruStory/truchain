package users

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	tcmn "github.com/tendermint/tendermint/libs/common"
)

// User is the externally-facing account object
type User struct {
	Address       string        `json:"address"`
	AccountNumber uint64        `json:"account_number"`
	Coins         sdk.Coins     `json:"coins"`
	Sequence      uint64        `json:"sequence"`
	Pubkey        tcmn.HexBytes `json:"pubkey"`
}

// NewUser creates a new User struct from an auth.Account (like AppAccount)
func NewUser(acc auth.Account) User {
	var pubKey []byte

	// GetPubKey can return nil and Bytes() will panic due to nil pointer
	if acc.GetPubKey() != nil {
		pubKey = acc.GetPubKey().Bytes()
	}

	return User{
		Address:       acc.GetAddress().String(),
		AccountNumber: acc.GetAccountNumber(),
		Coins:         acc.GetCoins(),
		Sequence:      acc.GetSequence(),
		Pubkey:        tcmn.HexBytes(pubKey),
	}
}
