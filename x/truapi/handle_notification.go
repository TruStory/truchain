package truapi

import (
	"encoding/json"
	"net/http"

	"github.com/go-pg/pg"

	"github.com/TruStory/truchain/x/chttp"
	"github.com/TruStory/truchain/x/db"
	"github.com/TruStory/truchain/x/truapi/cookies"
)

// UpdateNotificationEventRequest represents the JSON request
type UpdateNotificationEventRequest struct {
	NotificationID int64 `json:"notification_id"`
	Read           *bool `json:"read,omitempty"`
}

// HandleNotificationEvent takes a `UpdateNotificationEventRequest` and returns a 200 response
func (ta *TruAPI) HandleNotificationEvent(r *http.Request) chttp.Response {
	switch r.Method {
	case http.MethodPut:
		return ta.handleUpdateNotificationEvent(r)
	default:
		return chttp.SimpleErrorResponse(404, Err404ResourceNotFound)
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
	if request.Read == nil {
		return chttp.SimpleErrorResponse(400, Err400MissingParameter)
	}

	// if request was made to mark all notification as read
	if request.NotificationID == -1 && *request.Read == true {
		return markAllAsRead(ta, r)
	}

	notificationEvent := &db.NotificationEvent{ID: request.NotificationID}
	err = ta.DBClient.Find(notificationEvent)
	if err == pg.ErrNoRows {
		return chttp.SimpleErrorResponse(404, Err404ResourceNotFound)
	}
	if err != nil {
		return chttp.SimpleErrorResponse(401, err)
	}

	notificationEvent.Read = *request.Read
	err = ta.DBClient.UpdateModel(notificationEvent)
	if err != nil {
		return chttp.SimpleErrorResponse(500, err)
	}

	return chttp.SimpleResponse(200, nil)
}

func markAllAsRead(ta *TruAPI, r *http.Request) chttp.Response {
	user, err := cookies.GetAuthenticatedUser(r)
	if err != nil {
		return chttp.SimpleErrorResponse(401, Err401NotAuthenticated)
	}

	err = ta.DBClient.MarkAllNotificationEventsAsReadByAddress(user.Address)
	if err != nil {
		return chttp.SimpleErrorResponse(500, Err500InternalServerError)
	}

	return chttp.SimpleResponse(200, nil)
}
