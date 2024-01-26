package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	keepertest "andromedad/testutil/keeper"
	"andromedad/x/nibtdfvas/types"
)

func TestGetParams(t *testing.T) {
	k, ctx := keepertest.NibtdfvasKeeper(t)
	params := types.DefaultParams()

	require.NoError(t, k.SetParams(ctx, params))
	require.EqualValues(t, params, k.GetParams(ctx))
}
