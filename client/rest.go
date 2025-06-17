package client

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"freemasonry.cc/blockchain/core"
	"freemasonry.cc/blockchain/core/chainnet"
	"freemasonry.cc/blockchain/util"
	"freemasonry.cc/blockchain/x/gateway/types"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdktx "github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	xauthsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	evmhd "github.com/evmos/ethermint/crypto/hd"
	"github.com/sirupsen/logrus"
	"strconv"
	"strings"
	"time"
)

type RestClient struct {
	logPrefix string
}

type TxResponse struct {
	Height    uint64 `json:"height"`
	TxHash    string `json:"txhash"`
	Codespace string `json:"codespace"`
	Code      uint32 `json:"code"`
	Data      string `json:"data"`
	RawLog    string `json:"raw_log"`
	Info      string `json:"info"`
	GasWanted uint64 `json:"gas_wanted"`
	GasUsed   uint64 `json:"gas_used"`
	Timestamp string `json:"timestamp"`
}

type BroadcastTxResponse struct {
	TxResponse TxResponse `protobuf:"bytes,1,opt,name=tx_response,json=txResponse,proto3" json:"tx_response,omitempty"`
}

type SimulateResponse struct {
	GasInfo GasInfo `json:"gas_info"`
}

type GasInfo struct {
	GasWanted uint64 `json:"gas_wanted"`
	GasUsed   uint64 `json:"gas_used"`
}

func (this *RestClient) WaitTxConfirm(txhash string) (TxResponse, error) {
	reqTotal := 0
	txres := TxResponse{}
	var returnErr error
	for {
		<-time.After(time.Second * 2)
		reqTotal++
		
		if reqTotal >= 10 {
			if returnErr != nil {
				return txres, returnErr
			} else {
				return txres, errors.New("timed out waiting for tx")
			}
		}
		findTxRes, err := restClient.FindTx(txhash)
		if err != nil {
			if strings.Contains(err.Error(), "tx not found: ") {
				continue
			} else {
				returnErr = err
				continue
				
			}
		}
		return findTxRes, nil
	}
}

func (this *RestClient) FindTx(txhash string) (txres TxResponse, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	url := fmt.Sprintf("%s/cosmos/tx/v1beta1/txs/%s", chainnet.Current.GetApiServerUrl(), txhash)
	res, err := util.HttpGet(url, time.Second*10)
	if err != nil {
		log.WithError(err).Error("HttpGet")
		return txres, chainnet.Current.GetNodeConnectionErr()
	}

	log.Debug(res)

	type GetTxResponse struct {
		Code       int    `json:"code"`
		Message    string `json:"message"`
		TxResponse struct {
			Height    string `json:"height"`
			TxHash    string `json:"txhash"`
			Codespace string `json:"codespace"`
			Code      uint32 `json:"code"`
			Data      string `json:"data"`
			RawLog    string `json:"raw_log"`
			Info      string `json:"info"`
			GasWanted string `json:"gas_wanted"`
			GasUsed   string `json:"gas_used"`
			Timestamp string `json:"timestamp"`
		} `json:"tx_response"`
	}
	getRes := GetTxResponse{}
	err = json.Unmarshal([]byte(res), &getRes)
	if err != nil {
		log.WithError(err).Error("Unmarshal")
		return txres, err
	}
	if getRes.Code != 0 {
		log.WithField("tx", txhash).Info(getRes.Message)
		return txres, errors.New(getRes.Message)
	}
	getTxRes := getRes.TxResponse
	height, _ := strconv.ParseUint(getTxRes.Height, 10, 64)
	gasWanted, _ := strconv.ParseUint(getTxRes.GasWanted, 10, 64)
	gasUsed, _ := strconv.ParseUint(getTxRes.GasUsed, 10, 64)
	txres = TxResponse{
		Height:    height,
		TxHash:    getTxRes.TxHash,
		Codespace: getTxRes.Codespace,
		Code:      getTxRes.Code,
		Data:      getTxRes.Data,
		RawLog:    getTxRes.RawLog,
		Info:      getTxRes.Info,
		GasWanted: gasWanted,
		GasUsed:   gasUsed,
		Timestamp: getTxRes.Timestamp,
	}
	
	return txres, err
}

