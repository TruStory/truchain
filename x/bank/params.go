package bank

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
)

var (
	ParamKeyRewardBrokerAddress = []byte("rewardBrokerAddress")
)

type Params struct {
	RewardBrokerAddress sdk.AccAddress `json:"reward_broker_address"`
}

func DefaultParams() Params {
	return Params{RewardBrokerAddress: nil}
}

func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		{Key: ParamKeyRewardBrokerAddress, Value: &p.RewardBrokerAddress},
	}
}

// ParamKeyTable for bank module
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&Params{})
}

// GetParams gets the genesis params for the bank module
func (k Keeper) GetParams(ctx sdk.Context) Params {
	var paramSet Params
	k.paramStore.GetParamSet(ctx, &paramSet)
	return paramSet
}

// SetParams sets the params for bank module
func (k Keeper) SetParams(ctx sdk.Context, params Params) {
	k.paramStore.SetParamSet(ctx, &params)
}
