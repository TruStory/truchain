package stake

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
)

const (
	// StoreKey is string representation of the store key
	StoreKey = "stake"
)

// Keeper data type storing keys to the key-value store
type Keeper struct {
	paramStore params.Subspace
}

// NewKeeper creates a new keeper with write and read access
func NewKeeper(paramStore params.Subspace) Keeper {
	return Keeper{
		paramStore.WithTypeTable(ParamTypeTable()),
	}
}

// ValidateArgument validates the length of an argument
func (k Keeper) ValidateArgument(ctx sdk.Context, argument string) sdk.Error {
	len := len([]rune(argument))
	minArgumentLength := k.GetParams(ctx).MinArgumentLength
	maxArgumentLength := k.GetParams(ctx).MaxArgumentLength

	if len > 0 && (len < minArgumentLength || len > maxArgumentLength) {
		return ErrInvalidArgumentMsg(argument)
	}

	return nil
}
