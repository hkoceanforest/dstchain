package keeper

import (
	sdkmath "cosmossdk.io/math"
	"encoding/json"
	"errors"
	"freemasonry.cc/blockchain/core"
	"freemasonry.cc/blockchain/util"
	"freemasonry.cc/blockchain/x/dao/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/tx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/group"
	"github.com/shopspring/decimal"
	abci "github.com/tendermint/tendermint/abci/types"
	"strconv"
	"strings"
)

func NewQuerier(k Keeper, legacyQuerierCdc *codec.LegacyAmino) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		var (
			res []byte
			err error
		)
		switch path[0] {
		case types.QueryBurnLevels:
			return queryBurnLevels(ctx, req, k, legacyQuerierCdc)
		case types.QueryBurnLevel:
			return queryBurnLevel(ctx, req, k, legacyQuerierCdc)
		case types.QueryPersonClusterInfo:
			return queryPersonClusterInfo(ctx, req, k, legacyQuerierCdc)
		case types.QueryClusterInfoById:
			return queryClusterInfoById(ctx, req, k, legacyQuerierCdc)
		case types.QueryGatewayClusters:
			return queryGatewayClusters(ctx, req, k, legacyQuerierCdc)
		case types.QueryClusterGasReward:
			return queryClustersGasReward(ctx, req, k, legacyQuerierCdc)
		case types.QueryClusterOwnerReward:
			return queryClustersOwnerReward(ctx, req, k, legacyQuerierCdc)
		case types.QueryClusterDeviceReward:
			return queryClustersDeviceReward(ctx, req, k, legacyQuerierCdc)
		case types.QueryDeductionFee:
			return queryClustersDeductionFee(ctx, req, k, legacyQuerierCdc)
		case types.QueryInClusters:
			return queryInClusters(ctx, req, k, legacyQuerierCdc)
		case types.QueryClusterInfo:
			return queryClusterInfo(ctx, req, k, legacyQuerierCdc)
		case types.QueryClusterPersonInfo:
			return queryClusterPersonInfo(ctx, req, k, legacyQuerierCdc)
		case types.QueryCluster:
			return queryCluster(ctx, req, k, legacyQuerierCdc)
		case types.QueryClusterPersonals:
			return queryClusterProposals(ctx, req, k, legacyQuerierCdc)
		case types.QueryClusterPersonalInfo:
			return queryClusterPersonalInfo(ctx, req, k, legacyQuerierCdc)
		case types.QueryClusterProposalVoters:
			return queryClusterProposalVoters(ctx, req, k, legacyQuerierCdc)
		case types.QueryClusterProposalVoter:
			return queryClusterProposalVoter(ctx, req, k, legacyQuerierCdc)
		case types.QueryGroupMembers:
			return queryGroupMembers(ctx, req, k, legacyQuerierCdc)
		case types.QueryGroupInfo:
			return queryGroupInfo(ctx, req, k, legacyQuerierCdc)
		case types.QueryDvmList:
			return queryDvmList(ctx, req, k, legacyQuerierCdc)
		case types.QueryDaoParams:
			return QueryDaoParams(ctx, req, k, legacyQuerierCdc)
		case types.QueryClusterApproveInfo:
			return queryClusterApproveInfo(ctx, req, k, legacyQuerierCdc)
		case types.QueryClusterRelation:
			return queryClusterRelation(ctx, req, k, legacyQuerierCdc)
		case types.QueryDaoStatistic:
			return queryDaoStatistic(ctx, req, k, legacyQuerierCdc)
		case types.QueryNoDvm:
			return queryNoDvm(ctx, req, k, legacyQuerierCdc)
		case types.QueryNoDvmChat:
			return QueryNoDvmChat(ctx, req, k, legacyQuerierCdc)
		case types.QueryClusterProposalTallyResult:
			return queryClusterProposalTallyResult(ctx, req, k, legacyQuerierCdc)
		case types.QueryClusterAd:
			return queryClusterAd(ctx, req, k, legacyQuerierCdc)
		case types.QueryAllParams:
			return queryAllParams(ctx, req, k, legacyQuerierCdc)
		case types.QueryRedPacketInfo:
			return queryRedPacketInfo(ctx, req, k, legacyQuerierCdc)
		case types.QueryRedPacketContractAddr:
			return queryRedPacketContractAddr(ctx, req, k, legacyQuerierCdc)
		case types.QueryNextAccountNumber:
			return queryNextAccountNumber(ctx, req, k, legacyQuerierCdc)
		case types.QueryGroupDeposit:
			return queryGroupDeposit(ctx, req, k, legacyQuerierCdc)
		case types.QueryGroupVotesByVoter:
			return queryGroupVotesByVoter(ctx, req, k, legacyQuerierCdc)
		case types.QueryDaoQueue:
			return queryDaoQueue(ctx, req, k, legacyQuerierCdc)
		case types.QueryAllClusterReward:
			return queryAllClustersReward(ctx, req, k, legacyQuerierCdc)
		case types.QueryCutRewards:
			return queryCutRewards(ctx, req, k, legacyQuerierCdc)
		case types.QuerySupplyStatistics:
			return querySupplyStatistics(ctx, req, k, legacyQuerierCdc)
		case types.QueryGasDeduction:
			return queryGasDeduction(ctx, req, k, legacyQuerierCdc)
		case types.QueryMiningStatus:
			return queryMiningStatus(ctx, req, k, legacyQuerierCdc)
		case types.QueryClusterReceiveDao:
			return queryClusterReceiveDao(ctx, req, k, legacyQuerierCdc)
		default:
			err = sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unknown %s query endpoint: %s", types.ModuleName, path[0])
		}

		return res, err
	}
}

func queryClusterReceiveDao(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	log := core.BuildLog(core.GetPackageFuncName(), core.LmChainDaoQuery)
	var params types.QueryClusterParams

	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		log.WithError(err).Error("UnmarshalJSON")
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}
	clusterId := params.ClusterId
	if strings.Contains(clusterId, ".") {
		
		clusterId, err = k.GetClusterId(ctx, clusterId)
		if err != nil {
			return nil, err
		}
	}
	resp, err := k.GetClusterDaoRewardSum(ctx, clusterId)
	if err != nil {
		return nil, err
	}
	bz, err := codec.MarshalJSONIndent(legacyQuerierCdc, resp)
	if err != nil {
		return nil, err
	}
	return bz, nil
}

func queryMiningStatus(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	miningStatus := k.StartMint(ctx)
	idoEndMark := k.GetGenesisIdoEndMark(ctx)

	resp := types.QueryMiningStatusResp{
		MintStart:  miningStatus,
		IdoEndMark: idoEndMark,
	}
	bz, err := codec.MarshalJSONIndent(legacyQuerierCdc, resp)
	if err != nil {
		return nil, err
	}
	return bz, nil
}

func queryGasDeduction(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	log := core.BuildLog(core.GetPackageFuncName(), core.LmChainDaoQuery)
	var params types.GasDeductionParams

	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		log.WithError(err).Error("UnmarshalJSON")
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}
	payFee := sdk.NewCoins(sdk.NewCoin(core.BaseDenom, sdk.ZeroInt()))
	for _, msg := range params.Msg {
		calFee := sdk.NewCoins(sdk.NewCoin(core.BaseDenom, sdk.ZeroInt()))
		if obj, oks := msg.(*bankTypes.MsgSend); oks {
			calFee, err = k.CalculateSendFee(ctx, *obj, params.Fee)
			if err != nil {
				return nil, err
			}
		} else {
			calFee, err = k.CalculateFee(ctx, msg, params.Fee)
			if err != nil {
				return nil, err
			}
		}

		if calFee != nil && !calFee.IsZero() {
			payFee = payFee.Add(calFee...)
		}
	}
	bz, err := codec.MarshalJSONIndent(legacyQuerierCdc, payFee)
	if err != nil {
		return nil, err
	}
	return bz, nil
}

func querySupplyStatistics(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	log := core.BuildLog(core.GetPackageFuncName(), core.LmChainDaoQuery)
	mintSupply, err := k.GetMintSupply(ctx)
	if err != nil {
		log.WithError(err).Error("GetMintSupply")
		return nil, err
	}
	burnSupply, err := k.GetBurnSupply(ctx)
	if err != nil {
		log.WithError(err).Error("burnSupply")
		return nil, err
	}
	burnSupply = burnSupply.Add(core.BurnRepair)
	noMintSupply := core.DstSupply.TruncateInt().Sub(mintSupply)

	swapSupply, err := k.GetSwapDelegateSupply(ctx)
	if err != nil {
		log.WithError(err).Error("GetSwapDelegateSupply")
		return nil, err
	}

	resp := types.QuerySupplyStatisticsResp{
		BurnSupply:      burnSupply,
		NoMintSupply:    noMintSupply,
		CirculateSupply: core.DstSupply.TruncateInt().Sub(burnSupply).Sub(noMintSupply).Sub(swapSupply),
		SwapSupply:      swapSupply,
	}
	bz, err := codec.MarshalJSONIndent(legacyQuerierCdc, resp)
	if err != nil {
		return nil, err
	}
	return bz, nil
}

