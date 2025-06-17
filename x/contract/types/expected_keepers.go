package types

import (
	"context"
	chatTypes "freemasonry.cc/blockchain/x/chat/types"
	types2 "freemasonry.cc/blockchain/x/gateway/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/evmos/ethermint/x/evm/statedb"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
)

type EVMKeeper interface {
	GetParams(ctx sdk.Context) evmtypes.Params
	GetAccountWithoutBalance(ctx sdk.Context, addr common.Address) *statedb.Account
	GetCode(ctx sdk.Context, codeHash common.Hash) []byte
	EstimateGas(c context.Context, req *evmtypes.EthCallRequest) (*evmtypes.EstimateGasResponse, error)
	ApplyMessage(ctx sdk.Context, msg core.Message, tracer vm.EVMLogger, commit bool) (*evmtypes.MsgEthereumTxResponse, error)
	
}

type ChatKeeper interface {
	SetGatewayIssueToken(ctx sdk.Context, gatewayAddress string, tokenInfo GatewayTokenInfo) error
	GetGatewayIssueToken(ctx sdk.Context, gatewayAddress string) (*GatewayTokenInfo, error)
	GetRegisterInfo(ctx sdk.Context, fromAddress string) (userinfo chatTypes.UserInfo, err error)
	SetRegisterInfo(ctx sdk.Context, userInfo chatTypes.UserInfo) error
	RegisterMobile(ctx sdk.Context, nodeAddress, fromAddress, mobilePrefix string) (mobile string, err error)
}

type GatewayKeeper interface {
	GetGatewayInfo(ctx sdk.Context, gatewayAddress string) (*types2.Gateway, error)
}
