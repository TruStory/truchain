package bank

import (
	"encoding/json"
	"fmt"

	"github.com/TruStory/truchain/x/bank/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewHandler creates a new handler for bank module
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgPayReward:
			return handleMsgPayReward(ctx, keeper, msg)
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

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypePayGift,
			sdk.NewAttribute(sdk.AttributeKeyModule, ModuleName),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender.String()),
			sdk.NewAttribute(types.AttributeRecipient, msg.Recipient.String()),
		),
	)

	return sdk.Result{}
}

func handleMsgPayReward(ctx sdk.Context, keeper Keeper, msg MsgPayReward) sdk.Result {
	if err := msg.ValidateBasic(); err != nil {
		return err.Result()
	}
	err := keeper.payReward(ctx, msg.Sender, msg.Recipient, msg.Reward, msg.InviteID)
	if err != nil {
		return err.Result()
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypePayReward,
			sdk.NewAttribute(sdk.AttributeKeyModule, ModuleName),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender.String()),
			sdk.NewAttribute(types.AttributeRecipient, msg.Recipient.String()),
		),
	)

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
