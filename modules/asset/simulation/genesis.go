package simulation

// DONTCOVER

import (
	"fmt"
	"github.com/irisnet/irishub/modules/asset/types"
	"math/rand"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
)

// Simulation parameter constants
const (
	AssetTaxRate         = "asset_tax_rate"
	IssueTokenBaseFee    = "issue_token_base_fee"
	MintTokenFeeRatio    = "mint_token_fee_ratio"
	CreateGatewayBaseFee = "create_gateway_base_fee"
	GatewayAssetFeeRatio = "gateway_asset_fee_ratio"
	AssetFeeDenom        = "asset_fee_denom"
)

// RandomDec randomized sdk.RandomDec
func RandomDec(r *rand.Rand) sdk.Dec {
	return sdk.NewDec(r.Int63())
}

// RandomInt randomized sdk.Int
func RandomInt(r *rand.Rand) sdk.Int {
	return sdk.NewInt(r.Int63())
}

// RandomizedGenState generates a random GenesisState for bank
func RandomizedGenState(simState *module.SimulationState) {

	var assetTaxRate sdk.Dec
	var issueTokenBaseFee sdk.Int
	var mintTokenFeeRatio sdk.Dec
	var createGatewayBaseFee sdk.Int
	var gatewayAssetFeeRatio sdk.Dec
	simState.AppParams.GetOrGenerate(
		simState.Cdc, AssetTaxRate, &assetTaxRate, simState.Rand,
		func(r *rand.Rand) { assetTaxRate = RandomDec(r) },
	)

	simState.AppParams.GetOrGenerate(
		simState.Cdc, IssueTokenBaseFee, &issueTokenBaseFee, simState.Rand,
		func(r *rand.Rand) { issueTokenBaseFee = RandomInt(r) },
	)

	simState.AppParams.GetOrGenerate(
		simState.Cdc, MintTokenFeeRatio, &mintTokenFeeRatio, simState.Rand,
		func(r *rand.Rand) { mintTokenFeeRatio = RandomDec(r) },
	)

	simState.AppParams.GetOrGenerate(
		simState.Cdc, CreateGatewayBaseFee, &createGatewayBaseFee, simState.Rand,
		func(r *rand.Rand) { createGatewayBaseFee = RandomInt(r) },
	)

	simState.AppParams.GetOrGenerate(
		simState.Cdc, GatewayAssetFeeRatio, &gatewayAssetFeeRatio, simState.Rand,
		func(r *rand.Rand) { gatewayAssetFeeRatio = RandomDec(r) },
	)

	assetGenesis := types.NewGenesisState(
		types.NewParams(assetTaxRate, issueTokenBaseFee, mintTokenFeeRatio, ""),
		types.Tokens{},
	)

	fmt.Printf("Selected randomly generated bank parameters:\n%s\n", codec.MustMarshalJSONIndent(simState.Cdc, assetGenesis))
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(assetGenesis)
}
