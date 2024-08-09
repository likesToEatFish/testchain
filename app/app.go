package app

import (
	"io"

	_ "cosmossdk.io/api/cosmos/tx/config/v1" // import for side-effects
	clienthelpers "cosmossdk.io/client/v2/helpers"
	"cosmossdk.io/depinject"
	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"
	_ "cosmossdk.io/x/circuit" // import for side-effects
	circuitkeeper "cosmossdk.io/x/circuit/keeper"
	_ "cosmossdk.io/x/evidence" // import for side-effects
	evidencekeeper "cosmossdk.io/x/evidence/keeper"
	feegrantkeeper "cosmossdk.io/x/feegrant/keeper"
	_ "cosmossdk.io/x/feegrant/module" // import for side-effects
	nftkeeper "cosmossdk.io/x/nft/keeper"
	_ "cosmossdk.io/x/nft/module" // import for side-effects
	_ "cosmossdk.io/x/upgrade"    // import for side-effects
	upgradekeeper "cosmossdk.io/x/upgrade/keeper"
	abci "github.com/cometbft/cometbft/abci/types"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/server/api"
	"github.com/cosmos/cosmos-sdk/server/config"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authsims "github.com/cosmos/cosmos-sdk/x/auth/simulation"
	_ "github.com/cosmos/cosmos-sdk/x/auth/tx/config" // import for side-effects
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	_ "github.com/cosmos/cosmos-sdk/x/auth/vesting" // import for side-effects
	authzkeeper "github.com/cosmos/cosmos-sdk/x/authz/keeper"
	_ "github.com/cosmos/cosmos-sdk/x/authz/module" // import for side-effects
	_ "github.com/cosmos/cosmos-sdk/x/bank"         // import for side-effects
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	_ "github.com/cosmos/cosmos-sdk/x/consensus" // import for side-effects
	consensuskeeper "github.com/cosmos/cosmos-sdk/x/consensus/keeper"
	_ "github.com/cosmos/cosmos-sdk/x/crisis" // import for side-effects
	crisiskeeper "github.com/cosmos/cosmos-sdk/x/crisis/keeper"
	_ "github.com/cosmos/cosmos-sdk/x/distribution" // import for side-effects
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/cosmos/cosmos-sdk/x/gov"
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	groupkeeper "github.com/cosmos/cosmos-sdk/x/group/keeper"
	_ "github.com/cosmos/cosmos-sdk/x/group/module" // import for side-effects
	_ "github.com/cosmos/cosmos-sdk/x/mint"         // import for side-effects
	mintkeeper "github.com/cosmos/cosmos-sdk/x/mint/keeper"
	_ "github.com/cosmos/cosmos-sdk/x/params" // import for side-effects
	paramsclient "github.com/cosmos/cosmos-sdk/x/params/client"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	_ "github.com/cosmos/cosmos-sdk/x/slashing" // import for side-effects
	slashingkeeper "github.com/cosmos/cosmos-sdk/x/slashing/keeper"
	_ "github.com/cosmos/cosmos-sdk/x/staking" // import for side-effects
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	_ "github.com/cosmos/ibc-go/modules/capability" // import for side-effects
	capabilitykeeper "github.com/cosmos/ibc-go/modules/capability/keeper"
	_ "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts" // import for side-effects
	icacontrollerkeeper "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/controller/keeper"
	icahostkeeper "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/host/keeper"
	_ "github.com/cosmos/ibc-go/v8/modules/apps/29-fee" // import for side-effects
	ibcfeekeeper "github.com/cosmos/ibc-go/v8/modules/apps/29-fee/keeper"
	ibctransferkeeper "github.com/cosmos/ibc-go/v8/modules/apps/transfer/keeper"
	ibckeeper "github.com/cosmos/ibc-go/v8/modules/core/keeper"

	"cosmossdk.io/math"
	"github.com/cometbft/cometbft/crypto"
	"github.com/cometbft/cometbft/libs/bytes"
	tmos "github.com/cometbft/cometbft/libs/os"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	testchainmodulekeeper "testchain/x/testchain/keeper"
	// this line is used by starport scaffolding # stargate/app/moduleImport

	"testchain/docs"
)

