package truapi

import (
	"context"
	"net/http"

	"github.com/TruStory/truchain/x/category"
	"github.com/TruStory/truchain/x/chttp"
	"github.com/TruStory/truchain/x/game"
	"github.com/TruStory/truchain/x/graphql"
	"github.com/TruStory/truchain/x/story"
	"github.com/TruStory/truchain/x/users"
	sdk "github.com/cosmos/cosmos-sdk/types"
	thunder "github.com/samsarahq/thunder/graphql"
	"github.com/samsarahq/thunder/graphql/graphiql"
)

// TruAPI implements an HTTP server for TruStory functionality using `chttp.API`
type TruAPI struct {
	*chttp.API
	GraphQLClient *graphql.Client
}

// NewTruAPI returns a `TruAPI` instance populated with the existing app and a new GraphQL client
func NewTruAPI(aa *chttp.App) *TruAPI {
	ta := TruAPI{
		API:           chttp.NewAPI(aa, supported),
		GraphQLClient: graphql.NewGraphQLClient(),
	}

	return &ta
}

// RegisterRoutes applies the TruStory API routes to the `chttp.API` router
func (ta *TruAPI) RegisterRoutes() {
	ta.Use(chttp.JSONResponseMiddleware)
	http.Handle("/graphql", thunder.Handler(ta.GraphQLClient.Schema))
	http.Handle("/graphiql/", http.StripPrefix("/graphiql/", graphiql.Handler()))
	ta.HandleFunc("/graphql", ta.HandleGraphQL)
	ta.HandleFunc("/presigned", ta.HandlePresigned)
	ta.HandleFunc("/register", ta.HandleRegistration)
}

// RegisterResolvers builds the app's GraphQL schema from resolvers (declared in `resolver.go`)
func (ta *TruAPI) RegisterResolvers() {
	getUser := func(ctx context.Context, addr sdk.AccAddress) users.User {
		res := ta.usersResolver(ctx, users.QueryUsersByAddressesParams{Addresses: []string{addr.String()}})
		if len(res) > 0 {
			return res[0]
		}
		return users.User{}
	}

	ta.GraphQLClient.RegisterObjectResolver("Argument", story.Argument{}, map[string]interface{}{
		"creator": func(ctx context.Context, q story.Argument) users.User { return getUser(ctx, q.Creator) },
	})

	ta.GraphQLClient.RegisterQueryResolver("categories", ta.allCategoriesResolver)
	ta.GraphQLClient.RegisterQueryResolver("category", ta.categoryResolver)
	ta.GraphQLClient.RegisterObjectResolver("Category", category.Category{}, map[string]interface{}{
		"id":      func(_ context.Context, q category.Category) int64 { return q.ID },
		"stories": ta.categoryStoriesResolver,
		"creator": func(ctx context.Context, q category.Category) users.User { return getUser(ctx, q.Creator) },
	})

	ta.GraphQLClient.RegisterObjectResolver("Evidence", story.Evidence{}, map[string]interface{}{
		"creator": func(ctx context.Context, q story.Evidence) users.User { return getUser(ctx, q.Creator) },
		"url":     func(ctx context.Context, q story.Evidence) string { return q.URL.String() },
	})

	ta.GraphQLClient.RegisterObjectResolver("Game", game.Game{}, map[string]interface{}{
		"id":              func(_ context.Context, q game.Game) int64 { return q.ID },
		"creator":         func(ctx context.Context, q game.Game) users.User { return getUser(ctx, q.Creator) },
		"challengeDenom":  func(_ context.Context, q game.Game) string { return q.ChallengePool.Denom },
		"challengeAmount": func(_ context.Context, q game.Game) string { return q.ChallengePool.Amount.String() },
	})

	ta.GraphQLClient.RegisterQueryResolver("story", ta.storyResolver)
	ta.GraphQLClient.RegisterObjectResolver("Story", story.Story{}, map[string]interface{}{
		"id":       func(_ context.Context, q story.Story) int64 { return q.ID },
		"category": ta.storyCategoryResolver,
		"creator":  func(ctx context.Context, q story.Story) users.User { return getUser(ctx, q.Creator) },
		"source":   func(ctx context.Context, q story.Story) string { return q.Source.String() },
		"argument": func(ctx context.Context, q story.Story) []story.Argument { return q.Arguments },
		"evidence": func(ctx context.Context, q story.Story) []story.Evidence { return q.Evidence },
		"game":     ta.gameResolver,
	})

	ta.GraphQLClient.RegisterQueryResolver("users", ta.usersResolver)
	ta.GraphQLClient.RegisterObjectResolver("User", users.User{}, map[string]interface{}{
		"id":             func(_ context.Context, q users.User) string { return q.Address },
		"pubkey":         func(_ context.Context, q users.User) string { return q.Pubkey.String() },
		"twitterProfile": ta.twitterProfileResolver,
	})

	ta.GraphQLClient.RegisterObjectResolver("TwitterProfile", users.TwitterProfile{}, map[string]interface{}{
		"id": func(_ context.Context, q users.TwitterProfile) string { return q.ID },
	})

	ta.GraphQLClient.RegisterObjectResolver("Coin", sdk.Coin{}, map[string]interface{}{
		"amount": func(_ context.Context, q sdk.Coin) int64 { return q.Amount.Int64() },
	})

	ta.GraphQLClient.BuildSchema()
}
