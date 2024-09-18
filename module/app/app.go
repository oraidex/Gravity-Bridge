package app

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"cosmossdk.io/log"
	abci "github.com/cometbft/cometbft/abci/types"
	tmjson "github.com/cometbft/cometbft/libs/json"
	tmos "github.com/cometbft/cometbft/libs/os"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/server"
	ibcconnectiontypes "github.com/cosmos/ibc-go/v8/modules/core/03-connection/types"
	"github.com/gorilla/mux"
	"github.com/rakyll/statik/fs"
	"github.com/spf13/cast"

	// Cosmos SDK
	"cosmossdk.io/simapp"
	storetypes "cosmossdk.io/store/types"
	"cosmossdk.io/x/evidence"
	evidencekeeper "cosmossdk.io/x/evidence/keeper"
	evidencetypes "cosmossdk.io/x/evidence/types"
	"cosmossdk.io/x/upgrade"
	upgradekeeper "cosmossdk.io/x/upgrade/keeper"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/grpc/cmtservice"
	nodeservice "github.com/cosmos/cosmos-sdk/client/grpc/node"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	ccodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/server/api"
	"github.com/cosmos/cosmos-sdk/server/config"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/auth"
	sdkante "github.com/cosmos/cosmos-sdk/x/auth/ante"
	authcodec "github.com/cosmos/cosmos-sdk/x/auth/codec"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authsims "github.com/cosmos/cosmos-sdk/x/auth/simulation"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/auth/vesting"
	vestingtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	authzkeeper "github.com/cosmos/cosmos-sdk/x/authz/keeper"
	authzmodule "github.com/cosmos/cosmos-sdk/x/authz/module"
	"github.com/cosmos/cosmos-sdk/x/bank"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	consensusparamkeeper "github.com/cosmos/cosmos-sdk/x/consensus/keeper"
	consensusparamtypes "github.com/cosmos/cosmos-sdk/x/consensus/types"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	crisiskeeper "github.com/cosmos/cosmos-sdk/x/crisis/keeper"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/cosmos/cosmos-sdk/x/gov"
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	"github.com/cosmos/cosmos-sdk/x/mint"
	mintkeeper "github.com/cosmos/cosmos-sdk/x/mint/keeper"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	paramsclient "github.com/cosmos/cosmos-sdk/x/params/client"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	paramsproposal "github.com/cosmos/cosmos-sdk/x/params/types/proposal"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	slashingkeeper "github.com/cosmos/cosmos-sdk/x/slashing/keeper"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cosmos/ibc-go/modules/capability"
	capabilitykeeper "github.com/cosmos/ibc-go/modules/capability/keeper"
	capabilitytypes "github.com/cosmos/ibc-go/modules/capability/types"

	// Cosmos IBC-Go
	ica "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts"
	icahost "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/host"
	icahostkeeper "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/host/keeper"
	icahosttypes "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/host/types"
	icatypes "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/types"
	transfer "github.com/cosmos/ibc-go/v8/modules/apps/transfer"
	ibctransferkeeper "github.com/cosmos/ibc-go/v8/modules/apps/transfer/keeper"
	ibctransfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	ibc "github.com/cosmos/ibc-go/v8/modules/core"
	ibcclienttypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types"
	porttypes "github.com/cosmos/ibc-go/v8/modules/core/05-port/types"
	ibcexported "github.com/cosmos/ibc-go/v8/modules/core/exported"
	ibckeeper "github.com/cosmos/ibc-go/v8/modules/core/keeper"

	// Osmosis-Labs Bech32-IBC
	"github.com/Gravity-Bridge/Gravity-Bridge/module/x/bech32ibc"
	bech32ibckeeper "github.com/Gravity-Bridge/Gravity-Bridge/module/x/bech32ibc/keeper"
	bech32ibctypes "github.com/Gravity-Bridge/Gravity-Bridge/module/x/bech32ibc/types"

	// unnamed import of statik for swagger UI support
	// _ "github.com/cosmos/cosmos-sdk/client/docs/statik"
	_ "github.com/Gravity-Bridge/Gravity-Bridge/module/docs/statik"

	// Tharsis Ethermint
	ethante "github.com/evmos/ethermint/app/ante"

	"github.com/Gravity-Bridge/Gravity-Bridge/module/app/ante"
	gravityparams "github.com/Gravity-Bridge/Gravity-Bridge/module/app/params"
	"github.com/Gravity-Bridge/Gravity-Bridge/module/x/gravity"
	"github.com/Gravity-Bridge/Gravity-Bridge/module/x/gravity/exported"
	"github.com/Gravity-Bridge/Gravity-Bridge/module/x/gravity/keeper"
	gravitytypes "github.com/Gravity-Bridge/Gravity-Bridge/module/x/gravity/types"
)

const appName = "app"