const (
	AccountAddressPrefix = "cosmos"
	Name                 = "testchain"
)

var (
	// DefaultNodeHome default home directories for the application daemon
	DefaultNodeHome string
)

var (
	_ runtime.AppI            = (*App)(nil)
	_ servertypes.Application = (*App)(nil)
)

// App extends an ABCI application, but with most of its parameters exported.
// They are exported for convenience in creating helper functions, as object
// capabilities aren't needed for testing.
type App struct {
	*runtime.App
	legacyAmino       *codec.LegacyAmino
	appCodec          codec.Codec
	txConfig          client.TxConfig
	interfaceRegistry codectypes.InterfaceRegistry

	// keepers
	AccountKeeper         authkeeper.AccountKeeper
	BankKeeper            bankkeeper.Keeper
	StakingKeeper         *stakingkeeper.Keeper
	DistrKeeper           distrkeeper.Keeper
	ConsensusParamsKeeper consensuskeeper.Keeper

	SlashingKeeper       slashingkeeper.Keeper
	MintKeeper           mintkeeper.Keeper
	GovKeeper            *govkeeper.Keeper
	CrisisKeeper         *crisiskeeper.Keeper
	UpgradeKeeper        *upgradekeeper.Keeper
	ParamsKeeper         paramskeeper.Keeper
	AuthzKeeper          authzkeeper.Keeper
	EvidenceKeeper       evidencekeeper.Keeper
	FeeGrantKeeper       feegrantkeeper.Keeper
	GroupKeeper          groupkeeper.Keeper
	NFTKeeper            nftkeeper.Keeper
	CircuitBreakerKeeper circuitkeeper.Keeper

	// IBC
	IBCKeeper           *ibckeeper.Keeper // IBC Keeper must be a pointer in the app, so we can SetRouter on it correctly
	CapabilityKeeper    *capabilitykeeper.Keeper
	IBCFeeKeeper        ibcfeekeeper.Keeper
	ICAControllerKeeper icacontrollerkeeper.Keeper
	ICAHostKeeper       icahostkeeper.Keeper
	TransferKeeper      ibctransferkeeper.Keeper

	// Scoped IBC
	ScopedIBCKeeper           capabilitykeeper.ScopedKeeper
	ScopedIBCTransferKeeper   capabilitykeeper.ScopedKeeper
	ScopedICAControllerKeeper capabilitykeeper.ScopedKeeper
	ScopedICAHostKeeper       capabilitykeeper.ScopedKeeper
	ScopedKeepers             map[string]capabilitykeeper.ScopedKeeper

	TestchainKeeper testchainmodulekeeper.Keeper
	// this line is used by starport scaffolding # stargate/app/keeperDeclaration

	// simulation manager
	sm *module.SimulationManager
}

func init() {
	var err error
	clienthelpers.EnvPrefix = Name
	DefaultNodeHome, err = clienthelpers.GetNodeHomeDirectory("." + Name)
	if err != nil {
		panic(err)
	}
}

// getGovProposalHandlers return the chain proposal handlers.
func getGovProposalHandlers() []govclient.ProposalHandler {
	var govProposalHandlers []govclient.ProposalHandler
	// this line is used by starport scaffolding # stargate/app/govProposalHandlers

	govProposalHandlers = append(govProposalHandlers,
		paramsclient.ProposalHandler,
		// this line is used by starport scaffolding # stargate/app/govProposalHandler
	)

	return govProposalHandlers
}

// AppConfig returns the default app config.
func AppConfig() depinject.Config {
	return depinject.Configs(
		appConfig,
		// Alternatively, load the app config from a YAML file.
		// appconfig.LoadYAML(AppConfigYAML),
		depinject.Supply(
			// supply custom module basics
			map[string]module.AppModuleBasic{
				genutiltypes.ModuleName: genutil.NewAppModuleBasic(genutiltypes.DefaultMessageValidator),
				govtypes.ModuleName:     gov.NewAppModuleBasic(getGovProposalHandlers()),
				// this line is used by starport scaffolding # stargate/appConfig/moduleBasic
			},
		),
	)
}

