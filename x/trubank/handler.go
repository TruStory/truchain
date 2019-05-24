package trubank

import (
	app "github.com/TruStory/truchain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewHandler creates a new handler for all trubank messages
func NewHandler(k WriteKeeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case PayRewardMsg:
			return handlePayRewardMsg(ctx, k, msg)
		default:
			return app.ErrMsgHandler(msg)
		}
	}
}

// ============================================================================

func handlePayRewardMsg(ctx sdk.Context, k WriteKeeper, msg PayRewardMsg) sdk.Result {
	if err := msg.ValidateBasic(); err != nil {
		return err.Result()
	}

	// TODO: move this into genesis.json and import via parameters
	adminAccount, addrErr := sdk.AccAddressFromBech32("cosmos1xqc5gs2xfdryws6dtfvng3z32ftr2de56tksud")
	if addrErr != nil {
		return sdk.ErrInvalidAddress("Could not access admin account").Result()
	}

	if !msg.Creator.Equals(adminAccount) {
		return sdk.ErrInvalidAddress("You are not authorized to distribute rewards").Result()
	}

	_, err := k.AddCoin(
		ctx,
		msg.Recipient,
		msg.Reward,
		0, // No StoryID
		InviteAFriend,
		msg.InviteID,
	)
	if err != nil {
		return err.Result()
	}

	return app.Result(msg.InviteID)
}
