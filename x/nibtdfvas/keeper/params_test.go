package keeper_test

import (
	"testing"

	testkeeper "github.com/andromedaprotocol/andromedad/testutil/keeper"
	// testkeeper "andromedad/testutil/keeper"


	"github.com/andromedaprotocol/andromedad/x/nibtdfvas/types"
	// "andromedad/x/nibtdfvas/types"

	"github.com/stretchr/testify/require"
)

func TestGetParams(t *testing.T) {
	k, ctx := testkeeper.NibtdfvasKeeper(t)
	params := types.DefaultParams()

	k.SetParams(ctx, params)

	require.EqualValues(t, params, k.GetParams(ctx))
}
