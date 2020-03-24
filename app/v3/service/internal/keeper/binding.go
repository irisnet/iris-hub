package keeper

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/irisnet/irishub/app/v1/auth"
	"github.com/irisnet/irishub/app/v3/service/internal/types"
	sdk "github.com/irisnet/irishub/types"
)

// AddServiceBinding creates a new service binding
func (k Keeper) AddServiceBinding(
	ctx sdk.Context,
	serviceName string,
	provider sdk.AccAddress,
	deposit sdk.Coins,
	pricing string,
) sdk.Error {
	if _, found := k.GetServiceDefinition(ctx, serviceName); !found {
		return types.ErrUnknownServiceDefinition(k.codespace, serviceName)
	}

	if _, found := k.GetServiceBinding(ctx, serviceName, provider); found {
		return types.ErrServiceBindingExists(k.codespace)
	}

	parsedPricing, err := k.ParsePricing(ctx, pricing)
	if err != nil {
		return err
	}

	if err := types.ValidatePricing(parsedPricing); err != nil {
		return err
	}

	minDeposit := k.getMinDeposit(ctx, parsedPricing)
	if !deposit.IsAllGTE(minDeposit) {
		return types.ErrInvalidDeposit(k.codespace, fmt.Sprintf("insufficient deposit: minimal deposit %s, %s got", minDeposit, deposit))
	}

	// Send coins from the provider's account to ServiceDepositCoinsAccAddr
	_, err = k.bk.SendCoins(ctx, provider, auth.ServiceDepositCoinsAccAddr, deposit)
	if err != nil {
		return err
	}

	available := true
	disabledTime := time.Time{}

	svcBinding := types.NewServiceBinding(serviceName, provider, deposit, pricing, available, disabledTime)
	k.SetServiceBinding(ctx, svcBinding)

	k.SetPricing(ctx, serviceName, provider, parsedPricing)

	return nil
}

// UpdateServiceBinding updates the specified service binding
func (k Keeper) UpdateServiceBinding(
	ctx sdk.Context,
	serviceName string,
	provider sdk.AccAddress,
	deposit sdk.Coins,
	pricing string,
) (err sdk.Error) {
	binding, found := k.GetServiceBinding(ctx, serviceName, provider)
	if !found {
		return types.ErrUnknownServiceBinding(k.codespace)
	}

	updated := false

	// add the deposit
	if !deposit.Empty() {
		binding.Deposit = binding.Deposit.Add(deposit)
		updated = true
	}

	parsedPricing := k.GetPricing(ctx, serviceName, provider)

	// update the pricing
	if len(pricing) != 0 {
		parsedPricing, err = k.ParsePricing(ctx, pricing)
		if err != nil {
			return err
		}

		if err := types.ValidatePricing(parsedPricing); err != nil {
			return err
		}

		binding.Pricing = pricing
		k.SetPricing(ctx, serviceName, provider, parsedPricing)

		updated = true
	}

	// only check deposit when the binding is available and updated
	if binding.Available && updated {
		minDeposit := k.getMinDeposit(ctx, parsedPricing)
		if !binding.Deposit.IsAllGTE(minDeposit) {
			return types.ErrInvalidDeposit(k.codespace, fmt.Sprintf("insufficient deposit: minimal deposit %s, %s got", minDeposit, binding.Deposit))
		}
	}

	if !deposit.Empty() {
		// Send coins from the provider's account to ServiceDepositCoinsAccAddr
		_, err := k.bk.SendCoins(ctx, provider, auth.ServiceDepositCoinsAccAddr, deposit)
		if err != nil {
			return err
		}
	}

	if updated {
		k.SetServiceBinding(ctx, binding)
	}

	return nil
}

// DisableService disables the specified service binding
func (k Keeper) DisableService(ctx sdk.Context, serviceName string, provider sdk.AccAddress) sdk.Error {
	binding, found := k.GetServiceBinding(ctx, serviceName, provider)
	if !found {
		return types.ErrUnknownServiceBinding(k.codespace)
	}

	if !binding.Available {
		return types.ErrServiceBindingUnavailable(k.codespace)
	}

	binding.Available = false
	binding.DisabledTime = ctx.BlockHeader().Time

	k.SetServiceBinding(ctx, binding)

	return nil
}

