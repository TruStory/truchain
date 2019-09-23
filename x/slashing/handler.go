package slashing

import (
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewHandler creates a new handler for slashing module
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgSlashArgument:
			return handleMsgSlashArgument(ctx, keeper, msg)
		case MsgAddAdmin:
			return handleMsgAddAdmin(ctx, keeper, msg)
		case MsgRemoveAdmin:
			return handleMsgRemoveAdmin(ctx, keeper, msg)
		case MsgUpdateParams:
			return handleMsgUpdateParams(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized slashing message type: %T", msg)
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleMsgSlashArgument(ctx sdk.Context, keeper Keeper, msg MsgSlashArgument) sdk.Result {
	if err := msg.ValidateBasic(); err != nil {
		return err.Result()
	}

	slash, punishmentResults, err := keeper.CreateSlash(ctx, msg.ArgumentID, msg.SlashType, msg.SlashReason, msg.SlashDetailedReason, msg.Creator)
	if err != nil {
		return err.Result()
	}

	res, jsonErr := ModuleCodec.MarshalJSON(slash)
	if jsonErr != nil {
		return sdk.ErrInternal(fmt.Sprintf("Marshal result error: %s", jsonErr)).Result()
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(AttributeKeyMinSlashCountKey, fmt.Sprintf("%d", keeper.GetParams(ctx).MinSlashCount)),
		),
	)

	if len(punishmentResults) > 0 {
		_json, jsonErr := json.Marshal(punishmentResults)
		if jsonErr != nil {
			return sdk.ErrInternal(fmt.Sprintf("Marshal result error: %s", jsonErr)).Result()
		}

		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				sdk.EventTypeMessage,
				sdk.NewAttribute(AttributeKeySlashResults, string(_json)),
			),
		)
	}

	return sdk.Result{
		Data:   res,
		Events: ctx.EventManager().Events(),
	}
}

func handleMsgAddAdmin(ctx sdk.Context, k Keeper, msg MsgAddAdmin) sdk.Result {
	if err := msg.ValidateBasic(); err != nil {
		return err.Result()
	}

	err := k.AddAdmin(ctx, msg.Admin, msg.Creator)
	if err != nil {
		return err.Result()
	}

	res, jsonErr := ModuleCodec.MarshalJSON(true)
	if jsonErr != nil {
		return sdk.ErrInternal(fmt.Sprintf("Marshal result error: %s", jsonErr)).Result()
	}

	return sdk.Result{
		Data: res,
	}
}

func handleMsgRemoveAdmin(ctx sdk.Context, k Keeper, msg MsgRemoveAdmin) sdk.Result {
	if err := msg.ValidateBasic(); err != nil {
		return err.Result()
	}

	err := k.RemoveAdmin(ctx, msg.Admin, msg.Remover)
	if err != nil {
		return err.Result()
	}

	res, jsonErr := ModuleCodec.MarshalJSON(true)
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

	err := k.UpdateParams(ctx, msg.Updater, msg.Updates, msg.UpdatedFields)

	if err != nil {
		return err.Result()
	}

	res, jsonErr := ModuleCodec.MarshalJSON(true)

	if jsonErr != nil {
		return sdk.ErrInternal(fmt.Sprintf("Marshal result error: %s", jsonErr)).Result()
	}

	return sdk.Result{
		Data: res,
	}
}
