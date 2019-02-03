package truapi

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"sort"

	"github.com/TruStory/truchain/x/db"
	"github.com/TruStory/truchain/x/params"
	"github.com/TruStory/truchain/x/vote"

	"github.com/TruStory/truchain/x/challenge"

	"github.com/TruStory/truchain/x/backing"

	app "github.com/TruStory/truchain/types"
	"github.com/TruStory/truchain/x/category"
	"github.com/TruStory/truchain/x/game"
	"github.com/TruStory/truchain/x/story"
	"github.com/TruStory/truchain/x/users"
	sdk "github.com/cosmos/cosmos-sdk/types"
	amino "github.com/tendermint/go-amino"
)

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

	// sort in alphabetical order
	sort.Slice(*cs, func(i, j int) bool {
		return (*cs)[j].Title > (*cs)[i].Title
	})

	return *cs
}

func (ta *TruAPI) allStoriesResolver(ctx context.Context, q struct{}) []story.Story {
	res := ta.RunQuery("stories/all", struct{}{})

	if res.Code != 0 {
		fmt.Println("Resolver err: ", res)
		return []story.Story{}
	}

	stories := new([]story.Story)
	err := json.Unmarshal(res.Value, stories)
	if err != nil {
		panic(err)
	}

	return *stories
}

func (ta *TruAPI) backingResolver(
	_ context.Context, q app.QueryByStoryIDAndCreatorParams) backing.Backing {

	res := ta.RunQuery("backings/storyIDAndCreator", q)

	if res.Code != 0 {
		fmt.Println("Resolver err: ", res)
		return backing.Backing{}
	}

	backing := new(backing.Backing)
	err := json.Unmarshal(res.Value, backing)
	if err != nil {
		panic(err)
	}

	return *backing
}

func (ta *TruAPI) backingsResolver(
	_ context.Context, q app.QueryByIDParams) []backing.Backing {

	res := ta.RunQuery("backings/storyID", q)

	if res.Code != 0 {
		fmt.Println("Resolver err: ", res)
		return []backing.Backing{}
	}

	backings := new([]backing.Backing)
	err := json.Unmarshal(res.Value, backings)
	if err != nil {
		panic(err)
	}

	return *backings
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

func (ta *TruAPI) challengeResolver(
	_ context.Context, q app.QueryByStoryIDAndCreatorParams) challenge.Challenge {
	res := ta.RunQuery("challenges/storyIDAndCreator", q)

	if res.Code != 0 {
		fmt.Println("Resolver err: ", res)
		return challenge.Challenge{}
	}

	challenge := new(challenge.Challenge)
	err := json.Unmarshal(res.Value, challenge)
	if err != nil {
		panic(err)
	}

	return *challenge
}

func (ta *TruAPI) challengesResolver(
	_ context.Context, q app.QueryByIDParams) []challenge.Challenge {

	res := ta.RunQuery("challenges/gameID", q)

	if res.Code != 0 {
		fmt.Println("Resolver err: ", res)
		return []challenge.Challenge{}
	}

	challenges := new([]challenge.Challenge)
	err := json.Unmarshal(res.Value, challenges)
	if err != nil {
		panic(err)
	}

	return *challenges
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

func (ta *TruAPI) paramsResolver(_ context.Context) params.Params {
	res := ta.RunQuery("params", nil)

	if res.Code != 0 {
		fmt.Println("Resolver err: ", res)
		return params.Params{}
	}

	p := new(params.Params)
	err := json.Unmarshal(res.Value, p)
	if err != nil {
		panic(err)
	}

	return *p
}

func (ta *TruAPI) storyCategoryResolver(ctx context.Context, q story.Story) category.Category {
	return ta.categoryResolver(ctx, category.QueryCategoryByIDParams{ID: q.CategoryID})
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

func (ta *TruAPI) twitterProfileResolver(
	ctx context.Context, q users.User) db.TwitterProfile {

	addr := q.Address
	twitterProfile, err := ta.DBClient.TwitterProfileByAddress(addr)
	if err != nil {
		fmt.Println("Twitter profile not found. Returning mock data.")
		return db.TwitterProfile{
			ID:        1234,
			Username:  "someone",
			FullName:  "Some Person",
			Address:   addr,
			AvatarURI: fmt.Sprintf("https://randomuser.me/api/portraits/thumb/women/%d.jpg", rand.Intn(50)+1),
		}
	}

	return twitterProfile
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

func (ta *TruAPI) voteResolver(
	_ context.Context, q app.QueryByStoryIDAndCreatorParams) vote.TokenVote {
	res := ta.RunQuery("votes/storyIDAndCreator", q)

	if res.Code != 0 {
		fmt.Println("Resolver err: ", res)
		return vote.TokenVote{}
	}

	tokenVote := new(vote.TokenVote)
	err := json.Unmarshal(res.Value, tokenVote)
	if err != nil {
		panic(err)
	}

	return *tokenVote
}

func (ta *TruAPI) votesResolver(
	_ context.Context, q app.QueryByIDParams) []vote.TokenVote {

	res := ta.RunQuery("votes/gameID", q)

	if res.Code != 0 {
		fmt.Println("Resolver err: ", res)
		return []vote.TokenVote{}
	}

	tokenVotes := new([]vote.TokenVote)
	err := json.Unmarshal(res.Value, tokenVotes)
	if err != nil {
		panic(err)
	}

	return *tokenVotes
}

func (ta *TruAPI) votesTotalAmountResolver(_ context.Context, q game.Game) sdk.Coin {
	res := ta.RunQuery(
		"votes/totalVoteAmountByGameID", app.QueryByIDParams{ID: q.ID})

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
