package bank

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/TruStory/truchain/x/bank/tags"
)

// NewHandler creates a new handler for bank module
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgPayReward:
			return handleMsgPayReward(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized bank message type: %T", msg)
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleMsgPayReward(ctx sdk.Context, keeper Keeper, msg MsgPayReward) sdk.Result {
	if err := msg.ValidateBasic(); err != nil {
		return err.Result()
	}
	err := keeper.payReward(ctx, msg.Sender, msg.Recipient, msg.Reward, msg.InviteID)
	if err != nil {
		return err.Result()
	}
	tags := sdk.NewTags(
		tags.Category, tags.TxCategory,
		tags.Action, tags.ActionPayReward,
		tags.Sender, msg.Sender.String(),
		tags.Recipient, msg.Recipient.String(),
	)
	return sdk.Result{
		Tags: tags,
	}
}
