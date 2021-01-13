package app

import (
	"io"
	"os"
	"path/filepath"

	"github.com/CosmWasm/wasmd/x/wasm"
	"github.com/cosmos/cosmos-sdk/client"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/libs/log"
	tmos "github.com/tendermint/tendermint/libs/os"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	dbm "github.com/tendermint/tm-db"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client/grpc/tmservice"
	"github.com/cosmos/cosmos-sdk/client/rpc"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/server/api"
	"github.com/cosmos/cosmos-sdk/server/config"
	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	authrest "github.com/cosmos/cosmos-sdk/x/auth/client/rest"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authsims "github.com/cosmos/cosmos-sdk/x/auth/simulation"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/auth/vesting"
	"github.com/cosmos/cosmos-sdk/x/bank"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/capability"
	capabilitykeeper "github.com/cosmos/cosmos-sdk/x/capability/keeper"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	crisiskeeper "github.com/cosmos/cosmos-sdk/x/crisis/keeper"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	distrclient "github.com/cosmos/cosmos-sdk/x/distribution/client"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/cosmos/cosmos-sdk/x/evidence"
	evidencekeeper "github.com/cosmos/cosmos-sdk/x/evidence/keeper"
	evidencetypes "github.com/cosmos/cosmos-sdk/x/evidence/types"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/cosmos/cosmos-sdk/x/gov"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/cosmos/cosmos-sdk/x/ibc/applications/transfer"
	ibctransferkeeper "github.com/cosmos/cosmos-sdk/x/ibc/applications/transfer/keeper"
	ibctransfertypes "github.com/cosmos/cosmos-sdk/x/ibc/applications/transfer/types"
	ibc "github.com/cosmos/cosmos-sdk/x/ibc/core"
	ibcclient "github.com/cosmos/cosmos-sdk/x/ibc/core/02-client"
	porttypes "github.com/cosmos/cosmos-sdk/x/ibc/core/05-port/types"
	ibchost "github.com/cosmos/cosmos-sdk/x/ibc/core/24-host"
	ibckeeper "github.com/cosmos/cosmos-sdk/x/ibc/core/keeper"
	"github.com/cosmos/cosmos-sdk/x/params"
	paramsclient "github.com/cosmos/cosmos-sdk/x/params/client"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	paramproposal "github.com/cosmos/cosmos-sdk/x/params/types/proposal"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	slashingkeeper "github.com/cosmos/cosmos-sdk/x/slashing/keeper"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cosmos/cosmos-sdk/x/upgrade"
	upgradeclient "github.com/cosmos/cosmos-sdk/x/upgrade/client"
	upgradekeeper "github.com/cosmos/cosmos-sdk/x/upgrade/keeper"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	"github.com/irisnet/irismod/modules/coinswap"
	coinswapkeeper "github.com/irisnet/irismod/modules/coinswap/keeper"
	coinswaptypes "github.com/irisnet/irismod/modules/coinswap/types"
	"github.com/irisnet/irismod/modules/htlc"
	htlckeeper "github.com/irisnet/irismod/modules/htlc/keeper"
	htlctypes "github.com/irisnet/irismod/modules/htlc/types"
	"github.com/irisnet/irismod/modules/nft"
	nftkeeper "github.com/irisnet/irismod/modules/nft/keeper"
	nfttypes "github.com/irisnet/irismod/modules/nft/types"
	"github.com/irisnet/irismod/modules/oracle"
	oraclekeeper "github.com/irisnet/irismod/modules/oracle/keeper"
	oracletypes "github.com/irisnet/irismod/modules/oracle/types"
	"github.com/irisnet/irismod/modules/random"
	randomkeeper "github.com/irisnet/irismod/modules/random/keeper"
	randomtypes "github.com/irisnet/irismod/modules/random/types"
	"github.com/irisnet/irismod/modules/record"
	recordkeeper "github.com/irisnet/irismod/modules/record/keeper"
	recordtypes "github.com/irisnet/irismod/modules/record/types"
	"github.com/irisnet/irismod/modules/service"
	servicekeeper "github.com/irisnet/irismod/modules/service/keeper"
	servicetypes "github.com/irisnet/irismod/modules/service/types"
	"github.com/irisnet/irismod/modules/token"
	tokenkeeper "github.com/irisnet/irismod/modules/token/keeper"
	tokentypes "github.com/irisnet/irismod/modules/token/types"

	"github.com/irisnet/irishub/address"
	irisappparams "github.com/irisnet/irishub/app/params"
	"github.com/irisnet/irishub/lite"
	"github.com/irisnet/irishub/modules/guardian"
	guardiankeeper "github.com/irisnet/irishub/modules/guardian/keeper"
	guardiantypes "github.com/irisnet/irishub/modules/guardian/types"
	"github.com/irisnet/irishub/modules/legacy"
	"github.com/irisnet/irishub/modules/mint"
	mintkeeper "github.com/irisnet/irishub/modules/mint/keeper"
	minttypes "github.com/irisnet/irishub/modules/mint/types"

	store "github.com/cosmos/cosmos-sdk/store/types"
	sdkupgrade "github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

const appName = "IrisApp"

var (
	// DefaultNodeHome default home directories for the application daemon
	DefaultNodeHome string

	// ModuleBasics defines the module BasicManager is in charge of setting up basic,
	// non-dependant module elements, such as codec registration
	// and genesis verification.
	ModuleBasics = module.NewBasicManager(
		auth.AppModuleBasic{},
		genutil.AppModuleBasic{},
		bank.AppModuleBasic{},
		capability.AppModuleBasic{},
		staking.AppModuleBasic{},
		mint.AppModuleBasic{},
		distr.AppModuleBasic{},
		gov.NewAppModuleBasic(
			paramsclient.ProposalHandler, distrclient.ProposalHandler, upgradeclient.ProposalHandler, upgradeclient.CancelProposalHandler,
		),
		params.AppModuleBasic{},
		crisis.AppModuleBasic{},
		slashing.AppModuleBasic{},
		ibc.AppModuleBasic{},
		upgrade.AppModuleBasic{},
		evidence.AppModuleBasic{},
		transfer.AppModuleBasic{},
		vesting.AppModuleBasic{},
		guardian.AppModuleBasic{},
		token.AppModuleBasic{},
		record.AppModuleBasic{},
		nft.AppModuleBasic{},
		htlc.AppModuleBasic{},
		coinswap.AppModuleBasic{},
		service.AppModuleBasic{},
		oracle.AppModuleBasic{},
		random.AppModuleBasic{},
		legacy.AppModuleBasic{},
		wasm.AppModuleBasic{},
	)

	// module account permissions
	maccPerms = map[string][]string{
		authtypes.FeeCollectorName:     nil,
		distrtypes.ModuleName:          nil,
		minttypes.ModuleName:           {authtypes.Minter},
		stakingtypes.BondedPoolName:    {authtypes.Burner, authtypes.Staking},
		stakingtypes.NotBondedPoolName: {authtypes.Burner, authtypes.Staking},
		govtypes.ModuleName:            {authtypes.Burner},
		ibctransfertypes.ModuleName:    {authtypes.Minter, authtypes.Burner},
		tokentypes.ModuleName:          {authtypes.Minter, authtypes.Burner},
		htlctypes.ModuleName:           nil,
		coinswaptypes.ModuleName:       {authtypes.Minter, authtypes.Burner},
		servicetypes.DepositAccName:    {authtypes.Burner},
		servicetypes.RequestAccName:    nil,
	}

	// module accounts that are allowed to receive tokens
	allowedReceivingModAcc = map[string]bool{
		distrtypes.ModuleName: true,
	}

	nativeToken tokentypes.Token
)

// Verify app interface at compile time
var _ simapp.App = (*IrisApp)(nil)

// IrisApp extends an ABCI application, but with most of its parameters exported.
// They are exported for convenience in creating helper functions, as object
// capabilities aren't needed for testing.
type IrisApp struct {
	*baseapp.BaseApp
	cdc               *codec.LegacyAmino
	appCodec          codec.Marshaler
	interfaceRegistry types.InterfaceRegistry

	invCheckPeriod uint

	// keys to access the substores
	keys    map[string]*sdk.KVStoreKey
	tkeys   map[string]*sdk.TransientStoreKey
	memKeys map[string]*sdk.MemoryStoreKey

	// keepers
	accountKeeper    authkeeper.AccountKeeper
	bankKeeper       bankkeeper.Keeper
	capabilityKeeper *capabilitykeeper.Keeper
	stakingKeeper    stakingkeeper.Keeper
	slashingKeeper   slashingkeeper.Keeper
	mintKeeper       mintkeeper.Keeper
	distrKeeper      distrkeeper.Keeper
	govKeeper        govkeeper.Keeper
	crisisKeeper     crisiskeeper.Keeper
	upgradeKeeper    upgradekeeper.Keeper
	paramsKeeper     paramskeeper.Keeper
	ibcKeeper        *ibckeeper.Keeper // IBC Keeper must be a pointer in the app, so we can SetRouter on it correctly
	evidenceKeeper   evidencekeeper.Keeper
	transferKeeper   ibctransferkeeper.Keeper

	// make scoped keepers public for test purposes
	ScopedIBCKeeper      capabilitykeeper.ScopedKeeper
	ScopedTransferKeeper capabilitykeeper.ScopedKeeper
	ScopedIBCMockKeeper  capabilitykeeper.ScopedKeeper

	guardianKeeper guardiankeeper.Keeper
	tokenKeeper    tokenkeeper.Keeper
	recordKeeper   recordkeeper.Keeper
	nftKeeper      nftkeeper.Keeper
	htlcKeeper     htlckeeper.Keeper
	coinswapKeeper coinswapkeeper.Keeper
	serviceKeeper  servicekeeper.Keeper
	oracleKeeper   oraclekeeper.Keeper
	randomKeeper   randomkeeper.Keeper
	wasmKeeper     wasm.Keeper

	// the module manager
	mm *module.Manager

	// simulation manager
	sm *module.SimulationManager
}

func init() {
	address.ConfigureBech32Prefix()
	nativeToken = tokentypes.Token{
		Symbol:        "iris",
		Name:          "Irishub staking token",
		Scale:         6,
		MinUnit:       "uiris",
		InitialSupply: 2000000000,
		MaxSupply:     10000000000,
		Mintable:      true,
		Owner:         sdk.AccAddress(crypto.AddressHash([]byte(tokentypes.ModuleName))).String(),
	}

	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	DefaultNodeHome = filepath.Join(userHomeDir, ".iris")

	owner, err := sdk.AccAddressFromBech32(nativeToken.Owner)
	if err != nil {
		panic(err)
	}

	tokentypes.SetNativeToken(
		nativeToken.Symbol,
		nativeToken.Name,
		nativeToken.MinUnit,
		nativeToken.Scale,
		nativeToken.InitialSupply,
		nativeToken.MaxSupply,
		nativeToken.Mintable,
		owner,
	)
}

// NewIrisApp returns a reference to an initialized IrisApp.
func NewIrisApp(
	logger log.Logger, db dbm.DB, traceStore io.Writer, loadLatest bool, skipUpgradeHeights map[int64]bool,
	homePath string, invCheckPeriod uint, encodingConfig irisappparams.EncodingConfig,
	appOpts servertypes.AppOptions, baseAppOptions ...func(*baseapp.BaseApp),
) *IrisApp {

	// TODO: Remove cdc in favor of appCodec once all modules are migrated.
	appCodec := encodingConfig.Marshaler
	cdc := encodingConfig.Amino
	interfaceRegistry := encodingConfig.InterfaceRegistry

	bApp := baseapp.NewBaseApp(appName, logger, db, encodingConfig.TxConfig.TxDecoder(), baseAppOptions...)
	bApp.SetCommitMultiStoreTracer(traceStore)
	bApp.SetAppVersion(version.Version)
	bApp.SetInterfaceRegistry(interfaceRegistry)

	keys := sdk.NewKVStoreKeys(
		authtypes.StoreKey, banktypes.StoreKey, stakingtypes.StoreKey,
		minttypes.StoreKey, distrtypes.StoreKey, slashingtypes.StoreKey,
		govtypes.StoreKey, paramstypes.StoreKey, ibchost.StoreKey, upgradetypes.StoreKey,
		evidencetypes.StoreKey, ibctransfertypes.StoreKey, capabilitytypes.StoreKey,
		guardiantypes.StoreKey, tokentypes.StoreKey, nfttypes.StoreKey, htlctypes.StoreKey, recordtypes.StoreKey,
		coinswaptypes.StoreKey, servicetypes.StoreKey, oracletypes.StoreKey, randomtypes.StoreKey, wasm.StoreKey,
	)
	tkeys := sdk.NewTransientStoreKeys(paramstypes.TStoreKey)
	memKeys := sdk.NewMemoryStoreKeys(capabilitytypes.MemStoreKey)

	app := &IrisApp{
		BaseApp:           bApp,
		cdc:               cdc,
		appCodec:          appCodec,
		interfaceRegistry: interfaceRegistry,
		invCheckPeriod:    invCheckPeriod,
		keys:              keys,
		tkeys:             tkeys,
		memKeys:           memKeys,
	}

	app.paramsKeeper = initParamsKeeper(appCodec, cdc, keys[paramstypes.StoreKey], tkeys[paramstypes.TStoreKey])

	// set the BaseApp's parameter store
	bApp.SetParamStore(app.paramsKeeper.Subspace(baseapp.Paramspace).WithKeyTable(paramskeeper.ConsensusParamsKeyTable()))

	// add capability keeper and ScopeToModule for ibc module
	app.capabilityKeeper = capabilitykeeper.NewKeeper(appCodec, keys[capabilitytypes.StoreKey], memKeys[capabilitytypes.MemStoreKey])
	scopedIBCKeeper := app.capabilityKeeper.ScopeToModule(ibchost.ModuleName)
	scopedTransferKeeper := app.capabilityKeeper.ScopeToModule(ibctransfertypes.ModuleName)

	// add keepers
	app.accountKeeper = authkeeper.NewAccountKeeper(
		appCodec, keys[authtypes.StoreKey], app.GetSubspace(authtypes.ModuleName), authtypes.ProtoBaseAccount, maccPerms,
	)
	app.bankKeeper = bankkeeper.NewBaseKeeper(
		appCodec, keys[banktypes.StoreKey], app.accountKeeper, app.GetSubspace(banktypes.ModuleName), app.BlockedAddrs(),
	)
	stakingKeeper := stakingkeeper.NewKeeper(
		appCodec, keys[stakingtypes.StoreKey], app.accountKeeper, app.bankKeeper, app.GetSubspace(stakingtypes.ModuleName),
	)
	app.mintKeeper = mintkeeper.NewKeeper(
		appCodec, keys[minttypes.StoreKey], app.GetSubspace(minttypes.ModuleName),
		app.accountKeeper, app.bankKeeper, authtypes.FeeCollectorName,
	)
	app.distrKeeper = distrkeeper.NewKeeper(
		appCodec, keys[distrtypes.StoreKey], app.GetSubspace(distrtypes.ModuleName), app.accountKeeper, app.bankKeeper,
		&stakingKeeper, authtypes.FeeCollectorName, app.ModuleAccountAddrs(),
	)
	app.slashingKeeper = slashingkeeper.NewKeeper(
		appCodec, keys[slashingtypes.StoreKey], &stakingKeeper, app.GetSubspace(slashingtypes.ModuleName),
	)
	app.crisisKeeper = crisiskeeper.NewKeeper(
		app.GetSubspace(crisistypes.ModuleName), invCheckPeriod, app.bankKeeper, authtypes.FeeCollectorName,
	)
	app.upgradeKeeper = upgradekeeper.NewKeeper(skipUpgradeHeights, keys[upgradetypes.StoreKey], appCodec, homePath)

	// register the staking hooks
	// NOTE: stakingKeeper above is passed by reference, so that it will contain these hooks
	app.stakingKeeper = *stakingKeeper.SetHooks(
		stakingtypes.NewMultiStakingHooks(app.distrKeeper.Hooks(), app.slashingKeeper.Hooks()),
	)

	// Create IBC Keeper
	app.ibcKeeper = ibckeeper.NewKeeper(
		appCodec, keys[ibchost.StoreKey], app.GetSubspace(ibchost.ModuleName), app.stakingKeeper, scopedIBCKeeper,
	)

	// register the proposal types
	govRouter := govtypes.NewRouter()
	govRouter.AddRoute(govtypes.RouterKey, govtypes.ProposalHandler).
		AddRoute(paramproposal.RouterKey, params.NewParamChangeProposalHandler(app.paramsKeeper)).
		AddRoute(distrtypes.RouterKey, distr.NewCommunityPoolSpendProposalHandler(app.distrKeeper)).
		AddRoute(upgradetypes.RouterKey, upgrade.NewSoftwareUpgradeProposalHandler(app.upgradeKeeper)).
		AddRoute(wasm.RouterKey, wasm.NewWasmProposalHandler(app.wasmKeeper, wasm.EnableAllProposals)).
		AddRoute(ibchost.RouterKey, ibcclient.NewClientUpdateProposalHandler(app.ibcKeeper.ClientKeeper))
	app.govKeeper = govkeeper.NewKeeper(
		appCodec, keys[govtypes.StoreKey], app.GetSubspace(govtypes.ModuleName), app.accountKeeper, app.bankKeeper,
		&stakingKeeper, govRouter,
	)

	// Create Transfer Keepers
	app.transferKeeper = ibctransferkeeper.NewKeeper(
		appCodec, keys[ibctransfertypes.StoreKey], app.GetSubspace(ibctransfertypes.ModuleName),
		app.ibcKeeper.ChannelKeeper, &app.ibcKeeper.PortKeeper,
		app.accountKeeper, app.bankKeeper, scopedTransferKeeper,
	)
	transferModule := transfer.NewAppModule(app.transferKeeper)

	// Create static IBC router, add transfer route, then set and seal it
	ibcRouter := porttypes.NewRouter()
	ibcRouter.AddRoute(ibctransfertypes.ModuleName, transferModule)
	app.ibcKeeper.SetRouter(ibcRouter)

	// create evidence keeper with router
	evidenceKeeper := evidencekeeper.NewKeeper(
		appCodec, keys[evidencetypes.StoreKey], &app.stakingKeeper, app.slashingKeeper,
	)
	// If evidence needs to be handled for the app, set routes in router here and seal
	app.evidenceKeeper = *evidenceKeeper

	app.guardianKeeper = guardiankeeper.NewKeeper(appCodec, keys[guardiantypes.StoreKey])
	app.tokenKeeper = tokenkeeper.NewKeeper(
		appCodec, keys[tokentypes.StoreKey], app.GetSubspace(tokentypes.ModuleName),
		app.bankKeeper, authtypes.FeeCollectorName,
	)
	app.recordKeeper = recordkeeper.NewKeeper(appCodec, keys[recordtypes.StoreKey])
	app.nftKeeper = nftkeeper.NewKeeper(appCodec, keys[nfttypes.StoreKey])

	app.htlcKeeper = htlckeeper.NewKeeper(appCodec, keys[htlctypes.StoreKey], app.accountKeeper, app.bankKeeper)

	app.coinswapKeeper = coinswapkeeper.NewKeeper(
		appCodec, keys[coinswaptypes.StoreKey], app.GetSubspace(coinswaptypes.ModuleName),
		app.bankKeeper, app.accountKeeper,
	)

	app.serviceKeeper = servicekeeper.NewKeeper(
		appCodec, keys[servicetypes.StoreKey], app.accountKeeper, app.bankKeeper,
		app.GetSubspace(servicetypes.ModuleName), authtypes.FeeCollectorName,
	)

	app.oracleKeeper = oraclekeeper.NewKeeper(
		appCodec, keys[oracletypes.StoreKey], app.GetSubspace(oracletypes.ModuleName),
		app.serviceKeeper,
	)

	app.randomKeeper = randomkeeper.NewKeeper(appCodec, keys[randomtypes.StoreKey], app.bankKeeper, app.serviceKeeper)

	wasmDir := filepath.Join(homePath, "wasm")
	wasmConfig, err := wasm.ReadWasmConfig(appOpts)
	if err != nil {
		panic("error while reading wasm config: " + err.Error())
	}
	// The last arguments can contain custom message handlers, and custom query handlers,
	// if we want to allow any custom callbacks
	supportedFeatures := "staking"
	app.wasmKeeper = wasm.NewKeeper(
		appCodec,
		keys[wasm.StoreKey],
		app.GetSubspace(wasm.ModuleName),
		app.accountKeeper,
		app.bankKeeper,
		app.stakingKeeper,
		app.distrKeeper,
		bApp.Router(),
		wasmDir,
		wasmConfig,
		supportedFeatures,
		nil,
		nil,
	)

	/****  Module Options ****/
	var skipGenesisInvariants = false
	opt := appOpts.Get(crisis.FlagSkipGenesisInvariants)
	if opt, ok := opt.(bool); ok {
		skipGenesisInvariants = opt
	}

	wasmModule := wasm.NewAppModule(&app.wasmKeeper, app.stakingKeeper)

	// NOTE: Any module instantiated in the module manager that is later modified
	// must be passed by reference here.
	app.mm = module.NewManager(
		genutil.NewAppModule(
			app.accountKeeper, app.stakingKeeper, app.BaseApp.DeliverTx,
			encodingConfig.TxConfig,
		),
		auth.NewAppModule(appCodec, app.accountKeeper, authsims.RandomGenesisAccounts),
		vesting.NewAppModule(app.accountKeeper, app.bankKeeper),
		bank.NewAppModule(appCodec, app.bankKeeper, app.accountKeeper),
		capability.NewAppModule(appCodec, *app.capabilityKeeper),
		crisis.NewAppModule(&app.crisisKeeper, skipGenesisInvariants),
		gov.NewAppModule(appCodec, app.govKeeper, app.accountKeeper, app.bankKeeper),
		mint.NewAppModule(appCodec, app.mintKeeper),
		slashing.NewAppModule(appCodec, app.slashingKeeper, app.accountKeeper, app.bankKeeper, app.stakingKeeper),
		distr.NewAppModule(appCodec, app.distrKeeper, app.accountKeeper, app.bankKeeper, app.stakingKeeper),
		staking.NewAppModule(appCodec, app.stakingKeeper, app.accountKeeper, app.bankKeeper),
		upgrade.NewAppModule(app.upgradeKeeper),
		evidence.NewAppModule(app.evidenceKeeper),
		ibc.NewAppModule(app.ibcKeeper),
		params.NewAppModule(app.paramsKeeper),
		transferModule,
		guardian.NewAppModule(appCodec, app.guardianKeeper),
		token.NewAppModule(appCodec, app.tokenKeeper, app.accountKeeper, app.bankKeeper),
		record.NewAppModule(appCodec, app.recordKeeper, app.accountKeeper, app.bankKeeper),
		nft.NewAppModule(appCodec, app.nftKeeper, app.accountKeeper, app.bankKeeper),
		htlc.NewAppModule(appCodec, app.htlcKeeper, app.accountKeeper, app.bankKeeper),
		coinswap.NewAppModule(appCodec, app.coinswapKeeper, app.accountKeeper, app.bankKeeper),
		service.NewAppModule(appCodec, app.serviceKeeper, app.accountKeeper, app.bankKeeper),
		oracle.NewAppModule(appCodec, app.oracleKeeper),
		random.NewAppModule(appCodec, app.randomKeeper, app.accountKeeper, app.bankKeeper),
		legacy.NewAppModule(appCodec, app.bankKeeper),
		wasmModule,
	)

	// During begin block slashing happens after distr.BeginBlocker so that
	// there is nothing left over in the validator fee pool, so as to keep the
	// CanWithdrawInvariant invariant.
	// NOTE: staking module is required if HistoricalEntries param > 0
	app.mm.SetOrderBeginBlockers(
		upgradetypes.ModuleName, minttypes.ModuleName, distrtypes.ModuleName,
		slashingtypes.ModuleName, evidencetypes.ModuleName, stakingtypes.ModuleName,
		ibchost.ModuleName, htlctypes.ModuleName, randomtypes.ModuleName, wasm.ModuleName,
	)
	app.mm.SetOrderEndBlockers(
		crisistypes.ModuleName, govtypes.ModuleName, stakingtypes.ModuleName,
		servicetypes.ModuleName, wasm.ModuleName,
	)

	// NOTE: The genutils module must occur after staking so that pools are
	// properly initialized with tokens from genesis accounts.
	// NOTE: Capability module must occur first so that it can initialize any capabilities
	// so that other modules that want to create or claim capabilities afterwards in InitChain
	// can do so safely.
	app.mm.SetOrderInitGenesis(
		capabilitytypes.ModuleName, authtypes.ModuleName, banktypes.ModuleName, distrtypes.ModuleName, stakingtypes.ModuleName,
		slashingtypes.ModuleName, govtypes.ModuleName, minttypes.ModuleName, crisistypes.ModuleName,
		ibchost.ModuleName, genutiltypes.ModuleName, evidencetypes.ModuleName, ibctransfertypes.ModuleName,
		guardiantypes.ModuleName, tokentypes.ModuleName, nfttypes.ModuleName, htlctypes.ModuleName, recordtypes.ModuleName,
		coinswaptypes.ModuleName, servicetypes.ModuleName, oracletypes.ModuleName, randomtypes.ModuleName, wasm.ModuleName,
	)

	app.mm.RegisterInvariants(&app.crisisKeeper)
	app.mm.RegisterRoutes(app.Router(), app.QueryRouter(), encodingConfig.Amino)
	app.mm.RegisterServices(module.NewConfigurator(app.MsgServiceRouter(), app.GRPCQueryRouter()))

	// create the simulation manager and define the order of the modules for deterministic simulations
	//
	// NOTE: this is not required apps that don't use the simulator for fuzz testing
	// transactions
	app.sm = module.NewSimulationManager(
		auth.NewAppModule(appCodec, app.accountKeeper, authsims.RandomGenesisAccounts),
		bank.NewAppModule(appCodec, app.bankKeeper, app.accountKeeper),
		capability.NewAppModule(appCodec, *app.capabilityKeeper),
		gov.NewAppModule(appCodec, app.govKeeper, app.accountKeeper, app.bankKeeper),
		mint.NewAppModule(appCodec, app.mintKeeper),
		staking.NewAppModule(appCodec, app.stakingKeeper, app.accountKeeper, app.bankKeeper),
		distr.NewAppModule(appCodec, app.distrKeeper, app.accountKeeper, app.bankKeeper, app.stakingKeeper),
		slashing.NewAppModule(appCodec, app.slashingKeeper, app.accountKeeper, app.bankKeeper, app.stakingKeeper),
		params.NewAppModule(app.paramsKeeper),
		evidence.NewAppModule(app.evidenceKeeper),
		ibc.NewAppModule(app.ibcKeeper),
		transferModule,
		guardian.NewAppModule(appCodec, app.guardianKeeper),
		token.NewAppModule(appCodec, app.tokenKeeper, app.accountKeeper, app.bankKeeper),
		record.NewAppModule(appCodec, app.recordKeeper, app.accountKeeper, app.bankKeeper),
		nft.NewAppModule(appCodec, app.nftKeeper, app.accountKeeper, app.bankKeeper),
		htlc.NewAppModule(appCodec, app.htlcKeeper, app.accountKeeper, app.bankKeeper),
		coinswap.NewAppModule(appCodec, app.coinswapKeeper, app.accountKeeper, app.bankKeeper),
		service.NewAppModule(appCodec, app.serviceKeeper, app.accountKeeper, app.bankKeeper),
		oracle.NewAppModule(appCodec, app.oracleKeeper),
		random.NewAppModule(appCodec, app.randomKeeper, app.accountKeeper, app.bankKeeper),
	)

	app.sm.RegisterStoreDecoders()

	// initialize stores
	app.MountKVStores(keys)
	app.MountTransientStores(tkeys)
	app.MountMemoryStores(memKeys)

	// initialize BaseApp
	app.SetInitChainer(app.InitChainer)
	app.SetBeginBlocker(app.BeginBlocker)
	app.SetAnteHandler(NewAnteHandler(
		app.accountKeeper,
		app.bankKeeper,
		app.tokenKeeper,
		app.oracleKeeper,
		app.guardianKeeper,
		ante.DefaultSigVerificationGasConsumer,
		encodingConfig.TxConfig.SignModeHandler(),
	))
	app.SetEndBlocker(app.EndBlocker)

	// Set software upgrade execution logic
	app.RegisterUpgradePlan("v1",
		store.StoreUpgrades{
			Added: []string{wasm.StoreKey},
		},
		func(ctx sdk.Context, plan sdkupgrade.Plan) {
			wasm.InitGenesis(ctx, &app.wasmKeeper, wasm.GenesisState{
				Params: wasm.DefaultParams(),
			},
				app.stakingKeeper,
				wasmModule.Route().Handler(),
			)
		},
	)

	if loadLatest {
		if err := app.LoadLatestVersion(); err != nil {
			tmos.Exit(err.Error())
		}

		// Initialize and seal the capability keeper so all persistent capabilities
		// are loaded in-memory and prevent any further modules from creating scoped
		// sub-keepers.
		// This must be done during creation of baseapp rather than in InitChain so
		// that in-memory capabilities get regenerated on app restart.
		// Note that since this reads from the store, we can only perform it when
		// `loadLatest` is set to true.
		ctx := app.BaseApp.NewUncachedContext(true, tmproto.Header{})
		app.capabilityKeeper.InitializeAndSeal(ctx)
	}

	app.ScopedIBCKeeper = scopedIBCKeeper
	app.ScopedTransferKeeper = scopedTransferKeeper

	return app
}

// MakeCodecs constructs the *std.Codec and *codec.LegacyAmino instances used by
// irisapp. It is useful for tests and clients who do not want to construct the
// full irisapp
func MakeCodecs() (codec.Marshaler, *codec.LegacyAmino) {
	config := MakeEncodingConfig()
	return config.Marshaler, config.Amino
}

// Name returns the name of the App
func (app *IrisApp) Name() string { return app.BaseApp.Name() }

// BeginBlocker application updates every begin block
func (app *IrisApp) BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
	return app.mm.BeginBlock(ctx, req)
}

