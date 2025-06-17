package app

import (
	"context"
	"encoding/json"
	"fmt"
	v10 "freemasonry.cc/blockchain/app/upgrades/v10"
	v11 "freemasonry.cc/blockchain/app/upgrades/v11"
	v12 "freemasonry.cc/blockchain/app/upgrades/v12"
	v13 "freemasonry.cc/blockchain/app/upgrades/v13"
	v14 "freemasonry.cc/blockchain/app/upgrades/v14"
	v2 "freemasonry.cc/blockchain/app/upgrades/v2"
	v3 "freemasonry.cc/blockchain/app/upgrades/v3"
	v4 "freemasonry.cc/blockchain/app/upgrades/v4"
	v5 "freemasonry.cc/blockchain/app/upgrades/v5"
	v6 "freemasonry.cc/blockchain/app/upgrades/v6"
	v7 "freemasonry.cc/blockchain/app/upgrades/v7"
	v8 "freemasonry.cc/blockchain/app/upgrades/v8"
	v9 "freemasonry.cc/blockchain/app/upgrades/v9"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"

	contractkeeper "freemasonry.cc/blockchain/x/contract/keeper"
	contracttypes "freemasonry.cc/blockchain/x/contract/types"
	"github.com/cosmos/cosmos-sdk/x/group"
	groupkeeper "github.com/cosmos/cosmos-sdk/x/group/keeper"
	groupModule "github.com/cosmos/cosmos-sdk/x/group/module"
	"github.com/cosmos/cosmos-sdk/x/mint"
	mintkeeper "github.com/cosmos/cosmos-sdk/x/mint/keeper"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/evmos/evmos/v10/x/claims"
	claimstypes "github.com/evmos/evmos/v10/x/claims/types"
	"github.com/evmos/evmos/v10/x/erc20"
	"github.com/evmos/evmos/v10/x/ibc/transfer"
	"github.com/evmos/evmos/v10/x/recovery"
	"github.com/evmos/evmos/v10/x/vesting"
	vestingtypes "github.com/evmos/evmos/v10/x/vesting/types"

	"github.com/gorilla/mux"
	"github.com/rakyll/statik/fs"
	"github.com/spf13/cast"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	tmos "github.com/tendermint/tendermint/libs/os"
	dbm "github.com/tendermint/tm-db"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/grpc/node"
	"github.com/cosmos/cosmos-sdk/client/grpc/tmservice"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/server/api"
	"github.com/cosmos/cosmos-sdk/server/config"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/cosmos/cosmos-sdk/simapp"
	simappparams "github.com/cosmos/cosmos-sdk/simapp/params"
	"github.com/cosmos/cosmos-sdk/store/streaming"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authsims "github.com/cosmos/cosmos-sdk/x/auth/simulation"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	authzkeeper "github.com/cosmos/cosmos-sdk/x/authz/keeper"
	authzmodule "github.com/cosmos/cosmos-sdk/x/authz/module"
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
	"github.com/cosmos/cosmos-sdk/x/feegrant"
	feegrantkeeper "github.com/cosmos/cosmos-sdk/x/feegrant/keeper"
	feegrantmodule "github.com/cosmos/cosmos-sdk/x/feegrant/module"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/cosmos/cosmos-sdk/x/gov"
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	govv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
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
	ibctestingtypes "github.com/cosmos/ibc-go/v5/testing/types"

	ibctransfer "github.com/cosmos/ibc-go/v5/modules/apps/transfer"
	ibctransfertypes "github.com/cosmos/ibc-go/v5/modules/apps/transfer/types"
	ibc "github.com/cosmos/ibc-go/v5/modules/core"
	ibcclient "github.com/cosmos/ibc-go/v5/modules/core/02-client"
	ibcclientclient "github.com/cosmos/ibc-go/v5/modules/core/02-client/client"
	ibcclienttypes "github.com/cosmos/ibc-go/v5/modules/core/02-client/types"
	porttypes "github.com/cosmos/ibc-go/v5/modules/core/05-port/types"
	ibchost "github.com/cosmos/ibc-go/v5/modules/core/24-host"
	ibckeeper "github.com/cosmos/ibc-go/v5/modules/core/keeper"
	ibctesting "github.com/cosmos/ibc-go/v5/testing"

	"github.com/evmos/ethermint/encoding"
	"github.com/evmos/ethermint/ethereum/eip712"
	srvflags "github.com/evmos/ethermint/server/flags"
	ethermint "github.com/evmos/ethermint/types"
	"github.com/evmos/ethermint/x/evm"
	evmkeeper "github.com/evmos/ethermint/x/evm/keeper"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	"github.com/evmos/ethermint/x/evm/vm/geth"
	"github.com/evmos/ethermint/x/feemarket"
	feemarketkeeper "github.com/evmos/ethermint/x/feemarket/keeper"
	feemarkettypes "github.com/evmos/ethermint/x/feemarket/types"

	
	_ "github.com/evmos/evmos/v10/client/docs/statik"

	"freemasonry.cc/blockchain/app/ante"

	

	"freemasonry.cc/blockchain/x/chat"
	chatkeeper "freemasonry.cc/blockchain/x/chat/keeper"
	chattypes "freemasonry.cc/blockchain/x/chat/types"

	"freemasonry.cc/blockchain/x/dao"
	daokeeper "freemasonry.cc/blockchain/x/dao/keeper"
	daotypes "freemasonry.cc/blockchain/x/dao/types"

	"freemasonry.cc/blockchain/x/gateway"
	gatewaykeeper "freemasonry.cc/blockchain/x/gateway/keeper"
	gatewaytypes "freemasonry.cc/blockchain/x/gateway/types"

	"freemasonry.cc/blockchain/x/contract"
	erc20client "github.com/evmos/evmos/v10/x/erc20/client"
	erc20keeper "github.com/evmos/evmos/v10/x/erc20/keeper"
	erc20types "github.com/evmos/evmos/v10/x/erc20/types"

	transferkeeper "github.com/evmos/evmos/v10/x/ibc/transfer/keeper"

	claimskeeper "github.com/evmos/evmos/v10/x/claims/keeper"

	recoverykeeper "github.com/evmos/evmos/v10/x/recovery/keeper"
	recoverytypes "github.com/evmos/evmos/v10/x/recovery/types"
	vestingkeeper "github.com/evmos/evmos/v10/x/vesting/keeper"
)

