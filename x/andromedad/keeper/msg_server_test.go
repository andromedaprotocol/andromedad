package keeper_test

import (
	"context"
	"testing"

	keepertest "github.com/andromedaprotocol/andromedad/testutil/keeper"
	"github.com/andromedaprotocol/andromedad/x/andromedad/keeper"
	"github.com/andromedaprotocol/andromedad/x/andromedad/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func setupMsgServer(t testing.TB) (types.MsgServer, context.Context) {
	k, ctx := keepertest.AndromedadKeeper(t)
	return keeper.NewMsgServerImpl(*k), sdk.WrapSDKContext(ctx)
}
