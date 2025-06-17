package chainnet

import (
	"freemasonry.cc/blockchain/core"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/evmos/ethermint/crypto/ethsecp256k1"
)

const (
	
	ChainNetNameDst     = "dst"
	ChainIdDst          = "daodst_7777-3"
	BaseDenomDst        = core.BaseDenom
	GovDenomDst         = core.GovDenom
	AddressPrefixDst    = sdk.Bech32MainPrefix
	CoinTypeDst         = uint32(60)
	AlgoDst             = ethsecp256k1.KeyType
	ApiServerUrlDst     = "http://18.167.8.64:1317"
	RpcServerUrlDst     = "http://18.167.8.64:26657"
	JsonRpcServerUrlDst = "http://18.167.8.64:8545"
	DefaultGasPriceDst  = "3000000000dst" 

	
	ChainNetNameKava     = "kava"
	ChainIdKava          = "kava_2222-10"
	BaseDenomKava        = "kava"
	GovDenomKava         = "kava"
	AddressPrefixKava    = "kava"
	CoinTypeKava         = uint32(459)
	AlgoKava             = string(hd.Secp256k1Type)
	ApiServerUrlKava     = "https://api.data.kava.io"
	RpcServerUrlKava     = "https://rpc-kava-01.stakeflow.io:443"
	JsonRpcServerUrlKAva = ""
	DefaultGasPriceKava  = "0.0015ukava" 

	
	ChainNetNameKavaEvm = "kava_evm"
	AlgoKavaEvm         = ethsecp256k1.KeyType
	CoinTypeKavaEvm     = uint32(60)
)

var (
	ChainNetKava    = NewKavaChainNet(ChainNetNameKava, CoinTypeKava, AlgoKava, AddressPrefixKava, ChainIdKava, BaseDenomKava, ApiServerUrlKava, RpcServerUrlKava, JsonRpcServerUrlKAva, DefaultGasPriceKava)
	ChainNetKavaEvm = NewKavaEvmChainNet(ChainNetNameKavaEvm, CoinTypeKavaEvm, AlgoKavaEvm, AddressPrefixKava, ChainIdKava, BaseDenomKava, ApiServerUrlKava, RpcServerUrlKava, JsonRpcServerUrlKAva, DefaultGasPriceKava)
	ChainNetDst     = NewDstChainNet(ChainNetNameDst, CoinTypeDst, AlgoDst, AddressPrefixDst, ChainIdDst, BaseDenomDst, ApiServerUrlDst, RpcServerUrlDst, JsonRpcServerUrlDst, DefaultGasPriceDst)
)

func init() {
	register(ChainNetKava)
	register(ChainNetKavaEvm)
	register(ChainNetDst)
	Switch(ChainNetNameDst)
}
