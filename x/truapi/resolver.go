package truapi

import (
	"context"
	"encoding/json"
	"fmt"
	"path"
	"sort"

	app "github.com/TruStory/truchain/types"
	"github.com/TruStory/truchain/x/argument"
	"github.com/TruStory/truchain/x/backing"
	"github.com/TruStory/truchain/x/category"
	"github.com/TruStory/truchain/x/challenge"
	"github.com/TruStory/truchain/x/db"
	"github.com/TruStory/truchain/x/params"
	"github.com/TruStory/truchain/x/story"
	"github.com/TruStory/truchain/x/truapi/cookies"
	trubank "github.com/TruStory/truchain/x/trubank"
	"github.com/TruStory/truchain/x/users"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/kelseyhightower/envconfig"
	amino "github.com/tendermint/go-amino"
	abci "github.com/tendermint/tendermint/abci/types"
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

// Deprecated in favor of storiesResolver
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

	filteredStories, err := ta.filterFlaggedStories(stories)
	if err != nil {
		fmt.Println("Resolver err: ", err)
		panic(err)
	}

	return filteredStories
}

func (ta *TruAPI) storiesResolver(_ context.Context, q app.QueryByCategoryIDParams) []story.Story {
	var res abci.ResponseQuery
	if q.CategoryID == -1 {
		res = ta.RunQuery("stories/all", struct{}{})
	} else {
		res = ta.RunQuery("stories/category", story.QueryCategoryStoriesParams{CategoryID: q.CategoryID})
	}

	if res.Code != 0 {
		fmt.Println("Resolver err: ", res)
		return []story.Story{}
	}

	stories := new([]story.Story)
	err := json.Unmarshal(res.Value, stories)
	if err != nil {
		panic(err)
	}

	filteredStories, err := ta.filterFlaggedStories(stories)
	if err != nil {
		fmt.Println("Resolver err: ", err)
		panic(err)
	}

	return filteredStories
}

func (ta *TruAPI) argumentResolver(_ context.Context, q app.QueryByIDParams) argument.Argument {
	res := ta.RunQuery(
		path.Join(argument.QueryPath, argument.QueryArgumentByID),
		app.QueryByIDParams{ID: q.ID},
	)

	if res.Code != 0 {
		fmt.Println("Resolver err: ", res)
		return argument.Argument{}
	}

	argument := new(argument.Argument)
	err := json.Unmarshal(res.Value, argument)
	if err != nil {
		panic(err)
	}

	return *argument
}

func (ta *TruAPI) likesObjectResolver(_ context.Context, q app.QueryByIDParams) []argument.Like {
	query := path.Join(argument.QueryPath, argument.QueryLikesByArgumentID)
	res := ta.RunQuery(query, app.QueryByIDParams{ID: q.ID})

	if res.Code != 0 {
		fmt.Println("Resolver err: ", res)
		return []argument.Like{}
	}

	likes := new([]argument.Like)
	err := json.Unmarshal(res.Value, likes)
	if err != nil {
		panic(err)
	}

	return *likes
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

func (ta *TruAPI) backingPoolResolver(_ context.Context, q story.Story) sdk.Coin {
	res := ta.RunQuery(path.Join(backing.QueryPath, backing.QueryBackingAmountByStoryID), app.QueryByIDParams{ID: q.ID})

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

func (ta *TruAPI) challengePoolResolver(_ context.Context, q story.Story) sdk.Coin {
	res := ta.RunQuery(path.Join(challenge.QueryPath, challenge.QueryChallengeAmountByStoryID), app.QueryByIDParams{ID: q.ID})

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

// Deprecated in favor of storiesResolver
func (ta *TruAPI) categoryStoriesResolver(_ context.Context, q category.Category) []story.Story {
	res := ta.RunQuery("stories/category", story.QueryCategoryStoriesParams{CategoryID: q.ID})

	if res.Code != 0 {
		fmt.Println("Resolver err: ", res)
		return []story.Story{}
	}

	stories := new([]story.Story)
	err := json.Unmarshal(res.Value, stories)
	if err != nil {
		panic(err)
	}

	filteredStories, err := ta.filterFlaggedStories(stories)
	if err != nil {
		fmt.Println("Resolver err: ", err)
		panic(err)
	}

	return filteredStories
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

	res := ta.RunQuery(
		path.Join(challenge.QueryPath, challenge.QueryByStoryID), q)

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
		// TODO [shanev]: Add back after adding error handling to resolvers
		// fmt.Println("Resolver err: ", err)
		return db.TwitterProfile{}
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

func (ta *TruAPI) transactionsResolver(
	_ context.Context, q app.QueryByCreatorParams) []trubank.Transaction {

	res := ta.RunQuery(
		path.Join(trubank.QueryPath, trubank.QueryTransactionsByCreator), q)

	if res.Code != 0 {
		fmt.Println("Resolver err: ", res)
		return []trubank.Transaction{}
	}

	transactions := new([]trubank.Transaction)
	err := json.Unmarshal(res.Value, transactions)
	if err != nil {
		panic(err)
	}

	return *transactions
}

func (ta *TruAPI) notificationsResolver(ctx context.Context, q struct{}) []db.NotificationEvent {
	user, ok := ctx.Value(userContextKey).(*cookies.AuthenticatedUser)
	if !ok {
		return make([]db.NotificationEvent, 0)
	}
	evts, err := ta.DBClient.NotificationEventsByAddress(user.Address)
	if err != nil {
		panic(err)
	}
	return evts
}

func (ta *TruAPI) addressesWhoFlaggedResolver(ctx context.Context, q story.Story) []string {
	flaggedStories, err := ta.DBClient.FlaggedStoriesByStoryID(q.ID)
	if err != nil {
		return []string{}
	}
	var addressesWhoFlagged []string
	for _, story := range flaggedStories {
		addressesWhoFlagged = append(addressesWhoFlagged, story.Creator)
	}
	return addressesWhoFlagged
}

func (ta *TruAPI) filterFlaggedStories(stories *[]story.Story) ([]story.Story, error) {
	type FlagConfig struct {
		Limit int    `default:"4294967295"`
		Admin string `default:"cosmos1xqc5gwzpg3fyv5en2fzyx36z2se5ks33tt57e7"`
	}

	var flagConfig FlagConfig
	err := envconfig.Process("api_story_flag", &flagConfig)
	if err != nil {
		return nil, err
	}

	filteredStories := make([]story.Story, 0)
	for _, story := range *stories {
		storyFlags, err := ta.DBClient.FlaggedStoriesByStoryID(story.ID)
		if err != nil {
			return nil, err
		}
		if len(storyFlags) > 0 {
			if storyFlags[0].Creator == flagConfig.Admin {
				continue
			}
		}
		if len(storyFlags) < flagConfig.Limit {
			filteredStories = append(filteredStories, story)
		}
	}

	return filteredStories, nil
}

func (ta *TruAPI) commentsResolver(ctx context.Context, q argument.Argument) []db.Comment {
	comments, err := ta.DBClient.CommentsByArgumentID(q.ID)
	if err != nil {
		panic(err)
	}
	return comments
}
