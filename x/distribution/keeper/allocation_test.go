package keeper_test

import (
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"

	"github.com/cosmos/cosmos-sdk/codec/address"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"

	"github.com/andromedaprotocol/andromedad/x/distribution"
	"github.com/andromedaprotocol/andromedad/x/distribution/keeper"
	distrtestutil "github.com/andromedaprotocol/andromedad/x/distribution/testutil"
	disttypes "github.com/andromedaprotocol/andromedad/x/distribution/types"
)

func TestAllocateTokensToValidatorWithCommission(t *testing.T) {
	ctrl := gomock.NewController(t)
	key := storetypes.NewKVStoreKey(disttypes.StoreKey)
	storeService := runtime.NewKVStoreService(key)
	testCtx := testutil.DefaultContextWithDB(t, key, storetypes.NewTransientStoreKey("transient_test"))
	encCfg := moduletestutil.MakeTestEncodingConfig(distribution.AppModuleBasic{})
	ctx := testCtx.Ctx.WithBlockHeader(cmtproto.Header{Time: time.Now()})

	bankKeeper := distrtestutil.NewMockBankKeeper(ctrl)
	stakingKeeper := distrtestutil.NewMockStakingKeeper(ctrl)
	accountKeeper := distrtestutil.NewMockAccountKeeper(ctrl)

	valCodec := address.NewBech32Codec("cosmosvaloper")

	accountKeeper.EXPECT().GetModuleAddress("distribution").Return(distrAcc.GetAddress())
	// set rewards dripper module
	accountKeeper.EXPECT().GetModuleAddress("rewards_dripper").Return(rewardsAcc.GetAddress())
	stakingKeeper.EXPECT().ValidatorAddressCodec().Return(valCodec).AnyTimes()

	distrKeeper := keeper.NewKeeper(
		encCfg.Codec,
		storeService,
		accountKeeper,
		bankKeeper,
		stakingKeeper,
		"fee_collector",
		authtypes.NewModuleAddress("gov").String(),
	)

	// create validator with 50% commission
	val, err := distrtestutil.CreateValidator(valConsPk0, math.NewInt(100))
	require.NoError(t, err)
	val.Commission = stakingtypes.NewCommission(math.LegacyNewDecWithPrec(5, 1), math.LegacyNewDecWithPrec(5, 1), math.LegacyNewDec(0))
	stakingKeeper.EXPECT().ValidatorByConsAddr(gomock.Any(), sdk.GetConsAddress(valConsPk0)).Return(val, nil).AnyTimes()

	// allocate tokens
	tokens := sdk.DecCoins{
		{Denom: sdk.DefaultBondDenom, Amount: math.LegacyNewDec(10)},
	}
	require.NoError(t, distrKeeper.AllocateTokensToValidator(ctx, val, tokens))

	// check commission
	expected := sdk.DecCoins{
		{Denom: sdk.DefaultBondDenom, Amount: math.LegacyNewDec(5)},
	}

	valBz, err := valCodec.StringToBytes(val.GetOperator())
	require.NoError(t, err)

	valCommission, err := distrKeeper.GetValidatorAccumulatedCommission(ctx, valBz)
	require.NoError(t, err)
	require.Equal(t, expected, valCommission.Commission)

	// check current rewards
	currentRewards, err := distrKeeper.GetValidatorCurrentRewards(ctx, valBz)
	require.NoError(t, err)
	require.Equal(t, expected, currentRewards.Rewards)
}
