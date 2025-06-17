package keeper

import (
	sdkmath "cosmossdk.io/math"
	"errors"
	"fmt"
	"freemasonry.cc/blockchain/core"
	"freemasonry.cc/blockchain/util"
	"freemasonry.cc/blockchain/x/dao/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authType "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/ethereum/go-ethereum/common"
)

func (k Keeper) PeriodicRewards(ctx sdk.Context) {
	
	if !k.StartMint(ctx) {
		return
	}
	
	if ctx.BlockHeight()%k.GetParams(ctx).MintBlockInterval != 0 {
		return
	}
	
	total, err := k.GetTotalPowerAmount(ctx)
	if err != nil {
		panic(err)
	}
	if total.IsZero() {
		return
	}
	params := k.GetParams(ctx)
	
	dstDec := params.DayMintAmount.Quo(sdk.NewDec(core.DayBlockNum / k.GetParams(ctx).MintBlockInterval))
	dposCoinDec := dstDec.Mul(params.DposRewardPercent)
	err = k.BankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(sdk.NewCoin(core.BaseDenom, dstDec.TruncateInt())))
	if err != nil {
		panic(err)
	}
	err = k.BankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(sdk.NewCoin(core.BaseDenom, dposCoinDec.TruncateInt())))
	if err != nil {
		panic(err)
	}

	err = k.AddMintSupply(ctx, dstDec.TruncateInt().Add(dposCoinDec.TruncateInt()))
	if err != nil {
		panic(err)
	}

	
	deviceReward := dstDec.Mul(core.ClusterDeviceRate)
	
	err = k.AddSwapDelegateSupply(ctx, dstDec.Sub(deviceReward).TruncateInt())
	if err != nil {
		panic(err)
	}

	
	err = k.BankKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, authType.FeeCollectorName, sdk.NewCoins(sdk.NewCoin(core.BaseDenom, dposCoinDec.TruncateInt())))
	if err != nil {
		panic(err)
	}

	erc20Rewards, err := k.GetErc20Reward(ctx)
	if err != nil {
		panic(err)
	}

	
	clusters := k.GetAllClusters(ctx)
	for _, cluster := range clusters {
		
		rate := cluster.ClusterPower.Quo(total)
		if rate.GT(sdk.OneDec()) {
			panic(errors.New("cluster power greater than total power"))
		}
		
		share := sdk.NewDecCoin(core.BaseDenom, dstDec.Mul(rate).TruncateInt())
		
		currentReward := k.GetClusterCurrentRewards(ctx, cluster.ClusterId)
		
		outstanding := k.GetClusterOutstandingRewards(ctx, cluster.ClusterId)
		
		currentReward.Rewards = currentReward.Rewards.Add(share)
		
		outstanding.Rewards = outstanding.Rewards.Add(share)
		if erc20Rewards != nil {
			
			for _, coin := range erc20Rewards {
				am := coin.Amount.Mul(rate)
				erc20Share := sdk.NewDecCoin(coin.Denom, am.TruncateInt())
				currentReward.Rewards = currentReward.Rewards.Add(erc20Share)
				outstanding.Rewards = outstanding.Rewards.Add(erc20Share)
			}
			
			k.DeleteErc20Reward(ctx)
		}
		
		k.SetClusterCurrentRewards(ctx, cluster.ClusterId, currentReward)
		
		k.SetClusterOutstandingRewards(ctx, cluster.ClusterId, outstanding)
	}

}


func (k Keeper) StartMint(ctx sdk.Context) bool {
	logs := core.BuildLog(core.GetStructFuncName(k), core.LmChainDaoKeeper)
	flag := k.GetGenesisIdoEndMark(ctx)
	logs.Info("k.GetGenesisIdoEndMark(ctx)", flag)
	if !flag {
		return false
	}
	nowTime := ctx.BlockTime().Unix()
	startTime, err := k.GetGenesisIdoEndTime(ctx)
	if err != nil {
		panic(err)
	}
	diff := nowTime - startTime
	if diff/core.StartMintSeconds > 0 && flag {
		return true
	}
	return false
}


