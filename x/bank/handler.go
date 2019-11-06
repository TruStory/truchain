package bank

import (
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewHandler creates a new handler for bank module
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgSendGift:
			return handleMsgSendGift(ctx, keeper, msg)
		case MsgUpdateParams:
			return handleMsgUpdateParams(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized bank message type: %T", msg)
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleMsgSendGift(ctx sdk.Context, keeper Keeper, msg MsgSendGift) sdk.Result {
	if err := msg.ValidateBasic(); err != nil {
		return err.Result()
	}
	err := keeper.sendGift(ctx, msg.Sender, msg.Recipient, msg.Reward)
	if err != nil {
		fmt.Println("error", err)
		return err.Result()
	}

	return sdk.Result{}
}

func handleMsgUpdateParams(ctx sdk.Context, k Keeper, msg MsgUpdateParams) sdk.Result {
	if err := msg.ValidateBasic(); err != nil {
		return err.Result()
	}

	err := k.UpdateParams(ctx, msg.Updates, msg.UpdatedFields)
	if err != nil {
		return err.Result()
	}

	res, jsonErr := json.Marshal(true)
	if jsonErr != nil {
		return sdk.ErrInternal(fmt.Sprintf("Marshal result error: %s", jsonErr)).Result()
	}

	return sdk.Result{
		Data: res,
	}
}
