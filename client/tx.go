package client

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"freemasonry.cc/blockchain/app"
	"freemasonry.cc/blockchain/cmd/stcd/cmd"
	"freemasonry.cc/blockchain/core"
	"freemasonry.cc/blockchain/core/chainnet"
	"freemasonry.cc/blockchain/util"
	daoTypes "freemasonry.cc/blockchain/x/dao/types"
	"freemasonry.cc/blockchain/x/gateway/types"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkErrors "github.com/cosmos/cosmos-sdk/types/errors"
	sdktx "github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	"github.com/cosmos/cosmos-sdk/x/auth/migrations/legacytx"
	xauthsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	types2 "github.com/cosmos/cosmos-sdk/x/auth/types"
	types3 "github.com/cosmos/cosmos-sdk/x/bank/types"
	govTypes "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	evmhd "github.com/evmos/ethermint/crypto/hd"
	"github.com/evmos/ethermint/encoding"
	evmoskr "github.com/evmos/evmos/v10/crypto/keyring"
	"github.com/sirupsen/logrus"
	"github.com/tendermint/tendermint/crypto/tmhash"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	ttypes "github.com/tendermint/tendermint/types"
	"os"
	"regexp"
	"strings"
)

type TxClient struct {
	logPrefix string
}


func (txClient *TxClient) GasMsg(msgType, msgs string) (sdk.Msg, error) {
	log := core.BuildLog(core.GetStructFuncName(txClient), core.LmChainClient)
	msgByte := []byte(msgs)
	if unmashal, ok := msgUnmashalHandles[msgType]; ok {
		msg, err := unmashal(msgByte)
		if err != nil {
			log.WithError(err).WithField("msgType", msgType).Error("gasMsg.Unmarshal MsgSend")
			return nil, err
		}
		return msg, nil
	} else {
		return nil, errors.New("unregister unmashal msg type:" + msgType)
	}
}


func (txClient *TxClient) GasInfo(seqDetail core.AccountNumberSeqResponse, msg ...sdk.Msg) (sdk.Coin, uint64, sdk.DecCoin, sdk.Coins, error) {
	log := core.BuildLog(core.GetStructFuncName(txClient), core.LmChainClient)
	clientCtx := chainnet.ChainNetDst.GetClientCtx()
	clientFactory := chainnet.ChainNetDst.GetClientFactory()
	fee := sdk.NewCoin(core.BaseDenom, sdk.ZeroInt())
	gasPrice := sdk.NewDecCoin(core.BaseDenom, sdk.ZeroInt())
	minFee, err := sdk.ParseCoinNormalized(core.ChainDefaultFeeStr)
	if err != nil {
		return minFee, 0, gasPrice, sdk.NewCoins(minFee), err
	}

	
	
	
	
	
	mempoolClient := NewMempoolClient()
	pendingTxs, err := mempoolClient.UnconfirmedTxs()
	if err != nil {
		return minFee, 0, gasPrice, sdk.NewCoins(minFee), err
	}
	for _, tx := range pendingTxs {
		for _, tm := range (*tx).GetMsgs() {
			for _, acc := range tm.GetSigners() {
				for _, sm := range msg {
					for _, address := range sm.GetSigners() {
						if acc.Equals(address) {
							return minFee, 0, gasPrice, sdk.NewCoins(minFee), core.UnconfirmedTxsErr
						}
					}
				}
			}
		}
	}
	clientFactory = clientFactory.WithSequence(seqDetail.Sequence).WithSimulateAndExecute(true)
	gasInfo, _, err := tx.CalculateGas(clientCtx, clientFactory, msg...)

	
	gasAdd := uint64(0)
	for _, m := range msg {
		gasAdd = gasAdd + GatAddGasFromMsg(m)
	}

	if err != nil {
		sp := strings.Split(err.Error(), ": ")
		size := len(sp)
		if size-2 >= 0 {
			err = errors.New(sp[size-2])
		} else if size-1 >= 0 {
			err = errors.New(sp[size-1])
		}
		log.WithError(err).Error("tx.CalculateGas")
		return minFee, 0, gasPrice, sdk.NewCoins(minFee), err
	}
	gas := gasInfo.GasInfo.GasUsed*3 + gasAdd
	
	gasPrice = txClient.QueryGasPrice()
	gasDec := sdk.NewDec(int64(gas))

	log.WithFields(logrus.Fields{
		"gasPriceDec": gasPrice.String(),
		"gasDec":      gasDec.String(),
		"gas1":        gasPrice.Amount.Mul(gasDec),
		"gas2":        gasPrice.Amount.Mul(gasDec).TruncateInt(),
	}).Debug("GasInfo")
	fee = sdk.NewCoin(gasPrice.Denom, gasPrice.Amount.Mul(gasDec).TruncateInt())
	payFee, err := txClient.QueryGasDeduction(sdk.NewCoins(fee), msg...)
	if err != nil {
		return minFee, 0, gasPrice, sdk.NewCoins(minFee), err
	}
	return fee, gas, gasPrice, payFee, nil
}

