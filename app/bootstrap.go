package app

import (
	"fmt"
	"net/url"
	"time"

	tru "github.com/TruStory/truchain/types"
	"github.com/TruStory/truchain/x/story"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func createFakeStory(ctx sdk.Context, sk story.WriteKeeper) int64 {
	body := "Body of story."

	ctx = ctx.WithBlockHeader(abci.Header{Time: time.Now().UTC()})

	catID := int64(1)
	creator := sdk.AccAddress([]byte{1, 2})
	storyType := story.Default
	source := url.URL{}

	url, _ := url.Parse("http://shanesbrain.net")
	e := story.Evidence{
		Creator:   creator,
		URL:       *url,
		Timestamp: tru.NewTimestamp(ctx.BlockHeader()),
	}
	evidence := []story.Evidence{e}
	argument := []story.Argument{}

	storyID, _ := sk.NewStory(ctx, argument, body, catID, creator, evidence, source, storyType)

	return storyID
}

func loadTestDB(ctx sdk.Context, storyKeeper story.WriteKeeper) {
	storyID := createFakeStory(ctx, storyKeeper)
	fmt.Println(storyID)
}