func init() {
	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	DefaultNodeHome = filepath.Join(pwd, ".stcd")

	
	sdk.DefaultPowerReduction = ethermint.PowerReduction
	
	feemarkettypes.DefaultMinGasPrice = MainnetMinGasPrices
	feemarkettypes.DefaultMinGasMultiplier = MainnetMinGasMultiplier
	
	stakingtypes.DefaultMinCommissionRate = sdk.NewDecWithPrec(5, 2)

	claimstypes.DefaultEVMChannels = []string{
		"channel-0", 
	}
}

const Name = "stcd"

var (
	
	DefaultNodeHome string

	EncodingConfig = encoding.MakeConfig(ModuleBasics)
	
	
	
	ModuleBasics = module.NewBasicManager(
		auth.AppModuleBasic{},
		genutil.AppModuleBasic{},
		bank.AppModuleBasic{},
		capability.AppModuleBasic{},
		staking.AppModuleBasic{},
		distr.AppModuleBasic{},
		gov.NewAppModuleBasic(
			[]govclient.ProposalHandler{
				paramsclient.ProposalHandler, distrclient.ProposalHandler, upgradeclient.LegacyProposalHandler, upgradeclient.LegacyCancelProposalHandler,
				ibcclientclient.UpdateClientProposalHandler, ibcclientclient.UpgradeProposalHandler,
				
				erc20client.RegisterCoinProposalHandler, erc20client.RegisterERC20ProposalHandler, erc20client.ToggleTokenConversionProposalHandler,
			},
		),
		params.AppModuleBasic{},
		crisis.AppModuleBasic{},
		slashing.AppModuleBasic{},
		ibc.AppModuleBasic{},
		authzmodule.AppModuleBasic{},
		feegrantmodule.AppModuleBasic{},
		upgrade.AppModuleBasic{},
		evidence.AppModuleBasic{},
		transfer.AppModuleBasic{AppModuleBasic: &ibctransfer.AppModuleBasic{}},
		vesting.AppModuleBasic{},
		evm.AppModuleBasic{},
		feemarket.AppModuleBasic{},
		erc20.AppModuleBasic{},
		claims.AppModuleBasic{},
		recovery.AppModuleBasic{},
		chat.AppModuleBasic{},
		gateway.AppModuleBasic{},
		contract.AppModuleBasic{},
		dao.AppModuleBasic{},
		groupModule.AppModuleBasic{},
		mint.AppModuleBasic{},
	)

	
	maccPerms = map[string][]string{
		authtypes.FeeCollectorName:     {authtypes.Burner},
		distrtypes.ModuleName:          nil,
		stakingtypes.BondedPoolName:    {authtypes.Burner, authtypes.Staking},
		stakingtypes.NotBondedPoolName: {authtypes.Burner, authtypes.Staking},
		govtypes.ModuleName:            {authtypes.Burner},
		ibctransfertypes.ModuleName:    {authtypes.Minter, authtypes.Burner},
		evmtypes.ModuleName:            {authtypes.Minter, authtypes.Burner}, 

		erc20types.ModuleName:  {authtypes.Minter, authtypes.Burner},
		claimstypes.ModuleName: nil,

		chattypes.ModuleName:     {authtypes.Minter, authtypes.Burner},
		chattypes.ModuleBurnName: {authtypes.Burner},
		gatewaytypes.ModuleName:  nil,

		contracttypes.ModuleName:                 nil,
		contracttypes.StakeContractDeploy:        nil,
		contracttypes.TokenFactoryContractDeploy: nil,
		daotypes.ModuleName:                      {authtypes.Burner, authtypes.Staking, authtypes.Minter},
		group.ModuleName:                         {authtypes.Burner, authtypes.Staking, authtypes.Minter},
		minttypes.ModuleName:                     {authtypes.Minter},
		contracttypes.GenesisIdoReward:           {authtypes.Burner},
	}

	
	allowedReceivingModAcc = map[string]bool{
		chattypes.ModuleName:           true,
		distrtypes.ModuleName:          true, 
		contracttypes.ModuleName:       true,
		erc20types.ModuleName:          true,
		group.ModuleName:               true,
		daotypes.ModuleName:            true,
		contracttypes.GenesisIdoReward: true,
	}
)

var (
	_ servertypes.Application = (*Evmos)(nil)
	_ ibctesting.TestingApp   = (*Evmos)(nil)
)

