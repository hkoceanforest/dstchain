package keeper

import (
	"fmt"
	"freemasonry.cc/blockchain/core"
	"freemasonry.cc/blockchain/x/dao/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) IncrementHistoricalRewards(ctx sdk.Context) {
	if ctx.BlockHeight()%core.DayBlockNum != 0 {
		return
	}

	clusters := k.GetAllClusters(ctx)
	for _, cluster := range clusters {
		k.IncrementDevicePeriod(ctx, cluster, true)
	}
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeIncrementHistoricalRewards,
			sdk.NewAttribute(types.AttributeKeyTime, ctx.BlockTime().String()),
		),
	)
}

func (k Keeper) initializeDevice(ctx sdk.Context, val types.DeviceCluster) {

	k.SetDeviceHistoricalRewards(ctx, val.ClusterId, 0, types.NewClusterHistoricalRewards(sdk.DecCoins{}, 1))

	k.SetDeviceCurrentRewards(ctx, val.ClusterId, types.NewClusterCurrentRewards(sdk.DecCoins{}, 1))

	k.SetDeviceOutstandingRewards(ctx, val.ClusterId, types.ClusterOutstandingRewards{Rewards: sdk.DecCoins{}})
}

func (k Keeper) IncrementDevicePeriod(ctx sdk.Context, val types.DeviceCluster, isNew bool) uint64 {
	logs := core.BuildLog(core.GetStructFuncName(k), core.LmChainDaoKeeper)
	logs.Debug("IncrementDevicePeriod Start isNew:", isNew)

	rewards := k.GetDeviceCurrentRewards(ctx, val.ClusterId)

	logs.Debug("GetDeviceCurrentRewards rewards:", rewards.String())

	var current sdk.DecCoins

	power := sdk.NewDec(int64(len(val.ClusterDeviceMembers)))
	current = rewards.Rewards.QuoDecTruncate(power)

	historical := k.GetDeviceHistoricalRewards(ctx, val.ClusterId, rewards.Period-1)

	logs.Debug("historical:", historical.String())

	hisRatio := historical.CumulativeRewardRatio
	logs.Debug("power:", power.String(), " current:", current.String(), " hisRatio:", hisRatio.String(), " rewards.period:", rewards.Period)
	if isNew {
		period := int(rewards.Period) - 2
		if period >= 0 {

			k.decrementDeviceReferenceCount(ctx, val.ClusterId, uint64(period))
		}
		logs.Debug("isNew current:", current.String())

		k.SetDeviceHistoricalRewards(ctx, val.ClusterId, rewards.Period-1, types.NewClusterHistoricalRewardsWithHis(sdk.DecCoins{}, historical.GetHisReward(), 1, historical.ReceiveCount))

		k.SetDeviceHistoricalRewards(ctx, val.ClusterId, rewards.Period, types.NewClusterHistoricalRewardsWithHis(current, rewards.Rewards, 1, 0))

		nowHis := k.GetDeviceHistoricalRewards(ctx, val.ClusterId, rewards.Period)
		logs.Debug("isNew nowHis:", nowHis.String())

		outstanding := k.GetDeviceOutstandingRewards(ctx, val.ClusterId)

		k.SetDeviceCurrentRewards(ctx, val.ClusterId, types.NewClusterCurrentRewards(outstanding.Rewards, rewards.Period+1))
	} else {
		power = power.Sub(sdk.NewDec(historical.ReceiveCount))
		logs.Debug("false power:", power.String())
		current = sdk.DecCoins{}
		if !power.IsZero() {
			current = historical.HisReward.QuoDecTruncate(power)
		}
		logs.Debug("false historical.HisReward:", historical.HisReward.String())
		logs.Debug("false current:", current.String())
		logs.Debug("false rewards.Rewards:", rewards.Rewards.String())
		k.SetDeviceHistoricalRewards(ctx, val.ClusterId, rewards.Period-1, types.NewClusterHistoricalRewardsWithHis(current, historical.GetHisReward(), 1, historical.ReceiveCount))

	}
	return rewards.Period
}