func queryCutRewards(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	log := core.BuildLog(core.GetPackageFuncName(), core.LmChainDaoQuery)
	addr := string(req.Data)

	cutInfo, err := k.GetPowerRewardCycleInfo(ctx, addr)
	if err != nil {
		log.WithError(err).Error("GetPowerRewardCycleInfo")
		return nil, err
	}

	
	isStart := false

	
	allRedemption := sdk.ZeroInt()
	
	allReceive := sdk.ZeroInt()
	
	lastCycle := int64(0)
	
	remainDays := int64(0)

	nowTime := ctx.BlockTime().Unix()

	for _, info := range cutInfo.CycleInfo {
		if info.Cycle > lastCycle {
			lastCycle = info.Cycle
		}

		
		if info.StartTime+86400*core.CutPowerRewardTimes > nowTime && info.Status == 1 {
			
			redemptionDays := (nowTime - info.StartTime) / 86400
			allRedemption = allRedemption.Add(info.CutPerReward.MulRaw(core.CutPowerRewardTimes - redemptionDays))
			isStart = true
			remainDays = core.CutPowerRewardTimes - info.ReceiveTimes
		}

		if info.Status == 1 {
			diff := nowTime - info.StartTime
			shouldReceiveTimes := diff / 86400 
			if shouldReceiveTimes >= core.CutPowerRewardTimes {
				shouldReceiveTimes = core.CutPowerRewardTimes
			}

			receivedTimes := info.ReceiveTimes 

			receiveTimes := shouldReceiveTimes - receivedTimes 

			receiveRewardAmount := info.CutPerReward.Mul(sdk.NewInt(receiveTimes)) 

			allReceive = allReceive.Add(receiveRewardAmount)
		}

	}

	
	cycle := (nowTime - k.GetStartTime(ctx)) / core.CutProductionSeconds

	nextCycleStartTime := k.GetStartTime(ctx) + (cycle+1)*core.CutProductionSeconds - nowTime

	resp := types.GetCutRewardInfoResp{
		IsStart:          isStart,
		RemainDays:       remainDays,
		NextStartTime:    nextCycleStartTime,
		RedemptionAmount: allRedemption,
		ReceiveAmount:    allReceive,
	}

	bz, err := codec.MarshalJSONIndent(legacyQuerierCdc, resp)
	if err != nil {
		return nil, err
	}
	return bz, nil
}

func queryDaoQueue(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	log := core.BuildLog(core.GetPackageFuncName(), core.LmChainDaoQuery)
	var params types.QueryDaoQueueParams

	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		log.WithError(err).Error("UnmarshalJSON")
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}
	clusterId := params.ClusterId
	if strings.Contains(clusterId, ".") {
		
		clusterId, err = k.GetClusterId(ctx, clusterId)
		if err != nil {
			return nil, err
		}
	}
	start := (params.Page - 1) * params.PageSize
	end := params.Page * params.PageSize
	queue, err := k.GetClusterDaoRewardQueue(ctx, clusterId)
	if err != nil {
		return nil, err
	}
	
	if params.Member != "" {
		newQueue := make([]types.ClusterMemberDaoReward, 0)
		for _, item := range queue {
			if item.Address == params.Member {
				newQueue = append(newQueue, item)
			}
		}
		queue = newQueue
	}
	if start >= len(queue) {
		return nil, nil
	}
	if end > len(queue) {
		end = len(queue)
	}
	param := k.GetParams(ctx)
	queue = queue[start:end]
	resp := types.QueryDaoQueueResp{
		Queues: queue,
		Total:  len(queue),
		Rate:   param.ReceiveDaoRatio,
	}
	bz, err := codec.MarshalJSONIndent(legacyQuerierCdc, resp)
	if err != nil {
		return nil, err
	}
	return bz, nil
}

func queryGroupDeposit(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	params := k.groupKeeper.GetParams(ctx)
	depositByte, err := codec.MarshalJSONIndent(legacyQuerierCdc, params.Deposit)
	if err != nil {
		return nil, err
	}
	return depositByte, nil
}

func queryNextAccountNumber(ctx sdk.Context, _ abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	newAcNum := k.accountKeeper.GetNextAccountNumber(ctx)
	newAcNumByte, err := codec.MarshalJSONIndent(legacyQuerierCdc, newAcNum)
	if err != nil {
		return nil, err
	}
	return newAcNumByte, nil
}


func queryRedPacketContractAddr(ctx sdk.Context, _ abci.RequestQuery, k Keeper, _ *codec.LegacyAmino) ([]byte, error) {
	redPacketContractAddr := k.contractKeeper.GetRedPacketContractAddress(ctx)
	return []byte(redPacketContractAddr), nil
}


func queryRedPacketInfo(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	log := core.BuildLog(core.GetPackageFuncName(), core.LmChainDaoQuery)
	var params types.QueryRedPacketInfoParams

	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		log.WithError(err).Error("UnmarshalJSON")
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	redPacketInfo, err := k.GetRedPacketInfo(ctx, params.RedPacketId)
	if err != nil {
		return nil, err
	}

	bz, err := codec.MarshalJSONIndent(legacyQuerierCdc, redPacketInfo)
	if err != nil {
		return nil, err
	}
	return bz, nil
}

func queryAllParams(ctx sdk.Context, _ abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	params := k.GetParams(ctx)

	resp, err := codec.MarshalJSONIndent(legacyQuerierCdc, params)
	if err != nil {
		return nil, err
	}

	return resp, nil
}


func queryClusterAd(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	log := core.BuildLog(core.GetPackageFuncName(), core.LmChainDaoQuery)
	var params types.QueryClusterAdParams

	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		log.WithError(err).Error("UnmarshalJSON")
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}
	startTime, err := strconv.ParseInt(params.StartTime, 10, 64)
	if err != nil {
		return nil, err
	}
	endTime, err := strconv.ParseInt(params.EndTime, 10, 64)
	if err != nil {
		return nil, err
	}
	level, err := strconv.ParseInt(params.Level, 10, 64)
	if err != nil {
		return nil, err
	}
	clusters := k.ClustersFilter(ctx, startTime, endTime, level, params.Rate)
	param := k.GetParams(ctx)
	rate, err := k.GetExchangeRate(ctx, authtypes.NewModuleAddress(types.ModuleName), param)
	if err != nil {
		return nil, err
	}
	amount := sdkmath.ZeroInt()
	resp := types.QueryClusterAdResp{}
	quantity := 0
	for _, cluster := range clusters {
		resp.ClusterId = append(resp.ClusterId, cluster.ClusterChatId)
		adAmount := rate.MulRaw(int64(len(cluster.ClusterDeviceMembers)))
		quantity = quantity + len(cluster.ClusterDeviceMembers)
		amount = amount.Add(adAmount)
	}
	resp.Amount = amount.String()
	resp.Quantity = quantity
	bz, err := codec.MarshalJSONIndent(legacyQuerierCdc, resp)
	if err != nil {
		return nil, err
	}
	return bz, nil
}

func QueryNoDvmChat(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	log := core.BuildLog(core.GetPackageFuncName(), core.LmChainDaoQuery)
	var params types.QueryNoDvmParams

	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		log.WithError(err).Error("UnmarshalJSON")
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	addrs, err := k.GetNoDvm(ctx, params.ClusterId, params.Addrs)
	if err != nil {
		return nil, err
	}

	res := make([]string, 0)
	for _, addr := range addrs {
		
		resInfo, err := k.ChatKeeper.GetRegisterInfo(ctx, addr)
		if err != nil {
			continue
		}

		gatewayInfo, err := k.gatewayKeeper.GetGatewayInfo(ctx, resInfo.NodeAddress)
		if err != nil {
			continue
		}

		if len(gatewayInfo.GatewayNum) > 0 {
			firstNum := ""
			for _, index := range gatewayInfo.GatewayNum {
				if index.IsFirst == true {
					firstNum = index.NumberIndex
				}
			}

			if firstNum == "" {
				continue
			}
			res = append(res, "@"+addr+":"+firstNum+"."+core.GovDenom)
		} else {
			continue
		}
	}

	bz, err := codec.MarshalJSONIndent(legacyQuerierCdc, res)
	if err != nil {
		return nil, err
	}
	return bz, nil
}