var (
	// DefaultNodeHome sets the folder where the applcation data and configuration will be stored
	DefaultNodeHome string

	// ModuleBasics The module BasicManager is in charge of setting up basic,
	// non-dependant module elements, such as codec registration
	// and genesis verification.

	ModuleBasics = module.NewBasicManager(
		auth.AppModuleBasic{},
		authzmodule.AppModuleBasic{},
		genutil.AppModuleBasic{},
		bank.AppModuleBasic{},
		capability.AppModuleBasic{},
		staking.AppModuleBasic{},
		mint.AppModuleBasic{},
		distr.AppModuleBasic{},
		gov.NewAppModuleBasic(
			[]govclient.ProposalHandler{
				paramsclient.ProposalHandler,
			},
		),
		params.AppModuleBasic{},
		crisis.AppModuleBasic{},
		slashing.AppModuleBasic{},
		ibc.AppModuleBasic{},
		upgrade.AppModuleBasic{},
		evidence.AppModuleBasic{},
		transfer.AppModuleBasic{},
		vesting.AppModuleBasic{},
		gravity.AppModuleBasic{},
		bech32ibc.AppModuleBasic{},
		ica.AppModuleBasic{},
	)

	// module account permissions
	// NOTE: We believe that this is giving various modules access to functions of the supply module? We will probably need to use this.
	maccPerms = map[string][]string{
		authtypes.FeeCollectorName:     nil,
		distrtypes.ModuleName:          nil,
		minttypes.ModuleName:           {authtypes.Minter},
		stakingtypes.BondedPoolName:    {authtypes.Burner, authtypes.Staking},
		stakingtypes.NotBondedPoolName: {authtypes.Burner, authtypes.Staking},
		govtypes.ModuleName:            {authtypes.Burner},
		ibctransfertypes.ModuleName:    {authtypes.Minter, authtypes.Burner},
		gravitytypes.ModuleName:        {authtypes.Minter, authtypes.Burner},
		icatypes.ModuleName:            nil,
	}

	// module accounts that are allowed to receive tokens
	allowedReceivingModAcc = map[string]bool{
		distrtypes.ModuleName: true,
	}

	// verify app interface at compile time
	_ runtime.AppI            = (*Gravity)(nil)
	_ servertypes.Application = (*Gravity)(nil)

	// enable checks that run on the first BeginBlocker execution after an upgrade/genesis init/node restart
	firstBlock sync.Once
)

// MakeCodec creates the application codec. The codec is sealed before it is
// returned.
func MakeCodec() *codec.LegacyAmino {
	var cdc = codec.NewLegacyAmino()
	ModuleBasics.RegisterLegacyAminoCodec(cdc)
	vesting.AppModuleBasic{}.RegisterLegacyAminoCodec(cdc)
	sdk.RegisterLegacyAminoCodec(cdc)
	ccodec.RegisterCrypto(cdc)
	cdc.Seal()
	return cdc
}

// Gravity extended ABCI application
type Gravity struct {
	*baseapp.BaseApp
	legacyAmino       *codec.LegacyAmino
	appCodec          codec.Codec
	interfaceRegistry types.InterfaceRegistry

	invCheckPeriod uint

	// keys to access the substores
	keys    map[string]*storetypes.KVStoreKey
	tKeys   map[string]*storetypes.TransientStoreKey
	memKeys map[string]*storetypes.MemoryStoreKey

	// keepers
	// NOTE: If you add anything to this struct, add a nil check to ValidateMembers below!
	accountKeeper         *authkeeper.AccountKeeper
	authzKeeper           *authzkeeper.Keeper
	bankKeeper            *bankkeeper.BaseKeeper
	capabilityKeeper      *capabilitykeeper.Keeper
	stakingKeeper         *stakingkeeper.Keeper
	slashingKeeper        *slashingkeeper.Keeper
	mintKeeper            *mintkeeper.Keeper
	distrKeeper           *distrkeeper.Keeper
	govKeeper             *govkeeper.Keeper
	crisisKeeper          *crisiskeeper.Keeper
	upgradeKeeper         *upgradekeeper.Keeper
	paramsKeeper          *paramskeeper.Keeper
	ibcKeeper             *ibckeeper.Keeper
	evidenceKeeper        *evidencekeeper.Keeper
	ibcTransferKeeper     *ibctransferkeeper.Keeper
	gravityKeeper         *keeper.Keeper
	bech32IbcKeeper       *bech32ibckeeper.Keeper
	icaHostKeeper         *icahostkeeper.Keeper
	consensusParamsKeeper consensusparamkeeper.Keeper

	// make scoped keepers public for test purposes
	// NOTE: If you add anything to this struct, add a nil check to ValidateMembers below!
	ScopedIBCKeeper      *capabilitykeeper.ScopedKeeper
	ScopedTransferKeeper *capabilitykeeper.ScopedKeeper
	ScopedIcaHostKeeper  *capabilitykeeper.ScopedKeeper

	// Module Manager
	mm *module.Manager

	// simulation manager
	sm *module.SimulationManager

	// configurator
	configurator *module.Configurator
}