func GatAddGasFromMsg(msg sdk.Msg) uint64 {
	
	
	
	return 100000
}


func (txClient *TxClient) QueryGasDeduction(fee sdk.Coins, msg ...sdk.Msg) (sdk.Coins, error) {
	log := core.BuildLog(core.GetStructFuncName(txClient), core.LmChainClient)
	params := daoTypes.GasDeductionParams{
		Msg: msg,
		Fee: fee,
	}
	clientCtx := chainnet.ChainNetDst.GetClientCtx()
	reqData, err := clientCtx.LegacyAmino.MarshalJSON(params)
	if err != nil {
		return nil, err
	}
	resBytes, _, err := util.QueryWithDataWithUnwrapErr(clientCtx, "custom/dao/query_gas_deduction", reqData)
	if err != nil {
		log.WithError(err).Error("QueryWithData")
		return nil, err
	}
	payFee := sdk.Coins{}
	if resBytes != nil {
		err := util.Json.Unmarshal(resBytes, &payFee)
		if err != nil {
			return nil, err
		}
	}
	return payFee, nil
}

func (txClient *TxClient) QueryGasPrice() sdk.DecCoin {
	log := core.BuildLog(core.GetStructFuncName(txClient), core.LmChainClient)
	defaultGasPrice := chainnet.ChainNetDst.GetDefaultGasPrice()
	clientCtx := chainnet.ChainNetDst.GetClientCtx()
	resBytes, _, err := util.QueryWithDataWithUnwrapErr(clientCtx, "custom/"+types.ModuleName+"/"+types.QueryGasPrice, nil)
	if err != nil {
		log.WithError(err).Error("QueryWithData")
		return defaultGasPrice
	}
	gasPrices := sdk.DecCoins{}
	if resBytes != nil {
		err := util.Json.Unmarshal(resBytes, &gasPrices)
		if err != nil {
			return defaultGasPrice
		}
	}
	if len(gasPrices) > 0 {
		if !gasPrices[0].IsZero() {
			return gasPrices[0]
		}
	}
	return defaultGasPrice
}


func (txClient *TxClient) GetTranserTypeConfig() map[string]string {
	return core.GetTranserTypeConfig()
}

func (txClient *TxClient) ConvertTxToStdTx(cosmosTx sdk.Tx) (*legacytx.StdTx, error) {
	signingTx, ok := cosmosTx.(xauthsigning.Tx)
	if !ok {
		return nil, errors.New("tx to stdtx error")
	}

	stdTx, err := tx.ConvertTxToStdTx(app.EncodingConfig.Amino, signingTx)
	if err != nil {
		return nil, err
	}
	return &stdTx, nil
}

func (txClient *TxClient) TermintTx2CosmosTx(signTxs ttypes.Tx) (sdk.Tx, error) {
	return app.EncodingConfig.TxConfig.TxDecoder()(signTxs)
}

func (txClient *TxClient) SignTx2Bytes(signTxs xauthsigning.Tx) ([]byte, error) {
	return app.EncodingConfig.TxConfig.TxEncoder()(signTxs)
}

func (txClient *TxClient) SetFee(signTxs xauthsigning.Tx) ([]byte, error) {
	return app.EncodingConfig.TxConfig.TxEncoder()(signTxs)
}

func (txClient *TxClient) FindByByte(txhash []byte) (resultTx *ctypes.ResultTx, notFound bool, err error) {
	notFound = false
	log := core.BuildLog(core.GetStructFuncName(txClient), core.LmChainClient)
	clientCtx := chainnet.ChainNetDst.GetClientCtx()
	node, err := clientCtx.GetNode()
	if err != nil {
		log.WithError(err).Error("GetNode")
		return
	}
	resultTx, err = node.Tx(context.Background(), txhash, true)
	if err != nil {
		
		notFound = txClient.isTxNotFoundError(err.Error())
		if notFound {
			err = nil
		} else {
			log.WithError(err).WithField("txhash", hex.EncodeToString(txhash)).Error("node.Tx")
		}
		return
	}
	return
}

