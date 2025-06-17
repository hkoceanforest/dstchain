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

type KavaChainNet struct {
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

func (this *KavaChainNet) SetApiServerUrl(url string) {
	this.apiServerUrl = url
}

func (this *KavaChainNet) SetRpcServerUrl(url string) {
	this.rpcServerUrl = url
}

func (this *KavaChainNet) SetJsonRpcServerUrl(url string) {
	this.jsonRpcServerUrl = url
}

func (this *KavaChainNet) SetClientCtx(ctx client.Context) {
	this.clientCtx = ctx
}

func (this *KavaChainNet) SetChainId(chainid string) {
	this.chainId = chainid
}

func (this *KavaChainNet) GetJsonRpcServerUrl() string {
	return this.jsonRpcServerUrl
}

func (this *KavaChainNet) GetClientCtx() client.Context {
	return this.clientCtx
}

func (this *KavaChainNet) GetClientFactory() tx.Factory {
	return this.clientFactory
}

func (this *KavaChainNet) GetRpcServerUrl() string {
	return this.rpcServerUrl
}

func (this *KavaChainNet) GetApiServerUrl() string {
	return strings.TrimSuffix(this.apiServerUrl, "/")
}

func (this *KavaChainNet) GetName() string {
	return this.name
}

func (this *KavaChainNet) GetCoinType() uint32 {
	return this.coinType
}

func (this *KavaChainNet) GetAlgo() string {
	return this.algo
}

func (this *KavaChainNet) GetAddressPrefix() string {
	return this.addressPrefix
}

func (this *KavaChainNet) GetDefaultGasPrice() sdk.DecCoin {
	return this.defaultGasPrice
}

func (this *KavaChainNet) GetChainId() string {
	return this.chainId
}

func (this *KavaChainNet) GetBaseDenom() string {
	return this.baseDenom
}


func (this *KavaChainNet) AccAddressToBech32(addr sdk.AccAddress) (address string) {
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

func (this *KavaChainNet) Bech32ToAccAddress(address string) (addr sdk.AccAddress, err error) {
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

func (this *KavaChainNet) SetSdkConfig() {
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(this.addressPrefix, "kavapub")
	config.SetBech32PrefixForValidator("kavavaloper", "kavavaloperpub")
	config.SetBech32PrefixForConsensusNode("kavavalcons", "kavavalconspub")
	config.SetCoinType(this.coinType)
	config.SetFullFundraiserPath(fmt.Sprintf("m/%d'/%d'/0'/0/0", sdk.Purpose, this.coinType))
	config.SetPurpose(sdk.Purpose)
}

func (this *KavaChainNet) GetNodeConnectionErr() error {
	return core.KavaChainNodeErr
}

func NewKavaChainNet(name string, coinType uint32, algo, addressPrefix, chainId, baseDenom, apiServerUrl, rpcServerUrl, jsonRpcServerUrl, defaultGasPrice string) ChainNet {
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

	return &KavaChainNet{
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
