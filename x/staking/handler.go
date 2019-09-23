package staking

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewHandler creates a new handler for staking module
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgSubmitArgument:
			return handleMsgSubmitArgument(ctx, keeper, msg)
		case MsgSubmitUpvote:
			return handleMsgSubmitUpvote(ctx, keeper, msg)
		case MsgEditArgument:
			return handleMsgEditArgument(ctx, keeper, msg)
		case MsgAddAdmin:
			return handleMsgAddAdmin(ctx, keeper, msg)
		case MsgRemoveAdmin:
			return handleMsgRemoveAdmin(ctx, keeper, msg)
		case MsgUpdateParams:
			return handleMsgUpdateParams(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized staking message type: %T", msg)
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleMsgSubmitArgument(ctx sdk.Context, keeper Keeper, msg MsgSubmitArgument) sdk.Result {
	if err := msg.ValidateBasic(); err != nil {
		return err.Result()
	}
	argument, err := keeper.SubmitArgument(ctx, msg.Body, msg.Summary, msg.Creator, msg.ClaimID, msg.StakeType)
	if err != nil {
		return err.Result()
	}
	res, codecErr := ModuleCodec.MarshalJSON(argument)
	if codecErr != nil {
		return sdk.ErrInternal(fmt.Sprintf("Marshal result error: %s", codecErr)).Result()
	}

	return sdk.Result{
		Data: res,
	}
}

func handleMsgSubmitUpvote(ctx sdk.Context, keeper Keeper, msg MsgSubmitUpvote) sdk.Result {
	if err := msg.ValidateBasic(); err != nil {
		return err.Result()
	}
	stake, err := keeper.SubmitUpvote(ctx, msg.ArgumentID, msg.Creator)
	if err != nil {
		return err.Result()
	}
	res, codecErr := ModuleCodec.MarshalJSON(stake)
	if codecErr != nil {
		return sdk.ErrInternal(fmt.Sprintf("Marshal result error: %s", codecErr)).Result()
	}

	return sdk.Result{
		Data: res,
	}
}

func handleMsgEditArgument(ctx sdk.Context, keeper Keeper, msg MsgEditArgument) sdk.Result {
	if err := msg.ValidateBasic(); err != nil {
		return err.Result()
	}
	argument, err := keeper.EditArgument(ctx, msg.Body, msg.Summary, msg.Creator, msg.ArgumentID)
	if err != nil {
		return err.Result()
	}
	res, codecErr := ModuleCodec.MarshalJSON(argument)
	if codecErr != nil {
		return sdk.ErrInternal(fmt.Sprintf("Marshal result error: %s", codecErr)).Result()
	}
	return sdk.Result{
		Data: res,
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
