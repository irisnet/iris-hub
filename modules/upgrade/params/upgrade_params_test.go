package upgradeparams

import (
	"testing"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/stretchr/testify/require"
	"github.com/cosmos/cosmos-sdk/store"

	abci "github.com/tendermint/tendermint/abci/types"
	dbm "github.com/tendermint/tendermint/libs/db"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/libs/log"
)

func defaultContext(key sdk.StoreKey) sdk.Context {
	db := dbm.NewMemDB()
	cms := store.NewCommitMultiStore(db)
	cms.MountStoreWithDB(key, sdk.StoreTypeIAVL, db)
	cms.LoadLatestVersion()
	ctx := sdk.NewContext(cms, abci.Header{}, false, log.NewNopLogger())
	return ctx
}

func TestCurrentUpgradeProposalIdParameter(t *testing.T) {
	skey := sdk.NewKVStoreKey("params")
	ctx := defaultContext(skey)
	paramKeeper := params.NewKeeper(wire.NewCodec(), skey)

	CurrentUpgradeProposalIdParameter.SetReadWriter(paramKeeper.Setter())
	find := CurrentUpgradeProposalIdParameter.LoadValue(ctx)
	require.Equal(t, find, false)

	CurrentUpgradeProposalIdParameter.InitGenesis(nil)
	require.Equal(t, int64(-1), CurrentUpgradeProposalIdParameter.Value)

	CurrentUpgradeProposalIdParameter.LoadValue(ctx)
	require.Equal(t, int64(-1), CurrentUpgradeProposalIdParameter.Value)

	CurrentUpgradeProposalIdParameter.Value = 3
	CurrentUpgradeProposalIdParameter.SaveValue(ctx)

	CurrentUpgradeProposalIdParameter.LoadValue(ctx)
	require.Equal(t, int64(3), CurrentUpgradeProposalIdParameter.Value)
}

func TestProposalAcceptHeightParameter(t *testing.T) {
	skey := sdk.NewKVStoreKey("params")
	ctx := defaultContext(skey)
	paramKeeper := params.NewKeeper(wire.NewCodec(), skey)

	ProposalAcceptHeightParameter.SetReadWriter(paramKeeper.Setter())
	find := ProposalAcceptHeightParameter.LoadValue(ctx)
	require.Equal(t, find, false)

	ProposalAcceptHeightParameter.InitGenesis(nil)
	require.Equal(t, int64(-1), ProposalAcceptHeightParameter.Value)

	ProposalAcceptHeightParameter.LoadValue(ctx)
	require.Equal(t, int64(-1), ProposalAcceptHeightParameter.Value)

	ProposalAcceptHeightParameter.Value = 3
	ProposalAcceptHeightParameter.SaveValue(ctx)

	ProposalAcceptHeightParameter.LoadValue(ctx)
	require.Equal(t, int64(3), ProposalAcceptHeightParameter.Value)
}

func TestUpgradeParameterSetAndGet(t *testing.T) {
	skey := sdk.NewKVStoreKey("params")
	ctx := defaultContext(skey)
	paramKeeper := params.NewKeeper(wire.NewCodec(), skey)

	CurrentUpgradeProposalIdParameter.SetReadWriter(paramKeeper.Setter())
	find := CurrentUpgradeProposalIdParameter.LoadValue(ctx)
	require.Equal(t, find, false)

	ProposalAcceptHeightParameter.SetReadWriter(paramKeeper.Setter())
	find = ProposalAcceptHeightParameter.LoadValue(ctx)
	require.Equal(t, find, false)

	SetCurrentUpgradeProposalId(ctx,5)
	require.Equal(t,int64(5),GetCurrentUpgradeProposalId(ctx))
	SetProposalAcceptHeight(ctx,100)
	require.Equal(t, int64(100),GetProposalAcceptHeight(ctx) )
}
