package keeper

import (
	sdkmath "cosmossdk.io/math"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"freemasonry.cc/blockchain/core"
	"freemasonry.cc/blockchain/util"
	"freemasonry.cc/blockchain/x/dao/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authTypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/group"
	"github.com/ethereum/go-ethereum/common"
	"strconv"
	"strings"
)

func (k Keeper) ValidateSalaryRatio(params types.Params, salaryRatio sdk.Dec) error {
	if salaryRatio.GT(params.SalaryRewardRatio.MaxRatio) || salaryRatio.LT(params.SalaryRewardRatio.MinRatio) {
		return core.ErrSalaryRatio
	}
	return nil
}
func (k Keeper) ValidateDvmRatio(params types.Params, dvmRatio sdk.Dec) error {
	if dvmRatio.GT(params.DvmRewardRatio.MaxRatio) || dvmRatio.LT(params.DvmRewardRatio.MinRatio) {
		return core.ErrDvmRatio
	}
	return nil
}

func (k Keeper) ValidateDaoRatio(params types.Params, daoRatio sdk.Dec) error {
	if daoRatio.GT(params.DaoRewardRatio.MaxRatio) || daoRatio.LT(params.DaoRewardRatio.MinRatio) {
		return core.ErrDaoRatio
	}
	return nil
}


func (k Keeper) GetAllClusters(ctx sdk.Context) (clusters []types.DeviceCluster) {
	k.IterateAllCluster(ctx, func(cluster types.DeviceCluster) bool {
		clusters = append(clusters, cluster)
		return false
	})

	return clusters
}


func (k Keeper) GetAllPersonClusters(ctx sdk.Context) (clusters []types.PersonalClusterInfo) {
	k.IterateAllPersonClusters(ctx, func(cluster types.PersonalClusterInfo) bool {
		clusters = append(clusters, cluster)
		return false
	})

	return clusters
}

func (k Keeper) GetAllClusterAddress(ctx sdk.Context) (addressMap []types.ClusterChatId2ClusterId) {
	addressMap = make([]types.ClusterChatId2ClusterId, 0)
	k.IterateAllClusterAddress(ctx, func(cluster types.ClusterChatId2ClusterId) bool {
		addressMap = append(addressMap, cluster)
		return false
	})

	return addressMap
}

func (k Keeper) IterateAllClusterAddress(ctx sdk.Context, cb func(clusterChatId2ClusterId types.ClusterChatId2ClusterId) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.ClusterIdKey)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		r := types.ClusterChatId2ClusterId{
			ClusterChatId: k.GetClusterIdFromKey(string(iterator.Key())),
			ClusterId:     string(iterator.Value()),
		}
		if cb(r) {
			break
		}
	}
}


func (k Keeper) GetAllGatewayClusters(ctx sdk.Context) (clusters []types.GatewayClusterExport) {
	k.IterateAllGatewayClusters(ctx, func(cluster types.GatewayClusterExport) bool {
		clusters = append(clusters, cluster)
		return false
	})

	return clusters
}


func (k Keeper) IterateAllGatewayClusters(ctx sdk.Context, cb func(cluster types.GatewayClusterExport) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.ClusterForGateway)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		gatewayClusters := make(map[string]struct{})
		err := util.Json.Unmarshal(iterator.Value(), &gatewayClusters)
		if err != nil {
			panic(err)
		}

		r := types.GatewayClusterExport{
			GatewayAddress: k.GetGatewayAddressFromKey(string(iterator.Key())),
			Clusters:       gatewayClusters,
		}

		if cb(r) {
			break
		}
	}
}


func (k Keeper) GetAllClusterVotePolicy(ctx sdk.Context) (addressMap []types.ClusterStrategyAddress) {
	addressMap = make([]types.ClusterStrategyAddress, 0)
	k.IterateAllClusterVotePolicy(ctx, func(cluster types.ClusterStrategyAddress) bool {
		addressMap = append(addressMap, cluster)
		return false
	})

	return addressMap
}


func (k Keeper) IterateAllClusterVotePolicy(ctx sdk.Context, cb func(cluster types.ClusterStrategyAddress) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.ClusterPolicyPrefix)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		r := types.ClusterStrategyAddress{
			ClusterId:       string(iterator.Value()),
			StrategyAddress: k.GetClusterIdFromPolicyKey(string(iterator.Key())),
		}
		if cb(r) {
			break
		}
	}
}


