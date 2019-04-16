package truapi

import (
	"encoding/json"
	"net/http"

	"github.com/TruStory/truchain/x/chttp"
	"github.com/TruStory/truchain/x/db"
)

// UpdateNotificationEventRequest represents the JSON request for.. come on
type UpdateNotificationEventRequest struct {
	NotificationID int64 `json:"notification_id"`
	Read           bool  `json:"read"`
}

// HandleNotificationEvent takes a `UpdateNotificationEventRequest` and returns a 200 response
func (ta *TruAPI) HandleNotificationEvent(r *http.Request) chttp.Response {
	switch r.Method {
	case http.MethodPut:
		return ta.handleUpdateNotificationEvent(r)
	default:
		return chttp.SimpleErrorResponse(404, Err404)
	}
}

func (ta *TruAPI) handleUpdateNotificationEvent(r *http.Request) chttp.Response {
	// check if we have a user before doing anything
	user := r.Context().Value(userContextKey)
	if user == nil {
		return chttp.SimpleErrorResponse(401, Err401NotAuthenticated)
	}

	request := &UpdateNotificationEventRequest{}
	err := json.NewDecoder(r.Body).Decode(request)
	if err != nil {
		return chttp.SimpleErrorResponse(400, err)
	}

	notificationEvent := &db.NotificationEvent{ID: request.NotificationID}
	err = ta.DBClient.Find(notificationEvent)
	if err != nil {
		return chttp.SimpleErrorResponse(401, err)
	}

	notificationEvent.Read = request.Read
	err = ta.DBClient.Update(notificationEvent)
	if err != nil {
		return chttp.SimpleErrorResponse(500, err)
	}

	return chttp.SimpleResponse(200, nil)
}
