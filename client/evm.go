package client

import (
	"bytes"
	"context"
	sdkmath "cosmossdk.io/math"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"freemasonry.cc/blockchain/client/evm"
	"freemasonry.cc/blockchain/core"
	"freemasonry.cc/blockchain/core/chainnet"
	"freemasonry.cc/blockchain/util"
	"freemasonry.cc/blockchain/x/contract/types"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common/hexutil"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	erc20types "github.com/evmos/evmos/v10/x/erc20/types"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net/http"
	"time"
)

func (this *EvmClient) GetTokenPair(token string) (erc20types.TokenPair, error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	clientCtx := chainnet.ChainNetDst.GetClientCtx()
	tokenPair := erc20types.TokenPair{}
	params := erc20types.QueryTokenPairRequest{Token: token}
	bz, err := clientCtx.LegacyAmino.MarshalJSON(params)
	if err != nil {
		log.WithError(err).Error("MarshalJSON")
		return tokenPair, err
	}
	resBytes, _, err := util.QueryWithDataWithUnwrapErr(clientCtx, "custom/contract/"+types.QueryTokenPair, bz)
	if err != nil {
		log.WithError(err).Error("QueryWithData")
		return tokenPair, err
	}
	if resBytes != nil {
		err = clientCtx.LegacyAmino.UnmarshalJSON(resBytes, &tokenPair)
		if err != nil {
			log.WithError(err).Error("UnmarshalJSON")
			return tokenPair, err
		}
	}
	return tokenPair, nil
}

func (this *EvmClient) GetContractCode(contract string) (string, error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	params := evmtypes.QueryCodeRequest{Address: contract}
	clientCtx := chainnet.ChainNetDst.GetClientCtx()
	bz, err := clientCtx.LegacyAmino.MarshalJSON(params)
	if err != nil {
		log.WithError(err).Error("MarshalJSON")
		return "", err
	}
	resBytes, _, err := util.QueryWithDataWithUnwrapErr(clientCtx, "custom/contract/"+types.QueryContractCode, bz)
	if err != nil {
		log.WithError(err).Error("QueryWithData")
		return "", err
	}
	if resBytes != nil {
		return hexutil.Encode(resBytes), nil
	}
	return "", nil
}

func (this *EvmClient) TraceTx(tx string) (string, error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	clientCtx := chainnet.ChainNetDst.GetClientCtx()
	data2, _ := base64.StdEncoding.DecodeString("cywylwAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA3gtrOnZAAAAAAAAAAAAAAAAAAAkKeABU83Luxw3vPxyWwbvNOiEzY=")
	r, _ := base64.StdEncoding.DecodeString("KCnJEYNfgQxKmpLd5KaZ/uSQFWCf0WoybRjNmU4xrEE=")
	s, _ := base64.StdEncoding.DecodeString("GIIxWznN/1WMbfmOMmhbQw7A8JhxUo78h77ky4S73/8=")
	v, _ := base64.StdEncoding.DecodeString("AQ==")

	ChainID := sdk.NewInt(7777)
	GasTipCap := sdkmath.NewInt(1500000000)
	GasFeeCap := sdkmath.NewInt(1500000008)
	Amount := sdkmath.NewInt(0)
	feetx := codectypes.UnsafePackAny(&evmtypes.DynamicFeeTx{
		ChainID:   &ChainID,
		Nonce:     uint64(12),
		GasTipCap: &GasTipCap,
		GasFeeCap: &GasFeeCap,
		GasLimit:  uint64(58915),
		To:        "0x2DE216cd9D6684870FEa8944E3967f09639733a9",
		Amount:    &Amount,
		Data:      data2,
		Accesses:  nil,
		V:         v,
		R:         r,
		S:         s,
	})
	req := evmtypes.QueryTraceTxRequest{
		Msg: &evmtypes.MsgEthereumTx{
			Data: feetx, Size_: 0, Hash: "0x8fe52b42a953dc7fb8d71171006a4194455d00785716cb60703f3475e7de3c11", From: "",
		},
		TraceConfig: &evmtypes.TraceConfig{
			
			DisableStack:     false,
			DisableStorage:   false,
			Reexec:           99,
			Debug:            true,
			EnableMemory:     true,
			EnableReturnData: true,
			Limit:            0,
		},
	}
	conn, err := grpc.Dial(
		"127.0.0.1:9090",
		grpc.WithInsecure(), 
		grpc.WithDefaultCallOptions(grpc.ForceCodec(codec.NewProtoCodec(clientCtx.InterfaceRegistry).GRPCCodec())),
	)
	if err != nil {
		log.WithError(err).Error("grpc.Dial")
		return "", err
	}

	newCtx, _ := context.WithTimeout(context.Background(), time.Second*5)
	queryClient := evmtypes.NewQueryClient(conn)
	resp, err := queryClient.TraceTx(newCtx, &req)
	if err != nil {
		log.WithError(err).Error("TraceBlock")
		return "", err
	}
	if resp != nil {
		return resp.String(), nil
	}
	return "", nil
}


