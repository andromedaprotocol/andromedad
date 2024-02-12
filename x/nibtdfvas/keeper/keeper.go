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
		// memKey   storetypes.StoreKey
		// the address capable of executing a MsgUpdateParams message. Typically, this should be the x/gov module account.
		authority sdk.AccAddress
		accountKeeper    types.AccountKeeper
		bankKeeper    types.BankKeeper
		stakingKeeper types.StakingKeeper
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
	// memKey storetypes.StoreKey,
	authority sdk.AccAddress,
	accountKeeper    types.AccountKeeper,
	bankKeeper    types.BankKeeper,
	stakingKeeper types.StakingKeeper,
) Keeper {
	// ensure gov module account is set and is not nil
	if err := sdk.VerifyAddressFormat(authority); err != nil {
		panic(err)
	}

	return Keeper{
		cdc:       cdc,
		storeKey:  storeKey,
		// memKey:    memKey,
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

func (k Keeper) MintCoins(ctx sdk.Context, moduleName string, mintCoins sdk.Coins) error {
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

	// denom := k.stakingKeeper.BondDenom(ctx)

	// k.Logger(ctx).Error(denom)

	currentParams := k.GetParams(ctx)


	blockReward := math.LegacyNewDec(currentParams.TokenOutflowPerBlock)

	// k.Logger(ctx).Error(string(blockReward.RoundInt64()))
	// directToValidator := ((currentParams.DirectToValidatorPercent) * (blockReward)) / 100
	directToValidator := (math.LegacyNewDec(currentParams.DirectToValidatorPercent).Mul(blockReward)).Quo(math.LegacyNewDec(100))
	toStakers := blockReward.Sub(directToValidator)

	moduleAddress := k.accountKeeper.GetModuleAddress(types.ModuleName)

	balanceOfModule := k.bankKeeper.GetBalance(ctx, moduleAddress, "andr")

	// k.Logger(ctx).Error(moduleAcc.GetAddress().String())

	// k.Logger(ctx).Error(string(balanceOfModule.Amount.Int64()))

	if !(blockReward.GT(math.LegacyNewDecFromInt(balanceOfModule.Amount))) {

		k.DistributeTokensToValidators(ctx, req, directToValidator)

		k.DistributeTokensToStakers(ctx, req, toStakers)
	} else {
		k.Logger(ctx).Error("Insufficient balance for token distribution")
	}
}


func (k Keeper) DistributeTokensToValidators(ctx sdk.Context, req abci.RequestBeginBlock, amount math.LegacyDec) {

	// determine the total power signing the block
	var previousTotalPower int64

	for _, voteInfo := range req.LastCommitInfo.GetVotes() {
		previousTotalPower += voteInfo.Validator.Power
	}

	bondedVotes := req.LastCommitInfo.GetVotes()	

	// denom := k.stakingKeeper.BondDenom(ctx)

	for _, vote := range bondedVotes {

		validator, _ := k.stakingKeeper.GetValidatorByConsAddr(ctx, vote.Validator.Address)

		// fmt.Println("validator value:", validator)

		powerFraction := math.LegacyNewDec(vote.Validator.Power).QuoTruncate(math.LegacyNewDec(previousTotalPower))

		rewardToValidator := amount.MulTruncate(powerFraction)
		fmt.Println("rewardToValidator value:", rewardToValidator)

		k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, sdk.AccAddress(validator.GetOperator()), sdk.NewCoins(sdk.NewCoin("andr", rewardToValidator.TruncateInt())))
		moduleAddress := k.accountKeeper.GetModuleAddress(types.ModuleName)
		fmt.Println("balance value:", k.bankKeeper.GetBalance(ctx, moduleAddress, "andr"))
		fmt.Println("balance value:", k.bankKeeper.GetBalance(ctx, sdk.AccAddress(validator.GetOperator()), "andr"))
	} 
}


func (k Keeper) DistributeTokensToStakers(ctx sdk.Context, req abci.RequestBeginBlock, amount math.LegacyDec) {

	var previousTotalPower int64

	for _, voteInfo := range req.LastCommitInfo.GetVotes() {
		previousTotalPower += voteInfo.Validator.Power
	}

	bondedVotes := req.LastCommitInfo.GetVotes()
	// // denom := k.stakingKeeper.BondDenom(ctx)

	for _, vote := range bondedVotes {

		validator, _ := k.stakingKeeper.GetValidatorByConsAddr(ctx, vote.Validator.Address)

		powerFraction := math.LegacyNewDec(vote.Validator.Power).QuoTruncate(math.LegacyNewDec(previousTotalPower))

		rewardToOneValidatorDelegators := amount.MulTruncate(powerFraction)
		
		delegations := k.stakingKeeper.GetValidatorDelegations(ctx, validator.GetOperator())

		totalShares := math.LegacyNewDec(0)

		for _, delegation := range delegations {
			totalShares = totalShares.Add(delegation.GetShares())
		}

		for _, delegation := range delegations {

			delegatorAddr := delegation.GetDelegatorAddr()

			delegatorShares := delegation.GetShares()

			toDelegatorAmount := (delegatorShares.QuoTruncate(totalShares)).Mul(rewardToOneValidatorDelegators)

			k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, delegatorAddr, sdk.NewCoins(sdk.NewCoin("andr", toDelegatorAmount.TruncateInt())))
		}
	}
}