package story

import (
	"net/url"

	app "github.com/TruStory/truchain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewHandler creates a new handler for all TruStory messages
func NewHandler(k WriteKeeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case SubmitStoryMsg:
			return handleSubmitStoryMsg(ctx, k, msg)
		default:
			return app.ErrMsgHandler(msg)
		}
	}
}

// ============================================================================

func handleSubmitStoryMsg(ctx sdk.Context, k WriteKeeper, msg SubmitStoryMsg) sdk.Result {
	if err := msg.ValidateBasic(); err != nil {
		return err.Result()
	}

	// parse url from string
	sourceURL, urlError := url.ParseRequestURI(msg.Source)
	if urlError != nil {
		return ErrInvalidSourceURL(msg.Source).Result()
	}

	// create evidence type from url
	var evidence []Evidence
	for _, url := range msg.Evidence {
		e := Evidence{
			Creator:   msg.Creator,
			URL:       url,
			Timestamp: app.NewTimestamp(ctx.BlockHeader()),
		}
		evidence = append(evidence, e)
	}

	id, err := k.NewStory(
		ctx, msg.Body, msg.CategoryID, msg.Creator, evidence, *sourceURL, msg.StoryType)
	if err != nil {
		return err.Result()
	}

	return app.Result(id)
}