func (k Keeper) GetAllClusterCreateTime(ctx sdk.Context) []types.ClusterCreateTime {
	createTimes := make([]types.ClusterCreateTime, 0)
	k.IterateAllClusterCreateTime(ctx, func(cluster types.ClusterCreateTime) bool {
		createTimes = append(createTimes, cluster)
		return false
	})
	return createTimes
}

func (k Keeper) GetClusterIdFromPolicyKey(key string) string {
	return key[len(types.ClusterPolicyPrefix):]
}

func (k Keeper) GetGatewayAddressFromKey(key string) string {
	return key[len(types.ClusterForGateway):]
}

func (k Keeper) GetClusterChatIdFromKey(key string) string {
	return key[len(types.ClusterIdKey):]
}

func (k Keeper) GetClusterIdFromKey(key string) string {
	return key[len(types.ClusterIdKey):]
}


func (k Keeper) IterateAllPersonClusters(ctx sdk.Context, cb func(cluster types.PersonalClusterInfo) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.PersonClusterInfoKey)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		cluster := new(types.PersonalClusterInfo)
		err := util.Json.Unmarshal(iterator.Value(), cluster)
		if err != nil {
			panic(err)
		}
		if cb(*cluster) {
			break
		}
	}
}

func (k Keeper) IterateAllCluster(ctx sdk.Context, cb func(cluster types.DeviceCluster) (stop bool)) {
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, types.DeviceClusterKey)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		cluster := new(types.DeviceCluster)
		err := util.Json.Unmarshal(iterator.Value(), cluster)
		if err != nil {
			panic(err)
		}
		if cb(*cluster) {
			break
		}
	}
}