func (this *EvmClient) NetVersion() (string, error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	var res string

	rpcRes, err := this.Call("net_version", []string{})
	if err != nil {
		log.WithError(err).Error("call")
		return res, err
	}

	err = json.Unmarshal(rpcRes.Result, &res)
	if err != nil {
		log.WithField("result", rpcRes.Result).WithError(err).Error("Unmarshal")
		return res, err
	}
	return res, nil
}


func (this *EvmClient) NetListening() (bool, error) {
	var res bool
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	rpcRes, err := this.Call("net_listening", []string{})
	if err != nil {
		log.WithError(err).Error("call")
		return res, err
	}

	err = json.Unmarshal(rpcRes.Result, &res)
	if err != nil {
		log.WithField("result", rpcRes.Result).WithError(err).Error("Unmarshal")
		return res, err
	}
	return res, err
}


func (this *EvmClient) NetPeerCount() (int, error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)

	var res int
	rpcRes, err := this.Call("net_peerCount", []string{})
	if err != nil {
		log.WithError(err).Error("call")
		return res, err
	}

	err = json.Unmarshal(rpcRes.Result, &res)
	if err != nil {
		log.WithField("result", rpcRes.Result).WithError(err).Error("Unmarshal")
		return res, err
	}
	return res, err
}




func (this *EvmClient) GetBlockNumber(blockNumber string, fullTx bool) (map[string]interface{}, error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)

	var res map[string]interface{}
	rpcRes, err := this.Call("eth_getBlockByNumber", []interface{}{blockNumber, true})
	if err != nil {
		log.WithError(err).Error("call")
		return res, err
	}

	err = json.Unmarshal(rpcRes.Result, &res)
	if err != nil {
		log.WithField("result", rpcRes.Result).WithError(err).Error("Unmarshal")
		return res, err
	}
	return res, err
}


func (this *EvmClient) BlockNumber() (hexutil.Big, error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)

	var res hexutil.Big
	rpcRes, err := this.Call("eth_blockNumber", []interface{}{})
	if err != nil {
		log.WithError(err).Error("call")
		return res, err
	}
	err = json.Unmarshal(rpcRes.Result, &res)
	if err != nil {
		log.WithField("result", rpcRes.Result).WithError(err).Error("Unmarshal")
		return res, err
	}
	return res, err
}


func (this *EvmClient) GetBalance(addr string) (hexutil.Big, error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)

	var res hexutil.Big
	rpcRes, err := this.Call("eth_getBalance", []string{addr, "latest"})
	if err != nil {
		log.WithField("address", addr).WithError(err).Error("call")
		return res, err
	}
	if rpcRes.Error != nil {
		log.WithFields(logrus.Fields{"code": rpcRes.Error.Code, "message": rpcRes.Error.Message, "data": rpcRes.Error.Data}).Error("rpcError")
		return res, errors.New(rpcRes.Error.Message)
	}

	err = res.UnmarshalJSON(rpcRes.Result)
	if err != nil {
		log.WithField("result", rpcRes.Result).WithError(err).Error("UnmarshalJSON")
		return res, err
	}
	return res, nil
}

