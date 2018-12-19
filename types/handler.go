package types

import (
	"encoding/binary"
	"encoding/json"
	"net/url"
	"reflect"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Utilities for all module handlers

// IDResult is a common result struct containing the id of a created object
type IDResult struct {
	ID int64 `json:"id"`
}

// ErrMsgHandler returns an unknown Msg request error result
func ErrMsgHandler(msg sdk.Msg) sdk.Result {
	if mType := reflect.TypeOf(msg); mType != nil {
		errMsg := "Unrecognized Msg type: " + mType.Name()
		return sdk.ErrUnknownRequest(errMsg).Result()
	}

	return sdk.ErrUnknownRequest("Unrecognized Msg").Result()
}

// ParseEvidence parses evidence url strings into real URLs
func ParseEvidence(urlStrings []string) (
	urls []url.URL, err sdk.Error) {
	for _, urlString := range urlStrings {

		evidenceURL, err := url.ParseRequestURI(urlString)
		if err != nil {
			return urls, ErrInvalidEvidenceURL(urlString)
		}

		urls = append(urls, *evidenceURL)
	}

	return urls, nil
}

// Result returns a successful handler result with id of the type
// encoded as binary data
func Result(id int64) sdk.Result {
	bz, err := json.Marshal(IDResult{ID: id})

	if err != nil {
		panic(err)
	}

	return sdk.Result{Data: bz}
}

// i2b converts an int64 into a byte array
func i2b(x int64) []byte {
	var b [binary.MaxVarintLen64]byte
	return b[:binary.PutVarint(b[:], x)]
}