// EndBlocker application updates every end block
func (app *IrisApp) EndBlocker(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
	return app.mm.EndBlock(ctx, req)
}

// InitChainer application update at chain initialization
func (app *IrisApp) InitChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	var genesisState GenesisState
	app.cdc.MustUnmarshalJSON(req.AppStateBytes, &genesisState)

	// add system service at InitChainer, overwrite if it exists
	var serviceGenState servicetypes.GenesisState
	app.appCodec.MustUnmarshalJSON(genesisState[servicetypes.ModuleName], &serviceGenState)
	serviceGenState.Definitions = append(serviceGenState.Definitions, servicetypes.GenOraclePriceSvcDefinition())
	serviceGenState.Bindings = append(serviceGenState.Bindings, servicetypes.GenOraclePriceSvcBinding(nativeToken.MinUnit))
	serviceGenState.Definitions = append(serviceGenState.Definitions, randomtypes.GetSvcDefinition())
	genesisState[servicetypes.ModuleName] = app.appCodec.MustMarshalJSON(&serviceGenState)

	return app.mm.InitGenesis(ctx, app.appCodec, genesisState)
}

// LoadHeight loads a particular height
func (app *IrisApp) LoadHeight(height int64) error {
	return app.LoadVersion(height)
}

// ModuleAccountAddrs returns all the app's module account addresses.
func (app *IrisApp) ModuleAccountAddrs() map[string]bool {
	modAccAddrs := make(map[string]bool)
	for acc := range maccPerms {
		modAccAddrs[authtypes.NewModuleAddress(acc).String()] = true
	}

	return modAccAddrs
}

