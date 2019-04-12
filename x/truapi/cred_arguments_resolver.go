package truapi

import (
	"context"
	"encoding/json"
	"fmt"
	"path"

	"github.com/TruStory/truchain/x/stake"

	app "github.com/TruStory/truchain/types"
	"github.com/TruStory/truchain/x/backing"
	"github.com/TruStory/truchain/x/challenge"
	"github.com/TruStory/truchain/x/story"
	trubank "github.com/TruStory/truchain/x/trubank"
)

func (ta *TruAPI) credArguments(
	ctx context.Context, q app.QueryTrasanctionsByCreatorAndCategoryParams) []CredArgument {
	credArguments := make([]CredArgument, 0)
	transactions := make([]trubank.Transaction, 0)
	res := ta.RunQuery(
		path.Join(trubank.QueryPath, trubank.QueryLikeTransactionsByCreator), q)

	if res.Code != 0 {
		fmt.Println("Resolver err: ", res)
		return credArguments
	}

	err := json.Unmarshal(res.Value, &transactions)
	if err != nil {
		panic(err)
	}

	mappedCategories := make(map[string]int64)
	for _, c := range ta.allCategoriesResolver(ctx, struct{}{}) {
		mappedCategories[c.Slug] = c.ID
	}

	queryBacking := path.Join(backing.QueryPath, backing.QueryBackingByID)
	queryChallenge := path.Join(challenge.QueryPath, challenge.QueryChallengeByID)

	filteredTransactions := make(map[string]trubank.Transaction)
	for _, tx := range transactions {
		key := mapArgumentTransaction(tx.TransactionType, tx.GroupID, tx.ReferenceID)
		filteredTransactions[key] = tx
	}
	for _, tx := range filteredTransactions {
		// if category denom is sent filter by category
		if q.Denom != nil {
			filterCategoryID, ok := mappedCategories[*q.Denom]
			if !ok {
				continue
			}
			story := ta.storyResolver(ctx, story.QueryStoryByIDParams{ID: tx.GroupID})
			if story.CategoryID != filterCategoryID {
				continue
			}
		}
		var vote stake.Vote
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
			vote = *backing.Vote
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
			vote = *challenge.Vote
		}
		argument := ta.argumentResolver(ctx, app.QueryByIDParams{ID: vote.ArgumentID})
		credArgument := CredArgument{
			ID:        argument.ID,
			StoryID:   argument.StoryID,
			Body:      argument.Body,
			Creator:   argument.Creator,
			Timestamp: argument.Timestamp,
			Vote:      vote.Vote,
			Amount:    vote.Amount,
		}
		credArguments = append(credArguments, credArgument)
	}
	return credArguments
}

func mapArgumentTransaction(tType trubank.TransactionType, storyID, stakeID int64) string {
	return fmt.Sprintf("%d/%d/%d", storyID, tType, stakeID)
}
