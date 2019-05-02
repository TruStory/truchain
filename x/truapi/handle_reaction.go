package truapi

import (
	"encoding/json"
	"net/http"

	"github.com/TruStory/truchain/x/db"
	"github.com/TruStory/truchain/x/truapi/cookies"

	"github.com/TruStory/truchain/x/chttp"
)

// ReactionRequest represents the http request for a reaction
type ReactionRequest struct {
	ReactionType     db.ReactionType     `json:"reaction_type"`
	ReactionableType db.ReactionableType `json:"reactionable_type"`
	ReactionableID   int64               `json:"reactionable_id"`
}

// UnreactionRequest represents the http request to unreact an already created reaction
type UnreactionRequest struct {
	ID int64 `json:"id"`
}

// HandleReaction handles the mutations and queries about reactions
func (ta *TruAPI) HandleReaction(r *http.Request) chttp.Response {
	switch r.Method {
	case http.MethodPost:
		return ta.createReaction(r)
	case http.MethodDelete:
		return ta.deleteReaction(r)
	default:
		return chttp.SimpleErrorResponse(401, Err404ResourceNotFound)
	}
}

func (ta *TruAPI) createReaction(r *http.Request) chttp.Response {
	creator := r.Context().Value(userContextKey)
	if creator == nil {
		return chttp.SimpleErrorResponse(401, Err401NotAuthenticated)
	}

	request := &ReactionRequest{}
	err := json.NewDecoder(r.Body).Decode(request)
	if err != nil {
		return chttp.SimpleErrorResponse(400, Err400MissingParameter)
	}

	rxnable := db.Reactionable{
		Type: request.ReactionableType,
		ID:   request.ReactionableID,
	}
	err = ta.DBClient.ReactOnReactionable(
		creator.(*cookies.AuthenticatedUser).Address,
		request.ReactionType,
		rxnable,
	)
	if err != nil {
		return chttp.SimpleErrorResponse(500, err)
	}

	return chttp.SimpleResponse(200, nil)
}

func (ta *TruAPI) deleteReaction(r *http.Request) chttp.Response {
	creator := r.Context().Value(userContextKey)
	if creator == nil {
		return chttp.SimpleErrorResponse(401, Err401NotAuthenticated)
	}

	request := &UnreactionRequest{}
	err := json.NewDecoder(r.Body).Decode(request)
	if err != nil {
		return chttp.SimpleErrorResponse(400, Err400MissingParameter)
	}

	err = ta.DBClient.UnreactByAddressAndID(
		creator.(*cookies.AuthenticatedUser).Address,
		request.ID,
	)
	if err != nil {
		return chttp.SimpleErrorResponse(500, Err500InternalServerError)
	}

	return chttp.SimpleResponse(200, nil)
}
