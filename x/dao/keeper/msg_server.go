package keeper

import (
	"context"
	sdkmath "cosmossdk.io/math"
	"encoding/hex"
	"freemasonry.cc/blockchain/util"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/tendermint/tendermint/crypto/tmhash"
	"math/big"
	"strconv"
	"strings"

	"freemasonry.cc/blockchain/core"
	"freemasonry.cc/blockchain/x/dao/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/group"
	"github.com/sirupsen/logrus"
)

var _ types.MsgServer = &msgServer{}

type msgServer struct {
	Keeper
	logPrefix string
}

func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper, logPrefix: "dao | msgServer | "}
}

func (k msgServer) ReceivePowerCutReward(goCtx context.Context, msg *types.MsgReceivePowerCutReward) (*types.MsgEmptyResponse, error) {
	

	ctx := sdk.UnwrapSDKContext(goCtx)

	
	personCycleInfo, err := k.GetPowerRewardCycleInfo(ctx, msg.FromAddress)
	if err != nil {
		return nil, err
	}

	allReceiveAmount := sdk.ZeroInt()

	accFrom, err := sdk.AccAddressFromBech32(msg.FromAddress)
	if err != nil {
		return nil, core.ErrAddressFormat
	}

	nowTime := ctx.BlockTime().Unix()

	
	for i, cycleInfo := range personCycleInfo.CycleInfo {
		if cycleInfo.Status == 1 {

			status := int64(1)

			
			diff := nowTime - cycleInfo.StartTime
			shouldReceiveTimes := diff / 86400 
			if shouldReceiveTimes >= core.CutPowerRewardTimes {
				shouldReceiveTimes = core.CutPowerRewardTimes
				status = 2
			}

			receivedTimes := cycleInfo.ReceiveTimes 

			receiveTimes := shouldReceiveTimes - receivedTimes

			receiveRewardAmount := cycleInfo.CutPerReward.Mul(sdk.NewInt(receiveTimes))

			
			allReceiveAmount = allReceiveAmount.Add(receiveRewardAmount)

			
			personCycleInfo.CycleInfo[i] = types.CycleInfo{
				Cycle:             cycleInfo.Cycle,
				AllReward:         cycleInfo.AllReward,
				CutPerReward:      cycleInfo.CutPerReward,
				RemainReward:      cycleInfo.RemainReward.Sub(receiveRewardAmount),
				ReceiveTimes:      shouldReceiveTimes,
				AllCutReward:      cycleInfo.AllCutReward,
				StartTime:         cycleInfo.StartTime,
				ClusterRewardList: cycleInfo.ClusterRewardList,
				Status:            status,
			}

		}
	}

	if allReceiveAmount.IsZero() {
		return nil, core.ErrNoPowerCutReward
	}

	err = k.BankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, accFrom, sdk.NewCoins(sdk.NewCoin(core.BaseDenom, allReceiveAmount)))
	if err != nil {
		return nil, err
	}

	
	err = k.SetPowerRewardCycleInfo(ctx, msg.FromAddress, personCycleInfo)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeReceiveCutReward,
		sdk.NewAttribute(types.AttributeSendeer, msg.FromAddress),
		sdk.NewAttribute(types.AttributeKeyAmount, allReceiveAmount.String()),
	))

	return nil, nil

}

func (k msgServer) StartPowerRewardRedeem(goCtx context.Context, msg *types.MsgStartPowerRewardRedeem) (*types.MsgEmptyResponse, error) {
	logs := core.BuildLog(core.GetStructFuncName(k), core.LmChainDaoKeeper)
	ctx := sdk.UnwrapSDKContext(goCtx)

	
	nowTime := ctx.BlockTime().Unix()
	startTime := k.GetStartTime(ctx)
	diff := nowTime - startTime

	if diff/core.CutProductionSeconds == 0 {
		return nil, core.ErrCutProductionNotStart
	}

	
	curentCycle, err := k.GetCutPowerRewardCycle(ctx)
	if err != nil {
		logs.WithError(err).Error("Get cut power reward cycle error：", err)
		return nil, err
	}

	
	personCycleInfo, err := k.GetPowerRewardCycleInfo(ctx, msg.FromAddress)
	if err != nil {
		logs.WithError(err).Error("Get power reward cycle info error：", err)
		return nil, err
	}

	_, ok := personCycleInfo.CycleInfo[curentCycle]
	if !ok { 
		err = k.StartPowerRewardCut(ctx, personCycleInfo, msg.FromAddress, curentCycle)
		if err != nil {
			logs.WithError(err).Error("Start power reward cut error：", err)
			return nil, err
		}
	} else { 
		return nil, core.ErrCurrentCycleAlreadyStart
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeStartPowerRewardRedeem,
		sdk.NewAttribute(types.AttributeSendeer, msg.FromAddress),
	))

	return &types.MsgEmptyResponse{}, nil
}

func (k msgServer) AgreeJoinClusterApply(goCtx context.Context, msg *types.MsgAgreeJoinClusterApply) (*types.MsgEmptyResponse, error) {
	logs := core.BuildLog(core.GetStructFuncName(k), core.LmChainDaoKeeper)
	logs.Debug("AgreeJoinClusterApply Start------------")
	ctx := sdk.UnwrapSDKContext(goCtx)

	fromAddr, err := sdk.AccAddressFromBech32(msg.FromAddress)
	if err != nil {
		logs.WithError(err).Error("AccAddressFromBech32 error")
		return nil, core.ErrAddressFormat
	}
	logs.Debug("ap1")
	

	ok, err := k.ValidateGatewaySign(ctx, fromAddr, msg.MemberOnlineAmount, []string{msg.MemberAddress}, msg.GatewayAddress, msg.GatewaySign, true)
	if err != nil {
		logs.WithError(err).Error("AgreeJoinClusterApply ValidateGatewaySign Error")
		return nil, core.ErrGatewaySign
	}
	logs.Debug("ap2")
	if !ok {
		logs.Error("AgreeJoinClusterApply ValidateGatewaySign Error")
		return nil, core.ErrGatewaySign
	}
	logs.Debug("ap3")

	
	signBytes, err := hex.DecodeString(msg.Sign)
	if err != nil {
		logs.WithError(err).WithField("sign", msg.Sign).Error("Decode sign error")
		return nil, core.SignError
	}

	logs.Debug("ap4")
	
	dataUnSign := []byte(msg.MemberAddress + msg.ClusterId)

	
	comPub, err := util.GetPubKeyFromSign(signBytes, dataUnSign)
	if err != nil {
		logs.WithError(err).Error("Get pubkey from sign error")
		return nil, core.SignError
	}

	logs.Debug("ap5")
	
	ethMemberAddr := common.HexToAddress(comPub.Address().String())
	memberAddr := sdk.AccAddress(ethMemberAddr.Bytes())
	if memberAddr.String() != msg.MemberAddress {
		logs.WithFields(logrus.Fields{
			"memberAddr":        memberAddr.String(),
			"msg.MemberAddress": msg.MemberAddress,
		}).Error("Member signature address error")
		return &types.MsgEmptyResponse{}, core.SignError
	}
	logs.Debug("ap6")
	
	cluster, err := k.GetClusterByChatId(ctx, msg.ClusterId)
	if err != nil {
		logs.WithError(err).Error("Get cluster by chat cluster id error")
		return nil, err
	}
	logs.Debug("ap7")
	
	if _, ok := cluster.ClusterDeviceMembers[msg.MemberAddress]; ok {
		logs.WithField("msg.MemberAddress", msg.MemberAddress).Error("Member already in cluster")
		return nil, core.ErrErrMemberAlreadyInCluster
	}
	logs.Debug("ap8")
	
	members := make([]types.Members, 0)
	members = append(members, types.Members{
		MemberAddress: msg.MemberAddress,
		IndexNum:      msg.IndexNum,
		ChatAddress:   msg.ChatAddress,
	})
	logs.Debug("ap9")
	err = k.AddMember(ctx, cluster.ClusterId, msg.FromAddress, members, msg.MemberOnlineAmount, false)
	if err != nil {
		logs.WithError(err).Error("Add member error")
		return nil, err
	}
	logs.Debug("ap10")
	fromBalance := k.BankKeeper.GetAllBalances(ctx, fromAddr)

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeAgreeJoinClusterApply,
		sdk.NewAttribute(types.AttributeSendeer, msg.FromAddress),
		sdk.NewAttribute(types.AttributeSenderBalances, fromBalance.String()),
		sdk.NewAttribute(types.AttributeKeyClusterChatId, msg.ClusterId),
		sdk.NewAttribute(types.AttributeKeyClusterTrueId, cluster.ClusterId),
		sdk.NewAttribute(types.AttributeKeyMemberAddr, msg.MemberAddress),
		sdk.NewAttribute(types.AttributeKeyIndexNum, msg.IndexNum),
	))

	return &types.MsgEmptyResponse{}, nil
}

