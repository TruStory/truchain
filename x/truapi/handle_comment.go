package truapi

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/TruStory/truchain/x/chttp"
	"github.com/TruStory/truchain/x/db"
	"github.com/TruStory/truchain/x/truapi/cookies"
)

// AddCommentRequest represents the JSON request for adding a comment
type AddCommentRequest struct {
	ParentID   int64  `json:"parent_id,omitonempty"`
	ArgumentID int64  `json:"argument_id"`
	Body       string `json:"body"`
	// Deprecated: Creator comes from cookie
	Creator string `json:"creator"`
}

// HandleComment handles requests for comments
func (ta *TruAPI) HandleComment(r *http.Request) chttp.Response {
	switch r.Method {
	case http.MethodPost:
		return ta.handleCreateComment(r)
	default:
		return chttp.SimpleErrorResponse(404, Err404ResourceNotFound)
	}
}

func (ta *TruAPI) handleCreateComment(r *http.Request) chttp.Response {
	request := &AddCommentRequest{}
	err := json.NewDecoder(r.Body).Decode(request)
	if err != nil {
		return chttp.SimpleErrorResponse(400, err)
	}

	user := r.Context().Value(userContextKey).(*cookies.AuthenticatedUser)
	if user == nil {
		return chttp.SimpleErrorResponse(401, Err401NotAuthenticated)
	}

	comment := &db.Comment{
		ParentID:   request.ParentID,
		ArgumentID: request.ArgumentID,
		Body:       request.Body,
		Creator:    user.Address,
	}
	err = ta.DBClient.AddComment(comment)
	if err != nil {
		return chttp.SimpleErrorResponse(500, err)
	}
	respBytes, err := json.Marshal(comment)
	if err != nil {
		return chttp.SimpleErrorResponse(500, err)
	}
	ta.sendCommentNotification(CommentNotificationRequest{
		ID:         comment.ID,
		ArgumentID: comment.ArgumentID,
		Creator:    comment.Creator,
		Timestamp:  time.Now(),
	})
	return chttp.SimpleResponse(200, respBytes)
}
