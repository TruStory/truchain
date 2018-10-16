package truapi

import (
  "github.com/TruStory/truchain/x/chttp"
  "github.com/TruStory/truchain/x/graphql"
)

type TruApi struct {
  *chttp.Api
  GraphQLClient *graphql.Client
}

func NewTruApi(aa *chttp.App) *TruApi {
  ta := TruApi{
    Api: chttp.NewApi(aa, supported), 
    GraphQLClient: graphql.NewGraphqlClient(),
  }

  return &ta
}

func (ta *TruApi) RegisterRoutes() {
  ta.Use(chttp.ContentJSONMiddleware)
  ta.HandleFunc("/graphql", ta.HandleGraphQL)
  ta.HandleFunc("/presigned", ta.HandlePresigned)
  ta.HandleFunc("/register", ta.HandleRegistration)
}

func (ta *TruApi) RegisterResolvers() {
  ta.GraphQLClient.RegisterQueryResolver("story", ta.storyResolver)
  ta.GraphQLClient.BuildSchema()
}