// EnableService enables the specified service binding
func (k Keeper) EnableService(ctx sdk.Context, serviceName string, provider sdk.AccAddress, deposit sdk.Coins) sdk.Error {
	binding, found := k.GetServiceBinding(ctx, serviceName, provider)
	if !found {
		return types.ErrUnknownServiceBinding(k.codespace)
	}

	if binding.Available {
		return types.ErrServiceBindingAvailable(k.codespace)
	}

	// add the deposit
	if !deposit.Empty() {
		binding.Deposit = binding.Deposit.Add(deposit)
	}

	minDeposit := k.getMinDeposit(ctx, k.GetPricing(ctx, serviceName, provider))
	if !binding.Deposit.IsAllGTE(minDeposit) {
		return types.ErrInvalidDeposit(k.codespace, fmt.Sprintf("insufficient deposit: minimal deposit %s, %s got", minDeposit, binding.Deposit))
	}

	if !deposit.Empty() {
		// Send coins from the provider's account to ServiceDepositCoinsAccAddr
		_, err := k.bk.SendCoins(ctx, provider, auth.ServiceDepositCoinsAccAddr, deposit)
		if err != nil {
			return err
		}
	}

	binding.Available = true
	binding.DisabledTime = time.Time{}

	k.SetServiceBinding(ctx, binding)

	return nil
}

// RefundDeposit refunds the deposit from the specified service binding
func (k Keeper) RefundDeposit(ctx sdk.Context, serviceName string, provider sdk.AccAddress) sdk.Error {
	binding, found := k.GetServiceBinding(ctx, serviceName, provider)
	if !found {
		return types.ErrUnknownServiceBinding(k.codespace)
	}

	if binding.Available {
		return types.ErrServiceBindingAvailable(k.codespace)
	}

	if binding.Deposit.IsZero() {
		return types.ErrInvalidDeposit(k.codespace, "the deposit of the service binding is zero")
	}

	params := k.GetParamSet(ctx)
	refundableTime := binding.DisabledTime.Add(params.ArbitrationTimeLimit).Add(params.ComplaintRetrospect)

	currentTime := ctx.BlockHeader().Time
	if currentTime.Before(refundableTime) {
		return types.ErrIncorrectRefundTime(k.codespace, fmt.Sprintf("%v", refundableTime))
	}

	// Send coins from ServiceDepositCoinsAccAddr to the provider's account
	_, err := k.bk.SendCoins(ctx, auth.ServiceDepositCoinsAccAddr, binding.Provider, binding.Deposit)
	if err != nil {
		return err
	}

	binding.Deposit = sdk.Coins{}
	k.SetServiceBinding(ctx, binding)

	return nil
}

// RefundDeposits refunds the deposits of all the binding services
func (k Keeper) RefundDeposits(ctx sdk.Context) sdk.Error {
	iterator := k.AllServiceBindingsIterator(ctx)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var binding types.ServiceBinding
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &binding)

		_, err := k.bk.SendCoins(ctx, auth.ServiceDepositCoinsAccAddr, binding.Provider, binding.Deposit)
		if err != nil {
			return err
		}
	}

	return nil
}

// SetServiceBinding sets the service binding
func (k Keeper) SetServiceBinding(ctx sdk.Context, svcBinding types.ServiceBinding) {
	store := ctx.KVStore(k.storeKey)

	bz := k.cdc.MustMarshalBinaryLengthPrefixed(svcBinding)
	store.Set(GetServiceBindingKey(svcBinding.ServiceName, svcBinding.Provider), bz)
}

// GetServiceBinding retrieves the specified service binding
func (k Keeper) GetServiceBinding(ctx sdk.Context, serviceName string, provider sdk.AccAddress) (svcBinding types.ServiceBinding, found bool) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(GetServiceBindingKey(serviceName, provider))
	if bz == nil {
		return svcBinding, false
	}

	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &svcBinding)
	return svcBinding, true
}

