package params

import (
	"github.com/TruStory/truchain/x/community"
	"testing"


	"github.com/stretchr/testify/assert"
)

func TestParams_Success(t *testing.T) {
	ctx, keeper := mockDB()

	params := keeper.Params(ctx)
	t.Log(params)

	assert.Equal(t, params.CommunityParams.MinNameLength, community.DefaultParams().MinNameLength)
}
