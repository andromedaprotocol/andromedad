package nibtdfvas

import (
	"fmt"
	errorsmod "cosmossdk.io/errors"

	// "andromedad/x/nibtdfvas/keeper"
	// "andromedad/x/nibtdfvas/types"

	"github.com/andromedaprotocol/andromedad/x/nibtdfvas/keeper"
	"github.com/andromedaprotocol/andromedad/x/nibtdfvas/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	err := k.SetParams(ctx, genState.Params)
	if err != nil {
		panic(errorsmod.Wrapf(err, "error setting params"))
	}

	moduleAcc := k.GetNibtdfvasAccount(ctx)
	if moduleAcc == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.ModuleName))
	}

	// Mint 10,000 stake coins
	mintCoins := sdk.NewCoins(sdk.NewCoin("token", sdk.NewInt(1000000)))

	if err := k.Mint(ctx, mintCoins); err != nil {
		panic(err)
	}

	if err := k.SendCoinsFromModuleToModule(ctx, mintCoins); err != nil {
    panic(err)
  }

}

// ExportGenesis returns the module's exported genesis
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)

	return genesis
}