func (k Keeper) CreateCluster(ctx sdk.Context, msg *types.MsgCreateClusterAddMembers, metadata string) (string, error) {
	logs := core.BuildLog(core.GetStructFuncName(k), core.LmChainDaoKeeper)

	
	p, err := k.GetPersonClusterInfo(ctx, msg.FromAddress)
	if err != nil {
		return "", err
	}

	
	_, err = k.gatewayKeeper.GetGatewayInfo(ctx, msg.GateAddress)
	if err != nil {
		return "", err
	}

	goCtx := sdk.WrapSDKContext(ctx)

	
	groupMembers := []group.MemberRequest{
		{
			Address:  msg.FromAddress,
			Weight:   "20",
			Metadata: metadata,
		},
	}
	groupMsg := group.MsgCreateGroup{
		Admin: msg.FromAddress, Metadata: metadata, Members: groupMembers,
	}

	
	groupResp, err := k.groupKeeper.CreateGroup(goCtx, &groupMsg)
	if err != nil {
		logs.WithError(err).Error("CreateCluster CreateGroup Error")
		return "", err
	}
	params := k.GetParams(ctx)
	
	policy := group.NewPercentageDecisionPolicy(
		"0.5",
		params.GetVotingPeriod(), 
		0,
	)
	policyReq := &group.MsgCreateGroupPolicy{
		Admin:   msg.FromAddress,
		GroupId: groupResp.GroupId,
	}
	err = policyReq.SetDecisionPolicy(policy)
	if err != nil {
		logs.WithError(err).Error("CreateCluster SetDecisionPolicy Error")
		return "", err
	}
	policyRes, err := k.groupKeeper.CreateGroupPolicy(ctx, policyReq)
	if err != nil {
		logs.WithError(err).Error("CreateCluster CreateGroupPolicy Error")
		return "", err
	}

	
	clusterId := util.Md5String(msg.ClusterId)
	cluster := types.DeviceCluster{
		ClusterId:      clusterId,
		ClusterChatId:  msg.ClusterId,
		ClusterName:    msg.ClusterName,
		ClusterOwner:   msg.FromAddress,
		ClusterGateway: msg.GateAddress,
		ClusterLeader:  p.FirstPowerCluster,
		ClusterDeviceMembers: map[string]types.ClusterDeviceMember{ 
			msg.FromAddress: {
				Address:     msg.FromAddress,
				ActivePower: sdk.NewDec(1), 
			},
		},
		ClusterPowerMembers:    make(map[string]types.ClusterPowerMember),
		ClusterPower:           sdk.ZeroDec(),
		ClusterLevel:           1,             
		ClusterBurnAmount:      sdk.ZeroDec(), 
		ClusterActiveDevice:    1,             
		ClusterDaoPool:         authTypes.NewModuleAddress(clusterId).String(),
		ClusterRouteRewardPool: authTypes.NewModuleAddress(clusterId + "connectedPool").String(),
		ClusterDeviceRatio:     core.ClusterDeviceRate,            
		ClusterSalaryRatio:     params.SalaryRewardRatio.MaxRatio, 
		ClusterDvmRatio:        params.DvmRewardRatio.MaxRatio,    
		ClusterDaoRatio:        sdk.OneDec(),                      
		OnlineRatio:            sdk.MustNewDecFromStr("1"),
		OnlineRatioUpdateTime:  ctx.BlockTime().Unix(),
		ClusterAdminList:       make(map[string]struct{}),
		ClusterVoteId:          groupResp.GroupId,
		ClusterVotePolicy:      policyRes.Address,
		ClusterSalaryRatioUpdateHeight: types.ClusterChangeRatioHeight{
			SalaryRatioUpdateHeight: ctx.BlockHeight(),
			DvmRatioUpdateHeight:    ctx.BlockHeight(),
			DaoRatioUpdateHeight:    ctx.BlockHeight(),
		},
	}

	
	err = k.SetDeviceCluster(ctx, cluster)
	if err != nil {
		return "", err
	}

	
	newDivice := p.Device
	newDivice[clusterId] = struct{}{}

	newOwners := p.Owner
	newOwners[clusterId] = struct{}{}

	p.Device = newDivice 
	p.Owner = newOwners  

	
	err = k.SetPersonClusterInfo(ctx, p)
	if err != nil {
		return "", err
	}

	
	k.SetClusterId(ctx, clusterId, msg.ClusterId)

	
	k.SetClusterCreateTime(ctx, clusterId)

	
	k.SetClusterVotePolicy(ctx, clusterId, policyRes.Address)

	
	err = k.AddClusterToGateway(ctx, msg.GateAddress, clusterId)
	if err != nil {
		return "", core.ErrSetDeviceCluster
	}

	
	k.initializeCluster(ctx, cluster)

	
	k.initializeDevice(ctx, cluster)

	fromAccAddr, err := sdk.AccAddressFromBech32(msg.FromAddress)
	if err != nil {
		logs.WithError(err).WithField("from:", msg.FromAddress).Error("CreateCluster AccAddressFromBech32 Error")
		return "", core.ErrAddressFormat
	}
	fromBalance := k.BankKeeper.GetAllBalances(ctx, fromAccAddr)

	
	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeCreateCluster,
		sdk.NewAttribute(types.AttributeKeyClusterChatId, msg.ClusterId),      
		sdk.NewAttribute(types.AttributeKeyClusterOwner, msg.FromAddress),     
		sdk.NewAttribute(types.AttributeSenderBalances, fromBalance.String()), 
	))

	
	k.initializeDeviceDelegation(ctx, cluster, msg.FromAddress, 1)

	return clusterId, nil
}

func (k Keeper) GetClusterId(ctx sdk.Context, clusterChatId string) (string, error) {
	logs := core.BuildLog(core.GetStructFuncName(k), core.LmChainDaoKeeper)
	store := ctx.KVStore(k.storeKey)

	key := types.GetClusterIdKey(clusterChatId)

	if store.Has(key) {
		clusterIdByte := store.Get(key)
		if clusterIdByte == nil {
			logs.WithField("clusterChatId:", clusterChatId).Error("GetClusterId value nil")
			return "", core.ErrGetClusterId
		}

		return string(clusterIdByte), nil
	} else {
		logs.WithField("clusterChatId:", clusterChatId).Warning("GetClusterId not found")
		return "", core.ErrClusterIdNotFound
	}
}

func (k Keeper) SetClusterId(ctx sdk.Context, clusterId, clusterChatId string) {
	store := ctx.KVStore(k.storeKey)

	key := types.GetClusterIdKey(clusterChatId)

	store.Set(key, []byte(clusterId))
}

