package keeper

import (
	"fmt"

	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	// auth "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	// bank "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	// staking "github.com/cosmos/cosmos-sdk/x/staking/keeper"

	// "andromedad/x/nibtdfvas/types"
	"github.com/andromedaprotocol/andromedad/x/nibtdfvas/types"

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
