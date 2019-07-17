package asset

import (
	"fmt"

	"github.com/irisnet/irishub/app/v1/asset/tags"
	"github.com/irisnet/irishub/app/v1/bank"
	"github.com/irisnet/irishub/app/v1/params"
	"github.com/irisnet/irishub/codec"
	"github.com/irisnet/irishub/modules/guardian"
	sdk "github.com/irisnet/irishub/types"
)

type Keeper struct {
	storeKey sdk.StoreKey
	cdc      *codec.Codec
	bk       bank.Keeper
	gk       guardian.Keeper

	// codespace
	codespace sdk.CodespaceType
	// params subspace
	paramSpace params.Subspace
}

func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, bk bank.Keeper, gk guardian.Keeper, codespace sdk.CodespaceType, paramSpace params.Subspace) Keeper {
	return Keeper{
		storeKey:   key,
		cdc:        cdc,
		bk:         bk,
		gk:         gk,
		codespace:  codespace,
		paramSpace: paramSpace.WithTypeTable(ParamTypeTable()),
	}
}

// return the codespace
func (k Keeper) Codespace() sdk.CodespaceType {
	return k.codespace
}

// IssueToken issue a new token
func (k Keeper) IssueToken(ctx sdk.Context, token FungibleToken) (sdk.Tags, sdk.Error) {
	if token.GetSource() == GATEWAY {
		gateway, err := k.GetGateway(ctx, token.GetGateway())
		if err != nil {
			return nil, err
		}
		if !gateway.Owner.Equals(token.GetOwner()) {
			return nil, ErrUnauthorizedIssueGatewayAsset(k.codespace,
				fmt.Sprintf("Gateway %s token can only be created by %s, unauthorized creator %s",
					gateway.Moniker, gateway.Owner, token.GetOwner()))
		}
	}
	token, owner, err := k.AddToken(ctx, token)
	if err != nil {
		return nil, err
	}

	// for native and gateway tokens
	if owner != nil {
		initialSupply := sdk.NewCoin(token.GetDenom(), token.GetInitSupply())

		// Add coins into owner's account
		_, _, err := k.bk.AddCoins(ctx, owner, sdk.Coins{initialSupply})
		if err != nil {
			return nil, err
		}

		// Set total supply
		k.bk.SetTotalSupply(ctx, initialSupply)
		ctx.CoinFlowTags().AppendCoinFlowTag(ctx, owner.String(), owner.String(), initialSupply.String(), sdk.IssueTokenFlow, "")
	}

	createTags := sdk.NewTags(
		tags.Id, []byte(token.GetUniqueID()),
		tags.Denom, []byte(token.GetDenom()),
		tags.Source, []byte(token.GetSource().String()),
		tags.Gateway, []byte(token.GetGateway()),
		tags.Owner, []byte(token.GetOwner().String()),
	)

	return createTags, nil
}

// save a new token to keystore
func (k Keeper) AddToken(ctx sdk.Context, token FungibleToken) (FungibleToken, sdk.AccAddress, sdk.Error) {
	token = token.Sanitize()
	tokenId, err := GetTokenID(token.GetSource(), token.GetSymbol(), token.GetGateway())
	if err != nil {
		return token, nil, err
	}
	if k.HasToken(ctx, tokenId) {
		return token, nil, ErrAssetAlreadyExists(k.codespace, fmt.Sprintf("token already exists: %s", token.GetUniqueID()))
	}

	var owner sdk.AccAddress
	if token.GetSource() == GATEWAY {
		gateway, err := k.GetGateway(ctx, token.GetGateway())
		if err != nil {
			return token, nil, err
		}
		owner = gateway.Owner
	} else if token.GetSource() == NATIVE {
		owner = token.GetOwner()
		token.SymbolAtSource = ""
		token.Gateway = ""
	}

	err = k.SetToken(ctx, token)
	if err != nil {
		return token, nil, err
	}

	// Set token to be prefixed with owner and source
	if token.GetSource() == NATIVE {
		err = k.SetTokens(ctx, owner, token)
		if err != nil {
			return token, nil, err
		}
	}

	// Set token to be prefixed with source
	err = k.SetTokens(ctx, sdk.AccAddress{}, token)
	if err != nil {
		return token, nil, err
	}

	return token, owner, nil
}

func (k Keeper) HasToken(ctx sdk.Context, tokenId string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(KeyToken(tokenId))
}

func (k Keeper) SetToken(ctx sdk.Context, token FungibleToken) sdk.Error {
	if token.GetSource() == GATEWAY {
		token.Owner = nil
	}

	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(token)

	tokenId, err := GetTokenID(token.GetSource(), token.GetSymbol(), token.GetGateway())
	if err != nil {
		return err
	}

	store.Set(KeyToken(tokenId), bz)
	return nil
}