func (k msgServer) CreateClusterAddMembers(goCtx context.Context, msg *types.MsgCreateClusterAddMembers) (*types.MsgEmptyResponse, error) {
	logs := core.BuildLog(core.GetStructFuncName(k), core.LmChainDaoKeeper)
	ctx := sdk.UnwrapSDKContext(goCtx)

	fromAccAddr, err := sdk.AccAddressFromBech32(msg.FromAddress)
	if err != nil {
		return nil, core.ErrAddressFormat
	}

	
	membersAddrs := make([]string, 0)
	for _, member := range msg.Members {
		membersAddrs = append(membersAddrs, member.MemberAddress)
	}

	ok, err := k.ValidateGatewaySign(ctx, fromAccAddr, msg.MemberOnlineAmount, membersAddrs, msg.GateAddress, msg.GatewaySign, false)
	if err != nil {
		logs.WithError(err).Error("CreateClusterAddMembers ValidateGatewaySign Error")
		return nil, core.ErrGatewaySign
	}
	if !ok {
		logs.Error("CreateClusterAddMembers ValidateGatewaySign Error")
		return nil, core.ErrGatewaySign
	}

	
	c, err := k.GetClusterByChatId(ctx, msg.ClusterId)
	if err == nil || c.ClusterId != "" {
		return nil, core.GetClusterExisted
	}

	
	if msg.GateAddress != "" && msg.ChatAddress != "" {
		err = k.ChatKeeper.Register(ctx, msg.FromAddress, msg.ChatAddress, msg.GateAddress, nil)
		if err != nil {
			return nil, err
		}
	}

	
	params := k.GetParams(ctx)

	err = k.ValidateDaoRatio(params, msg.ClusterDaoRatio)
	if err != nil {
		return nil, err
	}

	
	clusterTrueId, err := k.Keeper.CreateCluster(ctx, msg, msg.Metadata)
	if err != nil {
		return nil, err
	}

	
	if msg.BurnAmount.Add(msg.FreezeAmount).LT(sdk.NewDecFromInt(k.GetParams(ctx).ClusterLevels[1].BurnAmount)) {
		return nil, core.ErrCreateClusterBurn
	}

	
	_, err = k.BurnGetPower(ctx, fromAccAddr, fromAccAddr, clusterTrueId, msg.BurnAmount, msg.FreezeAmount, true)
	if err != nil {
		return nil, err
	}

	
	
	clusterId, err := k.GetClusterId(ctx, msg.ClusterId)
	if err != nil {
		return nil, err
	}

	if msg.Members != nil {
		err = k.AddMember(ctx, clusterId, msg.FromAddress, msg.Members, msg.MemberOnlineAmount, true)
		if err != nil {
			logs.WithField("clusterId:", msg.ClusterId).Error("ClusterAddMembers AddMember Error")
			return nil, err
		}

		
		fromBalances := k.BankKeeper.GetAllBalances(ctx, fromAccAddr)

		
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeAddMembers,
				sdk.NewAttribute(types.AttributeSendeer, msg.FromAddress),
				sdk.NewAttribute(types.AttributeSenderBalances, fromBalances.String()),
				sdk.NewAttribute(types.AttributeClusterId, msg.ClusterId),
				sdk.NewAttribute(types.AttributeKeyClusterTrueId, clusterId),
			),
		)
	}

	return &types.MsgEmptyResponse{}, nil
}

func (k msgServer) ReturnRedPacket(goCtx context.Context, msg *types.MsgReturnRedPacket) (*types.MsgEmptyResponse, error) {
	logs := core.BuildLog(core.GetStructFuncName(k), core.LmChainDaoKeeper)

	ctx := sdk.UnwrapSDKContext(goCtx)

	
	redPacket, err := k.GetRedPacketInfo(ctx, msg.Redpacketid)
	if err != nil {
		return nil, err
	}

	
	if redPacket.Remain().IsZero() || int64(len(redPacket.Receive)) == redPacket.Count {
		logs.Error("ReturnRedPacket ErrRedPacketCollected")
		return nil, core.ErrRedPacketCollected
	}

	
	if ctx.BlockHeight() < redPacket.EndBlock {
		logs.Error("ReturnRedPacket EndBlock Error")
		return nil, core.ErrRedPacketEndBlock
	}

	
	if msg.Fromaddress != redPacket.Sender {
		logs.Error("ReturnRedPacket Sender Error")
		return nil, core.ErrRedPacketSender
	}

	fromAddr, err := sdk.AccAddressFromBech32(msg.Fromaddress)
	if err != nil {
		logs.WithError(err).Error("ReturnRedPacket From Address Format Error")
		return nil, core.ErrAddressFormat
	}

	
	returnAmount, err := k.ReturnRedPacketLogic(ctx, fromAddr, redPacket)
	if err != nil {
		logs.WithError(err).Error("ReturnRedPacketLogic  Error")
		return nil, err
	}

	
	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeReturnRedPacket,
		sdk.NewAttribute(types.AttributeSendeer, fromAddr.String()),
		sdk.NewAttribute(types.AttributeSenderBalances, k.BankKeeper.GetAllBalances(ctx, fromAddr).String()),
		sdk.NewAttribute(types.AttributeKeyAmount, returnAmount.String()), 
		sdk.NewAttribute(types.AttributeKeyRedPacketId, redPacket.Id),     
	))

	return &types.MsgEmptyResponse{}, nil
}

func (k msgServer) OpenRedPacket(goCtx context.Context, msg *types.MsgOpenRedPacket) (*types.MsgEmptyResponse, error) {

	logs := core.BuildLog(core.GetStructFuncName(k), core.LmChainDaoKeeper)

	ctx := sdk.UnwrapSDKContext(goCtx)

	
	redPacket, err := k.GetRedPacketInfo(ctx, msg.Redpacketid)
	if err != nil {
		logs.WithError(err).Error("OpenRedPacket GetRedPacketInfo Error")
		return nil, err
	}

	
	fromAddr, err := sdk.AccAddressFromBech32(msg.Fromaddress)
	if err != nil {
		logs.WithError(err).Error("OpenRedPacket From Address Error")
		return &types.MsgEmptyResponse{}, core.ErrAddressFormat
	}

	
	cluster, err := k.GetCluster(ctx, redPacket.ClusterTrueId)
	if err != nil {
		logs.WithError(err).Error("OpenRedPacket Get Cluster Error")
		return nil, err
	}

	if _, ok := cluster.ClusterDeviceMembers[msg.Fromaddress]; !ok {
		return nil, core.ErrRedPacketClusters
	}

	
	if ctx.BlockHeight() > redPacket.EndBlock {
		return nil, core.ErrRedPacketExpired
	}

	
	if int64(len(redPacket.Receive)) >= redPacket.Count {
		return nil, core.ErrRedPacketCollected
	}

	
	if redPacket.IsReceived(msg.Fromaddress) {
		return nil, core.ErrRedPacketRepeat
	}

	
	if redPacket.Remain().LTE(sdk.ZeroInt()) {
		return nil, core.ErrRedPacketInsufficient
	}

	
	_, err = k.ReceiveRedPacket(ctx, fromAddr, redPacket, cluster.ClusterChatId)

	return &types.MsgEmptyResponse{}, nil
}

