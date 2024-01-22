package keeper

import (
	"github.com/andromedaprotocol/andromedad/x/nibtdfvas/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetParams get all parameters as types.Params
func (k Keeper) GetDAOParams(ctx sdk.Context) types.DAOParams {
	var daoParams types.DAOParams
	store := ctx.KVStore(k.storeKey)
	bz := store.Get([]byte("DAOParams"))
	if bz == nil {
		return types.DefaultDAOParams()
	}
	sdk.cdc.MustUnmarshalBinaryBare(bz, &daoParams)
	return daoParams
}

// SetParams set the params
func (k Keeper) SetDAOParams(ctx sdk.Context, daoParams types.DAOParams) {
	// k.paramstore.SetParamSet(ctx, &params)
	store := ctx.KVStore(k.storeKey)
	store.SetParams([]byte("DAOParams"), sdk.cdc.MustUnmarshalBinaryBare(&daoParams))
}