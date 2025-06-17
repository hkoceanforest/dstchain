package client

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"freemasonry.cc/blockchain/core"
	"freemasonry.cc/blockchain/core/chainnet"
	"freemasonry.cc/blockchain/util"
	"freemasonry.cc/blockchain/x/dao/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/migrations/legacytx"
	"github.com/cosmos/cosmos-sdk/x/group"
	"github.com/sirupsen/logrus"
	ttypes "github.com/tendermint/tendermint/types"
)

type ClusterClient struct {
	TxClient  *TxClient
	logPrefix string
}

func (cc ClusterClient) GetRedPacketContractAddr() (string, error) {
	logs := core.BuildLog(core.GetStructFuncName(cc), core.LmChainClient)
	clientCtx := chainnet.ChainNetDst.GetClientCtx()
	resBytes, _, err := util.QueryWithDataWithUnwrapErr(clientCtx, "custom/dao/"+types.QueryRedPacketContractAddr, nil)
	if err != nil {
		logs.WithError(err).Error("QueryWithData")
		return "", err
	}

	return string(resBytes), nil
}

func (cc ClusterClient) QueryRedPacketInfo(redPacketId string) (types.RedPacket, error) {
	logs := core.BuildLog(core.GetStructFuncName(cc), core.LmChainClient)

	req := types.QueryRedPacketInfoParams{
		RedPacketId: redPacketId,
	}

	resp := types.RedPacket{}
	clientCtx := chainnet.ChainNetDst.GetClientCtx()
	reqData, err := clientCtx.LegacyAmino.MarshalJSON(req)
	if err != nil {
		return resp, err
	}

	resBytes, _, err := util.QueryWithDataWithUnwrapErr(clientCtx, "custom/dao/"+types.QueryRedPacketInfo, reqData)
	if err != nil {
		logs.WithError(err).Error("QueryWithData")
		return resp, err
	}
	err = clientCtx.LegacyAmino.UnmarshalJSON(resBytes, &resp)
	if err != nil {
		logs.WithError(err).Error("UnmarshalJSON")
		return resp, err
	}
	return resp, nil
}

func (cc ClusterClient) QueryClusterInfo(clusterId string) (*types.ClusterInfo, error) {
	logs := core.BuildLog(core.GetStructFuncName(cc), core.LmChainClient).WithFields(logrus.Fields{"clusterId": clusterId})
	clientCtx := chainnet.ChainNetDst.GetClientCtx()
	resp := &types.ClusterInfo{}

	param := []byte(clusterId)

	resBytes, _, err := util.QueryWithDataWithUnwrapErr(clientCtx, "custom/dao/"+types.QueryClusterInfo, param)
	if err != nil {
		logs.WithError(err).Error("QueryWithData")
		return nil, err
	}

	if resBytes == nil {
		return nil, nil
	}

	err = clientCtx.LegacyAmino.UnmarshalJSON(resBytes, &resp)
	if err != nil {
		logs.WithError(err).Error("UnmarshalJSON")
		return nil, err
	}
	return resp, nil
}

func (cc ClusterClient) QueryInClusters(fromAddress string) ([]types.InClusters, error) {
	logs := core.BuildLog(core.GetStructFuncName(cc), core.LmChainClient).WithFields(logrus.Fields{"fromAddress": fromAddress})

	resp := make([]types.InClusters, 0)
	clientCtx := chainnet.ChainNetDst.GetClientCtx()
	_, err := sdk.AccAddressFromBech32(fromAddress)
	if err != nil {
		logs.WithError(err).Error("AccAddressFromBech32")
		return resp, core.ErrAddressFormat
	}

	params := []byte(fromAddress)

	resBytes, _, err := util.QueryWithDataWithUnwrapErr(clientCtx, "custom/dao/"+types.QueryInClusters, params)
	if err != nil {
		logs.WithError(err).Error("QueryWithData")
		return resp, err
	}
	err = clientCtx.LegacyAmino.UnmarshalJSON(resBytes, &resp)
	if err != nil {
		logs.WithError(err).Error("UnmarshalJSON")
		return resp, err
	}
	return resp, nil
}