func (k msgServer) RedPacket(goCtx context.Context, msg *types.MsgRedPacket) (*types.MsgRedPacketResponse, error) {

	logs := core.BuildLog(core.GetStructFuncName(k), core.LmChainDaoKeeper)

	ctx := sdk.UnwrapSDKContext(goCtx)

	
	fromAddr, err := sdk.AccAddressFromBech32(msg.Fromaddress)
	if err != nil {
		return &types.MsgRedPacketResponse{}, core.ErrAddressFormat
	}

	
	cluster, err := k.GetClusterByChatId(ctx, msg.Clusterid)
	if err != nil {
		return &types.MsgRedPacketResponse{}, err
	}

	
	amountLimit, _ := new(big.Int).SetString(core.MaxRedPacketAmount, 10)
	if msg.Amount.Amount.GT(sdk.NewIntFromBigInt(amountLimit)) {
		return &types.MsgRedPacketResponse{}, core.ErrRedPacketAmountMax
	}

	
	if msg.Redtype == types.RedPacketTypeLucky && msg.Amount.Amount.Quo(sdk.NewInt(msg.Count)).LT(sdk.NewInt(1)) {
		return &types.MsgRedPacketResponse{}, core.ErrRedPacketAmountMin
	}

	
	
	if msg.Count <= 0 || msg.Count > int64(len(cluster.ClusterDeviceMembers)) {
		return &types.MsgRedPacketResponse{}, core.ErrRedPacketCount
	}

	
	txBytes := ctx.TxBytes()
	txHash := hex.EncodeToString(tmhash.Sum(txBytes))

	redPacketTotalAmount := sdk.ZeroInt()

	if msg.Redtype == types.RedPacketTypeNormal { 
		redPacketTotalAmount = msg.Amount.Amount.Mul(sdk.NewInt(msg.Count))
	} else if msg.Redtype == types.RedPacketTypeLucky { 
		redPacketTotalAmount = msg.Amount.Amount
	} else {
		return &types.MsgRedPacketResponse{}, core.ErrRedPacketType
	}

	
	err = k.BankKeeper.SendCoinsFromAccountToModule(ctx, fromAddr, types.ModuleName, sdk.NewCoins(sdk.NewCoin(msg.Amount.Denom, redPacketTotalAmount)))
	if err != nil {
		logs.WithError(err).Error("RedPacket amount insufficient")
		return &types.MsgRedPacketResponse{}, err
	}

	err = k.InitRedPacketInfo(ctx, fromAddr, cluster.ClusterId, txHash, msg.Amount, msg.Count, msg.Redtype)
	if err != nil {
		logs.WithError(err).Error("SetRedPacketInfo Error")
		return &types.MsgRedPacketResponse{}, err
	}

	fromBalance := k.BankKeeper.GetAllBalances(ctx, fromAddr)

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeRedPacket,
		sdk.NewAttribute(types.AttributeSendeer, msg.Fromaddress),
		sdk.NewAttribute(types.AttributeSenderBalances, fromBalance.String()),
		sdk.NewAttribute(types.AttributeKeyClusterChatId, msg.Clusterid),
		sdk.NewAttribute(types.AttributeKeyRedPacketId, txHash),
		sdk.NewAttribute(types.AttributeKeyAmount, msg.Amount.Amount.String()),
		sdk.NewAttribute(types.AttributeKeyDenom, msg.Amount.Denom),
		sdk.NewAttribute(types.AttributeKeyRedPacketCount, strconv.FormatInt(msg.Count, 10)),
		sdk.NewAttribute(types.AttributeKeyRedPacketType, strconv.FormatInt(msg.Redtype, 10)),
		sdk.NewAttribute(types.AttributeKeyRedPacketEndBlock, strconv.FormatInt(ctx.BlockHeight()+core.RedPacketTimeOut, 10)),
	))

	return &types.MsgRedPacketResponse{
		Clusterid: txHash,
	}, nil
}



func (k msgServer) AgreeJoinCluster(goCtx context.Context, msg *types.MsgAgreeJoinCluster) (*types.MsgEmptyResponse, error) {
	logs := core.BuildLog(core.GetStructFuncName(k), core.LmChainDaoKeeper)

	logs.Debug("AgreeJoinCluster start-------------")
	logs.Debug("params:")
	logs.Debug("msg.FromAddress", msg.FromAddress)
	logs.Debug("msg.ClusterId", msg.ClusterId)
	logs.Debug("msg.IndexNum", msg.IndexNum)
	logs.Debug("msg.ChatAddress", msg.ChatAddress)
	logs.Debug("msg.MemberOnlineAmount", msg.MemberOnlineAmount)
	logs.Debug("msg.GatewayAddress", msg.GatewayAddress)
	logs.Debug("msg.GatewaySign", msg.GatewaySign)
	logs.Debug("msg.Sign", msg.Sign)

	ctx := sdk.UnwrapSDKContext(goCtx)

	fromAccAddr, err := sdk.AccAddressFromBech32(msg.FromAddress)
	if err != nil {
		return &types.MsgEmptyResponse{}, core.ErrAddressFormat
	}

	logs.Debug("Join1")

	
	ok, err := k.ValidateGatewaySign(ctx, fromAccAddr, msg.MemberOnlineAmount, []string{msg.FromAddress}, msg.GatewayAddress, msg.GatewaySign, true)
	if err != nil {
		logs.WithError(err).Error("AgreeJoinCluster ValidateGatewaySign Error")
		return nil, core.ErrGatewaySign
	}

	logs.Debug("Join2")
	if !ok {
		logs.Error("AgreeJoinCluster ValidateGatewaySign Error")
		return nil, core.ErrGatewaySign
	}

	logs.Debug("Join3")
	
	signBytes, err := hex.DecodeString(msg.Sign)
	if err != nil {
		logs.WithError(err).Error("sign DecodeString error")
		return &types.MsgEmptyResponse{}, err
	}

	logs.Debug("Join4")

	
	clusterTrueId, err := k.GetClusterId(ctx, msg.ClusterId)
	if err != nil {
		logs.WithError(err).Error("GetClusterId error")
		return &types.MsgEmptyResponse{}, err
	}

	logs.Debug("Join5")

	
	cluster, err := k.GetCluster(ctx, clusterTrueId)
	if err != nil {
		logs.WithError(err).Error("GetCluster error")
		return &types.MsgEmptyResponse{}, err
	}

	logs.Debug("Join6")

	
	if _, ok := cluster.ClusterDeviceMembers[msg.FromAddress]; ok {
		return &types.MsgEmptyResponse{}, core.ErrMemberAlreadyExist
	}

	logs.Debug("Join7")

	

	
	dataUnSign := []byte(msg.FromAddress + "," + msg.ClusterId)

	
	
	comPub, err := util.GetPubKeyFromSign(signBytes, dataUnSign)
	if err != nil {
		logs.WithError(err).Error("GetPubKeyFromSign error")
		return &types.MsgEmptyResponse{}, err
	}

	logs.Debug("Join8")

	
	ethOwnerAddr := common.HexToAddress(comPub.Address().String())
	logs.Debug("Join9")

	dstOwnerAddr := sdk.AccAddress(ethOwnerAddr.Bytes())
	logs.Debug("Join10")
	if dstOwnerAddr.String() != cluster.ClusterOwner {
		return &types.MsgEmptyResponse{}, core.ErrClusterOwnerErr
	}

	logs.Debug("Join11")

	
	err = k.AddMember(
		ctx,
		clusterTrueId,
		cluster.ClusterOwner,
		[]types.Members{
			{
				MemberAddress: msg.FromAddress,
				IndexNum:      msg.IndexNum,
				ChatAddress:   msg.ChatAddress,
			},
		},
		msg.MemberOnlineAmount,
		false,
	)
	if err != nil {
		logs.WithError(err).Error("AddMember error")
		return &types.MsgEmptyResponse{}, err
	}

	logs.Debug("Join12")

	fromBalance := k.BankKeeper.GetAllBalances(ctx, fromAccAddr)

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeAgreeJoinCluster,
		sdk.NewAttribute(types.AttributeSendeer, msg.FromAddress),
		sdk.NewAttribute(types.AttributeSenderBalances, fromBalance.String()),
		sdk.NewAttribute(types.AttributeKeyClusterChatId, msg.ClusterId),
		sdk.NewAttribute(types.AttributeKeyClusterTrueId, cluster.ClusterId),
		sdk.NewAttribute(types.AttributeKeyClusterOwner, cluster.ClusterOwner),
	))

	return &types.MsgEmptyResponse{}, nil
}

