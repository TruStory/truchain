package stake

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
)

// store keys for voting params
var (
	KeyMinArgumentLength = []byte("minArgumentLength")
	KeyMaxArgumentLength = []byte("maxArgumentLength")
	KeyMinInterestRate   = []byte("minInterestRate")
	KeyMaxInterestRate   = []byte("maxInterestRate")
	KeyAmountWeight      = []byte("amountWeight")
	KeyPeriodWeight      = []byte("periodWeight")
)

// Params holds parameters for voting
type Params struct {
	MinArgumentLength int     `json:"min_argument_length"`
	MaxArgumentLength int     `json:"max_argument_length"`
	MinInterestRate   sdk.Dec `json:"min_interest_rate"`
	MaxInterestRate   sdk.Dec `json:"max_interest_rate"`
	AmountWeight      sdk.Dec `json:"amount_weight"`
	PeriodWeight      sdk.Dec `json:"period_weight"`
}

// DefaultParams is the default parameters for voting
func DefaultParams() Params {
	return Params{
		MinArgumentLength: 10,
		MaxArgumentLength: 1000,
		AmountWeight:      sdk.NewDecWithPrec(333, 3), // 33.3%
		PeriodWeight:      sdk.NewDecWithPrec(667, 3), // 66.7%
		MinInterestRate:   sdk.ZeroDec(),              // 0%
		MaxInterestRate:   sdk.NewDecWithPrec(10, 2),  // 10%
	}
}

// KeyValuePairs implements params.ParamSet
func (p *Params) KeyValuePairs() params.KeyValuePairs {
	return params.KeyValuePairs{
		{Key: KeyMinArgumentLength, Value: &p.MinArgumentLength},
		{Key: KeyMaxArgumentLength, Value: &p.MaxArgumentLength},
		{Key: KeyMinInterestRate, Value: &p.MinInterestRate},
		{Key: KeyMaxInterestRate, Value: &p.MaxInterestRate},
		{Key: KeyAmountWeight, Value: &p.AmountWeight},
		{Key: KeyPeriodWeight, Value: &p.PeriodWeight},
	}
}

// ParamTypeTable for story module
func ParamTypeTable() params.TypeTable {
	return params.NewTypeTable().RegisterParamSet(&Params{})
}

// GetParams gets the genesis params for the type
func (k Keeper) GetParams(ctx sdk.Context) Params {
	var paramSet Params
	k.paramStore.GetParamSet(ctx, &paramSet)
	return paramSet
}

// SetParams sets the params for the module
func (k Keeper) SetParams(ctx sdk.Context, params Params) {
	logger := ctx.Logger().With("module", "staking")
	k.paramStore.SetParamSet(ctx, &params)
	logger.Info(fmt.Sprintf("Loaded staking module params: %+v", params))
}

// Msg defines data common to backing, challenge, and
// token vote messages.
type Msg struct {
	StoryID  int64          `json:"story_id"`
	Amount   sdk.Coin       `json:"amount"`
	Argument string         `json:"argument,omitempty"`
	Creator  sdk.AccAddress `json:"creator"`
}
