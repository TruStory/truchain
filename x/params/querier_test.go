package params

import (
	"strings"
	"testing"

	"github.com/TruStory/truchain/x/account"
	"github.com/TruStory/truchain/x/slashing"

	"github.com/stretchr/testify/assert"
	abci "github.com/tendermint/tendermint/abci/types"
)

const custom = "custom"

func TestQueryPath_Success(t *testing.T) {
	ctx, keeper := mockDB()

	query := abci.RequestQuery{
		Path: strings.Join([]string{custom, QueryPath}, "/"),
	}

	querier := NewQuerier(keeper)
	resBytes, err := querier(ctx, []string{QueryPath}, query)
	assert.Nil(t, err)

	var returnedParams Params
	jsonErr := ModuleCodec.UnmarshalJSON(resBytes, &returnedParams)
	assert.Nil(t, jsonErr)
	assert.Equal(t, returnedParams.AccountParams.MaxSlashCount, account.DefaultParams().MaxSlashCount)
	assert.Equal(t, returnedParams.SlashingParams.MinSlashCount, slashing.DefaultParams().MinSlashCount)
}