// ValidateMembers checks for nil members
func (app Gravity) ValidateMembers() {
	if app.legacyAmino == nil {
		panic("Nil legacyAmino!")
	}

	// keepers
	if app.accountKeeper == nil {
		panic("Nil accountKeeper!")
	}
	if app.authzKeeper == nil {
		panic("Nil authzKeeper!")
	}
	if app.bankKeeper == nil {
		panic("Nil bankKeeper!")
	}
	if app.capabilityKeeper == nil {
		panic("Nil capabilityKeeper!")
	}
	if app.stakingKeeper == nil {
		panic("Nil stakingKeeper!")
	}
	if app.slashingKeeper == nil {
		panic("Nil slashingKeeper!")
	}
	if app.mintKeeper == nil {
		panic("Nil mintKeeper!")
	}
	if app.distrKeeper == nil {
		panic("Nil distrKeeper!")
	}
	if app.govKeeper == nil {
		panic("Nil govKeeper!")
	}
	if app.crisisKeeper == nil {
		panic("Nil crisisKeeper!")
	}
	if app.upgradeKeeper == nil {
		panic("Nil upgradeKeeper!")
	}
	if app.paramsKeeper == nil {
		panic("Nil paramsKeeper!")
	}
	if app.ibcKeeper == nil {
		panic("Nil ibcKeeper!")
	}
	if app.evidenceKeeper == nil {
		panic("Nil evidenceKeeper!")
	}
	if app.ibcTransferKeeper == nil {
		panic("Nil ibcTransferKeeper!")
	}
	if app.gravityKeeper == nil {
		panic("Nil gravityKeeper!")
	}
	if app.bech32IbcKeeper == nil {
		panic("Nil bech32IbcKeeper!")
	}
	if app.icaHostKeeper == nil {
		panic("Nil icaHostKeeper!")
	}

	// scoped keepers
	if app.ScopedIBCKeeper == nil {
		panic("Nil ScopedIBCKeeper!")
	}
	if app.ScopedTransferKeeper == nil {
		panic("Nil ScopedTransferKeeper!")
	}

	// managers
	if app.mm == nil {
		panic("Nil ModuleManager!")
	}
	if app.sm == nil {
		panic("Nil ModuleManager!")
	}
}

func init() {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	DefaultNodeHome = filepath.Join(userHomeDir, ".gravity")
}

