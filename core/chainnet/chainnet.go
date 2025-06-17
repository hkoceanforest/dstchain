package chainnet

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"sync"
)

var _ ChainNet = &DstChainNet{}
var _ ChainNet = &KavaChainNet{}
var chainnets = sync.Map{}

var Current ChainNet


func register(net ChainNet) {
	chainnets.Store(net.GetName(), net)
}


func Switch(netname string) {
	_chainNet, ok := chainnets.Load(netname)
	if ok {
		Current = _chainNet.(ChainNet)
	} else {
		panic(fmt.Sprintf("chainnet %s not register", netname))
	}
	Current.SetSdkConfig() 
}

type ChainNet interface {
	SetClientCtx(ctx client.Context)
	GetClientCtx() client.Context
	GetClientFactory() tx.Factory
	GetName() string
	GetCoinType() uint32
	GetAlgo() string
	GetAddressPrefix() string       
	GetApiServerUrl() string        
	SetApiServerUrl(url string)     
	GetRpcServerUrl() string        
	SetRpcServerUrl(url string)     
	GetJsonRpcServerUrl() string    
	SetJsonRpcServerUrl(url string) 
	GetChainId() string
	SetChainId(chainid string)       
	GetBaseDenom() string            
	GetDefaultGasPrice() sdk.DecCoin 
	AccAddressToBech32(addr sdk.AccAddress) (address string)
	Bech32ToAccAddress(address string) (addr sdk.AccAddress, err error)
	SetSdkConfig()               
	GetNodeConnectionErr() error 
}
