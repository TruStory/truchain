package users

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/tendermint/tendermint/crypto/ed25519"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/stretchr/testify/assert"
)

func Test_NewUserEmptyPublicKey(t *testing.T) {
	// Check that NewUser don't panick by checking that recover() result is nil
	defer func() {
		assert.Nil(t, recover())
	}()
	key := ed25519.GenPrivKey()
	pub := key.PubKey()
	addr := sdk.AccAddress(pub.Address())
	acc := auth.ProtoBaseAccount()
	acc.SetAddress(addr)
	user := NewUser(acc)
	assert.Nil(t, user.Pubkey)
}
