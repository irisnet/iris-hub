package keeper

import (
	"encoding/json"
	"github.com/irisnet/irishub/tests"
	"testing"

	"github.com/irisnet/irishub/app/v1/asset/internal/types"
	"github.com/irisnet/irishub/app/v1/auth"
	"github.com/irisnet/irishub/app/v1/bank"
	"github.com/irisnet/irishub/app/v1/params"
	"github.com/irisnet/irishub/codec"
	sdk "github.com/irisnet/irishub/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
)

func TestKeeper_IssueToken(t *testing.T) {
	ms, accountKey, assetKey, paramskey, paramsTkey := tests.SetupMultiStore()

	cdc := codec.New()
	types.RegisterCodec(cdc)
	auth.RegisterBaseAccount(cdc)

	ctx := sdk.NewContext(ms, abci.Header{}, false, log.NewNopLogger())
	pk := params.NewKeeper(cdc, paramskey, paramsTkey)
	ak := auth.NewAccountKeeper(cdc, accountKey, auth.ProtoBaseAccount)
	bk := bank.NewBaseKeeper(cdc, ak)
	keeper := NewKeeper(cdc, assetKey, bk, types.DefaultCodespace, pk.Subspace(types.DefaultParamSpace))
	addr := sdk.AccAddress([]byte("addr1"))

	acc := ak.NewAccountWithAddress(ctx, addr)

	ft := types.NewFungibleToken(types.NATIVE, "", "btc", "btc", 1, "", "satoshi", sdk.NewIntWithDecimal(1, 0), sdk.NewIntWithDecimal(1, 0), true, acc.GetAddress())
	_, err := keeper.IssueToken(ctx, ft)
	assert.NoError(t, err)

	assert.True(t, keeper.HasToken(ctx, "btc"))

	token, found := keeper.getToken(ctx, "btc")
	assert.True(t, found)

	assert.Equal(t, ft.GetDenom(), token.GetDenom())
	assert.Equal(t, ft.Owner, ft.Owner)

	msgJson, _ := json.Marshal(ft)
	assetJson, _ := json.Marshal(token)
	assert.Equal(t, msgJson, assetJson)
}

func TestKeeper_IssueGatewayToken(t *testing.T) {
	ms, accountKey, assetKey, paramskey, paramsTkey := tests.SetupMultiStore()

	cdc := codec.New()
	types.RegisterCodec(cdc)
	auth.RegisterBaseAccount(cdc)

	ctx := sdk.NewContext(ms, abci.Header{}, false, log.NewNopLogger())
	pk := params.NewKeeper(cdc, paramskey, paramsTkey)
	ak := auth.NewAccountKeeper(cdc, accountKey, auth.ProtoBaseAccount)
	bk := bank.NewBaseKeeper(cdc, ak)
	keeper := NewKeeper(cdc, assetKey, bk, types.DefaultCodespace, pk.Subspace(types.DefaultParamSpace))

	owner := ak.NewAccountWithAddress(ctx, []byte("owner"))
	gatewayOwner := ak.NewAccountWithAddress(ctx, []byte("gatewayOwner"))

	moniker := "moniker"
	identity := "identity"
	details := "details"
	website := "website"

	gateway := types.NewGateway(gatewayOwner.GetAddress(), moniker, identity, details, website)
	gatewayToken := types.NewFungibleToken(types.GATEWAY, "test", "btc", "btc", 1, "btc", "satoshi", sdk.NewIntWithDecimal(1, 0), sdk.NewIntWithDecimal(1, 0), true, owner.GetAddress())
	gatewayToken1 := types.NewFungibleToken(types.GATEWAY, "moniker", "btc", "btc", 1, "btc", "satoshi", sdk.NewIntWithDecimal(1, 0), sdk.NewIntWithDecimal(1, 0), true, gatewayOwner.GetAddress())

	// unknown gateway moniker
	_, err := keeper.IssueToken(ctx, gatewayToken)
	assert.Error(t, err)
	token, found := keeper.getToken(ctx, "test.btc")
	assert.False(t, found)

	// unauthorized creator
	keeper.SetGateway(ctx, gateway)
	_, err = keeper.IssueToken(ctx, gatewayToken)
	assert.Error(t, err)
	token, found = keeper.getToken(ctx, "moniker.btc")
	assert.False(t, found)

	_, err = keeper.IssueToken(ctx, gatewayToken1)
	assert.NoError(t, err)
	token, found = keeper.getToken(ctx, "moniker.btc")
	assert.True(t, found)
	assert.Equal(t, "moniker.btc", token.GetUniqueID())
}