// BlockedAddrs returns all the app's module account addresses that are not
// allowed to receive external tokens.
func (app *IrisApp) BlockedAddrs() map[string]bool {
	blockedAddrs := make(map[string]bool)
	for acc := range maccPerms {
		blockedAddrs[authtypes.NewModuleAddress(acc).String()] = !allowedReceivingModAcc[acc]
	}

	return blockedAddrs
}

// LegacyAmino returns IrisApp's amino codec.
//
// NOTE: This is solely to be used for testing purposes as it may be desirable
// for modules to register their own custom testing types.
func (app *IrisApp) LegacyAmino() *codec.LegacyAmino {
	return app.cdc
}

// AppCodec returns IrisApp's app codec.
//
// NOTE: This is solely to be used for testing purposes as it may be desirable
// for modules to register their own custom testing types.
func (app *IrisApp) AppCodec() codec.Marshaler {
	return app.appCodec
}

// InterfaceRegistry returns IrisApp's InterfaceRegistry
func (app *IrisApp) InterfaceRegistry() types.InterfaceRegistry {
	return app.interfaceRegistry
}

// GetKey returns the KVStoreKey for the provided store key.
//
// NOTE: This is solely to be used for testing purposes.
func (app *IrisApp) GetKey(storeKey string) *sdk.KVStoreKey {
	return app.keys[storeKey]
}