func (k Keeper) CutProduct(ctx sdk.Context) {
	nowTime := ctx.BlockTime().Unix()
	startTime, err := k.GetCutProductTime(ctx)
	if err != nil {
		panic(err)
	}
	diff := nowTime - startTime
	if diff/core.CutProductionSeconds > 0 {
		params := k.GetParams(ctx)
		
		params.DayMintAmount = params.DayMintAmount.Quo(sdk.NewDec(2))
		
		
		
		
		k.SetParams(ctx, params)
		k.SetCutProductTime(ctx, nowTime)
	}
}

func (k Keeper) initializeCluster(ctx sdk.Context, val types.DeviceCluster) {
	
	k.SetClusterHistoricalRewards(ctx, val.ClusterId, 0, types.NewClusterHistoricalRewards(sdk.DecCoins{}, 1))

	
	k.SetClusterCurrentRewards(ctx, val.ClusterId, types.NewClusterCurrentRewards(sdk.DecCoins{}, 1))

	
	k.SetClusterOutstandingRewards(ctx, val.ClusterId, types.ClusterOutstandingRewards{Rewards: sdk.DecCoins{}})
}

func (k Keeper) IncrementClusterPeriod(ctx sdk.Context, val types.DeviceCluster) uint64 {
	
	rewards := k.GetClusterCurrentRewards(ctx, val.ClusterId)

	
	var current sdk.DecCoins

	if val.ClusterPower.IsZero() {

		
		
		remainderPool := k.GetRemainderPool(ctx)
		outstanding := k.GetClusterOutstandingRewards(ctx, val.ClusterId)
		remainderPool.CommunityPool = remainderPool.CommunityPool.Add(rewards.Rewards...)
		outstanding.Rewards = outstanding.GetRewards().Sub(rewards.Rewards)
		k.SetRemainderPool(ctx, remainderPool)
		k.SetClusterOutstandingRewards(ctx, val.ClusterId, outstanding)

		current = sdk.DecCoins{}
	} else {
		
		current = rewards.Rewards.QuoDecTruncate(val.ClusterPower)
	}

	
	historical := k.GetClusterHistoricalRewards(ctx, val.ClusterId, rewards.Period-1).CumulativeRewardRatio

	
	k.decrementReferenceCount(ctx, val.ClusterId, rewards.Period-1)

	
	k.SetClusterHistoricalRewards(ctx, val.ClusterId, rewards.Period, types.NewClusterHistoricalRewards(historical.Add(current...), 1))

	
	k.SetClusterCurrentRewards(ctx, val.ClusterId, types.NewClusterCurrentRewards(sdk.DecCoins{}, rewards.Period+1))

	return rewards.Period
}

func (k Keeper) decrementReferenceCount(ctx sdk.Context, clusterId string, period uint64) {
	historical := k.GetClusterHistoricalRewards(ctx, clusterId, period)
	if historical.ReferenceCount == 0 {
		panic("cannot set negative reference count")
	}
	historical.ReferenceCount--
	if historical.ReferenceCount == 0 {
		k.DeleteClusterHistoricalReward(ctx, clusterId, period)
	} else {
		k.SetClusterHistoricalRewards(ctx, clusterId, period, historical)
	}
}

func (k Keeper) incrementReferenceCount(ctx sdk.Context, clusterId string, period uint64) {
	historical := k.GetClusterHistoricalRewards(ctx, clusterId, period)
	if historical.ReferenceCount > 2 {
		panic("reference count should never exceed 2")
	}
	historical.ReferenceCount++
	k.SetClusterHistoricalRewards(ctx, clusterId, period, historical)
}

func (k Keeper) InitializeGasDelegation(ctx sdk.Context, cluster types.DeviceCluster, memberAddress string) {
	
	previousPeriod := k.GetClusterCurrentRewards(ctx, cluster.ClusterId).Period - 1
	
	k.incrementReferenceCount(ctx, cluster.ClusterId, previousPeriod)
	
	stake := cluster.ClusterPowerMembers[memberAddress].ActivePower
	
	
	k.SetBurnStartingInfo(ctx, cluster.ClusterId, memberAddress, types.NewBurnStartingInfo(previousPeriod, stake, uint64(ctx.BlockHeight())))
}