func (this *EvmClient) GetTransactionReceipt(hash string) (map[string]interface{}, error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	var res map[string]interface{}
	rpcRes, err := this.Call("eth_getTransactionReceipt", []interface{}{hash})
	if err != nil {
		log.WithField("hash", hash).WithError(err).Error("call")
		return res, err
	}
	if rpcRes.Error != nil {
		log.WithFields(logrus.Fields{"code": rpcRes.Error.Code, "message": rpcRes.Error.Message, "data": rpcRes.Error.Data}).Error("rpcError")
		return res, errors.New(rpcRes.Error.Message)
	}
	err = json.Unmarshal(rpcRes.Result, &res)
	if err != nil {
		log.WithField("result", rpcRes.Result).WithError(err).Error("Unmarshal")
		return res, err
	}
	return res, err
}

func (this *EvmClient) GetTransactionByHash(hash string) (map[string]interface{}, error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	var res map[string]interface{}
	rpcRes, err := this.Call("eth_getTransactionByHash", []interface{}{hash})
	if err != nil {
		log.WithField("hash", hash).WithError(err).Error("call")
		return res, err
	}
	if rpcRes.Error != nil {
		log.WithFields(logrus.Fields{"code": rpcRes.Error.Code, "message": rpcRes.Error.Message, "data": rpcRes.Error.Data}).Error("rpcError")
		return res, errors.New(rpcRes.Error.Message)
	}
	err = json.Unmarshal(rpcRes.Result, &res)
	if err != nil {
		log.WithField("result", rpcRes.Result).WithError(err).Error("Unmarshal")
		return res, err
	}
	return res, err
}


func (this *EvmClient) GetAddress() ([]hexutil.Bytes, error) {
	rpcRes, err := this.CallWithError("eth_accounts", []string{})
	if err != nil {
		return nil, err
	}
	var res []hexutil.Bytes
	err = json.Unmarshal(rpcRes.Result, &res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (this *EvmClient) CreateRequest(method string, params interface{}) evm.Request {
	return evm.Request{
		Version: "2.0",
		Method:  method,
		Params:  params,
		ID:      1,
	}
}

func (this *EvmClient) CallWithError(method string, params interface{}) (*evm.Response, error) {
	req, err := json.Marshal(this.CreateRequest(method, params))
	if err != nil {
		return nil, err
	}

	var rpcRes *evm.Response
	time.Sleep(1 * time.Second)

	httpReq, err := http.NewRequestWithContext(context.Background(), "POST", chainnet.Current.GetJsonRpcServerUrl(), bytes.NewBuffer(req))
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	res, err := client.Do(httpReq)
	if err != nil {
		return nil, errors.New("Could not perform request")
	}

	decoder := json.NewDecoder(res.Body)
	rpcRes = new(evm.Response)
	err = decoder.Decode(&rpcRes)
	if err != nil {
		return nil, err
	}

	err = res.Body.Close()
	if err != nil {
		return nil, err
	}

	if rpcRes.Error != nil {
		return nil, fmt.Errorf(rpcRes.Error.Message)
	}

	return rpcRes, nil
}

func (this *EvmClient) Call(method string, params interface{}) (*evm.Response, error) {
	req, err := json.Marshal(this.CreateRequest(method, params))
	if err != nil {
		return nil, err
	}

	var rpcRes *evm.Response
	time.Sleep(1 * time.Second)

	httpReq, err := http.NewRequestWithContext(context.Background(), "POST", chainnet.Current.GetJsonRpcServerUrl(), bytes.NewBuffer(req))
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	res, err := client.Do(httpReq)
	if err != nil {
		return nil, err
	}

	decoder := json.NewDecoder(res.Body)
	rpcRes = new(evm.Response)
	err = decoder.Decode(&rpcRes)
	if err != nil {
		return nil, err
	}

	err = res.Body.Close()
	if err != nil {
		return nil, err
	}
	return rpcRes, nil
}

type EvmClient struct {
}
