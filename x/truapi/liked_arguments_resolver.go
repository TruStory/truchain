package truapi

import (
	"context"
	"encoding/json"
	"fmt"
	"path"

	app "github.com/TruStory/truchain/types"
	"github.com/TruStory/truchain/x/backing"
	"github.com/TruStory/truchain/x/challenge"
	"github.com/TruStory/truchain/x/story"
	trubank "github.com/TruStory/truchain/x/trubank"
)

func (ta *TruAPI) likedArguments(
	ctx context.Context, q app.QueryTrasanctionsByCreatorAndCategoryParams) []LikedArgument {
	likedArguments := make([]LikedArgument, 0)
	transactions := make([]trubank.Transaction, 0)
	res := ta.RunQuery(
		path.Join(trubank.QueryPath, trubank.QueryLikeTransactionsByCreator), q)

	if res.Code != 0 {
		fmt.Println("Resolver err: ", res)
		return likedArguments
	}

	err := json.Unmarshal(res.Value, &transactions)
	if err != nil {
		panic(err)
	}
	queryBacking := path.Join(backing.QueryPath, backing.QueryBackingByID)
	queryChallenge := path.Join(challenge.QueryPath, challenge.QueryChallengeByID)

	for _, tx := range transactions {
		likedArgument := LikedArgument{
			Transaction: tx,
		}
		// if category Id is sent filter by category
		if q.CategoryID != nil {
			story := ta.storyResolver(ctx, story.QueryStoryByIDParams{ID: tx.GroupID})
			if story.CategoryID != *q.CategoryID {
				continue
			}
		}
		switch tx.TransactionType {
		case trubank.BackingLike:
			var backing backing.Backing
			res := ta.RunQuery(queryBacking, app.QueryByIDParams{ID: tx.ReferenceID})
			if res.Code != 0 {
				fmt.Println("error getting backing", res)
				continue
			}
			err := json.Unmarshal(res.Value, &backing)
			if err != nil {
				panic(err)
			}
			likedArgument.Stake = *backing.Vote
		case trubank.ChallengeLike:
			var challenge challenge.Challenge
			res := ta.RunQuery(queryChallenge, app.QueryByIDParams{ID: tx.ReferenceID})
			if res.Code != 0 {
				fmt.Println("error getting challenge", res)
				continue
			}
			err := json.Unmarshal(res.Value, &challenge)
			if err != nil {
				panic(err)
			}
			likedArgument.Stake = *challenge.Vote
		}
		argument := ta.argumentResolver(ctx, app.QueryByIDParams{ID: likedArgument.Stake.ArgumentID})
		likedArgument.Argument = argument
		likedArguments = append(likedArguments, likedArgument)
	}
	return likedArguments
}