func (this *RestClient) FindAccountNumberSeq(accountAddr string) (core.AccountNumberSeqResponse, error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	seq := core.AccountNumberSeqResponse{}
	url := fmt.Sprintf("%s/cosmos/auth/v1beta1/accounts/%s", chainnet.Current.GetApiServerUrl(), accountAddr)
	
	resp, err := util.HttpGet(url, time.Second*6)
	if err != nil {
		log.WithError(err).Error("HttpGet")
		return seq, chainnet.Current.GetNodeConnectionErr()
	}
	log.Debug(resp)

	type QueryAccountResponse struct {
		Code    int                    `json:"code"`
		Message string                 `json:"message"`
		Account map[string]interface{} `json:"account"`
	}
	
	queryTesp := QueryAccountResponse{}

	err = json.Unmarshal([]byte(resp), &queryTesp)
	if err != nil {
		log.WithError(err).Error("UnmarshalJSON")
		return seq, errors.New(queryTesp.Message)
	}

	if queryTesp.Code == 5 && strings.Contains(queryTesp.Message, "not found") {
		log.WithField("address", accountAddr).Info("not found")
		seq.AccountNumber = 0
		seq.Sequence = 0
		seq.NotFound = true
		return seq, nil
	}

	if queryTesp.Code != 0 {
		seq.AccountNumber = 0
		seq.Sequence = 0
		seq.NotFound = true
		return seq, errors.New(queryTesp.Message)
	}
	log.WithField("data", queryTesp.Account).Debug("account info")
	
	seq.NotFound = false
	typeUrl := queryTesp.Account["@type"].(string)
	if typeUrl == "/cosmos.auth.v1beta1.BaseAccount" {
		switch queryTesp.Account["account_number"].(type) {
		case uint64:
			seq.AccountNumber = queryTesp.Account["account_number"].(uint64)
		case string:
			seq.AccountNumber, _ = strconv.ParseUint(queryTesp.Account["account_number"].(string), 10, 64)
		}

		switch queryTesp.Account["sequence"].(type) {
		case uint64:
			seq.Sequence = queryTesp.Account["sequence"].(uint64)
		case string:
			seq.Sequence, _ = strconv.ParseUint(queryTesp.Account["sequence"].(string), 10, 64)
		}
	} else if typeUrl == "/ethermint.types.v1.EthAccount" {
		baseAccount := queryTesp.Account["base_account"].(map[string]interface{})

		switch baseAccount["account_number"].(type) {
		case int64:
			seq.AccountNumber = uint64(baseAccount["account_number"].(int64))
		case uint64:
			seq.AccountNumber = baseAccount["account_number"].(uint64)
		case string:
			seq.AccountNumber, _ = strconv.ParseUint(baseAccount["account_number"].(string), 10, 64)
		default:
			seq.AccountNumber = 0
		}

		switch baseAccount["sequence"].(type) {
		case int64:
			seq.Sequence = uint64(baseAccount["sequence"].(int64))
		case uint64:
			seq.Sequence = baseAccount["sequence"].(uint64)
		case string:
			seq.Sequence, _ = strconv.ParseUint(baseAccount["sequence"].(string), 10, 64)
		default:
			seq.Sequence = 0
		}
	} else {
		return seq, errors.New("not support typeurl " + typeUrl)
	}
	return seq, err
}

func (this *RestClient) GetBlockWithTx(height int64) (txres sdktx.GetBlockWithTxsResponse, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	url := fmt.Sprintf("%s/cosmos/tx/v1beta1/txs/block/%d", chainnet.Current.GetApiServerUrl(), height)
	res, err := util.HttpGet(url, time.Second*5)
	if err != nil {
		log.WithError(err).Error("HttpGet")
		return txres, chainnet.Current.GetNodeConnectionErr()
	}
	
	err = chainnet.Current.GetClientCtx().Codec.UnmarshalJSON([]byte(res), &txres)
	if err != nil {
		log.WithError(err).Error("UnmarshalJSON")
		return txres, err
	}
	return txres, err
}