// GetTKey returns the TransientStoreKey for the provided store key.
//
// NOTE: This is solely to be used for testing purposes.
func (app *IrisApp) GetTKey(storeKey string) *sdk.TransientStoreKey {
	return app.tkeys[storeKey]
}

// GetMemKey returns the MemStoreKey for the provided mem key.
//
// NOTE: This is solely used for testing purposes.
func (app *IrisApp) GetMemKey(storeKey string) *sdk.MemoryStoreKey {
	return app.memKeys[storeKey]
}

// GetSubspace returns a param subspace for a given module name.
//
// NOTE: This is solely to be used for testing purposes.
func (app *IrisApp) GetSubspace(moduleName string) paramstypes.Subspace {
	subspace, _ := app.paramsKeeper.GetSubspace(moduleName)
	return subspace
}

// SimulationManager implements the SimulationApp interface
func (app *IrisApp) SimulationManager() *module.SimulationManager {
	return app.sm
}

// RegisterAPIRoutes registers all application module routes with the provided
// API server.
func (app *IrisApp) RegisterAPIRoutes(apiSvr *api.Server, apiConfig config.APIConfig) {
	clientCtx := apiSvr.ClientCtx
	rpc.RegisterRoutes(clientCtx, apiSvr.Router)
	// Register legacy tx routes.
	authrest.RegisterTxRoutes(clientCtx, apiSvr.Router)

	// Register new tendermint queries routes from grpc-gateway.
	tmservice.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)

	// Register new tx routes from grpc-gateway.
	authtx.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)
	// Register legacy and grpc-gateway routes for all modules.
	ModuleBasics.RegisterRESTRoutes(clientCtx, apiSvr.Router)
	ModuleBasics.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)

	if apiConfig.Swagger {
		lite.RegisterSwaggerAPI(clientCtx, apiSvr.Router)
	}
}

