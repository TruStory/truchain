package chttp

import (
	"encoding/json"
)

// Response describes the information that should be available in the HTTP response body to any API request
type Response interface {
	HTTPCode() int
	Data() []byte
	Error() string
	Marshal() ([]byte, error)
}

// JSONResponse is an implementation of Response which encodes the required data as JSON
type JSONResponse struct {
	status int
	Body   *json.RawMessage `json:"data"`
	Err    string           `json:"error"`
}

// NewResponse returns a JSONResponse with the given status/data/err
func NewResponse(status int, data json.RawMessage, err error) Response {
	errs := ""

	if err != nil {
		errs = err.Error()
	}

	return Response(JSONResponse{status: status, Body: &data, Err: errs})
}

// SimpleResponse is a helper to return a response with the given status and data (err is nil)
func SimpleResponse(status int, data json.RawMessage) Response {
	return NewResponse(status, data, nil)
}

// SimpleDataResponse is a helper to return a response whose data is an object
func SimpleDataResponse(status int, data map[string]interface{}) Response {
	bz, err := json.Marshal(data)

	if err != nil {
		panic(err)
	}

	return NewResponse(status, bz, nil)
}

// SimpleErrorResponse is a helper to return a response with an error (data is nil)
func SimpleErrorResponse(status int, err error) Response {
	return NewResponse(status, []byte("{}"), err)
}

// HTTPCode implements Response.HTTPCode
func (r JSONResponse) HTTPCode() int {
	if r.status != 0 {
		return r.status
	}

	if r.Err != "" {
		return 500
	}

	return 200
}

// Data implements Response.Data
func (r JSONResponse) Data() []byte {
	return *r.Body
}

// Error implements Response.Error
func (r JSONResponse) Error() string {
	return r.Err
}

// Marshal implements Response.Marshal
func (r JSONResponse) Marshal() ([]byte, error) {
	return json.Marshal(r)
}
