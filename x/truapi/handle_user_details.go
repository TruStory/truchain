package truapi

import (
	"encoding/json"
	"net/http"

	"github.com/TruStory/truchain/x/chttp"
	"github.com/TruStory/truchain/x/cookies"
)

// UserResponse is a JSON response body representing the result of User
type UserResponse struct {
	Address string `json:"address"`
}

// HandleUserDetails takes a `UserRequest` and returns a `UserResponse`
func (ta *TruAPI) HandleUserDetails(r *http.Request) chttp.Response {

	// Get the user context
	truUser, err := cookies.GetUserFromCookie(r)
	if err == http.ErrNoCookie {
		return chttp.SimpleErrorResponse(401, err)
	}
	if err != nil {
		return chttp.SimpleErrorResponse(401, err)
	}

	responseBytes, _ := json.Marshal(UserResponse{
		Address: truUser["address"],
	})

	return chttp.SimpleResponse(200, responseBytes)
}
