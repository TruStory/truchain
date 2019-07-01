package account

import (
	"fmt"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tendermint/tendermint/crypto/secp256k1"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewHandler creates a new handler for auth module
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgRegisterKey:
			return handleMsgRegisterKey(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized auth message type: %T", msg)
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleMsgRegisterKey(ctx sdk.Context, k Keeper, msg MsgRegisterKey) sdk.Result {
	pubKey, err := toPubKey(msg.PubKeyAlgo, msg.PubKey.Bytes())
	if err != nil {
		return sdk.Result{Code: 1, Data: []byte("Registration Error: parsing public key: " + err.Error())}
	}

	appAccount, sdkErr := k.CreateAppAccount(ctx, msg.Address, msg.Coins, pubKey)
	if sdkErr != nil {
		return sdkErr.Result()
	}

	res, jsonErr := k.codec.MarshalJSON(appAccount)
	if jsonErr != nil {
		return sdk.ErrInternal(fmt.Sprintf("Marshal result error: %s", jsonErr)).Result()
	}

	return sdk.Result{
		Data: res,
	}
}

// toPubKey returns an instance of `crypto.PubKey` using the given algorithm
func toPubKey(algo string, rawPubKeyBytes []byte) (crypto.PubKey, error) {
	switch algo {
	case "ed25519":
		ek := ed25519.PubKeyEd25519{}
		copy(ek[:], rawPubKeyBytes)
		return ek, nil
	case "secp256k1":
		sk := secp256k1.PubKeySecp256k1{}
		copy(sk[:], rawPubKeyBytes)
		return sk, nil
	default:
		return secp256k1.PubKeySecp256k1{}, unsupportedAlgoError(algo, []string{"ed25519", "secp256k1"})
	}
}

func unsupportedAlgoError(name string, supported []string) error {
	s := "Tx Error: Unsupported public key algorithm \"%s\" (supported: %v)"
	return fmt.Errorf(s, name, supported)
}