func (cc ClusterClient) QueryClusterGasReward(clusterId, member string) (sdk.DecCoins, error) {
	log := core.BuildLog(core.GetStructFuncName(cc), core.LmChainClient).WithFields(logrus.Fields{"clusterId": clusterId})
	params := types.QueryClusterRewardParams{ClusterId: clusterId, Member: member}
	bz, err := util.Json.Marshal(params)
	if err != nil {
		log.WithError(err).Error("MarshalJSON")
		return sdk.DecCoins{}, err
	}
	clientCtx := chainnet.ChainNetDst.GetClientCtx()
	resBytes, _, err := util.QueryWithDataWithUnwrapErr(clientCtx, "custom/dao/"+types.QueryClusterGasReward, bz)
	if err != nil {
		log.WithError(err).Error("QueryWithData")
		return sdk.DecCoins{}, nil
	}
	res := sdk.DecCoins{}
	err = clientCtx.LegacyAmino.UnmarshalJSON(resBytes, &res)
	if err != nil {
		log.WithError(err).Error("UnmarshalJSON")
		return sdk.DecCoins{}, nil
	}
	return res, nil
}


func (cc ClusterClient) QueryClusterDeviceReward(clusterId, member string) (sdk.DecCoins, error) {
	log := core.BuildLog(core.GetStructFuncName(cc), core.LmChainClient).WithFields(logrus.Fields{"clusterId": clusterId})
	params := types.QueryClusterRewardParams{ClusterId: clusterId, Member: member}
	clientCtx := chainnet.ChainNetDst.GetClientCtx()
	bz, err := util.Json.Marshal(params)
	if err != nil {
		log.WithError(err).Error("MarshalJSON")
		return sdk.DecCoins{}, err
	}
	resBytes, _, err := util.QueryWithDataWithUnwrapErr(clientCtx, "custom/dao/"+types.QueryClusterDeviceReward, bz)
	if err != nil {
		log.WithError(err).Error("QueryWithData")
		return sdk.DecCoins{}, nil
	}
	res := sdk.DecCoins{}
	err = clientCtx.LegacyAmino.UnmarshalJSON(resBytes, &res)
	if err != nil {
		log.WithError(err).Error("UnmarshalJSON")
		return sdk.DecCoins{}, nil
	}
	return res, nil
}

func (cc ClusterClient) QueryPersonClusterInfo(from string) (types.PersonClusterStatisticsInfo, error) {
	log := core.BuildLog(core.GetStructFuncName(cc), core.LmChainClient).WithFields(logrus.Fields{"address": from})
	clientCtx := chainnet.ChainNetDst.GetClientCtx()
	res := types.PersonClusterStatisticsInfo{}
	params := types.QueryPersonClusterInfoRequest{
		From: from,
	}
	bz, err := util.Json.Marshal(params)
	if err != nil {
		log.WithError(err).Error("MarshalJSON")
		return types.PersonClusterStatisticsInfo{}, err
	}

	resBytes, _, err := util.QueryWithDataWithUnwrapErr(clientCtx, "custom/dao/"+types.QueryPersonClusterInfo, bz)
	if err != nil {
		log.WithError(err).Error("QueryWithData")
		return types.PersonClusterStatisticsInfo{}, nil
	}

	err = clientCtx.LegacyAmino.UnmarshalJSON(resBytes, &res)
	if err != nil {
		log.WithError(err).Error("UnmarshalJSON")
		return types.PersonClusterStatisticsInfo{}, nil
	}

	return res, nil
}

func (cc ClusterClient) QueryClusterInfoById(clusterId string) (types.DeviceCluster, error) {
	log := core.BuildLog(core.GetStructFuncName(cc), core.LmChainClient).WithFields(logrus.Fields{"chatClusterId": clusterId})
	clientCtx := chainnet.ChainNetDst.GetClientCtx()
	bz := []byte(clusterId)
	resBytes, _, err := util.QueryWithDataWithUnwrapErr(clientCtx, "custom/dao/"+types.QueryCluster, bz)
	if err != nil {
		log.WithError(err).Error("QueryWithData")
		return types.DeviceCluster{}, err
	}

	res := types.DeviceCluster{}
	err = clientCtx.LegacyAmino.UnmarshalJSON(resBytes, &res)
	if err != nil {
		log.WithError(err).Error("QueryClusterInfoById UnmarshalJSON Error")
		return types.DeviceCluster{}, err
	}

	return res, nil
}

