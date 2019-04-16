package truapi

import (
	"encoding/json"
	"net/http"

	"github.com/TruStory/truchain/x/chttp"
	"github.com/TruStory/truchain/x/db"
)

// MarkNotificationAsReadRequest represents the JSON request for.. come on
type MarkNotificationAsReadRequest struct {
	NotificationID int64 `json:"notification_id"`
}

// HandleMarkNotificationAsRead takes a `MarkNotificationAsReadRequest` and returns a 200 response
func (ta *TruAPI) HandleMarkNotificationAsRead(r *http.Request) chttp.Response {
	request := &MarkNotificationAsReadRequest{}
	err := json.NewDecoder(r.Body).Decode(request)
	if err != nil {
		return chttp.SimpleErrorResponse(400, err)
	}

	notificationEvent := &db.NotificationEvent{ID: request.NotificationID}
	err = ta.DBClient.Find(notificationEvent)
	if err != nil {
		return chttp.SimpleErrorResponse(401, err)
	}

	notificationEvent.Read = true
	err = ta.DBClient.Update(notificationEvent)
	if err != nil {
		return chttp.SimpleErrorResponse(500, err)
	}

	return chttp.SimpleResponse(200, nil)
}
