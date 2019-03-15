package truapi

import (
	"encoding/json"
	"net/http"

	"github.com/TruStory/truchain/x/chttp"
)

// UserResponse is a JSON response body representing the result of User
type UserResponse struct {
	Address string `json:"address"`
}

// HandleUserDetails takes a `UserRequest` and returns a `UserResponse`
func (ta *TruAPI) HandleUserDetails(r *http.Request) chttp.Response {

	// Get the user context
	truUser, err := GetTruUserFromCookie(r)
	if err != nil {
		panic(err)
	}

	responseBytes, _ := json.Marshal(UserResponse{
		Address: truUser["address"],
	})

	return chttp.SimpleResponse(200, responseBytes)
}