func TestCreateGatewayKeeper(t *testing.T) {
	ms, accountKey, assetKey, paramskey, paramsTkey := tests.SetupMultiStore()

	cdc := codec.New()
	types.RegisterCodec(cdc)
	auth.RegisterBaseAccount(cdc)

	ctx := sdk.NewContext(ms, abci.Header{}, false, log.NewNopLogger())
	paramsKeeper := params.NewKeeper(cdc, paramskey, paramsTkey)
	ak := auth.NewAccountKeeper(cdc, accountKey, auth.ProtoBaseAccount)
	bk := bank.NewBaseKeeper(cdc, ak)
	keeper := NewKeeper(cdc, assetKey, bk, types.DefaultCodespace, paramsKeeper.Subspace(types.DefaultParamSpace))

	// define variables
	owner := ak.NewAccountWithAddress(ctx, []byte("owner"))
	moniker := "moniker"
	identity := "identity"
	details := "details"
	website := "website"

	// construct a test gateway
	gateway := types.NewGateway(owner.GetAddress(), moniker, identity, details, website)

	// assert the gateway of the given moniker does not exist at the beginning
	require.False(t, keeper.HasGateway(ctx, moniker))

	// create a gateway and assert that the gateway exists now
	keeper.SetGateway(ctx, gateway)
	require.True(t, keeper.HasGateway(ctx, moniker))

	// assert GetGateway will return the previous gateway
	res, _ := keeper.GetGateway(ctx, moniker)
	require.Equal(t, gateway, res)
}

func TestEditGatewayKeeper(t *testing.T) {
	ms, accountKey, assetKey, paramskey, paramsTkey := tests.SetupMultiStore()

	cdc := codec.New()
	types.RegisterCodec(cdc)
	auth.RegisterBaseAccount(cdc)

	ctx := sdk.NewContext(ms, abci.Header{}, false, log.NewNopLogger())
	paramsKeeper := params.NewKeeper(cdc, paramskey, paramsTkey)
	ak := auth.NewAccountKeeper(cdc, accountKey, auth.ProtoBaseAccount)
	bk := bank.NewBaseKeeper(cdc, ak)
	keeper := NewKeeper(cdc, assetKey, bk, types.DefaultCodespace, paramsKeeper.Subspace(types.DefaultParamSpace))

	// define variables
	owner := ak.NewAccountWithAddress(ctx, []byte("owner")).GetAddress()
	moniker := "moniker"
	identity := "identity"
	details := "details"
	website := "website"
	newIdentity := "new identity"
	newDetails := "new details"
	newWebsite := "new website"

	// build a MsgCreateGateway
	createMsg := types.NewMsgCreateGateway(owner, moniker, identity, details, website)

	// create a gateway and assert that the gateway exists now
	_, err := keeper.CreateGateway(ctx, createMsg)
	require.Nil(t, err)
	require.True(t, keeper.HasGateway(ctx, moniker))

	// assert GetGateway will return the previous gateway
	res, _ := keeper.GetGateway(ctx, moniker)
	require.Equal(t, identity, res.Identity)
	require.Equal(t, details, res.Details)
	require.Equal(t, website, res.Website)

	// build a MsgEditGateway
	editMsg := types.NewMsgEditGateway(owner, moniker, newIdentity, newDetails, newWebsite)

	// edit the gateway
	_, err = keeper.EditGateway(ctx, editMsg)
	require.Nil(t, err)

	// assert GetGateway will return the new filed values
	res, _ = keeper.GetGateway(ctx, moniker)
	require.Equal(t, newIdentity, res.Identity)
	require.Equal(t, newDetails, res.Details)
	require.Equal(t, newWebsite, res.Website)

	// build another MsgEditGateway with details and website not updated
	editMsg = types.NewMsgEditGateway(owner, moniker, identity, types.DoNotModify, types.DoNotModify)

	// edit the gateway again
	_, err = keeper.EditGateway(ctx, editMsg)
	require.Nil(t, err)

	// assert GetGateway will return the gateway with only identity updated
	res, _ = keeper.GetGateway(ctx, moniker)
	require.Equal(t, identity, res.Identity)
	require.Equal(t, newDetails, res.Details)
	require.Equal(t, newWebsite, res.Website)
}

