package keeper

import (
	"freemasonry.cc/blockchain/core"
	"freemasonry.cc/blockchain/util"
	"freemasonry.cc/blockchain/x/dao/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/group"
)


func (k Keeper) CalculateFee(ctx sdk.Context, msg sdk.Msg, fee sdk.Coins) (sdk.Coins, error) {
	
	clusterId, addr := k.deductionTxMsg(ctx, msg)
	if clusterId == "" {
		return fee, nil
	}
	feeCoin := sdk.Coin{}
	for _, coin := range fee {
		if coin.Denom == core.BaseDenom {
			feeCoin = coin
		}
	}
	cluster, err := k.GetClusterByChatId(ctx, clusterId)
	if err != nil {
		return fee, nil
	}
	cycleFee, err := k.getCycleFee(ctx, clusterId)
	if err != nil {
		return nil, err
	}
	
	powerDec, err := k.CalculateBurnGetPower(ctx, sdk.OneDec())
	if err != nil {
		return nil, err
	}
	
	limit := cluster.ClusterPowerMembers[cluster.ClusterOwner].ActivePower.Quo(k.GetParams(ctx).PowerGasRatio).Quo(powerDec).TruncateInt()
	
	balance := limit.Sub(cycleFee.Amount)
	if balance.IsNegative() || balance.IsZero() {
		return fee, nil
	}
	
	if balance.LT(feeCoin.Amount) {
		
		err = k.setCycleFee(ctx, clusterId, sdk.NewCoin(core.BaseDenom, balance))
		if err != nil {
			return nil, err
		}
		
		payFee := sdk.NewCoin(feeCoin.Denom, feeCoin.Amount.Sub(balance))
		events := sdk.Events{
			sdk.NewEvent(
				types.EventTypeDeductionFee,
				sdk.NewAttribute(sdk.AttributeKeyFee, feeCoin.String()),
				sdk.NewAttribute(types.AttributeDeductionFee, sdk.NewCoin(core.BaseDenom, balance).String()),
				sdk.NewAttribute(sdk.AttributeKeyFeePayer, addr),
				sdk.NewAttribute(types.AttributeClusterId, clusterId),
			),
		}
		ctx.EventManager().EmitEvents(events)
		
		return sdk.NewCoins(payFee), nil
	}
	err = k.setCycleFee(ctx, clusterId, feeCoin)
	if err != nil {
		return nil, err
	}
	events := sdk.Events{
		sdk.NewEvent(
			types.EventTypeDeductionFee,
			sdk.NewAttribute(sdk.AttributeKeyFee, feeCoin.String()),
			sdk.NewAttribute(types.AttributeDeductionFee, feeCoin.String()),
			sdk.NewAttribute(sdk.AttributeKeyFeePayer, addr),
			sdk.NewAttribute(types.AttributeClusterId, clusterId),
		),
	}
	ctx.EventManager().EmitEvents(events)
	return nil, nil
}

func (k Keeper) setCycleFee(ctx sdk.Context, clusterId string, fee sdk.Coin) error {
	store := ctx.KVStore(k.storeKey)
	data := make(map[int64]sdk.Coin)
	key := types.GetClusterDeductionFeeKey(clusterId)
	if store.Has(key) {
		bz := store.Get(key)
		err := util.Json.Unmarshal(bz, &data)
		if err != nil {
			return err
		}
	}
	cycle := ctx.BlockHeight() / core.DayBlockNum
	if _, ok := data[cycle]; ok {
		data[cycle] = data[cycle].Add(fee)
	} else {
		data[cycle] = fee
	}
	dataByte, err := util.Json.Marshal(data)
	if err != nil {
		return err
	}
	store.Set(key, dataByte)
	return nil
}

func (k Keeper) getCycleFee(ctx sdk.Context, clusterId string) (sdk.Coin, error) {
	store := ctx.KVStore(k.storeKey)
	zeroCoin := sdk.NewCoin(core.BaseDenom, sdk.ZeroInt())
	data := make(map[int64]sdk.Coin)
	key := types.GetClusterDeductionFeeKey(clusterId)
	cycle := ctx.BlockHeight() / core.DayBlockNum
	if store.Has(key) {
		bz := store.Get(key)
		err := util.Json.Unmarshal(bz, &data)
		if err != nil {
			return zeroCoin, err
		}
		if _, ok := data[cycle]; ok {
			return data[cycle], nil
		}
		return zeroCoin, nil
	}
	return zeroCoin, nil
}


