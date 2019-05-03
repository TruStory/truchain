package truapi

import (
	"encoding/json"
	"net/http"

	"github.com/TruStory/truchain/x/truapi/render"
)

// TranslateMentionsRequest represents the JSON request for translating mentions.
type TranslateMentionsRequest struct {
	Body string `json:"body"`
}

// HandleTranslateCosmosMentions returns a string body with cosmos addresses mentions.
func (ta *TruAPI) HandleTranslateCosmosMentions(w http.ResponseWriter, r *http.Request) {
	request := &TranslateMentionsRequest{}
	err := json.NewDecoder(r.Body).Decode(request)
	if err != nil {
		render.Error(w, r, err.Error(), http.StatusBadRequest)
		return
	}
	b, err := ta.DBClient.TranslateToCosmosMentions(request.Body)
	if err != nil {
		render.Error(w, r, err.Error(), http.StatusInternalServerError)
		return
	}
	request.Body = b
	render.Response(w, r, request, http.StatusOK)
}