// New returns a reference to an initialized App.
func New(
	logger log.Logger,
	db dbm.DB,
	traceStore io.Writer,
	loadLatest bool,
	appOpts servertypes.AppOptions,
	baseAppOptions ...func(*baseapp.BaseApp),
) (*App, error) {
	var (
		app        = &App{ScopedKeepers: make(map[string]capabilitykeeper.ScopedKeeper)}
		appBuilder *runtime.AppBuilder

		// merge the AppConfig and other configuration in one config
		appConfig = depinject.Configs(
			AppConfig(),
			depinject.Supply(
				appOpts, // supply app options
				logger,  // supply logger
				// Supply with IBC keeper getter for the IBC modules with App Wiring.
				// The IBC Keeper cannot be passed because it has not been initiated yet.
				// Passing the getter, the app IBC Keeper will always be accessible.
				// This needs to be removed after IBC supports App Wiring.
				app.GetIBCKeeper,
				app.GetCapabilityScopedKeeper,

				// here alternative options can be supplied to the DI container.
				// those options can be used f.e to override the default behavior of some modules.
				// for instance supplying a custom address codec for not using bech32 addresses.
				// read the depinject documentation and depinject module wiring for more information
				// on available options and how to use them.
			),
		)
	)

	if err := depinject.Inject(appConfig,
		&appBuilder,
		&app.appCodec,
		&app.legacyAmino,
		&app.txConfig,
		&app.interfaceRegistry,
		&app.AccountKeeper,
		&app.BankKeeper,
		&app.StakingKeeper,
		&app.DistrKeeper,
		&app.ConsensusParamsKeeper,
		&app.SlashingKeeper,
		&app.MintKeeper,
		&app.GovKeeper,
		&app.CrisisKeeper,
		&app.UpgradeKeeper,
		&app.ParamsKeeper,
		&app.AuthzKeeper,
		&app.EvidenceKeeper,
		&app.FeeGrantKeeper,
		&app.NFTKeeper,
		&app.GroupKeeper,
		&app.CircuitBreakerKeeper,
		&app.TestchainKeeper,
		// this line is used by starport scaffolding # stargate/app/keeperDefinition
	); err != nil {
		panic(err)
	}

	// add to default baseapp options
	// enable optimistic execution
	baseAppOptions = append(baseAppOptions, baseapp.SetOptimisticExecution())

	// build app
	app.App = appBuilder.Build(db, traceStore, baseAppOptions...)

	// register legacy modules
	if err := app.registerIBCModules(appOpts); err != nil {
		return nil, err
	}

	// register streaming services
	if err := app.RegisterStreamingServices(appOpts, app.kvStoreKeys()); err != nil {
		return nil, err
	}

	/****  Module Options ****/

	app.ModuleManager.RegisterInvariants(app.CrisisKeeper)

	// create the simulation manager and define the order of the modules for deterministic simulations
	overrideModules := map[string]module.AppModuleSimulation{
		authtypes.ModuleName: auth.NewAppModule(app.appCodec, app.AccountKeeper, authsims.RandomGenesisAccounts, app.GetSubspace(authtypes.ModuleName)),
	}
	app.sm = module.NewSimulationManagerFromAppModules(app.ModuleManager.Modules, overrideModules)
	app.sm.RegisterStoreDecoders()

	// A custom InitChainer sets if extra pre-init-genesis logic is required.
	// This is necessary for manually registered modules that do not support app wiring.
	// Manually set the module version map as shown below.
	// The upgrade module will automatically handle de-duplication of the module version map.
	app.SetInitChainer(func(ctx sdk.Context, req *abci.RequestInitChain) (*abci.ResponseInitChain, error) {
		if err := app.UpgradeKeeper.SetModuleVersionMap(ctx, app.ModuleManager.GetVersionMap()); err != nil {
			return nil, err
		}
		return app.App.InitChainer(ctx, req)
	})

	if err := app.Load(loadLatest); err != nil {
		return nil, err
	}

	return app, nil
}

// LegacyAmino returns App's amino codec.
//
// NOTE: This is solely to be used for testing purposes as it may be desirable
// for modules to register their own custom testing types.
func (app *App) LegacyAmino() *codec.LegacyAmino {
	return app.legacyAmino
}

// AppCodec returns App's app codec.
//
// NOTE: This is solely to be used for testing purposes as it may be desirable
// for modules to register their own custom testing types.
func (app *App) AppCodec() codec.Codec {
	return app.appCodec
}

