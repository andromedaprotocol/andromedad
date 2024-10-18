package keeper

import (
	"testing"

	"github.com/stretchr/testify/require"

	"cosmossdk.io/log"
	"cosmossdk.io/store"
	storemetrics "cosmossdk.io/store/metrics"
	storetypes "cosmossdk.io/store/types"

	tmdb "github.com/cosmos/cosmos-db"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"

	"github.com/andromedaprotocol/andromedad/x/feeburn/keeper"
	"github.com/andromedaprotocol/andromedad/x/feeburn/types"
)

func FeeburnKeeper(t testing.TB) (*keeper.Keeper, sdk.Context) {
	storeKey := storetypes.NewKVStoreKey(types.StoreKey)
	memStoreKey := storetypes.NewMemoryStoreKey(types.MemStoreKey)

	logger := log.NewNopLogger()

	db := tmdb.NewMemDB()
	stateStore := store.NewCommitMultiStore(db, logger, storemetrics.NewNoOpMetrics())
	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(memStoreKey, storetypes.StoreTypeMemory, nil)
	require.NoError(t, stateStore.LoadLatestVersion())

	registry := codectypes.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(registry)

	k := keeper.NewKeeper(
		cdc,
		storeKey,
		memStoreKey,
		authtypes.NewModuleAddress(govtypes.ModuleName),
	)

	ctx := sdk.NewContext(stateStore, cmtproto.Header{}, false, log.NewNopLogger())

	// Initialize params
	_ = k.SetParams(ctx, types.DefaultParams())

	return k, ctx
}