func (cc *ClusterClient) CreateCluster(from string, fee legacytx.StdFee, gatewayAddress, clusterId, clusterName, chatAddress string, salaryRatio, daoRatio, burnAmount sdk.Dec, privateKey string) (tx ttypes.Tx, resp *core.BroadcastTxResponse, err error) {

	log := core.BuildLog(core.GetStructFuncName(cc), core.LmChainClient)

	msg := types.NewMsgCreateCluster(from, gatewayAddress, clusterId, chatAddress, clusterName, salaryRatio, daoRatio, burnAmount, sdk.ZeroDec())
	if err != nil {
		log.Error("NewMsgTransfer")
		return
	}
	var result *core.BaseResponse
	
	tx, result, err = cc.TxClient.SignAndSendMsg(from, privateKey, fee, "txBase.Memo", msg)
	if err != nil {
		return
	}
	resp = new(core.BroadcastTxResponse)
	
	if result.Status == 1 {
		dataByte, err1 := util.Json.Marshal(result.Data)
		if err1 != nil {
			err = err1
			return
		}
		err = util.Json.Unmarshal(dataByte, resp)
		if err != nil {
			return
		}
		return tx, resp, nil
	} else {
		
		resp.TxHash = hex.EncodeToString(tx.Hash())
		return tx, resp, errors.New(result.Info)
	}
}

func (cc *ClusterClient) ChangeClusterId(from string, fee legacytx.StdFee, clusterId, newClusterChatid string, privateKey string) (tx ttypes.Tx, resp *core.BroadcastTxResponse, err error) {

	log := core.BuildLog(core.GetStructFuncName(cc), core.LmChainClient)

	msg := types.NewMsgClusterChangeId(from, clusterId, newClusterChatid)
	if err != nil {
		log.Error("NewMsgTransfer")
		return
	}
	var result *core.BaseResponse
	
	tx, result, err = cc.TxClient.SignAndSendMsg(from, privateKey, fee, "txBase.Memo", msg)
	if err != nil {
		return
	}
	resp = new(core.BroadcastTxResponse)
	
	if result.Status == 1 {
		dataByte, err1 := util.Json.Marshal(result.Data)
		if err1 != nil {
			err = err1
			return
		}
		err = util.Json.Unmarshal(dataByte, resp)
		if err != nil {
			return
		}
		return tx, resp, nil
	} else {
		
		resp.TxHash = hex.EncodeToString(tx.Hash())
		return tx, resp, errors.New(result.Info)
	}
}

func (cc ClusterClient) QueryDaoParams() (types.DaoParams, error) {
	logs := core.BuildLog(core.GetStructFuncName(cc), core.LmChainClient)

	resp := types.DaoParams{}
	clientCtx := chainnet.ChainNetDst.GetClientCtx()
	resBytes, _, err := util.QueryWithDataWithUnwrapErr(clientCtx, "custom/dao/"+types.QueryDaoParams, nil)
	if err != nil {
		logs.WithError(err).Error("QueryWithData")
		return resp, err
	}
	err = clientCtx.LegacyAmino.UnmarshalJSON(resBytes, &resp)
	if err != nil {
		logs.WithError(err).Error("UnmarshalJSON")
		return resp, err
	}
	return resp, nil
}

func (cc ClusterClient) QueryAllParams() (types.Params, error) {
	logs := core.BuildLog(core.GetStructFuncName(cc), core.LmChainClient)
	clientCtx := chainnet.ChainNetDst.GetClientCtx()
	resp := types.Params{}

	resBytes, _, err := util.QueryWithDataWithUnwrapErr(clientCtx, "custom/dao/"+types.QueryAllParams, nil)
	if err != nil {
		logs.WithError(err).Error("QueryWithData")
		return resp, err
	}
	err = clientCtx.LegacyAmino.UnmarshalJSON(resBytes, &resp)
	if err != nil {
		logs.WithError(err).Error("UnmarshalJSON")
		return resp, err
	}
	return resp, nil
}

