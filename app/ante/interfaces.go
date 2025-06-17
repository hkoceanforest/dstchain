package ante

import (
	bankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/evmos/ethermint/x/evm/statedb"
	erc20types "github.com/evmos/evmos/v10/x/erc20/types"
	"math/big"

	daoTypes "freemasonry.cc/blockchain/x/dao/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/params"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	evm "github.com/evmos/ethermint/x/evm/vm"
)

type EvmKeeper interface {
	GetParams(ctx sdk.Context) (params evmtypes.Params)
	ChainID() *big.Int
	GetBaseFee(ctx sdk.Context, ethCfg *params.ChainConfig) *big.Int
}
type DaoKeeper interface {
	CalculateFee(ctx sdk.Context, msg sdk.Msg, fee sdk.Coins) (sdk.Coins, error)
	CalculateSendFee(ctx sdk.Context, msg bankTypes.MsgSend, fee sdk.Coins) (sdk.Coins, error)
	FreeTxMsg(ctx sdk.Context, msg sdk.Msg) bool
	CalculateEvmFee(ctx sdk.Context, tx sdk.Tx, toAddress common.Address, fee sdk.Coins) (sdk.Coins, error)
	GetClusterByChatId(ctx sdk.Context, cluserChatId string) (res daoTypes.DeviceCluster, err error)
	GetCluster(ctx sdk.Context, cluserId string) (res daoTypes.DeviceCluster, err error)
	GetRidClusterId(ctx sdk.Context, redpacketId string) (string, error)
	GetParams(ctx sdk.Context) (params daoTypes.Params)
	DeductionToken(ctx sdk.Context, contractAddress, sender string) error
}

type ContractsKeeper interface {
	GetRedPacketContractAddress(ctx sdk.Context) string
	GetTokenPairID(ctx sdk.Context, token string) []byte
	GetTokenPair(ctx sdk.Context, id []byte) (erc20types.TokenPair, bool)
}

type DynamicFeeEVMKeeper interface {
	ChainID() *big.Int
	GetParams(ctx sdk.Context) evmtypes.Params
	GetBaseFee(ctx sdk.Context, ethCfg *params.ChainConfig) *big.Int
}

type EVMKeeper interface {
	statedb.Keeper
	DynamicFeeEVMKeeper

	NewEVM(ctx sdk.Context, msg core.Message, cfg *evmtypes.EVMConfig, tracer vm.EVMLogger, stateDB vm.StateDB) evm.EVM
	DeductTxCostsFromUserBalance(ctx sdk.Context, fees sdk.Coins, from common.Address) error
	GetBalance(ctx sdk.Context, addr common.Address) *big.Int
	ResetTransientGasUsed(ctx sdk.Context)
	GetTxIndexTransient(ctx sdk.Context) uint64
	GetChainConfig(ctx sdk.Context) evmtypes.ChainConfig
	GetEVMDenom(ctx sdk.Context) string
	GetEnableCreate(ctx sdk.Context) bool
	GetEnableCall(ctx sdk.Context) bool
	GetAllowUnprotectedTxs(ctx sdk.Context) bool
}
