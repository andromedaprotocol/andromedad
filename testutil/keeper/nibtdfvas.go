package keeper

import (
	"testing"
	// "fmt"

	"github.com/andromedaprotocol/andromedad/x/nibtdfvas/keeper"
	"github.com/andromedaprotocol/andromedad/x/nibtdfvas/types"
	// "andromedad/x/nibtdfvas/keeper"
	// "andromedad/x/nibtdfvas/types"

	tmdb "github.com/cometbft/cometbft-db"
	"github.com/cometbft/cometbft/libs/log"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/store"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	// "github.com/cosmos/cosmos-sdk/x/auth"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/stretchr/testify/require"

	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cosmos/cosmos-sdk/x/nft"
	ibctransfertypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	alliancemoduletypes "github.com/terra-money/alliance/x/alliance/types"
	ibcfeetypes "github.com/cosmos/ibc-go/v7/modules/apps/29-fee/types"
	icatypes "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/types"

	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	bankkeeper "github.com/terra-money/alliance/custom/bank/keeper"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
)

func NibtdfvasKeeper(t testing.TB) (*keeper.Keeper, sdk.Context) {

	const Bech32Prefix = "andr"
	storeKey := sdk.NewKVStoreKey(types.StoreKey)
	memStoreKey := storetypes.NewMemoryStoreKey(types.MemStoreKey)

	db := tmdb.NewMemDB()
	stateStore := store.NewCommitMultiStore(db)
	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(memStoreKey, storetypes.StoreTypeMemory, nil)
	require.NoError(t, stateStore.LoadLatestVersion())

	registry := codectypes.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(registry)

	maccPerms := map[string][]string{
		authtypes.FeeCollectorName:          {authtypes.Burner},
		distrtypes.ModuleName:               nil,
		minttypes.ModuleName:                {authtypes.Minter},
		stakingtypes.BondedPoolName:         {authtypes.Burner, authtypes.Staking},
		stakingtypes.NotBondedPoolName:      {authtypes.Burner, authtypes.Staking},
		govtypes.ModuleName:                 {authtypes.Burner},
		nft.ModuleName:                      nil,
		ibctransfertypes.ModuleName:         {authtypes.Minter, authtypes.Burner},
		ibcfeetypes.ModuleName:              nil,
		icatypes.ModuleName:                 nil,
		wasmtypes.ModuleName:                {authtypes.Burner},
		alliancemoduletypes.ModuleName:      {authtypes.Minter, authtypes.Burner},
		alliancemoduletypes.RewardsPoolName: nil,
		types.ModuleName:		 {authtypes.Minter, authtypes.Burner, authtypes.Staking},
	}

	modAccAddrs := make(map[string]bool)
	for acc := range maccPerms {
		modAccAddrs[authtypes.NewModuleAddress(acc).String()] = true
	}

	accountKeeper := authkeeper.NewAccountKeeper(
		cdc, storeKey, authtypes.ProtoBaseAccount, maccPerms, Bech32Prefix, authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)
	bankKeeper := bankkeeper.NewBaseKeeper(
		cdc, storeKey, accountKeeper, modAccAddrs, authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)
	stakingKeeper := stakingkeeper.NewKeeper(
		cdc, storeKey, accountKeeper, bankKeeper, authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	k := keeper.NewKeeper(
		cdc,
		storeKey,
		memStoreKey,
		authtypes.NewModuleAddress("nibtdfvas"),
		accountKeeper,
		bankKeeper,
		stakingKeeper,
	)

	ctx := sdk.NewContext(stateStore, tmproto.Header{}, false, log.NewNopLogger())

	// Initialize params
	_ = k.SetParams(ctx, types.DefaultParams())

	// moduleAccount := auth.NewModuleAccount(types.ModuleName, auth.Burner, auth.Minter)

	// moduleAcc := k.GetNibtdfvasAccount(ctx)
	// if moduleAcc == nil {
	// 	panic(fmt.Sprintf("%s module account has not been set", types.ModuleName))
	// }

	return k, ctx
}
