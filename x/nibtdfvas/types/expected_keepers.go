package types

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	abci "github.com/cometbft/cometbft/abci/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// AccountKeeper defines the expected account keeper used for simulations (noalias)
type AccountKeeper interface {
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) types.AccountI
	GetModuleAccount(ctx sdk.Context, name string) authtypes.ModuleAccountI
	// Methods imported from account should be defined here
}

// BankKeeper defines the expected interface needed to retrieve account balances.
type BankKeeper interface {
	SpendableCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
	// Methods imported from bank should be defined here
	MintCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) error
	SendCoinsFromModuleToAccount(
		ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins,
	) error
	SendCoinsFromModuleToModule(
		ctx sdk.Context, senderModule, recipientModule string, amt sdk.Coins,
	) error
	GetBalance(ctx sdk.Context, addr sdk.AccAddress, denom string) sdk.Coin
}

type StakingKeeper interface {
	BlockValidatorUpdates(ctx sdk.Context) []abci.ValidatorUpdate
	BondDenom(ctx sdk.Context) string
	GetBondedValidatorsByPower(ctx sdk.Context) []stakingtypes.Validator
	TotalBondedTokens(ctx sdk.Context) math.Int
	GetValidatorDelegations(ctx sdk.Context, valAddr sdk.ValAddress) []stakingtypes.Delegation
}