func (k Keeper) calculateBurnRewardsBetween(ctx sdk.Context, val types.DeviceCluster,
	startingPeriod, endingPeriod uint64, stake sdk.Dec,
) (rewards sdk.DecCoins) {
	
	if startingPeriod > endingPeriod {
		panic("startingPeriod cannot be greater than endingPeriod")
	}
	
	if stake.IsNegative() {
		panic("stake should not be negative")
	}
	
	starting := k.GetClusterHistoricalRewards(ctx, val.ClusterId, startingPeriod)
	ending := k.GetClusterHistoricalRewards(ctx, val.ClusterId, endingPeriod)
	difference := ending.CumulativeRewardRatio.Sub(starting.CumulativeRewardRatio)
	if difference.IsAnyNegative() {
		panic("negative rewards should not be possible")
	}
	
	rewards = difference.MulDecTruncate(stake)
	return
}

func (k Keeper) CalculateBurnRewards(ctx sdk.Context, val types.DeviceCluster, memberAddress string, endingPeriod uint64) (rewards sdk.DecCoins) {
	
	startingInfo := k.GetBurnStartingInfo(ctx, val.ClusterId, memberAddress)
	if startingInfo.Height == uint64(ctx.BlockHeight()) {
		
		return
	}
	
	if startingInfo.Stake.IsNil() {
		return
	}

	startingPeriod := startingInfo.PreviousPeriod

	stake := startingInfo.Stake

	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	

	
	
	
	
	

	currentStake := val.ClusterPowerMembers[memberAddress].ActivePower

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

	

	rewards = rewards.Add(k.calculateBurnRewardsBetween(ctx, val, startingPeriod, endingPeriod, stake)...)
	return rewards
}


func (k Keeper) addClusterMemberReward(ctx sdk.Context, clusterId, memberAddress string, finalRewards sdk.Coins) error {
	store := ctx.KVStore(k.storeKey)
	key := types.GetClusterMemberRewardKey(clusterId, memberAddress)
	var reward types.HisClusterMemberReward
	if store.Has(key) {
		rewardByte := store.Get(key)
		err := util.Json.Unmarshal(rewardByte, &reward)
		if err != nil {
			return err
		}
		reward.DeviceReward = reward.DeviceReward.Add(finalRewards.AmountOf(core.BaseDenom))
	} else {
		reward.DeviceReward = finalRewards.AmountOf(core.BaseDenom)
	}
	for _, rewards := range finalRewards {
		if rewards.Denom != core.BaseDenom {
			reward.Erc20Reward = reward.Erc20Reward.Add(rewards)
		}
	}
	rewardByte, err := util.Json.Marshal(reward)
	if err != nil {
		return err
	}
	store.Set(key, rewardByte)
	return nil
}

func (k Keeper) setClusterMemberReward(ctx sdk.Context, clusterId, memberAddress string, hisReward types.HisClusterMemberReward) error {
	store := ctx.KVStore(k.storeKey)
	key := types.GetClusterMemberRewardKey(clusterId, memberAddress)
	rewardByte, err := util.Json.Marshal(hisReward)
	if err != nil {
		return err
	}
	store.Set(key, rewardByte)
	return nil
}

func (k Keeper) deleteClusterMemberReward(ctx sdk.Context, clusterId, memberAddress string) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetClusterMemberRewardKey(clusterId, memberAddress))
}


func (k Keeper) GetClusterMemberReward(ctx sdk.Context, clusterId, memberAddress string) (types.HisClusterMemberReward, error) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetClusterMemberRewardKey(clusterId, memberAddress)
	reward := types.HisClusterMemberReward{DeviceReward: sdkmath.ZeroInt(), HisReward: sdkmath.ZeroInt()}
	if store.Has(key) {
		rewardByte := store.Get(key)
		err := util.Json.Unmarshal(rewardByte, &reward)
		if err != nil {
			return reward, err
		}
		return reward, nil
	}
	return reward, nil
}


