package challenge

import (
	"net/url"

	app "github.com/TruStory/truchain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewHandler creates a new handler for all challenge messages
func NewHandler(k WriteKeeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case CreateChallengeMsg:
			return handleCreateChallengeMsg(ctx, k, msg)
		default:
			return app.ErrMsgHandler(msg)
		}
	}
}

// ============================================================================

// handles a message to create a challenge
func handleCreateChallengeMsg(
	ctx sdk.Context, k WriteKeeper, msg CreateChallengeMsg) sdk.Result {

	if err := msg.ValidateBasic(); err != nil {
		return err.Result()
	}

	var evidence []url.URL
	for _, urlString := range msg.Evidence {

		evidenceURL, urlError := url.ParseRequestURI(urlString)
		if urlError != nil {
			return ErrInvalidEvidenceURL(urlString).Result()
		}

		evidence = append(evidence, *evidenceURL)
	}

	id, err := k.Create(
		ctx, msg.StoryID, msg.Amount, msg.Argument, msg.Creator, evidence)
	if err != nil {
		return err.Result()
	}

	return app.Result(id)
}
