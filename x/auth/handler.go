package auth

import (
	"fmt"

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
	if err := msg.ValidateBasic(); err != nil {
		return err.Result()
	}

	appAccount := k.NewAppAccount(ctx, msg.Address, msg.Coins, msg.PubKey)

	res, jsonErr := k.codec.MarshalJSON(appAccount)
	if jsonErr != nil {
		return sdk.ErrInternal(fmt.Sprintf("Marshal result error: %s", jsonErr)).Result()
	}

	return sdk.Result{
		Data: res,
	}
}
