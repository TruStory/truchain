package challenge

import (
	app "github.com/TruStory/truchain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewHandler creates a new handler for all challenge messages
func NewHandler(k WriteKeeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case SubmitChallengeMsg:
			return handleSubmitChallengeMsg(ctx, k, msg)
		default:
			return app.ErrMsgHandler(msg)
		}
	}
}

// ============================================================================

// handleStartChallengeMsg handles a message to start a challenge
func handleSubmitChallengeMsg(ctx sdk.Context, k WriteKeeper, msg SubmitChallengeMsg) sdk.Result {
	if err := msg.ValidateBasic(); err != nil {
		return err.Result()
	}

	id, err := k.Create(
		ctx, msg.StoryID, msg.Amount,
		msg.Argument, msg.Creator, msg.Evidence)
	if err != nil {
		return err.Result()
	}

	return app.Result(id)
}
