package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"

	testkeeper "github.com/andromedaprotocol/andromedad/testutil/keeper"
	"github.com/andromedaprotocol/andromedad/x/feeburn/types"
)

func TestParamsQuery(t *testing.T) {
	keeper, ctx := testkeeper.FeeburnKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	params := types.DefaultParams()
	keeper.SetParams(ctx, params)

	response, err := keeper.Params(wctx, &types.QueryParamsRequest{})
	require.NoError(t, err)
	require.Equal(t, &types.QueryParamsResponse{Params: params}, response)
}
