package chainnet

import (
	"bytes"
	"errors"
	"fmt"
	"freemasonry.cc/blockchain/app"
	"freemasonry.cc/blockchain/core"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	authType "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/spf13/pflag"
	rpchttp "github.com/tendermint/tendermint/rpc/client/http"
	"strings"
)

type KavaEvmChainNet struct {
	name             string
	coinType         uint32
	algo             string
	addressPrefix    string
	chainId          string
	baseDenom        string
	apiServerUrl     string
	rpcServerUrl     string
	clientCtx        client.Context
	clientFactory    tx.Factory
	jsonRpcServerUrl string
	defaultGasPrice  sdk.DecCoin
}

func (this *KavaEvmChainNet) SetApiServerUrl(url string) {
	this.apiServerUrl = url
}

func (this *KavaEvmChainNet) SetRpcServerUrl(url string) {
	this.rpcServerUrl = url
}

func (this *KavaEvmChainNet) SetJsonRpcServerUrl(url string) {
	this.jsonRpcServerUrl = url
}

func (this *KavaEvmChainNet) SetClientCtx(ctx client.Context) {
	this.clientCtx = ctx
}

func (this *KavaEvmChainNet) SetChainId(chainid string) {
	this.chainId = chainid
}

func (this *KavaEvmChainNet) GetJsonRpcServerUrl() string {
	return this.jsonRpcServerUrl
}

func (this *KavaEvmChainNet) GetClientCtx() client.Context {
	return this.clientCtx
}

func (this *KavaEvmChainNet) GetClientFactory() tx.Factory {
	return this.clientFactory
}

func (this *KavaEvmChainNet) GetRpcServerUrl() string {
	return this.rpcServerUrl
}

func (this *KavaEvmChainNet) GetApiServerUrl() string {
	return strings.TrimSuffix(this.apiServerUrl, "/")
}

func (this *KavaEvmChainNet) GetName() string {
	return this.name
}

func (this *KavaEvmChainNet) GetCoinType() uint32 {
	return this.coinType
}

func (this *KavaEvmChainNet) GetAlgo() string {
	return this.algo
}

func (this *KavaEvmChainNet) GetAddressPrefix() string {
	return this.addressPrefix
}

func (this *KavaEvmChainNet) GetDefaultGasPrice() sdk.DecCoin {
	return this.defaultGasPrice
}

func (this *KavaEvmChainNet) GetChainId() string {
	return this.chainId
}

func (this *KavaEvmChainNet) GetBaseDenom() string {
	return this.baseDenom
}


func (this *KavaEvmChainNet) AccAddressToBech32(addr sdk.AccAddress) (address string) {
	if addr.Empty() {
		return ""
	}
	if bytes.Equal(addr, []byte(sdk.BlackHoleAddress)) {
		return this.addressPrefix + (sdk.BlackHoleAddress + sdk.BlackHoleAddress[1:])
	}
	bech32Addr, err := bech32.ConvertAndEncode(this.addressPrefix, addr)
	if err != nil {
		panic(err)
	}
	return bech32Addr
}

func (this *KavaEvmChainNet) Bech32ToAccAddress(address string) (addr sdk.AccAddress, err error) {
	if len(strings.TrimSpace(address)) == 0 {
		return sdk.AccAddress{}, errors.New("empty address string is not allowed")
	}
	bech32PrefixAccAddr := this.addressPrefix
	bz, err := sdk.GetFromBech32(address, bech32PrefixAccAddr)
	if err != nil {
		return nil, err
	}
	err = sdk.VerifyAddressFormat(bz)
	if err != nil {
		return nil, err
	}
	return bz, nil
}

func (this *KavaEvmChainNet) SetSdkConfig() {
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(this.addressPrefix, "kavapub")
	config.SetBech32PrefixForValidator("kavavaloper", "kavavaloperpub")
	config.SetBech32PrefixForConsensusNode("kavavalcons", "kavavalconspub")
	config.SetCoinType(this.coinType)
	config.SetFullFundraiserPath(fmt.Sprintf("m/%d'/%d'/0'/0/0", sdk.Purpose, this.coinType))
	config.SetPurpose(sdk.Purpose)
}

func (this *KavaEvmChainNet) GetNodeConnectionErr() error {
	return core.KavaChainNodeErr
}
func NewKavaEvmChainNet(name string, coinType uint32, algo, addressPrefix, chainId, baseDenom, apiServerUrl, rpcServerUrl, jsonRpcServerUrl, defaultGasPrice string) ChainNet {
	log := core.BuildLog(core.GetPackageFuncName(), core.LmChainNet)
	rpcClient, err := rpchttp.New(rpcServerUrl, "/websocket")
	if err != nil {
		log.WithError(err).Error("rpchttp.New")
		panic(err)
	}

	clientCtx := client.Context{}.WithChainID(chainId).
		WithCodec(app.EncodingConfig.Codec).
		WithTxConfig(app.EncodingConfig.TxConfig).
		WithLegacyAmino(app.EncodingConfig.Amino).
		WithOffline(true).
		WithNodeURI(rpcServerUrl).
		WithClient(rpcClient).
		WithAccountRetriever(authType.AccountRetriever{})

	flags := pflag.NewFlagSet("chat", pflag.ContinueOnError)

	flagErrorBuf := new(bytes.Buffer)

	flags.SetOutput(flagErrorBuf)

	
	clientFactory := tx.NewFactoryCLI(clientCtx, flags).WithChainID(chainId).
		WithAccountRetriever(clientCtx.AccountRetriever).
		WithTxConfig(clientCtx.TxConfig)

	defaultGasPriceCoin, err := sdk.ParseDecCoin(defaultGasPrice)
	if err != nil {
		log.WithError(err).Error("sdk.ParseDecCoin")
		panic(err)
	}

	return &KavaEvmChainNet{
		name:             name,
		coinType:         coinType,
		algo:             algo,
		addressPrefix:    addressPrefix,
		chainId:          chainId,
		baseDenom:        baseDenom,
		apiServerUrl:     apiServerUrl,
		rpcServerUrl:     rpcServerUrl,
		jsonRpcServerUrl: jsonRpcServerUrl,
		clientFactory:    clientFactory,
		clientCtx:        clientCtx,
		defaultGasPrice:  defaultGasPriceCoin,
	}
}
