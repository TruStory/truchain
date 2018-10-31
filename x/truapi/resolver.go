package truapi

import (
	"context"
	"fmt"

	"github.com/TruStory/truchain/x/story"
	"github.com/tendermint/go-amino"
)

func (ta *TruAPI) storyResolver(_ context.Context, q story.QueryCategoryStoriesParams) []story.Story {
	res := ta.RunQuery("stories/category", q)

	if res.Code != 0 {
		fmt.Println("Resolver err: ", res)
		return []story.Story{}
	}

	s := new([]story.Story)
	err := amino.UnmarshalJSON(res.Value, s)

	if err != nil {
		panic(err)
	}

	return *s
}
