package nibftdfvas_test

import (
	"testing"

	keepertest "github.com/andromedaprotocol/andromedad/testutil/keeper"
	"github.com/andromedaprotocol/andromedad/testutil/nullify"
	"github.com/andromedaprotocol/andromedad/x/nibftdfvas"
	"github.com/andromedaprotocol/andromedad/x/nibftdfvas/types"
	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params:	types.DefaultParams(),
		
		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.NibftdfvasKeeper(t)
	nibftdfvas.InitGenesis(ctx, *k, genesisState)
	got := nibftdfvas.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	

	// this line is used by starport scaffolding # genesis/test/assert
}
