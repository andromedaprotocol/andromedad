package nibftdfvas

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/andromedaprotocol/andromedad/x/nibftdfvas/keeper"
	"github.com/andromedaprotocol/andromedad/x/nibftdfvas/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// NewHandler ...
func NewHandler(k keeper.Keeper) sdk.Handler {
	// this line is used by starport scaffolding # handler/msgServer

	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		// this line is used by starport scaffolding # 1
		case types.MsgTokenDistribution:
			return handleMsgTokenDistribution(ctx, k, msg)
		default:
			errMsg := fmt.Sprintf("unrecognized %s message type: %T", types.ModuleName, msg)
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}
	}
}

// func NewHandler(k keeper.Keeper) sdk.Handler {
// 	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
// 		ctx = ctx.WithEventManager(sdk.NewEventManager())

// 		switch msg := msg.(type) {
// 		case types.MsgTokenDistribution:
// 			return handleMsgTokenDistribution(ctx, k, msg)
// 		default:
// 			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
// 		}
// 	}
// }

func handleMsgTokenDistribution(ctx sdk.Context, k keeper.Keeper, msg types.MsgTokenDistribution) (*sdk.Result, error) {
	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}

	// Implement logic to handle token distribution message
	// You may need to validate the message and distribute tokens accordingly
	k.DistributeTokens(ctx, types.MyModuleParams{
		TokenOutflowPerBlock:     3,
		DirectToValidatorPercent: 20,
	})

	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}