func (k Keeper) changeClusterChatId(ctx sdk.Context, oldChatId, newChatId, clusterId string) {
	store := ctx.KVStore(k.storeKey)

	k.SetClusterId(ctx, clusterId, newChatId)

	oldKey := types.GetClusterIdKey(oldChatId)
	store.Delete(oldKey)
}

func (k Keeper) GetPersonClusterInfo(ctx sdk.Context, address string) (res types.PersonalClusterInfo, err error) {
	log := core.BuildLog(core.GetStructFuncName(k), core.LmChainDaoKeeper)

	store := ctx.KVStore(k.storeKey)

	key := types.GetPersonClusterInfoKey(address)

	if store.Has(key) { 
		resByte := store.Get(key)
		err = util.Json.Unmarshal(resByte, &res)
		if err != nil {
			log.WithError(err).WithField("address:", address).Error("GetPersonClusterInfo Json.Unmarshal Error")
			return res, core.ErrGetPersonClusterInfo
		}
		return res, nil
	} else { 
		zeroDec := sdk.ZeroDec()
		return types.PersonalClusterInfo{
			Address:           address,
			Device:            make(map[string]struct{}),
			Owner:             make(map[string]struct{}),
			BePower:           make(map[string]struct{}),
			AllBurn:           zeroDec,
			ActivePower:       zeroDec,
			FreezePower:       zeroDec,
			FirstPowerCluster: "",
		}, nil
	}
}

func (k Keeper) GetClusterVotePolicy(ctx sdk.Context, policyAddress string) (string, error) {
	logs := core.BuildLog(core.GetStructFuncName(k), core.LmChainDaoKeeper)
	store := ctx.KVStore(k.storeKey)

	key := types.GetClusterPolicyKey(policyAddress)

	if store.Has(key) {
		clusterIdByte := store.Get(key)
		if clusterIdByte == nil {
			logs.WithField("policyAddress:", policyAddress).Error("GetClusterId value nil")
			return "", core.ErrGetClusterId
		}
		return string(clusterIdByte), nil
	} else {
		logs.WithField("policyAddress:", policyAddress).Error("GetClusterId not found")
		return "", core.ErrClusterIdNotFound
	}
}

func (k Keeper) SetClusterVotePolicy(ctx sdk.Context, clusterId, policyAddress string) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetClusterPolicyKey(policyAddress), []byte(clusterId))
}

func (k Keeper) SetDeviceCluster(ctx sdk.Context, cluster types.DeviceCluster) error {
	log := core.BuildLog(core.GetStructFuncName(k), core.LmChainDaoKeeper)
	store := ctx.KVStore(k.storeKey)
	clusterByte, err := util.Json.Marshal(cluster)
	if err != nil {
		log.WithError(err).WithField("cluster struct", cluster).Error("StorageDeviceCluster json marshal error")
		return core.ErrSetCluster
	}

	store.Set(types.GetDeviceClusterKey(cluster.ClusterId), clusterByte)

	return nil
}

func (k Keeper) SetPersonClusterInfo(ctx sdk.Context, p types.PersonalClusterInfo) error {
	log := core.BuildLog(core.GetStructFuncName(k), core.LmChainDaoKeeper)

	store := ctx.KVStore(k.storeKey)

	_, err := sdk.AccAddressFromBech32(p.Address)

	if err != nil {
		log.WithField("address: ", p.Address).WithError(err).Error(" SetPersonClusterInfo p.Address Error")
		return core.ErrAddressFormat
	}

	pByte, err := util.Json.Marshal(p)
	if err != nil {
		log.WithError(err).WithField("p struct", p).Error("SetPersonClusterInfo json marshal error")
		return core.ErrSetPersonClusterInfo
	}

	store.Set(types.GetPersonClusterInfoKey(p.Address), pByte)

	return nil
}

func (k Keeper) GetCluster(ctx sdk.Context, cluserId string) (res types.DeviceCluster, err error) {
	log := core.BuildLog(core.GetStructFuncName(k), core.LmChainDaoKeeper)

	store := ctx.KVStore(k.storeKey)

	key := types.GetDeviceClusterKey(cluserId)

	if store.Has(key) { 
		resByte := store.Get(key)
		err = util.Json.Unmarshal(resByte, &res)
		if err != nil {
			log.WithError(err).WithField("cluserId:", cluserId).Error("GetCluster Json.Unmarshal Error")
			return res, core.ErrGetCluster
		}
		return res, nil
	} else { 
		return res, core.GetClusterNotFound
	}
}

