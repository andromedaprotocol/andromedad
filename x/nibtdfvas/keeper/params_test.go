package keeper_test

import (
	"testing"

	testkeeper "andromedad/testutil/keeper"
	"andromedad/x/nibtdfvas/types"
	"github.com/stretchr/testify/require"
)

func TestGetParams(t *testing.T) {
	k, ctx := testkeeper.NibtdfvasKeeper(t)
	params := types.DefaultParams()

	k.SetParams(ctx, params)

	require.EqualValues(t, params, k.GetParams(ctx))
}
