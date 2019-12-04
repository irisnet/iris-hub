package keeper

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/irisnet/irishub/modules/rand/internal/types"
)

type Keeper struct {
	storeKey sdk.StoreKey
	cdc      *codec.Codec

	// codespace
	codespace sdk.CodespaceType
}

func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, codespace sdk.CodespaceType) Keeper {
	return Keeper{
		storeKey:  key,
		cdc:       cdc,
		codespace: codespace,
	}
}

// Codespace returns the codespace
func (k Keeper) Codespace() sdk.CodespaceType {
	return k.codespace
}

// GetCdc returns the cdc
func (k Keeper) GetCdc() *codec.Codec {
	return k.cdc
}

// RequestRand requests a random number
func (k Keeper) RequestRand(ctx sdk.Context, consumer sdk.AccAddress, blockInterval uint64) (types.Request, sdk.Error) {
	currentHeight := ctx.BlockHeight()
	destHeight := currentHeight + int64(blockInterval)

	// get tx hash
	txHash := types.SHA256(ctx.TxBytes())

	// build request
	request := types.NewRequest(currentHeight, consumer, txHash)

	// generate the request id
	reqID := types.GenerateRequestID(request)

	// add to the queue
	k.EnqueueRandRequest(ctx, destHeight, reqID, request)

	return request, nil
}

// SetRand stores the random number
func (k Keeper) SetRand(ctx sdk.Context, reqID []byte, rand types.Rand) {
	store := ctx.KVStore(k.storeKey)

	bz := k.cdc.MustMarshalBinaryLengthPrefixed(rand)
	store.Set(types.KeyRand(reqID), bz)
}

// EnqueueRandRequest enqueue the random number request
func (k Keeper) EnqueueRandRequest(ctx sdk.Context, height int64, reqID []byte, request types.Request) {
	store := ctx.KVStore(k.storeKey)

	bz := k.cdc.MustMarshalBinaryLengthPrefixed(request)
	store.Set(types.KeyRandRequestQueue(height, reqID), bz)
}

// DequeueRandRequest removes the random number request by the specified height and request id
func (k Keeper) DequeueRandRequest(ctx sdk.Context, height int64, reqID []byte) {
	store := ctx.KVStore(k.storeKey)

	// delete the key
	store.Delete(types.KeyRandRequestQueue(height, reqID))
}

// GetRand retrieves the random number by the specified request id
func (k Keeper) GetRand(ctx sdk.Context, reqID []byte) (types.Rand, sdk.Error) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.KeyRand(reqID))
	if bz == nil {
		return types.Rand{}, types.ErrInvalidReqID(k.codespace, fmt.Sprintf("request id does not exist: %s", hex.EncodeToString(reqID)))
	}

	var rand types.Rand
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &rand)

	return rand, nil
}

// IterateRands iterates through all the random numbers
func (k Keeper) IterateRands(ctx sdk.Context, op func(r types.Rand) (stop bool)) {
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, types.PrefixRand)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var rand types.Rand
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &rand)

		if stop := op(rand); stop {
			break
		}
	}
}

// IterateRandRequestQueueByHeight iterates the random number request queue by the specified height
func (k Keeper) IterateRandRequestQueueByHeight(ctx sdk.Context, height int64) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, types.KeyRandRequestQueueSubspace(height))
}

// IterateRandRequestQueue iterates through the random number request queue
func (k Keeper) IterateRandRequestQueue(ctx sdk.Context, op func(h int64, r types.Request) (stop bool)) {
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, types.PrefixRandRequestQueue)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		keyParts := bytes.Split(iterator.Key(), types.KeyDelimiter)
		height, _ := strconv.ParseInt(string(keyParts[1]), 10, 64)

		var request types.Request
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &request)

		if stop := op(height, request); stop {
			break
		}
	}
}
