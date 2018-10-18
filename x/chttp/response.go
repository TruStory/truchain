package chttp

import (
	"encoding/json"
)

type Response interface {
	HTTPCode() int
	Data() []byte
	Error() string
	Marshal() ([]byte, error)
}

type JsonResponse struct {
	status int
	Body   *json.RawMessage `json:"data"`
	Err    string           `json:"error"`
}

func NewResponse(status int, data json.RawMessage, err error) Response {
	errs := ""

	if err != nil {
		errs = err.Error()
	}

	return Response(JsonResponse{status: status, Body: &data, Err: errs})
}

func SimpleResponse(status int, data json.RawMessage) Response {
	return NewResponse(status, data, nil)
}

func SimpleDataResponse(status int, data map[string]interface{}) Response {
	bz, err := json.Marshal(data)

	if err != nil {
		panic(err)
	}

	return NewResponse(status, bz, nil)
}

func SimpleErrorResponse(status int, err error) Response {
	return NewResponse(status, []byte("{}"), err)
}

func (r JsonResponse) HTTPCode() int {
	if r.status != 0 {
		return r.status
	}

	if r.Err != "" {
		return 500
	}

	return 200
}

func (r JsonResponse) Data() []byte {
	return *r.Body
}

func (r JsonResponse) Error() string {
	return r.Err
}

func (r JsonResponse) Marshal() ([]byte, error) {
	return json.Marshal(r)
}
