package truapi

import (
  "encoding/json"
  "io/ioutil"
  "net/http"

  "github.com/TruStory/truchain/x/chttp"
  "github.com/TruStory/truchain/x/graphql"
)

func (ta *TruApi) HandleGraphQL(r *http.Request) chttp.Response {
  gr := new(graphql.Request)
  jsonBytes, err := ioutil.ReadAll(r.Body)

  if err != nil {
    return chttp.SimpleErrorResponse(500, err)
  }

  err = json.Unmarshal(jsonBytes, &gr)

  if err != nil {
    return chttp.SimpleErrorResponse(400, err)
  }

  return ta.GraphQLClient.Query(r.Context(), *gr)
}
