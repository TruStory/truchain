package truapi

import (
	"encoding/json"
	"net/http"

	"github.com/TruStory/truchain/x/chttp"
	"github.com/TruStory/truchain/x/truapi/cookies"
)

// UserResponse is a JSON response body representing the result of User
type UserResponse struct {
	UserID   int64  `json:"userId"`
	Username string `json:"username"`
	Fullname string `json:"fullname"`
	Address  string `json:"address"`
}

// HandleUserDetails takes a `UserRequest` and returns a `UserResponse`
func (ta *TruAPI) HandleUserDetails(r *http.Request) chttp.Response {
	user, err := cookies.GetAuthenticatedUser(r)
	if err == http.ErrNoCookie {
		return chttp.SimpleErrorResponse(401, err)
	}
	if err != nil {
		return chttp.SimpleErrorResponse(401, err)
	}

	twitterProfile, err := ta.DBClient.TwitterProfileByID(user.TwitterProfileID)
	if err != nil {
		return chttp.SimpleErrorResponse(401, err)
	}

	// Chain was restarted and DB was wiped so Address and TwitterProfileID contained in cookie is stale.
	if twitterProfile.ID == 0 {
		return chttp.SimpleErrorResponse(401, err)
	}

	responseBytes, _ := json.Marshal(UserResponse{
		UserID:   twitterProfile.ID,
		Fullname: twitterProfile.FullName,
		Username: twitterProfile.Username,
		Address:  twitterProfile.Address,
	})

	return chttp.SimpleResponse(200, responseBytes)
}