type Evmos struct {
	*baseapp.BaseApp

	
	cdc               *codec.LegacyAmino
	appCodec          codec.Codec
	interfaceRegistry types.InterfaceRegistry

	invCheckPeriod uint

	
	keys    map[string]*storetypes.KVStoreKey
	tkeys   map[string]*storetypes.TransientStoreKey
	memKeys map[string]*storetypes.MemoryStoreKey

	
	AccountKeeper    authkeeper.AccountKeeper
	BankKeeper       bankkeeper.Keeper
	CapabilityKeeper *capabilitykeeper.Keeper
	StakingKeeper    stakingkeeper.Keeper
	SlashingKeeper   slashingkeeper.Keeper
	DistrKeeper      distrkeeper.Keeper
	GovKeeper        govkeeper.Keeper
	CrisisKeeper     crisiskeeper.Keeper
	UpgradeKeeper    upgradekeeper.Keeper
	ParamsKeeper     paramskeeper.Keeper
	FeeGrantKeeper   feegrantkeeper.Keeper
	AuthzKeeper      authzkeeper.Keeper
	IBCKeeper        *ibckeeper.Keeper 
	EvidenceKeeper   evidencekeeper.Keeper
	TransferKeeper   transferkeeper.Keeper

	
	ScopedIBCKeeper      capabilitykeeper.ScopedKeeper
	ScopedTransferKeeper capabilitykeeper.ScopedKeeper

	
	EvmKeeper       *evmkeeper.Keeper
	FeeMarketKeeper feemarketkeeper.Keeper

	
	ClaimsKeeper   *claimskeeper.Keeper
	Erc20Keeper    erc20keeper.Keeper
	VestingKeeper  vestingkeeper.Keeper
	RecoveryKeeper *recoverykeeper.Keeper

	MintKeeper     mintkeeper.Keeper
	ChatKeeper     chatkeeper.Keeper
	GatewayKeeper  gatewaykeeper.Keeper
	ContractKeeper contractkeeper.Keeper
	DaoKeeper      daokeeper.Keeper
	GroupKeeper    groupkeeper.Keeper

	
	mm *module.Manager

	
	configurator module.Configurator

	tpsCounter *tpsCounter
}