func NewGravityApp(
	logger log.Logger, db dbm.DB, traceStore io.Writer, loadLatest bool, skipUpgradeHeights map[int64]bool,
	homePath string, invCheckPeriod uint, encodingConfig gravityparams.EncodingConfig,
	appOpts servertypes.AppOptions, baseAppOptions ...func(*baseapp.BaseApp),
) *Gravity {
	appCodec := encodingConfig.Marshaler
	legacyAmino := encodingConfig.Amino
	interfaceRegistry := encodingConfig.InterfaceRegistry

	bApp := *baseapp.NewBaseApp(appName, logger, db, encodingConfig.TxConfig.TxDecoder(), baseAppOptions...)
	bApp.SetCommitMultiStoreTracer(traceStore)
	bApp.SetVersion(version.Version)
	bApp.SetInterfaceRegistry(interfaceRegistry)

	keys := storetypes.NewKVStoreKeys(
		authtypes.StoreKey, authzkeeper.StoreKey, banktypes.StoreKey,
		stakingtypes.StoreKey, minttypes.StoreKey, distrtypes.StoreKey,
		slashingtypes.StoreKey, govtypes.StoreKey, paramstypes.StoreKey,
		ibcexported.StoreKey, upgradetypes.StoreKey, evidencetypes.StoreKey,
		ibctransfertypes.StoreKey, capabilitytypes.StoreKey,
		gravitytypes.StoreKey, bech32ibctypes.StoreKey,
		icahosttypes.StoreKey, consensusparamtypes.StoreKey,
	)
	tKeys := storetypes.NewTransientStoreKeys(paramstypes.TStoreKey)
	memKeys := storetypes.NewMemoryStoreKeys(capabilitytypes.MemStoreKey)

	// nolint: exhaustruct
	var app = &Gravity{
		BaseApp:           &bApp,
		legacyAmino:       legacyAmino,
		appCodec:          appCodec,
		interfaceRegistry: interfaceRegistry,
		invCheckPeriod:    invCheckPeriod,
		keys:              keys,
		tKeys:             tKeys,
		memKeys:           memKeys,
	}

	Authority := authtypes.NewModuleAddress(govtypes.ModuleName)
	AuthorityAddr := Authority.String()

	paramsKeeper := initParamsKeeper(appCodec, legacyAmino, keys[paramstypes.StoreKey], tKeys[paramstypes.TStoreKey])
	// for the sake of param migration
	paramsKeeper.Subspace(exported.DefaultParamspace)
	app.paramsKeeper = &paramsKeeper

	// set the BaseApp's parameter store
	app.consensusParamsKeeper = consensusparamkeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[consensusparamtypes.StoreKey]),
		AuthorityAddr,
		runtime.EventService{},
	)

	bApp.SetParamStore(app.consensusParamsKeeper.ParamsStore)

	capabilityKeeper := *capabilitykeeper.NewKeeper(
		appCodec,
		keys[capabilitytypes.StoreKey],
		memKeys[capabilitytypes.MemStoreKey],
	)
	app.capabilityKeeper = &capabilityKeeper

	scopedIBCKeeper := capabilityKeeper.ScopeToModule(ibcexported.ModuleName)
	app.ScopedIBCKeeper = &scopedIBCKeeper

	scopedTransferKeeper := capabilityKeeper.ScopeToModule(ibctransfertypes.ModuleName)
	app.ScopedTransferKeeper = &scopedTransferKeeper

	scopedIcaHostKeeper := capabilityKeeper.ScopeToModule(icahosttypes.SubModuleName)
	app.ScopedIcaHostKeeper = &scopedIcaHostKeeper

	// Applications that wish to enforce statically created ScopedKeepers should call `Seal` after creating
	// their scoped modules in `NewApp` with `ScopeToModule`
	capabilityKeeper.Seal()

	accountKeeper := authkeeper.NewAccountKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[authtypes.StoreKey]),
		authtypes.ProtoBaseAccount,
		maccPerms,
		authcodec.NewBech32Codec(sdk.GetConfig().GetBech32AccountAddrPrefix()),
		sdk.GetConfig().GetBech32AccountAddrPrefix(),
		AuthorityAddr,
	)
	app.accountKeeper = &accountKeeper

	authzKeeper := authzkeeper.NewKeeper(
		runtime.NewKVStoreService(keys[authzkeeper.StoreKey]),
		appCodec,
		bApp.MsgServiceRouter(),
		app.accountKeeper,
	)
	app.authzKeeper = &authzKeeper

	bankKeeper := bankkeeper.NewBaseKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[banktypes.StoreKey]),
		accountKeeper,
		app.BlockedAddrs(),
		AuthorityAddr,
		logger,
	)
	app.bankKeeper = &bankKeeper

	stakingKeeper := stakingkeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[stakingtypes.StoreKey]),
		accountKeeper,
		bankKeeper,
		AuthorityAddr,
		authcodec.NewBech32Codec(sdk.GetConfig().GetBech32ValidatorAddrPrefix()),
		authcodec.NewBech32Codec(sdk.GetConfig().GetBech32ConsensusAddrPrefix()),
	)
	app.stakingKeeper = stakingKeeper

	distrKeeper := distrkeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[distrtypes.StoreKey]),
		accountKeeper,
		bankKeeper,
		stakingKeeper,
		authtypes.FeeCollectorName,
		AuthorityAddr,
	)
	app.distrKeeper = &distrKeeper

	slashingKeeper := slashingkeeper.NewKeeper(
		appCodec,
		legacyAmino,
		runtime.NewKVStoreService(keys[slashingtypes.StoreKey]),
		stakingKeeper,
		AuthorityAddr,
	)
	app.slashingKeeper = &slashingKeeper

	upgradeKeeper := upgradekeeper.NewKeeper(
		skipUpgradeHeights,
		runtime.NewKVStoreService(keys[upgradetypes.StoreKey]),
		appCodec,
		homePath,
		&bApp,
		AuthorityAddr,
	)
	app.upgradeKeeper = upgradeKeeper

	ibcKeeper := *ibckeeper.NewKeeper(
		appCodec,
		keys[ibcexported.StoreKey],
		app.GetSubspace(ibcexported.ModuleName),
		stakingKeeper,
		upgradeKeeper,
		scopedIBCKeeper,
		AuthorityAddr,
	)
	app.ibcKeeper = &ibcKeeper

	ibcTransferKeeper := ibctransferkeeper.NewKeeper(
		appCodec, keys[ibctransfertypes.StoreKey], app.GetSubspace(ibctransfertypes.ModuleName),
		ibcKeeper.ChannelKeeper, ibcKeeper.ChannelKeeper, ibcKeeper.PortKeeper,
		accountKeeper, bankKeeper, scopedTransferKeeper, AuthorityAddr,
	)
	app.ibcTransferKeeper = &ibcTransferKeeper

	bech32IbcKeeper := *bech32ibckeeper.NewKeeper(
		ibcKeeper.ChannelKeeper, appCodec, keys[bech32ibctypes.StoreKey],
		ibcTransferKeeper,
	)
	app.bech32IbcKeeper = &bech32IbcKeeper

	icaHostKeeper := icahostkeeper.NewKeeper(
		appCodec, keys[icahosttypes.StoreKey], app.GetSubspace(icahosttypes.SubModuleName),
		nil,
		ibcKeeper.ChannelKeeper, ibcKeeper.PortKeeper,
		accountKeeper, scopedIcaHostKeeper, app.MsgServiceRouter(), AuthorityAddr,
	)
	app.icaHostKeeper = &icaHostKeeper

	gravityKeeper := keeper.NewKeeper(
		keys[gravitytypes.StoreKey],
		app.GetSubspace(gravitytypes.DefaultParamspace),
		appCodec,
		&bankKeeper,
		stakingKeeper,
		&slashingKeeper,
		&distrKeeper,
		&accountKeeper,
		&ibcTransferKeeper,
		&bech32IbcKeeper,
		ibcKeeper.ChannelKeeper, // interface default is reference
	)
	app.gravityKeeper = &gravityKeeper

	// Add the staking hooks from distribution, slashing, and gravity to staking
	stakingKeeper.SetHooks(
		stakingtypes.NewMultiStakingHooks(
			distrKeeper.Hooks(),
			slashingKeeper.Hooks(),
			gravityKeeper.Hooks(),
		),
	)

	mintKeeper := mintkeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[minttypes.StoreKey]),
		stakingKeeper,
		accountKeeper,
		bankKeeper,
		authtypes.FeeCollectorName,
		AuthorityAddr,
	)
	app.mintKeeper = &mintKeeper

	crisisKeeper := crisiskeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[crisistypes.StoreKey]),
		invCheckPeriod,
		bankKeeper,
		authtypes.FeeCollectorName,
		AuthorityAddr,
		app.accountKeeper.AddressCodec(),
	)
	app.crisisKeeper = crisisKeeper

	govRouter := govv1beta1.NewRouter()
	govRouter.AddRoute(govtypes.RouterKey, govv1beta1.ProposalHandler).
		AddRoute(paramsproposal.RouterKey, params.NewParamChangeProposalHandler(paramsKeeper)).
		AddRoute(gravitytypes.RouterKey, keeper.NewGravityProposalHandler(gravityKeeper)).
		AddRoute(bech32ibctypes.RouterKey, bech32ibc.NewBech32IBCProposalHandler(*app.bech32IbcKeeper))

	govConfig := govtypes.DefaultConfig()
	govKeeper := govkeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[govtypes.StoreKey]),
		accountKeeper,
		bankKeeper,
		stakingKeeper,
		distrKeeper,
		app.MsgServiceRouter(),
		govConfig,
		AuthorityAddr,
	)
	app.govKeeper = govKeeper

	ibcTransferAppModule := transfer.NewAppModule(ibcTransferKeeper)

	// create IBC module from top to bottom of stack
	var ibcTransferIBCModule porttypes.IBCModule
	ibcTransferIBCModule = transfer.NewIBCModule(ibcTransferKeeper)
	ibcTransferIBCModule = gravity.NewIBCMiddleware(ibcTransferIBCModule, *app.gravityKeeper)
	icaAppModule := ica.NewAppModule(nil, &icaHostKeeper)
	icaHostIBCModule := icahost.NewIBCModule(icaHostKeeper)

	ibcRouter := porttypes.NewRouter()
	ibcRouter.AddRoute(ibctransfertypes.ModuleName, ibcTransferIBCModule).
		AddRoute(icahosttypes.SubModuleName, icaHostIBCModule)
	ibcKeeper.SetRouter(ibcRouter)

	evidenceKeeper := evidencekeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[evidencetypes.StoreKey]),
		stakingKeeper,
		slashingKeeper,
		app.accountKeeper.AddressCodec(),
		runtime.ProvideCometInfoService(),
	)
	app.evidenceKeeper = evidenceKeeper

	var skipGenesisInvariants = cast.ToBool(appOpts.Get(crisis.FlagSkipGenesisInvariants))

	// app.registerStoreLoaders()

	mm := module.NewManager(
		genutil.NewAppModule(
			accountKeeper,
			stakingKeeper,
			app,
			encodingConfig.TxConfig,
		),
		auth.NewAppModule(
			appCodec,
			accountKeeper,
			nil,
			app.GetSubspace(authtypes.ModuleName),
		),
		authzmodule.NewAppModule(
			appCodec,
			authzKeeper,
			accountKeeper,
			bankKeeper,
			app.InterfaceRegistry(),
		),
		vesting.NewAppModule(
			accountKeeper,
			bankKeeper,
		),
		bank.NewAppModule(
			appCodec,
			bankKeeper,
			accountKeeper,
			app.GetSubspace(banktypes.ModuleName),
		),
		capability.NewAppModule(
			appCodec,
			capabilityKeeper,
			false,
		),
		crisis.NewAppModule(
			crisisKeeper,
			skipGenesisInvariants,
			app.GetSubspace(crisistypes.ModuleName),
		),
		gov.NewAppModule(
			appCodec,
			govKeeper,
			accountKeeper,
			bankKeeper,
			app.GetSubspace(govtypes.ModuleName),
		),
		mint.NewAppModule(
			appCodec,
			mintKeeper,
			accountKeeper,
			nil,
			app.GetSubspace(minttypes.ModuleName),
		),
		slashing.NewAppModule(
			appCodec,
			slashingKeeper,
			accountKeeper,
			bankKeeper,
			stakingKeeper,
			app.GetSubspace(slashingtypes.ModuleName),
			app.interfaceRegistry,
		),
		distr.NewAppModule(
			appCodec,
			distrKeeper,
			accountKeeper,
			bankKeeper,
			stakingKeeper,
			app.GetSubspace(distrtypes.ModuleName),
		),
		staking.NewAppModule(appCodec,
			stakingKeeper,
			accountKeeper,
			bankKeeper,
			app.GetSubspace(stakingtypes.ModuleName),
		),
		upgrade.NewAppModule(upgradeKeeper, accountKeeper.AddressCodec()),
		evidence.NewAppModule(*evidenceKeeper),
		ibc.NewAppModule(&ibcKeeper),
		params.NewAppModule(paramsKeeper),
		ibcTransferAppModule,
		gravity.NewAppModule(
			gravityKeeper,
			bankKeeper,
			app.GetSubspace(exported.DefaultParamspace),
		),
		bech32ibc.NewAppModule(
			appCodec,
			bech32IbcKeeper,
		),
		icaAppModule,
	)
	app.mm = mm

	// NOTE: capability module's BeginBlocker must come before any modules using capabilities (e.g. IBC)
	mm.SetOrderBeginBlockers(
		upgradetypes.ModuleName,
		capabilitytypes.ModuleName,
		minttypes.ModuleName,
		distrtypes.ModuleName,
		slashingtypes.ModuleName,
		evidencetypes.ModuleName,
		stakingtypes.ModuleName,
		ibcexported.ModuleName,
		banktypes.ModuleName,
		crisistypes.ModuleName,
		authtypes.ModuleName,
		vestingtypes.ModuleName,
		ibctransfertypes.ModuleName,
		bech32ibctypes.ModuleName,
		gravitytypes.ModuleName,
		genutiltypes.ModuleName,
		authz.ModuleName,
		govtypes.ModuleName,
		paramstypes.ModuleName,
		icatypes.ModuleName,
	)
	mm.SetOrderEndBlockers(
		crisistypes.ModuleName,
		govtypes.ModuleName,
		stakingtypes.ModuleName,
		icatypes.ModuleName,
		gravitytypes.ModuleName,
		upgradetypes.ModuleName,
		capabilitytypes.ModuleName,
		minttypes.ModuleName,
		distrtypes.ModuleName,
		slashingtypes.ModuleName,
		evidencetypes.ModuleName,
		ibcexported.ModuleName,
		banktypes.ModuleName,
		authtypes.ModuleName,
		vestingtypes.ModuleName,
		ibctransfertypes.ModuleName,
		bech32ibctypes.ModuleName,
		genutiltypes.ModuleName,
		authz.ModuleName,
		paramstypes.ModuleName,
	)
	mm.SetOrderInitGenesis(
		capabilitytypes.ModuleName,
		authtypes.ModuleName,
		banktypes.ModuleName,
		distrtypes.ModuleName,
		stakingtypes.ModuleName,
		slashingtypes.ModuleName,
		govtypes.ModuleName,
		minttypes.ModuleName,
		upgradetypes.ModuleName,
		ibcexported.ModuleName,
		genutiltypes.ModuleName,
		evidencetypes.ModuleName,
		ibctransfertypes.ModuleName,
		authz.ModuleName,
		bech32ibctypes.ModuleName, // Must go before gravity so that pending ibc auto forwards can be restored
		gravitytypes.ModuleName,
		crisistypes.ModuleName,
		vestingtypes.ModuleName,
		paramstypes.ModuleName,
		icatypes.ModuleName,
	)

	mm.RegisterInvariants(crisisKeeper)
	configurator := module.NewConfigurator(appCodec, app.MsgServiceRouter(), app.GRPCQueryRouter())
	app.configurator = &configurator
	mm.RegisterServices(*app.configurator)

	sm := *module.NewSimulationManager(
		auth.NewAppModule(appCodec, accountKeeper, authsims.RandomGenesisAccounts, app.GetSubspace(authtypes.ModuleName)),
		bank.NewAppModule(appCodec, bankKeeper, accountKeeper, app.GetSubspace(banktypes.ModuleName)),
		capability.NewAppModule(appCodec, capabilityKeeper, false),
		gov.NewAppModule(appCodec, govKeeper, accountKeeper, bankKeeper, app.GetSubspace(govtypes.ModuleName)),
		mint.NewAppModule(appCodec, mintKeeper, accountKeeper, nil, app.GetSubspace(minttypes.ModuleName)),
		staking.NewAppModule(appCodec, stakingKeeper, accountKeeper, bankKeeper, app.GetSubspace(stakingtypes.ModuleName)),
		distr.NewAppModule(appCodec, distrKeeper, accountKeeper, bankKeeper, stakingKeeper, app.GetSubspace(distrtypes.ModuleName)),
		slashing.NewAppModule(appCodec, slashingKeeper, accountKeeper, bankKeeper, stakingKeeper, app.GetSubspace(slashingtypes.ModuleName), app.interfaceRegistry),
		params.NewAppModule(paramsKeeper),
		evidence.NewAppModule(*evidenceKeeper),
		ibc.NewAppModule(&ibcKeeper),
		ibcTransferAppModule,
	)
	app.sm = &sm

	sm.RegisterStoreDecoders()

	app.MountKVStores(keys)
	app.MountTransientStores(tKeys)
	app.MountMemoryStores(memKeys)

	app.SetInitChainer(app.InitChainer)
	app.SetBeginBlocker(app.BeginBlocker)

	options := sdkante.HandlerOptions{
		AccountKeeper:   accountKeeper,
		BankKeeper:      bankKeeper,
		FeegrantKeeper:  nil, // Note: If feegrant keeper is added, change the bridge_fees.go antehandler to use it as well
		SignModeHandler: encodingConfig.TxConfig.SignModeHandler(),
		SigGasConsumer:  ethante.DefaultSigVerificationGasConsumer,
	}

	// Note: If feegrant keeper is added, add it to the NewAnteHandler call instead of nil
	ah, err := ante.NewAnteHandler(options, &gravityKeeper, &accountKeeper, &bankKeeper, nil, &ibcKeeper, appCodec)
	if err != nil {
		panic("invalid antehandler created")
	}
	app.SetAnteHandler(*ah)

	app.SetEndBlocker(app.EndBlocker)

	app.registerUpgradeHandlers()

	if loadLatest {
		if err := app.LoadLatestVersion(); err != nil {
			tmos.Exit(err.Error())
		}
	}

	keeper.RegisterProposalTypes()

	// We don't allow anything to be nil
	app.ValidateMembers()
	return app
}

