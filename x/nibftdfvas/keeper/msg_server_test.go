package keeper_test

import (
	"context"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
    "github.com/andromedaprotocol/andromedad/x/nibftdfvas/types"
    "github.com/andromedaprotocol/andromedad/x/nibftdfvas/keeper"
    keepertest "github.com/andromedaprotocol/andromedad/testutil/keeper"
)

func setupMsgServer(t testing.TB) (types.MsgServer, context.Context) {
	k, ctx := keepertest.NibftdfvasKeeper(t)
	return keeper.NewMsgServerImpl(*k), sdk.WrapSDKContext(ctx)
}