func (k Keeper) sendRewards(ctx sdk.Context, val types.DeviceCluster, memberAddress string, finalRewards sdk.Coins, hisReward types.HisClusterMemberReward, daoPay sdkmath.Int) error {
	logs := core.BuildLog(core.GetStructFuncName(k), core.LmChainDaoKeeper)
	
	if memberAddress == val.ClusterDaoPool {
		ownerAddr, err := sdk.AccAddressFromBech32(val.ClusterOwner)
		if err != nil {
			return err
		}
		
		poolAddr, err := sdk.AccAddressFromBech32(val.ClusterVotePolicy)
		if err != nil {
			return err
		}
		
		ownerAmount := sdk.NewCoins()
		
		daoPoolAmount := sdk.NewCoins()
		
		finalRewards = finalRewards.Add(sdk.NewCoin(core.BaseDenom, hisReward.DeviceReward)).Add(sdk.NewCoin(core.BaseDenom, hisReward.HisReward))
		if finalRewards.AmountOf(core.BaseDenom).IsZero() {
			return nil
		}
		for _, reward := range finalRewards {
			if reward.Denom != core.BaseDenom {
				continue
			}
			
			am := sdk.NewDecFromInt(reward.Amount).Mul(val.ClusterDvmRatio).TruncateInt()
			ownerAmount = ownerAmount.Add(sdk.NewCoin(reward.GetDenom(), am))
			daoPoolAmount = daoPoolAmount.Add(sdk.NewCoin(reward.GetDenom(), reward.Amount.Sub(am)))
		}
		
		err = k.BankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, ownerAddr, ownerAmount)
		if err != nil {
			return err
		}

		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeWithdrawSalary,
				sdk.NewAttribute(sdk.AttributeKeyAmount, ownerAmount.AmountOf(core.BaseDenom).String()),
				sdk.NewAttribute(types.AttributeKeyCluster, val.ClusterChatId),
				sdk.NewAttribute(types.AttributeKeyMember, ownerAddr.String()),
				sdk.NewAttribute(types.AttributeClusterId, val.ClusterChatId),
				sdk.NewAttribute(types.AttributeSendeer, core.ContractAddressDao.String()),
			),
		)
		
		err = k.BankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, poolAddr, daoPoolAmount)
		if err != nil {
			return err
		}
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeWithdrawSwapDpos,
				sdk.NewAttribute(sdk.AttributeKeyAmount, daoPoolAmount.AmountOf(core.BaseDenom).String()),
				sdk.NewAttribute(types.AttributeKeyCluster, val.ClusterChatId),
				sdk.NewAttribute(types.AttributeKeyMember, poolAddr.String()),
				sdk.NewAttribute(types.AttributeClusterId, val.ClusterChatId),
			),
		)
		
		k.deleteClusterMemberReward(ctx, val.ClusterId, memberAddress)
	} else { 
		params := k.GetParams(ctx)
		daoRate := params.BurnRewardFeeRate
		addr, err := sdk.AccAddressFromBech32(memberAddress)
		if err != nil {
			return err
		}
		
		ownerAmount := sdk.NewCoins()
		
		deviceAmount := sdk.NewDecCoins()

		
		erc20Coins := k.GetErc20RewardCoins(hisReward.Erc20Reward, finalRewards)

		for _, reward := range finalRewards {
			if reward.Denom != core.BaseDenom {
				continue
			}
			
			am := sdk.NewDecFromInt(reward.Amount).Mul(core.ClusterDeviceRate).TruncateInt()
			deviceAmount = deviceAmount.Add(sdk.NewDecCoin(reward.GetDenom(), am))
			ownerAmount = ownerAmount.Add(sdk.NewCoin(reward.GetDenom(), reward.Amount.Sub(am)))
		}
		if !hisReward.DeviceReward.IsZero() {
			
			am := sdk.NewDecFromInt(hisReward.DeviceReward).Mul(core.ClusterDeviceRate).TruncateInt()
			deviceAmount = deviceAmount.Add(sdk.NewDecCoin(core.BaseDenom, am))
			
			hisReward.HisReward = hisReward.HisReward.Add(hisReward.DeviceReward.Sub(am))
			hisReward.DeviceReward = sdkmath.ZeroInt()
		}
		logs.Debug(memberAddress, " make device reward: ", core.MustParseLedgerDec(deviceAmount.AmountOf(core.BaseDenom)))
		ownerAmount = ownerAmount.Add(sdk.NewCoin(core.BaseDenom, hisReward.HisReward))
		
		powerReward, err := k.GetAllUnreceivedPowerReward(ctx, memberAddress, val.ClusterId)
		if err != nil {
			return err
		}
		if !powerReward.IsZero() {
			
			ownerAmount = ownerAmount.Sub(sdk.NewCoin(core.BaseDenom, powerReward))
			
			err = k.UpdatePowerRewardCycleInfo(ctx, memberAddress, val.ClusterId)
			if err != nil {
				return err
			}
		}
		
		daoFee := sdk.NewDecFromInt(ownerAmount.AmountOf(core.BaseDenom)).Mul(daoRate).TruncateInt()
		
		if daoFee.GT(daoPay) {
			daoFee = daoPay
			
			receiveCoin := sdk.NewDecFromInt(daoPay).Quo(params.BurnRewardFeeRate).TruncateInt()
			hisReward.HisReward = ownerAmount.Sub(sdk.NewCoin(core.BaseDenom, receiveCoin)).AmountOf(core.BaseDenom)
			
			hisReward.Erc20Reward = nil
			
			err = k.setClusterMemberReward(ctx, val.ClusterId, memberAddress, hisReward)
			if err != nil {
				return err
			}
			ownerAmount = sdk.NewCoins(sdk.NewCoin(core.BaseDenom, receiveCoin))
		} else {
			
			k.deleteClusterMemberReward(ctx, val.ClusterId, memberAddress)
		}
		daoFeeCoins := sdk.NewCoins(sdk.NewCoin(core.BurnRewardFeeDenom, daoFee))
		
		err = k.BankKeeper.SendCoinsFromAccountToModule(ctx, addr, types.ModuleName, daoFeeCoins)
		if err != nil {
			return err
		}
		
		err = k.BankKeeper.BurnCoins(ctx, types.ModuleName, daoFeeCoins)
		if err != nil {
			return err
		}
		err = k.BankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, addr, ownerAmount)
		if err != nil {
			return err
		}
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeWithdrawSwapDpos,
				sdk.NewAttribute(sdk.AttributeKeyAmount, ownerAmount.AmountOf(core.BaseDenom).String()),
				sdk.NewAttribute(types.AttributeKeyCluster, val.ClusterChatId),
				sdk.NewAttribute(types.AttributeKeyMember, memberAddress),
				sdk.NewAttribute(types.AttributeClusterId, val.ClusterChatId),
				sdk.NewAttribute(types.AttributeDaoFee, daoFee.String()),
			),
		)
		
		deviceCurrentReward := k.GetDeviceCurrentRewards(ctx, val.ClusterId)
		deviceCurrentReward.Rewards = deviceCurrentReward.Rewards.Add(deviceAmount...)
		k.SetDeviceCurrentRewards(ctx, val.ClusterId, deviceCurrentReward)
		deviceOutstanding := k.GetDeviceOutstandingRewards(ctx, val.ClusterId)
		deviceOutstanding.Rewards = deviceOutstanding.Rewards.Add(deviceAmount...)
		k.SetDeviceOutstandingRewards(ctx, val.ClusterId, deviceOutstanding)
		
		if erc20Coins != nil && !erc20Coins.IsZero() {
			logs.Info("send erc20 rewards", " cluster:", val.ClusterChatId, " member:", memberAddress, " erc20Coins:", erc20Coins)
			for _, coin := range erc20Coins {
				ethAddr := common.BytesToAddress(addr.Bytes())
				err = k.SendErc20Reward(ctx, coin.Denom, ethAddr, coin.Amount)
				if err != nil {
					logs.Info("send erc20 reward failed :", err)
					return err
				}
			}
		}
		
		err = k.SubSwapDelegateSupply(ctx, ownerAmount.AmountOf(core.BaseDenom))
		if err != nil {
			return err
		}
	}
	return nil
}


