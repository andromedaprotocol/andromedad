package nibtdfvas_test

import (
	"testing"

	keepertest "github.com/andromedaprotocol/andromedad/testutil/keeper"
	"github.com/andromedaprotocol/andromedad/testutil/nullify"
	// keepertest "andromedad/testutil/keeper"
	// "andromedad/testutil/nullify"

	// "andromedad/x/nibtdfvas"
	// "andromedad/x/nibtdfvas/types"

	"github.com/andromedaprotocol/andromedad/x/nibtdfvas"
	"github.com/andromedaprotocol/andromedad/x/nibtdfvas/types"
	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),
	}

	k, ctx := keepertest.NibtdfvasKeeper(t)
	nibtdfvas.InitGenesis(ctx, *k, genesisState)
	got := nibtdfvas.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)
}
