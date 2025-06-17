package client

import (
	"context"
	"encoding/hex"
	"freemasonry.cc/blockchain/app"
	"freemasonry.cc/blockchain/cmd/stcd/cmd"
	"freemasonry.cc/blockchain/core"
	"freemasonry.cc/blockchain/core/chainnet"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/migrations/legacytx"
	types2 "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/evmos/ethermint/encoding"
	evmoskr "github.com/evmos/evmos/v10/crypto/keyring"
	"github.com/tendermint/tendermint/crypto/tmhash"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	ttypes "github.com/tendermint/tendermint/types"
	"os"
)


type MempoolClient struct {
	logPrefix string
}

func (this *MempoolClient) GetAccountSequence(accountAddr string) (core.AccountNumberSeqResponse, error) {
	seq := core.AccountNumberSeqResponse{}
	accAddr, err := sdk.AccAddressFromBech32(accountAddr)
	if err != nil {
		return seq, err
	}
	encodingConfig := encoding.MakeConfig(app.ModuleBasics)
	clientCtx := client.Context{}.WithChainID(chainnet.ChainNetDst.GetChainId()).WithBroadcastMode(flags.BroadcastSync).
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
		seq.NotFound = true
		return seq, err
	}

	seq.AccountNumber = accountNumber
	seq.Sequence = sequence

	pendingTxs, err := this.UnconfirmedTxs()
	if err != nil {
		seq.NotFound = false
		return seq, err
	}
	for _, tx := range pendingTxs {
		for _, msg := range (*tx).GetMsgs() {
			for _, acc := range msg.GetSigners() {
				if acc.Equals(accAddr) {
					seq.Sequence = seq.Sequence + 1
				}
			}
		}
	}
	return seq, err
}

func (this *MempoolClient) TxExisted(tx ttypes.Tx) (res *ctypes.ResultTxExisted, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	clientCtx := chainnet.ChainNetDst.GetClientCtx()
	node, err := clientCtx.GetNode()
	if err != nil {
		log.WithError(err).Error("GetNode")
		return res, err
	}
	return node.TxExisted(context.Background(), tx)
}

func (this *MempoolClient) CheckTx(tx ttypes.Tx) (resultTx *ctypes.ResultCheckTx, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	clientCtx := chainnet.ChainNetDst.GetClientCtx()
	node, err := clientCtx.GetNode()
	if err != nil {
		log.WithError(err).Error("GetNode")
		return resultTx, err
	}
	return node.CheckTx(context.Background(), tx)
}


func (this *MempoolClient) UnconfirmedTxs() (list []*legacytx.StdTx, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	clientCtx := chainnet.ChainNetDst.GetClientCtx()
	node, err := clientCtx.GetNode()
	if err != nil {
		log.WithError(err).Error("GetNode")
		return list, err
	}
	limit := 100
	resp, err := node.UnconfirmedTxs(context.Background(), &limit)
	if err != nil {
		log.WithError(err).Error("UnconfirmedTxs")
		return list, err
	}
	for _, tx := range resp.Txs {
		cosmosTx, err := termintTx2CosmosTx(tx)
		if err != nil {
			log.WithError(err).Error("TermintTx2CosmosTx1")
			return list, err
		}
		stdTx, err := convertTxToStdTx(cosmosTx)
		if err != nil {
			log.WithError(err).Error("ConvertTxToStdTx1")
			return list, err
		}
		stdTx.Memo = hex.EncodeToString(tmhash.Sum(tx))
		list = append(list, stdTx)
	}
	return list, nil
}


func (this *MempoolClient) UnconfirmedTx() (list *ctypes.ResultUnconfirmedTxs, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	clientCtx := chainnet.ChainNetDst.GetClientCtx()
	node, err := clientCtx.GetNode()
	if err != nil {
		log.WithError(err).Error("GetNode")
		return list, err
	}
	limit := 100
	resp, err := node.UnconfirmedTxs(context.Background(), &limit)
	if err != nil {
		log.WithError(err).Error("UnconfirmedTxs")
		return list, err
	}
	return resp, nil
}


func (this *MempoolClient) GetPoolFirstAndLastTxs() (frist *legacytx.StdTx, last *legacytx.StdTx, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	clientCtx := chainnet.ChainNetDst.GetClientCtx()
	node, err := clientCtx.GetNode()
	if err != nil {
		log.WithError(err).Error("GetNode")
		return frist, last, err
	}
	resp, err := node.PoolFirstAndLastTxs(context.Background())
	if err != nil {
		log.WithError(err).Error("UnconfirmedTxs")
		return frist, last, err
	}
	var cosmosTx sdk.Tx
	if resp.First != nil {
		cosmosTx, err = termintTx2CosmosTx(resp.First)
		if err != nil {
			log.WithError(err).Error("TermintTx2CosmosTx1")
			return frist, last, err
		}
		frist, err = convertTxToStdTx(cosmosTx)
		if err != nil {
			log.WithError(err).Error("ConvertTxToStdTx1")
			return frist, last, err
		}
	}
	if resp.Last != nil {
		cosmosTx, err = termintTx2CosmosTx(resp.Last)
		if err != nil {
			log.WithError(err).Error("TermintTx2CosmosTx2")
			return frist, last, err
		}
		last, err = convertTxToStdTx(cosmosTx)
		if err != nil {
			log.WithError(err).Error("ConvertTxToStdTx2")
			return frist, last, err
		}
	}

	return frist, last, nil
}
