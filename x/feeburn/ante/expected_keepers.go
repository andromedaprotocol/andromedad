package ante

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// BankKeeper defines the contract needed for supply related APIs (noalias)
type BankKeeper interface {
	SendCoinsFromModuleToModule(ctx context.Context, senderModule, recipientModule string, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx context.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	BurnCoins(ctx context.Context, moduleName string, amt sdk.Coins) error
}
