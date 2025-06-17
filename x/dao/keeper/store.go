package keeper

import (
	"freemasonry.cc/blockchain/x/dao/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) SetClusterHistoricalRewards(ctx sdk.Context, clusterId string, period uint64, rewards types.ClusterHistoricalRewards) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshal(&rewards)
	store.Set(types.GetClusterHistoricalRewardsKey(clusterId, period), b)
}

func (k Keeper) SetDeviceHistoricalRewards(ctx sdk.Context, clusterId string, period uint64, rewards types.ClusterHistoricalRewards) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshal(&rewards)
	store.Set(types.GetDeviceHistoricalRewardsKey(clusterId, period), b)
}

func (k Keeper) SetClusterCurrentRewards(ctx sdk.Context, clusterId string, rewards types.ClusterCurrentRewards) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshal(&rewards)
	store.Set(types.GetClusterCurrentRewardsKey(clusterId), b)
}

func (k Keeper) SetDeviceCurrentRewards(ctx sdk.Context, clusterId string, rewards types.ClusterCurrentRewards) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshal(&rewards)
	store.Set(types.GetDeviceCurrentRewardsKey(clusterId), b)
}

func (k Keeper) SetClusterOutstandingRewards(ctx sdk.Context, clusterId string, rewards types.ClusterOutstandingRewards) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshal(&rewards)
	store.Set(types.GetClusterOutstandingRewardsKey(clusterId), b)
}

func (k Keeper) SetDeviceOutstandingRewards(ctx sdk.Context, clusterId string, rewards types.ClusterOutstandingRewards) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshal(&rewards)
	store.Set(types.GetDeviceOutstandingRewardsKey(clusterId), b)
}

func (k Keeper) GetClusterHistoricalRewards(ctx sdk.Context, clusterId string, period uint64) (rewards types.ClusterHistoricalRewards) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.GetClusterHistoricalRewardsKey(clusterId, period))
	k.cdc.MustUnmarshal(b, &rewards)
	return
}

func (k Keeper) GetDeviceHistoricalRewards(ctx sdk.Context, clusterId string, period uint64) (rewards types.ClusterHistoricalRewards) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.GetDeviceHistoricalRewardsKey(clusterId, period))
	k.cdc.MustUnmarshal(b, &rewards)
	return
}

func (k Keeper) GetClusterCurrentRewards(ctx sdk.Context, clusterId string) (rewards types.ClusterCurrentRewards) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.GetClusterCurrentRewardsKey(clusterId))
	k.cdc.MustUnmarshal(b, &rewards)
	return
}

func (k Keeper) GetDeviceCurrentRewards(ctx sdk.Context, clusterId string) (rewards types.ClusterCurrentRewards) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.GetDeviceCurrentRewardsKey(clusterId))
	k.cdc.MustUnmarshal(b, &rewards)
	return
}

func (k Keeper) GetClusterOutstandingRewards(ctx sdk.Context, clusterId string) (rewards types.ClusterOutstandingRewards) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetClusterOutstandingRewardsKey(clusterId))
	k.cdc.MustUnmarshal(bz, &rewards)
	return
}

func (k Keeper) GetDeviceOutstandingRewards(ctx sdk.Context, clusterId string) (rewards types.ClusterOutstandingRewards) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetDeviceOutstandingRewardsKey(clusterId))
	k.cdc.MustUnmarshal(bz, &rewards)
	return
}

func (k Keeper) DeleteClusterHistoricalReward(ctx sdk.Context, clusterId string, period uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetClusterHistoricalRewardsKey(clusterId, period))
}

func (k Keeper) DeleteDeviceHistoricalReward(ctx sdk.Context, clusterId string, period uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetDeviceHistoricalRewardsKey(clusterId, period))
}

func (k Keeper) SetBurnStartingInfo(ctx sdk.Context, clusterId, memberAddress string, period types.BurnStartingInfo) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshal(&period)
	store.Set(types.GetBurnStartingInfoKey(clusterId, memberAddress), b)
}

func (k Keeper) SetDeviceStartingInfo(ctx sdk.Context, clusterId, memberAddress string, period types.BurnStartingInfo) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshal(&period)
	store.Set(types.GetDeviceStartingInfoKey(clusterId, memberAddress), b)
}

func (k Keeper) GetBurnStartingInfo(ctx sdk.Context, clusterId, memberAddress string) (period types.BurnStartingInfo) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.GetBurnStartingInfoKey(clusterId, memberAddress))
	k.cdc.MustUnmarshal(b, &period)
	return
}

func (k Keeper) GetDeviceStartingInfo(ctx sdk.Context, clusterId, memberAddress string) (period types.BurnStartingInfo) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.GetDeviceStartingInfoKey(clusterId, memberAddress))
	k.cdc.MustUnmarshal(b, &period)
	return
}

func (k Keeper) HasBurnStartingInfo(ctx sdk.Context, clusterId, memberAddress string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.GetBurnStartingInfoKey(clusterId, memberAddress))
}

func (k Keeper) HasDeviceStartingInfo(ctx sdk.Context, clusterId, memberAddress string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.GetDeviceStartingInfoKey(clusterId, memberAddress))
}

func (k Keeper) DeleteBurnStartingInfo(ctx sdk.Context, clusterId, memberAddress string) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetBurnStartingInfoKey(clusterId, memberAddress))
}

func (k Keeper) DeleteDeviceStartingInfo(ctx sdk.Context, clusterId, memberAddress string) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetDeviceStartingInfoKey(clusterId, memberAddress))
}