func TestQueryGatewayKeeper(t *testing.T) {
	ms, accountKey, assetKey, paramskey, paramsTkey := tests.SetupMultiStore()

	cdc := codec.New()
	types.RegisterCodec(cdc)
	auth.RegisterBaseAccount(cdc)

	ctx := sdk.NewContext(ms, abci.Header{}, false, log.NewNopLogger())
	paramsKeeper := params.NewKeeper(cdc, paramskey, paramsTkey)
	ak := auth.NewAccountKeeper(cdc, accountKey, auth.ProtoBaseAccount)
	bk := bank.NewBaseKeeper(cdc, ak)
	keeper := NewKeeper(cdc, assetKey, bk, types.DefaultCodespace, paramsKeeper.Subspace(types.DefaultParamSpace))

	// define variables
	var (
		owners     = []sdk.AccAddress{ak.NewAccountWithAddress(ctx, []byte("owner1")).GetAddress(), ak.NewAccountWithAddress(ctx, []byte("owner2")).GetAddress()}
		monikers   = []string{"moni", "ker"}
		identities = []string{"id1", "id2"}
		details    = []string{"details1", "details2"}
		websites   = []string{"website1", "website2"}
	)

	// construct gateways
	gateway1 := types.NewGateway(owners[0], monikers[0], identities[0], details[0], websites[0])
	gateway2 := types.NewGateway(owners[1], monikers[1], identities[1], details[1], websites[1])

	// create gateways
	keeper.SetGateway(ctx, gateway1)
	keeper.SetOwnerGateway(ctx, gateway1.Owner, gateway1.Moniker)

	keeper.SetGateway(ctx, gateway2)
	keeper.SetOwnerGateway(ctx, gateway2.Owner, gateway2.Moniker)

	// query gateway
	res1, _ := keeper.GetGateway(ctx, gateway1.Moniker)
	require.Equal(t, gateway1, res1)

	res2, _ := keeper.GetGateway(ctx, gateway2.Moniker)
	require.Equal(t, gateway2, res2)

	// query gateways with a specified owner
	var gateways1 []types.Gateway
	iter1 := keeper.GetGateways(ctx, gateway1.Owner)
	defer iter1.Close()

	for ; iter1.Valid(); iter1.Next() {
		var moniker string
		keeper.cdc.MustUnmarshalBinaryLengthPrefixed(iter1.Value(), &moniker)

		gateway, err := keeper.GetGateway(ctx, moniker)
		if err != nil {
			continue
		}

		gateways1 = append(gateways1, gateway)
	}

	require.Equal(t, []types.Gateway{gateway1}, gateways1)

	var gateways2 []types.Gateway
	iter2 := keeper.GetGateways(ctx, gateway2.Owner)
	defer iter2.Close()

	for ; iter2.Valid(); iter2.Next() {
		var moniker string
		keeper.cdc.MustUnmarshalBinaryLengthPrefixed(iter2.Value(), &moniker)

		gateway, err := keeper.GetGateway(ctx, moniker)
		if err != nil {
			continue
		}

		gateways2 = append(gateways2, gateway)
	}

	require.Equal(t, []types.Gateway{gateway2}, gateways2)

	// query all gateways
	var gateways3 []types.Gateway
	keeper.IterateGateways(ctx, func(gw types.Gateway) (stop bool) {
		gateways3 = append(gateways3, gw)
		return false
	})

	require.Equal(t, []types.Gateway{gateway2, gateway1}, gateways3)
}

