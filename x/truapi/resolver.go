package truapi

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"

	app "github.com/TruStory/truchain/types"
	"github.com/TruStory/truchain/x/category"
	"github.com/TruStory/truchain/x/game"
	"github.com/TruStory/truchain/x/story"
	"github.com/TruStory/truchain/x/users"
	sdk "github.com/cosmos/cosmos-sdk/types"
	amino "github.com/tendermint/go-amino"
)

func (ta *TruAPI) categoryStoriesResolver(_ context.Context, q category.Category) []story.Story {
	res := ta.RunQuery("stories/category", story.QueryCategoryStoriesParams{CategoryID: q.ID})

	if res.Code != 0 {
		fmt.Println("Resolver err: ", res)
		return []story.Story{}
	}

	s := new([]story.Story)
	err := json.Unmarshal(res.Value, s)

	if err != nil {
		panic(err)
	}

	return *s
}

func (ta *TruAPI) storyResolver(_ context.Context, q story.QueryStoryByIDParams) story.Story {
	res := ta.RunQuery("stories/id", q)

	if res.Code != 0 {
		fmt.Println("Resolver err: ", res)
		return story.Story{}
	}

	s := new(story.Story)
	err := json.Unmarshal(res.Value, s)

	if err != nil {
		panic(err)
	}

	return *s
}

func (ta *TruAPI) allCategoriesResolver(ctx context.Context, q struct{}) []category.Category {
	res := ta.RunQuery("categories/all", struct{}{})

	if res.Code != 0 {
		fmt.Println("Resolver err: ", res)
		return []category.Category{}
	}

	cs := new([]category.Category)
	err := json.Unmarshal(res.Value, cs)

	if err != nil {
		panic(err)
	}

	return *cs
}

func (ta *TruAPI) categoryResolver(ctx context.Context, q category.QueryCategoryByIDParams) category.Category {
	res := ta.RunQuery("categories/id", q)

	if res.Code != 0 {
		fmt.Println("Resolver err: ", res)
		return category.Category{}
	}

	c := new(category.Category)
	err := json.Unmarshal(res.Value, c)

	if err != nil {
		panic(err)
	}

	return *c
}

func (ta *TruAPI) storyCategoryResolver(ctx context.Context, q story.Story) category.Category {
	return ta.categoryResolver(ctx, category.QueryCategoryByIDParams{ID: q.CategoryID})
}

func (ta *TruAPI) usersResolver(ctx context.Context, q users.QueryUsersByAddressesParams) []users.User {
	res := ta.RunQuery("users/addresses", q)

	if res.Code != 0 {
		fmt.Println("Resolver err: ", res)
		return []users.User{}
	}

	u := new([]users.User)

	err := amino.UnmarshalJSON(res.Value, u)

	if err != nil {
		panic(err)
	}

	return *u
}

func (ta *TruAPI) gameResolver(_ context.Context, q story.Story) game.Game {
	res := ta.RunQuery("games/id", game.QueryGameByIDParams{ID: q.GameID})

	if res.Code != 0 {
		fmt.Println("Resolver err: ", res)
		return game.Game{}
	}

	g := new(game.Game)

	err := json.Unmarshal(res.Value, g)

	if err != nil {
		panic(err)
	}

	return *g
}

func (ta *TruAPI) backingTotalResolver(_ context.Context, q story.Story) sdk.Coin {
	res := ta.RunQuery("backings/totalAmountByStoryID", app.QueryByIDParams{ID: q.ID})

	if res.Code != 0 {
		fmt.Println("Resolver err: ", res)
		return sdk.Coin{}
	}

	amount := new(sdk.Coin)
	err := amino.UnmarshalJSON(res.Value, amount)
	if err != nil {
		panic(err)
	}

	return *amount
}

func (ta *TruAPI) challengeThresholdResolver(_ context.Context, q game.Game) sdk.Coin {
	res := ta.RunQuery("games/challengeThresholdByGameID", app.QueryByIDParams{ID: q.ID})

	if res.Code != 0 {
		fmt.Println("Resolver err: ", res)
		return sdk.Coin{}
	}

	amount := new(sdk.Coin)
	err := json.Unmarshal(res.Value, amount)
	if err != nil {
		panic(err)
	}

	return *amount
}

// TODO: [shanev/truted] Handle this when working on user profiles
// https://github.com/TruStory/truchain/issues/196
func (ta *TruAPI) twitterProfileResolver(ctx context.Context, q users.User) users.TwitterProfile {
	addr := q.Address
	fmt.Println("Mocking ('fetching') Twitter profile for address: " + addr)
	return users.TwitterProfile{
		ID:        "1234567890123456789",
		Username:  "someone",
		FullName:  "Some Person",
		Address:   addr,
		AvatarURI: fmt.Sprintf("https://randomuser.me/api/portraits/thumb/men/%d.jpg", rand.Intn(50)+1),
	}
}