func (k Keeper) GetClusterByChatId(ctx sdk.Context, cluserChatId string) (res types.DeviceCluster, err error) {
	logs := core.BuildLog(core.GetStructFuncName(k), core.LmChainDaoKeeper)

	cluserId, err := k.GetClusterId(ctx, cluserChatId)
	if err != nil {
		logs.WithError(err).Warning("GetClusterByChatId GetClusterId Error")
		return types.DeviceCluster{}, err
	}

	store := ctx.KVStore(k.storeKey)

	key := types.GetDeviceClusterKey(cluserId)

	if store.Has(key) { 
		resByte := store.Get(key)
		err = util.Json.Unmarshal(resByte, &res)
		if err != nil {
			logs.WithError(err).WithField("cluserId:", cluserId).Error("GetCluster Json.Unmarshal Error")
			return res, core.ErrGetCluster
		}
		return res, nil
	} else { 
		return res, core.GetClusterNotFound
	}
}

func (k Keeper) CalculateClusterUpgrade(ctx sdk.Context, cluster types.DeviceCluster, addBurnAmount sdk.Dec, addActiveMembers int64) int64 {
	
	newBurn := cluster.ClusterBurnAmount.Add(addBurnAmount)
	
	
	
	levels := k.GetParams(ctx).ClusterLevels
	
	for i := len(levels) - 1; i >= 0; i-- {
		
		
		if newBurn.GTE(sdk.NewDecFromInt(levels[i].BurnAmount)) {
			if levels[i].Level > cluster.ClusterLevel {
				return levels[i].Level
			}
		}
	}
	return cluster.ClusterLevel
}

func (k Keeper) AddClusterToGateway(ctx sdk.Context, gateway, clusterId string) error {
	logs := core.BuildLog(core.GetStructFuncName(k), core.LmChainDaoKeeper)
	store := ctx.KVStore(k.storeKey)

	key := types.GetClusterForGatewayKey(gateway)

	clustersForGateway := make(map[string]struct{})
	
	if !store.Has(key) {
		clustersForGateway[clusterId] = struct{}{}
	} else {
		oldByte := store.Get(key)
		err := util.Json.Unmarshal(oldByte, &clustersForGateway)
		if err != nil {
			logs.WithError(err).Error("AddClusterToGateway clustersForGateway Unmarshal Error")
			return core.ErrAddClusterToGateway
		}

		clustersForGateway[clusterId] = struct{}{}
	}

	clustersForGatewayByte, err := util.Json.Marshal(clustersForGateway)
	if err != nil {
		logs.WithError(err).Error("AddClusterToGateway clustersForGateway Marshal Error")
		return core.ErrAddClusterToGateway
	}

	store.Set(key, clustersForGatewayByte)

	return nil
}

func (k Keeper) GateClusterToGateway(ctx sdk.Context, gateway string) (map[string]map[string]types.ClusterDeviceMember, error) {
	logs := core.BuildLog(core.GetStructFuncName(k), core.LmChainDaoKeeper)
	store := ctx.KVStore(k.storeKey)

	key := types.GetClusterForGatewayKey(gateway)

	if !store.Has(key) {
		return map[string]map[string]types.ClusterDeviceMember{}, core.EmptyGetClustertoGate
	}

	
	res := make(map[string]map[string]types.ClusterDeviceMember)
	clusters := make(map[string]struct{})

	clustersByte := store.Get(key)

	err := util.Json.Unmarshal(clustersByte, &clusters)
	if err != nil {
		logs.WithError(err).Error("get --> json.Unmarshal")
		return map[string]map[string]types.ClusterDeviceMember{}, core.ErrGetClusterToGateway
	}

	
	for clusterId := range clusters {
		cluster, err := k.GetCluster(ctx, clusterId)
		if err != nil {
			logs.WithError(err).WithField("cluster_id", clusterId).Error("GateClusterToGateway GetCluster Error")
			return map[string]map[string]types.ClusterDeviceMember{}, err
		}

		res[clusterId] = cluster.ClusterDeviceMembers
	}

	return res, nil
}

