package truapi

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/TruStory/truchain/x/chttp"
	"github.com/TruStory/truchain/x/db"
)

// AddCommentRequest represents the JSON request for adding a comment
type AddCommentRequest struct {
	ParentID   int64  `json:"parent_id.omitonempty"`
	ArgumentID int64  `json:"argument_id"`
	Body       string `json:"body"`
	Creator    string `json:"creator"`
}

// HandleAddComment takes a `AddCommentRequest` and returns a 200 response
func (ta *TruAPI) HandleAddComment(r *http.Request) chttp.Response {
	request := &AddCommentRequest{}
	err := json.NewDecoder(r.Body).Decode(request)
	if err != nil {
		return chttp.SimpleErrorResponse(400, err)
	}

	user := r.Context().Value(userContextKey)
	if user == nil {
		return chttp.SimpleErrorResponse(401, err)
	}

	comment := &db.Comment{
		ParentID:   request.ParentID,
		ArgumentID: request.ArgumentID,
		Body:       request.Body,
		Creator:    request.Creator,
		CreatedAt:  time.Now(),
	}
	err = ta.DBClient.AddComment(comment)
	if err != nil {
		return chttp.SimpleErrorResponse(400, err)
	}

	return chttp.SimpleResponse(200, nil)
}