func (cc ClusterClient) QueryClusterPersonInfo(clusterId, fromAddress string) (types.ClusterPersonalInfo, error) {
	logs := core.BuildLog(core.GetStructFuncName(cc), core.LmChainClient)

	req := types.QueryClusterPersonalInfoParams{
		ClusterId:   clusterId,
		FromAddress: fromAddress,
	}
	clientCtx := chainnet.ChainNetDst.GetClientCtx()
	resp := types.ClusterPersonalInfo{}

	reqData, err := util.Json.Marshal(req)
	if err != nil {
		return resp, err
	}

	resBytes, _, err := util.QueryWithDataWithUnwrapErr(clientCtx, "custom/dao/"+types.QueryClusterPersonInfo, reqData)
	if err != nil {
		logs.WithError(err).Error("QueryWithData")
		return resp, err
	}
	err = clientCtx.LegacyAmino.UnmarshalJSON(resBytes, &resp)
	if err != nil {
		logs.WithError(err).Error("UnmarshalJSON")
		return resp, err
	}
	return resp, nil
}

func (cc ClusterClient) QueryDaoStatistic() (types.QueryDaoStatisticResp, error) {
	logs := core.BuildLog(core.GetStructFuncName(cc), core.LmChainClient)
	resp := types.QueryDaoStatisticResp{}
	clientCtx := chainnet.ChainNetDst.GetClientCtx()
	resBytes, _, err := util.QueryWithDataWithUnwrapErr(clientCtx, "custom/dao/"+types.QueryDaoStatistic, nil)
	if err != nil {
		logs.WithError(err).Error("QueryWithData")
		return resp, err
	}
	err = clientCtx.LegacyAmino.UnmarshalJSON(resBytes, &resp)
	if err != nil {
		logs.WithError(err).Error("UnmarshalJSON")
		return resp, err
	}
	return resp, nil
}

func (cc ClusterClient) QueryDeviceNoDvm(clusterId string, addrs []string) ([]string, error) {
	logs := core.BuildLog(core.GetStructFuncName(cc), core.LmChainClient)
	resp := make([]string, 0)
	clientCtx := chainnet.ChainNetDst.GetClientCtx()
	params := types.QueryNoDvmParams{
		ClusterId: clusterId,
		Addrs:     addrs,
	}

	reqData, err := clientCtx.LegacyAmino.MarshalJSON(params)
	if err != nil {
		return resp, err
	}

	resBytes, _, err := util.QueryWithDataWithUnwrapErr(clientCtx, "custom/dao/"+types.QueryNoDvm, reqData)
	if err != nil {
		logs.WithError(err).Error("QueryWithData")
		return resp, err
	}
	err = clientCtx.LegacyAmino.UnmarshalJSON(resBytes, &resp)
	if err != nil {
		logs.WithError(err).Error("UnmarshalJSON")
		return resp, err
	}
	return resp, nil
}

func (cc ClusterClient) QueryDeviceNoDvmFormatChat(clusterId string, addrs []string) ([]string, error) {
	logs := core.BuildLog(core.GetStructFuncName(cc), core.LmChainClient)
	resp := make([]string, 0)
	clientCtx := chainnet.ChainNetDst.GetClientCtx()
	params := types.QueryNoDvmParams{
		ClusterId: clusterId,
		Addrs:     addrs,
	}

	reqData, err := clientCtx.LegacyAmino.MarshalJSON(params)
	if err != nil {
		return resp, err
	}

	resBytes, _, err := util.QueryWithDataWithUnwrapErr(clientCtx, "custom/dao/"+types.QueryNoDvmChat, reqData)
	if err != nil {
		logs.WithError(err).Error("QueryWithData")
		return resp, err
	}
	err = clientCtx.LegacyAmino.UnmarshalJSON(resBytes, &resp)
	if err != nil {
		logs.WithError(err).Error("UnmarshalJSON")
		return resp, err
	}
	return resp, nil
}

