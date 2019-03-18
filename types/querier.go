package types

import (
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// These are helpers used by Queriers

// QueryByIDParams is query params for any ID
type QueryByIDParams struct {
	ID int64
}

// QueryByStoryIDAndCreatorParams is query params for backing,
// challenge, and token votes by story id and creator
type QueryByStoryIDAndCreatorParams struct {
	StoryID int64
	Creator string
}

// QueryByCreatorParams returns the query params for getting any query by the creator
type QueryByCreatorParams struct {
	Creator string
}

// UnmarshalQueryParams unmarshals the request query from a client
func UnmarshalQueryParams(req abci.RequestQuery, params interface{}) (sdkErr sdk.Error) {
	parseErr := json.Unmarshal(req.Data, params)
	if parseErr != nil {
		sdkErr = sdk.ErrUnknownRequest(fmt.Sprintf("Incorrectly formatted request data - %s", parseErr.Error()))
		return
	}
	return
}

// MustMarshal marshals a struct into JSON bytes
func MustMarshal(v interface{}) (res []byte) {
	res, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		panic("Could not marshal result to JSON")
	}
	return
}
