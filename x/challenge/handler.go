package challenge

import (
	app "github.com/TruStory/truchain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewHandler creates a new handler for all TruStory messages
func NewHandler(k WriteKeeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case StartChallengeMsg:
			return handleStartChallengeMsg(ctx, k, msg)
		default:
			return app.ErrMsgHandler(msg)
		}
	}
}

// // ============================================================================

func handleStartChallengeMsg(ctx sdk.Context, k WriteKeeper, msg StartChallengeMsg) sdk.Result {
	if err := msg.ValidateBasic(); err != nil {
		return err.Result()
	}

	id, err := k.NewChallenge(
		ctx,
		msg.StoryID,
		msg.Amount,
		msg.Argument,
		msg.Creator,
		msg.Evidence,
		msg.Reason)
	if err != nil {
		return err.Result()
	}

	return app.Result(id)
}