func (cc ClusterClient) QueryDvmList(addr string) ([]types.DvmInfo, error) {
	logs := core.BuildLog(core.GetStructFuncName(cc), core.LmChainClient)
	resp := make([]types.DvmInfo, 0)

	params := []byte(addr)
	clientCtx := chainnet.ChainNetDst.GetClientCtx()
	resBytes, _, err := util.QueryWithDataWithUnwrapErr(clientCtx, "custom/dao/"+types.QueryDvmList, params)
	if err != nil {
		logs.WithError(err).Error("QueryWithData")
		return resp, err
	}
	err = clientCtx.LegacyAmino.UnmarshalJSON(resBytes, &resp)
	if err != nil {
		logs.WithError(err).Error("UnmarshalJSON")
		return resp, err
	}
	return resp, nil
}

func (cc *ClusterClient) RedPacket(from string, fee legacytx.StdFee, privateKey, clusterChatId string, amount sdk.Coin, count int64, redPacketType int64) (tx ttypes.Tx, resp *core.BroadcastTxResponse, err error) {

	log := core.BuildLog(core.GetStructFuncName(cc), core.LmChainClient)

	msg := types.NewMsgRedPacket(from, clusterChatId, amount, count, redPacketType)
	if err != nil {
		log.Error("NewMsgTransfer")
		return
	}
	var result *core.BaseResponse
	
	tx, result, err = cc.TxClient.SignAndSendMsg(from, privateKey, fee, "txBase.Memo", msg)
	if err != nil {
		return
	}
	resp = new(core.BroadcastTxResponse)
	
	if result.Status == 1 {
		dataByte, err1 := util.Json.Marshal(result.Data)
		if err1 != nil {
			err = err1
			return
		}
		err = util.Json.Unmarshal(dataByte, resp)
		if err != nil {
			return
		}
		return tx, resp, nil
	} else {
		
		resp.TxHash = hex.EncodeToString(tx.Hash())
		return tx, resp, errors.New(result.Info)
	}
}

func (cc *ClusterClient) OpenRedPacket(from string, fee legacytx.StdFee, privateKey, redPacketId string) (tx ttypes.Tx, resp *core.BroadcastTxResponse, err error) {

	log := core.BuildLog(core.GetStructFuncName(cc), core.LmChainClient)

	msg := types.NewMsgOpenRedPacket(from, redPacketId)
	if err != nil {
		log.Error("NewMsgTransfer")
		return
	}
	var result *core.BaseResponse
	
	tx, result, err = cc.TxClient.SignAndSendMsg(from, privateKey, fee, "txBase.Memo", msg)
	if err != nil {
		return
	}
	resp = new(core.BroadcastTxResponse)
	
	if result.Status == 1 {
		dataByte, err1 := util.Json.Marshal(result.Data)
		if err1 != nil {
			err = err1
			return
		}
		err = util.Json.Unmarshal(dataByte, resp)
		if err != nil {
			return
		}
		return tx, resp, nil
	} else {
		
		resp.TxHash = hex.EncodeToString(tx.Hash())
		return tx, resp, errors.New(result.Info)
	}
}

func (cc *ClusterClient) ReturnRedPacket(from string, fee legacytx.StdFee, privateKey, redPacketId string) (tx ttypes.Tx, resp *core.BroadcastTxResponse, err error) {

	log := core.BuildLog(core.GetStructFuncName(cc), core.LmChainClient)

	msg := types.NewMsgReturnRedPacket(from, redPacketId)
	if err != nil {
		log.Error("NewMsgTransfer")
		return
	}
	var result *core.BaseResponse
	
	tx, result, err = cc.TxClient.SignAndSendMsg(from, privateKey, fee, "txBase.Memo", msg)
	if err != nil {
		return
	}
	resp = new(core.BroadcastTxResponse)
	
	if result.Status == 1 {
		dataByte, err1 := util.Json.Marshal(result.Data)
		if err1 != nil {
			err = err1
			return
		}
		err = util.Json.Unmarshal(dataByte, resp)
		if err != nil {
			return
		}
		return tx, resp, nil
	} else {
		
		resp.TxHash = hex.EncodeToString(tx.Hash())
		return tx, resp, errors.New(result.Info)
	}
}