func (this *RestClient) GetBalances(address string) (coins core.BalanceResponse, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	url := fmt.Sprintf("%s/cosmos/bank/v1beta1/balances/%s", chainnet.Current.GetApiServerUrl(), address)
	res, err := util.HttpGet(url, time.Second*5)
	if err != nil {
		log.WithError(err).Error("HttpGet")
		return coins, chainnet.Current.GetNodeConnectionErr()
	}
	err = json.Unmarshal([]byte(res), &coins)
	if err != nil {
		log.WithError(err).Error("UnmarshalJSON")
		return coins, err
	}
	return coins, err
}


func (this *RestClient) GetIBCLightClientLatestHeight(clientId string) (latestHeight ClientLatestHeight, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	url := fmt.Sprintf("%s/ibc/core/client/v1/client_states/%s", chainnet.Current.GetApiServerUrl(), clientId)
	res, err := util.HttpGet(url, time.Second*5)
	if err != nil {
		log.WithError(err).Error("HttpGet")
		return latestHeight, chainnet.Current.GetNodeConnectionErr()
	}
	err = json.Unmarshal([]byte(res), &latestHeight)
	if err != nil {
		log.WithError(err).Error("UnmarshalJSON")
		return latestHeight, err
	}
	return latestHeight, err
}


func (this *RestClient) GetLatestHeight() (height int64, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	url := fmt.Sprintf("%s/cosmos/base/tendermint/v1beta1/blocks/latest", chainnet.Current.GetApiServerUrl())
	res, err := util.HttpGet(url, time.Second*5)
	if err != nil {
		log.WithError(err).Error("HttpGet")
		return height, chainnet.Current.GetNodeConnectionErr()
	}
	blockData := BlockResponse{}
	err = json.Unmarshal([]byte(res), &blockData)
	if err != nil {
		log.WithError(err).Error("UnmarshalJSON")
		return height, err
	}
	height, err = strconv.ParseInt(blockData.Block.Header.Height, 10, 64)
	if err != nil {
		log.WithError(err).Error("ParseInt")
		return height, err
	}
	return height, err
}

func (this *RestClient) GetAcknowledgements(channels, portId, sequence string) (ack AcknowledgementResponse, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	url := fmt.Sprintf("%s/ibc/core/channel/v1/channels/%s/ports/%s/packet_acks/%s", chainnet.Current.GetApiServerUrl(), channels, portId, sequence)
	res, err := util.HttpGet(url, time.Second*5)
	if err != nil {
		log.WithError(err).Error("HttpGet")
		return ack, chainnet.Current.GetNodeConnectionErr()
	}
	err = json.Unmarshal([]byte(res), &ack)
	if err != nil {
		log.WithError(err).Error("UnmarshalJSON")
		return ack, err
	}
	return ack, err
}

