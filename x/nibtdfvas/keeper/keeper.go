package keeper

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/log"
	// "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"andromedad/x/nibtdfvas/types"
)

type (
	Keeper struct {
		cdc        codec.BinaryCodec
		storeKey   sdk.StoreKey
		memKey     sdk.StoreKey
		paramstore paramtypes.Subspace

		stakingKeeper types.StakingKeeper
		bankKeeper types.BankKeeper
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey,
	memKey sdk.StoreKey,
	ps paramtypes.Subspace,

	// stakingKeeper types.StakingKeeper,
) *Keeper {
	// set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	return &Keeper{

		cdc:           cdc,
		storeKey:      storeKey,
		memKey:        memKey,
		paramstore:    ps,
		// stakingKeeper: stakingKeeper,
		// bankKeeper:		 bankKeeper,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) DistributeTokens(ctx sdk.Context, params types.Params) {
	// Get the DAO parameters from the store or use default values
	currentParams := k.GetParams(ctx)

	if currentParams.TokenOutflowPerBlock != params.TokenOutflowPerBlock{
		currentParams.TokenOutflowPerBlock = params.TokenOutflowPerBlock
	}

	if currentParams.DirectToValidatorPercent != params.DirectToValidatorPercent{
		currentParams.DirectToValidatorPercent = params.DirectToValidatorPercent
	}

	// k.SetParams(ctx, currentParams)

	// Distribute tokens per block
	blockReward := currentParams.TokenOutflowPerBlock
	directToValidator := ((currentParams.DirectToValidatorPercent) * (blockReward)) / 100
	toStakers := blockReward - directToValidator

	// Distribute tokens to validators
	k.DistributeTokensToValidators(ctx, directToValidator)

	// Distribute tokens to stakers
	k.DistributeTokensToStakers(ctx, toStakers)
}

func (k Keeper) DistributeTokensToValidators(ctx sdk.Context, amount int64) {
	// Implement logic to distribute tokens to validators
	// daoParams := k.GetDAOParams(ctx)
	totalStakingTokens := k.stakingKeeper.TotalBondedTokens(ctx)

	// // Iterate over validators and distribute tokens
	validatorIterator := k.stakingKeeper.GetBondedValidatorsByPower(ctx)
	for _, validator := range validatorIterator {
		validatorAddr := validator.OperatorAddress

		// Calculate the amount based on the staking ratio
		validatorStake := validator.BondedTokens()
		validatorAmount := amount * (validatorStake.Int64()) / (totalStakingTokens.Int64())

		// x := validatorAmount.Int64() / amount
		// validatorAmount := amount.Int64()

		// Send tokens to the validator
		k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, sdk.AccAddress(validatorAddr), sdk.NewCoins(sdk.NewCoin("ANDR", sdk.NewInt(validatorAmount))))
	}
}

func (k Keeper) DistributeTokensToStakers(ctx sdk.Context, amount int64) {
	// Implement logic to distribute tokens to stakers
	// daoParams := k.GetDAOParams(ctx)

	// // Iterate over validators and distribute tokens
	validatorIterator := k.stakingKeeper.GetBondedValidatorsByPower(ctx)
	for _, validator := range validatorIterator {

		delegations := k.stakingKeeper.GetValidatorDelegations(ctx, validator.GetOperator())
		totalDelegated := k.sumDelegations(delegations)

		for _, delegation := range delegations {
			delegatorReward := ((delegation.GetShares().MulInt64(amount)).Quo(totalDelegated))
			k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, sdk.AccAddress(delegation.GetDelegatorAddr()), sdk.NewCoins(sdk.NewCoin("ANDR", sdk.NewInt(delegatorReward.RoundInt64()))))
		}

		// Calculate the amount based on the staking ratio
		// validatorStake := validator.BondedTokens()
		// validatorAmount := validatorStake.QuoRaw(totalStakingTokens.Int64())
		// // validatorAmount := amount.Int64()

		// // Send tokens to the validator
		// k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, sdk.AccAddress(validatorAddr), sdk.NewCoins(sdk.NewCoin("ANDR", sdk.NewInt(validatorAmount.Int64()))))
	}
}

func (k Keeper) sumDelegations (delegations []stakingtypes.Delegation) (sum sdk.Dec) {
	for _, delegation := range delegations {
		sum = sum.Add(delegation.GetShares())
	}
	return sum
}