func queryNoDvm(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	log := core.BuildLog(core.GetPackageFuncName(), core.LmChainDaoQuery)
	var params types.QueryNoDvmParams

	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		log.WithError(err).Error("UnmarshalJSON")
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	res, err := k.GetNoDvm(ctx, params.ClusterId, params.Addrs)
	if err != nil {
		return nil, err
	}

	bz, err := codec.MarshalJSONIndent(legacyQuerierCdc, res)
	if err != nil {
		return nil, err
	}
	return bz, nil
}

func queryDaoStatistic(ctx sdk.Context, _ abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {

	resp := types.QueryDaoStatisticResp{}
	resp.Statistic = make([]types.DaoStatistic, 0)
	resp.ClusterIds = make([]string, 0)
	
	clusters := k.GetAllClusters(ctx)
	for _, cluster := range clusters {

		
		onlineRatio := sdk.NewDec(cluster.ClusterActiveDevice).Quo(sdk.NewDec(int64(len(cluster.ClusterDeviceMembers))))

		
		connectivityRate := k.GetClusterConnectivityRate(ctx, cluster.ClusterId)

		
		deviceNoEvmAmount := k.GetDeviceNoEvm(cluster)

		resp.Statistic = append(resp.Statistic, types.DaoStatistic{
			ClusterId:              cluster.ClusterId,
			ClusterChatId:          cluster.ClusterChatId,
			DeviceOnlienRatio:      onlineRatio,
			DeviceConnectivityRate: connectivityRate,
			DeviceAmount:           int64(len(cluster.ClusterDeviceMembers)),
			DvmAmount:              int64(len(cluster.ClusterPowerMembers)),
			DeviceNoEvm:            deviceNoEvmAmount,
			ClusterBurnAmount:      cluster.ClusterBurnAmount,
		})

		resp.ClusterIds = append(resp.ClusterIds, cluster.ClusterChatId)
	}

	bz, err := codec.MarshalJSONIndent(legacyQuerierCdc, resp)
	if err != nil {
		return nil, err
	}
	return bz, nil
}

func queryClusterRelation(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	log := core.BuildLog(core.GetPackageFuncName(), core.LmChainDaoQuery)
	var params types.QueryClusterParams

	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		log.WithError(err).Error("UnmarshalJSON")
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}
	clusterId := params.ClusterId

	if strings.Contains(clusterId, ".") {
		
		clusterId, err = k.GetClusterId(ctx, clusterId)
		if err != nil {
			return nil, err
		}
	}
	cluster, err := k.GetCluster(ctx, clusterId)
	if err != nil {
		return nil, err
	}
	if cluster.ClusterPowerMembers == nil {
		return nil, nil
	}
	resp := []types.QueryClusterRelationResp{}
	for _, member := range cluster.ClusterPowerMembers {
		person, err := k.GetPersonClusterInfo(ctx, member.Address)
		if err != nil {
			return nil, err
		}
		if person.Owner != nil && len(person.Owner) > 0 {
			
			userInfo, err := k.ChatKeeper.GetRegisterInfo(ctx, member.Address)
			if err != nil {
				return nil, err
			}
			gateway, err := k.gatewayKeeper.GetGatewayInfo(ctx, userInfo.NodeAddress)
			if err != nil {
				return nil, err
			}
			indexNum := ""
			for _, index := range gateway.GatewayNum {
				if index.IsFirst {
					indexNum = index.NumberIndex
					break
				}
			}
			relation := types.QueryClusterRelationResp{Address: member.Address, IndexNum: indexNum}
			resp = append(resp, relation)
		}
	}
	bz, err := codec.MarshalJSONIndent(legacyQuerierCdc, resp)
	if err != nil {
		return nil, err
	}
	return bz, nil
}

func queryBurnLevels(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	log := core.BuildLog(core.GetPackageFuncName(), core.LmChainDaoQuery)
	var params types.QueryBurnLevelsParams

	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		log.WithError(err).Error("UnmarshalJSON")
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	burnLevels, err := k.GetBurnLevelByAccAddresses(ctx, params.Addresses)
	if err != nil {
		log.WithError(err).Info("QueryBurnLevels Err")
		return nil, err
	}

	burnLevelsByte, err := codec.MarshalJSONIndent(legacyQuerierCdc, burnLevels)
	if err != nil {
		log.WithError(err).Info("QueryBurnLevels Marshal Err")
		return nil, err
	}

	if burnLevelsByte == nil {
		log.Warning("QueryBurnLevels Not Fount")
	}

	return burnLevelsByte, nil
}

func queryBurnLevel(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	logs := core.BuildLog(core.GetPackageFuncName(), core.LmChainDaoQuery)

	addr := string(req.Data)

	accAddr, err := sdk.AccAddressFromBech32(addr)
	if err != nil {
		return nil, core.ErrAddressFormat
	}

	burnLevelInfo, err := k.GetBurnLevelByAccAddress(ctx, accAddr)
	if err != nil {
		logs.WithError(err).Error("queryBurnLevel GetBurnLevelByAccAddress err")
		return nil, err
	}

	burnLevelsByte, err := codec.MarshalJSONIndent(legacyQuerierCdc, burnLevelInfo)
	if err != nil {
		logs.WithError(err).Info("QueryBurnLevel MarshalJSONIndent Err")
		return nil, err
	}

	if burnLevelsByte == nil {
		logs.Warning("queryBurnLevel Not Fount")
	}

	return burnLevelsByte, nil
}

func queryDvmList(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	from := string(req.Data)

	_, err := sdk.AccAddressFromBech32(from)
	if err != nil {
		return nil, core.ErrAddressFormat
	}

	
	pInfo, err := k.GetPersonClusterInfo(ctx, from)
	if err != nil {
		return nil, err
	}

	res := make([]types.DvmInfo, 0)

	
	for clusterId := range pInfo.BePower {
		cluster, err := k.GetCluster(ctx, clusterId)
		if err != nil {
			if err == core.GetClusterNotFound {
				continue
			}
			return nil, err
		}

		var powerDvm sdk.Dec
		powerReward := sdk.ZeroDec()
		powerInfo, ok := cluster.ClusterPowerMembers[from]
		if !ok {
			powerDvm = sdk.ZeroDec()
		} else {
			powerDvm = powerInfo.ActivePower
			endingPeriod := k.IncrementClusterPeriod(ctx, cluster)
			rewards := k.CalculateBurnRewards(ctx, cluster, from, endingPeriod)
			
			powerReward = rewards.AmountOf(core.BaseDenom).Mul(sdk.OneDec().Sub(core.ClusterDeviceRate))
			
			hisReward, err := k.GetClusterMemberReward(ctx, clusterId, from)
			if err != nil {
				return nil, err
			}
			
			if !hisReward.DeviceReward.IsZero() {
				powerReward = powerReward.Add(sdk.NewDecFromInt(hisReward.DeviceReward).Mul(sdk.OneDec().Sub(core.ClusterDeviceRate)))
			}
			if !hisReward.HisReward.IsZero() {
				powerReward = powerReward.Add(sdk.NewDecFromInt(hisReward.HisReward))
			}
		}
		approveInfo, err := k.GetClusterApproveInfo(ctx, cluster.ClusterId, from)
		if err != nil {
			return nil, err
		}

		
		isOwner := false
		if cluster.ClusterOwner == from {
			isOwner = true
		}
		
		powerDec, err := k.CalculateBurnGetPower(ctx, sdk.OneDec())
		if err != nil {
			return nil, err
		}
		dvm := types.DvmInfo{
			ClusterChatId: cluster.ClusterChatId,
			ClusterId:     clusterId,
			PowerReward:   powerReward.TruncateInt(),
			PowerDvm:      powerDvm.TruncateInt(),
			GasDayDvm:     powerDvm.Quo(k.GetParams(ctx).PowerGasRatio).Quo(powerDec).TruncateInt(),
			AuthContract:  approveInfo.ApproveAddress,
			AuthHeight:    approveInfo.EndBlock,
			ClusterName:   cluster.ClusterName,
			IsOwner:       isOwner,
		}
		if approveInfo.ApproveAddress != "" && approveInfo.EndBlock < ctx.BlockHeight() {
			dvm.AuthContract = ""
			dvm.AuthHeight = 0
		}
		res = append(res, dvm)
	}

	bz, err := codec.MarshalJSONIndent(legacyQuerierCdc, res)

	return bz, nil
}