func (cc *ClusterClient) GetCutRewardInfo(from string) (types.GetCutRewardInfoResp, error) {

	logs := core.BuildLog(core.GetStructFuncName(cc), core.LmChainClient)

	params := []byte(from)

	queryInfo := types.GetCutRewardInfoResp{}

	clientCtx := chainnet.ChainNetDst.GetClientCtx()
	resBytes, _, err := util.QueryWithDataWithUnwrapErr(clientCtx, "custom/dao/"+types.QueryCutRewards, params)
	if err != nil {
		logs.WithError(err).Error("QueryWithData")
		return types.GetCutRewardInfoResp{}, err
	}
	err = clientCtx.LegacyAmino.UnmarshalJSON(resBytes, &queryInfo)
	if err != nil {
		logs.WithError(err).Error("UnmarshalJSON")
		return types.GetCutRewardInfoResp{}, err
	}

	return queryInfo, nil

}

func (cc *ClusterClient) QueryMintStatus() (types.QueryMiningStatusResp, error) {
	logs := core.BuildLog(core.GetStructFuncName(cc), core.LmChainClient)
	res := types.QueryMiningStatusResp{}
	clientCtx := chainnet.ChainNetDst.GetClientCtx()
	resBytes, _, err := util.QueryWithDataWithUnwrapErr(clientCtx, "custom/dao/"+types.QueryMiningStatus, nil)
	if err != nil {
		logs.WithError(err).Error("QueryWithData")
		return types.QueryMiningStatusResp{}, err
	}
	err = clientCtx.LegacyAmino.UnmarshalJSON(resBytes, &res)
	if err != nil {
		logs.WithError(err).Error("UnmarshalJSON")
		return types.QueryMiningStatusResp{}, err
	}
	return res, nil
}


func (cc *ClusterClient) QueryClusterDeposit() ([]sdk.Coin, error) {
	logs := core.BuildLog(core.GetStructFuncName(cc), core.LmChainClient)

	res := []sdk.Coin{}
	clientCtx := chainnet.ChainNetDst.GetClientCtx()
	resBytes, _, err := util.QueryWithDataWithUnwrapErr(clientCtx, "custom/dao/"+types.QueryGroupDeposit, nil)
	if err != nil {
		logs.WithError(err).Error("QueryWithData")
		return []sdk.Coin{}, err
	}
	err = clientCtx.LegacyAmino.UnmarshalJSON(resBytes, &res)
	if err != nil {
		logs.WithError(err).Error("UnmarshalJSON")
		return []sdk.Coin{}, err
	}

	return res, nil
}


func (cc *ClusterClient) QueryAllClusterProposals(clusterid, voter string) ([]types.ClusterProposals, error) {
	logs := core.BuildLog(core.GetStructFuncName(cc), core.LmChainClient)

	params := types.QueryClusterProposalsParams{
		ClusterId: clusterid,
		Voter:     voter,
	}
	clientCtx := chainnet.ChainNetDst.GetClientCtx()
	paramsJson, err := clientCtx.LegacyAmino.MarshalJSON(params)
	if err != nil {
		return []types.ClusterProposals{}, err
	}

	resBytes, _, err := util.QueryWithDataWithUnwrapErr(clientCtx, "custom/dao/"+types.QueryClusterPersonals, paramsJson)
	if err != nil {
		logs.WithError(err).Error("QueryWithData")
		return []types.ClusterProposals{}, err
	}
	res := make([]types.ClusterProposals, 0)
	err = json.Unmarshal(resBytes, &res)
	if err != nil {
		logs.WithError(err).Error("UnmarshalJSON")
		return []types.ClusterProposals{}, err
	}

	return res, nil
}


func (cc *ClusterClient) QueryClusterProposalInfo(id string) (group.ProposalNew, error) {
	logs := core.BuildLog(core.GetStructFuncName(cc), core.LmChainClient)
	clientCtx := chainnet.ChainNetDst.GetClientCtx()
	params := []byte(id)

	resBytes, _, err := util.QueryWithDataWithUnwrapErr(clientCtx, "custom/dao/"+types.QueryClusterPersonalInfo, params)
	if err != nil {
		logs.WithError(err).Error("QueryWithData")
		return group.ProposalNew{}, err
	}
	res := group.ProposalNew{}
	err = json.Unmarshal(resBytes, &res)
	if err != nil {
		logs.WithError(err).Error("UnmarshalJSON")
		return group.ProposalNew{}, err
	}

	return res, nil

}