// MakeCodecs constructs the *std.Codec and *codec.LegacyAmino instances used by
// simapp. It is useful for tests and clients who do not want to construct the
// full simapp
func MakeCodecs() (codec.Codec, *codec.LegacyAmino) {
	config := MakeEncodingConfig()
	return config.Marshaler, config.Amino
}

// Name returns the name of the App
func (app *Gravity) Name() string { return app.BaseApp.Name() }

// BeginBlocker application updates every begin block
func (app *Gravity) BeginBlocker(ctx sdk.Context) (sdk.BeginBlock, error) {
	out, err := app.mm.BeginBlock(ctx)
	if err != nil {
		return out, err
	}
	firstBlock.Do(func() { // Run the startup firstBeginBlocker assertions only once
		app.firstBeginBlocker(ctx)
	})

	return out, nil
}

// Perform necessary checks at the start of this node's first BeginBlocker execution
// Note: This should ONLY be called once, it should be called at the top of BeginBlocker guarded by firstBlock
func (app *Gravity) firstBeginBlocker(ctx sdk.Context) {
	app.assertBech32PrefixMatches(ctx)

	// The following call should panic if any invalid ERC20 addresses exist in types/const.go
	monitoredErc20s := app.gravityKeeper.MonitoredERC20Tokens(ctx)
	if len(monitoredErc20s) > 0 {
		ctx.Logger().Info(
			"Loaded Monitored ERC20 Tokens, your orchestrator is required to monitor the Gravity.sol balance of the following tokens: %v",
			monitoredErc20s,
		)
	} else {
		ctx.Logger().Info("Monitored ERC20 Tokens not yet set, your orchestrator is currently not required to monitor any Gravity.sol balances")
	}
}

