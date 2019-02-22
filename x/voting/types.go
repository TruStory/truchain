package voting

import (
	"fmt"

	app "github.com/TruStory/truchain/types"
)

type poll struct {
	trueVotes  []app.Voter
	falseVotes []app.Voter
}

func (p poll) String() string {
	return fmt.Sprintf(
		"Poll results:\n True votes: %v\n False votes: %v",
		p.trueVotes, p.falseVotes)
}
