package keeper

import (
	"context"

	"cosmossdk.io/math"
)

// GetCommunityTax returns the current distribution community tax.
func (k Keeper) GetCommunityTax(ctx context.Context) (math.LegacyDec, error) {
	params, err := k.Params.Get(ctx)
	if err != nil {
		return math.LegacyDec{}, err
	}

	return params.CommunityTax, nil
}

// GetWithdrawAddrEnabled returns the current distribution withdraw address
// enabled parameter.
func (k Keeper) GetWithdrawAddrEnabled(ctx context.Context) (enabled bool, err error) {
	params, err := k.Params.Get(ctx)
	if err != nil {
		return false, err
	}

	return params.WithdrawAddrEnabled, nil
}

// Get RewardsPerBlock returns the current distribution rewards per block
func (k Keeper) GetRewardsPerBlock(ctx context.Context) (math.LegacyDec, error) {
	// return k.GetParams(ctx).RewardsPerBlock
	params, err := k.Params.Get(ctx)
	if err != nil {
		return math.LegacyDec{}, err
	}

	return params.RewardsPerBlock, nil
}