func NewEvmos(
	logger log.Logger,
	db dbm.DB,
	traceStore io.Writer,
	loadLatest bool,
	skipUpgradeHeights map[int64]bool,
	homePath string,
	invCheckPeriod uint,
	encodingConfig simappparams.EncodingConfig,
	appOpts servertypes.AppOptions,
	baseAppOptions ...func(*baseapp.BaseApp),
) *Evmos {
	appCodec := encodingConfig.Codec
	cdc := encodingConfig.Amino
	interfaceRegistry := encodingConfig.InterfaceRegistry

	eip712.SetEncodingConfig(encodingConfig)

	
	bApp := baseapp.NewBaseApp(
		Name,
		logger,
		db,
		encodingConfig.TxConfig.TxDecoder(),
		baseAppOptions...,
	)
	bApp.SetCommitMultiStoreTracer(traceStore)
	bApp.SetVersion(version.Version)
	bApp.SetInterfaceRegistry(interfaceRegistry)

	keys := sdk.NewKVStoreKeys(
		
		authtypes.StoreKey, banktypes.StoreKey, stakingtypes.StoreKey,
		distrtypes.StoreKey, slashingtypes.StoreKey,
		govtypes.StoreKey, paramstypes.StoreKey, upgradetypes.StoreKey,
		evidencetypes.StoreKey, capabilitytypes.StoreKey,
		feegrant.StoreKey, authzkeeper.StoreKey,
		
		ibchost.StoreKey, ibctransfertypes.StoreKey,
		
		evmtypes.StoreKey, feemarkettypes.StoreKey,
		
		erc20types.StoreKey,
		claimstypes.StoreKey, vestingtypes.StoreKey,

		chattypes.StoreKey,
		gatewaytypes.StoreKey,
		contracttypes.StoreKey,
		daotypes.StoreKey,
		group.StoreKey,
		minttypes.StoreKey,
	)

	
	tkeys := sdk.NewTransientStoreKeys(paramstypes.TStoreKey, evmtypes.TransientKey, feemarkettypes.TransientKey)
	memKeys := sdk.NewMemoryStoreKeys(capabilitytypes.MemStoreKey)

	
	if _, _, err := streaming.LoadStreamingServices(bApp, appOpts, appCodec, keys); err != nil {
		fmt.Printf("failed to load state streaming: %s", err)
		os.Exit(1)
	}

	app := &Evmos{
		BaseApp:           bApp,
		cdc:               cdc,
		appCodec:          appCodec,
		interfaceRegistry: interfaceRegistry,
		invCheckPeriod:    invCheckPeriod,
		keys:              keys,
		tkeys:             tkeys,
		memKeys:           memKeys,
	}

	
	app.ParamsKeeper = initParamsKeeper(appCodec, cdc, keys[paramstypes.StoreKey], tkeys[paramstypes.TStoreKey])
	
	bApp.SetParamStore(app.ParamsKeeper.Subspace(baseapp.Paramspace).WithKeyTable(paramstypes.ConsensusParamsKeyTable()))

	
	app.CapabilityKeeper = capabilitykeeper.NewKeeper(appCodec, keys[capabilitytypes.StoreKey], memKeys[capabilitytypes.MemStoreKey])

	scopedIBCKeeper := app.CapabilityKeeper.ScopeToModule(ibchost.ModuleName)
	scopedTransferKeeper := app.CapabilityKeeper.ScopeToModule(ibctransfertypes.ModuleName)

	
	
	app.CapabilityKeeper.Seal()

	
	app.AccountKeeper = authkeeper.NewAccountKeeper(
		appCodec, keys[authtypes.StoreKey], app.GetSubspace(authtypes.ModuleName), ethermint.ProtoAccount, maccPerms, sdk.GetConfig().GetBech32AccountAddrPrefix(),
	)
	bankKeeper := bankkeeper.NewBaseKeeper(
		appCodec, keys[banktypes.StoreKey], app.AccountKeeper, app.GetSubspace(banktypes.ModuleName), app.BlockedAddrs(),
	)
	app.BankKeeper = bankKeeper
	stakingKeeper := stakingkeeper.NewKeeper(
		appCodec, keys[stakingtypes.StoreKey], app.AccountKeeper, app.BankKeeper, app.GetSubspace(stakingtypes.ModuleName),
	)
	app.DistrKeeper = distrkeeper.NewKeeper(
		appCodec, keys[distrtypes.StoreKey], app.GetSubspace(distrtypes.ModuleName), app.AccountKeeper, app.BankKeeper,
		&stakingKeeper, authtypes.FeeCollectorName,
	)
	app.SlashingKeeper = slashingkeeper.NewKeeper(
		appCodec, keys[slashingtypes.StoreKey], &stakingKeeper, app.GetSubspace(slashingtypes.ModuleName),
	)
	app.CrisisKeeper = crisiskeeper.NewKeeper(
		app.GetSubspace(crisistypes.ModuleName), invCheckPeriod, app.BankKeeper, authtypes.FeeCollectorName,
	)
	app.FeeGrantKeeper = feegrantkeeper.NewKeeper(appCodec, keys[feegrant.StoreKey], app.AccountKeeper)
	app.UpgradeKeeper = upgradekeeper.NewKeeper(skipUpgradeHeights, keys[upgradetypes.StoreKey], appCodec, homePath, app.BaseApp, authtypes.NewModuleAddress(govtypes.ModuleName).String())

	app.AuthzKeeper = authzkeeper.NewKeeper(keys[authzkeeper.StoreKey], appCodec, app.MsgServiceRouter(), app.AccountKeeper)

	tracer := cast.ToString(appOpts.Get(srvflags.EVMTracer))

	
	app.FeeMarketKeeper = feemarketkeeper.NewKeeper(
		appCodec, app.GetSubspace(feemarkettypes.ModuleName), keys[feemarkettypes.StoreKey], tkeys[feemarkettypes.TransientKey],
	)

	app.EvmKeeper = evmkeeper.NewKeeper(
		appCodec, keys[evmtypes.StoreKey], tkeys[evmtypes.TransientKey], app.GetSubspace(evmtypes.ModuleName),
		app.AccountKeeper, app.BankKeeper, &stakingKeeper, app.FeeMarketKeeper,
		nil, geth.NewEVM, tracer, nil,
	)

	
	app.IBCKeeper = ibckeeper.NewKeeper(
		appCodec, keys[ibchost.StoreKey], app.GetSubspace(ibchost.ModuleName), &stakingKeeper, app.UpgradeKeeper, scopedIBCKeeper,
	)

	
	govRouter := govv1beta1.NewRouter()
	govRouter.AddRoute(govtypes.RouterKey, govv1beta1.ProposalHandler).
		AddRoute(paramproposal.RouterKey, params.NewParamChangeProposalHandler(app.ParamsKeeper)).
		AddRoute(distrtypes.RouterKey, distr.NewCommunityPoolSpendProposalHandler(app.DistrKeeper)).
		AddRoute(upgradetypes.RouterKey, upgrade.NewSoftwareUpgradeProposalHandler(app.UpgradeKeeper)).
		AddRoute(ibcclienttypes.RouterKey, ibcclient.NewClientProposalHandler(app.IBCKeeper.ClientKeeper)).
		AddRoute(erc20types.RouterKey, erc20.NewErc20ProposalHandler(&app.Erc20Keeper))

	govConfig := govtypes.DefaultConfig()
	
	govKeeper := govkeeper.NewKeeper(
		appCodec, keys[govtypes.StoreKey], app.GetSubspace(govtypes.ModuleName), app.AccountKeeper, app.BankKeeper,
		&stakingKeeper, govRouter, app.MsgServiceRouter(), govConfig,
	)

	

	app.ClaimsKeeper = claimskeeper.NewKeeper(
		appCodec, keys[claimstypes.StoreKey], app.GetSubspace(claimstypes.ModuleName),
		app.AccountKeeper, app.BankKeeper, &stakingKeeper, app.DistrKeeper,
	)

	
	app.GatewayKeeper = gatewaykeeper.NewKeeper(keys[gatewaytypes.StoreKey], appCodec, app.GetSubspace(gatewaytypes.ModuleName),
		app.AccountKeeper, app.BankKeeper, &stakingKeeper)

	
	
	
	app.StakingKeeper = *stakingKeeper.SetHooks(
		stakingtypes.NewMultiStakingHooks(
			app.DistrKeeper.Hooks(),
			app.SlashingKeeper.Hooks(),
			app.ClaimsKeeper.Hooks(),
			app.GatewayKeeper.Hooks(),
		),
	)

	app.VestingKeeper = vestingkeeper.NewKeeper(
		keys[vestingtypes.StoreKey], appCodec,
		app.AccountKeeper, app.BankKeeper, app.StakingKeeper,
	)

	app.Erc20Keeper = erc20keeper.NewKeeper(
		keys[erc20types.StoreKey], appCodec, app.GetSubspace(erc20types.ModuleName),
		app.AccountKeeper, app.BankKeeper, app.EvmKeeper, app.StakingKeeper, app.ClaimsKeeper,
	)

	app.GovKeeper = *govKeeper.SetHooks(
		govtypes.NewMultiGovHooks(
			app.ClaimsKeeper.Hooks(),
		),
	)

	app.GroupKeeper = groupkeeper.NewKeeper(keys[group.StoreKey], appCodec, app.MsgServiceRouter(), app.AccountKeeper, app.BankKeeper, group.DefaultConfig(), app.GetSubspace(group.ModuleName))

	app.TransferKeeper = transferkeeper.NewKeeper(
		appCodec, keys[ibctransfertypes.StoreKey], app.GetSubspace(ibctransfertypes.ModuleName),
		app.ClaimsKeeper, 
		app.IBCKeeper.ChannelKeeper, &app.IBCKeeper.PortKeeper,
		app.AccountKeeper, app.BankKeeper, scopedTransferKeeper,
		app.Erc20Keeper, 
	)

	app.RecoveryKeeper = recoverykeeper.NewKeeper(
		app.GetSubspace(recoverytypes.ModuleName),
		app.AccountKeeper,
		app.BankKeeper,
		app.IBCKeeper.ChannelKeeper,
		app.TransferKeeper,
		app.ClaimsKeeper,
	)

	app.ChatKeeper = chatkeeper.NewKeeper(keys[chattypes.StoreKey], appCodec, app.GetSubspace(chattypes.ModuleName),
		app.AccountKeeper, app.BankKeeper, app.GatewayKeeper)

	app.ContractKeeper = contractkeeper.NewKeeper(keys[contracttypes.StoreKey], appCodec, app.GetSubspace(contracttypes.ModuleName),
		app.AccountKeeper, app.BankKeeper, &stakingKeeper, app.EvmKeeper, app.ChatKeeper, app.GatewayKeeper, app.Erc20Keeper)

	app.DaoKeeper = daokeeper.NewKeeper(keys[daotypes.StoreKey], appCodec, app.GetSubspace(daotypes.ModuleName), app.AccountKeeper, app.BankKeeper, &stakingKeeper, app.GatewayKeeper, app.DistrKeeper, app.ChatKeeper, app.GroupKeeper, app.ContractKeeper)

	app.MintKeeper = mintkeeper.NewKeeper(app.appCodec, keys[minttypes.StoreKey], app.GetSubspace(minttypes.ModuleName), &stakingKeeper, app.AccountKeeper, app.BankKeeper, authtypes.FeeCollectorName)

	app.EvmKeeper = app.EvmKeeper.SetDaoKeeper(app.DaoKeeper)
	app.EvmKeeper = app.EvmKeeper.SetHooks(
		evmkeeper.NewMultiEvmHooks(
			app.Erc20Keeper.Hooks(),
			app.ClaimsKeeper.Hooks(),
			app.DaoKeeper.EVMHooks(),
		),
	)

	app.BankKeeper = app.BankKeeper.SetHooks(
		banktypes.NewMultiBankHooks(
			app.DaoKeeper.Hooks(),
		), &bankKeeper,
	)
	

	
	app.RecoveryKeeper.SetICS4Wrapper(app.IBCKeeper.ChannelKeeper)
	app.ClaimsKeeper.SetICS4Wrapper(app.RecoveryKeeper)

	
	transferModule := transfer.NewAppModule(app.TransferKeeper)

	

	
	var transferStack porttypes.IBCModule

	transferStack = transfer.NewIBCModule(app.TransferKeeper)
	transferStack = claims.NewIBCMiddleware(*app.ClaimsKeeper, transferStack)
	transferStack = recovery.NewIBCMiddleware(*app.RecoveryKeeper, transferStack)
	transferStack = erc20.NewIBCMiddleware(app.Erc20Keeper, transferStack)
	transferStack = dao.NewIBCMiddleware(app.DaoKeeper, transferStack)

	
	ibcRouter := porttypes.NewRouter()
	ibcRouter.AddRoute(ibctransfertypes.ModuleName, transferStack)
	app.IBCKeeper.SetRouter(ibcRouter)

	
	evidenceKeeper := evidencekeeper.NewKeeper(
		appCodec, keys[evidencetypes.StoreKey], &app.StakingKeeper, app.SlashingKeeper,
	)
	
	app.EvidenceKeeper = *evidenceKeeper

	

	
	
	skipGenesisInvariants := cast.ToBool(appOpts.Get(crisis.FlagSkipGenesisInvariants))

	
	
	app.mm = module.NewManager(
		
		genutil.NewAppModule(
			app.AccountKeeper, app.StakingKeeper, app.BaseApp.DeliverTx,
			encodingConfig.TxConfig,
		),
		auth.NewAppModule(appCodec, app.AccountKeeper, authsims.RandomGenesisAccounts),
		bank.NewAppModule(appCodec, app.BankKeeper, app.AccountKeeper),
		capability.NewAppModule(appCodec, *app.CapabilityKeeper),
		crisis.NewAppModule(&app.CrisisKeeper, skipGenesisInvariants),
		gov.NewAppModule(appCodec, app.GovKeeper, app.AccountKeeper, app.BankKeeper),
		slashing.NewAppModule(appCodec, app.SlashingKeeper, app.AccountKeeper, app.BankKeeper, app.StakingKeeper),
		distr.NewAppModule(appCodec, app.DistrKeeper, app.AccountKeeper, app.BankKeeper, app.StakingKeeper),
		staking.NewAppModule(appCodec, app.StakingKeeper, app.AccountKeeper, app.BankKeeper),
		upgrade.NewAppModule(app.UpgradeKeeper),
		evidence.NewAppModule(app.EvidenceKeeper),
		params.NewAppModule(app.ParamsKeeper),
		feegrantmodule.NewAppModule(appCodec, app.AccountKeeper, app.BankKeeper, app.FeeGrantKeeper, app.interfaceRegistry),
		authzmodule.NewAppModule(appCodec, app.AuthzKeeper, app.AccountKeeper, app.BankKeeper, app.interfaceRegistry),

		
		ibc.NewAppModule(app.IBCKeeper),
		transferModule,
		
		evm.NewAppModule(app.EvmKeeper, app.AccountKeeper),
		feemarket.NewAppModule(app.FeeMarketKeeper),
		
		erc20.NewAppModule(app.Erc20Keeper, app.AccountKeeper),
		claims.NewAppModule(appCodec, *app.ClaimsKeeper),
		vesting.NewAppModule(app.VestingKeeper, app.AccountKeeper, app.BankKeeper, app.StakingKeeper),
		recovery.NewAppModule(*app.RecoveryKeeper),

		chat.NewAppModule(app.ChatKeeper, app.AccountKeeper),
		gateway.NewAppModule(app.GatewayKeeper, app.AccountKeeper),
		contract.NewAppModule(app.ContractKeeper, app.AccountKeeper),
		dao.NewAppModule(app.DaoKeeper, app.AccountKeeper),
		groupModule.NewAppModule(appCodec, app.GroupKeeper, app.AccountKeeper, app.BankKeeper, app.interfaceRegistry),
		mint.NewAppModule(appCodec, app.MintKeeper, app.AccountKeeper, nil),
	)

	
	
	
	
	
	
	app.mm.SetOrderBeginBlockers(
		upgradetypes.ModuleName,
		capabilitytypes.ModuleName,
		feemarkettypes.ModuleName,
		evmtypes.ModuleName,
		distrtypes.ModuleName,
		slashingtypes.ModuleName,
		evidencetypes.ModuleName,
		stakingtypes.ModuleName,
		ibchost.ModuleName,
		
		ibctransfertypes.ModuleName,
		authtypes.ModuleName,
		banktypes.ModuleName,
		govtypes.ModuleName,
		crisistypes.ModuleName,
		genutiltypes.ModuleName,
		authz.ModuleName,
		feegrant.ModuleName,
		paramstypes.ModuleName,
		vestingtypes.ModuleName,
		erc20types.ModuleName,
		claimstypes.ModuleName,
		recoverytypes.ModuleName,
		chattypes.ModuleName,
		gatewaytypes.ModuleName,
		contracttypes.ModuleName,
		daotypes.ModuleName,
		group.ModuleName,
		minttypes.ModuleName,
	)

	
	app.mm.SetOrderEndBlockers(
		crisistypes.ModuleName,
		govtypes.ModuleName,
		stakingtypes.ModuleName,
		evmtypes.ModuleName,
		feemarkettypes.ModuleName,
		claimstypes.ModuleName,
		
		ibchost.ModuleName,
		ibctransfertypes.ModuleName,
		capabilitytypes.ModuleName,
		authtypes.ModuleName,
		banktypes.ModuleName,
		distrtypes.ModuleName,
		slashingtypes.ModuleName,
		genutiltypes.ModuleName,
		evidencetypes.ModuleName,
		authz.ModuleName,
		feegrant.ModuleName,
		paramstypes.ModuleName,
		upgradetypes.ModuleName,
		
		vestingtypes.ModuleName,
		erc20types.ModuleName,
		recoverytypes.ModuleName,

		chattypes.ModuleName,
		gatewaytypes.ModuleName,
		contracttypes.ModuleName,
		daotypes.ModuleName,
		group.ModuleName,
		minttypes.ModuleName,
	)

	
	
	
	
	
	app.mm.SetOrderInitGenesis(
		
		capabilitytypes.ModuleName,
		authtypes.ModuleName,
		banktypes.ModuleName,
		distrtypes.ModuleName,
		
		claimstypes.ModuleName,
		stakingtypes.ModuleName,
		slashingtypes.ModuleName,
		govtypes.ModuleName,
		ibchost.ModuleName,
		
		evmtypes.ModuleName,
		
		
		feemarkettypes.ModuleName,
		genutiltypes.ModuleName,
		evidencetypes.ModuleName,
		ibctransfertypes.ModuleName,
		authz.ModuleName,
		feegrant.ModuleName,
		paramstypes.ModuleName,
		upgradetypes.ModuleName,
		
		vestingtypes.ModuleName,
		erc20types.ModuleName,
		recoverytypes.ModuleName,
		
		crisistypes.ModuleName,

		chattypes.ModuleName,
		gatewaytypes.ModuleName,
		contracttypes.ModuleName,
		daotypes.ModuleName,
		group.ModuleName,
		minttypes.ModuleName,
	)

	app.mm.RegisterInvariants(&app.CrisisKeeper)
	app.mm.RegisterRoutes(app.Router(), app.QueryRouter(), encodingConfig.Amino)
	app.configurator = module.NewConfigurator(app.appCodec, app.MsgServiceRouter(), app.GRPCQueryRouter())
	app.mm.RegisterServices(app.configurator)

	
	

	
	app.MountKVStores(keys)
	app.MountTransientStores(tkeys)
	app.MountMemoryStores(memKeys)

	
	app.SetInitChainer(app.InitChainer)
	app.SetBeginBlocker(app.BeginBlocker)

	maxGasWanted := cast.ToUint64(appOpts.Get(srvflags.EVMMaxTxGasWanted))
	options := ante.HandlerOptions{
		AccountKeeper:          app.AccountKeeper,
		BankKeeper:             app.BankKeeper,
		ExtensionOptionChecker: nil,
		EvmKeeper:              app.EvmKeeper,
		StakingKeeper:          app.StakingKeeper,
		FeegrantKeeper:         app.FeeGrantKeeper,
		IBCKeeper:              app.IBCKeeper,
		FeeMarketKeeper:        app.FeeMarketKeeper,
		SignModeHandler:        encodingConfig.TxConfig.SignModeHandler(),
		SigGasConsumer:         SigVerificationGasConsumer,
		Cdc:                    appCodec,
		MaxTxGasWanted:         maxGasWanted,
		DaoKeeper:              app.DaoKeeper,
		ContractsKeeper:        app.ContractKeeper,
	}

	if err := options.Validate(); err != nil {
		panic(err)
	}

	app.SetAnteHandler(ante.NewAnteHandler(options))
	app.SetEndBlocker(app.EndBlocker)
	app.setupUpgradeHandlers()

	if loadLatest {
		if err := app.LoadLatestVersion(); err != nil {
			tmos.Exit(err.Error())
		}
	}

	app.ScopedIBCKeeper = scopedIBCKeeper
	app.ScopedTransferKeeper = scopedTransferKeeper

	
	app.tpsCounter = newTPSCounter(logger)
	go func() {
		
		
		_ = app.tpsCounter.start(context.Background())
	}()

	return app
}

