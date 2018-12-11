package upgrade

import (
	sdk "github.com/irisnet/irishub/types"
	"github.com/irisnet/irishub/app/protocol"
	"fmt"
	"github.com/irisnet/irishub/modules/upgrade/params"
	"github.com/irisnet/irishub/modules/params"
)

// GenesisState - all upgrade state that must be provided at genesis
type GenesisState struct {
	SwitchPeriod int64    `json:"switch_period"`
}

// InitGenesis - build the genesis version For first Version
func InitGenesis(ctx sdk.Context, k Keeper, router protocol.Router, data GenesisState) {

	RegisterModuleList(router)

	moduleList, found := GetModuleListFromBucket(0)
	fmt.Println(moduleList)
	if !found {
		panic("No module list info found for genesis version")
	}

	genesisVersion := NewVersion(0, 0, 0, moduleList)
	k.AddNewVersion(ctx, genesisVersion)

	params.InitGenesisParameter(&upgradeparams.ProposalAcceptHeightParameter, ctx, -1)
	params.InitGenesisParameter(&upgradeparams.CurrentUpgradeProposalIdParameter, ctx, 0)
	params.InitGenesisParameter(&upgradeparams.SwitchPeriodParameter, ctx, data.SwitchPeriod)

	InitGenesis_commitID(ctx, k)
}


// WriteGenesis - output genesis parameters
func WriteGenesis(ctx sdk.Context) GenesisState {

	return GenesisState{
		SwitchPeriod: upgradeparams.GetSwitchPeriod(ctx),
	}
}

// get raw genesis raw message for testing
func DefaultGenesisState() GenesisState {
	return GenesisState{
		SwitchPeriod: 57600,
	}
}

// get raw genesis raw message for testing
func DefaultGenesisStateForTest() GenesisState {
	return GenesisState{
		SwitchPeriod: 15,
	}
}
