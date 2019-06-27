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
		default:
			errMsg := fmt.Sprintf("Unrecognized bank message type: %T", msg)
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
	res, codecErr := ModuleCodec.MarshalBinaryBare(argument)
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
	res, codecErr := ModuleCodec.MarshalBinaryBare(stake)
	if codecErr != nil {
		return sdk.ErrInternal(fmt.Sprintf("Marshal result error: %s", codecErr)).Result()
	}
	return sdk.Result{
		Data: res,
	}
}
