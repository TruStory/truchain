package community

import (
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewHandler creates a new handler for all community messages
func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgNewCommunity:
			return handleMsgNewCommunity(ctx, k, msg)
		case MsgUpdateParams:
			return handleMsgUpdateParams(ctx, k, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized community message type: %T", msg)
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleMsgNewCommunity(ctx sdk.Context, k Keeper, msg MsgNewCommunity) sdk.Result {
	if err := msg.ValidateBasic(); err != nil {
		return err.Result()
	}

	community, err := k.NewCommunity(ctx, msg.ID, msg.Name, msg.Description, msg.Creator)
	if err != nil {
		return err.Result()
	}

	res, jsonErr := json.Marshal(community)
	if jsonErr != nil {
		return sdk.ErrInternal(fmt.Sprintf("Marshal result error: %s", jsonErr)).Result()
	}

	return sdk.Result{
		Data: res,
	}
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