func (this *RestClient) QueryGasPrice() sdk.DecCoin {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	defaultGasPrice := chainnet.Current.GetDefaultGasPrice()
	log.Debug("defaultGasPrice:", defaultGasPrice)
	if chainnet.Current.GetName() == chainnet.ChainNetNameDst {
		
		resBytes, _, err := util.QueryWithDataWithUnwrapErr(chainnet.Current.GetClientCtx(), "custom/"+types.ModuleName+"/"+types.QueryGasPrice, nil)
		if err != nil {
			log.WithError(err).Error("QueryWithDataWithUnwrapErr")
			return defaultGasPrice
		}
		gasPrices := sdk.DecCoins{}
		if resBytes != nil {
			err := util.Json.Unmarshal(resBytes, &gasPrices)
			if err != nil {
				log.WithError(err).Error("Unmarshal")
				return defaultGasPrice
			}
		}
		log.Debug("NodeGasPrice:", gasPrices)
		if len(gasPrices) > 0 {
			if !gasPrices[0].IsZero() {
				return gasPrices[0]
			}
		}
	}
	return defaultGasPrice
}
func (this *RestClient) SignTx(privateKey string, seqDetail core.AccountNumberSeqResponse, memo string, msgs ...sdk.Msg) (xauthsigning.Tx, error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	privKeyBytes, err := hex.DecodeString(privateKey)
	if err != nil {
		log.WithError(err).Error("hex.DecodeString")
		return nil, err
	}
	keyringAlgos := keyring.SigningAlgoList{hd.Secp256k1, evmhd.EthSecp256k1}
	algo, err := keyring.NewSigningAlgoFromString(chainnet.Current.GetAlgo(), keyringAlgos)
	if err != nil {
		log.WithError(err).Error("NewSigningAlgoFromString")
		return nil, err
	}
	privKey := algo.Generate()(privKeyBytes)

	
	simRes, err := this.Simulate(seqDetail, msgs...)
	if err != nil {
		log.WithError(err).Error("Simulate")
		return nil, util.ErrEegularFilter(err)
	}
	gasPrice := this.QueryGasPrice()
	gasUsedDec := sdk.NewDec(int64(simRes.GasInfo.GasUsed))
	gasLimit := gasUsedDec.Mul(sdk.NewDec(5)) 

	feeCoin := sdk.NewCoin(gasPrice.Denom, gasPrice.Amount.Mul(gasLimit).TruncateInt())

	log.WithFields(logrus.Fields{"gasUsed": simRes.GasInfo.GasUsed, "gasLimit": gasLimit, "gasUrice": gasPrice}).Info("CulGas")

	
	feeCoins := sdk.NewCoins(feeCoin)

	clientCtx := chainnet.Current.GetClientCtx()
	clientFactory := chainnet.Current.GetClientFactory()

	signMode := clientCtx.TxConfig.SignModeHandler().DefaultMode()
	signerData := xauthsigning.SignerData{
		ChainID:       clientCtx.ChainID,
		AccountNumber: seqDetail.AccountNumber,
		Sequence:      seqDetail.Sequence,
	}
	txBuild, err := tx.BuildUnsignedTx(clientFactory, msgs...)
	if err != nil {
		log.WithError(err).Error("tx.BuildUnsignedTx")
		return nil, err
	}

	txBuild.SetGasLimit(gasLimit.TruncateInt().Uint64()) 
	txBuild.SetFeeAmount(feeCoins)                       
	txBuild.SetMemo(memo)                                
	sigData := signing.SingleSignatureData{
		SignMode:  signMode,
		Signature: nil,
	}
	sig := signing.SignatureV2{
		PubKey:   privKey.PubKey(),
		Data:     &sigData,
		Sequence: seqDetail.Sequence,
	}
	
	if err = txBuild.SetSignatures(sig); err != nil {
		log.WithError(err).Error("SetSignatures")
		return nil, err
	}
	signV2, err := tx.SignWithPrivKey(signMode, signerData, txBuild, privKey, clientCtx.TxConfig, seqDetail.Sequence)
	if err != nil {
		log.WithError(err).Error("SignWithPrivKey")
		return nil, err
	}
	err = txBuild.SetSignatures(signV2)
	if err != nil {
		log.WithError(err).Error("SetSignatures")
		return nil, err
	}

	signedTx := txBuild.GetTx()
	return signedTx, nil
}

func (this *RestClient) SignTxBytes(seqDetail core.AccountNumberSeqResponse, privateKey string, memo string, msg ...sdk.Msg) (txBytes []byte, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	
	signedTx, err := this.SignTx(privateKey, seqDetail, memo, msg...)
	if err != nil {
		log.WithError(err).Error("SignTx")
		return nil, err
	}
	
	signPubkey, err := signedTx.GetPubKeys()
	if err != nil {
		log.WithError(err).Error("GetPubKeys")
		return nil, err
	}

	signV2, _ := signedTx.GetSignaturesV2()
	senderAddrBytes := signV2[0].PubKey.Address().Bytes()
	signAddrBytes := signPubkey[0].Address().Bytes()
	if !bytes.Equal(signAddrBytes, senderAddrBytes) {
		return nil, errors.New("sign error")
	}

	
	txBytes, err = chainnet.Current.GetClientCtx().TxConfig.TxEncoder()(signedTx)
	if err != nil {
		log.WithError(err).Error("TxEncoder")
		return nil, err
	}
	return
}

