package simulation

import (
	"encoding/json"
	"math/rand"
	"testing"

	"github.com/tendermint/tendermint/crypto"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/irisnet/irishub/simulation/mock"
	"github.com/irisnet/irishub/simulation/mock/simulation"
)

func TestBankWithRandomMessages(t *testing.T) {
	mapp := mock.NewApp()

	bank.RegisterCodec(mapp.Cdc)
	mapper := mapp.AccountKeeper
	coinKeeper := bank.NewKeeper(mapper)
	mapp.Router().AddRoute("bank", []*sdk.KVStoreKey{mapp.KeyAccount}, bank.NewHandler(coinKeeper))

	err := mapp.CompleteSetup([]*sdk.KVStoreKey{})
	if err != nil {
		panic(err)
	}

	appStateFn := func(r *rand.Rand, keys []crypto.PrivKey, accs []sdk.AccAddress) json.RawMessage {
		mock.RandomSetGenesis(r, mapp, accs, []string{"iris-atto"})
		return json.RawMessage("{}")
	}

	simulation.Simulate(
		t, mapp.BaseApp, appStateFn,
		[]simulation.TestAndRunTx{
			TestAndRunSingleInputMsgSend(mapper),
		},
		[]simulation.RandSetup{},
		[]simulation.Invariant{
			NonnegativeBalanceInvariant(mapper),
			TotalCoinsInvariant(mapper, func() sdk.Coins { return mapp.TotalCoinsSupply }),
		},
		30, 30,
		false,
	)
}