// RegisterTxService implements the Application.RegisterTxService method.
func (app *IrisApp) RegisterTxService(clientCtx client.Context) {
	authtx.RegisterTxService(app.BaseApp.GRPCQueryRouter(), clientCtx, app.BaseApp.Simulate, app.interfaceRegistry)
}

// RegisterTendermintService implements the Application.RegisterTendermintService method.
func (app *IrisApp) RegisterTendermintService(clientCtx client.Context) {
	tmservice.RegisterTendermintService(app.BaseApp.GRPCQueryRouter(), clientCtx, app.interfaceRegistry)
}

// RegisterUpgradePlan implements the upgrade execution logic of the upgrade module
func (app *IrisApp) RegisterUpgradePlan(planName string,
	upgrades store.StoreUpgrades, upgradeHandler sdkupgrade.UpgradeHandler) {
	upgradeInfo, err := app.upgradeKeeper.ReadUpgradeInfoFromDisk()
	if err != nil {
		app.Logger().Info("not found upgrade plan", "planName", planName, "err", err.Error())
		return
	}

	if upgradeInfo.Name == planName && !app.upgradeKeeper.IsSkipHeight(upgradeInfo.Height) {
		// this configures a no-op upgrade handler for the planName upgrade
		app.upgradeKeeper.SetUpgradeHandler(planName, upgradeHandler)
		// configure store loader that checks if version+1 == upgradeHeight and applies store upgrades
		app.SetStoreLoader(sdkupgrade.UpgradeStoreLoader(upgradeInfo.Height, &upgrades))
	}
}

