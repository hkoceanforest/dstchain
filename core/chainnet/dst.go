package chainnet

import (
	"bytes"
	"errors"
	"freemasonry.cc/blockchain/app"
	"freemasonry.cc/blockchain/core"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	authType "github.com/cosmos/cosmos-sdk/x/auth/types"
	evmoskr "github.com/evmos/evmos/v10/crypto/keyring"
	"github.com/spf13/pflag"
	rpchttp "github.com/tendermint/tendermint/rpc/client/http"
	"strings"
)

type DstChainNet struct {
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

func (this *DstChainNet) SetApiServerUrl(url string) {
	this.apiServerUrl = url
}

func (this *DstChainNet) SetRpcServerUrl(url string) {
	this.rpcServerUrl = url
}

func (this *DstChainNet) SetJsonRpcServerUrl(url string) {
	this.jsonRpcServerUrl = url
}

func (this *DstChainNet) SetClientCtx(ctx client.Context) {
	this.clientCtx = ctx
}

func (this *DstChainNet) SetChainId(chainid string) {
	this.chainId = chainid
	this.clientCtx = this.clientCtx.WithChainID(chainid)
}

func (this *DstChainNet) GetJsonRpcServerUrl() string {
	return this.jsonRpcServerUrl
}

func (this *DstChainNet) GetClientCtx() client.Context {
	return this.clientCtx
}

func (this *DstChainNet) GetClientFactory() tx.Factory {
	return this.clientFactory
}

func (this *DstChainNet) GetRpcServerUrl() string {
	return this.rpcServerUrl
}

func (this *DstChainNet) GetApiServerUrl() string {
	return strings.TrimSuffix(this.apiServerUrl, "/")
}

func (this *DstChainNet) GetName() string {
	return this.name
}

func (this *DstChainNet) GetCoinType() uint32 {
	return this.coinType
}

func (this *DstChainNet) GetAlgo() string {
	return this.algo
}

func (this *DstChainNet) GetAddressPrefix() string {
	return this.addressPrefix
}

func (this *DstChainNet) GetChainId() string {
	return this.chainId
}

func (this *DstChainNet) GetBaseDenom() string {
	return this.baseDenom
}

func (this *DstChainNet) GetDefaultGasPrice() sdk.DecCoin {
	return this.defaultGasPrice
}

func (this *DstChainNet) AccAddressToBech32(addr sdk.AccAddress) (address string) {
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
func (this *DstChainNet) Bech32ToAccAddress(address string) (addr sdk.AccAddress, err error) {
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

func (this *DstChainNet) SetSdkConfig() {
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(sdk.Bech32PrefixAccAddr, sdk.Bech32PrefixAccPub)
	config.SetBech32PrefixForValidator(sdk.Bech32PrefixValAddr, sdk.Bech32PrefixValPub)
	config.SetBech32PrefixForConsensusNode(sdk.Bech32PrefixConsAddr, sdk.Bech32PrefixConsPub)
	config.SetCoinType(this.coinType)
	config.SetFullFundraiserPath(sdk.FullFundraiserPath)
	config.SetPurpose(sdk.Purpose)
}

func (this *DstChainNet) GetNodeConnectionErr() error {
	return core.DstChainNodeErr
}

func NewDstChainNet(name string, coinType uint32, algo, addressPrefix, chainId, baseDenom, apiServerUrl, rpcServerUrl, jsonRpcServerUrl, defaultGasPrice string) ChainNet {
	log := core.BuildLog(core.GetPackageFuncName(), core.LmChainNet)
	rpcClient, err := rpchttp.New(rpcServerUrl, "/websocket")
	if err != nil {
		log.WithError(err).Error("rpchttp.New")
		panic(err)
	}

	clientCtx := client.Context{}.WithChainID(chainId).
		WithInterfaceRegistry(app.EncodingConfig.InterfaceRegistry).
		WithCodec(app.EncodingConfig.Codec).
		WithTxConfig(app.EncodingConfig.TxConfig).
		WithLegacyAmino(app.EncodingConfig.Amino).
		WithOffline(true).
		WithNodeURI(rpcServerUrl).
		WithClient(rpcClient).
		WithKeyringOptions(evmoskr.Option()).
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

	return &DstChainNet{
		name:             name,
		coinType:         coinType,
		algo:             algo,
		addressPrefix:    addressPrefix,
		chainId:          chainId,
		baseDenom:        baseDenom,
		apiServerUrl:     apiServerUrl,
		rpcServerUrl:     rpcServerUrl,
		jsonRpcServerUrl: jsonRpcServerUrl,
		clientCtx:        clientCtx,
		clientFactory:    clientFactory,
		defaultGasPrice:  defaultGasPriceCoin,
	}
}
