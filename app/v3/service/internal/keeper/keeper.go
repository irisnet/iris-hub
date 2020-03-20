package keeper

import (
	"github.com/irisnet/irishub/app/v1/params"
	"github.com/irisnet/irishub/app/v3/service/internal/types"
	"github.com/irisnet/irishub/codec"
	sdk "github.com/irisnet/irishub/types"
)

// Keeper
type Keeper struct {
	storeKey sdk.StoreKey
	cdc      *codec.Codec
	bk       types.BankKeeper
	ak       types.AssetKeeper
	gk       types.GuardianKeeper

	// codespace
	codespace sdk.CodespaceType
	// params subspace
	paramSpace params.Subspace
	// metrics
	metrics *types.Metrics

	// used to map the module name to response callback
	respCallbacks map[string]types.ResponseCallback
}

// NewKeeper
func NewKeeper(
	cdc *codec.Codec,
	key sdk.StoreKey,
	bk types.BankKeeper,
	ak types.AssetKeeper,
	gk types.GuardianKeeper,
	codespace sdk.CodespaceType,
	paramSpace params.Subspace,
	metrics *types.Metrics,
) Keeper {
	keeper := Keeper{
		storeKey:   key,
		cdc:        cdc,
		bk:         bk,
		ak:         ak,
		gk:         gk,
		codespace:  codespace,
		paramSpace: paramSpace.WithTypeTable(types.ParamTypeTable()),
		metrics:    metrics,
	}

	keeper.respCallbacks = make(map[string]types.ResponseCallback)

	return keeper
}

// Codespace return the codespace
func (k Keeper) Codespace() sdk.CodespaceType {
	return k.codespace
}

// GetCdc returns the cdc
func (k Keeper) GetCdc() *codec.Codec {
	return k.cdc
}

func (k Keeper) GetMetrics() *types.Metrics {
	return k.metrics
}

// InitMetrics
func (k Keeper) InitMetrics(store sdk.KVStore) {
	iterator := k.AllActiveRequestsIterator(store)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		k.metrics.ActiveRequests.Add(1)
	}
}

// get service params from the global param store
func (k Keeper) GetParamSet(ctx sdk.Context) types.Params {
	var params types.Params
	k.paramSpace.GetParamSet(ctx, &params)
	return params
}

// set service params from the global param store
func (k Keeper) SetParamSet(ctx sdk.Context, params types.Params) {
	k.paramSpace.SetParamSet(ctx, &params)
}
