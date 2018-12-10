package app

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/irisnet/irishub/codec"
	sdk "github.com/irisnet/irishub/types"
	"github.com/irisnet/irishub/modules/auth"
	"github.com/irisnet/irishub/modules/bank"
	distr "github.com/irisnet/irishub/modules/distribution"
	"github.com/irisnet/irishub/modules/mint"
	"github.com/irisnet/irishub/modules/params"
	"github.com/irisnet/irishub/modules/slashing"
	"github.com/irisnet/irishub/modules/stake"
	bam "github.com/irisnet/irishub/baseapp"
	"github.com/irisnet/irishub/modules/arbitration"
	"github.com/irisnet/irishub/modules/arbitration/params"
	"github.com/irisnet/irishub/modules/gov"
	"github.com/irisnet/irishub/modules/gov/params"
	"github.com/irisnet/irishub/modules/record"
	"github.com/irisnet/irishub/modules/service"
	"github.com/irisnet/irishub/modules/service/params"
	"github.com/irisnet/irishub/modules/upgrade"
	"github.com/irisnet/irishub/modules/upgrade/params"
	"github.com/spf13/viper"
	abci "github.com/tendermint/tendermint/abci/types"
	cmn "github.com/tendermint/tendermint/libs/common"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
	"time"
	"github.com/irisnet/irishub/modules/guardian"
)

const (
	appName    = "IrisApp"
	FlagReplay = "replay"
)

// default home directories for expected binaries
var (
	DefaultLCDHome  = os.ExpandEnv("$HOME/.irislcd")
	DefaultCLIHome  = os.ExpandEnv("$HOME/.iriscli")
	DefaultNodeHome = os.ExpandEnv("$HOME/.iris")
)

// Extended ABCI application
type IrisApp struct {
	*bam.BaseApp
	cdc *codec.Codec

	// keys to access the substores
	keyMain          *sdk.KVStoreKey
	keyAccount       *sdk.KVStoreKey
	keyStake         *sdk.KVStoreKey
	tkeyStake        *sdk.TransientStoreKey
	keySlashing      *sdk.KVStoreKey
	keyMint          *sdk.KVStoreKey
	keyDistr         *sdk.KVStoreKey
	tkeyDistr        *sdk.TransientStoreKey
	keyGov           *sdk.KVStoreKey
	keyFeeCollection *sdk.KVStoreKey
	keyParams        *sdk.KVStoreKey
	tkeyParams       *sdk.TransientStoreKey
	keyUpgrade       *sdk.KVStoreKey
	keyService       *sdk.KVStoreKey
	keyGuardian      *sdk.KVStoreKey
	keyRecord        *sdk.KVStoreKey

	// Manage getting and setting accounts
	accountMapper       auth.AccountKeeper
	feeCollectionKeeper auth.FeeCollectionKeeper
	bankKeeper          bank.Keeper
	stakeKeeper         stake.Keeper
	slashingKeeper      slashing.Keeper
	mintKeeper          mint.Keeper
	distrKeeper         distr.Keeper
	govKeeper           gov.Keeper
	paramsKeeper        params.Keeper
	upgradeKeeper       upgrade.Keeper
	serviceKeeper       service.Keeper
	guardianKeeper      guardian.Keeper
	recordKeeper        record.Keeper

	// fee manager
	feeManager auth.FeeManager
	hookHub    HookHub // handle Hook callback of any version modules
}

func NewIrisApp(logger log.Logger, db dbm.DB, traceStore io.Writer, baseAppOptions ...func(*bam.BaseApp)) *IrisApp {
	cdc := MakeCodec()

	bApp := bam.NewBaseApp(appName, logger, db, auth.DefaultTxDecoder(cdc), baseAppOptions...)
	bApp.SetCommitMultiStoreTracer(traceStore)

	// create your application object
	var app = &IrisApp{
		BaseApp:          bApp,
		cdc:              cdc,
		keyMain:          sdk.NewKVStoreKey("main"),
		keyAccount:       sdk.NewKVStoreKey("acc"),
		keyStake:         sdk.NewKVStoreKey("stake"),
		tkeyStake:        sdk.NewTransientStoreKey("transient_stake"),
		keyMint:          sdk.NewKVStoreKey("mint"),
		keyDistr:         sdk.NewKVStoreKey("distr"),
		tkeyDistr:        sdk.NewTransientStoreKey("transient_distr"),
		keySlashing:      sdk.NewKVStoreKey("slashing"),
		keyGov:           sdk.NewKVStoreKey("gov"),
		keyRecord:        sdk.NewKVStoreKey("record"),
		keyFeeCollection: sdk.NewKVStoreKey("fee"),
		keyParams:        sdk.NewKVStoreKey("params"),
		tkeyParams:       sdk.NewTransientStoreKey("transient_params"),
		keyUpgrade:       sdk.NewKVStoreKey("upgrade"),
		keyService:       sdk.NewKVStoreKey("service"),
		keyGuardian:      sdk.NewKVStoreKey("guardian"),
	}

	var lastHeight int64
	if viper.GetBool(FlagReplay) {
		lastHeight = bam.Replay(app.Logger)
	}

	app.initKeeper()
	app.wireRouterForAllVersion()
	app.mountStoreAndSetupBaseApp(lastHeight)
	app.registerParams()

	return app
}