func (k msgServer) ClusterChangeDaoRatio(goCtx context.Context, msg *types.MsgClusterChangeDaoRatio) (*types.MsgEmptyResponse, error) {
	logs := core.BuildLog(core.GetStructFuncName(k), core.LmChainDaoKeeper)
	ctx := sdk.UnwrapSDKContext(goCtx)
	cluster, err := k.GetClusterByChatId(ctx, msg.GetClusterId())
	if err != nil {
		logs.WithError(err).WithField("clusterId", msg.ClusterId).Error("GetClusterByChatId err")
		return nil, err
	}
	
	if msg.FromAddress != cluster.ClusterOwner {
		logs.Error("ClusterChangeSalaryRatio ClusterOwner Error")
		return nil, core.ErrClusterOwnerErr
	}

	
	if cluster.ClusterSalaryRatioUpdateHeight.DaoRatioUpdateHeight != 0 && ctx.BlockHeight()-cluster.ClusterSalaryRatioUpdateHeight.DaoRatioUpdateHeight < core.DayBlockNum {
		return nil, core.ErrClusterConfigChange
	}

	
	params := k.GetParams(ctx)
	err = k.ValidateDaoRatio(params, msg.DaoRatio)
	if err != nil {
		logs.WithError(err).WithField("ratio", msg.DaoRatio).Error("dao ratio check err")
		return nil, err
	}

	cluster.ClusterDaoRatio = msg.DaoRatio
	cluster.ClusterSalaryRatioUpdateHeight.DaoRatioUpdateHeight = ctx.BlockHeight()
	err = k.SetDeviceCluster(ctx, cluster)
	if err != nil {
		logs.WithError(err).WithField("clusterId", msg.ClusterId).Error("SetClusterByChatId err")
		return nil, err
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.ChangeClusterDaoRatio,
		sdk.NewAttribute(types.AttributeSendeer, msg.FromAddress),
	))

	return &types.MsgEmptyResponse{}, nil
}













func (k msgServer) ClusterAd(goCtx context.Context, msg *types.MsgClusterAd) (*types.MsgEmptyResponse, error) {
	logs := core.BuildLog(core.GetStructFuncName(k), core.LmChainDaoKeeper)
	ctx := sdk.UnwrapSDKContext(goCtx)
	fromAddr, err := sdk.AccAddressFromBech32(msg.FromAddress)
	if err != nil {
		return nil, err
	}
	err = k.clusterAd(ctx, msg.ClusterId, fromAddr)
	if err != nil {
		logs.WithError(err).Error("clusterAd Err")
		return nil, err
	}
	return nil, nil
}

func (k msgServer) PersonDvmApprove(goCtx context.Context, msg *types.MsgPersonDvmApprove) (*types.MsgEmptyResponse, error) {
	logs := core.BuildLog(core.GetStructFuncName(k), core.LmChainDaoKeeper)
	ctx := sdk.UnwrapSDKContext(goCtx)

	err := k.DvmApprove(ctx, msg.ApproveAddress, msg.ClusterId, msg.ApproveEndBlock, msg.FromAddress)
	if err != nil {
		logs.WithError(err).Error("PersonDvmApprove Error")
		return nil, err
	}

	return nil, nil
}


func (k msgServer) ClusterPowerApprove(goCtx context.Context, msg *types.MsgClusterPowerApprove) (*types.MsgEmptyResponse, error) {
	logs := core.BuildLog(core.GetStructFuncName(k), core.LmChainDaoKeeper)
	ctx := sdk.UnwrapSDKContext(goCtx)
	var err error
	
	if !k.contractKeeper.IsContract(ctx, msg.ApproveAddress) {
		return nil, core.ErrContractAddress
	}
	clusterId := msg.ClusterId
	if strings.Contains(msg.ClusterId, ".") {
		
		clusterId, err = k.GetClusterId(ctx, msg.ClusterId)
		if err != nil {
			return nil, err
		}
	}
	cluster, err := k.GetCluster(ctx, clusterId)
	if err != nil {
		logs.WithError(err).Error("ClusterPowerApprove GetCluster Err")
		return nil, err
	}
	if cluster.ClusterVotePolicy != msg.FromAddress {
		return nil, core.ErrVotePolicyAddress
	}
	
	if _, ok := cluster.ClusterPowerMembers[cluster.ClusterDaoPool]; !ok {
		return nil, core.ErrMemberNotExist
	}
	curApprove, err := k.GetClusterApproveInfo(ctx, cluster.ClusterId, cluster.ClusterDaoPool)
	if err != nil {
		logs.WithError(err).Error("ClusterPowerApprove GetClusterIdApproveInfo Err")
		return nil, err
	}
	if curApprove.EndBlock > ctx.BlockHeight() {
		return nil, core.ErrApproveNotEnd
	}
	blockNum, err := strconv.ParseInt(msg.ApproveEndBlock, 10, 64)
	if err != nil {
		logs.WithError(err).Error("ClusterPowerApprove ApproveEndBlock Err")
		return nil, err
	}
	endBlock := ctx.BlockHeight() + blockNum
	approvePower := types.ApprovePower{
		ClusterId: clusterId,
		Address:   cluster.ClusterDaoPool,
		IsDaoPool: true,
		EndBlock:  endBlock,
	}
	err = k.SetContractApproveInfo(ctx, msg.ApproveAddress, approvePower)
	if err != nil {
		return nil, err
	}
	approve := types.ClusterCurApprove{
		ApproveAddress: strings.ToLower(msg.ApproveAddress),
		EndBlock:       endBlock,
	}
	err = k.SetClusterIdApproveInfo(ctx, clusterId, cluster.ClusterDaoPool, approve)
	if err != nil {
		return nil, err
	}
	
	err = k.SettlementRouteBridgingReward(ctx, clusterId, msg.FromAddress, true)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (k msgServer) UpdateAdmin(goCtx context.Context, msg *types.MsgUpdateAdmin) (*types.MsgEmptyResponse, error) {

	logs := core.BuildLog(core.GetStructFuncName(k), core.LmChainDaoKeeper)
	ctx := sdk.UnwrapSDKContext(goCtx)

	
	clusterId, err := k.GetClusterId(ctx, msg.ClusterId)
	if err != nil {
		logs.WithError(err).Error("UpdateAdmin GetClusterId Err")
		return nil, err
	}

	
	cluster, err := k.GetCluster(ctx, clusterId)
	if err != nil {
		logs.WithError(err).Error("UpdateAdmin GetCluster Err")
		return nil, err
	}

	
	err = k.ValidateClusterAdmin(ctx, cluster, msg.ClusterAdminList)
	if err != nil {
		logs.WithError(err).Error("UpdateAdmin ValidateClusterAdmin Err")
		return nil, err
	}

	
	newAdminMap := make(map[string]struct{})
	for _, admin := range msg.ClusterAdminList {
		newAdminMap[admin] = struct{}{}
	}
	cluster.ClusterAdminList = newAdminMap

	
	err = k.SetDeviceCluster(ctx, cluster)
	if err != nil {
		logs.WithError(err).Error("UpdateAdmin SetDeviceCluster Err")
		return nil, err
	}

	
	fromAccAddr, err := sdk.AccAddressFromBech32(msg.FromAddress)
	if err != nil {
		return nil, core.ErrAddressFormat
	}

	fromBalance := k.BankKeeper.GetAllBalances(ctx, fromAccAddr)

	
	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypUpdateAdmin,
		sdk.NewAttribute(types.AttributeSendeer, msg.FromAddress),
		sdk.NewAttribute(types.AttributeSenderBalances, fromBalance.String()),
	))
	
	err = k.SettlementRouteBridgingReward(ctx, clusterId, msg.FromAddress, true)
	if err != nil {
		return nil, err
	}
	return nil, nil
}


