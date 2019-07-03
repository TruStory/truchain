package exported

import (
	"testing"

	"github.com/magiconair/properties/assert"
)

func Test_getFilters(t *testing.T) {
	filters := GetFilters(
		FilterByTransactionType(TransactionBacking, TransactionChallenge, TransactionUpvote),
		SortOrder(SortDesc),
		Limit(10),
		Offset(100),
	)
	assert.Equal(t, filters.TransactionTypes, []TransactionType{
		TransactionBacking,
		TransactionChallenge,
		TransactionUpvote})
	assert.Equal(t, filters.SortOrder, SortDesc)
	assert.Equal(t, filters.Limit, 10)
	assert.Equal(t, filters.Offset, 100)

	filters = GetFilters(
		Limit(-100),
		Offset(-1),
	)
	assert.Equal(t, filters.Limit, 0)
	assert.Equal(t, filters.Offset, 0)
}
