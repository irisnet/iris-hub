package simulation

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/irisnet/irishub/baseapp"
	"github.com/irisnet/irishub/simulation/mock/simulation"
)

// AllInvariants tests all slashing invariants
func AllInvariants() simulation.Invariant {
	return func(t *testing.T, app *baseapp.BaseApp, log string) {
		// TODO Any invariants to check here?
		require.Nil(t, nil)
	}
}