// custom tx codec
func MakeCodec() *codec.Codec {
	var cdc = codec.New()
	bank.RegisterCodec(cdc)
	stake.RegisterCodec(cdc)
	distr.RegisterCodec(cdc)
	slashing.RegisterCodec(cdc)
	gov.RegisterCodec(cdc)
	record.RegisterCodec(cdc)
	upgrade.RegisterCodec(cdc)
	service.RegisterCodec(cdc)
	guardian.RegisterCodec(cdc)
	auth.RegisterCodec(cdc)
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	return cdc
}

func (app *IrisApp) initKeeper() {
	// define the AccountKeeper
	app.accountMapper = auth.NewAccountKeeper(
		app.cdc,
		app.keyAccount,        // target store
		auth.ProtoBaseAccount, // prototype
	)

	// add handlers
	app.guardianKeeper = guardian.NewKeeper(
		app.cdc,
		app.keyGuardian,
		guardian.DefaultCodespace,
	)
	app.bankKeeper = bank.NewBaseKeeper(app.accountMapper)
	app.feeCollectionKeeper = auth.NewFeeCollectionKeeper(
		app.cdc,
		app.keyFeeCollection,
	)
	app.paramsKeeper = params.NewKeeper(
		app.cdc,
		app.keyParams, app.tkeyParams,
	)
	stakeKeeper := stake.NewKeeper(
		app.cdc,
		app.keyStake, app.tkeyStake,
		app.bankKeeper, app.paramsKeeper.Subspace(stake.DefaultParamspace),
		stake.DefaultCodespace,
	)
	app.mintKeeper = mint.NewKeeper(app.cdc, app.keyMint,
		app.paramsKeeper.Subspace(mint.DefaultParamspace),
		&stakeKeeper, app.feeCollectionKeeper,
	)
	app.distrKeeper = distr.NewKeeper(
		app.cdc,
		app.keyDistr,
		app.paramsKeeper.Subspace(distr.DefaultParamspace),
		app.bankKeeper, &stakeKeeper, app.feeCollectionKeeper,
		distr.DefaultCodespace,
	)
	app.slashingKeeper = slashing.NewKeeper(
		app.cdc,
		app.keySlashing,
		&stakeKeeper, app.paramsKeeper.Subspace(slashing.DefaultParamspace),
		slashing.DefaultCodespace,
	)

	app.govKeeper = gov.NewKeeper(
		app.cdc,
		app.keyGov,
		app.bankKeeper, &stakeKeeper,
		gov.DefaultCodespace,
	)

	app.recordKeeper = record.NewKeeper(
		app.cdc,
		app.keyRecord,
		record.DefaultCodespace,
	)
	app.serviceKeeper = service.NewKeeper(
		app.cdc,
		app.keyService,
		app.bankKeeper,
		app.guardianKeeper,
		service.DefaultCodespace,
	)
	app.upgradeKeeper = upgrade.NewKeeper(
		app.cdc,
		app.keyUpgrade, app.stakeKeeper,
	)

	app.hookHub = NewHooksHub(app.upgradeKeeper)
	// register the staking hookHub
	// NOTE: stakeKeeper above are passed by reference,
	// so that it can be modified like below:
	app.stakeKeeper = *stakeKeeper.SetHooks(app.hookHub)
}

func (app *IrisApp) mountStoreAndSetupBaseApp(lastHeight int64) {
	app.feeManager = auth.NewFeeManager(app.paramsKeeper.Subspace("Fee"))

	// initialize BaseApp
	app.MountStoresIAVL(app.keyMain, app.keyAccount, app.keyStake, app.keySlashing, app.keyGov, app.keyMint, app.keyDistr,
		app.keyFeeCollection, app.keyParams, app.keyUpgrade, app.keyRecord, app.keyService, app.keyGuardian)
	app.SetInitChainer(app.initChainer)
	app.SetBeginBlocker(app.BeginBlocker)
	app.SetAnteHandler(auth.NewAnteHandler(app.accountMapper, app.feeCollectionKeeper))
	app.MountStoresTransient(app.tkeyParams, app.tkeyStake, app.tkeyDistr)
	app.SetFeeRefundHandler(auth.NewFeeRefundHandler(app.accountMapper, app.feeCollectionKeeper, app.feeManager))
	app.SetFeePreprocessHandler(auth.NewFeePreprocessHandler(app.feeManager))
	app.SetEndBlocker(app.EndBlocker)
	app.SetRunMsg(app.runMsgs)

	var err error
	if viper.GetBool(FlagReplay) {
		err = app.LoadVersion(lastHeight, app.keyMain, true)
	} else {
		err = app.LoadLatestVersion(app.keyMain)
	}
	if err != nil {
		cmn.Exit(err.Error())
	}

	upgrade.RegisterModuleList(app.Router())
	app.upgradeKeeper.RefreshVersionList(app.GetKVStore(app.keyUpgrade))
}

