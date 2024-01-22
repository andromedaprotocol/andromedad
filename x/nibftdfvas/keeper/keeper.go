package keeper

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/andromedaprotocol/andromedad/x/nibtdfvas/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

type (
	Keeper struct {
		cdc        codec.BinaryCodec
		storeKey   sdk.StoreKey
		memKey     sdk.StoreKey
		paramstore paramtypes.Subspace
		bankKeeper bank.Keeper
		authKeeper auth.AccountKeeper
		paramsKeeper  sdk.ParamsKeeper
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey sdk.StoreKey,
	memKey sdk.StoreKey,
	ps paramtypes.Subspace,
	bankKeeper bank.Keeper,
	authKeeper auth.AccountKeeper
	paramsKeeper  sdk.ParamsKeeper
) *Keeper {
	// set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	return &Keeper{
		cdc:        cdc,
		storeKey:   storeKey,
		memKey:     memKey,
		paramstore: ps,
		bankKeeper: bankKeeper,
		authKeeper: authKeeper,
		paramsKeeper: paramsKeeper
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) DistributeTokens(ctx sdk.Context, params types.MyModuleParams) {
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
	directToValidator := (currentParams.DirectToValidatorPercent * blockReward) / 100
	toStakers := blockReward - directToValidator

	// Distribute tokens to validators
	k.DistributeTokensToValidators(ctx, directToValidator)

	// Distribute tokens to stakers
	k.DistributeTokensToValidators(ctx, toStakers)
}

func (k Keeper) DistributeTokensToValidators(ctx sdk.Context, amount int) {
	// Implement logic to distribute tokens to validators
	params := k.GetParams(ctx)

	// Iterate over validators and distribute tokens
	validatorIterator := k.bankKeeper.GetStakingKeeper().validatorIterator(ctx)
	for ; validatorIterator.Valid(); validatorIterator.Next() {
		validatorAddr := validatorIterator.Key()

		// Implement your logic to calculate the amount for each validator
		validatorAmount := amount / validatorIterator.Count()

		// Send tokens to the validator
		k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, sdk.AccAddress(validatorAddr), sdk.NewCoins(sdk.NewCoin("ANDR", sdk.NewInt(int64(validatorAmount)))))
	}
}

func (k Keeper) DistributeTokensToStakers(ctx sdk.Context, amount int) {
	// Implement logic to distribute tokens to stakers
	params := k.GetParams(ctx)

	// Iterate over stakers and distribute tokens
	stakerIterator := k.bankKeeper.GetStakingKeeper().StakeIterator(ctx)
	for ; stakerIterator.Valid(); stakerIterator.Next() {
		stakerAddr := stakerIterator.Key()

		// Implement your logic to calculate the amount for each staker
		stakerAmount := amount / stakerIterator.Count()

		// Send tokens to the staker
		k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, sdk.AccAddress(stakerAddr), sdk.NewCoins(sdk.NewCoin("ANDR", sdk.NewInt(int64(stakerAmount)))))
	}
}

// GetValidators returns the list of validators
func (k Keeper) GetValidators(ctx sdk.Context) []sdk.ValAddress {
	// Implement logic to get validators from the state
	return k.bankKeeper.GetStakingKeeper().GetValidators(ctx)
}

// GetStakers returns the list of stakers
func (k Keeper) GetStakers(ctx sdk.Context) []sdk.AccAddress {
	// Implement logic to get stakers from the state
	return k.bankKeeper.GetStakingKeeper().GetAllSDKAddresses(ctx)
}

// GetDistributionAuthority returns the authority address for initiating token distribution
func (k Keeper) GetDistributionAuthority(ctx sdk.Context) sdk.AccAddress {
	// Implement logic to get the distribution authority address from the state
	return k.paramsKeeper.GetDistributionAuthority(ctx)
}
