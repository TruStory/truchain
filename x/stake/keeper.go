package stake

import "github.com/cosmos/cosmos-sdk/x/params"

const (
	// StoreKey is string representation of the store key
	StoreKey = "staking"
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