// GetMaccPerms returns a copy of the module account permissions
func GetMaccPerms() map[string][]string {
	dupMaccPerms := make(map[string][]string)
	for k, v := range maccPerms {
		dupMaccPerms[k] = v
	}
	return dupMaccPerms
}

// initParamsKeeper init params keeper and its subspaces
func initParamsKeeper(appCodec codec.BinaryMarshaler, legacyAmino *codec.LegacyAmino, key, tkey sdk.StoreKey) paramskeeper.Keeper {
	paramsKeeper := paramskeeper.NewKeeper(appCodec, legacyAmino, key, tkey)

	paramsKeeper.Subspace(authtypes.ModuleName)
	paramsKeeper.Subspace(banktypes.ModuleName)
	paramsKeeper.Subspace(stakingtypes.ModuleName)
	paramsKeeper.Subspace(minttypes.ModuleName)
	paramsKeeper.Subspace(distrtypes.ModuleName)
	paramsKeeper.Subspace(slashingtypes.ModuleName)
	paramsKeeper.Subspace(govtypes.ModuleName).WithKeyTable(govtypes.ParamKeyTable())
	paramsKeeper.Subspace(crisistypes.ModuleName)
	paramsKeeper.Subspace(ibctransfertypes.ModuleName)
	paramsKeeper.Subspace(tokentypes.ModuleName)
	paramsKeeper.Subspace(recordtypes.ModuleName)
	paramsKeeper.Subspace(htlctypes.ModuleName)
	paramsKeeper.Subspace(coinswaptypes.ModuleName)
	paramsKeeper.Subspace(servicetypes.ModuleName)
	paramsKeeper.Subspace(ibchost.ModuleName)
	paramsKeeper.Subspace(wasm.ModuleName)

	return paramsKeeper
}
