package keeper_test

import (
	"context"
	"testing"

	keepertest "github.com/andromedaprotocol/andromedad/testutil/keeper"
	// keepertest "andromedad/testutil/keeper"


	"github.com/andromedaprotocol/andromedad/x/nibtdfvas/keeper"
	"github.com/andromedaprotocol/andromedad/x/nibtdfvas/types"

	// "andromedad/x/nibtdfvas/keeper"
	// "andromedad/x/nibtdfvas/types"
	
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func setupMsgServer(t testing.TB) (types.MsgServer, context.Context) {
	k, ctx := keepertest.NibtdfvasKeeper(t)
	return keeper.NewMsgServerImpl(k), sdk.WrapSDKContext(ctx)
}