// EndBlocker application updates every end block
func (app *Gravity) EndBlocker(ctx sdk.Context) (sdk.EndBlock, error) {
	return app.mm.EndBlock(ctx)
}

// InitChainer application update at chain initialization
func (app *Gravity) InitChainer(ctx sdk.Context, req *abci.RequestInitChain) (*abci.ResponseInitChain, error) {
	var genesisState simapp.GenesisState
	if err := tmjson.Unmarshal(req.AppStateBytes, &genesisState); err != nil {
		panic(err)
	}

	app.upgradeKeeper.SetModuleVersionMap(ctx, app.mm.GetVersionMap())

	return app.mm.InitGenesis(ctx, app.appCodec, genesisState)
}

// LoadHeight loads a particular height
func (app *Gravity) LoadHeight(height int64) error {
	return app.LoadVersion(height)
}

// ModuleAccountAddrs returns all the app's module account addresses.
func (app *Gravity) ModuleAccountAddrs() map[string]bool {
	modAccAddrs := make(map[string]bool)
	for acc := range maccPerms {
		modAccAddrs[authtypes.NewModuleAddress(acc).String()] = true
	}

	return modAccAddrs
}

// BlockedAddrs returns all the app's module account addresses that are not
// allowed to receive external tokens.
func (app *Gravity) BlockedAddrs() map[string]bool {
	blockedAddrs := make(map[string]bool)
	for acc := range maccPerms {
		blockedAddrs[authtypes.NewModuleAddress(acc).String()] = !allowedReceivingModAcc[acc]
	}

	return blockedAddrs
}