func (txClient *TxClient) FindByHex(txhashStr string) (resultTx *ctypes.ResultTx, notFound bool, err error) {
	var txhash []byte
	notFound = false
	log := core.BuildLog(core.GetStructFuncName(txClient), core.LmChainClient)
	txhash, err = hex.DecodeString(txhashStr)
	if err != nil {
		log.WithError(err).WithField("txhash", txhashStr).Error("hex.DecodeString")
		return
	}
	return txClient.FindByByte(txhash)
}

func (txClient *TxClient) isTxNotFoundError(errContent string) (ok bool) {
	errRegexp := `tx\ \([0-9A-Za-z]{64}\)\ not\ found`
	r, err := regexp.Compile(errRegexp)
	if err != nil {
		return false
	}
	if r.Match([]byte(errContent)) {
		return true
	} else {
		return false
	}
}


func (txClient *TxClient) SignAndSendMsg(address string, privateKey string, fee legacytx.StdFee, memo string, msg ...sdk.Msg) (tx ttypes.Tx, txRes *core.BaseResponse, err error) {
	log := core.BuildLog(core.GetStructFuncName(txClient), core.LmChainClient)
	
	seqDetail, err := txClient.FindAccountNumberSeq(address)
	if err != nil {
		return
	}

	
	

	
	signedTx, err := txClient.SignTx(privateKey, seqDetail, fee, memo, msg...)
	if err != nil {
		log.WithError(err).Error("SignTx")
		return
	}
	
	signPubkey, err := signedTx.GetPubKeys()
	if err != nil {
		log.WithError(err).Error("GetPubKeys")
		return
	}

	signV2, _ := signedTx.GetSignaturesV2()
	senderAddrBytes := signV2[0].PubKey.Address().Bytes()
	signAddrBytes := signPubkey[0].Address().Bytes()
	if !bytes.Equal(signAddrBytes, senderAddrBytes) {
		return nil, nil, errors.New("sign error")
	}

	
	tx, err = txClient.SignTx2Bytes(signedTx)
	if err != nil {
		log.WithError(err).Error("SignTx2Bytes")
		return
	}
	fmt.Println("txhash:", hex.EncodeToString(tmhash.Sum(tx)))
	
	txRes, err = txClient.Send(tx)
	if err != nil {
		log.WithError(err).Error("Send")
		return
	}
	if txRes != nil && txRes.Data != nil {
		var bytes []byte
		broadcastTxResponse := core.BroadcastTxResponse{}
		bytes, err = json.Marshal(txRes.Data)
		if err != nil {
			log.WithError(err).Error("Marshal")
			return
		}
		err = json.Unmarshal(bytes, &broadcastTxResponse)
		if err != nil {
			log.WithError(err).Error("Unmarshal")
			return
		}
		broadcastTxResponse.SignedTxStr = hex.EncodeToString(tx)
		txRes.Data = broadcastTxResponse
	}
	return
}

func (txClient *TxClient) FindAccountNumberSeq(accountAddr string) (core.AccountNumberSeqResponse, error) {
	logs := core.BuildLog(core.GetStructFuncName(txClient), core.LmChainClient)
	seq := core.AccountNumberSeqResponse{}
	accAddr, err := sdk.AccAddressFromBech32(accountAddr)
	if err != nil {
		return seq, err
	}
	encodingConfig := encoding.MakeConfig(app.ModuleBasics)
	clientCtx := chainnet.ChainNetDst.GetClientCtx().WithBroadcastMode(flags.BroadcastBlock).
		WithInterfaceRegistry(encodingConfig.InterfaceRegistry).
		WithCodec(encodingConfig.Codec).
		WithTxConfig(encodingConfig.TxConfig).
		WithLegacyAmino(encodingConfig.Amino).
		WithInput(os.Stdin).
		WithAccountRetriever(types2.AccountRetriever{}).
		WithHomeDir(app.DefaultNodeHome).
		WithKeyringOptions(evmoskr.Option()).
		WithViper(cmd.EnvPrefix).
		WithLedgerHasProtobuf(true)
	accountNumber, sequence, err := clientCtx.AccountRetriever.GetAccountNumberSequence(clientCtx, accAddr)
	if err != nil {
		
		accClient := NewAccountClient(txClient)
		nextAccNum, err := accClient.GetNextAccountNumber()
		if err != nil {
			return seq, err
		}
		logs.Info("NewAccNum:", accountNumber)
		seq.NotFound = true
		seq.AccountNumber = nextAccNum
		return seq, err
	}

	seq.AccountNumber = accountNumber
	seq.Sequence = sequence
	return seq, err

}


