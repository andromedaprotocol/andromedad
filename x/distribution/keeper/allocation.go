package keeper

import (
	"cosmossdk.io/math"
	abci "github.com/cometbft/cometbft/abci/types"

	"github.com/andromedaprotocol/andromedad/x/distribution/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// AllocateTokens performs reward and fee distribution to all validators based
// on the F1 fee distribution specification.
func (k Keeper) AllocateTokens(ctx sdk.Context, totalPreviousPower int64, bondedVotes []abci.VoteInfo) {
	// fetch and clear the collected fees for distribution, since this is
	// called in BeginBlock, collected fees will be from the previous block
	// (and distributed to the previous proposer)
	feeCollector := k.authKeeper.GetModuleAccount(ctx, k.feeCollectorName)
	feesCollectedInt := k.bankKeeper.GetAllBalances(ctx, feeCollector.GetAddress())
	feesCollected := sdk.NewDecCoinsFromCoins(feesCollectedInt...)

	// Fetch the RewardsDripper module account
	rewardsDripper := k.authKeeper.GetModuleAccount(ctx, types.RewardsDripperName)
	// Fetch the rewards dripper balance
	rewardsDripperBalance := k.bankKeeper.GetAllBalances(ctx, rewardsDripper.GetAddress())
	// Convert rewardsDripperBalance to DecCoins
	rewardsDripperCollected := sdk.NewDecCoinsFromCoins(rewardsDripperBalance...)

	// transfer collected fees to the distribution module account
	err := k.bankKeeper.SendCoinsFromModuleToModule(ctx, k.feeCollectorName, types.ModuleName, feesCollectedInt)
	if err != nil {
		panic(err)
	}

	// Calculate rewards to be dripped this block from Param set
	rewardsToDrip := k.GetRewardsPerBlock(ctx)

	// Create new coins with the denoms of the rewardsDripperBalance and the amount of rewards to be dripped
	rewardsCoins := make(sdk.Coins, len(rewardsDripperBalance))
	for i, coin := range rewardsDripperBalance {
		rewardsCoins[i] = sdk.NewCoin(coin.Denom, rewardsToDrip.TruncateInt())
	}

	// Convert to DecCoins
	rewardsToDripDec := sdk.NewDecCoinsFromCoins(rewardsCoins...)

	// Intersect balance of rewardsDripper with rewardsToDripDec to find the amount to be dripped
	rewardsToDripDec = rewardsToDripDec.Intersect(rewardsDripperCollected)

	// Convert rewardsToDripDec to Coins
	rewardsToDripInt, _ := rewardsToDripDec.TruncateDecimal()

	// transfer rewards to be dripped to the distribution module account
	if err := k.bankKeeper.SendCoinsFromModuleToModule(ctx, types.RewardsDripperName, types.ModuleName, rewardsToDripInt); err != nil {
		panic(err)
	}
	// temporary workaround to keep CanWithdrawInvariant happy
	// general discussions here: https://github.com/cosmos/cosmos-sdk/issues/2906#issuecomment-441867634
	feePool := k.GetFeePool(ctx)
	if totalPreviousPower == 0 {
		feePool.CommunityPool = feePool.CommunityPool.Add(feesCollected...)
		k.SetFeePool(ctx, feePool)
		return
	}

	// Combine all rewards
	allCollected := feesCollected.Add(rewardsToDripDec...)
	// calculate fraction allocated to validators
	remaining := allCollected
	communityTax := k.GetCommunityTax(ctx)
	voteMultiplier := math.LegacyOneDec().Sub(communityTax)
	feeMultiplier := feesCollected.MulDecTruncate(voteMultiplier)
	// To avoid adding a community tax to rewards to be dripped we add the rewardsToDripDec to the feeMultiplier
	// We DO NOT want to re-tax funds that already come from the pool as these are basely rewards
	feeMultiplier = feeMultiplier.Add(rewardsToDripDec...)

	// allocate tokens proportionally to voting power
	//
	// TODO: Consider parallelizing later
	//
	// Ref: https://github.com/cosmos/cosmos-sdk/pull/3099#discussion_r246276376
	for _, vote := range bondedVotes {
		validator := k.stakingKeeper.ValidatorByConsAddr(ctx, vote.Validator.Address)

		// TODO: Consider micro-slashing for missing votes.
		//
		// Ref: https://github.com/cosmos/cosmos-sdk/issues/2525#issuecomment-430838701
		powerFraction := math.LegacyNewDec(vote.Validator.Power).QuoTruncate(math.LegacyNewDec(totalPreviousPower))
		reward := feeMultiplier.MulDecTruncate(powerFraction)

		k.AllocateTokensToValidator(ctx, validator, reward)
		remaining = remaining.Sub(reward)
	}

	// allocate community funding
	feePool.CommunityPool = feePool.CommunityPool.Add(remaining...)
	k.SetFeePool(ctx, feePool)
}

// AllocateTokensToValidator allocate tokens to a particular validator,
// splitting according to commission.
func (k Keeper) AllocateTokensToValidator(ctx sdk.Context, val stakingtypes.ValidatorI, tokens sdk.DecCoins) {
	// split tokens between validator and delegators according to commission
	commission := tokens.MulDec(val.GetCommission())
	shared := tokens.Sub(commission)

	// update current commission
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeCommission,
			sdk.NewAttribute(sdk.AttributeKeyAmount, commission.String()),
			sdk.NewAttribute(types.AttributeKeyValidator, val.GetOperator().String()),
		),
	)
	currentCommission := k.GetValidatorAccumulatedCommission(ctx, val.GetOperator())
	currentCommission.Commission = currentCommission.Commission.Add(commission...)
	k.SetValidatorAccumulatedCommission(ctx, val.GetOperator(), currentCommission)

	// update current rewards
	currentRewards := k.GetValidatorCurrentRewards(ctx, val.GetOperator())
	currentRewards.Rewards = currentRewards.Rewards.Add(shared...)
	k.SetValidatorCurrentRewards(ctx, val.GetOperator(), currentRewards)

	// update outstanding rewards
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeRewards,
			sdk.NewAttribute(sdk.AttributeKeyAmount, tokens.String()),
			sdk.NewAttribute(types.AttributeKeyValidator, val.GetOperator().String()),
		),
	)

	outstanding := k.GetValidatorOutstandingRewards(ctx, val.GetOperator())
	outstanding.Rewards = outstanding.Rewards.Add(tokens...)
	k.SetValidatorOutstandingRewards(ctx, val.GetOperator(), outstanding)
}
