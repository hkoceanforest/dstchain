package testenv

import (
	sdkmath "cosmossdk.io/math"
	"fmt"
	"freemasonry.cc/blockchain/app"
	"freemasonry.cc/blockchain/contracts"
	"freemasonry.cc/blockchain/core"
	chatkeeper "freemasonry.cc/blockchain/x/chat/keeper"
	chattypes "freemasonry.cc/blockchain/x/chat/types"
	contractkeeper "freemasonry.cc/blockchain/x/contract/keeper"
	daokeeper "freemasonry.cc/blockchain/x/dao/keeper"
	daoTypes "freemasonry.cc/blockchain/x/dao/types"
	gatewaykeeper "freemasonry.cc/blockchain/x/gateway/keeper"
	gatewaytypes "freemasonry.cc/blockchain/x/gateway/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distributionkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	groupkeeper "github.com/cosmos/cosmos-sdk/x/group/keeper"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/evmos/ethermint/encoding"
	feemarkettypes "github.com/evmos/ethermint/x/feemarket/types"
	"github.com/stretchr/testify/suite"
	abci "github.com/tendermint/tendermint/abci/types"
	ce "github.com/tendermint/tendermint/crypto/encoding"
	"github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	sm "github.com/tendermint/tendermint/state"
	dbm "github.com/tendermint/tm-db"
	"math/big"
	"time"
)

const (
	genesisAccountKey = "genesis"
	moduleAccountName = "chat"
	moduleAccountDao  = "dao"
)

var (
	chatModuleAccount *authtypes.ModuleAccount
	genesisAccount    sdk.AccAddress
	acc1              authtypes.AccountI
	
	genesisAmountDst = sdk.MustNewDecFromStr("57700000000000000000000000").TruncateInt() 
	genesisAmountNxn = sdk.MustNewDecFromStr("2100000000000000000000000").TruncateInt()  
	
	genesisAccountTokens = sdk.NewCoins(
		sdk.NewCoin(core.BaseDenom, sdk.MustNewDecFromStr("5770000000000000000000000000").TruncateInt()),
		sdk.NewCoin(core.GovDenom, sdk.MustNewDecFromStr("2100000000000000000000000").TruncateInt()),
	)

	priv1 = secp256k1.GenPrivKey()
	pk1   = priv1.PubKey()

	defaultBurnAmount, _ = sdk.NewDecFromStr("50000000000000000000000")

	
	priv_dalidator_key = `{
  "address": "58B0E2F8F6AA7815C5F9BECF0467E30F65C70D63",
  "pub_key": {
    "type": "tendermint/PubKeyEd25519",
    "value": "i7reVVcmfiL4GwB+x2go6hH0HD7jySfdxDMjlL+xcg4="
  },
  "priv_key": {
    "type": "tendermint/PrivKeyEd25519",
    "value": "LaVxXL7R6Unyib788SwbVj1dVya8tedRm6UlE+tBRoGLut5VVyZ+IvgbAH7HaCjqEfQcPuPJJ93EMyOUv7FyDg=="
  }
}`
)

type IntegrationTestSuite struct {
	suite.Suite
	Env        *TestAppEnv
	bankKep    bankkeeper.Keeper
	accountKep *authkeeper.AccountKeeper

	gatewayKep  *gatewaykeeper.Keeper
	gatewayServ gatewaytypes.MsgServer

	stakingKep  *stakingkeeper.Keeper
	stakingServ stakingtypes.MsgServer

	chatKep  *chatkeeper.Keeper
	chatServ chattypes.MsgServer

	daoServ  daoTypes.MsgServer
	daoKep   *daokeeper.Keeper
	distrKep *distributionkeeper.Keeper
	groupKep *groupkeeper.Keeper

	contractKep *contractkeeper.Keeper
}


func (this *IntegrationTestSuite) SetupTest() {
	this.T().Log("SetupTest()")
	
	
	
	env, err := InitTestEnv()
	if err != nil {
		this.T().Fatal(err)
	}
	this.Env = env

	this.gatewayKep = &env.App.GatewayKeeper

	this.gatewayServ = gatewaykeeper.NewMsgServerImpl(*this.gatewayKep)

	this.chatKep = &env.App.ChatKeeper

	this.chatServ = chatkeeper.NewMsgServerImpl(*this.chatKep)

	this.bankKep = env.App.BankKeeper
	this.accountKep = &env.App.AccountKeeper

	this.stakingKep = &env.App.StakingKeeper
	this.stakingServ = stakingkeeper.NewMsgServerImpl(*this.stakingKep)

	this.distrKep = &env.App.DistrKeeper

	this.daoKep = &env.App.DaoKeeper

	this.daoServ = daokeeper.NewMsgServerImpl(*this.daoKep)

	this.groupKep = &env.App.GroupKeeper

	this.contractKep = &env.App.ContractKeeper
}

type TestAppEnv struct {
	App    *app.Evmos
	Ctx    sdk.Context
	Height int64
}

