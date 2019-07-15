package params

import (
	"testing"

	"github.com/TruStory/truchain/x/slashing"

	"github.com/stretchr/testify/assert"
	abci "github.com/tendermint/tendermint/abci/types"
)

const custom = "custom"

func TestQueryPath_Success(t *testing.T) {
	ctx, keeper := mockDB()

	query := abci.RequestQuery{
		Path: QueryPath,
	}

	querier := NewQuerier(keeper)
	resBytes, err := querier(ctx, []string{QueryPath}, query)
	assert.Nil(t, err)

	var returnedParams Params
	jsonErr := ModuleCodec.UnmarshalJSON(resBytes, &returnedParams)
	assert.Nil(t, jsonErr)
	assert.Equal(t, returnedParams.SlashingParams.MinSlashCount, slashing.DefaultParams().MinSlashCount)
}
