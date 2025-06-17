package keeper

import (
	sdkmath "cosmossdk.io/math"
	"encoding/json"
	"fmt"
	"freemasonry.cc/blockchain/core"
	"freemasonry.cc/blockchain/util"
	"freemasonry.cc/blockchain/x/dao/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) GetAllUnreceivedPowerReward(ctx sdk.Context, addr, clusterId string) (sdkmath.Int, error) {
	logs := core.BuildLog(core.GetStructFuncName(k), core.LmChainDaoKeeper)

	key := types.GetPowerRewardCutInfoKey(addr)

	store := ctx.KVStore(k.storeKey)

	var info types.PowerRewardCycleInfo
	if store.Has(key) {
		bz := store.Get(key)
		err := json.Unmarshal(bz, &info)
		if err != nil {
			logs.Error("Failed to unmarshal power reward cycle info", "error:", err)
			return sdk.ZeroInt(), core.ErrNoPowerRewardCycleInfo
		}
	} else {
		logs.Error("No power reward cycle info found")
		return sdk.ZeroInt(), nil
	}

	reward := sdk.ZeroInt()
	for _, cycleInfo := range info.CycleInfo {
		for _, clusterInfo := range cycleInfo.ClusterRewardList {
			if clusterInfo.ClusterId == clusterId && clusterInfo.IsReceive == false {
				reward = reward.Add(clusterInfo.CutReward)
			}
		}
	}

	return reward, nil

}

func (k Keeper) UpdatePowerRewardCycleInfo(ctx sdk.Context, addr, clusterId string) error {
	logs := core.BuildLog(core.GetStructFuncName(k), core.LmChainDaoKeeper)

	store := ctx.KVStore(k.storeKey)

	key := types.GetPowerRewardCutInfoKey(addr)

	var info types.PowerRewardCycleInfo
	if store.Has(key) {
		bz := store.Get(key)
		err := json.Unmarshal(bz, &info)
		if err != nil {
			logs.Error("Failed to unmarshal power reward cycle info", "error:", err)
			return core.ErrNoPowerRewardCycleInfo
		}
	} else {
		logs.Error("No power reward cycle info found")
		return core.ErrNoPowerRewardCycleInfo
	}

	for i, v := range info.CycleInfo {
		for j, clusterInfo := range v.ClusterRewardList {
			newClusterInfo := clusterInfo
			if clusterInfo.ClusterId == clusterId {
				newClusterInfo.IsReceive = true
			}
			info.CycleInfo[i].ClusterRewardList[j] = newClusterInfo
		}
	}

	fmt.Println("info:", info)

	bz, err := json.Marshal(info)
	if err != nil {
		return err
	}

	store.Set(key, bz)

	return nil
}

func (k Keeper) GetPowerRewardCycleInfo(ctx sdk.Context, addr string) (types.PowerRewardCycleInfo, error) {
	logs := core.BuildLog(core.GetStructFuncName(k), core.LmChainDaoKeeper)
	store := ctx.KVStore(k.storeKey)

	key := types.GetPowerRewardCutInfoKey(addr)

	if store.Has(key) {
		bz := store.Get(key)
		var info types.PowerRewardCycleInfo
		err := json.Unmarshal(bz, &info)
		if err != nil {
			logs.Error("Failed to unmarshal power reward cycle info", "error:", err)
			return types.PowerRewardCycleInfo{}, core.ErrNoPowerRewardCycleInfo
		}
		return info, nil
	} else {
		return types.PowerRewardCycleInfo{
			Address:   addr,
			CycleInfo: make(map[int64]types.CycleInfo),
		}, nil
	}
}

func (k Keeper) SetPowerRewardCycleInfo(ctx sdk.Context, addr string, cycleInfo types.PowerRewardCycleInfo) error {
	store := ctx.KVStore(k.storeKey)

	key := types.GetPowerRewardCutInfoKey(addr)

	bz, err := json.Marshal(cycleInfo)
	if err != nil {
		return err
	}

	store.Set(key, bz)

	return nil
}

func (k Keeper) AddPowerRewardCycleInfo(ctx sdk.Context, addr string) (int64, error) {
	return 0, nil
}

func (k Keeper) GetCutPowerRewardCycle(ctx sdk.Context) (int64, error) {
	blockTime := ctx.BlockTime().Unix()

	cycle := (blockTime - k.GetStartTime(ctx) + core.CutProductionSeconds) / core.CutProductionSeconds
	return cycle, nil
}

func (k Keeper) StartPowerRewardCut(ctx sdk.Context, cycleInfo types.PowerRewardCycleInfo, addr string, currentCycle int64) error {
	
	clusterPersonInfo, err := k.GetPersonClusterInfo(ctx, addr)
	if err != nil {
		return err
	}

	clusterRewardList := make(map[string]types.ClusterCutReward)

	
	allReward := sdk.ZeroDec()
	
	allCutReward := sdk.ZeroInt()
	
	allPerCutReward := sdk.ZeroInt()

	
	for clusterId, _ := range clusterPersonInfo.BePower {
		cluster, err := k.GetCluster(ctx, clusterId)
		if err != nil {
			return err
		}

		endingPeriod := k.IncrementClusterPeriod(ctx, cluster)
		rewards := k.CalculateBurnRewards(ctx, cluster, addr, endingPeriod)

		clusterReward := rewards.AmountOf(core.BaseDenom)

		allReward = allReward.Add(clusterReward)

		cutReward := clusterReward.Mul(sdk.MustNewDecFromStr(core.CutPowerRewardRatio))

		perCutRewardDec := cutReward.Quo(sdk.NewDecFromInt(sdk.NewInt(core.CutPowerRewardTimes)))

		perCutReward := perCutRewardDec.TruncateInt()

		allPerCutReward = allPerCutReward.Add(perCutReward)

		allCutReward = allCutReward.Add(perCutReward.Mul(sdk.NewInt(core.CutPowerRewardTimes)))

		clusterRewardList[clusterId] = types.ClusterCutReward{
			ClusterId: clusterId,
			AllReward: clusterReward,
			CutReward: allCutReward,
			IsReceive: false,
		}
	}

	cycleInfo.CycleInfo[currentCycle] = types.CycleInfo{
		Cycle:             currentCycle,
		AllReward:         allReward,
		CutPerReward:      allPerCutReward,
		RemainReward:      allCutReward,
		ReceiveTimes:      0,
		AllCutReward:      allCutReward,
		StartTime:         ctx.BlockTime().Unix(),
		ClusterRewardList: clusterRewardList,
		Status:            1,
	}

	fmt.Println("：", allCutReward.String())
	fmt.Println("：", allPerCutReward.String())

	store := ctx.KVStore(k.storeKey)

	bz, err := json.Marshal(cycleInfo)
	if err != nil {
		return err
	}

	key := types.GetPowerRewardCutInfoKey(addr)

	store.Set(key, bz)

	return nil

}



func (k Keeper) RecordStartTime(ctx sdk.Context, timeStamp int64) {
	store := ctx.KVStore(k.storeKey)

	key := types.GetStartTimeKey()

	store.Set(key, util.Int64ToBytes(timeStamp))
}

func (k Keeper) GetStartTime(ctx sdk.Context) int64 {
	store := ctx.KVStore(k.storeKey)

	key := types.GetStartTimeKey()

	if store.Has(key) {
		bz := store.Get(key)
		return util.BytesToInt64(bz)
	} else {
		return 1713375236
	}
}
