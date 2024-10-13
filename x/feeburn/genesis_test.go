package feeburn_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	keepertest "github.com/andromedaprotocol/andromedad/testutil/keeper"
	"github.com/andromedaprotocol/andromedad/testutil/nullify"
	"github.com/andromedaprotocol/andromedad/x/feeburn"
	"github.com/andromedaprotocol/andromedad/x/feeburn/types"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),
	}

	k, ctx := keepertest.FeeburnKeeper(t)
	feeburn.InitGenesis(ctx, *k, genesisState)
	got := feeburn.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)
}
