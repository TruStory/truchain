package slashing

import (
	"github.com/tendermint/tendermint/crypto"
)

// Slasher defines the admin account who can slash
type Slasher struct {
	PublicKey crypto.PubKey
}