func (app *IrisApp) registerParams() {
	params.SetParamReadWriter(app.paramsKeeper.Subspace(params.SignalParamspace).WithTypeTable(
		params.NewTypeTable(
			upgradeparams.CurrentUpgradeProposalIdParameter.GetStoreKey(), uint64((0)),
			upgradeparams.ProposalAcceptHeightParameter.GetStoreKey(), int64(0),
			upgradeparams.SwitchPeriodParameter.GetStoreKey(), int64(0),
		)),
		&upgradeparams.CurrentUpgradeProposalIdParameter,
		&upgradeparams.ProposalAcceptHeightParameter,
		&upgradeparams.SwitchPeriodParameter)

	params.SetParamReadWriter(app.paramsKeeper.Subspace(params.GovParamspace).WithTypeTable(
		params.NewTypeTable(
			govparams.DepositProcedureParameter.GetStoreKey(), govparams.DepositProcedure{},
			govparams.VotingProcedureParameter.GetStoreKey(), govparams.VotingProcedure{},
			govparams.TallyingProcedureParameter.GetStoreKey(), govparams.TallyingProcedure{},
			serviceparams.MaxRequestTimeoutParameter.GetStoreKey(), int64(0),
			serviceparams.MinDepositMultipleParameter.GetStoreKey(), int64(0),
			arbitrationparams.ComplaintRetrospectParameter.GetStoreKey(), time.Duration(0),
			arbitrationparams.ArbitrationTimelimitParameter.GetStoreKey(), time.Duration(0),
		)),
		&govparams.DepositProcedureParameter,
		&govparams.VotingProcedureParameter,
		&govparams.TallyingProcedureParameter,
		&serviceparams.MaxRequestTimeoutParameter,
		&serviceparams.MinDepositMultipleParameter,
		&arbitrationparams.ComplaintRetrospectParameter,
		&arbitrationparams.ArbitrationTimelimitParameter)

	params.RegisterGovParamMapping(
		&govparams.DepositProcedureParameter,
		&govparams.VotingProcedureParameter,
		&govparams.TallyingProcedureParameter,
		&serviceparams.MaxRequestTimeoutParameter,
		&serviceparams.MinDepositMultipleParameter)
}

func (app *IrisApp) LoadHeight(height int64) error {
	return app.LoadVersion(height, app.keyMain, false)
}

// application updates every end block
func (app *IrisApp) BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
	tags := slashing.BeginBlocker(ctx, req, app.slashingKeeper)

	// distribute rewards from previous block
	distr.BeginBlocker(ctx, req, app.distrKeeper)

	// mint new tokens for this new block
	mint.BeginBlocker(ctx, app.mintKeeper)

	return abci.ResponseBeginBlock{
		Tags: tags.ToKVPairs(),
	}
}

// application updates every end block
func (app *IrisApp) EndBlocker(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
	tags := gov.EndBlocker(ctx, app.govKeeper)
	validatorUpdates := stake.EndBlocker(ctx, app.stakeKeeper)
	tags = tags.AppendTags(upgrade.EndBlocker(ctx, app.upgradeKeeper))
	tags = tags.AppendTags(service.EndBlocker(ctx, app.serviceKeeper))

	app.assertRuntimeInvariants()

	return abci.ResponseEndBlock{
		ValidatorUpdates: validatorUpdates,
		Tags:             tags,
	}
}

