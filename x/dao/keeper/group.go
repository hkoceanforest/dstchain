package keeper

import (
	"strconv"

	"freemasonry.cc/blockchain/core"
	"freemasonry.cc/blockchain/util"
	"freemasonry.cc/blockchain/x/dao/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/group"
)


func (k Keeper) UpdateGroupMembers(ctx sdk.Context, cluster types.DeviceCluster, members []group.MemberRequest) error {
	logs := core.BuildLog(core.GetStructFuncName(k), core.LmChainDaoKeeper)

	goCtx := sdk.WrapSDKContext(ctx)

	groupId := cluster.ClusterVoteId
	groupInfo, err := k.groupKeeper.GetGroupInfo(goCtx, groupId)

	if err != nil {
		logs.WithError(err).WithField("group", groupId).Error("AddGroupMembers GetGroupInfo Error")
		return err
	}

	groupAdmin := groupInfo.GetAdmin()

	req := &group.MsgUpdateGroupMembers{
		GroupId:       groupId,
		Admin:         groupAdmin,
		MemberUpdates: members,
	}
	_, err = k.groupKeeper.UpdateGroupMembers(ctx, req)
	if err != nil {
		logs.WithError(err).WithField("group", groupId).Error("AddGroupMembers UpdateGroupMembers Error")
		return err
	}
	updateGroupInfo, err := k.groupKeeper.GetGroupInfo(goCtx, groupId)

	if err != nil {
		logs.WithError(err).WithField("group", groupId).Error("AddGroupMembers GetGroupInfo Error")
		return err
	}

	adminMember, err := k.groupKeeper.GetGroupMember(ctx, &group.GroupMember{
		GroupId: groupId,
		Member:  &group.Member{Address: groupAdmin},
	})
	if err != nil {
		logs.WithError(err).WithField("group", groupId).Error("AddGroupMembers GetGroupMember Error")
		return err
	}
	totalWeight, err := sdk.NewDecFromStr(updateGroupInfo.TotalWeight)
	if err != nil {
		logs.WithError(err).WithField("group", groupId).Error("AddGroupMembers totalWeight Error")
		return err
	}
	adminWeight, err := sdk.NewDecFromStr(adminMember.Member.Weight)
	if err != nil {
		logs.WithError(err).WithField("group", groupId).Error("AddGroupMembers adminWeight Error")
		return err
	}
	newAdminWeight := totalWeight.Sub(adminWeight).Quo(sdk.NewDec(4))
	newAdminWeightStr := core.RemoveDecLastZero(newAdminWeight)
	newAdminWeightNum, _ := strconv.ParseInt(newAdminWeightStr, 10, 64)
	if newAdminWeightNum <= 0 {
		newAdminWeightStr = "20"
	}
	reqAdmin := &group.MsgUpdateGroupMembers{
		GroupId: groupId,
		Admin:   groupAdmin,
		MemberUpdates: []group.MemberRequest{
			{Address: groupAdmin, Weight: newAdminWeightStr},
		},
	}
	
	_, err = k.groupKeeper.UpdateGroupMembers(ctx, reqAdmin)
	if err != nil {
		logs.WithError(err).WithField("group", groupId).Error("AddGroupMembers UpdateGroupAdmin Error")
		return err
	}
	return nil
}


func (k Keeper) LeaveGroup(ctx sdk.Context, cluster types.DeviceCluster, member string) error {
	logs := core.BuildLog(core.GetStructFuncName(k), core.LmChainDaoKeeper)

	goCtx := sdk.WrapSDKContext(ctx)

	groupId := cluster.ClusterVoteId

	req := &group.MsgLeaveGroup{
		GroupId: groupId,
		Address: member,
	}
	_, err := k.groupKeeper.LeaveGroup(ctx, req)
	if err != nil {
		logs.WithError(err).WithField("group", groupId).Error("ClusterMemberExit LeaveGroup Error")
		return err
	}
	updateGroupInfo, err := k.groupKeeper.GetGroupInfo(goCtx, groupId)
	groupAdmin := updateGroupInfo.GetAdmin()
	if err != nil {
		logs.WithError(err).WithField("group", groupId).Error("AddGroupMembers GetGroupInfo Error")
		return err
	}

	adminMember, err := k.groupKeeper.GetGroupMember(ctx, &group.GroupMember{
		GroupId: groupId,
		Member:  &group.Member{Address: groupAdmin},
	})
	if err != nil {
		logs.WithError(err).WithField("group", groupId).Error("AddGroupMembers GetGroupMember Error")
		return err
	}
	totalWeight, err := sdk.NewDecFromStr(updateGroupInfo.TotalWeight)
	if err != nil {
		logs.WithError(err).WithField("group", groupId).Error("AddGroupMembers totalWeight Error")
		return err
	}
	adminWeight, err := sdk.NewDecFromStr(adminMember.Member.Weight)
	if err != nil {
		logs.WithError(err).WithField("group", groupId).Error("AddGroupMembers adminWeight Error")
		return err
	}
	newAdminWeight := totalWeight.Sub(adminWeight)
	newAdminWeightStr := core.RemoveDecLastZero(newAdminWeight)
	newAdminWeightNum, err := strconv.ParseInt(newAdminWeightStr, 10, 64)
	if newAdminWeightNum <= 0 {
		newAdminWeightStr = "20"
	}
	reqAdmin := &group.MsgUpdateGroupMembers{
		GroupId: groupId,
		Admin:   groupAdmin,
		MemberUpdates: []group.MemberRequest{
			group.MemberRequest{Address: groupAdmin, Weight: newAdminWeightStr},
		},
	}
	
	_, err = k.groupKeeper.UpdateGroupMembers(ctx, reqAdmin)
	if err != nil {
		logs.WithError(err).WithField("group", groupId).Error("AddGroupMembers UpdateGroupAdmin Error")
		return err
	}
	return nil
}


