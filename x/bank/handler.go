package bank

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
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
	err := keeper.payReward(ctx, msg.Creator, msg.Recipient, msg.Reward, msg.InviteID)
	if err != nil {
		return err.Result()
	}
	return sdk.Result{}
}
