package truapi

import (
	"encoding/json"
	"net/http"

	"github.com/TruStory/truchain/x/db"
	"github.com/TruStory/truchain/x/truapi/cookies"
	"github.com/TruStory/truchain/x/truapi/render"
)

// DeviceTokenRegistrationRequest represents the JSON request of registeren a device token
// for push notifications.
type DeviceTokenRegistrationRequest struct {
	Address  string `json:"address"`
	Platform string `json:"platform"`
	Token    string `json:"token"`
}

// DeviceTokenUnregistration represents the JSON request to remove a device token.
type DeviceTokenUnregistration struct {
	Platform string `json:"platform"`
	Token    string `json:"token"`
}

// HandleDeviceTokenRegistration takes a `DeviceTokenRegistrationRequest` and returns a `DeviceToken`
func (ta *TruAPI) HandleDeviceTokenRegistration(w http.ResponseWriter, r *http.Request) {
	// check if request comes from an authenticated user.
	auth, err := cookies.GetAuthenticatedUser(r)
	if err != nil {
		render.Error(w, r, err.Error(), http.StatusBadRequest)
		return
	}
	request := &DeviceTokenRegistrationRequest{}
	err = json.NewDecoder(r.Body).Decode(request)
	if err != nil {
		render.Error(w, r, "bad payload", http.StatusBadRequest)
		return
	}

	// check if logged user matches the sent address
	if auth.Address != request.Address {
		render.Error(w, r, "invalid address", http.StatusBadRequest)
		return
	}
	deviceToken := &db.DeviceToken{
		Token:    request.Token,
		Address:  request.Address,
		Platform: request.Platform,
	}
	err = ta.DBClient.UpsertDeviceToken(deviceToken)
	if err == db.ErrInvalidAddress {
		render.Error(w, r, err.Error(), http.StatusNotFound)
		return
	}
	if err != nil {
		render.Error(w, r, err.Error(), http.StatusInternalServerError)
		return
	}
	render.Response(w, r, deviceToken, http.StatusOK)
}

// HandleUnregisterDeviceToken takes a `UnregisterDeviceTokenRequest`
func (ta *TruAPI) HandleUnregisterDeviceToken(w http.ResponseWriter, r *http.Request) {
	// check if request comes from an authenticated user.
	auth, err := cookies.GetAuthenticatedUser(r)
	if err != nil {
		render.Error(w, r, err.Error(), http.StatusBadRequest)
		return
	}
	request := &DeviceTokenUnregistration{}
	err = json.NewDecoder(r.Body).Decode(request)
	if err != nil {
		render.Error(w, r, "bad payload", http.StatusBadRequest)
		return
	}
	err = ta.DBClient.RemoveDeviceToken(auth.Address, request.Token, request.Platform)

	if err != nil {
		render.Error(w, r, err.Error(), http.StatusInternalServerError)
		return
	}
	render.Response(w, r, request, http.StatusOK)
}