func MakeAppAndGenesis(genDocFile string) *app.Evmos {
	if genDocFile == "" {
		feemarketTypes := feemarkettypes.GenesisState{
			Params: feemarkettypes.Params{
				NoBaseFee:                false,
				BaseFeeChangeDenominator: 8,
				ElasticityMultiplier:     2,
				EnableHeight:             0,
				BaseFee:                  sdk.NewInt(1000000000),
				MinGasMultiplier:         sdk.NewDecWithPrec(50, 2),
				MinGasPrice:              sdk.ZeroDec(),
			},
			BlockGas: 0,
		}
		return app.Setup(false, &feemarketTypes)
	}
	db := dbm.NewMemDB()
	app := app.NewEvmos(log.NewNopLogger(), db, nil, true, map[int64]bool{}, app.DefaultNodeHome, 5, encoding.MakeConfig(app.ModuleBasics), simapp.EmptyAppOptions{})

	genDoc, err := sm.MakeGenesisDocFromFile(genDocFile)
	if err != nil {
		panic("MakeGenesisDocFromFile:" + err.Error())
	}

	req := abci.RequestInitChain{
		Time:    genDoc.GenesisTime,
		ChainId: genDoc.ChainID,
		ConsensusParams: &abci.ConsensusParams{
			Block: &abci.BlockParams{
				MaxBytes: genDoc.ConsensusParams.Block.MaxBytes,
				MaxGas:   genDoc.ConsensusParams.Block.MaxGas,
			},
			Evidence: &tmproto.EvidenceParams{
				MaxAgeNumBlocks: genDoc.ConsensusParams.Evidence.MaxAgeNumBlocks,
				MaxAgeDuration:  genDoc.ConsensusParams.Evidence.MaxAgeDuration,
				MaxBytes:        genDoc.ConsensusParams.Evidence.MaxBytes,
			},
			Validator: &tmproto.ValidatorParams{
				PubKeyTypes: genDoc.ConsensusParams.Validator.PubKeyTypes,
			},
			Version: &tmproto.VersionParams{
				AppVersion: genDoc.ConsensusParams.Version.AppVersion,
			},
		},
		Validators:    []abci.ValidatorUpdate{},
		AppStateBytes: genDoc.AppState,
		InitialHeight: genDoc.InitialHeight,
	}

	for _, v := range genDoc.Validators {
		publickey, err := ce.PubKeyToProto(v.PubKey)
		if err != nil {
			panic(err)
		}
		req.Validators = append(req.Validators, abci.ValidatorUpdate{
			PubKey: publickey,
			Power:  v.Power,
		})
	}

	res := app.InitChain(req)
	fmt.Println("InitChainer response :", res)

	return app
}



func InitTestEnv() (*TestAppEnv, error) {

	
	chatApp := MakeAppAndGenesis("")

	
	

	ctx := chatApp.BaseApp.NewContext(false, tmproto.Header{Height: 1})

	
	

	genesisAccount = authtypes.NewModuleAddress(genesisAccountKey) 
	chatModuleAccount = authtypes.NewEmptyModuleAccount(moduleAccountName, authtypes.Minter, authtypes.Burner)
	env := &TestAppEnv{
		Ctx: ctx,
		App: chatApp,
	}

	chatApp.AccountKeeper.SetParams(ctx, authtypes.DefaultParams())
	chatApp.BankKeeper.SetParams(ctx, banktypes.DefaultParams())
	chatApp.StakingKeeper.SetParams(ctx, stakingtypes.DefaultParams())

	chatApp.AccountKeeper.SetModuleAccount(ctx, chatModuleAccount)

	
	err := chatApp.BankKeeper.MintCoins(ctx, moduleAccountName, genesisAccountTokens)
	if err != nil {
		return env, err
	}
	err = chatApp.BankKeeper.SendCoinsFromModuleToAccount(ctx, moduleAccountName, genesisAccount, genesisAccountTokens)
	if err != nil {
		return env, err
	}
	
	balance0 := chatApp.BankKeeper.GetAllBalances(ctx, genesisAccount)
	fmt.Println("Genesis Account:", genesisAccount.String(), " Balance:", balance0)

	priv := secp256k1.GenPrivKey()
	pk := priv.PubKey()
	acc1 = chatApp.AccountKeeper.NewAccountWithAddress(ctx, sdk.AccAddress(pk.Address()))
	chatApp.AccountKeeper.SetAccount(ctx, acc1) 
	return env, nil
}

func (this *IntegrationTestSuite) Commit() {
	_ = this.Env.App.Commit()
	header := this.Env.Ctx.BlockHeader()
	header.Height += 1
	header.Time = header.Time.Add(time.Nanosecond)
	this.Env.App.BeginBlock(abci.RequestBeginBlock{
		Header: header,
	})
	
	this.Env.Ctx = this.Env.App.BaseApp.NewContext(false, header)
}

func (this *IntegrationTestSuite) GetErc20Balance(ctx sdk.Context, contractAddress string, address sdk.AccAddress) (sdk.Int, error) {
	contractAddr := common.HexToAddress(contractAddress)
	addr := common.BytesToAddress(address.Bytes())
	erc20Abi := contracts.UsdtContract.ABI
	resp, err := this.contractKep.CallEVM(ctx, erc20Abi, addr, contractAddr, false, "balanceOf", addr)
	if err != nil {
		return sdkmath.ZeroInt(), err
	}
	unPackData, err := erc20Abi.Unpack("balanceOf", resp.Ret)
	if err != nil {
		return sdkmath.ZeroInt(), err
	}
	balance := unPackData[0].(*big.Int)
	return sdk.NewIntFromBigInt(balance), nil
}
func (this *IntegrationTestSuite) Erc20Transfer(ctx sdk.Context, contractAddress string) error {
	contractAddr := common.HexToAddress(contractAddress)
	erc20RewardAccount := common.BytesToAddress(core.Erc20Reward.Bytes())
	fromAccount := common.HexToAddress("0x833DbecaECa9E771633e9B3e9270293bd1ab9FD1")
	erc20Abi := contracts.UsdtContract.ABI
	amount, _ := sdkmath.NewIntFromString("1000000000000000000000")
	resp, err := this.contractKep.CallEVM(ctx, erc20Abi, fromAccount, contractAddr, true, "transfer", erc20RewardAccount, amount.BigInt())
	if err != nil {
		return err
	}
	if resp.Failed() {
		return err
	}
	return nil
}
