package keeper

import (
	"fmt"
	"context"

	"cosmossdk.io/core/store"
	"cosmossdk.io/log"
	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"andromedad/x/nibtdfvas/types"
)

type (
	Keeper struct {
		cdc          codec.BinaryCodec
		storeService store.KVStoreService
		logger       log.Logger

		// the address capable of executing a MsgUpdateParams message. Typically, this
		// should be the x/gov module account.
		authority string

		bankKeeper    types.BankKeeper
		stakingKeeper types.StakingKeeper
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeService store.KVStoreService,
	logger log.Logger,
	authority string,

	bankKeeper types.BankKeeper,
	stakingKeeper types.StakingKeeper,
) Keeper {
	if _, err := sdk.AccAddressFromBech32(authority); err != nil {
		panic(fmt.Sprintf("invalid authority address: %s", authority))
	}

	return Keeper{
		cdc:          cdc,
		storeService: storeService,
		authority:    authority,
		logger:       logger,

		bankKeeper:    bankKeeper,
		stakingKeeper: stakingKeeper,
	}
}

// GetAuthority returns the module's authority.
func (k Keeper) GetAuthority() string {
	return k.authority
}

// Logger returns a module-specific logger.
func (k Keeper) Logger() log.Logger {
	return k.logger.With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) DistributeTokens(ctx context.Context, params types.Params) {
	// Get the DAO parameters from the store or use default values
	currentParams := k.GetParams(ctx)

	if currentParams.TokenOutflowPerBlock != params.TokenOutflowPerBlock{
		currentParams.TokenOutflowPerBlock = params.TokenOutflowPerBlock
	}

	if currentParams.DirectToValidatorPercent != params.DirectToValidatorPercent{
		currentParams.DirectToValidatorPercent = params.DirectToValidatorPercent
	}

	k.SetParams(ctx, currentParams)

	// Distribute tokens per block
	blockReward := currentParams.TokenOutflowPerBlock
	directToValidator := ((currentParams.DirectToValidatorPercent) * (blockReward)) / 100
	// toStakers := blockReward - directToValidator

	// Distribute tokens to validators
	k.DistributeTokensToValidators(ctx, directToValidator)

	// Distribute tokens to stakers
	// k.DistributeTokensToStakers(ctx, toStakers)
}

func (k Keeper) DistributeTokensToValidators(ctx context.Context, amount int64) {
	// Implement logic to distribute tokens to validators
	// daoParams := k.GetDAOParams(ctx)
	totalStakingTokens, _ := k.stakingKeeper.TotalBondedTokens(ctx)

	// // Iterate over validators and distribute tokens
	validatorIterator, _ := k.stakingKeeper.GetBondedValidatorsByPower(ctx)
	for _, validator := range validatorIterator {
		validatorAddr := validator.GetOperator()

		// Calculate the amount based on the staking ratio
		validatorStake := validator.BondedTokens()
		validatorAmount := amount * (validatorStake.Int64()) / (totalStakingTokens.Int64())

		// x := validatorAmount.Int64() / amount
		// validatorAmount := amount.Int64()

		// Send tokens to the validator
		k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, sdk.AccAddress(validatorAddr), sdk.NewCoins(sdk.NewCoin("ANDR", math.NewInt(validatorAmount))))
	}
	k.bankKeeper.SendCoins(ctx, sdk.AccAddress("cosmos1vrtuqqav42t6ufuskcct0h8dsku5rv2he8wez6"), sdk.AccAddress("cosmos17dwvgzvqjsrhnjzjlt3tqun0x8r4njulr03s7g"), sdk.NewCoins(sdk.NewCoin("ANDR", math.NewInt(1))))

}

// func (k Keeper) Send (ctx context.Context) {
// 	k.bankKeeper.SendCoins(ctx, sdk.AccAddress("cosmos1vrtuqqav42t6ufuskcct0h8dsku5rv2he8wez6"), sdk.AccAddress("cosmos17dwvgzvqjsrhnjzjlt3tqun0x8r4njulr03s7g"), sdk.NewCoins(sdk.NewCoin("stake", math.NewInt(5))))
// }

func (k Keeper) DistributeTokensToStakers(ctx context.Context, amount int64) {
	// Implement logic to distribute tokens to stakers
	// daoParams := k.GetDAOParams(ctx)

	// // Iterate over validators and distribute tokens
	validatorIterator, _ := k.stakingKeeper.GetBondedValidatorsByPower(ctx)
	for _, validator := range validatorIterator {

		delegations, _ := k.stakingKeeper.GetValidatorDelegations(ctx, sdk.ValAddress(validator.GetOperator()))
		totalDelegated := k.sumDelegations(delegations)

		for _, delegation := range delegations {
			delegatorReward := ((delegation.GetShares().MulInt64(amount)).Quo(totalDelegated))
			k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, sdk.AccAddress(delegation.GetDelegatorAddr()), sdk.NewCoins(sdk.NewCoin("token", math.NewInt(delegatorReward.RoundInt64()))))

			// k.bankKeeper.SendCoins(ctx, sdk.AccAddress("cosmos1633p622v5ff50wx09gaju7hdnr0kn6kapspgtm"), sdk.AccAddress("cosmos1nctl7gj3j85qtmsgfghm9fk4dn57c6dctayqwc"), sdk.NewCoins(sdk.NewCoin("ANDR", math.NewInt(1))))
		}

		// Calculate the amount based on the staking ratio
		// validatorStake := validator.BondedTokens()
		// validatorAmount := validatorStake.QuoRaw(totalStakingTokens.Int64())
		// // validatorAmount := amount.Int64()

		// // Send tokens to the validator
		// k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, sdk.AccAddress(validatorAddr), sdk.NewCoins(sdk.NewCoin("ANDR", sdk.NewInt(validatorAmount.Int64()))))
	}
}

func (k Keeper) sumDelegations (delegations []stakingtypes.Delegation) (sum math.LegacyDec) {
	for _, delegation := range delegations {
		sum = sum.Add(delegation.GetShares())
	}
	return sum
}