// custom logic for iris initialization
func (app *IrisApp) initChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	stateJSON := req.AppStateBytes

	var genesisFileState GenesisFileState
	err := app.cdc.UnmarshalJSON(stateJSON, &genesisFileState)
	if err != nil {
		panic(err)
	}
	genesisState := convertToGenesisState(genesisFileState)
	// sort by account number to maintain consistency
	sort.Slice(genesisState.Accounts, func(i, j int) bool {
		return genesisState.Accounts[i].AccountNumber < genesisState.Accounts[j].AccountNumber
	})

	// load the accounts
	for _, gacc := range genesisState.Accounts {
		acc := gacc.ToAccount()
		acc.AccountNumber = app.accountMapper.GetNextAccountNumber(ctx)
		app.accountMapper.SetAccount(ctx, acc)
	}

	upgrade.InitGenesis(ctx, app.upgradeKeeper, app.Router(), genesisState.UpgradeData)

	// load the initial stake information
	validators, err := stake.InitGenesis(ctx, app.stakeKeeper, genesisState.StakeData)
	if err != nil {
		panic(err)
	}
	gov.InitGenesis(ctx, app.govKeeper, genesisState.GovData)

	feeTokenGensisConfig := auth.FeeGenesisStateConfig{
		FeeTokenNative:    IrisCt.MinUnit.Denom,
		GasPriceThreshold: 20000000000, // 20(glue), 20*10^9, 1 glue = 10^9 lue/gas, 1 iris = 10^18 lue
	}

	// load the address to pubkey map
	auth.InitGenesis(ctx, app.feeCollectionKeeper, genesisState.AuthData, app.feeManager, feeTokenGensisConfig)
	slashing.InitGenesis(ctx, app.slashingKeeper, genesisState.SlashingData, genesisState.StakeData)
	mint.InitGenesis(ctx, app.mintKeeper, genesisState.MintData)
	distr.InitGenesis(ctx, app.distrKeeper, genesisState.DistrData)
	err = IrisValidateGenesisState(genesisState)
	if err != nil {
		panic(err) // TODO find a way to do this w/o panics
	}

	if len(genesisState.GenTxs) > 0 {
		for _, genTx := range genesisState.GenTxs {
			var tx auth.StdTx
			err = app.cdc.UnmarshalJSON(genTx, &tx)
			if err != nil {
				panic(err)
			}
			bz := app.cdc.MustMarshalBinaryLengthPrefixed(tx)
			res := app.BaseApp.DeliverTx(bz)
			if !res.IsOK() {
				panic(res.Log)
			}
		}

		validators = app.stakeKeeper.ApplyAndReturnValidatorSetUpdates(ctx)
	}

	// sanity check
	if len(req.Validators) > 0 {
		if len(req.Validators) != len(validators) {
			panic(fmt.Errorf("len(RequestInitChain.Validators) != len(validators) (%d != %d)",
				len(req.Validators), len(validators)))
		}
		sort.Sort(abci.ValidatorUpdates(req.Validators))
		sort.Sort(abci.ValidatorUpdates(validators))
		for i, val := range validators {
			if !val.Equal(req.Validators[i]) {
				panic(fmt.Errorf("validators[%d] != req.Validators[%d] ", i, i))
			}
		}
	}

	service.InitGenesis(ctx, app.serviceKeeper, genesisState.ServiceData)
	arbitration.InitGenesis(ctx, genesisState.ArbitrationData)
	guardian.InitGenesis(ctx, app.guardianKeeper, genesisState.GuardianData)

	return abci.ResponseInitChain{
		Validators: validators,
	}
}

// Iterates through msgs and executes them
func (app *IrisApp) runMsgs(ctx sdk.Context, msgs []sdk.Msg, mode bam.RunTxMode) (result sdk.Result) {
	// accumulate results
	logs := make([]string, 0, len(msgs))
	var data []byte   // NOTE: we just append them all (?!)
	var tags sdk.Tags // also just append them all
	var code sdk.CodeType
	var codespace sdk.CodespaceType
	for msgIdx, msg := range msgs {
		// Match route.
		var msgType string
		var err sdk.Error
		if ctx.BlockHeight() != 0 {
			msgType, err = app.upgradeKeeper.GetMsgTypeInCurrentVersion(ctx, msg)

			if err != nil {
				return err.Result()
			}

		} else {
			msgType = msg.Route()
		}

		handler := app.Router().Route(msgType)
		if handler == nil {
			return sdk.ErrUnknownRequest("Unrecognized Msg type: " + msgType).Result()
		}

		var msgResult sdk.Result
		if mode != bam.RunTxModeCheck {
			msgResult = handler(ctx, msg)
		}

		// NOTE: GasWanted is determined by ante handler and
		// GasUsed by the GasMeter

		// Append Data and Tags
		data = append(data, msgResult.Data...)
		tags = append(tags, msgResult.Tags...)

		// Stop execution and return on first failed message.
		if !msgResult.IsOK() {
			logs = append(logs, fmt.Sprintf("Msg %d failed: %s", msgIdx, msgResult.Log))
			code = msgResult.Code
			codespace = msgResult.Codespace
			break
		}

		// Construct usable logs in multi-message transactions.
		logs = append(logs, fmt.Sprintf("Msg %d: %s", msgIdx, msgResult.Log))
	}

	// Set the final gas values.
	result = sdk.Result{
		Code:      code,
		Codespace: codespace,
		Data:      data,
		Log:       strings.Join(logs, "\n"),
		GasUsed:   ctx.GasMeter().GasConsumed(),
		// TODO: FeeAmount/FeeDenom
		Tags: tags,
	}

	return result
}
