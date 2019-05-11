package users

import (
	"encoding/json"
	"fmt"

	app "github.com/TruStory/truchain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

// NewHandler returns a handler for messages of type RegisterKeyMsg
func NewHandler(ak auth.AccountKeeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case RegisterKeyMsg:
			return handleRegisterKeyMsg(ctx, ak, msg)
		default:
			return app.ErrMsgHandler(msg)
		}
	}
}

// ============================================================================

func handleRegisterKeyMsg(ctx sdk.Context, ak auth.AccountKeeper, msg RegisterKeyMsg) sdk.Result {
	bacc := auth.NewBaseAccountWithAddress(msg.Address)
	key, err := stdKey(msg.PubKeyAlgo, msg.PubKey.Bytes())

	if err != nil {
		return sdk.Result{Code: 1, Data: []byte("Registration Error: parsing public key: " + err.Error())}
	}

	err = bacc.SetPubKey(key)
	if err != nil {
		return sdk.Result{Code: 1, Data: []byte("Registration Error: setting public key: " + err.Error())}
	}

	err = bacc.SetCoins(msg.Coins)
	if err != nil {
		return sdk.Result{Code: 1, Data: []byte("Registration Error: setting coins: " + err.Error())}
	}

	acc := app.NewAppAccount(bacc)
	ak.SetAccount(ctx, acc)
	bz, err := json.Marshal(*acc)
	if err != nil {
		return sdk.Result{Code: 1, Data: []byte("Registration Error: marshaling account: " + err.Error())}
	}

	return sdk.Result{Code: 0, Data: bz}
}

// stdKey returns an instance of `crypto.PubKey` using the given algorithm
func stdKey(algo string, bytes []byte) (crypto.PubKey, error) {
	switch algo {
	case "ed25519":
		ek := ed25519.PubKeyEd25519{}
		copy(ek[:], bytes)
		return ek, nil
	case "secp256k1":
		sk := secp256k1.PubKeySecp256k1{}
		copy(sk[:], bytes)
		return sk, nil
	default:
		return secp256k1.PubKeySecp256k1{}, unsupportedAlgoError(algo, []string{"ed25519", "secp256k1"})
	}
}

func unsupportedAlgoError(name string, supported []string) error {
	s := "Tx Error: Unsupported public key algorithm \"%s\" (supported: %v)"
	return fmt.Errorf(s, name, supported)
}
