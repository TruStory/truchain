package registration

import (
  "encoding/json"

  "github.com/cosmos/cosmos-sdk/x/auth"
  "github.com/TruStory/truchain/x/chttp"

  app "github.com/TruStory/truchain/types"
  sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewHandler(am auth.AccountMapper) sdk.Handler {
  return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
    switch msg := msg.(type) {
    case RegisterKeyMsg:
      return handleRegisterKeyMsg(ctx, am, msg)
    default:
      return app.ErrMsgHandler(msg)
    }
  }
}

// ============================================================================

func handleRegisterKeyMsg(ctx sdk.Context, am auth.AccountMapper, msg RegisterKeyMsg) sdk.Result {
  bacc := auth.NewBaseAccountWithAddress(msg.Address)
  key, err := chttp.StdKey(msg.PubKeyAlgo, msg.PubKey.Bytes())

  if err != nil {
    return sdk.Result{Code: 1, Data: []byte("Registration Error: parsing public key: " + err.Error())}
  }

  err = bacc.SetPubKey(*key)

  if err != nil {
    return sdk.Result{Code: 1, Data: []byte("Registration Error: setting public key: " + err.Error())}
  }

  err = bacc.SetCoins(msg.Coins)

  if err != nil {
    return sdk.Result{Code: 1, Data: []byte("Registration Error: setting coins: " + err.Error())}
  }

  acc := app.NewAppAccount(string(msg.Address), bacc)

  am.SetAccount(ctx, *acc)
  
  bz, _ := json.Marshal(*acc)

  return sdk.Result{Code: 0, Data: bz}
}
