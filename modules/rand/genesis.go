package rand

import (
	"fmt"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis stores genesis data
func InitGenesis(ctx sdk.Context, k Keeper, data GenesisState) {
	for height, requests := range data.PendingRandRequests {
		for _, request := range requests {
			h, err := strconv.ParseInt(height, 10, 64)
			if err != nil {
				continue
			}

			reqID := GenerateRequestID(request)
			k.EnqueueRandRequest(ctx, h, reqID, request)
		}
	}
}

// ExportGenesis outputs genesis data
func ExportGenesis(ctx sdk.Context, k Keeper) GenesisState {
	pendingRequests := make(map[string][]Request)

	k.IterateRandRequestQueue(ctx, func(height int64, request Request) bool {
		leftHeight := fmt.Sprintf("%d", height-ctx.BlockHeight()+1)
		pendingRequests[leftHeight] = append(pendingRequests[leftHeight], request)

		return false
	})

	return GenesisState{
		PendingRandRequests: pendingRequests,
	}
}

// DefaultGenesisState gets the default genesis state
func DefaultGenesisState() GenesisState {
	return GenesisState{
		PendingRandRequests: map[string][]Request{},
	}
}

// DefaultGenesisStateForTest gets the default genesis state for test
func DefaultGenesisStateForTest() GenesisState {
	return GenesisState{
		PendingRandRequests: map[string][]Request{},
	}
}