func (k Keeper) calculateWithdrawRewards(ctx sdk.Context, val types.DeviceCluster, memberAddress string) (sdk.Coins, error) {
	if _, ok := val.ClusterPowerMembers[memberAddress]; !ok {
		return nil, core.ErrNoBurn
	}
	finalRewards, err := k.withdrawSwapDpos(ctx, val, memberAddress)
	if err != nil {
		return nil, err
	}
	if !finalRewards.IsZero() {
		err = k.addClusterMemberReward(ctx, val.ClusterId, memberAddress, finalRewards)
		if err != nil {
			return nil, err
		}
		return finalRewards, nil
	}
	return nil, nil
}


func (k Keeper) withdrawAndSendRewards(ctx sdk.Context, val types.DeviceCluster, memberAddress string, daoPay sdkmath.Int) error {
	if _, ok := val.ClusterPowerMembers[memberAddress]; !ok {
		return core.ErrNoBurn
	}
	finalRewards, err := k.withdrawSwapDpos(ctx, val, memberAddress)
	if err != nil {
		return err
	}
	
	hisReward, err := k.GetClusterMemberReward(ctx, val.ClusterId, memberAddress)
	if err != nil {
		return err
	}
	if !finalRewards.IsZero() || !hisReward.HisReward.IsZero() || !hisReward.DeviceReward.IsZero() {
		err = k.sendRewards(ctx, val, memberAddress, finalRewards, hisReward, daoPay)
		if err != nil {
			return err
		}
	}
	
	k.InitializeGasDelegation(ctx, val, memberAddress)
	
	err = k.withdrawAllLpRewards(ctx, memberAddress, val.ClusterId)
	if err != nil {
		return err
	}
	return nil
}

