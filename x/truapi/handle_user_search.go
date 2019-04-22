package truapi

import (
	"encoding/json"
	"net/http"

	"github.com/TruStory/truchain/x/chttp"
)

// UsernameSearchResponse is a JSON response body representing the result of username search
type UsernameSearchResponse struct {
	Usernames []string `json:"usernames"`
}

// HandleUsernameSearch takes a `UsernameSearchRequest` and returns a `UsernameSearchResponse`
func (ta *TruAPI) HandleUsernameSearch(r *http.Request) chttp.Response {
	switch r.Method {
	case http.MethodGet:
		return ta.handleUsernameSearch(r)
	default:
		return chttp.SimpleErrorResponse(401, Err404ResourceNotFound)
	}
}

func (ta *TruAPI) handleUsernameSearch(r *http.Request) chttp.Response {
	err := r.ParseForm()
	if err != nil {
		return chttp.SimpleErrorResponse(500, Err400MissingParameter)
	}
	prefix := r.Form["username_prefix"][0]

	usernames, err := ta.DBClient.UsernamesByPrefix(prefix)
	if err != nil {
		return chttp.SimpleErrorResponse(500, err)
	}
	if len(usernames) == 0 {
		return chttp.SimpleErrorResponse(404, Err404ResourceNotFound)
	}

	responseBytes, _ := json.Marshal(UsernameSearchResponse{
		Usernames: usernames,
	})

	return chttp.SimpleResponse(200, responseBytes)
}
