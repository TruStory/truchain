package truapi

import (
	"github.com/TruStory/truchain/x/argument"
	"github.com/TruStory/truchain/x/stake"
	"github.com/TruStory/truchain/x/trubank"
)

// LikedArgument ..
type LikedArgument struct {
	Stake       stake.Vote
	Argument    argument.Argument
	Transaction trubank.Transaction
}