func (k Keeper) withdrawAllLpRewards(ctx sdk.Context, memberAddress, currentClusterId string) error {
	logs := core.BuildLog(core.GetStructFuncName(k), core.LmChainDaoKeeper)
	addr, err := sdk.AccAddressFromBech32(memberAddress)
	if err != nil {
		return err
	}
	
	personalInfo, err := k.GetPersonClusterInfo(ctx, memberAddress)
	if err != nil {
		return err
	}
	
	var erc20Coins sdk.Coins
	for clusterId, _ := range personalInfo.BePower {
		
		if clusterId == currentClusterId {
			continue
		}
		cluster, err := k.GetCluster(ctx, clusterId)
		if err != nil {
			return err
		}
		_, err = k.calculateWithdrawRewards(ctx, cluster, memberAddress)
		if err != nil {
			return err
		}
		
		k.InitializeGasDelegation(ctx, cluster, memberAddress)
		
		hisReward, err := k.GetClusterMemberReward(ctx, clusterId, memberAddress)
		if err != nil {
			return err
		}
		
		erc20Lp := k.GetErc20RewardCoins(hisReward.Erc20Reward, nil)
		if erc20Lp != nil && !erc20Lp.IsZero() {
			erc20Coins = erc20Coins.Add(erc20Lp...)
			
			hisReward.Erc20Reward = nil
			
			err = k.setClusterMemberReward(ctx, clusterId, memberAddress, hisReward)
			if err != nil {
				return err
			}
		}
	}
	logs.Info("withdraw all lp rewards ", "member:", memberAddress, " erc20Coins:", erc20Coins)
	
	if erc20Coins != nil && !erc20Coins.IsZero() {
		for _, coin := range erc20Coins {
			ethAddr := common.BytesToAddress(addr.Bytes())
			err = k.SendErc20Reward(ctx, coin.Denom, ethAddr, coin.Amount)
			if err != nil {
				logs.Error("send erc20 reward failed :", err)
				return err
			}
		}
	}
	return nil
}

