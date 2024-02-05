package keeper

import (
	"fmt"

	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	// banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	// auth "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	// bank "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	// staking "github.com/cosmos/cosmos-sdk/x/staking/keeper"

	// "andromedad/x/nibtdfvas/types"
	"github.com/andromedaprotocol/andromedad/x/nibtdfvas/types"
	"cosmossdk.io/math"
)

type (
	Keeper struct {
		cdc      codec.BinaryCodec
		storeKey storetypes.StoreKey
		memKey   storetypes.StoreKey
		// the address capable of executing a MsgUpdateParams message. Typically, this should be the x/gov module account.
		authority sdk.AccAddress
		accountKeeper    types.AccountKeeper
		bankKeeper    types.BankKeeper
		stakingKeeper types.StakingKeeper
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey,
	memKey storetypes.StoreKey,
	authority sdk.AccAddress,
	accountKeeper    types.AccountKeeper,
	bankKeeper    types.BankKeeper,
	stakingKeeper types.StakingKeeper,
) *Keeper {
	// ensure gov module account is set and is not nil
	if err := sdk.VerifyAddressFormat(authority); err != nil {
		panic(err)
	}

	return &Keeper{
		cdc:       cdc,
		storeKey:  storeKey,
		memKey:    memKey,
		authority: authority,
		accountKeeper: accountKeeper,
		bankKeeper: bankKeeper,
		stakingKeeper: stakingKeeper,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) GetNibtdfvasAccount(ctx sdk.Context) authtypes.ModuleAccountI {
	return k.accountKeeper.GetModuleAccount(ctx, types.ModuleName)
}

func (k Keeper) Mint(ctx sdk.Context, moduleName string, mintCoins sdk.Coins) error {
	if !mintCoins.IsValid() {
		return fmt.Errorf("invalid mint coins: %s", mintCoins.String())
	}
	return k.bankKeeper.MintCoins(ctx, moduleName, mintCoins)
}

func (k Keeper) SendCoinsFromModuleToModule(ctx sdk.Context, mintCoins sdk.Coins) error {
	return k.bankKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, types.ModuleName, mintCoins)
}

func (k Keeper) ValidatorUpdate(ctx sdk.Context) []abci.ValidatorUpdate {
	return k.stakingKeeper.BlockValidatorUpdates(ctx)
}

func (k Keeper) DistributeTokens(ctx sdk.Context, req abci.RequestBeginBlock) {


	// // determine the total power signing the block
	// var previousTotalPower int64

	// for _, voteInfo := range req.LastCommitInfo.GetVotes() {
	// 	previousTotalPower += voteInfo.Validator.Power
	// }

	denom := k.stakingKeeper.BondDenom(ctx)

	currentParams := k.GetParams(ctx)


	blockReward := math.LegacyNewDec(currentParams.TokenOutflowPerBlock)
	// directToValidator := ((currentParams.DirectToValidatorPercent) * (blockReward)) / 100
	directToValidator := (math.LegacyNewDec(currentParams.DirectToValidatorPercent).Quo(blockReward)).Quo(math.LegacyNewDec(100))
	toStakers := blockReward.Sub(directToValidator)

	moduleAcc := k.accountKeeper.GetModuleAccount(ctx, types.ModuleName)

	balanceOfModule := k.bankKeeper.GetBalance(ctx, moduleAcc.GetAddress(), denom)

	if !(blockReward.GT(math.LegacyNewDecFromInt(balanceOfModule.Amount))) {

		k.DistributeTokensToValidators(ctx, req, directToValidator)

		k.DistributeTokensToStakers(ctx, req, toStakers)
	} else {
		k.Logger(ctx).Error("Insufficient balance for token distribution")
	}
}


func (k Keeper) DistributeTokensToValidators(ctx sdk.Context, req abci.RequestBeginBlock, amount math.LegacyDec) {

	denom := k.stakingKeeper.BondDenom(ctx)

	totalStakingTokens := k.stakingKeeper.TotalBondedTokens(ctx)

	if totalStakingTokens.IsZero() {
			k.Logger(ctx).Error("Total staking tokens are zero, cannot distribute rewards")
			return
	}

	// Iterate over validators and distribute tokens
	validatorIterator := k.stakingKeeper.GetBondedValidatorsByPower(ctx)

	for _, validator := range validatorIterator {

		validatorAddr := validator.GetOperator()

		// Calculate the amount based on the staking ratio
		validatorStake := validator.BondedTokens()

		// validatorAmount := amount * (validatorStake.Int64()) / (totalStakingTokens.Int64())
		validatorAmount := (amount.Mul(math.LegacyNewDecFromInt(validatorStake))).Quo(math.LegacyNewDecFromInt(totalStakingTokens))

		k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, sdk.AccAddress(validatorAddr), sdk.NewCoins(sdk.NewCoin(denom, validatorAmount.TruncateInt())))
	}
}


func (k Keeper) DistributeTokensToStakers(ctx sdk.Context, req.RequestBeginBlock, amount math.LegacyDec) {

	denom := k.stakingKeeper.BondDenom(ctx)

	validatorIterator := k.stakingKeeper.GetBondedValidatorsByPower(ctx)

	for _, validator := range validatorIterator {
		delegations := k.stakingKeeper.GetValidatorDelegations(ctx, validator.GetOperator())

		for _, delegation := range delegations {
			delegatorAddr := delegation.GetDelegatorAddr()
			delegatorShares := delegation.GetShares()
			toDelegatorAmount := delegatorShares.Mul(amount)
			k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, delegatorAddr, sdk.NewCoins(sdk.NewCoin(denom, toDelegatorAmount.TruncateInt())))
		}
	}
}