func (k Keeper) decrementDeviceReferenceCount(ctx sdk.Context, clusterId string, period uint64) {
	historical := k.GetDeviceHistoricalRewards(ctx, clusterId, period)
	if historical.ReferenceCount == 0 {
		panic("cannot set negative reference count")
	}
	historical.ReferenceCount--
	if historical.ReferenceCount == 0 {
		k.DeleteDeviceHistoricalReward(ctx, clusterId, period)
	} else {
		k.SetDeviceHistoricalRewards(ctx, clusterId, period, historical)
	}
}

func (k Keeper) incrementDeviceReferenceCount(ctx sdk.Context, clusterId string, period uint64) {
	historical := k.GetDeviceHistoricalRewards(ctx, clusterId, period)
	if historical.ReferenceCount > 2 {
		panic("reference count should never exceed 2")
	}
	historical.ReferenceCount++
	k.SetDeviceHistoricalRewards(ctx, clusterId, period, historical)
}

func (k Keeper) initializeDeviceDelegation(ctx sdk.Context, cluster types.DeviceCluster, memberAddress string, periodDiff uint64) {

	previousPeriod := int(k.GetDeviceCurrentRewards(ctx, cluster.ClusterId).Period) - int(periodDiff)
	if previousPeriod < 0 {
		previousPeriod = 0
	}

	stake := cluster.ClusterDeviceMembers[memberAddress].ActivePower

	k.SetDeviceStartingInfo(ctx, cluster.ClusterId, memberAddress, types.NewBurnStartingInfo(uint64(previousPeriod), stake, uint64(ctx.BlockHeight())))
}

func (k Keeper) calculateDeviceRewardsBetween(ctx sdk.Context, val types.DeviceCluster,
	startingPeriod, endingPeriod uint64, stake sdk.Dec,
) (rewards sdk.DecCoins) {
	logs := core.BuildLog(core.GetStructFuncName(k), core.LmChainDaoKeeper)

	if startingPeriod > endingPeriod {
		panic("startingPeriod cannot be greater than endingPeriod")
	}

	if stake.IsNegative() {
		panic("stake should not be negative")
	}

	if endingPeriod-startingPeriod > 1 {
		startingPeriod = endingPeriod - 1
	}

	starting := k.GetDeviceHistoricalRewards(ctx, val.ClusterId, startingPeriod)
	ending := k.GetDeviceHistoricalRewards(ctx, val.ClusterId, endingPeriod)
	logs.Debug("clusterId:", val.ClusterId, " startingPeriod:", startingPeriod, " endingPeriod:", endingPeriod)
	logs.Debug("starting.CumulativeRewardRatio:", starting.CumulativeRewardRatio.String())
	logs.Debug("ending.CumulativeRewardRatio:", ending.CumulativeRewardRatio.String())
	_, hasNeg := ending.CumulativeRewardRatio.SafeSub(starting.CumulativeRewardRatio)
	if hasNeg {
		rewards = sdk.DecCoins{}
		return
	}
	difference := ending.CumulativeRewardRatio.Sub(starting.CumulativeRewardRatio)
	if difference.IsAnyNegative() {
		panic("negative rewards should not be possible")
	}

	rewards = difference.MulDecTruncate(stake)
	return
}

func (k Keeper) CalculateDeviceRewards(ctx sdk.Context, val types.DeviceCluster, memberAddress string, endingPeriod uint64) (rewards sdk.DecCoins) {

	startingInfo := k.GetDeviceStartingInfo(ctx, val.ClusterId, memberAddress)
	endingPeriod = endingPeriod - 1
	startingPeriod := startingInfo.PreviousPeriod
	stake := startingInfo.Stake
	currentStake := val.ClusterDeviceMembers[memberAddress].ActivePower

	if stake.GT(currentStake) {

		marginOfErr := sdk.SmallestDec().MulInt64(3)
		if stake.LTE(currentStake.Add(marginOfErr)) {
			stake = currentStake
		} else {
			panic(fmt.Sprintf("calculated final stake for gasPower %s greater than current stake"+
				"\n\tfinal stake:\t%s"+
				"\n\tcurrent stake:\t%s",
				memberAddress, stake, currentStake))
		}
	}

	rewards = rewards.Add(k.calculateDeviceRewardsBetween(ctx, val, startingPeriod, endingPeriod, stake)...)
	return rewards
}

