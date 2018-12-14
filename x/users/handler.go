package users

import (
	"encoding/json"

	app "github.com/TruStory/truchain/types"
	"github.com/TruStory/truchain/x/chttp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
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
	key, err := chttp.StdKey(msg.PubKeyAlgo, msg.PubKey.Bytes())

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