func (k msgServer) ClusterChangeName(goCtx context.Context, msg *types.MsgClusterChangeName) (*types.MsgEmptyResponse, error) {
	logs := core.BuildLog(core.GetStructFuncName(k), core.LmChainDaoKeeper)
	ctx := sdk.UnwrapSDKContext(goCtx)

	
	cluster, err := k.GetClusterByChatId(ctx, msg.ClusterId)
	if err != nil {
		logs.WithError(err).WithField("clusterChatId:", msg.ClusterId).Error("ClusterChangeName GetClusterByChatId Error")
		return nil, err
	}

	
	if cluster.ClusterOwner != msg.FromAddress {
		logs.WithFields(
			logrus.Fields{
				"from":         msg.FromAddress,
				"ClusterOwner": cluster.ClusterOwner,
			},
		).Error("ClusterChangeName Permmision Error")
		return nil, core.ErrOwnerPermission
	}

	
	cluster.ClusterName = msg.ClusterName

	
	err = k.SetDeviceCluster(ctx, cluster)
	if err != nil {
		logs.WithError(err).Error("ClusterChangeName SetDeviceCluster Error")
		return nil, core.ErrSetCluster
	}

	
	fromAccAddr, err := sdk.AccAddressFromBech32(msg.FromAddress)
	if err != nil {
		return nil, core.ErrAddressFormat
	}

	fromBalance := k.BankKeeper.GetAllBalances(ctx, fromAccAddr)

	
	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeChangeName,
		sdk.NewAttribute(types.AttributeSendeer, msg.FromAddress),             
		sdk.NewAttribute(types.AttributeSenderBalances, fromBalance.String()), 
	))

	return nil, nil
}

func (k msgServer) ClusterDeleteMembers(goCtx context.Context, msg *types.MsgDeleteMembers) (*types.MsgEmptyResponse, error) {
	logs := core.BuildLog(core.GetStructFuncName(k), core.LmChainDaoKeeper)

	ctx := sdk.UnwrapSDKContext(goCtx)

	fromAccAddr, err := sdk.AccAddressFromBech32(msg.FromAddress)
	if err != nil {
		logs.WithError(err).WithField("from:", msg.FromAddress).Error("delete members AccAddressFromBech32 Error")
		return nil, core.ErrAddressFormat
	}

	
	ok, err := k.ValidateGatewaySign(ctx, fromAccAddr, msg.MemberOnlineAmount, msg.Members, msg.GatewayAddress, msg.GatewaySign, false)
	if err != nil {
		logs.WithError(err).WithField("from:", msg.FromAddress).Error("delete members ValidateGatewaySign Error")
		return nil, err
	}
	if !ok {
		logs.WithField("from:", msg.FromAddress).Error("delete members gateway sign error")
		return nil, core.ErrGatewaySign
	}

	
	
	clusterId, err := k.GetClusterId(ctx, msg.ClusterId)
	if err != nil {
		logs.
			WithError(err).
			WithField("cluster_id:", msg.ClusterId).
			WithField("from:", msg.FromAddress).
			Error("ClusterMemberDelete GetClusterChatId Error")
		return nil, err
	}

	
	cluster, err := k.GetCluster(ctx, clusterId)
	if err != nil {
		logs.
			WithError(err).
			WithField("cluster_id:", msg.ClusterId).
			WithField("from:", msg.FromAddress).
			Error("ClusterMemberDelete GetClusterByChatId Error")
		return nil, err
	}

	
	clusterAuth := k.IsOwnerOrAdmin(cluster, msg.FromAddress)
	if !clusterAuth {
		return nil, core.ErrClusterPermission
	}
	
	groupMembers := []group.MemberRequest{}

	for _, member := range msg.Members {
		
		if member == cluster.ClusterOwner {
			logs.
				WithError(err).
				WithField("cluster_id:", msg.ClusterId).
				WithField("from:", msg.FromAddress).
				WithField("ClusterOwner:", cluster.ClusterOwner).
				Error("ClusterMemberDelete owner delete")
			return nil, core.ErrOwnerCannotExit
		}

		
		if k.IsAdmin(cluster, msg.FromAddress) && k.IsAdmin(cluster, member) {
			return nil, core.ErrClusterPermission
		}

		
		if _, ok := cluster.ClusterDeviceMembers[member]; !ok {
			logs.
				WithError(err).
				WithField("cluster_id:", msg.ClusterId).
				WithField("from:", msg.FromAddress).
				Error("ClusterMemberDelete member not found")
			return nil, core.ErrMemberNotExist
		}

		
		delete(cluster.ClusterDeviceMembers, member)

		
		p, err := k.GetPersonClusterInfo(ctx, member)
		if err != nil {
			logs.
				WithError(err).
				WithField("cluster_id:", msg.ClusterId).
				WithField("from:", msg.FromAddress).
				Error("ClusterMemberDelete GetPersonClusterInfo Error")
			return nil, err
		}

		
		delete(p.Device, cluster.ClusterId)

		
		k.DeleteDeviceStartingInfo(ctx, cluster.ClusterId, member)
		
		err = k.SetPersonClusterInfo(ctx, p)
		if err != nil {
			logs.
				WithError(err).
				WithField("from:", msg.FromAddress).
				Error("ClusterMemberDelete SetPersonClusterInfo Error")
			return nil, err
		}

		groupMembers = append(groupMembers, group.MemberRequest{Address: member, Weight: "0"})
	}

	
	cluster.ClusterActiveDevice = cluster.ClusterActiveDevice - msg.MemberOnlineAmount
	if cluster.ClusterActiveDevice < 0 {
		cluster.ClusterActiveDevice = 0
	}

	cluster.OnlineRatio = sdk.NewDec(cluster.ClusterActiveDevice).Quo(sdk.NewDec(int64(len(cluster.ClusterDeviceMembers))))

	
	err = k.SetDeviceCluster(ctx, cluster)
	if err != nil {
		logs.
			WithError(err).
			WithField("cluster_id:", msg.ClusterId).
			WithField("from:", msg.FromAddress).
			Error("ClusterMemberDelete SetDeviceCluster Error")
		return nil, err
	}

	fromBalance := k.BankKeeper.GetAllBalances(ctx, fromAccAddr)

	
	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeDeleteMembers,
		sdk.NewAttribute(types.AttributeSendeer, msg.FromAddress),             
		sdk.NewAttribute(types.AttributeSenderBalances, fromBalance.String()), 
	))
	if len(groupMembers) > 0 {
		err := k.UpdateGroupMembers(ctx, cluster, groupMembers)
		if err != nil {
			logs.WithError(err).WithField("group", clusterId).Error("delete members UpdateGroupMembers Error")
			return nil, err
		}
	}
	
	err = k.SettlementRouteBridgingReward(ctx, clusterId, msg.FromAddress, true)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (k msgServer) ClusterChangeSalaryRatio(goCtx context.Context, msg *types.MsgClusterChangeSalaryRatio) (*types.MsgEmptyResponse, error) {
	logs := core.BuildLog(core.GetStructFuncName(k), core.LmChainDaoKeeper)
	ctx := sdk.UnwrapSDKContext(goCtx)

	
	clusterId, err := k.GetClusterId(ctx, msg.ClusterId)
	if err != nil {
		logs.WithError(err).WithField("clusterChatId:", msg.ClusterId).Error("ClusterChangeSalaryRatio GetClusterId Error")
		return nil, err
	}

	
	cluster, err := k.GetCluster(ctx, clusterId)
	if err != nil {
		logs.WithError(err).WithField("clusterId:", clusterId).Error("ClusterChangeSalaryRatio GetCluster Error")
		return nil, err
	}
	
	if cluster.ClusterSalaryRatioUpdateHeight.SalaryRatioUpdateHeight != 0 && ctx.BlockHeight()-cluster.ClusterSalaryRatioUpdateHeight.SalaryRatioUpdateHeight < core.DayBlockNum {
		return nil, core.ErrClusterConfigChange
	}
	
	if msg.FromAddress != cluster.ClusterOwner {
		logs.Error("ClusterChangeSalaryRatio ClusterOwner Error")
		return nil, core.ErrClusterOwnerErr
	}
	
	params := k.GetParams(ctx)
	err = k.ValidateSalaryRatio(params, msg.SalaryRatio)
	if err != nil {
		logs.
			WithError(err).
			WithField("clusterChatId:", msg.ClusterId).
			WithField("from:", msg.FromAddress).
			WithField("SalaryRatio:", msg.SalaryRatio.String()).
			Error("ClusterChangeSalaryRatio Error")
		return nil, err
	}
	
	cluster.ClusterSalaryRatio = msg.SalaryRatio
	cluster.ClusterSalaryRatioUpdateHeight.SalaryRatioUpdateHeight = ctx.BlockHeight()
	err = k.SetDeviceCluster(ctx, cluster)
	if err != nil {
		logs.
			WithError(err).
			WithField("clusterChatId:", msg.ClusterId).
			WithField("from:", msg.FromAddress).
			Error("ClusterChangeSalaryRatio SetDeviceCluster Error")
		return nil, err
	}

	fromAcc, err := sdk.AccAddressFromBech32(msg.FromAddress)
	if err != nil {
		return nil, core.ErrAddressFormat
	}
	fromBalance := k.BankKeeper.GetAllBalances(ctx, fromAcc)
	
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeChangeSalaryRatio,
			sdk.NewAttribute(types.AttributeSendeer, msg.FromAddress),
			sdk.NewAttribute(types.AttributeSenderBalances, fromBalance.String()),
		),
	)

	return nil, nil
}

