package keeper_test

import (
	"context"
	"testing"

	keepertest "andromedad/testutil/keeper"
	"andromedad/x/nibtdfvas/keeper"
	"andromedad/x/nibtdfvas/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func setupMsgServer(t testing.TB) (types.MsgServer, context.Context) {
	k, ctx := keepertest.NibtdfvasKeeper(t)
	return keeper.NewMsgServerImpl(*k), sdk.WrapSDKContext(ctx)
}
