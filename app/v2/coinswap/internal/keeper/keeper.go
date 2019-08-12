package keeper

import (
	"fmt"

	"github.com/irisnet/irishub/app/v1/params"
	"github.com/irisnet/irishub/app/v2/coinswap/internal/types"
	"github.com/irisnet/irishub/codec"
	sdk "github.com/irisnet/irishub/types"
	"github.com/tendermint/tendermint/crypto"
)

// Keeper of the coinswap store
type Keeper struct {
	cdc        *codec.Codec
	storeKey   sdk.StoreKey
	bk         types.BankKeeper
	ak         types.AuthKeeper
	paramSpace params.Subspace
}

// NewKeeper returns a coinswap keeper. It handles:
// - creating new ModuleAccounts for each trading pair
// - burning minting liquidity coins
// - sending to and from ModuleAccounts
func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, bk types.BankKeeper, ak types.AuthKeeper, paramSpace params.Subspace) Keeper {
	return Keeper{
		storeKey:   key,
		bk:         bk,
		ak:         ak,
		cdc:        cdc,
		paramSpace: paramSpace.WithTypeTable(types.ParamTypeTable()),
	}
}

func (k Keeper) HandleSwap(ctx sdk.Context, msg types.MsgSwapOrder) (sdk.Tags, sdk.Error) {
	tags := sdk.EmptyTags()
	var amount sdk.Int
	var err sdk.Error
	var isDoubleSwap = msg.Input.Coin.Denom != sdk.IrisAtto && msg.Output.Coin.Denom != sdk.IrisAtto

	if isDoubleSwap && msg.IsBuyOrder {
		amount, err = k.doubleTradeInputForExactOutput(ctx, msg.Input, msg.Output)
	} else if isDoubleSwap && !msg.IsBuyOrder {
		amount, err = k.doubleTradeExactInputForOutput(ctx, msg.Input, msg.Output)
	} else if !isDoubleSwap && msg.IsBuyOrder {
		amount, err = k.tradeInputForExactOutput(ctx, msg.Input, msg.Output)
	} else if !isDoubleSwap && !msg.IsBuyOrder {
		amount, err = k.tradeExactInputForOutput(ctx, msg.Input, msg.Output)
	}
	if err != nil {
		return nil, err
	}
	tags.AppendTag(types.TagAmount, []byte(amount.String()))
	return tags, nil
}

func (k Keeper) HandleAddLiquidity(ctx sdk.Context, msg types.MsgAddLiquidity) sdk.Error {
	reservePoolName, err := k.GetReservePoolName(sdk.IrisAtto, msg.MaxToken.Denom)
	if err != nil {
		return err
	}
	reservePool := k.GetReservePool(ctx, reservePoolName)
	irisReserveAmt := reservePool.AmountOf(sdk.IrisAtto)
	tokenReserveAmt := reservePool.AmountOf(msg.MaxToken.Denom)
	liquidity := reservePool.AmountOf(reservePoolName)

	var mintLiquidityAmt sdk.Int
	var depositToken sdk.Coin
	var irisCoin = sdk.NewCoin(sdk.IrisAtto, msg.ExactIrisAmt)

	// calculate amount of UNI to be minted for sender
	// and coin amount to be deposited
	if liquidity.IsZero() {
		mintLiquidityAmt = msg.ExactIrisAmt
		depositToken = sdk.NewCoin(msg.MaxToken.Denom, msg.MaxToken.Amount)
	} else {
		mintLiquidityAmt = (liquidity.Mul(msg.ExactIrisAmt)).Div(irisReserveAmt)
		if mintLiquidityAmt.LT(msg.MinLiquidity) {
			return types.ErrConstraintNotMet(fmt.Sprintf("liquidity[%s] is less than user 's min reward[%s]", mintLiquidityAmt.String(), msg.MinLiquidity.String()))
		}
		depositAmt := (tokenReserveAmt.Mul(msg.ExactIrisAmt)).Div(irisReserveAmt)
		depositToken = sdk.NewCoin(msg.MaxToken.Denom, depositAmt)

		if depositAmt.GT(msg.MaxToken.Amount) {
			return types.ErrConstraintNotMet(fmt.Sprintf("amount[%s] of depositToken depositd is greater than user 's max deposited amount[%s]", depositToken.String(), msg.MaxToken.String()))
		}
	}
	return k.addLiquidity(ctx, msg.Sender, irisCoin, depositToken, reservePoolName, mintLiquidityAmt)
}

func (k Keeper) addLiquidity(ctx sdk.Context, sender sdk.AccAddress, irisCoin, token sdk.Coin, reservePoolName string, mintLiquidityAmt sdk.Int) sdk.Error {
	depositedTokens := sdk.NewCoins(irisCoin, token)
	poolAddr := getReservePoolAddr(reservePoolName)
	// transfer deposited token into coinswaps Account
	_, err := k.bk.SendCoins(ctx, sender, poolAddr, depositedTokens)
	if err != nil {
		return err
	}
	// mint liquidity vouchers for reserve Pool
	mintToken := sdk.NewCoins(sdk.NewCoin(reservePoolName, mintLiquidityAmt))
	k.bk.AddCoins(ctx, poolAddr, mintToken)
	// mint liquidity vouchers for sender
	k.bk.AddCoins(ctx, sender, mintToken)
	return nil
}