func queryClusterPersonInfo(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {

	log := core.BuildLog(core.GetPackageFuncName(), core.LmChainDaoQuery)
	var params types.QueryClusterPersonalInfoParams
	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		log.WithError(err).Error("UnmarshalJSON")
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	clusterId := params.ClusterId

	if strings.Contains(clusterId, ".") {
		
		clusterId, err = k.GetClusterId(ctx, clusterId)
		if err != nil {
			return nil, err
		}
	}

	cluster, err := k.GetCluster(ctx, clusterId)
	if err != nil {
		return nil, err
	}

	powerAmount := sdk.ZeroDec()
	burnAmount := sdk.ZeroDec()
	receiveDao := sdk.ZeroDec()
	powerInfo, ok := cluster.ClusterPowerMembers[params.FromAddress]
	powerReward := sdk.ZeroDec()

	endingPeriod := k.IncrementClusterPeriod(ctx, cluster)
	daoParam := k.GetParams(ctx)
	if ok {
		powerAmount = powerInfo.ActivePower
		burnAmount = powerInfo.BurnAmount
		rewards := k.CalculateBurnRewards(ctx, cluster, params.FromAddress, endingPeriod)
		
		powerReward = rewards.AmountOf(core.BaseDenom).Mul(sdk.OneDec().Sub(core.ClusterDeviceRate))
		
		hisReward, err := k.GetClusterMemberReward(ctx, clusterId, params.FromAddress)
		if err != nil {
			return nil, err
		}
		
		if !hisReward.DeviceReward.IsZero() {
			powerReward = powerReward.Add(sdk.NewDecFromInt(hisReward.DeviceReward).Mul(sdk.OneDec().Sub(core.ClusterDeviceRate)))
		}
		if !hisReward.HisReward.IsZero() {
			powerReward = powerReward.Add(sdk.NewDecFromInt(hisReward.HisReward))
		}
		rate := daoParam.ReceiveDaoRatio
		receiveDao = powerInfo.PowerCanReceiveDao.Mul(rate)
	}

	var isDevice bool
	_, ok = cluster.ClusterDeviceMembers[params.FromAddress]
	deviceReward := sdk.ZeroDec()

	if ok {
		isDevice = true
		endingDevicePeriod := k.GetDeviceCurrentRewards(ctx, cluster.ClusterId).Period
		rewards := k.CalculateDeviceRewards(ctx, cluster, params.FromAddress, endingDevicePeriod)
		deviceReward = rewards.AmountOf(core.BaseDenom)
	}

	burnRatio := sdk.ZeroDec()
	if !cluster.ClusterBurnAmount.IsZero() {
		burnRatio = burnAmount.Quo(cluster.ClusterBurnAmount)
	}

	isAdmin := k.IsAdmin(cluster, params.FromAddress)

	var isOwner bool
	if params.FromAddress == cluster.ClusterOwner {
		isOwner = true
	}
	
	powerDec, err := k.CalculateBurnGetPower(ctx, sdk.OneDec())
	if err != nil {
		return nil, err
	}
	gasDay := powerAmount.Quo(k.GetParams(ctx).PowerGasRatio).Quo(powerDec)

	approveInfo, err := k.GetClusterApproveInfo(ctx, cluster.ClusterId, params.FromAddress)

	ownerReward := sdk.ZeroDec()
	if params.FromAddress == cluster.ClusterOwner {
		log.Info(cluster.ClusterDaoPool)
		
		_, ok = cluster.ClusterPowerMembers[cluster.ClusterDaoPool]
		if ok {
			
			rewards := k.CalculateBurnRewards(ctx, cluster, cluster.ClusterDaoPool, endingPeriod)
			finalRewards, _ := rewards.TruncateDecimal()
			hisReward, err := k.GetClusterMemberReward(ctx, cluster.ClusterId, cluster.ClusterDaoPool)
			if err != nil {
				return nil, err
			}
			finalRewards = finalRewards.Add(sdk.NewCoin(core.BaseDenom, hisReward.DeviceReward)).Add(sdk.NewCoin(core.BaseDenom, hisReward.HisReward))
			
			ownerAmount := sdk.NewCoins()
			
			for _, reward := range finalRewards {
				
				am := sdk.NewDecFromInt(reward.Amount).Mul(cluster.ClusterSalaryRatio).TruncateInt()
				ownerAmount = ownerAmount.Add(sdk.NewCoin(reward.GetDenom(), am))
			}
			ownerReward = sdk.NewDecFromInt(ownerAmount.AmountOf(core.BaseDenom))
		}
	}

	res := types.ClusterPersonalInfo{
		PowerAmount:       powerAmount.TruncateInt(),
		GasDay:            gasDay.TruncateInt(),
		BurnAmount:        burnAmount.TruncateInt(),
		IsDevice:          isDevice,
		IsAdmin:           isAdmin,
		IsOwner:           isOwner,
		BurnRatio:         decimal.RequireFromString(burnRatio.String()).RoundFloor(3).String(),
		PowerReward:       powerReward.TruncateInt(),
		DeviceReward:      deviceReward.TruncateInt(),
		OwnerReward:       ownerReward.TruncateInt(),
		AuthContract:      approveInfo.ApproveAddress,
		AuthHeight:        approveInfo.EndBlock,
		ClusterOwner:      cluster.ClusterOwner,
		ClusterName:       cluster.ClusterName,
		ClusterId:         cluster.ClusterId,
		ClusterChatId:     cluster.ClusterChatId,
		DaoReceive:        receiveDao.TruncateInt(),
		BurnRewardFeeRate: core.RemoveStringLastZero(daoParam.BurnRewardFeeRate.String()),
	}
	if approveInfo.ApproveAddress != "" && approveInfo.EndBlock < ctx.BlockHeight() {
		res.AuthContract = ""
		res.AuthHeight = 0
	}
	bz, err := codec.MarshalJSONIndent(legacyQuerierCdc, res)
	return bz, nil
}
func queryCluster(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	clusterId := string(req.Data)
	var err error
	if strings.Contains(clusterId, ".") {
		
		clusterId, err = k.GetClusterId(ctx, clusterId)
		if err != nil {
			return nil, err
		}
	}
	cluster, err := k.GetCluster(ctx, clusterId)
	if err != nil {
		return nil, err
	}
	bz, err := codec.MarshalJSONIndent(legacyQuerierCdc, cluster)

	return bz, nil
}

func queryClusterInfo(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	logs := core.BuildLog(core.GetPackageFuncName(), core.LmChainDaoQuery)
	clusterId := string(req.Data)
	var err error
	if strings.Contains(clusterId, ".") {
		
		clusterId, err = k.GetClusterId(ctx, clusterId)
		if err != nil {
			return nil, err
		}
	}

	cluster, err := k.GetCluster(ctx, clusterId)
	if err != nil {
		if errors.Is(err, core.GetClusterNotFound) {
			return nil, nil
		}
		return nil, err
	}

	
	ownerPowerInfo, ok := cluster.ClusterPowerMembers[cluster.ClusterOwner]
	if !ok {
		return nil, core.ErrClusterOwnerPower
	}

	clusterDaoPoolPower := sdk.ZeroDec()
	clusterDaoPoolReward := sdkmath.ZeroInt()
	
	powerDec, err := k.CalculateBurnGetPower(ctx, sdk.OneDec())
	if err != nil {
		return nil, err
	}
	
	cycleFee, err := k.getCycleFee(ctx, cluster.ClusterChatId)
	if err != nil {
		return nil, err
	}
	clusterDayFreeGas := ownerPowerInfo.ActivePower.Quo(k.GetParams(ctx).PowerGasRatio).Quo(powerDec).TruncateInt()
	
	clusterDayFreeGasBalance := clusterDayFreeGas.Sub(cycleFee.Amount)

	clusterDaoPoolPowerInfo, ok := cluster.ClusterPowerMembers[cluster.ClusterDaoPool]
	if ok {
		clusterDaoPoolPower = clusterDaoPoolPowerInfo.ActivePower
		
		endingPeriod := k.IncrementClusterPeriod(ctx, cluster)
		
		rewards, _ := k.CalculateBurnRewards(ctx, cluster, cluster.ClusterDaoPool, endingPeriod).TruncateDecimal()
		if rewards != nil {
			for _, reward := range rewards {
				if reward.Denom == core.BaseDenom {
					clusterDaoPoolReward = clusterDaoPoolReward.Add(reward.Amount)
				}
			}
		}
	}
	
	daoPoolAddr, err := sdk.AccAddressFromBech32(cluster.ClusterDaoPool)
	if err != nil {
		logs.WithError(err).WithField("DaoPool:", cluster.ClusterDaoPool).Error("queryInClusterInfo AccAddressFromBech32 daoPool Error")
		return nil, core.ErrAddressFormat
	}
	daoPoolBalance := k.BankKeeper.GetAllBalances(ctx, daoPoolAddr)
	
	votePolicyAddr, err := sdk.AccAddressFromBech32(cluster.ClusterVotePolicy)
	if err != nil {
		logs.WithError(err).WithField("votePolicyAddr:", cluster.ClusterVotePolicy).Error("queryInClusterInfo AccAddressFromBech32 ClusterVotePolicy Error")
		return nil, core.ErrAddressFormat
	}
	votePolicyBalance := k.BankKeeper.GetAllBalances(ctx, votePolicyAddr)
	
	routeRewardPoolAddr, err := sdk.AccAddressFromBech32(cluster.ClusterRouteRewardPool)
	if err != nil {
		logs.WithError(err).WithField("routeRewardPoolAddr:", cluster.ClusterRouteRewardPool).Error("queryInClusterInfo AccAddressFromBech32 ClusterRouteRewardPool Error")
		return nil, core.ErrAddressFormat
	}
	routeRewardPoolBalance := k.BankKeeper.GetAllBalances(ctx, routeRewardPoolAddr)

	

	level := cluster.ClusterLevel

	params := k.GetParams(ctx)

	var toLevel int64
	if level == params.ClusterLevels[len(params.ClusterLevels)-1].Level {
		toLevel = level
	} else {
		toLevel = level + 1
	}

	clusterLevelInfo, err := k.GetLevelInfoByLevel(params, toLevel)
	if err != nil {
		return nil, err
	}

	levelInfo := types.ClusterLevelInfo{
		Level:                 level,
		BurnAmountNextLevel:   clusterLevelInfo.BurnAmount,
		ActiveAmountNextLevel: clusterLevelInfo.MemberAmount,
	}

	
	approveInfo, err := k.GetClusterApproveInfo(ctx, cluster.ClusterId, cluster.ClusterDaoPool)
	if err != nil {
		return nil, err
	}
	
	userInfo, err := k.ChatKeeper.GetRegisterInfo(ctx, cluster.ClusterOwner)
	if err != nil {
		return nil, err
	}
	gateway, err := k.gatewayKeeper.GetGatewayInfo(ctx, userInfo.NodeAddress)
	if err != nil {
		return nil, err
	}
	indexNum := ""
	for _, index := range gateway.GatewayNum {
		if index.IsFirst {
			indexNum = index.NumberIndex
			break
		}
	}
	resp := types.ClusterInfo{
		ClusterId:              cluster.ClusterId,
		ClusterChatId:          cluster.ClusterChatId,
		ClusterOwner:           cluster.ClusterOwner,
		ClusterName:            cluster.ClusterName,
		ClusterAllBurn:         cluster.ClusterBurnAmount.TruncateInt(),
		ClusterAllPower:        cluster.ClusterPower.TruncateInt(),
		OnlineRatio:            core.RemoveDecLastZero(cluster.OnlineRatio),
		ClusterActiveDevice:    cluster.ClusterActiveDevice,
		ClusterDeviceAmount:    int64(len(cluster.ClusterDeviceMembers)),
		DeviceConnectivityRate: core.RemoveDecLastZero(k.GetClusterConnectivityRate(ctx, clusterId)),
		ClusterDeviceRatio:     core.RemoveDecLastZero(core.ClusterDeviceRate),
		ClusterSalaryRatio:     core.RemoveDecLastZero(cluster.ClusterSalaryRatio),
		ClusterDaoRatio:        core.RemoveDecLastZero(sdk.OneDec().Sub(cluster.ClusterDaoRatio)),
		ClusterDvmRatio:        core.RemoveDecLastZero(cluster.ClusterDvmRatio),
		ClusterDayFreeGas:      clusterDayFreeGasBalance,
		ClusterDaoPoolPower:    clusterDaoPoolPower.TruncateInt(),
		DaoPoolDayFreeGas:      clusterDaoPoolPower.Quo(k.GetParams(ctx).PowerGasRatio).Quo(powerDec).TruncateInt(),
		DaoPoolAvailableAmount: votePolicyBalance.AmountOf(core.BaseDenom),
		DaoLicensingContract:   approveInfo.ApproveAddress,
		DaoLicensingHeight:     approveInfo.EndBlock,
		LevelInfo:              levelInfo,
		GatewayAddress:         userInfo.NodeAddress,
		IndexNum:               indexNum,
		ClusterVotePolicy:      cluster.ClusterVotePolicy,
		DaoPoolBalance:         daoPoolBalance.AmountOf(core.BurnRewardFeeDenom),
		RoutePoolBalance:       routeRewardPoolBalance.AmountOf(core.BurnRewardFeeDenom),
		ClusterDaoPoolReward:   clusterDaoPoolReward,
		ClusterDaoPool:         cluster.ClusterDaoPool,
	}
	
	if approveInfo.ApproveAddress != "" && approveInfo.EndBlock < ctx.BlockHeight() {
		resp.DaoLicensingContract = ""
		resp.DaoLicensingHeight = 0
	}

	bz, err := codec.MarshalJSONIndent(legacyQuerierCdc, resp)

	return bz, nil
}

func queryInClusters(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {

	addr := string(req.Data)

	pInfo, err := k.GetPersonClusterInfo(ctx, addr)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	var clusters []types.InClusters

	for clusterId := range pInfo.Device {

		inClusters := types.InClusters{
			ClusterId: clusterId,
			IsOwner:   false,
		}

		if _, ok := pInfo.Owner[clusterId]; ok {
			inClusters.IsOwner = true
		}

		clusters = append(clusters, inClusters)
	}

	bz, err := codec.MarshalJSONIndent(legacyQuerierCdc, clusters)

	return bz, nil
}


func queryClustersDeductionFee(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	log := core.BuildLog(core.GetPackageFuncName(), core.LmChainDaoQuery)
	var params types.QueryClusterParams
	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		log.WithError(err).Error("UnmarshalJSON")
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}
	clusterId := params.ClusterId
	if strings.Contains(clusterId, ".") {
		
		clusterId, err = k.GetClusterId(ctx, clusterId)
		if err != nil {
			return nil, err
		}
	}
	cluster, err := k.GetCluster(ctx, clusterId)
	if err != nil {
		return nil, err
	}
	cycleFee, err := k.getCycleFee(ctx, cluster.ClusterChatId)
	if err != nil {
		return nil, err
	}
	
	powerDec, err := k.CalculateBurnGetPower(ctx, sdk.OneDec())
	if err != nil {
		return nil, err
	}
	
	limit := cluster.ClusterPowerMembers[cluster.ClusterOwner].ActivePower.Quo(k.GetParams(ctx).PowerGasRatio).Quo(powerDec).TruncateInt()
	
	balance := limit.Sub(cycleFee.Amount)
	bz, err := codec.MarshalJSONIndent(legacyQuerierCdc, balance)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return bz, nil

}


func queryClustersOwnerReward(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	log := core.BuildLog(core.GetPackageFuncName(), core.LmChainDaoQuery)
	var params types.QueryClusterRewardParams
	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		log.WithError(err).Error("UnmarshalJSON")
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}
	clusterId := params.ClusterId
	if strings.Contains(clusterId, ".") {
		
		clusterId, err = k.GetClusterId(ctx, clusterId)
		if err != nil {
			return nil, err
		}
	}
	cluster, err := k.GetCluster(ctx, clusterId)
	if err != nil {
		return nil, err
	}

	
	ownerAmount := sdk.NewCoins()
	_, ok := cluster.ClusterPowerMembers[cluster.ClusterDaoPool]
	if ok {
		endingPeriod := k.IncrementClusterPeriod(ctx, cluster)
		
		rewards, _ := k.CalculateBurnRewards(ctx, cluster, cluster.ClusterDaoPool, endingPeriod).TruncateDecimal()
		if rewards == nil {
			rewards = sdk.NewCoins(sdk.NewCoin(core.BaseDenom, sdk.ZeroInt()))
		}
		
		for _, reward := range rewards {
			
			am := sdk.NewDecFromInt(reward.Amount).Mul(cluster.ClusterSalaryRatio).TruncateInt()
			ownerAmount = append(ownerAmount, sdk.NewCoin(reward.GetDenom(), am))
		}
	} else {
		ownerAmount = sdk.NewCoins(sdk.NewCoin(core.BaseDenom, sdk.ZeroInt()))
	}

	bz, err := codec.MarshalJSONIndent(legacyQuerierCdc, ownerAmount)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return bz, nil
}