func (app *Evmos) Name() string { return app.BaseApp.Name() }

func (app *Evmos) BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
	
	app.ScheduleForkUpgrade(ctx)
	return app.mm.BeginBlock(ctx, req)
}

func (app *Evmos) EndBlocker(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
	return app.mm.EndBlock(ctx, req)
}

func (app *Evmos) DeliverTx(req abci.RequestDeliverTx) (res abci.ResponseDeliverTx) {
	defer func() {
		
		
		if res.IsErr() {
			app.tpsCounter.incrementFailure()
		} else {
			app.tpsCounter.incrementSuccess()
		}
	}()
	return app.BaseApp.DeliverTx(req)
}

func (app *Evmos) InitChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	var genesisState simapp.GenesisState
	if err := json.Unmarshal(req.AppStateBytes, &genesisState); err != nil {
		panic(err)
	}

	app.UpgradeKeeper.SetModuleVersionMap(ctx, app.mm.GetVersionMap())

	return app.mm.InitGenesis(ctx, app.appCodec, genesisState)
}

func (app *Evmos) LoadHeight(height int64) error {
	return app.LoadVersion(height)
}

func (app *Evmos) ModuleAccountAddrs() map[string]bool {
	modAccAddrs := make(map[string]bool)

	accs := make([]string, 0, len(maccPerms))
	for k := range maccPerms {
		accs = append(accs, k)
	}
	sort.Strings(accs)

	for _, acc := range accs {
		modAccAddrs[authtypes.NewModuleAddress(acc).String()] = true
	}

	return modAccAddrs
}