func (k Keeper) ValidateClusterAdmin(ctx sdk.Context, cluster types.DeviceCluster, list []string) error {
	logs := core.BuildLog(core.GetStructFuncName(k), core.LmChainDaoKeeper)

	
	if int64(len(list)) > k.GetParams(ctx).MaxClusterMembers {
		return core.ErrMaxClusterMember
	}

	for _, addr := range list {
		
		if _, err := sdk.AccAddressFromBech32(addr); err != nil {
			logs.WithError(err).WithField("addr", addr).Error("ValidateClusterAdmin AccAddressFromBech32 Error")
			return core.ErrAddressFormat
		}

		
		if _, ok := cluster.ClusterDeviceMembers[addr]; !ok {
			logs.WithField("addr", addr).Error("ValidateClusterAdmin ClusterDeviceMembers Error")
			return core.ErrNotInCluster
		}

	}
	return nil
}

func (k Keeper) IsOwnerOrAdmin(cluster types.DeviceCluster, addr string) (res bool) {

	if addr == cluster.ClusterOwner {
		return true
	}

	if k.IsAdmin(cluster, addr) {
		return true
	}

	return false
}

func (k Keeper) IsAdmin(cluster types.DeviceCluster, addr string) (res bool) {

	if len(cluster.ClusterAdminList) > 0 {
		if _, ok := cluster.ClusterAdminList[addr]; ok {
			return true
		}
	}

	return false
}

func (k Keeper) GetClusterConnectivityRate(ctx sdk.Context, clusterId string) sdk.Dec {

	var cluster types.DeviceCluster
	var err error
	var res sdk.Dec
	cluster, err = k.GetCluster(ctx, clusterId)
	if err != nil {
		return sdk.ZeroDec()
	}

	res = sdk.MustNewDecFromStr("0.25")

	for i := 0; i < 2; i++ {
		newCluster, err := k.GetCluster(ctx, cluster.ClusterLeader)
		if err != nil {
			return res
		}

		cluster = newCluster
		res = res.Add(sdk.MustNewDecFromStr("0.5"))
	}

	return res
}

func (k Keeper) GetLevelInfoByLevel(params types.Params, level int64) (types.ClusterLevel, error) {
	levels := params.ClusterLevels

	for _, clusterLevel := range levels {
		if clusterLevel.Level == level {
			return clusterLevel, nil
		}
	}

	return types.ClusterLevel{}, core.ErrLevelInfo
}

func (k Keeper) GetDeviceNoEvm(c types.DeviceCluster) int64 {

	n := int64(0)

	for _, m := range c.ClusterDeviceMembers {
		if _, ok := c.ClusterPowerMembers[m.Address]; !ok {
			n = n + 1
		}
	}

	return n
}

func (k Keeper) GetNoDvm(ctx sdk.Context, clusterId string, addrs []string) ([]string, error) {
	res := make([]string, 0)
	cluster, err := k.GetClusterByChatId(ctx, clusterId)
	if err != nil {
		return res, err
	}

	for _, addr := range addrs {
		if _, ok := cluster.ClusterPowerMembers[addr]; ok {
			continue
		}

		res = append(res, addr)
	}

	return res, nil
}

