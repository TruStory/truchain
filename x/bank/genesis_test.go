package bank

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultGenesisState(t *testing.T) {
	state := DefaultGenesisState()
	assert.Len(t, state.Transactions, 0)
}

func TestExportGenesis(t *testing.T) {

}
