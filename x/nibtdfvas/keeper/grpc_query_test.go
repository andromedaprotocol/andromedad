package keeper_test

import (
	"testing"

	testkeeper "github.com/andromedaprotocol/andromedad/testutil/keeper"
	// testkeeper "andromedad/testutil/keeper"


	"github.com/andromedaprotocol/andromedad/x/nibtdfvas/types"
	// "andromedad/x/nibtdfvas/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestParamsQuery(t *testing.T) {
	keeper, ctx := testkeeper.NibtdfvasKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	params := types.DefaultParams()
	keeper.SetParams(ctx, params)

	response, err := keeper.Params(wctx, &types.QueryParamsRequest{})
	require.NoError(t, err)
	require.Equal(t, &types.QueryParamsResponse{Params: params}, response)
}
