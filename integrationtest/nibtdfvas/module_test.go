package nibtdfvas_test

import (
	"testing"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/stretchr/testify/require"

	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/andromedaprotocol/andromedad/x/nibtdfvas/testutil"
	"github.com/andromedaprotocol/andromedad/x/nibtdfvas/types"
	nibtdfvaskeeper "github.com/andromedaprotocol/andromedad/x/nibtdfvas/keeper"
)

func TestItCreatesModuleAccountOnInitBlock(t *testing.T) {
	var accountKeeper authkeeper.AccountKeeper
	var bankKeeper bankkeeper.Keeper
	var stakingKeeper *stakingkeeper.Keeper
	var nibtdfvasKeeper nibtdfvaskeeper.Keeper

	app, err := simtestutil.Setup(testutil.AppConfig, &accountKeeper, &bankKeeper, stakingKeeper, &nibtdfvasKeeper)
	require.NoError(t, err)

	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	acc := accountKeeper.GetAccount(ctx, authtypes.NewModuleAddress(types.ModuleName))
	require.NotNil(t, acc)
}