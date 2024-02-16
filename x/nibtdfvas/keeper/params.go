package keeper

import (
	"github.com/andromedaprotocol/andromedad/x/nibtdfvas/types"
	// "andromedad/x/nibtdfvas/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetParams get all parameters as types.Params
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.ParamsKey)
	if len(bz) == 0 {
		return params
	}

	k.cdc.MustUnmarshal(bz, &params)
	return params
}

// SetParams set the params
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) error {
	if err := params.Validate(); err != nil {
		return err
	}

	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&params)
	store.Set(types.ParamsKey, bz)

	return nil
}

func (k Keeper) GetTokenOutflowPerBlock(ctx sdk.Context) int64 {
	params := k.GetParams(ctx)
	return params.TokenOutflowPerBlock
}

func (k Keeper) GetDirectToValidator(ctx sdk.Context) int64 {
	params := k.GetParams(ctx)
	return params.DirectToValidatorPercent
}
