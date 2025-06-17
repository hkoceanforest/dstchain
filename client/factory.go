package client

import (
	"freemasonry.cc/blockchain/app"
	"freemasonry.cc/blockchain/core"
	"freemasonry.cc/blockchain/core/chainnet"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var key = NewSecretKey()

var restClient = NewRestClient()

func NewEvmClient() EvmClient {
	return EvmClient{}
}

func NewRestClient() RestClient {
	return RestClient{"RestClient"}
}

func NewTxClient() TxClient {
	return TxClient{"TxClient"}
}

func NewMempoolClient() MempoolClient {
	return MempoolClient{"TxClient"}
}

func NewBlockClient() BlockClient {
	return BlockClient{"sc-BlockClient"}
}

func NewChatClient(txClient *TxClient, accClient *AccountClient) *ChatClient {
	return &ChatClient{txClient, accClient, "ChatClient"}
}

func NewAccountClient(txClient *TxClient) AccountClient {
	return AccountClient{txClient, key, "AccountClient"}
}

func NewGatewayClinet(txClient *TxClient) GatewayClient {
	acountClient := NewAccountClient(txClient)
	return GatewayClient{txClient, &acountClient, "GatewayClient"}
}

func NewDposClinet(txClient *TxClient) DposClient {
	return DposClient{txClient, "DposClient"}
}

func NewClusterClient(txClient *TxClient) ClusterClient {
	return ClusterClient{
		TxClient:  txClient,
		logPrefix: "ClusterClient",
	}
}

func init() {
	chainnet.Switch(chainnet.ChainNetNameDst)
}


func MsgToStruct(msg sdk.Msg, obj interface{}) error {
	log := core.BuildLog(core.GetPackageFuncName(), core.LmChainClient)
	msgByte, err := app.EncodingConfig.Amino.Marshal(msg)
	if err != nil {
		log.WithError(err).Error("MarshalBinaryBare")
		return err
	}
	err = app.EncodingConfig.Amino.Unmarshal(msgByte, obj)
	if err != nil {
		log.WithError(err).Error("UnmarshalBinaryBare")
		return err
	}
	return nil
}