func (k Keeper) withdrawSwapDpos(ctx sdk.Context, val types.DeviceCluster, memberAddress string) (sdk.Coins, error) {
	
	if !k.HasBurnStartingInfo(ctx, val.ClusterId, memberAddress) {
		return nil, core.ErrEmptyBurnStartInfo
	}

	
	endingPeriod := k.IncrementClusterPeriod(ctx, val)
	rewardsRaw := k.CalculateBurnRewards(ctx, val, memberAddress, endingPeriod)
	outstanding := k.GetClusterOutstandingRewards(ctx, val.ClusterId).Rewards
	
	
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
	
	
	
	
	
	
	
	
	
	
	
	

	
	
	k.SetClusterOutstandingRewards(ctx, val.ClusterId, types.ClusterOutstandingRewards{Rewards: outstanding.Sub(rewards)})
	remainderPool := k.GetRemainderPool(ctx)
	remainderPool.CommunityPool = remainderPool.CommunityPool.Add(remainder...)
	k.SetRemainderPool(ctx, remainderPool)

	
	startingInfo := k.GetBurnStartingInfo(ctx, val.ClusterId, memberAddress)
	startingPeriod := startingInfo.PreviousPeriod
	k.decrementReferenceCount(ctx, val.ClusterId, startingPeriod)

	
	k.DeleteBurnStartingInfo(ctx, val.ClusterId, memberAddress)

	return finalRewards, nil
}


func (k Keeper) sendClusterSalary(ctx sdk.Context, burnAmount sdk.Dec, cluster types.DeviceCluster, params types.Params) error {
	
	reward := burnAmount.Mul(params.DaoRewardPercent).TruncateInt()
	
	err := k.BankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(sdk.NewCoin(core.BaseDenom, reward)))
	if err != nil {
		return err
	}
	err = k.AddMintSupply(ctx, reward)
	if err != nil {
		return err
	}
	
	ownerAmount := sdk.NewDecFromInt(reward).Mul(cluster.ClusterSalaryRatio).TruncateInt()
	ownerCoin := sdk.NewCoin(core.BaseDenom, ownerAmount)
	ownerAddr, err := sdk.AccAddressFromBech32(cluster.ClusterOwner)
	
	daoPoolCoin := sdk.NewCoin(core.BaseDenom, reward.Sub(ownerAmount))
	if err != nil {
		return err
	}
	
	err = k.BankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, ownerAddr, sdk.NewCoins(ownerCoin))
	if err != nil {
		return err
	}
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeWithdrawSalary,
			sdk.NewAttribute(sdk.AttributeKeyAmount, ownerAmount.String()),
			sdk.NewAttribute(types.AttributeKeyCluster, cluster.ClusterChatId),
			sdk.NewAttribute(types.AttributeKeyMember, ownerAddr.String()),
			sdk.NewAttribute(types.AttributeClusterId, cluster.ClusterChatId),
			sdk.NewAttribute(types.AttributeSendeer, core.ContractAddressDao.String()),
		),
	)
	poolAddr, err := sdk.AccAddressFromBech32(cluster.ClusterVotePolicy)
	if err != nil {
		return err
	}
	
	err = k.BankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, poolAddr, sdk.NewCoins(daoPoolCoin))
	if err != nil {
		return err
	}
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeWithdrawSwapDpos,
			sdk.NewAttribute(sdk.AttributeKeyAmount, daoPoolCoin.Amount.String()),
			sdk.NewAttribute(types.AttributeKeyCluster, cluster.ClusterChatId),
			sdk.NewAttribute(types.AttributeKeyMember, poolAddr.String()),
			sdk.NewAttribute(types.AttributeClusterId, cluster.ClusterChatId),
		),
	)
	return nil
}

