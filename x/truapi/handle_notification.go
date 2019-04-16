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

// HandleUpdateNotificationEvent takes a `MarkNotificationAsReadRequest` and returns a 200 response
func (ta *TruAPI) HandleUpdateNotificationEvent(r *http.Request) chttp.Response {
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
