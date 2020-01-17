package keeper

import (
	"testing"

	"github.com/irisnet/irishub/tests"

	"github.com/irisnet/irishub/app/v1/auth"
	"github.com/irisnet/irishub/app/v1/bank"
	"github.com/irisnet/irishub/app/v1/params"
	"github.com/irisnet/irishub/app/v3/asset/internal/types"
	"github.com/irisnet/irishub/codec"
	sdk "github.com/irisnet/irishub/types"
	"github.com/stretchr/testify/assert"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
)

func TestKeeperIssueToken(t *testing.T) {
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

	msg := types.NewMsgIssueToken("btc", "satoshi", "Bitcoin Network", 18, 21000000, 21000000, false, acc.GetAddress())
	_, err := keeper.IssueToken(ctx, msg)
	assert.NoError(t, err)

	assert.True(t, keeper.HasToken(ctx, types.GetTokenID(msg.Symbol)))

	token, err := keeper.getToken(ctx, types.GetTokenID(msg.Symbol))
	ft := types.NewFungibleToken(msg.Symbol, msg.Name, msg.MinUnitAlias, msg.Decimal, sdk.NewIntWithDecimal(int64(msg.InitialSupply), int(msg.Decimal)), sdk.NewIntWithDecimal(int64(msg.MaxSupply), int(msg.Decimal)), msg.Mintable, msg.Owner)
	assert.NoError(t, err)
	assert.Equal(t, ft, token)
}

func TestKeeperEditToken(t *testing.T) {
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

	msg := types.NewMsgIssueToken("btc", "satoshi", "Bitcoin Network", 18, 21000000, 21000000, false, acc.GetAddress())

	_, err := keeper.IssueToken(ctx, msg)
	assert.NoError(t, err)

	assert.True(t, keeper.HasToken(ctx, types.GetTokenID(msg.Symbol)))

	token, err := keeper.getToken(ctx, types.GetTokenID(msg.Symbol))
	ft := types.NewFungibleToken(msg.Symbol, msg.Name, msg.MinUnitAlias, msg.Decimal, sdk.NewIntWithDecimal(int64(msg.InitialSupply), int(msg.Decimal)), sdk.NewIntWithDecimal(int64(msg.MaxSupply), int(msg.Decimal)), msg.Mintable, msg.Owner)
	assert.NoError(t, err)
	assert.Equal(t, ft, token)

	//TODO:finish the edit token
	mintable := types.False
	msgEditToken := types.NewMsgEditToken("BTC Token", "btc", 0, mintable, acc.GetAddress())
	_, err = keeper.EditToken(ctx, msgEditToken)
	assert.NoError(t, err)

	token, err = keeper.getToken(ctx, types.GetTokenID(ft.Symbol))
	assert.NoError(t, err)
	ft.Name = msgEditToken.Name
	ft.Mintable = false

	assert.Equal(t, token, ft)
}

func TestMintToken(t *testing.T) {
	ms, accountKey, assetKey, paramskey, paramsTkey := tests.SetupMultiStore()

	cdc := codec.New()
	types.RegisterCodec(cdc)
	auth.RegisterBaseAccount(cdc)

	ctx := sdk.NewContext(ms, abci.Header{}, false, log.NewNopLogger())
	pk := params.NewKeeper(cdc, paramskey, paramsTkey)
	ak := auth.NewAccountKeeper(cdc, accountKey, auth.ProtoBaseAccount)
	bk := bank.NewBaseKeeper(cdc, ak)
	keeper := NewKeeper(cdc, assetKey, bk, types.DefaultCodespace, pk.Subspace(types.DefaultParamSpace))
	keeper.SetParamSet(ctx, types.DefaultParams())

	addr := sdk.AccAddress([]byte("addr1"))

	acc := ak.NewAccountWithAddress(ctx, addr)
	amtCoin, _ := sdk.NewIntFromString("1000000000000000000000000000")
	coin := sdk.Coins{sdk.NewCoin("iris-atto", amtCoin)}
	_, _, err := bk.AddCoins(ctx, addr, coin)
	assert.NoError(t, err)
	ak.IncreaseTotalLoosenToken(ctx, coin)

	msg := types.NewMsgIssueToken("btc", "satoshi", "Bitcoin Network", 0, 1000, 2100, true, acc.GetAddress())
	_, err = keeper.IssueToken(ctx, msg)
	assert.NoError(t, err)

	assert.True(t, keeper.HasToken(ctx, types.GetTokenID(msg.Symbol)))

	token, err := keeper.getToken(ctx, types.GetTokenID(msg.Symbol))
	ft := types.NewFungibleToken(msg.Symbol, msg.Name, msg.MinUnitAlias, msg.Decimal, sdk.NewIntWithDecimal(int64(msg.InitialSupply), int(msg.Decimal)), sdk.NewIntWithDecimal(int64(msg.MaxSupply), int(msg.Decimal)), msg.Mintable, msg.Owner)
	assert.NoError(t, err)
	assert.Equal(t, ft, token)

	msgMintToken := types.NewMsgMintToken(types.GetTokenID(ft.Symbol), acc.GetAddress(), nil, 1000)
	_, err = keeper.MintToken(ctx, msgMintToken)
	assert.NoError(t, err)

	balance := bk.GetCoins(ctx, acc.GetAddress())
	amt := balance.AmountOf(ft.GetDenom())
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

	issueMsg := types.NewMsgIssueToken("btc", "satoshi", "Bitcoin Network", 18, 21000000, 21000000, false, acc.GetAddress())
	_, err := keeper.IssueToken(ctx, issueMsg)
	assert.NoError(t, err)

	assert.True(t, keeper.HasToken(ctx, types.GetTokenID(issueMsg.Symbol)))

	tokenSrc, err := keeper.getToken(ctx, types.GetTokenID(issueMsg.Symbol))
	ft := types.NewFungibleToken(issueMsg.Symbol, issueMsg.Name, issueMsg.MinUnitAlias, issueMsg.Decimal, sdk.NewIntWithDecimal(int64(issueMsg.InitialSupply), int(issueMsg.Decimal)), sdk.NewIntWithDecimal(int64(issueMsg.MaxSupply), int(issueMsg.Decimal)), issueMsg.Mintable, issueMsg.Owner)
	assert.NoError(t, err)

	assert.Equal(t, ft, tokenSrc)

	dstOwner := sdk.AccAddress([]byte("TokenDstOwner"))
	msg := types.MsgTransferTokenOwner{
		SrcOwner: srcOwner,
		DstOwner: dstOwner,
		TokenId:  "btc",
	}
	_, err = keeper.TransferTokenOwner(ctx, msg)
	assert.NoError(t, err)

	token1, err := keeper.getToken(ctx, types.GetTokenID(ft.Symbol))
	assert.NoError(t, err)

	tokenSrc.Owner = dstOwner
	assert.Equal(t, tokenSrc, token1)

	keeper.iterateTokensWithOwner(ctx, dstOwner, func(token types.FungibleToken) (stop bool) {
		assert.Equal(t, token, tokenSrc)
		return false
	})
}
