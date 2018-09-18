package govparams

import (
	"encoding/json"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
)

var DepositProcedureParameter DepositProcedureParam

// Procedure around Deposits for governance
type DepositProcedure struct {
	MinDeposit       sdk.Coins `json:"min_deposit"`        //  Minimum deposit for a proposal to enter voting period.
	MaxDepositPeriod int64     `json:"max_deposit_period"` //  Maximum period for Atom holders to deposit on a proposal. Initial value: 2 months
}

type DepositProcedureParam struct {
	Value   DepositProcedure
	psetter params.Setter
	pgetter params.Getter
}

func (param *DepositProcedureParam) InitGenesis(genesisState interface{}) {
	if value, ok := genesisState.(DepositProcedure); ok {
		param.Value = value
	} else {
		param.Value = DepositProcedure{
			MinDeposit:       sdk.Coins{sdk.NewInt64Coin("iris", 10)},
			MaxDepositPeriod: 1440}
	}
}

func (param *DepositProcedureParam) SetReadWriter(setter params.Setter) {
	param.psetter = setter
	param.pgetter = setter.Getter
}

func (param *DepositProcedureParam) GetStoreKey() string {
	return "Gov/gov/depositProcedure"

}

func (param *DepositProcedureParam) SaveValue(ctx sdk.Context) {
	param.psetter.Set(ctx, param.GetStoreKey(), param.Value)
}

func (param *DepositProcedureParam) LoadValue(ctx sdk.Context) bool {
	err := param.pgetter.Get(ctx, param.GetStoreKey(), &param.Value)
	if err != nil {
		return false
	}
	return true
}

func (param *DepositProcedureParam) ToJson() string {
	jsonBytes, _ := json.Marshal(param.Value)
	return string(jsonBytes)
}

func (param *DepositProcedureParam) Update(ctx sdk.Context, jsonStr string) {
	if err := json.Unmarshal([]byte(jsonStr), &param.Value); err == nil {
		param.SaveValue(ctx)
	}
}

func (param *DepositProcedureParam) Valid(jsonStr string) sdk.Error {

	var err error

	if err = json.Unmarshal([]byte(jsonStr), &param.Value); err == nil {

		if param.Value.MinDeposit[0].Denom != "iris" {
			return sdk.NewError(DefaultCodespace, CodeInvalidMinDepositDenom, fmt.Sprintf("It should be iris "))
		}

		if param.Value.MinDeposit[0].Amount.GT(sdk.NewInt(10)) && param.Value.MinDeposit[0].Amount.LT(sdk.NewInt(20000)) {
			return sdk.NewError(DefaultCodespace, CodeInvalidMinDepositAmount, fmt.Sprintf("MinDepositAmount should be larger than 10 and less than 20000"))
		}

		if param.Value.MaxDepositPeriod > 20 && param.Value.MaxDepositPeriod < 20000 {
			return sdk.NewError(DefaultCodespace, CodeInvalidDepositPeriod, fmt.Sprintf("MaxDepositPeriod should be larger than 20 and less than 20000"))
		}

		return nil

	}
	return sdk.NewError(DefaultCodespace, CodeInvalidMinDeposit, fmt.Sprintf("Json is not valid"))
}
