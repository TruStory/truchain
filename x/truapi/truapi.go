package truapi

import (
	"context"
	"net/http"

	"github.com/TruStory/truchain/x/chttp"
	xgraphql "github.com/TruStory/truchain/x/graphql"
	"github.com/TruStory/truchain/x/story"
	"github.com/samsarahq/thunder/graphql"
	"github.com/samsarahq/thunder/graphql/graphiql"
)

// TruAPI implements an HTTP server for TruStory functionality using `chttp.API`
type TruAPI struct {
	*chttp.API
	GraphQLClient *xgraphql.Client
}

// NewTruAPI returns a `TruAPI` instance populated with the existing app and a new GraphQL client
func NewTruAPI(aa *chttp.App) *TruAPI {
	ta := TruAPI{
		API:           chttp.NewAPI(aa, supported),
		GraphQLClient: xgraphql.NewGraphQLClient(),
	}

	return &ta
}

// RegisterRoutes applies the TruStory API routes to the `chttp.API` router
func (ta *TruAPI) RegisterRoutes() {
	ta.Use(chttp.JSONResponseMiddleware)
	http.Handle("/graphql", graphql.Handler(ta.GraphQLClient.Schema))
	http.Handle("/graphiql/", http.StripPrefix("/graphiql/", graphiql.Handler()))
	ta.HandleFunc("/graphql", ta.HandleGraphQL)
	ta.HandleFunc("/presigned", ta.HandlePresigned)
	ta.HandleFunc("/register", ta.HandleRegistration)
}

// RegisterResolvers builds the app's GraphQL schema from resolvers (declared in `resolver.go`)
func (ta *TruAPI) RegisterResolvers() {
	ta.GraphQLClient.RegisterQueryResolver("stories", ta.storyResolver)
	ta.GraphQLClient.BuildSchema()
}
