package keeper_test

import (
	"testing"
	"fmt"

	"cosmossdk.io/math"
	abci "github.com/cometbft/cometbft/abci/types"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/stretchr/testify/require"

	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	// "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	// banktestutil "github.com/cosmos/cosmos-sdk/x/bank/testutil"
	// "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	// "github.com/cosmos/cosmos-sdk/x/distribution/testutil"
	// disttypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtestutil "github.com/cosmos/cosmos-sdk/x/staking/testutil"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	nibtdfvaskeeper "github.com/andromedaprotocol/andromedad/x/nibtdfvas/keeper"
	"github.com/andromedaprotocol/andromedad/x/nibtdfvas/testutil"
	"github.com/andromedaprotocol/andromedad/x/nibtdfvas/types"
)

func TestNibtdfvasKeeper(t *testing.T) {
	var (
		accountKeeper      authkeeper.AccountKeeper
		bankKeeper         bankkeeper.Keeper
		stakingKeeper      *stakingkeeper.Keeper
		nibtdfvasKeeper    nibtdfvaskeeper.Keeper
	)

	app, err := simtestutil.Setup(testutil.AppConfig,
		&accountKeeper,
		&bankKeeper,
		&stakingKeeper,
		&nibtdfvasKeeper,
	)
	require.NoError(t, err)

	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	addrs := simtestutil.AddTestAddrs(bankKeeper, stakingKeeper, ctx, 3, sdk.NewInt(1234))
	valAddrs := simtestutil.ConvertAddrsToValAddrs(addrs)
	tstaking := stakingtestutil.NewHelper(t, ctx, stakingKeeper)

	// create validator with 50% commission
	tstaking.Commission = stakingtypes.NewCommissionRates(sdk.NewDecWithPrec(5, 1), sdk.NewDecWithPrec(5, 1), math.LegacyNewDec(0))
	tstaking.CreateValidator(sdk.ValAddress(addrs[0]), valConsPk0, sdk.NewInt(300), true)

	// create second validator with 0% commission
	tstaking.Commission = stakingtypes.NewCommissionRates(math.LegacyNewDec(0), math.LegacyNewDec(0), math.LegacyNewDec(0))
	tstaking.CreateValidator(valAddrs[1], valConsPk1, sdk.NewInt(200), true)

	tokens := sdk.Coins{
		{Denom: "andr", Amount: math.NewInt(10000)},
	}

	moduleAddress := accountKeeper.GetModuleAddress(types.ModuleName)
	// fmt.Println("moduleAddress value:", moduleAddress)

	preModuleBalance := bankKeeper.GetBalance(ctx, moduleAddress, "andr").Amount
	// fmt.Println("preModuleBalance value:", preModuleBalance)
	require.Equal(t, preModuleBalance, sdk.NewInt(1000000))

	err = nibtdfvasKeeper.MintCoins(ctx, types.ModuleName, tokens)
	require.NoError(t, err)

	afterModuleBalance := bankKeeper.GetBalance(ctx, moduleAddress, "andr").Amount
	// fmt.Println("afterModuleBalance value:", afterModuleBalance)
	require.Equal(t, afterModuleBalance, sdk.NewInt(1010000))

	acc := nibtdfvasKeeper.GetNibtdfvasAccount(ctx)
	// fmt.Println("acc value:", acc)
	require.NotNil(t, acc)

	err = nibtdfvasKeeper.SendCoinsFromModuleToModule(ctx, tokens)
	require.NoError(t, err)

	updatedVal := nibtdfvasKeeper.ValidatorUpdate(ctx)
	// fmt.Println("updatedVal value:", updatedVal)
	require.Nil(t, updatedVal)

	

	abciValA := abci.Validator{
		Address: valConsPk0.Address(),
		Power:   60,
	}
	abciValB := abci.Validator{
		Address: valConsPk1.Address(),
		Power:   40,
	}

	votes := []abci.VoteInfo{
		{
			Validator: abciValA,
			SignedLastBlock: true,
		},
		{
			Validator:       abciValB,
			SignedLastBlock: true,
		},
	}

	info := abci.CommitInfo{
		Votes: votes,
	}

	req := abci.RequestBeginBlock{
		LastCommitInfo: info,
	}

	validator, _ := stakingKeeper.GetValidatorByConsAddr(ctx, abciValA.Address)

	preBalance := bankKeeper.GetBalance(ctx, sdk.AccAddress(validator.GetOperator()), "andr").Amount
	fmt.Println("preBalance value:", preBalance)

	nibtdfvasKeeper.DistributeTokensToValidators(ctx, req, math.LegacyNewDec(20))

	afterBalance := bankKeeper.GetBalance(ctx, sdk.AccAddress(validator.GetOperator()), "andr").Amount
	fmt.Println("afterBalance value:", afterBalance)
	fmt.Println("operatorAddr value:", sdk.AccAddress(validator.GetOperator()))

	nibtdfvasKeeper.DistributeTokensToStakers(ctx, req, math.LegacyNewDec(80))
	fmt.Println("afterBalance value:", bankKeeper.GetBalance(ctx, sdk.AccAddress(validator.GetOperator()), "andr").Amount)

	// require.Nil(t, preBalance)
}