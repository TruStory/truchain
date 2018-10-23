package challenge

import (
	app "github.com/TruStory/truchain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewHandler creates a new handler for all challenge messages
func NewHandler(k WriteKeeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case StartChallengeMsg:
			return handleStartChallengeMsg(ctx, k, msg)
		case UpdateChallengeMsg:
			return handleUpdateChallengeMsg(ctx, k, msg)
		default:
			return app.ErrMsgHandler(msg)
		}
	}
}

// ============================================================================

// handleStartChallengeMsg handles a message to start a challenge
func handleStartChallengeMsg(ctx sdk.Context, k WriteKeeper, msg StartChallengeMsg) sdk.Result {
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

// handleUpdateChallengeMsg handles a message to add a new challenger
func handleUpdateChallengeMsg(ctx sdk.Context, k WriteKeeper, msg UpdateChallengeMsg) sdk.Result {
	if err := msg.ValidateBasic(); err != nil {
		return err.Result()
	}

	id, err := k.Update(ctx, msg.ChallengeID, msg.Creator, msg.Amount)
	if err != nil {
		return err.Result()
	}

	return app.Result(id)
}
