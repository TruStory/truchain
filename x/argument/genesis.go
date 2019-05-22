package argument

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GenesisState - all story state that must be provided at genesis
type GenesisState struct {
	Arguments []Argument `json:"arguments"`
	Likes     []Like     `json:"likes"`
	Params    Params     `json:"params"`
}

// DefaultGenesisState for tests
func DefaultGenesisState() GenesisState {
	return GenesisState{
		Params: DefaultParams(),
	}
}

// InitGenesis initializes arguments and likes from genesis file
func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) {
	for _, argument := range data.Arguments {
		keeper.setArgument(ctx, argument)
	}
	for _, like := range data.Likes {
		keeper.setLike(ctx, like)
	}
	keeper.SetLen(ctx, int64(len(data.Arguments)))
	keeper.SetParams(ctx, data.Params)
}

// ExportGenesis exports the genesis state
func ExportGenesis(ctx sdk.Context, keeper Keeper) GenesisState {
	var arguments []Argument
	var likes []Like
	prefix := "argument:id:"
	err := keeper.EachPrefix(ctx, prefix, func(bz []byte) bool {
		var arg Argument
		keeper.GetCodec().MustUnmarshalBinaryLengthPrefixed(bz, &arg)
		arguments = append(arguments, arg)
		likesForArgument, err := keeper.LikesByArgumentID(ctx, arg.ID)
		if err != nil {
			panic(err)
		}
		likes = append(likes, likesForArgument...)
		return true
	})
	if err != nil {
		panic(err)
	}

	params := keeper.GetParams(ctx)

	return GenesisState{
		Arguments: arguments,
		Likes:     likes,
		Params:    params,
	}
}