func (app *Evmos) BlockedAddrs() map[string]bool {
	blockedAddrs := make(map[string]bool)

	accs := make([]string, 0, len(maccPerms))
	for k := range maccPerms {
		accs = append(accs, k)
	}
	sort.Strings(accs)

	for _, acc := range accs {
		blockedAddrs[authtypes.NewModuleAddress(acc).String()] = !allowedReceivingModAcc[acc]
	}

	return blockedAddrs
}


func (app *Evmos) LegacyAmino() *codec.LegacyAmino {
	return app.cdc
}


func (app *Evmos) AppCodec() codec.Codec {
	return app.appCodec
}

func (app *Evmos) InterfaceRegistry() types.InterfaceRegistry {
	return app.interfaceRegistry
}


func (app *Evmos) GetKey(storeKey string) *storetypes.KVStoreKey {
	return app.keys[storeKey]
}


func (app *Evmos) GetTKey(storeKey string) *storetypes.TransientStoreKey {
	return app.tkeys[storeKey]
}


func (app *Evmos) GetMemKey(storeKey string) *storetypes.MemoryStoreKey {
	return app.memKeys[storeKey]
}


func (app *Evmos) GetSubspace(moduleName string) paramstypes.Subspace {
	subspace, _ := app.ParamsKeeper.GetSubspace(moduleName)
	return subspace
}

