package dao

import (
	"encoding/json"
	"freemasonry.cc/blockchain/x/dao/keeper"
	"freemasonry.cc/blockchain/x/dao/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func InitGenesis(
	ctx sdk.Context,
	k keeper.Keeper,
	data types.GenesisState,
) {
	
	err := k.AddMintSupply(ctx, sdk.MustNewDecFromStr("77700000000000000000000000").TruncateInt())
	if err != nil {
		panic(err)
	}

	k.SetParams(ctx, data.Params)

	if len(data.Clusters) != 0 {
		for _, cluster := range data.Clusters {
			dm := make(map[string]types.ClusterDeviceMember)
			for dk, dv := range cluster.ClusterDeviceMembers {
				dm[dk] = types.ClusterDeviceMember{
					Address:     dv.Address,
					ActivePower: dv.ActivePower,
				}
			}

			dp := make(map[string]types.ClusterPowerMember)
			for dk, dv := range cluster.ClusterPowerMembers {
				dp[dk] = types.ClusterPowerMember{
					Address:            dv.Address,
					ActivePower:        dv.ActivePower,
					BurnAmount:         dv.BurnAmount,
					PowerCanReceiveDao: dv.PowerCanReceiveDao,
				}
			}

			
			adminList := map[string]struct{}{}
			err := json.Unmarshal([]byte(cluster.ClusterAdminList), &adminList)
			if err != nil {
				panic(err)
			}

			clusterInit := types.DeviceCluster{
				ClusterId:              cluster.ClusterId,
				ClusterChatId:          cluster.ClusterChatId,
				ClusterName:            cluster.ClusterName,
				ClusterOwner:           cluster.ClusterOwner,
				ClusterGateway:         cluster.ClusterGateway,
				ClusterLeader:          cluster.ClusterLeader,
				ClusterDeviceMembers:   dm,
				ClusterPowerMembers:    dp,
				ClusterPower:           cluster.ClusterPower,
				ClusterLevel:           cluster.ClusterLevel,
				ClusterBurnAmount:      cluster.ClusterBurnAmount,
				ClusterActiveDevice:    cluster.ClusterActiveDevice,
				ClusterDaoPool:         cluster.ClusterDaoPool,
				ClusterRouteRewardPool: cluster.ClusterRouteRewardPool,
				ClusterDeviceRatio:     cluster.ClusterDeviceRatio,
				ClusterSalaryRatioUpdateHeight: types.ClusterChangeRatioHeight{
					SalaryRatioUpdateHeight: cluster.ClusterSalaryRatioUpdateHeight.SalaryRatioUpdateHeight,
					DvmRatioUpdateHeight:    cluster.ClusterSalaryRatioUpdateHeight.DvmRatioUpdateHeight,
					DaoRatioUpdateHeight:    cluster.ClusterSalaryRatioUpdateHeight.DaoRatioUpdateHeight,
				},
				ClusterDeviceRatioUpdateHeight: cluster.ClusterDeviceRatioUpdateHeight,
				ClusterSalaryRatio:             cluster.ClusterSalaryRatio,
				ClusterDvmRatio:                cluster.ClusterDvmRatio,
				ClusterDaoRatio:                cluster.ClusterDaoRatio,
				OnlineRatio:                    cluster.OnlineRatio,
				OnlineRatioUpdateTime:          cluster.OnlineRatioUpdateTime,
				ClusterAdminList:               adminList,
				ClusterVoteId:                  cluster.ClusterVoteId,
				ClusterVotePolicy:              cluster.ClusterVotePolicy,
			}

			err = k.SetDeviceCluster(ctx, clusterInit)
			if err != nil {
				panic(err)
			}

			
			
			k.SetClusterHistoricalRewards(ctx, cluster.ClusterId, 0, types.NewClusterHistoricalRewards(sdk.DecCoins{}, 1))

			
			k.SetClusterCurrentRewards(ctx, cluster.ClusterId, types.NewClusterCurrentRewards(sdk.DecCoins{}, 1))
			
			k.SetClusterOutstandingRewards(ctx, cluster.ClusterId, types.ClusterOutstandingRewards{Rewards: sdk.DecCoins{}})
			
			k.SetDeviceHistoricalRewards(ctx, cluster.ClusterId, 0, types.NewClusterHistoricalRewards(sdk.DecCoins{}, 1))

			
			k.SetDeviceCurrentRewards(ctx, cluster.ClusterId, types.NewClusterCurrentRewards(sdk.DecCoins{}, 1))

			
			k.SetDeviceOutstandingRewards(ctx, cluster.ClusterId, types.ClusterOutstandingRewards{Rewards: sdk.DecCoins{}})

			
			previousPeriod := int(k.GetDeviceCurrentRewards(ctx, cluster.ClusterId).Period) - int(1)
			if previousPeriod < 0 {
				previousPeriod = 0
			}
			
			

			k.InitializeGasDelegation(ctx, clusterInit, clusterInit.ClusterOwner)

			for _, member := range clusterInit.ClusterDeviceMembers {
				
				stake := cluster.ClusterDeviceMembers[member.Address].ActivePower
				
				
				k.SetDeviceStartingInfo(ctx, cluster.ClusterId, member.Address, types.NewBurnStartingInfo(uint64(previousPeriod), stake, uint64(ctx.BlockHeight())))
			}

		}
	}

	
	for _, personalCluster := range data.PersonalClusters {
		device := map[string]struct{}{}
		err := json.Unmarshal([]byte(personalCluster.Device), &device)
		if err != nil {
			panic(err)
		}

		owner := map[string]struct{}{}
		err = json.Unmarshal([]byte(personalCluster.Owner), &owner)
		if err != nil {
			panic(err)
		}

		bePower := map[string]struct{}{}
		err = json.Unmarshal([]byte(personalCluster.BePower), &bePower)
		if err != nil {
			panic(err)
		}

		err = k.SetPersonClusterInfo(ctx, types.PersonalClusterInfo{
			Address:           personalCluster.Address,
			Device:            device,
			Owner:             owner,
			BePower:           bePower,
			AllBurn:           personalCluster.AllBurn,
			ActivePower:       personalCluster.ActivePower,
			FreezePower:       personalCluster.FreezePower,
			FirstPowerCluster: personalCluster.FirstPowerCluster,
		})
		if err != nil {
			panic(err)
		}
	}

	
	if data.GatewayCluster != "" {
		var gatewayCluster []types.GatewayClusterExport
		err := json.Unmarshal([]byte(data.GatewayCluster), &gatewayCluster)
		if err != nil {
			panic(err)
		}

		for _, info := range gatewayCluster {
			for clusterId := range info.Clusters {
				err = k.AddClusterToGateway(ctx, info.GatewayAddress, clusterId)
			}
		}
	}

	
	if data.ClusterChatIdReflection != nil {
		for _, clusterIdInfo := range data.ClusterChatIdReflection {
			k.SetClusterId(ctx, clusterIdInfo.ClusterId, clusterIdInfo.ClusterChatId)
		}
	}

	
	if len(data.ClusterStrategyAddress) != 0 {
		for _, policyAddrInfo := range data.ClusterStrategyAddress {
			k.SetClusterVotePolicy(ctx, policyAddrInfo.ClusterId, policyAddrInfo.StrategyAddress)
		}
	}

	
	if len(data.ClusterCreateTime) != 0 {
		for _, createTimeInfo := range data.ClusterCreateTime {
			err := k.InitClusterCreateTime(ctx, createTimeInfo.ClusterId, createTimeInfo.CreateTime)
			if err != nil {
				panic(err)
			}
		}
	}

	k.InitGenesis(ctx, &data)
}

