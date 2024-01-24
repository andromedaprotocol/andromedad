package nibtdfvas_test

import (
	"testing"

	keepertest "andromedad/testutil/keeper"
	"andromedad/testutil/nullify"
	"andromedad/x/nibtdfvas"
	"andromedad/x/nibtdfvas/types"
	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.NibtdfvasKeeper(t)
	nibtdfvas.InitGenesis(ctx, *k, genesisState)
	got := nibtdfvas.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	// this line is used by starport scaffolding # genesis/test/assert
}