func (app *Evmos) RegisterAPIRoutes(apiSvr *api.Server, apiConfig config.APIConfig) {
	clientCtx := apiSvr.ClientCtx

	
	authtx.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)
	
	tmservice.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)
	
	node.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)

	
	ModuleBasics.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)

	
	if apiConfig.Swagger {
		RegisterSwaggerAPI(clientCtx, apiSvr.Router)
	}
}

func (app *Evmos) RegisterTxService(clientCtx client.Context) {
	authtx.RegisterTxService(app.BaseApp.GRPCQueryRouter(), clientCtx, app.BaseApp.Simulate, app.interfaceRegistry)
}

func (app *Evmos) RegisterTendermintService(clientCtx client.Context) {
	tmservice.RegisterTendermintService(
		clientCtx,
		app.BaseApp.GRPCQueryRouter(),
		app.interfaceRegistry,
		app.Query,
	)
}

func (app *Evmos) RegisterNodeService(clientCtx client.Context) {
	node.RegisterNodeService(clientCtx, app.GRPCQueryRouter())
}


func (app *Evmos) GetBaseApp() *baseapp.BaseApp {
	return app.BaseApp
}

func (app *Evmos) GetStakingKeeper() ibctestingtypes.StakingKeeper {
	return app.StakingKeeper
}