func (k msgServer) ClusterChangeDvmRatio(goCtx context.Context, msg *types.MsgClusterChangeDvmRatio) (*types.MsgEmptyResponse, error) {
	logs := core.BuildLog(core.GetStructFuncName(k), core.LmChainDaoKeeper)
	ctx := sdk.UnwrapSDKContext(goCtx)

	
	clusterId, err := k.GetClusterId(ctx, msg.ClusterId)
	if err != nil {
		logs.WithError(err).WithField("clusterChatId:", msg.ClusterId).Error("MsgClusterChangeDvmRatio GetClusterId Error")
		return nil, err
	}

	
	cluster, err := k.GetCluster(ctx, clusterId)
	if err != nil {
		logs.WithError(err).WithField("clusterId:", clusterId).Error("MsgClusterChangeDvmRatio GetCluster Error")
		return nil, err
	}
	
	if cluster.ClusterSalaryRatioUpdateHeight.DvmRatioUpdateHeight != 0 && ctx.BlockHeight()-cluster.ClusterSalaryRatioUpdateHeight.DvmRatioUpdateHeight < core.DayBlockNum {
		return nil, core.ErrClusterConfigChange
	}
	
	if msg.FromAddress != cluster.ClusterOwner {
		logs.Error("ClusterChangeSalaryRatio ClusterOwner Error")
		return nil, core.ErrClusterOwnerErr
	}
	
	params := k.GetParams(ctx)
	err = k.ValidateDvmRatio(params, msg.DvmRatio)
	if err != nil {
		logs.
			WithError(err).
			WithField("clusterChatId:", msg.ClusterId).
			WithField("from:", msg.FromAddress).
			WithField("DvmRatio:", msg.DvmRatio.String()).
			Error("MsgClusterChangeDvmRatio Error")
		return nil, err
	}
	
	cluster.ClusterDvmRatio = msg.DvmRatio
	cluster.ClusterSalaryRatioUpdateHeight.DvmRatioUpdateHeight = ctx.BlockHeight()
	err = k.SetDeviceCluster(ctx, cluster)
	if err != nil {
		logs.
			WithError(err).
			WithField("clusterChatId:", msg.ClusterId).
			WithField("from:", msg.FromAddress).
			Error("MsgClusterChangeDvmRatio SetDeviceCluster Error")
		return nil, err
	}

	fromAcc, err := sdk.AccAddressFromBech32(msg.FromAddress)
	if err != nil {
		return nil, core.ErrAddressFormat
	}
	fromBalance := k.BankKeeper.GetAllBalances(ctx, fromAcc)
	
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeChangeDvmRatio,
			sdk.NewAttribute(types.AttributeSendeer, msg.FromAddress),
			sdk.NewAttribute(types.AttributeSenderBalances, fromBalance.String()),
		),
	)

	return nil, nil
}

func (k msgServer) ClusterChangeId(goCtx context.Context, msg *types.MsgClusterChangeId) (*types.MsgEmptyResponse, error) {
	logs := core.BuildLog(core.GetStructFuncName(k), core.LmChainDaoKeeper)
	ctx := sdk.UnwrapSDKContext(goCtx)

	
	clusterId, err := k.GetClusterId(ctx, msg.ClusterId)
	if err != nil {
		logs.WithError(err).WithField("clusterChatId:", msg.ClusterId).Error("ClusterChangeId GetClusterId Error")
		return nil, err
	}

	
	cluster, err := k.GetCluster(ctx, clusterId)
	if err != nil {
		logs.WithError(err).WithField("clusterTrueId:", clusterId).Error("ClusterChangeId GetCluster Error")
		return nil, err
	}

	
	if msg.FromAddress != cluster.ClusterOwner {
		logs.WithError(err).WithField("clusterTrueId:", clusterId).Error("ClusterChangeId GetCluster Error")
		return nil, core.ErrOwnerPermission
	}
	
	k.changeClusterChatId(ctx, cluster.ClusterChatId, msg.NewClusterId, cluster.ClusterId)

	
	cluster.ClusterChatId = msg.NewClusterId
	err = k.SetDeviceCluster(ctx, cluster)
	if err != nil {
		logs.WithError(err).Error("ClusterChangeId SetDeviceCluster Err")
		return nil, err
	}

	fromAccAddr, err := sdk.AccAddressFromBech32(msg.FromAddress)
	if err != nil {
		return nil, core.ErrAddressFormat

	}

	
	fromBalances := k.BankKeeper.GetAllBalances(ctx, fromAccAddr)

	
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeChangeClusterId,
			sdk.NewAttribute(types.AttributeSendeer, msg.FromAddress),
			sdk.NewAttribute(types.AttributeSenderBalances, fromBalances.String()),
		),
	)
	return nil, nil
}