func (k Keeper) DvmApprove(ctx sdk.Context, approveAddress, msgClusterId, approveEndBlock, fromAddress string) error {
	logs := core.BuildLog(core.GetStructFuncName(k), core.LmChainDaoKeeper)
	var err error
	
	if !k.contractKeeper.IsContract(ctx, approveAddress) {
		return core.ErrContractAddress
	}
	clusterId := msgClusterId
	if strings.Contains(msgClusterId, ".") {
		
		clusterId, err = k.GetClusterId(ctx, msgClusterId)
		if err != nil {
			return err
		}
	}
	cluster, err := k.GetCluster(ctx, clusterId)
	if err != nil {
		logs.WithError(err).Error("PersonDvmApprove GetCluster Err")
		return err
	}
	
	if _, ok := cluster.ClusterPowerMembers[fromAddress]; !ok {
		return core.ErrMemberNotExist
	}
	
	if fromAddress == cluster.ClusterOwner {
		return core.ErrAuthorization
	}

	curApprove, err := k.GetClusterApproveInfo(ctx, cluster.ClusterId, fromAddress)
	if err != nil {
		logs.WithError(err).Error("PersonDvmApprove GetClusterIdApproveInfo Err")
		return err
	}
	if curApprove.EndBlock > ctx.BlockHeight() {
		return core.ErrApproveNotEnd
	}
	blockNum, err := strconv.ParseInt(approveEndBlock, 10, 64)
	if err != nil {
		logs.WithError(err).Error("PersonDvmApprove ApproveEndBlock Err")
		return err
	}
	endBlock := ctx.BlockHeight() + blockNum
	approvePower := types.ApprovePower{
		ClusterId: clusterId,
		Address:   fromAddress,
		IsDaoPool: false,
		EndBlock:  endBlock,
	}
	err = k.SetContractApproveInfo(ctx, approveAddress, approvePower)
	if err != nil {
		return err
	}
	approve := types.ClusterCurApprove{
		ApproveAddress: strings.ToLower(approveAddress),
		EndBlock:       endBlock,
	}
	err = k.SetClusterIdApproveInfo(ctx, clusterId, fromAddress, approve)
	if err != nil {
		return err
	}
	return nil
}

func (k Keeper) SetClusterDaoRewardQueue(ctx sdk.Context, address, clusterId string, burnAmount sdk.Dec, params types.Params) error {
	store := ctx.KVStore(k.storeKey)
	key := types.GetClusterMemberDaoReward(clusterId)
	var queue []types.ClusterMemberDaoReward
	if burnAmount.IsZero() {
		return nil
	}
	if store.Has(key) {
		bz := store.Get(key)
		err := util.Json.Unmarshal(bz, &queue)
		if err != nil {
			return err
		}
	}
	
	daoReward := burnAmount.Mul(params.ReceiveDaoRatio).TruncateInt()
	if daoReward.IsZero() {
		return nil
	}
	member := types.ClusterMemberDaoReward{Address: address, DaoReward: daoReward, BurnAmount: burnAmount.TruncateInt(), Time: ctx.BlockTime().Unix()}
	queue = append(queue, member)
	bz, err := util.Json.Marshal(queue)
	if err != nil {
		return err
	}
	store.Set(key, bz)
	return nil
}

func (k Keeper) UpdateClusterDaoRewardQueue(ctx sdk.Context, queue []types.ClusterMemberDaoReward, clusterId string) error {
	store := ctx.KVStore(k.storeKey)
	key := types.GetClusterMemberDaoReward(clusterId)
	bz, err := util.Json.Marshal(queue)
	if err != nil {
		return err
	}
	store.Set(key, bz)
	return nil
}

func (k Keeper) GetClusterDaoRewardQueue(ctx sdk.Context, clusterId string) ([]types.ClusterMemberDaoReward, error) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetClusterMemberDaoReward(clusterId)
	var queue []types.ClusterMemberDaoReward
	if store.Has(key) {
		bz := store.Get(key)
		err := util.Json.Unmarshal(bz, &queue)
		if err != nil {
			return nil, err
		}
	}
	return queue, nil
}

func (k Keeper) GetClusterDaoRewardSum(ctx sdk.Context, clusterId string) (sdkmath.Int, error) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetClusterDaoRewardSumPrefixKey(clusterId)
	var sum sdkmath.Int
	if store.Has(key) {
		bz := store.Get(key)
		err := util.Json.Unmarshal(bz, &sum)
		if err != nil {
			return sum, err
		}
		return sum, nil
	}
	return sdkmath.ZeroInt(), nil
}


func (k Keeper) SetClusterDaoRewardSum(ctx sdk.Context, clusterId string, daoReceive sdkmath.Int) error {
	store := ctx.KVStore(k.storeKey)
	key := types.GetClusterDaoRewardSumPrefixKey(clusterId)
	sum, err := k.GetClusterDaoRewardSum(ctx, clusterId)
	if err != nil {
		return err
	}
	sum = sum.Add(daoReceive)
	sumByte, err := util.Json.Marshal(sum)
	if err != nil {
		return err
	}
	store.Set(key, sumByte)
	return nil
}