//TODO:finish the test
func TestKeeper_EditToken(t *testing.T) {
	ms, accountKey, assetKey, paramskey, paramsTkey := tests.SetupMultiStore()

	cdc := codec.New()
	types.RegisterCodec(cdc)
	auth.RegisterBaseAccount(cdc)

	ctx := sdk.NewContext(ms, abci.Header{}, false, log.NewNopLogger())
	pk := params.NewKeeper(cdc, paramskey, paramsTkey)
	ak := auth.NewAccountKeeper(cdc, accountKey, auth.ProtoBaseAccount)
	bk := bank.NewBaseKeeper(cdc, ak)
	keeper := NewKeeper(cdc, assetKey, bk, types.DefaultCodespace, pk.Subspace(types.DefaultParamSpace))
	addr := sdk.AccAddress([]byte("addr1"))

	acc := ak.NewAccountWithAddress(ctx, addr)

	ft := types.NewFungibleToken(types.NATIVE, "", "btc", "btc", 1, "", "satoshi", sdk.NewIntWithDecimal(1, 0), sdk.NewIntWithDecimal(21000000, 0), true, acc.GetAddress())

	_, err := keeper.IssueToken(ctx, ft)
	assert.NoError(t, err)

	assert.True(t, keeper.HasToken(ctx, "i.btc"))

	token, found := keeper.getToken(ctx, "i.btc")
	assert.True(t, found)

	assert.Equal(t, ft.GetDenom(), token.GetDenom())
	assert.Equal(t, ft.Owner, token.Owner)

	msgJson, _ := json.Marshal(ft)
	assetJson, _ := json.Marshal(token)
	assert.Equal(t, msgJson, assetJson)

	//TODO:finish the edit token
	mintable := false
	msgEditToken := types.NewMsgEditToken("BTC Token", "btc", "btc", "btc", 0, &mintable, acc.GetAddress())
	_, err = keeper.EditToken(ctx, msgEditToken)
	assert.NoError(t, err)
}

func TestTransferGatewayKeeper(t *testing.T) {
	ms, accountKey, assetKey, paramskey, paramsTkey := tests.SetupMultiStore()

	cdc := codec.New()
	types.RegisterCodec(cdc)
	auth.RegisterBaseAccount(cdc)

	ctx := sdk.NewContext(ms, abci.Header{}, false, log.NewNopLogger())
	paramsKeeper := params.NewKeeper(cdc, paramskey, paramsTkey)
	ak := auth.NewAccountKeeper(cdc, accountKey, auth.ProtoBaseAccount)
	bk := bank.NewBaseKeeper(cdc, ak)
	keeper := NewKeeper(cdc, assetKey, bk, types.DefaultCodespace, paramsKeeper.Subspace(types.DefaultParamSpace))

	// define variables
	originOwner := ak.NewAccountWithAddress(ctx, []byte("originOwner"))
	moniker := "moniker"
	identity := "identity"
	details := "details"
	website := "website"

	// construct a test gateway
	gateway := types.NewGateway(originOwner.GetAddress(), moniker, identity, details, website)

	// create a gateway
	keeper.SetGateway(ctx, gateway)

	// assert GetGateway will return the gateway with the previous owner
	res, _ := keeper.GetGateway(ctx, moniker)
	require.Equal(t, originOwner.GetAddress(), res.Owner)

	// build a msg for transferring the gateway owner
	newOwner := ak.NewAccountWithAddress(ctx, []byte("newOwner"))
	transferMsg := types.NewMsgTransferGatewayOwner(originOwner.GetAddress(), moniker, newOwner.GetAddress())

	// transfer
	_, err := keeper.TransferGatewayOwner(ctx, transferMsg)
	assert.NoError(t, err)

	// assert GetGateway will return the gateway with the new owner and the KeyOwnerGateway has been updated
	res, err = keeper.GetGateway(ctx, moniker)
	require.Equal(t, newOwner.GetAddress(), res.Owner)
	require.Equal(t, false, ctx.KVStore(keeper.storeKey).Has(KeyOwnerGateway(originOwner.GetAddress(), moniker)))
	require.Equal(t, true, ctx.KVStore(keeper.storeKey).Has(KeyOwnerGateway(newOwner.GetAddress(), moniker)))

	// transfer again and assert that the error will occur because of the ownership has been transferred
	_, err = keeper.TransferGatewayOwner(ctx, transferMsg)
	assert.Error(t, err)
}

