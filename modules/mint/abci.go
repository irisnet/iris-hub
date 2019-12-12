package mint

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/irisnet/irishub/modules/mint/internal/types"
)

// Called every block, process inflation on the first block of every hour
func BeginBlocker(ctx sdk.Context, k Keeper) {
	ctx = ctx.WithLogger(ctx.Logger().With("handler", "beginBlock").With("module", "iris/mint"))
	logger := ctx.Logger()
	// Get block BFT time and block height
	blockTime := ctx.BlockHeader().Time
	minter := k.GetMinter(ctx)
	if ctx.BlockHeight() <= 1 { // don't inflate token in the first block
		minter.LastUpdate = blockTime
		k.SetMinter(ctx, minter)
		return
	}

	// Calculate block mint amount
	params := k.GetParamSet(ctx)
	logger.Info("Mint parameters", "inflation_rate", params.Inflation.String(), "mint_denom", params.MintDenom)

	mintedCoin := minter.BlockProvision(params)
	logger.Info("Mint result", "block_provisions", mintedCoin.String(), "time", blockTime.String())

	mintedCoins := sdk.NewCoins(mintedCoin)
	// mint coins to submodule account
	if err := k.MintCoins(ctx, mintedCoins); err != nil {
		panic(err)
	}

	// send the minted coins to the fee collector account
	if err := k.AddCollectedFees(ctx, mintedCoins); err != nil {
		panic(err)
	}

	// Update last block BFT time
	lastInflationTime := minter.LastUpdate
	minter.LastUpdate = blockTime
	k.SetMinter(ctx, minter)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeMint,
			sdk.NewAttribute(types.AttributeKeyLastInflationTime, lastInflationTime.String()),
			sdk.NewAttribute(types.AttributeKeyInflationTime, blockTime.String()),
			sdk.NewAttribute(types.AttributeKeyMintCoin, mintedCoin.Amount.String()),
		),
	)
}
