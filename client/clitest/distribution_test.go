package clitest

import (
	"fmt"
	"testing"

	"github.com/cosmos/cosmos-sdk/tests"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/irisnet/irishub/app"
	"github.com/stretchr/testify/require"
)

func init() {
	irisHome, iriscliHome = getTestingHomeDirs()
}

func TestIrisCLIDistribution(t *testing.T) {
	chainID, servAddr, port := initializeFixtures(t)

	flags := fmt.Sprintf("--home=%s --node=%v --chain-id=%v", iriscliHome, servAddr, chainID)

	// start iris server
	proc := tests.GoExecuteTWithStdout(t, fmt.Sprintf("iris start --home=%s --rpc.laddr=%v", irisHome, servAddr))

	defer proc.Stop(false)
	tests.WaitForTMStart(port)
	tests.WaitForNextNBlocksTM(2, port)

	fooAddr, _ := executeGetAddrPK(t, fmt.Sprintf("iriscli keys show foo --output=json --home=%s", iriscliHome))
	barAddr, _ := executeGetAddrPK(t, fmt.Sprintf("iriscli keys show bar --output=json --home=%s", iriscliHome))

	fooAcc := executeGetAccount(t, fmt.Sprintf("iriscli bank account %s %v", fooAddr, flags))
	fooCoin := convertToIrisBaseAccount(t, fooAcc)
	require.Equal(t, "50iris", fooCoin)

	executeWrite(t, fmt.Sprintf("iriscli bank send %v --amount=2iris --to=%s --from=foo --gas=10000 --fee=10iris", flags, barAddr), app.DefaultKeyPass)
	tests.WaitForNextNBlocksTM(2, port)

	valAddr := sdk.ValAddress(fooAddr).String()

	withdrawAddress, err := tests.ExecuteT(t, fmt.Sprintf("iriscli distribution withdraw-address %s %s", fooAddr, flags), "")
	require.Empty(t, err)
	require.Equal(t, "No withdraw address specified. If the delegator does have valid delegations, then the withdraw address should be the same as the delegator address", withdrawAddress)

	executeWrite(t, fmt.Sprintf("iriscli distribution set-withdraw-addr %s --from=foo --fee=0.004iris %s", barAddr, flags), app.DefaultKeyPass)
	tests.WaitForNextNBlocksTM(2, port)

	withdrawAddress, err = tests.ExecuteT(t, fmt.Sprintf("iriscli distribution withdraw-address %s %s", fooAddr, flags), "")
	require.Empty(t, err)
	require.Equal(t, barAddr.String(), withdrawAddress)

	barAcc := executeGetAccount(t, fmt.Sprintf("iriscli bank account %s %v", barAddr, flags))
	barCoin := convertToIrisBaseAccount(t, barAcc)
	require.Equal(t, "2iris", barCoin)

	ddiList := executeGetDelegatorDistrInfo(t, fmt.Sprintf("iriscli distribution delegator-distr-info %s %s", fooAddr, flags))
	require.Equal(t, 1, len(ddiList))
	require.Equal(t, int64(0), ddiList[0].DelPoolWithdrawalHeight)
	require.Equal(t, fooAddr, ddiList[0].DelegatorAddr)
	require.Equal(t, valAddr, ddiList[0].ValOperatorAddr.String())

	ddi := executeGetDelegationDistrInfo(t, fmt.Sprintf("iriscli distribution delegation-distr-info --address-delegator=%s --address-validator=%s %s", fooAddr, valAddr, flags))
	require.Equal(t, int64(0), ddi.DelPoolWithdrawalHeight)
	require.Equal(t, fooAddr, ddi.DelegatorAddr)
	require.Equal(t, valAddr, ddi.ValOperatorAddr.String())

	vdi := executeGetValidatorDistrInfo(t, fmt.Sprintf("iriscli distribution validator-distr-info %s %s", valAddr, flags))
	require.Equal(t, valAddr, vdi.OperatorAddr.String())
	require.Equal(t, int64(0), vdi.FeePoolWithdrawalHeight)
	require.Equal(t, "0.2084262892iris", vdi.DelPool)
	require.Equal(t, "0.0000000000iris", vdi.ValCommission)

	executeWrite(t, fmt.Sprintf("iriscli distribution withdraw-rewards --from=foo --fee=0.004iris %s", flags), app.DefaultKeyPass)
	tests.WaitForNextNBlocksTM(2, port)

	barAcc = executeGetAccount(t, fmt.Sprintf("iriscli bank account %s %v", barAddr, flags))
	barCoin = convertToIrisBaseAccount(t, barAcc)
	num := getAmountFromCoinStr(barCoin)

	if num > 6.1 || num < 6.0 {
		t.Error("Test Failed: (6.0, 6.1) expected, recieved: {}", num)
	}
}