func (k Keeper) SetTokens(ctx sdk.Context, owner sdk.AccAddress, token FungibleToken) sdk.Error {
	store := ctx.KVStore(k.storeKey)

	tokenId, err := GetTokenID(token.GetSource(), token.GetSymbol(), token.GetGateway())
	if err != nil {
		return err
	}

	bz := k.cdc.MustMarshalBinaryLengthPrefixed(tokenId)

	store.Set(KeyTokens(owner, tokenId), bz)
	return nil
}

func (k Keeper) getToken(ctx sdk.Context, tokenId string) (token FungibleToken, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(KeyToken(tokenId))
	if bz == nil {
		return token, false
	}

	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &token)
	return token, true
}

func (k Keeper) getTokens(ctx sdk.Context, owner sdk.AccAddress, nonSymbolTokenId string) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, KeyTokens(owner, nonSymbolTokenId))
}

// CreateGateway creates a gateway
func (k Keeper) CreateGateway(ctx sdk.Context, msg MsgCreateGateway) (sdk.Tags, sdk.Error) {
	// check if the moniker already exists
	if k.HasGateway(ctx, msg.Moniker) {
		return nil, ErrGatewayAlreadyExists(k.codespace, fmt.Sprintf("the moniker already exists:%s", msg.Moniker))
	}

	var gateway = Gateway{
		Owner:    msg.Owner,
		Moniker:  msg.Moniker,
		Identity: msg.Identity,
		Details:  msg.Details,
		Website:  msg.Website,
	}

	// set the gateway and related keys
	k.SetGateway(ctx, gateway)
	k.SetOwnerGateway(ctx, msg.Owner, msg.Moniker)

	// TODO
	createTags := sdk.NewTags(
		"moniker", []byte(msg.Moniker),
	)

	return createTags, nil
}

// EditGateway edits the specified gateway
func (k Keeper) EditGateway(ctx sdk.Context, msg MsgEditGateway) (sdk.Tags, sdk.Error) {
	// get the destination gateway
	gateway, err := k.GetGateway(ctx, msg.Moniker)
	if err != nil {
		return nil, err
	}

	// check if the given owner matches with the owner of the destination gateway
	if !msg.Owner.Equals(gateway.Owner) {
		return nil, ErrInvalidOwner(k.codespace, fmt.Sprintf("the address %d is not the owner of the gateway %s", msg.Owner, msg.Moniker))
	}

	// update the gateway
	if msg.Identity != DoNotModify {
		gateway.Identity = msg.Identity
	}
	if msg.Details != DoNotModify {
		gateway.Details = msg.Details
	}
	if msg.Website != DoNotModify {
		gateway.Website = msg.Website
	}

	// set the new gateway
	k.SetGateway(ctx, gateway)

	// TODO
	editTags := sdk.NewTags(
		"moniker", []byte(msg.Moniker),
	)

	return editTags, nil
}

// EditToken edits the specified token
func (k Keeper) EditToken(ctx sdk.Context, msg MsgEditToken) (sdk.Tags, sdk.Error) {
	// get the destination token
	token, exist := k.getToken(ctx, msg.TokenId)
	if !exist {
		return nil, ErrAssetNotExists(k.codespace, fmt.Sprintf("token %s does not exist", msg.TokenId))
	}

	if token.Source == GATEWAY {
		gateway, _ := k.GetGateway(ctx, token.Gateway)
		token.Owner = gateway.Owner
	}

	if !msg.Owner.Equals(token.Owner) {
		return nil, ErrInvalidOwner(k.codespace, fmt.Sprintf("the address %d is not the owner of the token %s", msg.Owner, msg.TokenId))
	}

	hasIssuedAmt, found := k.bk.GetTotalSupply(ctx, token.GetDenom())
	if !found {
		return nil, ErrAssetNotExists(k.codespace, fmt.Sprintf("token denom %s does not exist", token.GetDenom()))
	}

	maxSupply := sdk.NewIntWithDecimal(int64(msg.MaxSupply), int(token.Decimal))
	if maxSupply.GT(sdk.ZeroInt()) && (maxSupply.LT(hasIssuedAmt.Amount) || maxSupply.GT(token.MaxSupply)) {
		return nil, ErrInvalidAssetMaxSupply(k.codespace, fmt.Sprintf("max supply must not be less than %s and greater than %s", hasIssuedAmt.Amount.String(), token.MaxSupply.String()))
	}

	if msg.Name != DoNotModify {
		token.Name = msg.Name
	}
	if msg.SymbolAtSource != DoNotModify {
		token.SymbolAtSource = msg.SymbolAtSource
	}
	if msg.SymbolMinAlias != DoNotModify {
		token.SymbolMinAlias = msg.SymbolMinAlias
	}
	if maxSupply.GT(sdk.ZeroInt()) {
		token.MaxSupply = maxSupply
	}
	if msg.Mintable != nil {
		token.Mintable = *msg.Mintable
	}

	if err := k.SetToken(ctx, token); err != nil {
		return nil, err
	}

	editTags := sdk.NewTags(
		tags.Id, []byte(msg.TokenId),
	)

	return editTags, nil
}