func (k Keeper) withdrawDeviceRewards(ctx sdk.Context, val types.DeviceCluster, memberAddress string) (sdk.Coins, error) {
	logs := core.BuildLog(core.GetStructFuncName(k), core.LmChainDaoKeeper)
	logs.Debug("withdrawDeviceRewards clusterId:", val.ClusterId, " memberAddress:", memberAddress)

	if !k.HasDeviceStartingInfo(ctx, val.ClusterId, memberAddress) {
		return nil, core.ErrEmptyBurnStartInfo
	}

	currentRewards := k.GetDeviceCurrentRewards(ctx, val.ClusterId)
	rewardsRaw := k.CalculateDeviceRewards(ctx, val, memberAddress, currentRewards.Period)
	outstanding := k.GetDeviceOutstandingRewards(ctx, val.ClusterId).Rewards

	rewards := rewardsRaw.Intersect(outstanding)
	if !rewards.IsEqual(rewardsRaw) {
		logger := k.Logger(ctx)
		logger.Info(
			"rounding error withdrawing rewards from validator",
			"member", memberAddress,
			"Cluster", val.ClusterId,
			"got", rewards.String(),
			"expected", rewardsRaw.String(),
		)
	}

	finalRewards, remainder := rewards.TruncateDecimal()

	if !finalRewards.IsZero() {

		addr, err := sdk.AccAddressFromBech32(memberAddress)
		if err != nil {
			return nil, err
		}
		err = k.BankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, addr, finalRewards)
		if err != nil {
			return nil, err
		}
	}

	balanceRewards := outstanding.Sub(rewards)
	logs.Debug("balanceRewards:", balanceRewards.String(), " rewards:", rewards.String())
	k.SetDeviceOutstandingRewards(ctx, val.ClusterId, types.ClusterOutstandingRewards{Rewards: balanceRewards})

	k.SetDeviceCurrentRewards(ctx, val.ClusterId, types.NewClusterCurrentRewards(balanceRewards, currentRewards.Period))

	historical := k.GetDeviceHistoricalRewards(ctx, val.ClusterId, currentRewards.Period-1)

	current := sdk.DecCoins{}
	if !rewards.IsZero() {
		power := sdk.NewDec(int64(len(val.ClusterDeviceMembers)) - historical.ReceiveCount - 1)
		logs.Debug("received after power:", power.String())
		if !power.IsZero() {
			current = balanceRewards.QuoDecTruncate(power)
		}
		logs.Debug("received  after current:", current.String())
		k.SetDeviceHistoricalRewards(ctx, val.ClusterId, currentRewards.Period-1, types.NewClusterHistoricalRewardsWithHis(current, balanceRewards, 1, historical.ReceiveCount+1))
	}
	remainderPool := k.GetRemainderPool(ctx)
	remainderPool.CommunityPool = remainderPool.CommunityPool.Add(remainder...)
	k.SetRemainderPool(ctx, remainderPool)

	k.DeleteDeviceStartingInfo(ctx, val.ClusterId, memberAddress)

	emittedRewards := finalRewards
	if finalRewards.IsZero() {
		baseDenom := core.BaseDenom

		emittedRewards = sdk.Coins{sdk.NewCoin(baseDenom, sdk.ZeroInt())}
	}

	fromAccAddr, err := sdk.AccAddressFromBech32(memberAddress)
	if err != nil {
		return nil, core.ErrAddressFormat

	}

	fromBalances := k.BankKeeper.GetAllBalances(ctx, fromAccAddr)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeWithdrawDeviceRewards,
			sdk.NewAttribute(sdk.AttributeKeyAmount, emittedRewards.AmountOf(core.BaseDenom).String()),
			sdk.NewAttribute(types.AttributeKeyCluster, val.ClusterId),
			sdk.NewAttribute(types.AttributeKeyMember, memberAddress),
			sdk.NewAttribute(types.AttributeSenderBalances, fromBalances.String()),
			sdk.NewAttribute(types.AttributeClusterId, val.ClusterChatId),
		),
	)

	return finalRewards, nil
}