func (this *RestClient) BroadcastTx(req sdktx.BroadcastTxRequest) (txres TxResponse, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)

	reqBytes, err := json.Marshal(req)
	if err != nil {
		log.WithError(err).Error("Marshal")
		return txres, err
	}
	url := fmt.Sprintf("%s/cosmos/tx/v1beta1/txs", chainnet.Current.GetApiServerUrl())
	res, err := util.HttpPostJson(url, reqBytes)
	if err != nil {
		log.WithError(err).Error("HttpPostJson")
		return txres, chainnet.Current.GetNodeConnectionErr()
	}
	log.WithField("res", res).Debug("BroadcastTx res")
	type GetBroadcastTxResponse struct {
		Code       int    `json:"code"`    
		Message    string `json:"message"` 
		TxResponse struct {
			Height    string `json:"height"`
			TxHash    string `json:"txhash"`
			Codespace string `json:"codespace"`
			Code      uint32 `json:"code"`
			Data      string `json:"data"`
			RawLog    string `json:"raw_log"`
			Info      string `json:"info"`
			GasWanted string `json:"gas_wanted"`
			GasUsed   string `json:"gas_used"`
			Timestamp string `json:"timestamp"`
		} `json:"tx_response"`
	}
	getRes := GetBroadcastTxResponse{}
	err = json.Unmarshal([]byte(res), &getRes)
	if err != nil {
		log.WithError(err).Error("Unmarshal")
		return txres, err
	}
	
	if getRes.Code != 0 {
		return txres, errors.New(getRes.Message)
	}
	getTxRes := getRes.TxResponse
	height, _ := strconv.ParseUint(getTxRes.Height, 10, 64)
	gasWanted, _ := strconv.ParseUint(getTxRes.GasWanted, 10, 64)
	gasUsed, _ := strconv.ParseUint(getTxRes.GasUsed, 10, 64)
	txres = TxResponse{
		Height:    height,
		TxHash:    getTxRes.TxHash,
		Codespace: getTxRes.Codespace,
		Code:      getTxRes.Code,
		Data:      getTxRes.Data,
		RawLog:    getTxRes.RawLog,
		Info:      getTxRes.Info,
		GasWanted: gasWanted,
		GasUsed:   gasUsed,
		Timestamp: getTxRes.Timestamp,
	}
	return txres, err
}

func (this *RestClient) Simulate(user core.AccountNumberSeqResponse, msg ...sdk.Msg) (txres SimulateResponse, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	clientCtx := chainnet.Current.GetClientCtx()
	clientFactory := chainnet.Current.GetClientFactory()
	txBuild, err := tx.BuildUnsignedTx(clientFactory, msg...)
	if err != nil {
		log.WithError(err).Error("BuildUnsignedTx")
		return
	}

	signMode := clientCtx.TxConfig.SignModeHandler().DefaultMode()
	sigData := signing.SingleSignatureData{
		SignMode:  signMode,
		Signature: nil,
	}
	sig := signing.SignatureV2{
		PubKey:   &secp256k1.PubKey{Key: []byte{}},
		Data:     &sigData,
		Sequence: user.Sequence,
	}
	
	if err = txBuild.SetSignatures(sig); err != nil {
		log.WithError(err).Error("SetSignatures")
		return
	}

	unsignedTx := txBuild.GetTx()
	txBytes, err := clientCtx.TxConfig.TxEncoder()(unsignedTx)
	if err != nil {
		log.WithError(err).Error("TxEncoder")
		return
	}
	req := sdktx.SimulateRequest{
		TxBytes: txBytes,
	}
	reqBytes, err := json.Marshal(req)
	if err != nil {
		log.WithError(err).Error("Marshal")
		return txres, err
	}
	url := fmt.Sprintf("%s/cosmos/tx/v1beta1/simulate", chainnet.Current.GetApiServerUrl())
	res, err := util.HttpPost(url, reqBytes, "application/json")
	if err != nil {
		log.WithError(err).Error("HttpPostJson")
		return txres, chainnet.Current.GetNodeConnectionErr()
	}

	log.Debug(res)

	type SimulateResponse2 struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		GasInfo struct {
			GasWanted string `json:"gas_wanted"`
			GasUsed   string `json:"gas_used"`
		} `json:"gas_info"`
	}
	txres2 := SimulateResponse2{}
	err = json.Unmarshal([]byte(res), &txres2)
	if err != nil {
		log.WithError(err).Error("UnmarshalJSON")
		return txres, err
	}
	if txres2.Code != 0 {
		err = errors.New(txres2.Message)
		return txres, err
	}
	txres.GasInfo.GasUsed, _ = strconv.ParseUint(txres2.GasInfo.GasUsed, 10, 64)
	txres.GasInfo.GasWanted, _ = strconv.ParseUint(txres2.GasInfo.GasWanted, 10, 64)
	return txres, err
}
