package nibtdfvas

import (
	"time"

	abci "github.com/cometbft/cometbft/abci/types"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/andromedaprotocol/andromedad/x/nibtdfvas/keeper"
	// "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	"github.com/andromedaprotocol/andromedad/x/nibtdfvas/types"
)


// Called every block, update validator set
// func EndBlocker(ctx sdk.Context, k types.StakingKeeper) []abci.ValidatorUpdate {
// 	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyEndBlocker)
// 	res := k.BlockValidatorUpdates(ctx)
// 	return res
// }
func BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock, k keeper.Keeper) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyBeginBlocker)
	// TokenOutflowPerBlock := k.GetTokenOutflowPerBlock(ctx)
	// DirectToValidatorPercent := k.GetDirectToValidator(ctx)
	k.DistributeTokens(ctx, req)

}


func EndBlocker(ctx sdk.Context, k keeper.Keeper) []abci.ValidatorUpdate {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyEndBlocker)
	return k.ValidatorUpdate(ctx)
}