func (txClient *TxClient) Send(req []byte) (txRes *core.BaseResponse, err error) {
	
	var txBytes []byte
	var res *sdk.TxResponse
	if req != nil {
		txBytes = req
	}
	encodingConfig := encoding.MakeConfig(app.ModuleBasics)
	clientCtx := chainnet.ChainNetDst.GetClientCtx().WithBroadcastMode(flags.BroadcastBlock).
		WithInterfaceRegistry(encodingConfig.InterfaceRegistry).
		WithCodec(encodingConfig.Codec).
		WithTxConfig(encodingConfig.TxConfig).
		WithLegacyAmino(encodingConfig.Amino).
		WithInput(os.Stdin).
		WithAccountRetriever(types2.AccountRetriever{}).
		WithHomeDir(app.DefaultNodeHome).
		WithKeyringOptions(evmoskr.Option()).
		WithViper(cmd.EnvPrefix).
		WithLedgerHasProtobuf(true)
	res, err = clientCtx.BroadcastTx(txBytes)

	txRes = &core.BaseResponse{}
	if res == nil && err == nil {
		txRes.Status = 0
		txRes.Info = core.BroadcastTimeOut.Error()
		txResponse := core.BroadcastTxResponse{}
		txResponse.Code = core.BroadcastTimeOut.ABCICode()
		txResponse.CodeSpace = core.BroadcastTimeOut.Codespace()
		txResponse.TxHash = hex.EncodeToString(txBytes)
		txRes.Data = txResponse
		return
	}

	if err != nil {
		sp := strings.Split(err.Error(), ": ")
		txRes.Info = sp[len(sp)-1]
		return
	}

	txResponse := core.BroadcastTxResponse{}
	txResponse.Code = res.Code
	txResponse.CodeSpace = res.Codespace
	txResponse.TxHash = res.TxHash
	txResponse.Height = res.Height

	if res.Code == 0 { 
		txRes.Status = 1
	} else {
		txRes.Status = 0
	}

	txRes.Info = ParseErrorCode(res.Code, res.Codespace, res.RawLog)
	txRes.Data = txResponse

	return
}

type SetChatInfo struct {
	FromAddress        string `json:"from_address,omitempty" yaml:"from_address"`
	NodeAddress        string `json:"node_address,omitempty" yaml:"node_address"`
	AddressBook        string `json:"address_book,omitempty" yaml:"address_book"`
	ChatBlacklist      string `json:"chat_blacklist,omitempty" yaml:"chat_blacklist"`
	ChatRestrictedMode string `json:"chat_restricted_mode,omitempty" yaml:"chat_limit"`
	ChatFeeAmount      string `json:"chat_fee_amount" yaml:"chat_fee_amount"`
	ChatFeeCoinSymbol  string `json:"chat_fee_coin_symbol" yaml:"chat_fee_coin_symbol"`
	ChatWhitelist      string `json:"chat_whitelist,omitempty" yaml:"chat_whitelist"`
	UpdateTime         int64  `json:"update_time,omitempty" yaml:"update_time"`
	ChatBlacklistEnc   string `json:"chat_blacklist_enc,omitempty" yaml:"chat_blacklist_enc"`
	ChatWhitelistEnc   string `json:"chat_whitelist_enc,omitempty" yaml:"chat_whitelist_enc"`
	Remarks            string `json:"remarks,omitempty" yaml:"remarks"`
}

type CrossChainOut struct {
	SendAddress string `json:"send_address"`
	ToAddress   string `json:"to_address"`
	CoinAmount  string `json:"coin_amount"`
	CoinSymbol  string `json:"coin_symbol"`
	ChainType   string `json:"chain_type"`
	Remark      string `json:"remark"`
}