func (k Keeper) deductionTxMsg(ctx sdk.Context, msg sdk.Msg) (string, string) {
	if obj, oks := msg.(*types.MsgClusterAddMembers); oks {
		return obj.ClusterId, obj.FromAddress
	}
	if obj, oks := msg.(*types.MsgDeleteMembers); oks {
		return obj.ClusterId, obj.FromAddress
	}
	if obj, oks := msg.(*types.MsgClusterChangeName); oks {
		return obj.ClusterId, obj.FromAddress
	}
	if obj, oks := msg.(*types.MsgClusterMemberExit); oks {
		return obj.ClusterId, obj.FromAddress
	}
	if obj, oks := msg.(*types.MsgBurnToPower); oks {
		return obj.ClusterId, obj.FromAddress
	}
	if obj, oks := msg.(*types.MsgClusterChangeSalaryRatio); oks {
		return obj.ClusterId, obj.FromAddress
	}
	if obj, oks := msg.(*types.MsgClusterChangeDvmRatio); oks {
		return obj.ClusterId, obj.FromAddress
	}
	if obj, oks := msg.(*types.MsgClusterChangeDaoRatio); oks {
		return obj.ClusterId, obj.FromAddress
	}
	if obj, oks := msg.(*types.MsgClusterChangeId); oks {
		return obj.ClusterId, obj.FromAddress
	}
	if obj, oks := msg.(*types.MsgWithdrawSwapDpos); oks {
		return obj.ClusterId, obj.MemberAddress
	}
	if obj, oks := msg.(*types.MsgWithdrawDeviceReward); oks {
		return obj.ClusterId, obj.MemberAddress
	}
	if obj, oks := msg.(*types.MsgWithdrawOwnerReward); oks {
		return obj.ClusterId, obj.Address
	}
	if obj, oks := msg.(*types.MsgAgreeJoinCluster); oks {
		return obj.ClusterId, obj.FromAddress
	}
	if obj, oks := msg.(*group.MsgSubmitProposal); oks {
		clusterId, err := k.GetClusterVotePolicy(ctx, obj.GroupPolicyAddress)
		if err != nil {
			return "", ""
		}
		return clusterId, obj.Proposers[0]
	}
	if obj, oks := msg.(*group.MsgVote); oks {
		proposal, err := k.groupKeeper.Proposal(ctx, &group.QueryProposalRequest{ProposalId: obj.ProposalId})
		if err != nil {
			return "", ""
		}
		clusterId, err := k.GetClusterVotePolicy(ctx, proposal.Proposal.GroupPolicyAddress)
		if err != nil {
			return "", ""
		}
		return clusterId, obj.Voter
	}

	
	if obj, oks := msg.(*types.MsgRedPacket); oks {
		return obj.Clusterid, obj.Fromaddress
	}

	
	if obj, oks := msg.(*types.MsgOpenRedPacket); oks {
		
		rInfo, err := k.GetRedPacketInfo(ctx, obj.Redpacketid)
		if err != nil {
			return "", ""
		}

		cluster, err := k.GetCluster(ctx, rInfo.ClusterTrueId)
		if err != nil {
			return "", ""
		}

		return cluster.ClusterChatId, obj.Fromaddress
	}

	
	if obj, oks := msg.(*types.MsgReturnRedPacket); oks {

		
		rInfo, err := k.GetRedPacketInfo(ctx, obj.Redpacketid)
		if err != nil {
			return "", ""
		}

		cluster, err := k.GetCluster(ctx, rInfo.ClusterTrueId)
		if err != nil {
			return "", ""
		}

		return cluster.ClusterChatId, obj.Fromaddress
	}
	return "", ""
}

func (k Keeper) FreeTxMsg(ctx sdk.Context, msg sdk.Msg) bool {
	if obj, oks := msg.(*types.MsgColonyRate); oks {
		
		gateway, err := k.gatewayKeeper.GetGatewayInfo(ctx, obj.GetGatewayAddress())
		if err != nil {
			return true
		}
		if gateway.MachineAddress != obj.Address {
			return true
		}
		return false
	}

	
	
	
	return true
}