func queryClustersGasReward(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	log := core.BuildLog(core.GetPackageFuncName(), core.LmChainDaoQuery)
	var params types.QueryClusterRewardParams
	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		log.WithError(err).Error("UnmarshalJSON")
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}
	clusterId := params.ClusterId
	
	if clusterId == "" {
		return queryAllClustersGasReward(ctx, req, k, legacyQuerierCdc)
	}
	if strings.Contains(clusterId, ".") {
		
		clusterId, err = k.GetClusterId(ctx, clusterId)
		if err != nil {
			return nil, err
		}
	}
	cluster, err := k.GetCluster(ctx, clusterId)
	if err != nil {
		return nil, err
	}
	powerReward := sdk.ZeroDec()
	rewardCoins := types.LpTokens{}
	erc20Coins := types.LpTokens{}
	_, ok := cluster.ClusterPowerMembers[params.Member]
	if ok {
		endingPeriod := k.IncrementClusterPeriod(ctx, cluster)
		rewards := k.CalculateBurnRewards(ctx, cluster, params.Member, endingPeriod)
		
		powerReward = rewards.AmountOf(core.BaseDenom).Mul(sdk.OneDec().Sub(core.ClusterDeviceRate))
		hisReward, err := k.GetClusterMemberReward(ctx, clusterId, params.Member)
		if err != nil {
			return nil, err
		}
		
		if !hisReward.DeviceReward.IsZero() {
			powerReward = powerReward.Add(sdk.NewDecFromInt(hisReward.DeviceReward).Mul(sdk.OneDec().Sub(core.ClusterDeviceRate)))
		}
		
		if !hisReward.HisReward.IsZero() {
			powerReward = powerReward.Add(sdk.NewDecFromInt(hisReward.HisReward))
		}
		dstToken := types.LpToken{LpToken: sdk.NewCoin(core.BaseDenom, powerReward.TruncateInt())}
		rewardCoins = append(rewardCoins, dstToken)
		for _, reward := range rewards {
			if reward.Denom == core.BaseDenom {
				continue
			}
			erc20Swap, err := k.GetErc20Symbol(ctx, reward.Denom)
			if err != nil || erc20Swap == nil || erc20Swap.Symbol == "" {
				continue
			}
			erc20Token := types.LpToken{
				LpToken:    sdk.NewCoin(erc20Swap.Symbol, reward.Amount.TruncateInt()),
				LpContract: erc20Swap.Contract,
				Token0:     erc20Swap.Token0,
				Token1:     erc20Swap.Token1,
			}
			erc20Coins = erc20Coins.Add(erc20Token)
		}
	}
	rewardCoins = append(rewardCoins, erc20Coins...)
	daoParams := k.GetParams(ctx)
	resp := types.QueryClusterRewardResp{Reward: rewardCoins, Rate: daoParams.BurnRewardFeeRate}
	bz, err := codec.MarshalJSONIndent(legacyQuerierCdc, resp)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return bz, nil
}

