package keeper

import (
	sdk "github.com/irisnet/irishub/types"
	"github.com/irisnet/irishub/modules/distribution/types"
)

// check whether a validator has distribution info
func (k Keeper) HasValidatorDistInfo(ctx sdk.Context,
	operatorAddr sdk.ValAddress) (exists bool) {
	store := ctx.KVStore(k.storeKey)
	return store.Has(GetValidatorDistInfoKey(operatorAddr))
}

// get the validator distribution info
func (k Keeper) GetValidatorDistInfo(ctx sdk.Context,
	operatorAddr sdk.ValAddress) (vdi types.ValidatorDistInfo) {

	store := ctx.KVStore(k.storeKey)

	b := store.Get(GetValidatorDistInfoKey(operatorAddr))
	if b == nil {
		panic("Stored validator-distribution info should not have been nil")
	}

	k.cdc.MustUnmarshalBinaryLengthPrefixed(b, &vdi)
	return
}

// set the validator distribution info
func (k Keeper) SetValidatorDistInfo(ctx sdk.Context, vdi types.ValidatorDistInfo) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshalBinaryLengthPrefixed(vdi)
	store.Set(GetValidatorDistInfoKey(vdi.OperatorAddr), b)
}

// remove a validator distribution info
func (k Keeper) RemoveValidatorDistInfo(ctx sdk.Context, valAddr sdk.ValAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(GetValidatorDistInfoKey(valAddr))
}

// Get the calculated accum of a validator at the current block
// without affecting the state.
func (k Keeper) GetValidatorAccum(ctx sdk.Context, operatorAddr sdk.ValAddress) (sdk.Dec, sdk.Error) {
	if !k.HasValidatorDistInfo(ctx, operatorAddr) {
		return sdk.Dec{}, types.ErrNoValidatorDistInfo(k.codespace)
	}

	// withdraw self-delegation
	height := ctx.BlockHeight()
	lastValPower := k.stakeKeeper.GetLastValidatorPower(ctx, operatorAddr)
	valInfo := k.GetValidatorDistInfo(ctx, operatorAddr)
	accum := valInfo.GetValAccum(height, sdk.NewDecFromInt(lastValPower))

	return accum, nil
}

// withdrawal all the validator rewards including the commission
func (k Keeper) WithdrawValidatorRewardsAll(ctx sdk.Context, operatorAddr sdk.ValAddress) (types.DecCoins, sdk.Tags, sdk.Error) {

	if !k.HasValidatorDistInfo(ctx, operatorAddr) {
		return nil, nil, types.ErrNoValidatorDistInfo(k.codespace)
	}

	// withdraw self-delegation
	accAddr := sdk.AccAddress(operatorAddr.Bytes())
	withdraw, resultTags := k.withdrawDelegationRewardsAll(ctx, accAddr)

	// withdrawal validator commission rewards
	valInfo := k.GetValidatorDistInfo(ctx, operatorAddr)
	wc := k.GetWithdrawContext(ctx, operatorAddr)
	valInfo, feePool, commission := valInfo.WithdrawCommission(wc)
	withdraw = withdraw.Plus(commission)
	k.SetValidatorDistInfo(ctx, valInfo)
	resultTags = resultTags.AppendTag(sdk.TagRewardCommission, []byte(commission.ToString()))
	k.WithdrawToDelegator(ctx, feePool, accAddr, withdraw)
	return withdraw, resultTags, nil
}

// get all the validator rewards including the commission
func (k Keeper) CurrentValidatorRewardsAll(ctx sdk.Context, operatorAddr sdk.ValAddress) (sdk.Coins, sdk.Error) {

	if !k.HasValidatorDistInfo(ctx, operatorAddr) {
		return sdk.Coins{}, types.ErrNoValidatorDistInfo(k.codespace)
	}

	// withdraw self-delegation
	accAddr := sdk.AccAddress(operatorAddr.Bytes())
	withdraw := k.CurrentDelegationRewardsAll(ctx, accAddr)

	// withdrawal validator commission rewards
	valInfo := k.GetValidatorDistInfo(ctx, operatorAddr)

	wc := k.GetWithdrawContext(ctx, operatorAddr)
	commission := valInfo.CurrentCommissionRewards(wc)
	withdraw = withdraw.Plus(commission)
	truncated, _ := withdraw.TruncateDecimal()
	return truncated, nil
}

// remove all validator distribution infos
func (k Keeper) RemoveValidatorDistInfos(ctx sdk.Context) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, ValidatorDistInfoKey)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		store.Delete(iter.Key())
	}
}