func (app *Evmos) GetStakingKeeperSDK() stakingkeeper.Keeper {
	return app.StakingKeeper
}

func (app *Evmos) GetIBCKeeper() *ibckeeper.Keeper {
	return app.IBCKeeper
}

func (app *Evmos) GetScopedIBCKeeper() capabilitykeeper.ScopedKeeper {
	return app.ScopedIBCKeeper
}

func (app *Evmos) GetTxConfig() client.TxConfig {
	cfg := encoding.MakeConfig(ModuleBasics)
	return cfg.TxConfig
}

func RegisterSwaggerAPI(_ client.Context, rtr *mux.Router) {
	statikFS, err := fs.New()
	if err != nil {
		panic(err)
	}

	staticServer := http.FileServer(statikFS)
	rtr.PathPrefix("/swagger/").Handler(http.StripPrefix("/swagger/", staticServer))
}

func GetMaccPerms() map[string][]string {
	dupMaccPerms := make(map[string][]string)
	for k, v := range maccPerms {
		dupMaccPerms[k] = v
	}

	return dupMaccPerms
}

func initParamsKeeper(
	appCodec codec.BinaryCodec, legacyAmino *codec.LegacyAmino, key, tkey storetypes.StoreKey,
) paramskeeper.Keeper {
	paramsKeeper := paramskeeper.NewKeeper(appCodec, legacyAmino, key, tkey)

	
	paramsKeeper.Subspace(authtypes.ModuleName)
	paramsKeeper.Subspace(banktypes.ModuleName)
	paramsKeeper.Subspace(stakingtypes.ModuleName)
	paramsKeeper.Subspace(distrtypes.ModuleName)
	paramsKeeper.Subspace(slashingtypes.ModuleName)
	paramsKeeper.Subspace(govtypes.ModuleName).WithKeyTable(govv1.ParamKeyTable())
	paramsKeeper.Subspace(crisistypes.ModuleName)
	paramsKeeper.Subspace(ibctransfertypes.ModuleName)
	paramsKeeper.Subspace(ibchost.ModuleName)
	
	paramsKeeper.Subspace(evmtypes.ModuleName)
	paramsKeeper.Subspace(feemarkettypes.ModuleName)
	
	paramsKeeper.Subspace(erc20types.ModuleName)
	paramsKeeper.Subspace(claimstypes.ModuleName)
	paramsKeeper.Subspace(recoverytypes.ModuleName)
	paramsKeeper.Subspace(chattypes.ModuleName)
	paramsKeeper.Subspace(gatewaytypes.ModuleName)
	paramsKeeper.Subspace(contracttypes.ModuleName)
	paramsKeeper.Subspace(daotypes.ModuleName)
	paramsKeeper.Subspace(minttypes.ModuleName)
	paramsKeeper.Subspace(group.ModuleName)
	return paramsKeeper
}

func (app *Evmos) setupUpgradeHandlers() {

	app.UpgradeKeeper.SetUpgradeHandler(v2.UpgradeName, v2.CreateUpgradeHandler(app.mm, app.configurator))
	app.UpgradeKeeper.SetUpgradeHandler(v3.UpgradeName, v3.CreateUpgradeHandler(app.mm, app.configurator, app.GatewayKeeper))
	app.UpgradeKeeper.SetUpgradeHandler(v4.UpgradeName, v4.CreateUpgradeHandler(app.mm, app.configurator))
	app.UpgradeKeeper.SetUpgradeHandler(v5.UpgradeName, v5.CreateUpgradeHandler(app.mm, app.configurator))
	app.UpgradeKeeper.SetUpgradeHandler(
		v6.UpgradeName,
		v6.CreateUpgradeHandler(app.mm, app.configurator, app.DaoKeeper, app.BankKeeper),
	)
	app.UpgradeKeeper.SetUpgradeHandler(
		v7.UpgradeName,
		v7.CreateUpgradeHandler(app.mm, app.configurator),
	)
	app.UpgradeKeeper.SetUpgradeHandler(
		v8.UpgradeName,
		v8.CreateUpgradeHandler(app.mm, app.configurator, app.FeeMarketKeeper),
	)
	app.UpgradeKeeper.SetUpgradeHandler(
		v9.UpgradeName,
		v9.CreateUpgradeHandler(app.mm, app.configurator),
	)
	app.UpgradeKeeper.SetUpgradeHandler(
		v10.UpgradeName,
		v10.CreateUpgradeHandler(app.mm, app.configurator),
	)
	app.UpgradeKeeper.SetUpgradeHandler(
		v11.UpgradeName,
		v11.CreateUpgradeHandler(app.mm, app.configurator),
	)
	app.UpgradeKeeper.SetUpgradeHandler(
		v12.UpgradeName,
		v12.CreateUpgradeHandler(app.mm, app.configurator, app.DaoKeeper),
	)
	app.UpgradeKeeper.SetUpgradeHandler(
		v13.UpgradeName,
		v13.CreateUpgradeHandler(app.mm, app.configurator, app.DaoKeeper),
	)
	app.UpgradeKeeper.SetUpgradeHandler(
		v14.UpgradeName,
		v14.CreateUpgradeHandler(app.mm, app.configurator),
	)
	
	
	upgradeInfo, err := app.UpgradeKeeper.ReadUpgradeInfoFromDisk()
	if err != nil {
		panic(fmt.Errorf("failed to read upgrade info from disk: %w", err))
	}

	if app.UpgradeKeeper.IsSkipHeight(upgradeInfo.Height) {
		return
	}
}
