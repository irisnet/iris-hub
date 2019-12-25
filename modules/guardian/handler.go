package guardian

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewHandler returns a handler for all "guardian" type messages.
func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case MsgAddProfiler:
			return handleMsgAddProfiler(ctx, k, msg)
		case MsgAddTrustee:
			return handleMsgAddTrustee(ctx, k, msg)
		case MsgDeleteProfiler:
			return handleMsgDeleteProfiler(ctx, k, msg)
		case MsgDeleteTrustee:
			return handleMsgDeleteTrustee(ctx, k, msg)
		default:
			errMsg := fmt.Sprintf("unrecognized guardian message type: %T", msg)
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// handleMsgAddProfiler handles MsgAddProfiler
func handleMsgAddProfiler(ctx sdk.Context, k Keeper, msg MsgAddProfiler) sdk.Result {
	if profiler, found := k.GetProfiler(ctx, msg.AddedBy); !found || profiler.GetAccountType() != Genesis {
		return ErrInvalidOperator(DefaultCodespace, msg.AddedBy).Result()
	}
	if _, found := k.GetProfiler(ctx, msg.Address); found {
		return ErrProfilerExists(DefaultCodespace, msg.Address).Result()
	}
	profiler := NewGuardian(msg.Description, Ordinary, msg.Address, msg.AddedBy)
	k.AddProfiler(ctx, profiler)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.AddedBy.String()),
		),
		sdk.NewEvent(
			EventTypeAddProfiler,
			sdk.NewAttribute(AttributeKeyProfilerAddress, msg.Address.String()),
			sdk.NewAttribute(AttributeKeyAddedBy, msg.AddedBy.String()),
		),
	})

	return sdk.Result{Events: ctx.EventManager().Events()}
}

// handleMsgAddTrustee handles MsgAddTrustee
func handleMsgAddTrustee(ctx sdk.Context, k Keeper, msg MsgAddTrustee) sdk.Result {
	if trustee, found := k.GetTrustee(ctx, msg.AddedBy); !found || trustee.GetAccountType() != Genesis {
		return ErrInvalidOperator(DefaultCodespace, msg.AddedBy).Result()
	}
	if _, found := k.GetTrustee(ctx, msg.Address); found {
		return ErrTrusteeExists(DefaultCodespace, msg.Address).Result()
	}
	trustee := NewGuardian(msg.Description, Ordinary, msg.Address, msg.AddedBy)
	k.AddTrustee(ctx, trustee)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.AddedBy.String()),
		),
		sdk.NewEvent(
			EventTypeAddTrustee,
			sdk.NewAttribute(AttributeKeyTrusteeAddress, msg.Address.String()),
			sdk.NewAttribute(AttributeKeyAddedBy, msg.AddedBy.String()),
		),
	})

	return sdk.Result{Events: ctx.EventManager().Events()}
}

// handleMsgDeleteProfiler handles MsgDeleteProfiler
func handleMsgDeleteProfiler(ctx sdk.Context, k Keeper, msg MsgDeleteProfiler) sdk.Result {
	if profiler, found := k.GetProfiler(ctx, msg.DeletedBy); !found || profiler.GetAccountType() != Genesis {
		return ErrInvalidOperator(DefaultCodespace, msg.DeletedBy).Result()
	}
	profiler, found := k.GetProfiler(ctx, msg.Address)
	if !found {
		return ErrProfilerNotExists(DefaultCodespace, msg.Address).Result()
	}
	if profiler.GetAccountType() == Genesis {
		return ErrDeleteGenesisProfiler(DefaultCodespace, msg.Address).Result()
	}

	k.DeleteProfiler(ctx, msg.Address)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.DeletedBy.String()),
		),
		sdk.NewEvent(
			EventTypeDeleteProfiler,
			sdk.NewAttribute(AttributeKeyProfilerAddress, msg.Address.String()),
			sdk.NewAttribute(AttributeKeyDeletedBy, msg.DeletedBy.String()),
		),
	})

	return sdk.Result{Events: ctx.EventManager().Events()}
}

// handleMsgDeleteTrustee handles MsgDeleteTrustee
func handleMsgDeleteTrustee(ctx sdk.Context, k Keeper, msg MsgDeleteTrustee) sdk.Result {
	if trustee, found := k.GetTrustee(ctx, msg.DeletedBy); !found || trustee.GetAccountType() != Genesis {
		return ErrInvalidOperator(DefaultCodespace, msg.DeletedBy).Result()
	}
	trustee, found := k.GetTrustee(ctx, msg.Address)
	if !found {
		return ErrTrusteeNotExists(DefaultCodespace, msg.Address).Result()
	}
	if trustee.GetAccountType() == Genesis {
		return ErrDeleteGenesisTrustee(DefaultCodespace, msg.Address).Result()
	}

	k.DeleteTrustee(ctx, msg.Address)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.DeletedBy.String()),
		),
		sdk.NewEvent(
			EventTypeDeleteTrustee,
			sdk.NewAttribute(AttributeKeyTrusteeAddress, msg.Address.String()),
			sdk.NewAttribute(AttributeKeyDeletedBy, msg.DeletedBy.String()),
		),
	})

	return sdk.Result{Events: ctx.EventManager().Events()}
}
