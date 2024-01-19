package keeper_test

import (
	"testing"

	testkeeper "github.com/andromedaprotocol/andromedad/testutil/keeper"
	"github.com/andromedaprotocol/andromedad/x/andromedad/types"
	"github.com/stretchr/testify/require"
)

func TestGetParams(t *testing.T) {
	k, ctx := testkeeper.AndromedadKeeper(t)
	params := types.DefaultParams()

	k.SetParams(ctx, params)

	require.EqualValues(t, params, k.GetParams(ctx))
}