func (txClient *TxClient) SignTx(privateKey string, seqDetail core.AccountNumberSeqResponse, fee legacytx.StdFee, memo string, msgs ...sdk.Msg) (xauthsigning.Tx, error) {
	log := core.BuildLog(core.GetStructFuncName(txClient), core.LmChainClient)
	clientCtx := chainnet.ChainNetDst.GetClientCtx()
	clientFactory := chainnet.ChainNetDst.GetClientFactory()
	privKeyBytes, err := hex.DecodeString(privateKey)
	if err != nil {
		log.WithError(err).Error("hex.DecodeString")
		return nil, err
	}
	keyringAlgos := keyring.SigningAlgoList{evmhd.EthSecp256k1}
	algo, err := keyring.NewSigningAlgoFromString("eth_secp256k1", keyringAlgos)
	if err != nil {
		return nil, err
	}
	privKey := algo.Generate()(privKeyBytes)
	
	if fee.Gas == flags.DefaultGasLimit {
		feeAmount, gas, _, _, err := txClient.GasInfo(seqDetail, msgs...)
		if err != nil {
			log.WithError(err).Error("CulGas")
			return nil, util.ErrFilter(err)
		}
		log.WithField("gas", gas).Info("CulGas:")
		fee.Gas = gas
		fee.Amount = sdk.NewCoins(feeAmount)
	}
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
	nodeInfo, err := clientCtx.Client.Status(context.Background())
	if err != nil {
		return nil, err
	}
	txBuild.SetGasLimit(fee.Gas)     
	txBuild.SetFeeAmount(fee.Amount) 
	txBuild.SetMemo(memo)            
	txBuild.SetTimeoutHeight(uint64(nodeInfo.SyncInfo.LatestBlockHeight + core.TxTimeoutHeight))
	sigData := signing.SingleSignatureData{
		SignMode:  signMode,
		Signature: nil,
	}
	sig := signing.SignatureV2{
		PubKey:   privKey.PubKey(),
		Data:     &sigData,
		Sequence: seqDetail.Sequence,
	}
	
	if err := txBuild.SetSignatures(sig); err != nil {
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

func (txClient *TxClient) QueryProposer(proposerId uint64) (resp *govTypes.Proposal, err error) {
	log := core.BuildLog(core.GetStructFuncName(txClient), core.LmChainClient)
	clientCtx := chainnet.ChainNetDst.GetClientCtx()
	params := govTypes.QueryProposalParams{ProposalID: proposerId}
	bz, err := clientCtx.LegacyAmino.MarshalJSON(params)
	if err != nil {
		log.WithError(err).Error("MarshalJSON")
		return
	}
	resBytes, _, err := util.QueryWithDataWithUnwrapErr(clientCtx, "custom/gov/"+govTypes.QueryProposal, bz)
	if err != nil {
		log.WithError(err).Error("QueryWithData")
		return
	}
	var res govTypes.Proposal
	if resBytes != nil {
		err = clientCtx.LegacyAmino.UnmarshalJSON(resBytes, &res)
		if err != nil {
			return
		}
		resp = &res
	}
	return
}

func (txClient *TxClient) QueryVoteDetail(proposerId uint64) (resp VoteDetail, err error) {
	log := core.BuildLog(core.GetStructFuncName(txClient), core.LmChainClient)
	clientCtx := chainnet.ChainNetDst.GetClientCtx()
	
	params := govTypes.QueryProposalParams{ProposalID: proposerId}
	bz, err := clientCtx.LegacyAmino.MarshalJSON(params)
	if err != nil {
		log.WithError(err).Error("MarshalJSON")
		return
	}
	resBytes, _, err := util.QueryWithDataWithUnwrapErr(clientCtx, "custom/gov/"+govTypes.QueryProposal, bz)
	if err != nil {
		log.WithError(err).Error("QueryWithData")
		return
	}
	var proposalBaseDetail govTypes.Proposal
	if resBytes != nil {
		err = clientCtx.LegacyAmino.UnmarshalJSON(resBytes, &proposalBaseDetail)
		if err != nil {
			return
		}
	}

	var govMetaData GovMetaData

	err = json.Unmarshal([]byte(proposalBaseDetail.Metadata), &govMetaData)
	if err != nil {
		return
	}

	proposalBaseDetail.Metadata = govMetaData.Description

	detailStr, err := clientCtx.LegacyAmino.MarshalJSON(proposalBaseDetail)
	if err != nil {
		return
	}

	resp.Detail = string(detailStr)

	resp.From = govMetaData.Proposer

	msgs, err := sdktx.GetMsgs(proposalBaseDetail.Messages, "sdk.MsgProposal")
	if err != nil {
		return
	}

	if msg, ok := msgs[0].(*govTypes.MsgExecLegacyContent); ok {
		content, err := govTypes.LegacyContentFromMessage(msg)
		if err != nil {
			return resp, err
		}

		resp.ProposalType = content.ProposalType()
	}

	
	depositCoin := sdk.Coins(proposalBaseDetail.TotalDeposit).AmountOf(core.GovDenom)
	resp.Deposit = sdk.NewCoin(core.GovDenom, depositCoin)

	
	resBytes, _, err = util.QueryWithDataWithUnwrapErr(clientCtx, "custom/gov/"+govTypes.QueryTally, bz)
	if err != nil {
		log.WithError(err).Error("QueryWithData")
		return
	}
	var resTally govTypes.TallyResult
	if resBytes != nil {
		err = clientCtx.LegacyAmino.UnmarshalJSON(resBytes, &resTally)
		if err != nil {
			return
		}
	}

	resp.Tally = resTally

	return
}

type GovMetaData struct {
	Description string `json:"description"`
	Proposer    string `json:"proposer"`
}

type VoteDetail struct {
	From         string               `json:"from"`
	Detail       string               `json:"detail"`
	Deposit      sdk.Coin             `json:"deposit"`
	Tally        govTypes.TallyResult `json:"tally"`
	ProposalType string               `json:"proposal_type"` 
}

func ParseErrorCode(code uint32, codeSpace string, rowlog string) string {
	if codeSpace == sdkErrors.RootCodespace {
		if code == sdkErrors.ErrInsufficientFee.ABCICode() { 
			return core.FeeIsTooLess
		} else if code == sdkErrors.ErrOutOfGas.ABCICode() { 
			return core.ErrorGasOut
		} else if code == sdkErrors.ErrUnauthorized.ABCICode() { 
			return core.ErrUnauthorized
		} else if code == sdkErrors.ErrWrongSequence.ABCICode() { 
			return core.ErrWrongSequence
		}
	}

	
	sp := strings.Split(rowlog, ": ")
	return sp[len(sp)-1]
}

func (txClient *TxClient) QueryProposalVotes(proposerId uint64, page, limit int) (resp govTypes.Votes, err error) {
	log := core.BuildLog(core.GetStructFuncName(txClient), core.LmChainClient)
	params := govTypes.QueryProposalVotesParams{
		ProposalID: proposerId,
		Page:       page,
		Limit:      limit,
	}
	clientCtx := chainnet.ChainNetDst.GetClientCtx()
	bz, err := clientCtx.LegacyAmino.MarshalJSON(params)
	if err != nil {
		log.WithError(err).Error("MarshalJSON")
		return
	}
	resBytes, _, err := util.QueryWithDataWithUnwrapErr(clientCtx, "custom/gov/"+govTypes.QueryVotes, bz)
	if err != nil {
		log.WithError(err).Error("QueryWithData")
		return
	}
	var proposalBaseDetail govTypes.Votes
	if resBytes != nil {
		err = clientCtx.LegacyAmino.UnmarshalJSON(resBytes, &proposalBaseDetail)
		if err != nil {
			return
		}
	}

	return proposalBaseDetail, nil
}

func (cc TxClient) Querybalances(fromAddress string) (sdk.Coins, error) {
	logs := core.BuildLog(core.GetStructFuncName(cc), core.LmChainClient)

	req := types3.QueryAllBalancesRequest{
		Address: fromAddress,
	}

	clientCtx := chainnet.ChainNetDst.GetClientCtx()
	paramsData, err := clientCtx.LegacyAmino.MarshalJSON(req)
	if err != nil {
		return nil, err
	}

	resBytes, _, err := util.QueryWithDataWithUnwrapErr(clientCtx, "custom/bank/"+types3.QueryAllBalances, paramsData)
	if err != nil {
		logs.WithError(err).Error("QueryWithData")
		return nil, err
	}

	resp := sdk.Coins{}
	err = clientCtx.LegacyAmino.UnmarshalJSON(resBytes, &resp)
	if err != nil {
		logs.WithError(err).Error("UnmarshalJSON")
		return resp, err
	}
	return resp, nil
}