// LegacyAmino returns SimApp's amino codec.
//
// NOTE: This is solely to be used for testing purposes as it may be desirable
// for modules to register their own custom testing types.
func (app *Gravity) LegacyAmino() *codec.LegacyAmino {
	return app.legacyAmino
}

// AppCodec returns SimApp's app codec.
//
// NOTE: This is solely to be used for testing purposes as it may be desirable
// for modules to register their own custom testing types.
func (app *Gravity) AppCodec() codec.Codec {
	return app.appCodec
}

// InterfaceRegistry returns SimApp's InterfaceRegistry
func (app *Gravity) InterfaceRegistry() types.InterfaceRegistry {
	return app.interfaceRegistry
}

// GetKey returns the KVStoreKey for the provided store key.
//
// NOTE: This is solely to be used for testing purposes.
func (app *Gravity) GetKey(storeKey string) *storetypes.KVStoreKey {
	return app.keys[storeKey]
}

// GetTKey returns the TransientStoreKey for the provided store key.
//
// NOTE: This is solely to be used for testing purposes.
func (app *Gravity) GetTKey(storeKey string) *storetypes.TransientStoreKey {
	return app.tKeys[storeKey]
}

// GetMemKey returns the MemStoreKey for the provided mem key.
//
// NOTE: This is solely used for testing purposes.
func (app *Gravity) GetMemKey(storeKey string) *storetypes.MemoryStoreKey {
	return app.memKeys[storeKey]
}

