package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/andromedaprotocol/andromedad/x/distribution/types"
)

// get outstanding rewards
func (k Keeper) GetValidatorOutstandingRewardsCoins(ctx context.Context, val sdk.ValAddress) (sdk.DecCoins, error) {
	rewards, err := k.GetValidatorOutstandingRewards(ctx, val)
	if err != nil {
		return nil, err
	}

	return rewards.Rewards, nil
}

// GetDistributionAccount returns the distribution ModuleAccount
func (k Keeper) GetDistributionAccount(ctx context.Context) sdk.ModuleAccountI {
	return k.authKeeper.GetModuleAccount(ctx, types.ModuleName)
}

// GetRewardsDripperAccount returns the rewards dripper ModuleAccount
func (k Keeper) GetRewardsDripperAccount(ctx sdk.Context) sdk.ModuleAccountI {
	return k.authKeeper.GetModuleAccount(ctx, types.RewardsDripperName)
}