// GetKey returns the KVStoreKey for the provided store key.
func (app *App) GetKey(storeKey string) *storetypes.KVStoreKey {
	kvStoreKey, ok := app.UnsafeFindStoreKey(storeKey).(*storetypes.KVStoreKey)
	if !ok {
		return nil
	}
	return kvStoreKey
}

// GetMemKey returns the MemoryStoreKey for the provided store key.
func (app *App) GetMemKey(storeKey string) *storetypes.MemoryStoreKey {
	key, ok := app.UnsafeFindStoreKey(storeKey).(*storetypes.MemoryStoreKey)
	if !ok {
		return nil
	}

	return key
}

// kvStoreKeys returns all the kv store keys registered inside App.
func (app *App) kvStoreKeys() map[string]*storetypes.KVStoreKey {
	keys := make(map[string]*storetypes.KVStoreKey)
	for _, k := range app.GetStoreKeys() {
		if kv, ok := k.(*storetypes.KVStoreKey); ok {
			keys[kv.Name()] = kv
		}
	}

	return keys
}

// GetSubspace returns a param subspace for a given module name.
func (app *App) GetSubspace(moduleName string) paramstypes.Subspace {
	subspace, _ := app.ParamsKeeper.GetSubspace(moduleName)
	return subspace
}

// GetIBCKeeper returns the IBC keeper.
func (app *App) GetIBCKeeper() *ibckeeper.Keeper {
	return app.IBCKeeper
}

// GetCapabilityScopedKeeper returns the capability scoped keeper.
func (app *App) GetCapabilityScopedKeeper(moduleName string) capabilitykeeper.ScopedKeeper {
	sk, ok := app.ScopedKeepers[moduleName]
	if !ok {
		sk = app.CapabilityKeeper.ScopeToModule(moduleName)
		app.ScopedKeepers[moduleName] = sk
	}
	return sk
}

// SimulationManager implements the SimulationApp interface.
func (app *App) SimulationManager() *module.SimulationManager {
	return app.sm
}

// RegisterAPIRoutes registers all application module routes with the provided
// API server.
func (app *App) RegisterAPIRoutes(apiSvr *api.Server, apiConfig config.APIConfig) {
	app.App.RegisterAPIRoutes(apiSvr, apiConfig)
	// register swagger API in app.go so that other applications can override easily
	if err := server.RegisterSwaggerAPI(apiSvr.ClientCtx, apiSvr.Router, apiConfig.Swagger); err != nil {
		panic(err)
	}

	// register app's OpenAPI routes.
	docs.RegisterOpenAPIService(Name, apiSvr.Router)
}

// GetMaccPerms returns a copy of the module account permissions
//
// NOTE: This is solely to be used for testing purposes.
func GetMaccPerms() map[string][]string {
	dup := make(map[string][]string)
	for _, perms := range moduleAccPerms {
		dup[perms.Account] = perms.Permissions
	}
	return dup
}

// BlockedAddresses returns all the app's blocked account addresses.
func BlockedAddresses() map[string]bool {
	result := make(map[string]bool)
	if len(blockAccAddrs) > 0 {
		for _, addr := range blockAccAddrs {
			result[addr] = true
		}
	} else {
		for addr := range GetMaccPerms() {
			result[addr] = true
		}
	}
	return result
}