func (k Keeper) AddSwapDelegateSupply(ctx sdk.Context, amount sdkmath.Int) error {
	store := ctx.KVStore(k.storeKey)
	key := types.GetDstDelegateSupplyPrefixKey()
	var supply sdkmath.Int
	if store.Has(key) {
		bz := store.Get(key)
		err := util.Json.Unmarshal(bz, &supply)
		if err != nil {
			return err
		}
		supply = supply.Add(amount)
	} else {
		supply = amount
	}
	supplyByte, err := util.Json.Marshal(supply)
	if err != nil {
		return err
	}
	store.Set(key, supplyByte)
	return nil
}

func (k Keeper) SubSwapDelegateSupply(ctx sdk.Context, amount sdkmath.Int) error {
	store := ctx.KVStore(k.storeKey)
	key := types.GetDstDelegateSupplyPrefixKey()
	var supply sdkmath.Int
	if store.Has(key) {
		bz := store.Get(key)
		err := util.Json.Unmarshal(bz, &supply)
		if err != nil {
			return err
		}
		supply = supply.Sub(amount)
		supplyByte, err := util.Json.Marshal(supply)
		if err != nil {
			return err
		}
		store.Set(key, supplyByte)
	}
	return nil
}

func (k Keeper) GetSwapDelegateSupply(ctx sdk.Context) (sdkmath.Int, error) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetDstDelegateSupplyPrefixKey()
	if store.Has(key) {
		bz := store.Get(key)
		var supply sdkmath.Int
		err := util.Json.Unmarshal(bz, &supply)
		if err != nil {
			return supply, err
		}
		return supply, nil
	}
	return sdkmath.NewInt(0), nil
}


func (k Keeper) SetRewardRatioYear(ctx sdk.Context) error {
	if !k.StartMint(ctx) {
		return nil
	}

	logs := core.BuildLog(core.GetStructFuncName(k), core.LmChainDaoKeeper)

	
	totalPowerAmount, err := k.GetTotalPowerAmount(ctx)
	if err != nil {
		return err
	}
	totalNotActiveAmount, err := k.GetNotActivePowerAmount(ctx)
	if err != nil {
		return err
	}
	totalPower := totalPowerAmount.Add(totalNotActiveAmount)
	if totalPower.IsZero() {
		totalPower = sdk.OneDec()
	}

	
	daymint := k.GetParams(ctx).DayMintAmount

	rewardPerPower := daymint.Quo(totalPower)

	
	err = k.SetDstPerPowerDay(ctx, rewardPerPower)
	if err != nil {
		logs.Error("set dst per power day error err", err.Error())
		return err
	}

	return nil
}

func (k Keeper) GetDstPerPowerDay(ctx sdk.Context) ([]sdk.Dec, error) {
	store := ctx.KVStore(k.storeKey)

	key := types.GetDstPerPowerDayKey()

	if !store.Has(key) {
		return make([]sdk.Dec, 0), nil
	} else {
		bz := store.Get(key)
		var dstPerPowerDay7 []sdk.Dec
		err := util.Json.Unmarshal(bz, &dstPerPowerDay7)
		if err != nil {
			return []sdk.Dec{}, err
		}

		return dstPerPowerDay7, nil
	}
}

func (k Keeper) SetDstPerPowerDay(ctx sdk.Context, newDstPerPower sdk.Dec) error {
	store := ctx.KVStore(k.storeKey)

	key := types.GetDstPerPowerDayKey()

	var dstPerPowerDay7 []sdk.Dec
	if !store.Has(key) {
		dstPerPowerDay7 = []sdk.Dec{
			newDstPerPower,
		}
	} else {
		bz := store.Get(key)
		err := util.Json.Unmarshal(bz, &dstPerPowerDay7)
		if err != nil {
			return err
		}

		
		if len(dstPerPowerDay7) < 7 {
			dstPerPowerDay7 = append([]sdk.Dec{newDstPerPower}, dstPerPowerDay7...)
		} else {
			dstPerPowerDay7 = append([]sdk.Dec{newDstPerPower}, dstPerPowerDay7[:6]...)
		}

	}

	bz, err := util.Json.Marshal(dstPerPowerDay7)
	if err != nil {
		return err
	}

	store.Set(key, bz)

	return nil
}
