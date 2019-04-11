package query

import (
	"github.com/TruStory/truchain/x/argument"
	"github.com/TruStory/truchain/x/trubank"
)

type LikedArgument struct {
	Argument    argument.Argument
	Transaction trubank.Transaction
}
