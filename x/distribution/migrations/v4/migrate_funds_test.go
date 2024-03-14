package v4_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/log"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	v4 "github.com/cosmos/cosmos-sdk/x/distribution/migrations/v4"
	distrtestutil "github.com/cosmos/cosmos-sdk/x/distribution/testutil"
	disttypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	pooltypes "github.com/cosmos/cosmos-sdk/x/protocolpool/types"

	addresscodec "github.com/cosmos/cosmos-sdk/codec/address"
	codectestutil "github.com/cosmos/cosmos-sdk/codec/testutil"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/testutil/integration"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
)

func TestFundsMigration(t *testing.T) {
	keys := storetypes.NewKVStoreKeys(
		authtypes.StoreKey, banktypes.StoreKey, disttypes.StoreKey,
	)
	logger := log.NewTestLogger(t)
	cms := integration.CreateMultiStore(keys, logger)
	encCfg := moduletestutil.MakeTestEncodingConfig(codectestutil.CodecOptions{}, auth.AppModule{}, bank.AppModule{}, distribution.AppModule{})
	ctx := sdk.NewContext(cms, true, logger)

	maccPerms := map[string][]string{
		pooltypes.ModuleName: nil,
		disttypes.ModuleName: {authtypes.Minter},
	}

	authority := authtypes.NewModuleAddress("gov")

	// create account keeper
	accountKeeper := authkeeper.NewAccountKeeper(
		runtime.NewEnvironment(runtime.NewKVStoreService(keys[authtypes.StoreKey]), log.NewNopLogger()),
		encCfg.Codec,
		authtypes.ProtoBaseAccount,
		maccPerms,
		addresscodec.NewBech32Codec(sdk.Bech32MainPrefix),
		sdk.Bech32MainPrefix,
		authority.String(),
	)

	// create bank keeper
	bankKeeper := bankkeeper.NewBaseKeeper(
		runtime.NewEnvironment(runtime.NewKVStoreService(keys[banktypes.StoreKey]), log.NewNopLogger()),
		encCfg.Codec,
		accountKeeper,
		map[string]bool{},
		authority.String(),
	)

	// gomock initializations
	ctrl := gomock.NewController(t)
	stakingKeeper := distrtestutil.NewMockStakingKeeper(ctrl)
	poolKeeper := distrtestutil.NewMockPoolKeeper(ctrl)

	// create distribution keeper
	distrKeeper := keeper.NewKeeper(
		encCfg.Codec,
		runtime.NewEnvironment(runtime.NewKVStoreService(keys[disttypes.StoreKey]), log.NewNopLogger()),
		accountKeeper,
		bankKeeper,
		stakingKeeper,
		poolKeeper,
		disttypes.ModuleName,
		authority.String(),
	)

	// Set feepool
	poolAmount := sdk.NewInt64Coin("test", 100000)
	feepool := disttypes.FeePool{
		CommunityPool: sdk.NewDecCoinsFromCoins(poolAmount),
	}
	err := distrKeeper.FeePool.Set(ctx, feepool)
	require.NoError(t, err)

	distrAcc := authtypes.NewEmptyModuleAccount(disttypes.ModuleName)

	// mint coins in distribution module account
	distrModBal := sdk.NewCoins(sdk.NewInt64Coin("test", 10000000))
	err = bankKeeper.MintCoins(ctx, distrAcc.GetName(), distrModBal)
	require.NoError(t, err)

	// Set pool module account
	poolAcc := authtypes.NewEmptyModuleAccount(pooltypes.ModuleName)

	// migrate feepool funds from distribution module account to pool module account
	_, err = v4.MigrateFunds(ctx, bankKeeper, feepool, distrAcc, poolAcc)
	require.NoError(t, err)

	// set distribution feepool as empty (since migration)
	err = distrKeeper.FeePool.Set(ctx, disttypes.FeePool{})
	require.NoError(t, err)

	// check pool module account balance equals pool amount
	poolMAccBal := bankKeeper.GetAllBalances(ctx, poolAcc.GetAddress())
	require.Equal(t, poolMAccBal, sdk.Coins{poolAmount})

	distrAccBal := bankKeeper.GetAllBalances(ctx, distrAcc.GetAddress())
	// check distribution module account balance is not same after migration
	require.NotEqual(t, distrModBal, distrAccBal)
	// check distribution module account balance is same as (current distrAccBal+poolAmount)
	require.Equal(t, distrModBal, distrAccBal.Add(poolAmount))
}