// InitAppForTestnet is broken down into two sections:
// Required Changes: Changes that, if not made, will cause the testnet to halt or panic
// Optional Changes: Changes to customize the testnet to one's liking (lower vote times, fund accounts, etc)
func InitAppForTestnet(app *App, newValAddr bytes.HexBytes, newValPubKey crypto.PubKey, newOperatorAddress, upgradeToTrigger string) *App {
	//
	// Required Changes:
	//

	ctx := app.BaseApp.NewUncachedContext(true, tmproto.Header{})
	pubkey := &ed25519.PubKey{Key: newValPubKey.Bytes()}
	pubkeyAny, err := codectypes.NewAnyWithValue(pubkey)
	if err != nil {
		tmos.Exit(err.Error())
	}

	// STAKING
	//

	// Create Validator struct for our new validator.
	_, bz, err := bech32.DecodeAndConvert(newOperatorAddress)
	if err != nil {
		tmos.Exit(err.Error())
	}
	bech32Addr, err := bech32.ConvertAndEncode("cosmosvaloper", bz)
	if err != nil {
		tmos.Exit(err.Error())
	}
	newVal := stakingtypes.Validator{
		OperatorAddress: bech32Addr,
		ConsensusPubkey: pubkeyAny,
		Jailed:          false,
		Status:          stakingtypes.Bonded,
		Tokens:          math.NewInt(900000000000000),
		DelegatorShares: math.LegacyMustNewDecFromStr("10000000"),
		Description: stakingtypes.Description{
			Moniker: "Testnet Validator",
		},
		Commission: stakingtypes.Commission{
			CommissionRates: stakingtypes.CommissionRates{
				Rate:          math.LegacyMustNewDecFromStr("0.05"),
				MaxRate:       math.LegacyMustNewDecFromStr("0.1"),
				MaxChangeRate: math.LegacyMustNewDecFromStr("0.05"),
			},
		},
		MinSelfDelegation: math.OneInt(),
	}

	validator, err := app.StakingKeeper.ValidatorAddressCodec().StringToBytes(newVal.GetOperator())
	if err != nil {
		tmos.Exit(err.Error())
	}

	// Remove all validators from power store
	stakingKey := app.GetKey(stakingtypes.ModuleName)
	stakingStore := ctx.KVStore(stakingKey)
	iterator, err := app.StakingKeeper.ValidatorsPowerStoreIterator(ctx)
	if err != nil {
		tmos.Exit(err.Error())
	}
	for ; iterator.Valid(); iterator.Next() {
		stakingStore.Delete(iterator.Key())
	}
	iterator.Close()

	// Remove all valdiators from last validators store
	iterator, err = app.StakingKeeper.LastValidatorsIterator(ctx)
	if err != nil {
		tmos.Exit(err.Error())
	}
	for ; iterator.Valid(); iterator.Next() {
		stakingStore.Delete(iterator.Key())
	}
	iterator.Close()

	// Remove all validators from validators store
	iterator = stakingStore.Iterator(stakingtypes.ValidatorsKey, storetypes.PrefixEndBytes(stakingtypes.ValidatorsKey))
	for ; iterator.Valid(); iterator.Next() {
		stakingStore.Delete(iterator.Key())
	}
	iterator.Close()

	// Remove all validators from unbonding queue
	iterator = stakingStore.Iterator(stakingtypes.ValidatorQueueKey, storetypes.PrefixEndBytes(stakingtypes.ValidatorQueueKey))
	for ; iterator.Valid(); iterator.Next() {
		stakingStore.Delete(iterator.Key())
	}
	iterator.Close()

	// Add our validator to power and last validators store
	app.StakingKeeper.SetValidator(ctx, newVal)
	err = app.StakingKeeper.SetValidatorByConsAddr(ctx, newVal)
	if err != nil {
		tmos.Exit(err.Error())
	}
	app.StakingKeeper.SetValidatorByPowerIndex(ctx, newVal)
	app.StakingKeeper.SetLastValidatorPower(ctx, validator, 0)
	if err := app.StakingKeeper.Hooks().AfterValidatorCreated(ctx, validator); err != nil {
		panic(err)
	}

	// DISTRIBUTION
	//

	// Initialize records for this validator across all distribution stores
	app.DistrKeeper.SetValidatorHistoricalRewards(ctx, validator, 0, distrtypes.NewValidatorHistoricalRewards(sdk.DecCoins{}, 1))
	app.DistrKeeper.SetValidatorCurrentRewards(ctx, validator, distrtypes.NewValidatorCurrentRewards(sdk.DecCoins{}, 1))
	app.DistrKeeper.SetValidatorAccumulatedCommission(ctx, validator, distrtypes.InitialValidatorAccumulatedCommission())
	app.DistrKeeper.SetValidatorOutstandingRewards(ctx, validator, distrtypes.ValidatorOutstandingRewards{Rewards: sdk.DecCoins{}})

	// SLASHING
	//

	// Set validator signing info for our new validator.
	newConsAddr := sdk.ConsAddress(newValAddr.Bytes())
	newValidatorSigningInfo := slashingtypes.ValidatorSigningInfo{
		Address:     newConsAddr.String(),
		StartHeight: app.LastBlockHeight() - 1,
		Tombstoned:  false,
	}
	app.SlashingKeeper.SetValidatorSigningInfo(ctx, newConsAddr, newValidatorSigningInfo)

	//
	// Optional Changes:
	//

	// BANK
	//

	defaultCoins := sdk.NewCoins(
		sdk.NewInt64Coin("stake", 1000000000))

	// account was generated from the following mnmenic:
	// ==============account 1==================================
	// mnemonic: "heart regret vivid nominee race scatter version coyote milk cherry axis mirror risk beef only elevator salon exact tomato lesson approve ridge cross capable"
	// address:"cosmos1f7twgcq4ypzg7y24wuywy06xmdet8pc4473tnq"
	// ==============account 2==================================
	// mnemonic:"soap step crash ceiling path virtual this armor accident pond share track spice woman vault discover share holiday inquiry oak shine scrub bulb arrive"
	// address:"cosmos1w7f3xx7e75p4l7qdym5msqem9rd4dyc4752spg"
	// ==============account 3==================================
	// mnemonic:"travel jelly basic visa apart kidney piano lumber elevator fat unknown guard matter used high drastic umbrella humble crush stock banner enlist mule unique"
	// address:"cosmos1g9v3zjt6rfkwm4s8sw9wu4jgz9me8pn27f8nyc"
	// ==============account 4==================================
	// mnemonic:"improve fun aim fringe machine shed repair olympic copper buddy road used trial liquid energy diamond orange lock time exact away change icon spike"
	// address:"cosmos1qvuhm5m644660nd8377d6l7yz9e9hhm9evmx3x"

	localAccounts := []sdk.AccAddress{
		sdk.MustAccAddressFromBech32("cosmos1f7twgcq4ypzg7y24wuywy06xmdet8pc4473tnq"),
		sdk.MustAccAddressFromBech32("cosmos1w7f3xx7e75p4l7qdym5msqem9rd4dyc4752spg"),
		sdk.MustAccAddressFromBech32("cosmos1g9v3zjt6rfkwm4s8sw9wu4jgz9me8pn27f8nyc"),
		sdk.MustAccAddressFromBech32("cosmos1qvuhm5m644660nd8377d6l7yz9e9hhm9evmx3x"),
	}

	// Fund local accounts
	for _, account := range localAccounts {
		err := app.BankKeeper.MintCoins(ctx, minttypes.ModuleName, defaultCoins)
		if err != nil {
			tmos.Exit(err.Error())
		}
		err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, account, defaultCoins)
		if err != nil {
			tmos.Exit(err.Error())
		}
	}

	// Fund edgenet faucet
	faucetCoins := sdk.NewCoins(
		sdk.NewInt64Coin("stake", 1000000000000))
	err = app.BankKeeper.MintCoins(ctx, minttypes.ModuleName, faucetCoins)
	if err != nil {
		tmos.Exit(err.Error())
	}
	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, sdk.MustAccAddressFromBech32("cosmos1f7twgcq4ypzg7y24wuywy06xmdet8pc4473tnq"), faucetCoins)
	if err != nil {
		tmos.Exit(err.Error())
	}

	// Mars bank account
	marsCoins := sdk.NewCoins(
		sdk.NewInt64Coin("stake", 10000000000000),
		sdk.NewInt64Coin("ibc/D189335C6E4A68B513C10AB227BF1C1D38C746766278BA3EEB4FB14124F1D858", 3000000000000),
		sdk.NewInt64Coin("ibc/EA1D43981D5C9A1C4AAEA9C23BB1D4FA126BA9BC7020A25E0AE4AA841EA25DC5", 3000000000000000000))
	err = app.BankKeeper.MintCoins(ctx, minttypes.ModuleName, marsCoins)
	if err != nil {
		tmos.Exit(err.Error())
	}
	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, sdk.MustAccAddressFromBech32("cosmos1f7twgcq4ypzg7y24wuywy06xmdet8pc4473tnq"), marsCoins)
	if err != nil {
		tmos.Exit(err.Error())
	}

	return app
}