func (k Keeper) SetClusterCreateTime(ctx sdk.Context, clusterId string) error {
	store := ctx.KVStore(k.storeKey)
	data := types.ClusterCreateTime{ClusterId: clusterId, CreateTime: ctx.BlockTime().Unix()}
	timeByte, err := json.Marshal(data)
	if err != nil {
		return err
	}
	store.Set(types.GetClusterTimeKey(clusterId), timeByte)
	return nil
}


func (k Keeper) InitClusterCreateTime(ctx sdk.Context, clusterId string, createTime int64) error {
	store := ctx.KVStore(k.storeKey)
	data := types.ClusterCreateTime{ClusterId: clusterId, CreateTime: createTime}
	timeByte, err := json.Marshal(data)
	if err != nil {
		return err
	}
	store.Set(types.GetClusterTimeKey(clusterId), timeByte)
	return nil
}


func (k Keeper) GetClusterCreateTime(ctx sdk.Context, clusterId string) (int64, error) {

	store := ctx.KVStore(k.storeKey)
	if !store.Has(types.GetClusterTimeKey(clusterId)) {
		return 0, core.ErrClusterIdNotFound
	}

	timeByte := store.Get(types.GetClusterTimeKey(clusterId))
	var data types.ClusterCreateTime
	err := json.Unmarshal(timeByte, &data)
	if err != nil {
		return 0, err
	}
	return data.CreateTime, nil
}

func (k Keeper) GetRemainderPool(ctx sdk.Context) (remainderPool types.RemainderPool) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.RemainderPoolKey)
	if b == nil {
		return types.RemainderPool{
			CommunityPool: sdk.DecCoins{},
		}
	}
	k.cdc.MustUnmarshal(b, &remainderPool)
	return
}

func (k Keeper) SetRemainderPool(ctx sdk.Context, remainderPool types.RemainderPool) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshal(&remainderPool)
	store.Set(types.RemainderPoolKey, b)
}


func (k Keeper) ValidateGatewaySign(ctx sdk.Context, fromAccAddr sdk.AccAddress, memberOnlineAmount int64, members []string, gatewayAddr, sign string, isAdd bool) (bool, error) {
	logs := core.BuildLog(core.GetStructFuncName(k), core.LmChainDaoKeeper)
	logs.Debug("fromAccAddr:", fromAccAddr.String())
	logs.Debug("memberOnlineAmount:", memberOnlineAmount)
	logs.Debug("members:", members)
	logs.Debug("gatewayAddr:", gatewayAddr)
	logs.Debug("sign:", sign)
	logs.Debug("isAdd:", isAdd)
	
	fromAccountI := k.accountKeeper.GetAccount(ctx, fromAccAddr)
	seq := fromAccountI.GetSequence() - 1

	
	types.SortSliceMembers(members)

	membersBytes, err := json.Marshal(members)
	if err != nil {
		logs.WithError(err).Error("Marshal members error:")
		return false, err
	}

	
	signData := types.GatewaySign{
		Members:      string(membersBytes),
		OnlineAmount: memberOnlineAmount,
		Seq:          int64(seq),
		MemberAdd:    isAdd,
	}
	signDataBytes, err := json.Marshal(signData)
	if err != nil {
		logs.WithError(err).Error("Marshal sign data error:")
		return false, err
	}

	
	signBytes, err := hex.DecodeString(sign)
	if err != nil {
		logs.WithError(err).Error("sign DecodeString error")
		return false, err
	}

	
	ethPubkey, err := util.GetPubKeyFromSign(signBytes, signDataBytes)
	if err != nil {
		logs.WithError(err).Error("Get pubkey from sign error")
		return false, err
	}

	
	ethOwnerAddr := common.HexToAddress(ethPubkey.Address().String())

	dstOwnerAddr := sdk.AccAddress(ethOwnerAddr.Bytes())
	fmt.Println(dstOwnerAddr.String())
	
	gateway, err := k.gatewayKeeper.GetGatewayInfo(ctx, gatewayAddr)
	if err != nil {
		return false, err
	}
	
	if dstOwnerAddr.String() != gateway.MachineAddress {
		return false, core.ErrMachineAddress
	}

	return true, nil
}