// TransferGatewayOwner transfers the owner of the specified gateway to a new one
func (k Keeper) TransferGatewayOwner(ctx sdk.Context, msg MsgTransferGatewayOwner) (sdk.Tags, sdk.Error) {
	// get the destination gateway
	gateway, err := k.GetGateway(ctx, msg.Moniker)
	if err != nil {
		return nil, err
	}

	// check if the given owner matches with the owner of the destination gateway
	if !msg.Owner.Equals(gateway.Owner) {
		return nil, ErrInvalidOwner(k.codespace, fmt.Sprintf("the address %d is not the owner of the gateway %s", msg.Owner, msg.Moniker))
	}

	// change the ownership
	gateway.Owner = msg.To

	// update the gateway and related keys
	k.SetGateway(ctx, gateway)
	k.UpdateOwnerGateway(ctx, gateway.Moniker, msg.Owner, msg.To)

	// TODO
	transferTags := sdk.NewTags(
		"moniker", []byte(msg.Moniker),
	)

	return transferTags, nil
}

// GetGateway retrieves the gateway of the given moniker
func (k Keeper) GetGateway(ctx sdk.Context, moniker string) (Gateway, sdk.Error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(KeyGateway(moniker))
	if bz == nil {
		return Gateway{}, ErrUnkwownGateway(k.codespace, fmt.Sprintf("unknown gateway moniker:%s", moniker))
	}

	var gateway Gateway
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &gateway)

	return gateway, nil
}

// HasGateway checks if the given gateway exists. Return true if exists, false otherwise
func (k Keeper) HasGateway(ctx sdk.Context, moniker string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(KeyGateway(moniker))
}

// SetGateway stores the given gateway into the underlying storage
func (k Keeper) SetGateway(ctx sdk.Context, gateway Gateway) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(gateway)

	// set KeyGateway
	store.Set(KeyGateway(gateway.Moniker), bz)
}

// SetOwnerGateway stores the gateway moniker into storage by the key KeyOwnerGateway. Intended for iteration on gateways of an owner
func (k Keeper) SetOwnerGateway(ctx sdk.Context, owner sdk.AccAddress, moniker string) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(moniker)

	// set KeyOwnerGateway
	store.Set(KeyOwnerGateway(owner, moniker), bz)
}

// UpdateOwnerGateway updates the KeyOwnerGateway key of the given moniker from an owner to another
func (k Keeper) UpdateOwnerGateway(ctx sdk.Context, moniker string, originOwner, newOwner sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)

	// delete the old key
	store.Delete(KeyOwnerGateway(originOwner, moniker))

	// add the new key
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(moniker)
	store.Set(KeyOwnerGateway(newOwner, moniker), bz)
}

// GetGateways retrieves all the gateways of the given owner
func (k Keeper) GetGateways(ctx sdk.Context, owner sdk.AccAddress) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, KeyGatewaysSubspace(owner))
}

// IterateGateways iterates through all existing gateways
func (k Keeper) IterateGateways(ctx sdk.Context, op func(gateway Gateway) (stop bool)) {
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, PrefixGateway)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var gateway Gateway
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &gateway)

		if stop := op(gateway); stop {
			break
		}
	}
}

// IterateTokens iterates through all existing tokens
func (k Keeper) IterateTokens(ctx sdk.Context, op func(token FungibleToken) (stop bool)) {
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, PrefixToken)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var token FungibleToken
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &token)

		if stop := op(token); stop {
			break
		}
	}
}

func (k Keeper) Init(ctx sdk.Context) {
	k.SetParamSet(ctx, DefaultParams())
}

