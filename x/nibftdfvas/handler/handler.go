package handler

import (
	"fmt"

	"github.com/andromedaprotocol/andromedad/x/nibtdfvas/keeper"
	"github.com/andromedaprotocol/andromedad/x/nibtdfvas/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// NewHandler ...
// func NewHandler(k keeper.Keeper) sdk.Handler {
// 	// this line is used by starport scaffolding # handler/msgServer

// 	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
// 		ctx = ctx.WithEventManager(sdk.NewEventManager())

// 		switch msg := msg.(type) {
// 		// this line is used by starport scaffolding # 1
// 		default:
// 			errMsg := fmt.Sprintf("unrecognized %s message type: %T", types.ModuleName, msg)
// 			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
// 		}
// 	}
// }

func NewHandler(k keeper.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		switch msg := msg.(type) {
		case types.MsgTokenDistribution:
			return handleMsgTokenDistribution(ctx, k, msg)
		default:
			return nil, sdk.ErrUnknownRequest("unknown nibtdfvas message")
		}
	}
}

func handleMsgTokenDistribution(ctx sdk.Context, k keeper.Keeper, msg types.MsgTokenDistribution) (*sdk.Result, error) {
	// Implement logic to handle token distribution message
	// need to validate the message and distribute tokens accordingly
	k.DistributeTokens(ctx, types.DAOParams{
		TokenOutflowPerBlock: 3,
		DirectToValidatorPercent: 20,
	})

	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}