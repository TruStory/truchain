package params

import (
	"testing"

	"github.com/TruStory/truchain/x/slashing"

	"github.com/stretchr/testify/assert"
)

func TestParams_Success(t *testing.T) {
	ctx, keeper := mockDB()

	params := keeper.Params(ctx)

	assert.Equal(t, params.SlashingParams.MinSlashCount, slashing.DefaultParams().MinSlashCount)
}