func TestMintTokenKeeper(t *testing.T) {
	ms, accountKey, assetKey, paramskey, paramsTkey := tests.SetupMultiStore()

	cdc := codec.New()
	types.RegisterCodec(cdc)
	auth.RegisterBaseAccount(cdc)

	ctx := sdk.NewContext(ms, abci.Header{}, false, log.NewNopLogger())
	pk := params.NewKeeper(cdc, paramskey, paramsTkey)
	ak := auth.NewAccountKeeper(cdc, accountKey, auth.ProtoBaseAccount)
	bk := bank.NewBaseKeeper(cdc, ak)
	keeper := NewKeeper(cdc, assetKey, bk, types.DefaultCodespace, pk.Subspace(types.DefaultParamSpace))
	keeper.Init(ctx)

	addr := sdk.AccAddress([]byte("addr1"))

	acc := ak.NewAccountWithAddress(ctx, addr)
	amtCoin, _ := sdk.NewIntFromString("1000000000000000000000000000")
	coin := sdk.Coins{sdk.NewCoin("iris-atto", amtCoin)}
	bk.AddCoins(ctx, addr, coin)
	ak.IncreaseTotalLoosenToken(ctx, coin)

	ft := types.NewFungibleToken(types.NATIVE, "", "btc", "btc", 0, "", "satoshi", sdk.NewIntWithDecimal(1000, 0), sdk.NewIntWithDecimal(10000, 0), true, acc.GetAddress())
	_, err := keeper.IssueToken(ctx, ft)
	assert.NoError(t, err)

	assert.True(t, keeper.HasToken(ctx, "btc"))

	token, found := keeper.getToken(ctx, "btc")
	assert.True(t, found)

	assert.Equal(t, ft.GetDenom(), token.GetDenom())
	assert.Equal(t, ft.Owner, ft.Owner)

	msgJson, _ := json.Marshal(ft)
	assetJson, _ := json.Marshal(token)
	assert.Equal(t, msgJson, assetJson)

	msgMintToken := types.NewMsgMintToken("btc", acc.GetAddress(), nil, 1000)
	_, err = keeper.MintToken(ctx, msgMintToken)
	assert.NoError(t, err)

	balance := bk.GetCoins(ctx, acc.GetAddress())
	amt := balance.AmountOf("btc-min")
	assert.Equal(t, "2000", amt.String())
}

func TestTransferOwnerKeeper(t *testing.T) {
	ms, accountKey, assetKey, paramskey, paramsTkey := tests.SetupMultiStore()

	cdc := codec.New()
	types.RegisterCodec(cdc)
	auth.RegisterBaseAccount(cdc)

	ctx := sdk.NewContext(ms, abci.Header{}, false, log.NewNopLogger())
	pk := params.NewKeeper(cdc, paramskey, paramsTkey)
	ak := auth.NewAccountKeeper(cdc, accountKey, auth.ProtoBaseAccount)
	bk := bank.NewBaseKeeper(cdc, ak)
	keeper := NewKeeper(cdc, assetKey, bk, types.DefaultCodespace, pk.Subspace(types.DefaultParamSpace))

	srcOwner := sdk.AccAddress([]byte("TokenSrcOwner"))

	acc := ak.NewAccountWithAddress(ctx, srcOwner)

	ft := types.NewFungibleToken(types.NATIVE, "", "btc", "btc", 1, "", "satoshi", sdk.NewIntWithDecimal(1, 0), sdk.NewIntWithDecimal(21000000, 0), true, acc.GetAddress())

	_, err := keeper.IssueToken(ctx, ft)
	assert.NoError(t, err)

	assert.True(t, keeper.HasToken(ctx, "i.btc"))

	token, found := keeper.getToken(ctx, "i.btc")
	assert.True(t, found)

	assert.Equal(t, ft.GetDenom(), token.GetDenom())
	assert.Equal(t, ft.Owner, token.Owner)

	msgJson, _ := json.Marshal(ft)
	assetJson, _ := json.Marshal(token)
	assert.Equal(t, msgJson, assetJson)

	dstOwner := sdk.AccAddress([]byte("TokenDstOwner"))
	msg := types.MsgTransferTokenOwner{
		SrcOwner: srcOwner,
		DstOwner: dstOwner,
		TokenId:  "btc",
	}
	_, err = keeper.TransferTokenOwner(ctx, msg)
	assert.NoError(t, err)

	token, found = keeper.getToken(ctx, "i.btc")
	assert.True(t, found)
	assert.Equal(t, dstOwner, token.Owner)
}