func (k msgServer) ClusterMemberExit(goCtx context.Context, msg *types.MsgClusterMemberExit) (*types.MsgEmptyResponse, error) {
	logs := core.BuildLog(core.GetStructFuncName(k), core.LmChainDaoKeeper)

	ctx := sdk.UnwrapSDKContext(goCtx)

	fromAccAddr, err := sdk.AccAddressFromBech32(msg.FromAddress)
	if err != nil {
		return nil, core.ErrAddressFormat

	}

	
	ok, err := k.ValidateGatewaySign(ctx, fromAccAddr, msg.MemberOnlineAmount, []string{msg.FromAddress}, msg.GatewayAddress, msg.GatewaySign, false)
	if err != nil {
		logs.
			WithError(err).
			WithField("from:", msg.FromAddress).
			WithField("cluster_id:", msg.ClusterId).
			Error("ClusterMemberExit ValidateGatewaySign Error")
		return nil, err
	}
	if !ok {
		logs.
			WithError(err).
			WithField("from:", msg.FromAddress).
			WithField("cluster_id:", msg.ClusterId).
			Error("ClusterMemberExit ValidateGatewaySign Error")
		return nil, core.ErrGatewaySign
	}

	
	
	clusterId, err := k.GetClusterId(ctx, msg.ClusterId)
	if err != nil {
		logs.
			WithError(err).
			WithField("cluster_id:", msg.ClusterId).
			WithField("from:", msg.FromAddress).
			Error("ClusterMemberExit GetClusterChatId Error")
		return nil, err
	}

	
	cluster, err := k.GetCluster(ctx, clusterId)
	if err != nil {
		logs.
			WithError(err).
			WithField("cluster_id:", msg.ClusterId).
			WithField("from:", msg.FromAddress).
			Error("ClusterMemberExit GetClusterByChatId Error")
		return nil, err
	}

	
	if msg.FromAddress == cluster.ClusterOwner {
		logs.
			WithError(err).
			WithField("cluster_id:", msg.ClusterId).
			WithField("from:", msg.FromAddress).
			WithField("ClusterOwner:", cluster.ClusterOwner).
			Error("ClusterMemberExit owner exit")
		return nil, core.ErrOwnerCannotExit
	}

	
	if _, ok := cluster.ClusterDeviceMembers[msg.FromAddress]; !ok {
		logs.
			WithError(err).
			WithField("cluster_id:", msg.ClusterId).
			WithField("from:", msg.FromAddress).
			Error("ClusterMemberExit member not found")
		return nil, core.ErrMemberNotExist
	}

	
	delete(cluster.ClusterDeviceMembers, msg.FromAddress)

	
	cluster.ClusterActiveDevice = cluster.ClusterActiveDevice - msg.MemberOnlineAmount
	if cluster.ClusterActiveDevice < 0 {
		cluster.ClusterActiveDevice = 0
	}

	cluster.OnlineRatio = sdk.NewDec(cluster.ClusterActiveDevice).Quo(sdk.NewDec(int64(len(cluster.ClusterDeviceMembers))))

	
	err = k.SetDeviceCluster(ctx, cluster)
	if err != nil {
		logs.
			WithError(err).
			WithField("cluster_id:", msg.ClusterId).
			WithField("from:", msg.FromAddress).
			Error("ClusterMemberExit SetDeviceCluster Error")
		return nil, err
	}
	
	k.DeleteDeviceStartingInfo(ctx, cluster.ClusterId, msg.FromAddress)

	
	p, err := k.GetPersonClusterInfo(ctx, msg.FromAddress)
	if err != nil {
		logs.
			WithError(err).
			WithField("cluster_id:", msg.ClusterId).
			WithField("from:", msg.FromAddress).
			Error("ClusterMemberExit GetPersonClusterInfo Error")
		return nil, err
	}

	
	delete(p.Device, cluster.ClusterId)

	
	err = k.SetPersonClusterInfo(ctx, p)
	if err != nil {
		logs.
			WithError(err).
			WithField("from:", msg.FromAddress).
			Error("ClusterMemberExit SetPersonClusterInfo Error")
		return nil, err
	}

	
	fromBalances := k.BankKeeper.GetAllBalances(ctx, fromAccAddr)

	
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeClusterExit,
			sdk.NewAttribute(types.AttributeSendeer, msg.FromAddress),
			sdk.NewAttribute(types.AttributeSenderBalances, fromBalances.String()),
		),
	)

	err = k.LeaveGroup(ctx, cluster, msg.FromAddress)
	if err != nil {
		logs.WithError(err).WithField("group", clusterId).Error("ClusterMemberExit LeaveGroup Error")
		return nil, err
	}
	return nil, nil
}

func (k msgServer) ThawFrozenPower(goCtx context.Context, msg *types.MsgThawFrozenPower) (*types.MsgEmptyResponse, error) {
	logs := core.BuildLog(core.GetStructFuncName(k), core.LmChainDaoKeeper)

	ctx := sdk.UnwrapSDKContext(goCtx)

	
	cluster, err := k.GetClusterByChatId(ctx, msg.ClusterId)
	if err != nil {
		logs.
			WithError(err).
			WithField("from:", msg.FromAddress).
			WithField("clusterChatId:", msg.ClusterId).
			Error("ThawFrozenPower GetClusterByChatId Error")
		return nil, err
	}

	
	if msg.GatewayAddress != "" && msg.ChatAddress != "" {
		err := k.ChatKeeper.Register(ctx, msg.FromAddress, msg.ChatAddress, msg.GatewayAddress, nil)
		if err != nil {
			return nil, err
		}
	}

	
	p, err := k.GetPersonClusterInfo(ctx, msg.FromAddress)
	if err != nil {
		logs.WithError(err).WithField("from:", msg.FromAddress).Error("ThawFrozenPower GetPersonClusterInfo Error")
		return nil, err
	}

	
	if p.FreezePower.LT(msg.ThawAmount) {
		logs.
			WithField("from:", msg.FromAddress).
			WithField("FreezeRamain:", p.FreezePower.String()).
			WithField("WantUse:", msg.ThawAmount.String()).
			Error("ThawFrozenPower insufficient freeze power")

		return nil, core.ErrFrozenInsufficient
	}

	
	if p.ActivePower.IsZero() {
		k.IncrementClusterPeriod(ctx, cluster)
		p.FirstPowerCluster = cluster.ClusterId
	} else {
		_, err = k.calculateWithdrawRewards(ctx, cluster, msg.FromAddress)
		if err != nil {
			logs.WithError(err).WithField("clusterId:", cluster.ClusterId).Error("withdrawBurnRewards error")
			return nil, err
		}
	}

	p.FreezePower = p.FreezePower.Sub(msg.ThawAmount)
	p.ActivePower = p.ActivePower.Add(msg.ThawAmount)
	p.AllBurn = p.AllBurn.Add(msg.ThawAmount)
	p.BePower[cluster.ClusterChatId] = struct{}{}

	
	err = k.SetPersonClusterInfo(ctx, p)
	if err != nil {
		logs.WithError(err).WithField("from:", msg.FromAddress).Error("ThawFrozenPower SetPersonClusterInfo Error")
		return nil, err
	}

	
	newActivePower := msg.ThawAmount
	newBurnAmount := msg.ThawAmount
	userPowerInfo, ok := cluster.ClusterPowerMembers[msg.FromAddress]
	if ok { 
		newActivePower = newActivePower.Add(userPowerInfo.ActivePower)
		newBurnAmount = newBurnAmount.Add(userPowerInfo.BurnAmount)
	}

	cluster.ClusterPowerMembers[msg.FromAddress] = types.ClusterPowerMember{
		Address:     msg.FromAddress,
		ActivePower: newActivePower,
		BurnAmount:  newBurnAmount,
	}

	
	err = k.SetDeviceCluster(ctx, cluster)
	if err != nil {
		logs.
			WithError(err).
			WithField("from:", msg.FromAddress).
			WithField("clusterChatId:", msg.ClusterId).
			Error("ThawFrozenPower SetDeviceCluster Error")

		return nil, err
	}

	
	k.InitializeGasDelegation(ctx, cluster, msg.FromAddress)

	
	err = k.AddTotalPowerAmount(ctx, msg.ThawAmount)
	if err != nil {
		logs.
			WithError(err).
			WithField("from:", msg.FromAddress).
			WithField("amount:", msg.ThawAmount).
			Error("ThawFrozenPower AddTotalPowerAmount Error")
		return nil, err
	}
	
	err = k.SubNotActivePowerAmount(ctx, msg.ThawAmount)
	if err != nil {
		logs.
			WithError(err).
			WithField("from:", msg.FromAddress).
			WithField("amount:", msg.ThawAmount).
			Error("ThawFrozenPower SubNotActivePowerAmount Error")
		return nil, err
	}

	
	err = k.AddTotalBurnAmount(ctx, msg.ThawAmount)
	if err != nil {
		logs.
			WithError(err).
			WithField("from:", msg.FromAddress).
			WithField("amount:", msg.ThawAmount).
			Error("ThawFrozenPower AddTotalBurnAmount Error")
		return nil, err
	}

	return nil, nil
}

func (k msgServer) WithdrawOwnerReward(goCtx context.Context, msg *types.MsgWithdrawOwnerReward) (*types.MsgEmptyResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	cluster, err := k.GetClusterByChatId(ctx, msg.ClusterId)
	if err != nil {
		return nil, err
	}
	err = k.withdrawAndSendRewards(ctx, cluster, cluster.ClusterDaoPool, sdk.ZeroInt())
	if err != nil {
		return nil, err
	}
	
	err = k.SettlementRouteBridgingReward(ctx, cluster.ClusterId, msg.Address, true)
	if err != nil {
		return nil, err
	}
	return &types.MsgEmptyResponse{}, nil
}