// TransferTokenOwner transfers the owner of the specified token to a new one
func (k Keeper) TransferTokenOwner(ctx sdk.Context, msg MsgTransferTokenOwner) (sdk.Tags, sdk.Error) {
	// get the destination token
	token, exist := k.getToken(ctx, msg.TokenId)
	if !exist {
		return nil, ErrAssetNotExists(k.codespace, fmt.Sprintf("token %s does not exist", msg.TokenId))
	}

	if token.Source != NATIVE {
		return nil, ErrInvalidAssetSource(k.codespace, fmt.Sprintf("only the token of which the source is native can be transferred,but the source of the current token is %s", token.Source.String()))
	}

	if !msg.SrcOwner.Equals(token.Owner) {
		return nil, ErrInvalidOwner(k.codespace, fmt.Sprintf("the address %s is not the owner of the token %s", msg.SrcOwner.String(), msg.TokenId))
	}

	token.Owner = msg.DstOwner

	// update token information
	if err := k.SetToken(ctx, token); err != nil {
		return nil, err
	}

	// reset all index for query-token
	if err := k.resetStoreKeyForQueryToken(ctx, msg, token); err != nil {
		return nil, err
	}
	tags := sdk.NewTags(
		tags.Id, []byte(msg.TokenId),
	)

	return tags, nil
}

// reset all index by DstOwner of token for query-token command
func (k Keeper) resetStoreKeyForQueryToken(ctx sdk.Context, msg MsgTransferTokenOwner, token FungibleToken) sdk.Error {
	store := ctx.KVStore(k.storeKey)

	tokenId, err := GetTokenID(token.GetSource(), token.GetSymbol(), token.GetGateway())
	if err != nil {
		return err
	}
	// delete the old key
	store.Delete(KeyTokens(msg.SrcOwner, tokenId))

	// add the new key
	return k.SetTokens(ctx, msg.DstOwner, token)
}

func (k Keeper) MintToken(ctx sdk.Context, msg MsgMintToken) (sdk.Tags, sdk.Error) {
	token, exist := k.getToken(ctx, msg.TokenId)
	if !exist {
		return nil, ErrAssetNotExists(k.codespace, fmt.Sprintf("token %s does not exist", msg.TokenId))
	}

	if token.Source == GATEWAY {
		gateway, _ := k.GetGateway(ctx, token.Gateway)
		token.Owner = gateway.Owner
	}

	if !msg.Owner.Equals(token.Owner) {
		return nil, ErrInvalidOwner(k.codespace, fmt.Sprintf("the address %s is not the owner of the token %s", msg.Owner.String(), msg.TokenId))
	}

	if !token.Mintable {
		return nil, ErrAssetNotMintable(k.codespace, fmt.Sprintf("the token %s is set to be non-mintable", msg.TokenId))
	}

	hasIssuedAmt, found := k.bk.GetTotalSupply(ctx, token.GetDenom())
	if !found {
		return nil, ErrAssetNotExists(k.codespace, fmt.Sprintf("token denom %s does not exist", token.GetDenom()))
	}

	//check the denom
	expDenom := token.GetDenom()
	if expDenom != hasIssuedAmt.Denom {
		return nil, ErrAssetNotExists(k.codespace, fmt.Sprintf("denom of mint token is not equal issued token,expected:%s,actual:%s", expDenom, hasIssuedAmt.Denom))
	}

	mintAmt := sdk.NewIntWithDecimal(int64(msg.Amount), int(token.Decimal))
	if mintAmt.Add(hasIssuedAmt.Amount).GT(token.MaxSupply) {
		exp := sdk.NewIntWithDecimal(1, int(token.Decimal))
		canAmt := token.MaxSupply.Sub(hasIssuedAmt.Amount).Div(exp)
		return nil, ErrInvalidAssetMaxSupply(k.codespace, fmt.Sprintf("The amount of mint tokens plus the total amount of issues has exceeded the maximum issue total,only accepts amount (0, %s]", canAmt.String()))
	}

	switch token.Source {
	case NATIVE:
		// handle fee for native token
		if err := TokenMintFeeHandler(ctx, k, msg.Owner, token.Symbol); err != nil {
			return nil, err
		}
		break
	case GATEWAY:
		// handle fee for gateway token
		if err := GatewayTokenMintFeeHandler(ctx, k, msg.Owner, token.Symbol); err != nil {
			return nil, err
		}
		break
	default:
		break
	}

	mintCoin := sdk.NewCoin(expDenom, mintAmt)
	//add TotalSupply
	if err := k.bk.IncreaseTotalSupply(ctx, mintCoin); err != nil {
		return nil, err
	}

	mintAcc := msg.To
	if mintAcc.Empty() {
		mintAcc = token.Owner
	}

	//add mintCoin to special account
	_, tags, err := k.bk.AddCoins(ctx, mintAcc, sdk.Coins{mintCoin})
	if err != nil {
		return nil, err
	}
	ctx.CoinFlowTags().AppendCoinFlowTag(ctx, msg.Owner.String(), mintAcc.String(), mintCoin.String(), sdk.MintTokenFlow, "")
	return tags, nil
}
