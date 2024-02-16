package keeper_test
import (
	"testing"
	"fmt"
	"time"

	// "github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	// "github.com/cosmos/cosmos-sdk/store/types"
	"cosmossdk.io/math"
	abci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/golang/mock/gomock"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"

	// simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	"github.com/cosmos/cosmos-sdk/testutil"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/andromedaprotocol/andromedad/x/nibtdfvas"
	"github.com/andromedaprotocol/andromedad/x/nibtdfvas/keeper"
	nibtdfvastestutil "github.com/andromedaprotocol/andromedad/x/nibtdfvas/testutil"
	nibtdfvastypes "github.com/andromedaprotocol/andromedad/x/nibtdfvas/types"
)

var (
	nibtdfvasAcct     = authtypes.NewModuleAddress("nibtdfvas")
)

func GetNibtdfvasKeeper(t *testing.T) (
	keeper.Keeper,
	*nibtdfvastestutil.MockAccountKeeper,
	*nibtdfvastestutil.MockBankKeeper,
	*nibtdfvastestutil.MockStakingKeeper,
	moduletestutil.TestEncodingConfig,
	sdk.Context,
) {
	ctrl := gomock.NewController(t)
	key := sdk.NewKVStoreKey(nibtdfvastypes.StoreKey)
	testCtx := testutil.DefaultContextWithDB(t, key, sdk.NewTransientStoreKey("transient_test"))
	encCfg := moduletestutil.MakeTestEncodingConfig(nibtdfvas.AppModuleBasic{})
	ctx := testCtx.Ctx.WithBlockHeader(tmproto.Header{Time: time.Now()})
	// ctx := sdk.NewContext(stateStore, tmproto.Header{}, false, log.NewNopLogger())

	accountKeeper := nibtdfvastestutil.NewMockAccountKeeper(ctrl)
	bankKeeper := nibtdfvastestutil.NewMockBankKeeper(ctrl)
	stakingKeeper := nibtdfvastestutil.NewMockStakingKeeper(ctrl)
	balances := make(map[string]sdk.Coins)
	balances[nibtdfvasAcct.String()] = sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(0)))

	accountKeeper.EXPECT().GetModuleAddress(nibtdfvastypes.ModuleName).Return(authtypes.NewModuleAddress(nibtdfvastypes.ModuleName)).AnyTimes()
	accountKeeper.EXPECT().GetModuleAccount(gomock.Any(), nibtdfvastypes.ModuleName).Return(authtypes.NewEmptyModuleAccount(nibtdfvastypes.ModuleName)).AnyTimes()

	stakingKeeper.EXPECT().BondDenom(ctx).Return("stake").AnyTimes()
	stakingKeeper.EXPECT().TotalBondedTokens(gomock.Any()).Return(math.NewInt(10000000)).AnyTimes()
	stakingKeeper.EXPECT().BlockValidatorUpdates(ctx).Return(nil).AnyTimes()
	stakingKeeper.EXPECT().GetBondedValidatorsByPower(ctx).Return(nil).AnyTimes()
	stakingKeeper.EXPECT().GetValidatorDelegations(ctx, gomock.Any()).Return(nil).AnyTimes()
	stakingKeeper.EXPECT().ValidatorByConsAddr(ctx, gomock.Any()).Return(nil).AnyTimes()
	stakingKeeper.EXPECT().GetValidatorByConsAddr(ctx, gomock.Any()).Return(stakingtypes.Validator{}, true).AnyTimes()

	bankKeeper.EXPECT().SendCoinsFromModuleToModule(gomock.Any(), nibtdfvastypes.ModuleName, nibtdfvastypes.ModuleName, sdk.NewCoins(sdk.NewCoin("andr", sdk.NewInt(100)))).AnyTimes()
	bankKeeper.EXPECT().SendCoinsFromModuleToAccount(gomock.Any(), nibtdfvastypes.ModuleName, gomock.Any(), gomock.Any()).AnyTimes()
	bankKeeper.EXPECT().MintCoins(ctx, "nibtdfvas", sdk.NewCoins(sdk.NewCoin("andr", sdk.NewInt(100)))).AnyTimes()
	bankKeeper.EXPECT().MintCoins(ctx, "nibtdfvas", sdk.NewCoins(sdk.NewCoin("andr", sdk.NewInt(1000000)))).AnyTimes()
	bankKeeper.EXPECT().GetBalance(ctx, gomock.Any(), "stake").Return(sdk.NewCoin("stake", sdk.NewInt(100))).AnyTimes()
	bankKeeper.EXPECT().GetBalance(ctx, gomock.Any(), "andr").Return(sdk.NewCoin("andr", sdk.NewInt(1000))).AnyTimes()

	bankKeeper.EXPECT().SendCoinsFromModuleToAccount(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(_ sdk.Context, module string, rcpt sdk.AccAddress, coins sdk.Coins) error {
		balances[rcpt.String()] = balances[rcpt.String()].Add(coins...)
		return nil
	}).AnyTimes()

	nibtdfvasKeeper := keeper.NewKeeper(
		encCfg.Codec,
		key,
		// key,
		// the address capable of executing a MsgUpdateParams message. Typically, this should be the x/gov module account.
		authtypes.NewModuleAddress("gov"),
		accountKeeper,
		bankKeeper,
		stakingKeeper,
	)

	return nibtdfvasKeeper, accountKeeper, bankKeeper, stakingKeeper, encCfg, ctx
}