func (k Keeper) SetContractApproveInfo(ctx sdk.Context, contract string, approve types.ApprovePower) error {
	log := core.BuildLog(core.GetStructFuncName(k), core.LmChainDaoKeeper)
	store := ctx.KVStore(k.storeKey)
	key := types.GetClusterApprovePowerInfoKey(contract)
	data := make(map[string]types.ApprovePower)
	if store.Has(key) {
		resByte := store.Get(key)
		err := util.Json.Unmarshal(resByte, &data)
		if err != nil {
			log.WithError(err).WithField("contract:", contract).Error("SetContractApproveInfo Json.Unmarshal Error")
			return core.ErrGetApprovePowerInfo
		}
	}
	data[approve.ClusterId+approve.Address] = approve
	clusterCurApproveByte, err := util.Json.Marshal(data)
	if err != nil {
		log.WithError(err).Error("SetContractApproveInfo Marshal Err")
		return err
	}
	store.Set(key, clusterCurApproveByte)
	
	return k.removeContractApproveInfo(ctx, contract, approve)
}


func (k Keeper) removeContractApproveInfo(ctx sdk.Context, contract string, approve types.ApprovePower) error {
	
	clusterApprove, err := k.GetClusterApproveInfo(ctx, approve.ClusterId, approve.Address)
	if err != nil {
		return err
	}
	
	if clusterApprove.ApproveAddress == "" || clusterApprove.ApproveAddress == contract {
		return nil
	}
	
	res, err := k.GetContractApproveInfo(ctx, clusterApprove.ApproveAddress)
	if err != nil {
		return err
	}
	if res != nil {
		
		if _, ok := res[approve.ClusterId+approve.Address]; ok {
			delete(res, approve.ClusterId+approve.Address)
			store := ctx.KVStore(k.storeKey)
			key := types.GetClusterApprovePowerInfoKey(contract)
			clusterCurApproveByte, err := util.Json.Marshal(res)
			if err != nil {
				return err
			}
			store.Set(key, clusterCurApproveByte)
		}
	}
	return nil
}


func (k Keeper) GetContractApproveInfo(ctx sdk.Context, address string) (res map[string]types.ApprovePower, err error) {
	log := core.BuildLog(core.GetStructFuncName(k), core.LmChainDaoKeeper)

	store := ctx.KVStore(k.storeKey)

	key := types.GetClusterApprovePowerInfoKey(address)

	if store.Has(key) { 
		resByte := store.Get(key)
		err = util.Json.Unmarshal(resByte, &res)
		if err != nil {
			log.WithError(err).WithField("address:", address).Error("GetContractApproveInfo Json.Unmarshal Error")
			return res, core.ErrGetApprovePowerInfo
		}
		return res, nil
	}
	return nil, nil
}


func (k Keeper) SetClusterIdApproveInfo(ctx sdk.Context, clusterId, address string, approve types.ClusterCurApprove) error {
	log := core.BuildLog(core.GetStructFuncName(k), core.LmChainDaoKeeper)
	store := ctx.KVStore(k.storeKey)

	key := types.GetClusterApproveInfoInfoKey(clusterId + address)

	clusterCurApproveByte, err := util.Json.Marshal(approve)
	if err != nil {
		log.WithError(err).Error("ClusterPowerApprove ApprovePower struct Err")
		return err
	}
	store.Set(key, clusterCurApproveByte)
	return nil
}


func (k Keeper) GetClusterApproveInfo(ctx sdk.Context, clusterId, address string) (res types.ClusterCurApprove, err error) {
	log := core.BuildLog(core.GetStructFuncName(k), core.LmChainDaoKeeper)

	store := ctx.KVStore(k.storeKey)

	key := types.GetClusterApproveInfoInfoKey(clusterId + address)

	if store.Has(key) { 
		resByte := store.Get(key)
		err = util.Json.Unmarshal(resByte, &res)
		if err != nil {
			log.WithError(err).WithField("clusterId:", clusterId).Error("GetClusterIdApproveInfo Json.Unmarshal Error")
			return res, core.ErrGetApprovePowerInfo
		}
		return res, nil
	} else { 
		return types.ClusterCurApprove{
			ApproveAddress: "",
			EndBlock:       0,
		}, nil
	}
}
