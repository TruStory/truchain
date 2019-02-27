package stake

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Msg defines data common to backing, challenge, and
// token vote messages.
type Msg struct {
	StoryID  int64          `json:"story_id"`
	Amount   sdk.Coin       `json:"amount"`
	Argument string         `json:"argument,omitempty"`
	Creator  sdk.AccAddress `json:"creator"`
}