func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	params := k.GetParams(ctx)

	c := k.GetAllClusters(ctx)
	clusters := make([]types.DeviceClusterExport, 0)
	for _, cluster := range c {
		adminList, err := json.Marshal(cluster.ClusterAdminList)
		if err != nil {
			panic(err)
		}

		dm := make(map[string]types.ClusterDeviceMemberExport)
		for dk, dv := range cluster.ClusterDeviceMembers {
			dm[dk] = types.ClusterDeviceMemberExport{
				Address:     dv.Address,
				ActivePower: dv.ActivePower,
			}
		}

		dp := make(map[string]types.ClusterPowerMemberExport)
		for dk, dv := range cluster.ClusterPowerMembers {

			dp[dk] = types.ClusterPowerMemberExport{
				Address:            dv.Address,
				ActivePower:        dv.ActivePower,
				BurnAmount:         dv.BurnAmount,
				PowerCanReceiveDao: dv.PowerCanReceiveDao,
			}
		}

		clusters = append(clusters, types.DeviceClusterExport{
			ClusterId:                      cluster.ClusterId,
			ClusterChatId:                  cluster.ClusterChatId,
			ClusterName:                    cluster.ClusterName,
			ClusterOwner:                   cluster.ClusterOwner,
			ClusterGateway:                 cluster.ClusterGateway,
			ClusterLeader:                  cluster.ClusterLeader,
			ClusterDeviceMembers:           dm,
			ClusterPowerMembers:            dp,
			ClusterPower:                   cluster.ClusterPower,
			ClusterLevel:                   cluster.ClusterLevel,
			ClusterBurnAmount:              cluster.ClusterBurnAmount,
			ClusterActiveDevice:            cluster.ClusterActiveDevice,
			ClusterDaoPool:                 cluster.ClusterDaoPool,
			ClusterRouteRewardPool:         cluster.ClusterRouteRewardPool,
			ClusterDeviceRatio:             cluster.ClusterDeviceRatio,
			ClusterDeviceRatioUpdateHeight: cluster.ClusterDeviceRatioUpdateHeight,
			ClusterSalaryRatio:             cluster.ClusterSalaryRatio,
			ClusterSalaryRatioUpdateHeight: types.ClusterChangeRatioHeightExport{
				SalaryRatioUpdateHeight: cluster.ClusterSalaryRatioUpdateHeight.SalaryRatioUpdateHeight,
				DvmRatioUpdateHeight:    cluster.ClusterSalaryRatioUpdateHeight.DvmRatioUpdateHeight,
				DaoRatioUpdateHeight:    cluster.ClusterSalaryRatioUpdateHeight.DaoRatioUpdateHeight,
			},
			ClusterDvmRatio:       cluster.ClusterDvmRatio,
			ClusterDaoRatio:       cluster.ClusterDaoRatio,
			OnlineRatio:           cluster.OnlineRatio,
			OnlineRatioUpdateTime: cluster.OnlineRatioUpdateTime,
			ClusterAdminList:      string(adminList),
			ClusterVoteId:         cluster.ClusterVoteId,
			ClusterVotePolicy:     cluster.ClusterVotePolicy,
		})
	}

	p := k.GetAllPersonClusters(ctx)
	personalClusters := make([]types.PersonalClusterInfoExport, 0)
	for _, personalCluster := range p {
		device, err := json.Marshal(personalCluster.Device)
		if err != nil {
			panic(err)
		}

		owner, err := json.Marshal(personalCluster.Owner)
		if err != nil {
			panic(err)
		}

		bePower, err := json.Marshal(personalCluster.BePower)
		if err != nil {
			panic(err)
		}

		personalClusters = append(personalClusters, types.PersonalClusterInfoExport{
			Address:           personalCluster.Address,
			Device:            string(device),
			Owner:             string(owner),
			BePower:           string(bePower),
			AllBurn:           personalCluster.AllBurn,
			ActivePower:       personalCluster.ActivePower,
			FreezePower:       personalCluster.FreezePower,
			FirstPowerCluster: personalCluster.FirstPowerCluster,
		})
	}

	
	clusterChatIdReflection := k.GetAllClusterAddress(ctx)

	
	gatewayCluster := k.GetAllGatewayClusters(ctx)
	gatewayClusterJson, err := json.Marshal(gatewayCluster)
	if err != nil {
		panic(err)
	}

	
	clusterStrategyAddress := k.GetAllClusterVotePolicy(ctx)

	
	clusterCreateTime := k.GetAllClusterCreateTime(ctx)

	timeExport := make([]types.ClusterCreateTimeExport, 0)
	for _, ctime := range clusterCreateTime {
		timeExport = append(timeExport, types.ClusterCreateTimeExport{
			ClusterId:  ctime.ClusterId,
			CreateTime: ctime.CreateTime,
		})
	}

	return &types.GenesisState{
		Params:                  params,
		Clusters:                clusters,
		PersonalClusters:        personalClusters,
		ClusterChatIdReflection: clusterChatIdReflection,
		GatewayCluster:          string(gatewayClusterJson),
		ClusterStrategyAddress:  clusterStrategyAddress,
		ClusterCreateTime:       timeExport,
	}
}
