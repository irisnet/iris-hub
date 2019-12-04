package keeper

import (
	"encoding/hex"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/irisnet/irishub/modules/rand/internal/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, sdk.Error) {
		switch path[0] {
		case types.QueryRand:
			return queryRand(ctx, req, k)
		case types.QueryRandRequestQueue:
			return queryRandRequestQueue(ctx, req, k)
		default:
			return nil, sdk.ErrUnknownRequest("unknown rand query endpoint")
		}
	}
}

func queryRand(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params types.QueryRandParams
	err := keeper.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}

	reqID, err := hex.DecodeString(params.ReqID)
	if err != nil {
		return nil, types.ErrInvalidReqID(types.DefaultCodespace, fmt.Sprintf("invalid request id: %s", err))
	}

	rand, err2 := keeper.GetRand(ctx, reqID)
	if err2 != nil {
		return nil, err2
	}

	bz, err := codec.MarshalJSONIndent(keeper.cdc, rand)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}

	return bz, nil
}

func queryRandRequestQueue(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params types.QueryRandRequestQueueParams
	err := keeper.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}

	if params.Height < 0 {
		return nil, types.ErrInvalidHeight(types.DefaultCodespace, fmt.Sprintf("the height must not be less than 0: %d", params.Height))
	}

	var requests []types.Request

	if params.Height == 0 {
		// query all pending requests
		requests = queryAllRandRequestsInQueue(ctx, keeper)
	} else {
		// query the pending requests by the specified height
		requests = queryRandRequestQueueByHeight(ctx, params.Height, keeper)
	}

	bz, err := codec.MarshalJSONIndent(keeper.cdc, requests)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}

	return bz, nil
}

func queryRandRequestQueueByHeight(ctx sdk.Context, height int64, keeper Keeper) []types.Request {
	var requests = make([]types.Request, 0)

	iterator := keeper.IterateRandRequestQueueByHeight(ctx, height)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var request types.Request
		keeper.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &request)

		requests = append(requests, request)
	}

	return requests
}

func queryAllRandRequestsInQueue(ctx sdk.Context, keeper Keeper) []types.Request {
	var requests = make([]types.Request, 0)

	keeper.IterateRandRequestQueue(ctx, func(h int64, r types.Request) (stop bool) {
		requests = append(requests, r)
		return false
	})

	return requests
}