func (k msgServer) WithdrawSwapDpos(goCtx context.Context, msg *types.MsgWithdrawSwapDpos) (*types.MsgEmptyResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	deviceCluster, err := k.GetClusterByChatId(ctx, msg.ClusterId)
	if err != nil {
		return nil, err
	}
	daoPay, ok := sdkmath.NewIntFromString(msg.DaoNum)
	if !ok {
		return nil, core.ParseCoinError
	}
	addr, err := sdk.AccAddressFromBech32(msg.MemberAddress)
	if err != nil {
		return nil, core.ParseAccountError
	}
	daoBalance := k.BankKeeper.GetBalance(ctx, addr, core.BurnRewardFeeDenom)
	if daoBalance.Amount.LT(daoPay) {
		return nil, sdkerrors.ErrInsufficientFunds
	}
	err = k.withdrawAndSendRewards(ctx, deviceCluster, msg.MemberAddress, daoPay)
	if err != nil {
		return nil, err
	}
	
	err = k.SettlementRouteBridgingReward(ctx, deviceCluster.ClusterId, msg.MemberAddress, true)
	if err != nil {
		return nil, err
	}
	return &types.MsgEmptyResponse{}, nil
}

func (k msgServer) WithdrawDeviceReward(goCtx context.Context, msg *types.MsgWithdrawDeviceReward) (*types.MsgEmptyResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	deviceCluster, err := k.GetClusterByChatId(ctx, msg.ClusterId)
	if err != nil {
		return nil, err
	}
	if _, ok := deviceCluster.ClusterDeviceMembers[msg.MemberAddress]; !ok {
		return nil, core.ErrNotInCluster
	}
	_, err = k.withdrawDeviceRewards(ctx, deviceCluster, msg.MemberAddress)
	if err != nil {
		return nil, err
	}
	
	k.initializeDeviceDelegation(ctx, deviceCluster, msg.MemberAddress, 1)
	
	err = k.SettlementRouteBridgingReward(ctx, deviceCluster.ClusterId, msg.MemberAddress, true)
	if err != nil {
		return nil, err
	}
	return &types.MsgEmptyResponse{}, nil
}

func (k msgServer) ClusterAddMembers(goCtx context.Context, msg *types.MsgClusterAddMembers) (*types.MsgEmptyResponse, error) {
	logs := core.BuildLog(core.GetStructFuncName(k), core.LmChainDaoKeeper)

	ctx := sdk.UnwrapSDKContext(goCtx)

	
	clusterId, err := k.GetClusterId(ctx, msg.ClusterId)
	if err != nil {
		return nil, err
	}

	err = k.AddMember(ctx, clusterId, msg.FromAddress, msg.Members, 1, false)
	if err != nil {
		logs.WithField("clusterId:", msg.ClusterId).Error("ClusterAddMembers AddMember Error")
		return nil, err
	}

	fromAccAddr, err := sdk.AccAddressFromBech32(msg.FromAddress)
	if err != nil {
		return nil, core.ErrAddressFormat

	}

	
	fromBalances := k.BankKeeper.GetAllBalances(ctx, fromAccAddr)

	
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeAddMembers,
			sdk.NewAttribute(types.AttributeSendeer, msg.FromAddress),
			sdk.NewAttribute(types.AttributeSenderBalances, fromBalances.String()),
			sdk.NewAttribute(types.AttributeClusterId, msg.ClusterId),
			sdk.NewAttribute(types.AttributeKeyClusterTrueId, clusterId),
		),
	)
	
	err = k.SettlementRouteBridgingReward(ctx, clusterId, msg.FromAddress, true)
	if err != nil {
		return nil, err
	}
	return &types.MsgEmptyResponse{}, nil
}

func (k msgServer) CreateCluster(goCtx context.Context, msg *types.MsgCreateCluster) (*types.MsgEmptyResponse, error) {
	
	ctx := sdk.UnwrapSDKContext(goCtx)

	
	fromAccAddr, err := sdk.AccAddressFromBech32(msg.FromAddress)
	if err != nil {
		return nil, core.ErrAddressFormat
	}

	
	c, err := k.GetClusterByChatId(ctx, msg.ClusterId)
	if err == nil && c.ClusterId != "" {
		return nil, core.GetClusterExisted
	}

	

	if msg.GateAddress != "" && msg.ChatAddress != "" {
		err = k.ChatKeeper.Register(ctx, msg.FromAddress, msg.ChatAddress, msg.GateAddress, nil)
		if err != nil {
			return nil, err
		}
	}

	
	params := k.GetParams(ctx)

	err = k.ValidateDaoRatio(params, msg.ClusterDaoRatio)
	if err != nil {
		return nil, err
	}

	
	msgCreate := &types.MsgCreateClusterAddMembers{
		FromAddress:     msg.FromAddress,
		GateAddress:     msg.GateAddress,
		ClusterId:       msg.ClusterId,
		SalaryRatio:     msg.SalaryRatio,
		BurnAmount:      msg.BurnAmount,
		ChatAddress:     msg.ChatAddress,
		ClusterName:     msg.ClusterName,
		FreezeAmount:    msg.FreezeAmount,
		ClusterDaoRatio: msg.ClusterDaoRatio,
	}

	clusterTrueId, err := k.Keeper.CreateCluster(ctx, msgCreate, msg.Metadata)
	if err != nil {
		return nil, err
	}

	
	if msg.BurnAmount.Add(msg.FreezeAmount).LT(sdk.NewDecFromInt(k.GetParams(ctx).ClusterLevels[1].BurnAmount)) {
		return nil, core.ErrCreateClusterBurn
	}

	
	_, err = k.BurnGetPower(ctx, fromAccAddr, fromAccAddr, clusterTrueId, msg.BurnAmount, msg.FreezeAmount, true)
	if err != nil {
		return nil, err
	}

	return &types.MsgEmptyResponse{}, nil
}

func (k msgServer) BurnToPower(goCtx context.Context, msg *types.MsgBurnToPower) (*types.MsgEmptyResponse, error) {
	logs := core.BuildLog(core.GetStructFuncName(k), core.LmChainDaoKeeper)

	ctx := sdk.UnwrapSDKContext(goCtx)

	
	if msg.UseFreezeAmount.IsPositive() {
		return nil, nil
	}

	
	if msg.GatewayAddress != "" && msg.ChatAddress != "" {
		err := k.ChatKeeper.Register(ctx, msg.FromAddress, msg.ChatAddress, msg.GatewayAddress, nil)
		if err != nil {
			return nil, err
		}
	}

	
	fromAccAddr, err := sdk.AccAddressFromBech32(msg.FromAddress)
	if err != nil {
		return nil, core.ErrAddressFormat
	}

	
	toAccAddr, err := sdk.AccAddressFromBech32(msg.ToAddress)
	if err != nil {
		return nil, core.ErrAddressFormat
	}

	
	clusterId, err := k.GetClusterId(ctx, msg.ClusterId)

	_, err = k.Keeper.BurnGetPower(ctx, fromAccAddr, toAccAddr, clusterId, msg.BurnAmount, msg.UseFreezeAmount, false)
	if err != nil {
		logs.WithError(err).Error("BurnToPower BurnGetPower Error")
		return nil, err
	}
	
	err = k.SettlementRouteBridgingReward(ctx, clusterId, msg.FromAddress, true)
	if err != nil {
		return nil, err
	}
	return nil, nil

}


func (k msgServer) ColonyRate(goCtx context.Context, msg *types.MsgColonyRate) (*types.MsgEmptyResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	
	gateway, err := k.gatewayKeeper.GetGatewayInfo(ctx, msg.GatewayAddress)
	if err != nil {
		return nil, err
	}
	
	if ctx.BlockTime().Unix()-gateway.MachineUpdateTime < core.OracleUpdateInterval {
		return nil, core.ErrMachineUpdateTime
	}
	
	if msg.Address != gateway.MachineAddress {
		return nil, core.ErrMachineAddress
	}
	for _, info := range msg.OnlineRate {
		err = k.UpdateClusterOnlineRate(ctx, info.Address, info.Rate)
		if err != nil {
			return nil, err
		}
	}
	
	gateway.MachineUpdateTime = ctx.BlockTime().Unix()
	err = k.gatewayKeeper.UpdateGatewayInfo(ctx, *gateway)
	if err != nil {
		return nil, err
	}

	

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Address),
		),
	)

	return &types.MsgEmptyResponse{}, nil
}