func (k Keeper) HandleRemoveLiquidity(ctx sdk.Context, msg types.MsgRemoveLiquidity) sdk.Error {
	reservePoolName, err := k.GetReservePoolName(sdk.IrisAtto, msg.MinToken.Denom)
	if err != nil {
		return err
	}

	// check if reserve pool exists
	reservePool := k.GetReservePool(ctx, reservePoolName)
	if reservePool == nil {
		return types.ErrReservePoolNotExists("")
	}

	irisReserveAmt := reservePool.AmountOf(sdk.IrisAtto)
	tokenReserveAmt := reservePool.AmountOf(msg.MinToken.Denom)
	liquidityReserve := reservePool.AmountOf(reservePoolName)
	if irisReserveAmt.LT(msg.MinIrisAmt) {
		return types.ErrInsufficientFunds(fmt.Sprintf("insufficient funds,actual:%s,expect:%s", irisReserveAmt.String(), msg.MinIrisAmt.String()))
	}
	if tokenReserveAmt.LT(msg.MinToken.Amount) {
		return types.ErrInsufficientFunds(fmt.Sprintf("insufficient funds,actual:%s,expect:%s", tokenReserveAmt.String(), msg.MinToken.Amount.String()))
	}
	if liquidityReserve.LT(msg.WithdrawLiquidity) {
		return types.ErrInsufficientFunds(fmt.Sprintf("insufficient funds,actual:%s,expect:%s", liquidityReserve.String(), msg.WithdrawLiquidity.String()))
	}

	// calculate amount of UNI to be burned for sender
	// and coin amount to be returned
	irisWithdrawnAmt := msg.WithdrawLiquidity.Mul(irisReserveAmt).Div(liquidityReserve)
	tokenWithdrawnAmt := msg.WithdrawLiquidity.Mul(tokenReserveAmt).Div(liquidityReserve)

	irisWithdrawCoin := sdk.NewCoin(sdk.IrisAtto, irisWithdrawnAmt)
	tokenWithdrawCoin := sdk.NewCoin(msg.MinToken.Denom, tokenWithdrawnAmt)
	deductUniCoin := sdk.NewCoin(reservePoolName, msg.WithdrawLiquidity)

	if irisWithdrawCoin.Amount.LT(msg.MinIrisAmt) {
		return types.ErrConstraintNotMet(fmt.Sprintf("The amount of iris available [%s] is less than the minimum amount specified [%s] by the user.", irisWithdrawCoin.String(), sdk.NewCoin(sdk.IrisAtto, msg.MinIrisAmt).String()))
	}
	if tokenWithdrawCoin.Amount.LT(msg.MinToken.Amount) {
		return types.ErrConstraintNotMet(fmt.Sprintf("The amount of token available [%s] is less than the minimum amount specified [%s] by the user.", tokenWithdrawCoin.String(), msg.MinToken.String()))
	}
	poolAddr := getReservePoolAddr(reservePoolName)
	return k.removeLiquidity(ctx, poolAddr, msg.Sender, deductUniCoin, irisWithdrawCoin, tokenWithdrawCoin)
}

func (k Keeper) removeLiquidity(ctx sdk.Context, poolAddr, sender sdk.AccAddress, deductUniCoin, irisWithdrawCoin, tokenWithdrawCoin sdk.Coin) sdk.Error {
	// burn liquidity from reserve Pool
	deltaCoins := sdk.NewCoins(deductUniCoin)
	_, _, err := k.bk.SubtractCoins(ctx, poolAddr, deltaCoins)
	if err != nil {
		return err
	}
	// burn liquidity from account
	_, _, err = k.bk.SubtractCoins(ctx, sender, deltaCoins)
	if err != nil {
		return err
	}
	// transfer withdrawn liquidity from coinswaps ModuleAccount to sender's account
	coins := sdk.NewCoins(irisWithdrawCoin, tokenWithdrawCoin)
	_, err = k.bk.SendCoins(ctx, poolAddr, sender, coins)
	return err
}

// GetReservePool returns the total balance of an reserve pool at the
// provided denomination.
func (k Keeper) GetReservePool(ctx sdk.Context, reservePoolName string) (coins sdk.Coins) {
	swapPoolAccAddr := getReservePoolAddr(reservePoolName)
	acc := k.ak.GetAccount(ctx, swapPoolAccAddr)
	if acc == nil {
		return nil
	}
	return acc.GetCoins()
}

// GetParams gets the parameters for the coinswap module.
func (k Keeper) GetParams(ctx sdk.Context) types.Params {
	var swapParams types.Params
	k.paramSpace.GetParamSet(ctx, &swapParams)
	return swapParams
}

// SetParams sets the parameters for the coinswap module.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramSpace.SetParamSet(ctx, &params)
}

func (k Keeper) Init(ctx sdk.Context) {
	paramSet := types.DefaultParams()
	k.paramSpace.SetParamSet(ctx, &paramSet)
}

func getReservePoolAddr(uniDenom string) sdk.AccAddress {
	return sdk.AccAddress(crypto.AddressHash([]byte(uniDenom)))
}