func queryAllClustersGasReward(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	log := core.BuildLog(core.GetPackageFuncName(), core.LmChainDaoQuery)
	var params types.QueryClusterRewardParams
	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		log.WithError(err).Error("UnmarshalJSON")
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	personalInfo, err := k.GetPersonClusterInfo(ctx, params.Member)
	if err != nil {
		return nil, err
	}

	powerReward := sdk.ZeroDec()
	rewardCoins := types.LpTokens{}
	dstCoins := types.LpTokens{}
	erc20Coins := types.LpTokens{}

	for clusterId, _ := range personalInfo.BePower {
		cluster, err := k.GetCluster(ctx, clusterId)
		if err != nil {
			return nil, err
		}
		_, ok := cluster.ClusterPowerMembers[params.Member]
		if ok {
			endingPeriod := k.IncrementClusterPeriod(ctx, cluster)
			rewards := k.CalculateBurnRewards(ctx, cluster, params.Member, endingPeriod)
			
			powerReward = rewards.AmountOf(core.BaseDenom).Mul(sdk.OneDec().Sub(core.ClusterDeviceRate))
			hisReward, err := k.GetClusterMemberReward(ctx, clusterId, params.Member)
			if err != nil {
				return nil, err
			}
			
			if !hisReward.DeviceReward.IsZero() {
				powerReward = powerReward.Add(sdk.NewDecFromInt(hisReward.DeviceReward).Mul(sdk.OneDec().Sub(core.ClusterDeviceRate)))
			}
			
			if !hisReward.HisReward.IsZero() {
				powerReward = powerReward.Add(sdk.NewDecFromInt(hisReward.HisReward))
			}
			dstToken := types.LpToken{LpToken: sdk.NewCoin(core.BaseDenom, powerReward.TruncateInt())}
			dstCoins = dstCoins.Add(dstToken)
			for _, reward := range rewards {
				if reward.Denom == core.BaseDenom {
					continue
				}
				erc20Swap, err := k.GetErc20Symbol(ctx, reward.Denom)
				if err != nil || erc20Swap == nil || erc20Swap.Symbol == "" {
					continue
				}
				erc20Token := types.LpToken{
					LpToken:    sdk.NewCoin(erc20Swap.Symbol, reward.Amount.TruncateInt()),
					LpContract: erc20Swap.Contract,
					Token0:     erc20Swap.Token0,
					Token1:     erc20Swap.Token1,
				}
				erc20Coins = erc20Coins.Add(erc20Token)
			}
		}
	}
	rewardCoins = append(rewardCoins, dstCoins...)
	rewardCoins = append(rewardCoins, erc20Coins...)
	daoParams := k.GetParams(ctx)
	resp := types.QueryClusterRewardResp{Reward: rewardCoins, Rate: daoParams.BurnRewardFeeRate}
	bz, err := codec.MarshalJSONIndent(legacyQuerierCdc, resp)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return bz, nil
}

func queryAllClustersReward(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	log := core.BuildLog(core.GetPackageFuncName(), core.LmChainDaoQuery)
	var params types.QueryClusterRewardParams
	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		log.WithError(err).Error("UnmarshalJSON")
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	personalInfo, err := k.GetPersonClusterInfo(ctx, params.Member)
	if err != nil {
		return nil, err
	}

	powerReward := sdk.ZeroDec()
	rewardMap := make(map[string]sdk.Coin)
	for clusterId, _ := range personalInfo.BePower {
		cluster, err := k.GetCluster(ctx, clusterId)
		if err != nil {
			return nil, err
		}
		_, ok := cluster.ClusterPowerMembers[params.Member]
		if ok {
			endingPeriod := k.IncrementClusterPeriod(ctx, cluster)
			rewards := k.CalculateBurnRewards(ctx, cluster, params.Member, endingPeriod)
			
			powerReward = rewards.AmountOf(core.BaseDenom).Mul(sdk.OneDec().Sub(core.ClusterDeviceRate))
			hisReward, err := k.GetClusterMemberReward(ctx, clusterId, params.Member)
			if err != nil {
				return nil, err
			}
			
			if !hisReward.DeviceReward.IsZero() {
				powerReward = powerReward.Add(sdk.NewDecFromInt(hisReward.DeviceReward).Mul(sdk.OneDec().Sub(core.ClusterDeviceRate)))
			}
			
			if !hisReward.HisReward.IsZero() {
				powerReward = powerReward.Add(sdk.NewDecFromInt(hisReward.HisReward))
			}
			dstToken := sdk.NewCoin(core.BaseDenom, powerReward.TruncateInt())
			rewardMap[cluster.ClusterChatId] = dstToken
		}
	}
	daoParams := k.GetParams(ctx)
	resp := types.QueryAllClusterRewardResp{Reward: rewardMap, Rate: daoParams.BurnRewardFeeRate}
	bz, err := codec.MarshalJSONIndent(legacyQuerierCdc, resp)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return bz, nil
}


func queryClustersDeviceReward(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	log := core.BuildLog(core.GetPackageFuncName(), core.LmChainDaoQuery)
	var params types.QueryClusterRewardParams
	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		log.WithError(err).Error("UnmarshalJSON")
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}
	clusterId := params.ClusterId
	if strings.Contains(clusterId, ".") {
		
		clusterId, err = k.GetClusterId(ctx, clusterId)
		if err != nil {
			return nil, err
		}
	}
	cluster, err := k.GetCluster(ctx, clusterId)
	if err != nil {
		return nil, err
	}
	rewards := sdk.NewCoins(sdk.NewCoin(core.BaseDenom, sdk.ZeroInt()))
	_, ok := cluster.ClusterDeviceMembers[params.Member]
	if ok {
		endingPeriod := k.GetDeviceCurrentRewards(ctx, cluster.ClusterId).Period
		rewards, _ = k.CalculateDeviceRewards(ctx, cluster, params.Member, endingPeriod).TruncateDecimal()
	}
	bz, err := codec.MarshalJSONIndent(legacyQuerierCdc, rewards)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return bz, nil
}

func queryGatewayClusters(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	log := core.BuildLog(core.GetPackageFuncName(), core.LmChainDaoQuery)
	var params types.QueryGatewayClustersParams
	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		log.WithError(err).Error("UnmarshalJSON")
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	clusters, err := k.GateClusterToGateway(ctx, params.GatewayAddress)
	if err != nil {
		return nil, err
	}

	res, err := util.Json.Marshal(clusters)
	if err != nil {
		log.WithError(err).Error("json marshal")
		return nil, err
	}

	return res, nil
}

func queryClusterInfoById(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	log := core.BuildLog(core.GetPackageFuncName(), core.LmChainDaoQuery)
	id := string(req.Data)

	p, err := k.GetClusterByChatId(ctx, id)
	if err != nil {
		log.WithError(err).Info("queryClusterInfoById GetClusterByChatId Err")
		return nil, err
	}

	pByte, err := codec.MarshalJSONIndent(legacyQuerierCdc, p)
	if err != nil {
		log.WithError(err).Info("queryClusterInfoById MarshalJSONIndent Err")
		return nil, err
	}

	return pByte, nil
}

