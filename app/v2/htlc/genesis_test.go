package htlc

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/require"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"

	"github.com/irisnet/irishub/app/v1/auth"
	"github.com/irisnet/irishub/app/v1/bank"
	"github.com/irisnet/irishub/codec"
	"github.com/irisnet/irishub/store"
	sdk "github.com/irisnet/irishub/types"
)

func setupMultiStore() (sdk.MultiStore, *sdk.KVStoreKey, *sdk.KVStoreKey) {
	db := dbm.NewMemDB()
	accountKey := sdk.NewKVStoreKey("accountkey")
	htlcKey := sdk.NewKVStoreKey("htlckey")

	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(accountKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(htlcKey, sdk.StoreTypeIAVL, db)
	_ = ms.LoadLatestVersion()

	return ms, accountKey, htlcKey
}

func TestExportHTLCGenesis(t *testing.T) {
	ms, accountKey, htlcKey := setupMultiStore()

	cdc := codec.New()
	RegisterCodec(cdc)
	auth.RegisterBaseAccount(cdc)

	ak := auth.NewAccountKeeper(cdc, accountKey, auth.ProtoBaseAccount)
	bk := bank.NewBaseKeeper(cdc, ak)
	keeper := NewKeeper(cdc, htlcKey, bk, DefaultCodespace)

	// build context
	currentBlockHeight := int64(100)
	ctx := sdk.NewContext(ms, abci.Header{}, false, log.NewNopLogger())
	ctx = ctx.WithBlockHeight(currentBlockHeight)

	// define variables
	senderAddrs := []sdk.AccAddress{sdk.AccAddress("sender1"), sdk.AccAddress("sender2")}
	receiverAddrs := []sdk.AccAddress{sdk.AccAddress("receiver1"), sdk.AccAddress("receiver2")}
	receiverOnOtherChain := "receiverOnOtherChain"
	amount := sdk.NewCoins(sdk.NewCoin(sdk.IrisAtto, sdk.NewInt(0)))
	secret := []byte("___abcdefghijklmnopqrstuvwxyz___")
	timestamps := []uint64{uint64(1580000000), 0}
	hashLocks := [][]byte{GetHashLock(secret, timestamps[0]), GetHashLock(secret, timestamps[1])}
	timeLocks := []uint64{50, 100}
	expireHeights := []uint64{timeLocks[0] + uint64(ctx.BlockHeight()), timeLocks[1] + uint64(ctx.BlockHeight())}
	state := OPEN
	initSecret := make([]byte, 0)

	// construct HTLCs
	htlc1 := NewHTLC(senderAddrs[0], receiverAddrs[0], receiverOnOtherChain, amount, initSecret, timestamps[0], expireHeights[0], state)
	htlc2 := NewHTLC(senderAddrs[1], receiverAddrs[1], receiverOnOtherChain, amount, initSecret, timestamps[1], expireHeights[1], state)

	// create HTLCs
	_, err := keeper.CreateHTLC(ctx, htlc1, hashLocks[0])
	require.Nil(t, err)
	_, err = keeper.CreateHTLC(ctx, htlc2, hashLocks[1])
	require.Nil(t, err)

	// preceed to the new block
	newBlockHeight := int64(150)
	ctx = ctx.WithBlockHeight(newBlockHeight)
	BeginBlocker(ctx, keeper)

	// export genesis
	exportedGenesis := ExportGenesis(ctx, keeper)
	exportedHTLCs := exportedGenesis.PendingHTLCs
	require.Equal(t, 1, len(exportedHTLCs))

	for hashLockHex, htlc := range exportedHTLCs {
		// assert the state must be OPEN
		require.True(t, htlc.State == OPEN)

		hashLock, err := hex.DecodeString(hashLockHex)
		require.Nil(t, err)

		// assert the HTLC with the given hash lock exists
		htlcInStore, err := keeper.GetHTLC(ctx, hashLock)
		require.Nil(t, err)

		// assert the expiration height is new
		newExpireHeight := htlcInStore.ExpireHeight - uint64(newBlockHeight) + 1
		require.Equal(t, newExpireHeight, htlc.ExpireHeight)

		// assert the exported HTLC is consistant with the HTLC in store except for the expiration height
		htlcInStore.ExpireHeight = newExpireHeight
		require.Equal(t, htlcInStore, htlc)
	}

	// assert the expired HTLCs(htlc1) have been refunded
	htlc, _ := keeper.GetHTLC(ctx, hashLocks[0])
	require.Equal(t, REFUNDED, htlc.State)
}