func TestIrisCLIWithdrawReward(t *testing.T) {
	chainID, servAddr, port := initializeFixtures(t)

	flags := fmt.Sprintf("--home=%s --node=%v --chain-id=%v", iriscliHome, servAddr, chainID)

	// start iris server
	proc := tests.GoExecuteTWithStdout(t, fmt.Sprintf("iris start --home=%s --rpc.laddr=%v", irisHome, servAddr))

	defer proc.Stop(false)
	tests.WaitForTMStart(port)
	tests.WaitForNextNBlocksTM(2, port)

	fooAddr, _ := executeGetAddrPK(t, fmt.Sprintf("iriscli keys show foo --output=json --home=%s", iriscliHome))
	barAddr, _ := executeGetAddrPK(t, fmt.Sprintf("iriscli keys show bar --output=json --home=%s", iriscliHome))

	fooAcc := executeGetAccount(t, fmt.Sprintf("iriscli bank account %s %v", fooAddr, flags))
	fooCoin := convertToIrisBaseAccount(t, fooAcc)
	require.Equal(t, "50iris", fooCoin)

	executeWrite(t, fmt.Sprintf("iriscli bank send %v --amount=2iris --to=%s --from=foo --gas=10000 --fee=30iris", flags, barAddr), app.DefaultKeyPass)
	tests.WaitForNextNBlocksTM(2, port)

	valAddr := sdk.ValAddress(fooAddr).String()

	executeWrite(t, fmt.Sprintf("iriscli distribution set-withdraw-addr %s --from=foo --fee=0.004iris %s", barAddr, flags), app.DefaultKeyPass)
	tests.WaitForNextNBlocksTM(2, port)

	vdi := executeGetValidatorDistrInfo(t, fmt.Sprintf("iriscli distribution validator-distr-info %s %s", valAddr, flags))
	require.Equal(t, valAddr, vdi.OperatorAddr.String())
	require.Equal(t, int64(0), vdi.FeePoolWithdrawalHeight)
	require.Equal(t, "0.6251262892iris", vdi.DelPool)
	require.Equal(t, "0.0000000000iris", vdi.ValCommission)

	executeWrite(t, fmt.Sprintf("iriscli distribution withdraw-rewards --only-from-validator=%s --from=foo --fee=0.004iris %s", valAddr, flags), app.DefaultKeyPass)
	tests.WaitForNextNBlocksTM(2, port)

	barAcc := executeGetAccount(t, fmt.Sprintf("iriscli bank account %s %v", barAddr, flags))
	barCoin := convertToIrisBaseAccount(t, barAcc)
	num := getAmountFromCoinStr(barCoin)

	if num > 14.3 || num <= 14.2 {
		t.Error("Test Failed: (14.2, 14.3) expected, recieved: {}", num)
	}

	vdi = executeGetValidatorDistrInfo(t, fmt.Sprintf("iriscli distribution validator-distr-info %s %s", valAddr, flags))
	require.Equal(t, valAddr, vdi.OperatorAddr.String())
	require.Equal(t, "0.0000000000iris", vdi.ValCommission)

	executeWrite(t, fmt.Sprintf("iriscli distribution withdraw-rewards --is-validator=true --from=foo --fee=0.004iris %s", flags), app.DefaultKeyPass)
	tests.WaitForNextNBlocksTM(2, port)

	barAcc = executeGetAccount(t, fmt.Sprintf("iriscli bank account %s %v", barAddr, flags))
	barCoin = convertToIrisBaseAccount(t, barAcc)
	numNew := getAmountFromCoinStr(barCoin)

	if numNew <= num {
		t.Error("Test Failed: if --is-validator is true, more reward should be return")
	}
}