func queryPersonClusterInfo(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	var err error
	log := core.BuildLog(core.GetPackageFuncName(), core.LmChainDaoQuery)
	var params types.QueryPersonClusterInfoRequest

	err = legacyQuerierCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		log.WithError(err).Error("UnmarshalJSON")
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	p, err := k.GetPersonClusterInfo(ctx, params.From)
	if err != nil {
		log.WithError(err).Info("queryPersonClusterInfo GetPersonClusterInfo Err")
		return nil, err
	}

	
	clusters := make(map[string]types.DeviceCluster)

	
	ownerSlice := make([]string, 0)
	powerSlice := make([]string, 0)
	
	deviceInfos := make([]types.DeviceInfo, 0)

	var cluster types.DeviceCluster
	for device := range p.Device {
		if clusterC, ok := clusters[device]; ok {
			cluster = clusterC
		} else {
			cluster, err = k.GetCluster(ctx, device)
			if err != nil {
				return nil, err
			}

			clusters[device] = cluster
		}

		
		createTime, err := k.GetClusterCreateTime(ctx, device)
		if err != nil {
			return nil, err
		}

		deviceInfo := types.DeviceInfo{
			ClusterChatId:     cluster.ClusterChatId,
			ClusterName:       cluster.ClusterName,
			ClusterLevel:      cluster.ClusterLevel,
			ClusterOwner:      cluster.ClusterOwner,
			ClusterCreateTime: createTime,
		}

		deviceInfos = append(deviceInfos, deviceInfo)
	}

	for owner := range p.Owner {
		if clusterC, ok := clusters[owner]; ok {
			cluster = clusterC
		} else {
			cluster, err = k.GetCluster(ctx, owner)
			if err != nil {
				return nil, err
			}

			clusters[owner] = cluster
		}

		ownerSlice = append(ownerSlice, cluster.ClusterChatId)
	}

	for power := range p.BePower {
		if clusterC, ok := clusters[power]; ok {
			cluster = clusterC
		} else {
			cluster, err = k.GetCluster(ctx, power)
			if err != nil {
				return nil, err
			}

			clusters[power] = cluster
		}
		powerSlice = append(powerSlice, cluster.ClusterChatId)
	}

	
	firstBurnCluster := ""

	if p.FirstPowerCluster != "" { 
		firstBurnCluster = p.FirstPowerCluster
	} else if len(p.BePower) > 0 { 
		createTimeArr := make([]struct {
			CreateTime int64
			ClusterId  string
		}, 0)
		for clusterId := range p.BePower {
			createTime, err := k.GetClusterCreateTime(ctx, clusterId)
			if err != nil {
				return nil, err
			}

			c, err := k.GetCluster(ctx, clusterId)
			if err != nil {
				return nil, err
			}

			createTimeArr = append(createTimeArr, struct {
				CreateTime int64
				ClusterId  string
			}{
				CreateTime: createTime,
				ClusterId:  c.ClusterChatId,
			})
		}

		minCreateTime := createTimeArr[0].CreateTime
		for _, v := range createTimeArr {
			if v.CreateTime < minCreateTime {
				minCreateTime = v.CreateTime
				firstBurnCluster = v.ClusterId
			}
		}
	}

	pResp := types.PersonClusterStatisticsInfo{
		Address:           params.From,
		Owner:             ownerSlice,
		BePower:           powerSlice,
		AllBurn:           p.AllBurn.TruncateInt(),
		ActivePower:       p.ActivePower.TruncateInt(),
		FreezePower:       p.FreezePower.TruncateInt(),
		DeviceInfo:        deviceInfos,
		FirstPowerCluster: firstBurnCluster,
	}

	pByte, err := codec.MarshalJSONIndent(legacyQuerierCdc, pResp)
	if err != nil {
		log.WithError(err).Info("queryPersonClusterInfo MarshalJSONIndent Err")
		return nil, err
	}

	return pByte, nil
}

func QueryDaoParams(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	logs := core.BuildLog(core.GetPackageFuncName(), core.LmChainDaoQuery)
	params := k.GetParams(ctx)

	
	burnGetPowerRatio, err := k.CalculateBurnGetPower(ctx, sdk.MustNewDecFromStr("1000000000000000000"))
	if err != nil {
		return nil, err
	}

	logs.Info("burnGetPowerRatio", burnGetPowerRatio)

	
	allPower, err := k.GetTotalPowerAmount(ctx)
	if err != nil {
		return nil, err
	}

	allBurn, err := k.GetTotalBurnAmount(ctx)
	if err != nil {
		return nil, err
	}

	logs.Info("allPower", allPower)

	var dayReward sdk.Dec
	if allPower.IsZero() {
		dayReward = params.DayMintAmount
	} else {
		dayReward = burnGetPowerRatio.Quo(allPower).Mul(params.DayMintAmount)
	}
	logs.Info("dayReward", dayReward)

	
	record, err := k.GetDstPerPowerDay(ctx)
	if err != nil {
		return nil, err
	}

	
	oneDstPower, err := k.CalculateBurnGetPower(ctx, sdk.MustNewDecFromStr("1000000000000000000"))

	var sevenDayYearRate string

	
	if len(record) >= 7 && k.StartMint(ctx) {
		sumDay7Reward := sdk.ZeroDec()
		for i := 0; i < 7; i++ {
			sumDay7Reward = sumDay7Reward.Add(record[i].Mul(oneDstPower))
		}

		sevenDayYearRateDec := sumDay7Reward.Quo(sdk.NewDec(7)).Mul(sdk.NewDec(365))
		b := sevenDayYearRateDec.Quo(sdk.MustNewDecFromStr("10000000000000000")).TruncateInt()
		c := sdk.NewDecFromInt(b).QuoInt64(100).String()
		sevenDayYearRate = core.RemoveStringLastZero(c)
	}

	
	var dayRewardStr string
	if k.StartMint(ctx) {
		dayRewardStr = core.RemoveStringLastZero(dayReward.String())
	}

	res := types.DaoParams{
		BurnGetPowerRatio: burnGetPowerRatio.TruncateInt().String(),
		SalaryRange: types.Range{
			Max: core.RemoveStringLastZero(params.SalaryRewardRatio.MaxRatio.String()),
			Min: core.RemoveStringLastZero(params.SalaryRewardRatio.MinRatio.String()),
		},
		DeviceRange: types.Range{
			Max: core.RemoveStringLastZero(core.ClusterDeviceRate.String()),
			Min: core.RemoveStringLastZero(core.ClusterDeviceRate.String()),
		},
		DvmRange: types.Range{
			Max: core.RemoveStringLastZero(params.DvmRewardRatio.MaxRatio.String()),
			Min: core.RemoveStringLastZero(params.DvmRewardRatio.MinRatio.String()),
		},
		CreateClusterMinBurn: params.ClusterLevels[1].BurnAmount,
		BurnAddress:          k.accountKeeper.GetModuleAddress(types.ModuleName).String(),
		DayBurnReward:        dayRewardStr,
		DaoRange: types.Range{
			Max: core.RemoveStringLastZero(params.DaoRewardRatio.MaxRatio.String()),
			Min: core.RemoveStringLastZero(params.DaoRewardRatio.MinRatio.String()),
		},
		ReceiveDaoRatio:   core.RemoveStringLastZero(params.ReceiveDaoRatio.String()),
		BurnRewardFeeRate: core.RemoveStringLastZero(params.BurnRewardFeeRate.String()),
		SevenDayYearRate:  sevenDayYearRate,
		TotalPower:        allPower.String(),
		TotalBurn:         allBurn.String(),
	}

	resp, err := codec.MarshalJSONIndent(legacyQuerierCdc, res)
	if err != nil {
		logs.WithError(err).Info("queryPersonClusterInfo MarshalJSONIndent Err")
		return nil, err
	}

	return resp, nil
}