func TestKeeperFunctions(t * testing.T) {
	nibtdfvasKeeper, accountKeeper, bankKeeper, _, _, ctx := GetNibtdfvasKeeper(t)
	nibtdfvasKeeper.SetParams(ctx, nibtdfvastypes.Params{
		TokenOutflowPerBlock: 3,
		DirectToValidatorPercent: 20,
	})

	nibtdfvasKeeper.ValidatorUpdate(ctx)

	currentParams := nibtdfvasKeeper.GetParams(ctx)
	fmt.Println("currentParams value:", currentParams)

	blockReward := math.LegacyNewDec(currentParams.TokenOutflowPerBlock)
	fmt.Println("blockReward value:", blockReward)

	directToValidator := (math.LegacyNewDec(currentParams.DirectToValidatorPercent).Mul(blockReward)).Quo(math.LegacyNewDec(100))
	fmt.Println("directToValidator value:", directToValidator)

	// nibtdfvasKeeper.MintCoins(ctx, "nibtdfvas", sdk.NewCoins(sdk.NewCoin("andr", sdk.NewInt(100))))
	// nibtdfvasKeeper.SendCoinsFromModuleToModule(ctx, sdk.NewCoins(sdk.NewCoin("andr", sdk.NewInt(100))))

	// bankKeeper.MintCoins(ctx, "nibtdfvas", sdk.NewCoins(sdk.NewCoin("andr", sdk.NewInt(100))))
	// bankKeeper.SendCoinsFromModuleToModule(ctx, "nibtdfvas", "nibtdfvas", sdk.NewCoins(sdk.NewCoin("andr", sdk.NewInt(100))))

	moduleAddress := accountKeeper.GetModuleAddress("nibtdfvas")
	fmt.Println("Parameter value:", moduleAddress)

	moduleBalance := bankKeeper.GetBalance(ctx, moduleAddress, "andr").Amount
	fmt.Println("Parameter value:", moduleBalance)
	// require.Equal(t, moduleBalance, sdk.NewInt(0))

	nibtdfvasKeeper.SendCoinsFromModuleToModule(ctx, nibtdfvastypes.ModuleName, nibtdfvastypes.ModuleName, sdk.NewCoins(sdk.NewCoin("andr", sdk.NewInt(100))))

	abciValA := abci.Validator{
		Address: valConsPk0.Address(),
		Power:   100,
	}

	votes := []abci.VoteInfo{
		{
			Validator: abciValA,
			SignedLastBlock: true,
		},
	}

	info := abci.CommitInfo{
		Votes: votes,
	}

	req := abci.RequestBeginBlock{
		LastCommitInfo: info,
	}

	preBalance := bankKeeper.GetBalance(ctx, abciValA.Address, "andr").Amount
	fmt.Println("Parameter value:", preBalance)

	nibtdfvasKeeper.DistributeTokensToValidators(ctx, req, math.LegacyNewDec(3))

	afterBalance := bankKeeper.GetBalance(ctx, abciValA.Address, "andr").Amount
	fmt.Println("Parameter value:", afterBalance)


	require.True(t, (afterBalance.Sub(preBalance)).LT(sdk.NewInt(1)))

	nibtdfvasKeeper.DistributeTokensToStakers(ctx, req, math.LegacyNewDec(2))
}