// ParsePricing parses the given string to Pricing
func (k Keeper) ParsePricing(ctx sdk.Context, pricing string) (p types.Pricing, err sdk.Error) {
	var rawPricing types.RawPricing
	if err := json.Unmarshal([]byte(pricing), &rawPricing); err != nil {
		return p, types.ErrInvalidPricing(k.codespace, fmt.Sprintf("failed to unmarshal the pricing: %s", err))
	}

	var coins sdk.Coins

	coinStrs := strings.Split(rawPricing.Price, ",")
	for _, coinStr := range coinStrs {
		unitName, amtStr, _ := sdk.ParseCoinParts(coinStr)
		if err != nil {
			return p, types.ErrInvalidPricing(k.codespace, fmt.Sprintf("failed to parse the pricing: %s", err))
		}

		amt, err := sdk.NewDecFromStr(amtStr)
		if err != nil {
			return p, types.ErrInvalidPricing(k.codespace, fmt.Sprintf("failed to parse the pricing: %s", err))
		}

		if unitName == sdk.Iris {
			coins = coins.Add(sdk.NewCoins(sdk.NewCoin(
				sdk.IrisAtto,
				amt.Mul(sdk.NewDecFromInt(sdk.NewIntWithDecimal(1, sdk.AttoScale))).TruncateInt(),
			)))
		} else {
			token, err := k.ak.GetToken(ctx, unitName)
			if err != nil {
				return p, types.ErrInvalidPricing(k.codespace, fmt.Sprintf("invalid price: %s", err))
			}

			coins = coins.Add(sdk.NewCoins(sdk.NewCoin(
				token.GetDenom(),
				amt.Mul(sdk.NewDecFromInt(sdk.NewIntWithDecimal(1, int(token.GetDecimal())))).TruncateInt(),
			)))
		}
	}

	p.Price = coins
	p.PromotionsByTime = rawPricing.PromotionsByTime
	p.PromotionsByVolume = rawPricing.PromotionsByVolume

	return p, nil
}

// SetPricing sets the pricing
func (k Keeper) SetPricing(
	ctx sdk.Context,
	serviceName string,
	provider sdk.AccAddress,
	pricing types.Pricing,
) {
	store := ctx.KVStore(k.storeKey)

	bz := k.cdc.MustMarshalBinaryLengthPrefixed(pricing)
	store.Set(GetPricingKey(serviceName, provider), bz)
}

// GetPricing retrieves the specified pricing
func (k Keeper) GetPricing(ctx sdk.Context, serviceName string, provider sdk.AccAddress) (pricing types.Pricing) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(GetPricingKey(serviceName, provider))
	if bz == nil {
		return
	}

	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &pricing)
	return pricing
}

// SetWithdrawAddress sets the withdrawal address for the specified provider
func (k Keeper) SetWithdrawAddress(ctx sdk.Context, provider, withdrawAddr sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Set(GetWithdrawAddrKey(provider), withdrawAddr.Bytes())
}

// GetWithdrawAddress gets the withdrawal address of the specified provider
func (k Keeper) GetWithdrawAddress(ctx sdk.Context, provider sdk.AccAddress) sdk.AccAddress {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(GetWithdrawAddrKey(provider))
	if bz == nil {
		return provider
	}

	return sdk.AccAddress(bz)
}

// IterateWithdrawAddresses iterates through all withdrawAddresses
func (k Keeper) IterateWithdrawAddresses(
	ctx sdk.Context,
	op func(provider sdk.AccAddress, withdrawAddress sdk.AccAddress) (stop bool),
) {
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, withdrawAddrKey)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		providerAddress := sdk.AccAddress(iterator.Key()[1:])
		withdrawAddress := sdk.AccAddress(iterator.Value())

		if stop := op(providerAddress, withdrawAddress); stop {
			break
		}
	}
}

// ServiceBindingsIterator returns an iterator for all bindings of the specified service
func (k Keeper) ServiceBindingsIterator(ctx sdk.Context, serviceName string) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, GetBindingsSubspace(serviceName))
}

// AllServiceBindingsIterator returns an iterator for all bindings
func (k Keeper) AllServiceBindingsIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, serviceBindingKey)
}

func (k Keeper) IterateServiceBindings(
	ctx sdk.Context,
	op func(binding types.ServiceBinding) (stop bool),
) {
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, serviceBindingKey)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var binding types.ServiceBinding
		k.GetCdc().MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &binding)

		if stop := op(binding); stop {
			break
		}
	}
}

// getMinDeposit gets the minimal deposit required for the service binding
func (k Keeper) getMinDeposit(ctx sdk.Context, pricing types.Pricing) sdk.Coins {
	params := k.GetParamSet(ctx)
	minDepositMultiple := sdk.NewInt(params.MinDepositMultiple)
	minDepositParam := params.MinDeposit

	price := pricing.Price.AmountOf(sdk.IrisAtto)

	// minimal deposit = max(price * minDepositMultiple, minDepositParam)
	minDeposit := sdk.NewCoins(sdk.NewCoin(sdk.IrisAtto, price.Mul(minDepositMultiple)))
	if minDeposit.IsAllLT(minDepositParam) {
		minDeposit = minDepositParam
	}

	return minDeposit
}