func queryClusterProposals(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	log := core.BuildLog(core.GetPackageFuncName(), core.LmChainDaoQuery)
	var params types.QueryClusterProposalsParams
	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		log.WithError(err).Error("UnmarshalJSON")
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}
	clusterId := params.ClusterId
	if strings.Contains(clusterId, ".") {
		
		clusterId, err = k.GetClusterId(ctx, clusterId)
		if err != nil {
			return nil, err
		}
	}
	cluster, err := k.GetCluster(ctx, clusterId)
	if err != nil {
		return nil, err
	}
	policyAddress := cluster.ClusterVotePolicy
	proposalsRes, err := k.groupKeeper.ProposalsByGroupPolicy(ctx, &group.QueryProposalsByGroupPolicyRequest{
		Address: policyAddress,
	})
	if err != nil {
		return nil, err
	}
	var votersRes *group.QueryVotesByVoterResponse
	if params.Voter != "" {
		votersRes, err = k.groupKeeper.VotesByVoter(ctx, &group.QueryVotesByVoterRequest{
			Voter: params.Voter,
		})
	}
	if err != nil {
		return nil, err
	}
	var proposals []types.ClusterProposals

	for _, p := range proposalsRes.Proposals {
		messages, _ := tx.GetMsgs(p.Messages, "")

		clusterProposals := types.ClusterProposals{}
		clusterProposals.Id = p.Id
		clusterProposals.GroupPolicyAddress = p.GroupPolicyAddress
		clusterProposals.Metadata = p.Metadata
		clusterProposals.Proposers = p.Proposers
		clusterProposals.SubmitTime = p.SubmitTime
		clusterProposals.GroupVersion = p.GroupVersion
		clusterProposals.GroupPolicyVersion = p.GroupPolicyVersion
		clusterProposals.Status = p.Status
		clusterProposals.FinalTallyResult = p.FinalTallyResult
		clusterProposals.VotingPeriodEnd = p.VotingPeriodEnd
		clusterProposals.ExecutorResult = p.ExecutorResult
		clusterProposals.Messages = messages
		if votersRes != nil {
			for _, v := range votersRes.Votes {
				if v.ProposalId == p.Id {
					clusterProposals.IsVoter = true
				}
			}
		}
		proposals = append(proposals, clusterProposals)
	}
	if err != nil {
		return nil, err
	}
	
	bz, err := json.Marshal(proposals)
	
	if err != nil {
		return nil, err
	}
	return bz, nil
}

func queryClusterPersonalInfo(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	var err error
	personalId, err := strconv.ParseUint(string(req.Data), 10, 64)
	if err != nil {
		return nil, err
	}
	proposalRes, err := k.groupKeeper.Proposal(ctx, &group.QueryProposalRequest{
		ProposalId: personalId,
	})
	messages, _ := tx.GetMsgs(proposalRes.Proposal.Messages, "")
	proposal := group.ProposalNew{
		Id:                 proposalRes.Proposal.Id,
		GroupPolicyAddress: proposalRes.Proposal.GroupPolicyAddress,
		Metadata:           proposalRes.Proposal.Metadata,
		Proposers:          proposalRes.Proposal.Proposers,
		SubmitTime:         proposalRes.Proposal.SubmitTime,
		GroupVersion:       proposalRes.Proposal.GroupVersion,
		GroupPolicyVersion: proposalRes.Proposal.GroupPolicyVersion,
		Status:             proposalRes.Proposal.Status,
		FinalTallyResult:   proposalRes.Proposal.FinalTallyResult,
		VotingPeriodEnd:    proposalRes.Proposal.VotingPeriodEnd,
		ExecutorResult:     proposalRes.Proposal.ExecutorResult,
		Messages:           messages,
	}
	if err != nil {
		return nil, err
	}
	
	
	bz, err := json.Marshal(proposal)
	if err != nil {
		return nil, err
	}
	return bz, nil
}

func queryClusterProposalVoters(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	var params types.QueryProposalVotersParams

	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}
	personalId, err := strconv.ParseUint(params.ProposalId, 10, 64)
	if err != nil {
		return nil, err
	}
	clusterId := params.ClusterId
	if strings.Contains(clusterId, ".") {
		
		clusterId, err = k.GetClusterId(ctx, clusterId)
		if err != nil {
			return nil, err
		}
	}
	cluster, err := k.GetCluster(ctx, clusterId)
	if err != nil {
		return nil, err
	}
	votersRes, err := k.groupKeeper.VotesByProposal(ctx, &group.QueryVotesByProposalRequest{
		ProposalId: personalId,
	})
	if err != nil {
		return nil, err
	}
	groupInfo, err := k.groupKeeper.GetGroupInfo(ctx, cluster.ClusterVoteId)
	if err != nil {
		return nil, err
	}
	groupMembersRes, err := k.groupKeeper.GroupMembers(ctx, &group.QueryGroupMembersRequest{
		GroupId: cluster.ClusterVoteId,
	})
	totalWeight, err := sdk.NewDecFromStr(groupInfo.TotalWeight)
	if err != nil {
		return nil, err
	}
	resp := types.QueryProposalVotersResp{}
	for _, v := range votersRes.Votes {
		vote := types.Vote{ProposalId: v.ProposalId, Voter: v.Voter, Option: v.Option, SubmitTime: v.SubmitTime, Weight: ""}
		for _, m := range groupMembersRes.Members {
			if m.Member.Address == v.Voter {
				memberWeight, err := sdk.NewDecFromStr(m.Member.Weight)
				if err != nil {
					continue
				}
				vote.Weight = memberWeight.Quo(totalWeight).String()
			}
		}
		resp.Votes = append(resp.Votes, vote)
	}
	
	bz, err := json.Marshal(resp)
	
	if err != nil {
		return nil, err
	}
	return bz, nil
}

func queryClusterProposalVoter(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	log := core.BuildLog(core.GetPackageFuncName(), core.LmChainDaoQuery)
	var params types.QueryClusterProposalVoterParams
	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		log.WithError(err).Error("UnmarshalJSON")
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}
	personalId, err := strconv.ParseUint(params.ProposalId, 10, 64)
	if err != nil {
		return nil, err
	}
	voterRes, err := k.groupKeeper.VoteByProposalVoter(ctx, &group.QueryVoteByProposalVoterRequest{
		ProposalId: personalId,
		Voter:      params.Voter,
	})

	if err != nil {
		return nil, err
	}
	
	
	bz, err := json.Marshal(voterRes)
	if err != nil {
		return nil, err
	}
	return bz, nil
}

func queryGroupMembers(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	clusterId := string(req.Data)
	var err error
	if strings.Contains(clusterId, ".") {
		
		clusterId, err = k.GetClusterId(ctx, clusterId)
		if err != nil {
			return nil, err
		}
	}
	cluster, err := k.GetCluster(ctx, clusterId)
	if err != nil {
		return nil, err
	}
	membersRes, err := k.groupKeeper.GroupMembers(ctx, &group.QueryGroupMembersRequest{
		GroupId: cluster.ClusterVoteId,
	})

	if err != nil {
		return nil, err
	}
	
	bz, err := json.Marshal(membersRes)
	
	if err != nil {
		return nil, err
	}
	return bz, nil
}

func queryGroupInfo(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	clusterId := string(req.Data)
	var err error
	if strings.Contains(clusterId, ".") {
		
		clusterId, err = k.GetClusterId(ctx, clusterId)
		if err != nil {
			return nil, err
		}
	}
	cluster, err := k.GetCluster(ctx, clusterId)
	if err != nil {
		return nil, err
	}
	groupRes, err := k.groupKeeper.GroupInfo(ctx, &group.QueryGroupInfoRequest{
		GroupId: cluster.ClusterVoteId,
	})

	if err != nil {
		return nil, err
	}
	
	
	bz, err := json.Marshal(groupRes)
	if err != nil {
		return nil, err
	}
	return bz, nil
}
func queryClusterApproveInfo(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	log := core.BuildLog(core.GetPackageFuncName(), core.LmChainDaoQuery)
	var params types.QueryClusterApproveParams
	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		log.WithError(err).Error("UnmarshalJSON")
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}
	clusterId := params.ClusterId
	if strings.Contains(clusterId, ".") {
		
		clusterId, err = k.GetClusterId(ctx, clusterId)
		if err != nil {
			return nil, err
		}
	}
	if params.Address == "" {
		cluster, err := k.GetCluster(ctx, clusterId)
		if err != nil {
			return nil, err
		}
		params.Address = cluster.ClusterDaoPool
	}
	clusterApproveInfo, err := k.GetClusterApproveInfo(ctx, clusterId, params.Address)
	if err != nil {
		return nil, err
	}
	
	bz, err := json.Marshal(clusterApproveInfo)

	return bz, nil
}
func queryClusterProposalTallyResult(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	var err error
	personalId, err := strconv.ParseUint(string(req.Data), 10, 64)
	if err != nil {
		return nil, err
	}
	tallyResult, err := k.groupKeeper.TallyResult(ctx, &group.QueryTallyResultRequest{
		ProposalId: personalId,
	})
	bz, err := json.Marshal(tallyResult)
	if err != nil {
		return nil, err
	}
	return bz, nil
}

func queryGroupVotesByVoter(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	var err error
	voter := string(req.Data)
	votersRes, err := k.groupKeeper.VotesByVoter(ctx, &group.QueryVotesByVoterRequest{
		Voter: voter,
	})

	if err != nil {
		return nil, err
	}
	
	bz, err := json.Marshal(votersRes)
	
	if err != nil {
		return nil, err
	}
	return bz, nil
}
