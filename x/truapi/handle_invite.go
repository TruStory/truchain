package truapi

import (
	"encoding/json"
	"errors"
	"net/http"
	"regexp"

	"github.com/TruStory/truchain/x/chttp"
	"github.com/TruStory/truchain/x/db"
	"github.com/TruStory/truchain/x/truapi/cookies"
)

// AddInviteRequest represents the JSON request for adding an invite
type AddInviteRequest struct {
	TwitterUsername string `json:"twitter_username,omitempty"`
	Email           string `json:"email"`
}

// HandleInvite handles requests for invites
func (ta *TruAPI) HandleInvite(r *http.Request) chttp.Response {
	switch r.Method {
	case http.MethodPost:
		return ta.handleCreateInvite(r)
	default:
		return chttp.SimpleErrorResponse(404, Err404ResourceNotFound)
	}
}

func (ta *TruAPI) handleCreateInvite(r *http.Request) chttp.Response {
	request := &AddInviteRequest{}
	err := json.NewDecoder(r.Body).Decode(request)
	if err != nil {
		return chttp.SimpleErrorResponse(400, err)
	}
	// check if valid email address
	re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	if !re.MatchString(request.Email) {
		return chttp.SimpleErrorResponse(422, errors.New("Invalid email address"))
	}

	user := r.Context().Value(userContextKey).(*cookies.AuthenticatedUser)

	if user == nil {
		return chttp.SimpleErrorResponse(401, Err401NotAuthenticated)
	}

	twitterProfile, err := ta.DBClient.TwitterProfileByUsername(request.TwitterUsername)
	if err != nil {
		return chttp.SimpleErrorResponse(500, err)
	}
	if twitterProfile != nil {
		return chttp.SimpleErrorResponse(422, errors.New("This user has already registered"))
	}

	invite := &db.Invite{
		Creator:               user.Address,
		FriendTwitterUsername: request.TwitterUsername,
		FriendEmail:           request.Email,
	}
	err = ta.DBClient.AddInvite(invite)
	// TODO: error on duplicate entry should return unique error code
	if err != nil {
		return chttp.SimpleErrorResponse(500, err)
	}
	if invite.ID == 0 {
		return chttp.SimpleErrorResponse(422, errors.New("This user has already been invited"))
	}
	respBytes, err := json.Marshal(invite)
	if err != nil {
		return chttp.SimpleErrorResponse(500, err)
	}
	return chttp.SimpleResponse(200, respBytes)
}
