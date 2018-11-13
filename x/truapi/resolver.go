package truapi

import (
	"context"
	"fmt"
	"strconv"

	"github.com/TruStory/truchain/x/category"
	"github.com/TruStory/truchain/x/story"
	"github.com/TruStory/truchain/x/users"
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

func (ta *TruAPI) storyCategoryResolver(ctx context.Context, q story.Story) category.Category {
	res := ta.RunQuery("categories/id", category.QueryCategoryParams{ID: strconv.FormatInt(q.CategoryID, 10)})

	if res.Code != 0 {
		fmt.Println("Resolver err: ", res)
		return category.Category{}
	}

	s := new(category.Category)
	err := amino.UnmarshalJSON(res.Value, s)

	if err != nil {
		panic(err)
	}

	return *s
}

func (ta *TruAPI) usersResolver(ctx context.Context, q users.QueryUsersByAddressesParams) []users.User {
	res := ta.RunQuery("users/addresses", q)

	if res.Code != 0 {
		fmt.Println("Resolver err: ", res)
		return []users.User{}
	}

	s := new([]users.User)

	err := amino.UnmarshalJSON(res.Value, s)

	if err != nil {
		panic(err)
	}

	return *s
}

func (ta *TruAPI) twitterProfileResolver(ctx context.Context, q users.User) users.TwitterProfile {
	addr := q.Address
	fmt.Println("\"Fetching\" Twitter profile for address: " + addr)
	return users.TwitterProfile{ID: "1234567890123456789", Username: "someone", FullName: "Some Person", Address: addr}
}