// GetSubspace returns a param subspace for a given module name.
func (app *Gravity) GetSubspace(moduleName string) paramstypes.Subspace {
	subspace, _ := app.paramsKeeper.GetSubspace(moduleName)
	return subspace
}

// SimulationManager implements the SimulationApp interface
func (app *Gravity) SimulationManager() *module.SimulationManager {
	return app.sm
}

// RegisterAPIRoutes registers all application module routes with the provided
// API server.
func (app *Gravity) RegisterAPIRoutes(apiSvr *api.Server, apiConfig config.APIConfig) {
	clientCtx := apiSvr.ClientCtx
	cmtservice.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)
	authtx.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)
	ModuleBasics.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)
	ModuleBasics.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)
	if apiConfig.Swagger {
		RegisterSwaggerAPI(clientCtx, apiSvr.Router)
	}
}

// RegisterSwaggerAPI registers swagger route with API Server
func RegisterSwaggerAPI(ctx client.Context, rtr *mux.Router) {
	statikFS, err := fs.New()
	if err != nil {
		panic(err)
	}

	staticServer := http.FileServer(statikFS)
	rtr.PathPrefix("/swagger/").Handler(http.StripPrefix("/swagger/", staticServer))
}

// RegisterTxService implements the Application.RegisterTxService method.
func (app *Gravity) RegisterTxService(clientCtx client.Context) {
	authtx.RegisterTxService(app.BaseApp.GRPCQueryRouter(), clientCtx, app.BaseApp.Simulate, app.interfaceRegistry)
}

// RegisterTendermintService implements the Application.RegisterTendermintService method.
func (app *Gravity) RegisterTendermintService(clientCtx client.Context) {
	cmtApp := server.NewCometABCIWrapper(app)
	cmtservice.RegisterTendermintService(
		clientCtx,
		app.BaseApp.GRPCQueryRouter(),
		app.interfaceRegistry,
		cmtApp.Query,
	)
}

func (app *Gravity) RegisterNodeService(clientCtx client.Context, cfg config.Config) {
	nodeservice.RegisterNodeService(clientCtx, app.GRPCQueryRouter(), cfg)
}

// GetMaccPerms returns a mapping of the application's module account permissions.
func GetMaccPerms() map[string][]string {
	modAccPerms := make(map[string][]string)
	for k, v := range maccPerms {
		modAccPerms[k] = v
	}
	return modAccPerms
}

// initParamsKeeper init params keeper and its subspaces
func initParamsKeeper(appCodec codec.BinaryCodec, legacyAmino *codec.LegacyAmino, key, tkey storetypes.StoreKey) paramskeeper.Keeper {
	paramsKeeper := paramskeeper.NewKeeper(appCodec, legacyAmino, key, tkey)

	paramsKeeper.Subspace(authtypes.ModuleName)
	paramsKeeper.Subspace(banktypes.ModuleName)
	paramsKeeper.Subspace(stakingtypes.ModuleName)
	paramsKeeper.Subspace(minttypes.ModuleName)
	paramsKeeper.Subspace(distrtypes.ModuleName)
	paramsKeeper.Subspace(slashingtypes.ModuleName)
	paramsKeeper.Subspace(govtypes.ModuleName)
	paramsKeeper.Subspace(crisistypes.ModuleName)
	paramsKeeper.Subspace(ibctransfertypes.ModuleName)
	paramsKeeper.Subspace(gravitytypes.DefaultParamspace)
	keyTable := ibcclienttypes.ParamKeyTable()
	keyTable.RegisterParamSet(&ibcconnectiontypes.Params{})
	paramsKeeper.Subspace(ibcexported.ModuleName).WithKeyTable(keyTable)
	paramsKeeper.Subspace(icahosttypes.SubModuleName)

	return paramsKeeper
}

// Registers handlers for all our upgrades
func (app *Gravity) registerUpgradeHandlers() {

}

// Sets up the StoreLoader for new, deleted, or renamed modules
func (app *Gravity) registerStoreLoaders() {
	// Read the upgrade height and name from previous execution
	upgradeInfo, err := app.upgradeKeeper.ReadUpgradeInfoFromDisk()
	if err != nil {
		panic(fmt.Sprintf("failed to read upgrade info from disk %s", err))
	}
	if app.upgradeKeeper.IsSkipHeight(upgradeInfo.Height) {
		return
	}

}
