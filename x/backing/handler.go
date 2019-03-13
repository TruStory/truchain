package backing

import (
	"encoding/json"

	app "github.com/TruStory/truchain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewHandler creates a new handler for all TruStory messages
func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case BackStoryMsg:
			return handleBackStoryMsg(ctx, k, msg)
		default:
			return app.ErrMsgHandler(msg)
		}
	}
}

// ============================================================================

func handleBackStoryMsg(ctx sdk.Context, k Keeper, msg BackStoryMsg) sdk.Result {
	if err := msg.ValidateBasic(); err != nil {
		return err.Result()
	}

	id, err := k.Create(
		ctx,
		msg.StoryID,
		msg.Amount,
		msg.Argument,
		msg.Creator,
		false)
	if err != nil {
		return err.Result()
	}

	// story, err := k.storyKeeper.Story(ctx, msg.StoryID)
	// if err != nil {
	// 	return err.Result()
	// }

	tags := sdk.NewTags(
		"tru.event", []byte("Push"),
		// "push.body", []byte(pushBody),
		// "push.from", msg.Creator.Bytes(),
		// "push.to", story.Creator.Bytes(),
	)

	bz, _ := json.Marshal(app.IDResult{ID: id})

	return sdk.Result{
		Data: bz,
		Tags: tags,
	}
}
