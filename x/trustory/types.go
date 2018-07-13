package trustory

type Story struct {
	Body        string      `json:"body"`
	Creator     sdk.Address `json:"creator"`
	SubmitBlock int64       `json:"submit_block`
	State       string      `json:"state"`
	YesVotes    int64       `json:"yes_votes`
	NoVotes     int64       `json:"no_votes"`
}

// NewStory creates a new story
func NewStory(body string, creator sdk.Address, blockHeight int64) Story {
	return Story{
		Body:        body,
		Creator:     creator,
		SubmitBlock: blockHeight,
		State:       "Created",
		YesVotes:    0,
		NoVotes:     0,
	}
}
