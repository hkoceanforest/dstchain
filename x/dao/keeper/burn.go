package keeper

import (
	sdkmath "cosmossdk.io/math"
	"freemasonry.cc/blockchain/core"
	"freemasonry.cc/blockchain/x/dao/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authType "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distributionTypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/sirupsen/logrus"
	"strconv"
)

func (k Keeper) BurnGetPower(ctx sdk.Context, fromAddr, toAddr sdk.AccAddress, clusterId string, burnAmount, useFreezeAmount sdk.Dec, isCreate bool) (sdk.Dec, error) {
	logs := core.BuildLog(core.GetStructFuncName(k), core.LmChainDaoKeeper)

	
	if !burnAmount.IsZero() && !useFreezeAmount.IsZero() {
		return sdk.ZeroDec(), core.ErrUseFreeze
	}

	
	if !useFreezeAmount.IsZero() && !fromAddr.Equals(toAddr) {
		return sdk.ZeroDec(), core.ErrFreezeBurn
	}
	
	if !useFreezeAmount.IsZero() {
		
		pf, err := k.GetPersonClusterInfo(ctx, fromAddr.String())
		if err != nil {
			logs.WithError(err).WithField("fromAddr", fromAddr.String()).Error("BurnGetPower GetPersonClusterInfo error")
			return sdk.ZeroDec(), err
		}
		if useFreezeAmount.GT(pf.FreezePower) {
			return sdk.ZeroDec(), core.ErrFreezePowerInsufficient
		}
	}
	
	cluster, err := k.GetCluster(ctx, clusterId)
	if err != nil {
		logs.WithError(err).WithField("clusterId:", clusterId).Error("BurnGetPower get cluster error")
		return sdk.ZeroDec(), err
	}
	
	if _, ok := cluster.ClusterDeviceMembers[fromAddr.String()]; !ok {
		return sdk.ZeroDec(), core.ErrMemberNotExist
	}

	
	params := k.GetParams(ctx)
	
	burnAmountInt := burnAmount.TruncateInt()
	burnGetPowerDec, err := k.CalculateBurnGetPower(ctx, sdk.NewDecFromInt(burnAmountInt))
	if err != nil {
		return sdk.ZeroDec(), err
	}
	if !burnAmount.IsZero() { 
		
		burnCoins := sdk.NewCoins(sdk.NewCoin(core.BaseDenom, burnAmountInt))
		err = k.BankKeeper.SendCoinsFromAccountToModule(ctx, fromAddr, types.ModuleName, burnCoins)
		if err != nil {
			logs.WithError(err).WithField("burnCoins", burnCoins.String()).Error("BurnGetPower SendCoinsFromAccountToModule error")
			return sdk.ZeroDec(), core.ErrDaoBurn
		}
		
		err = k.BankKeeper.BurnCoins(ctx, types.ModuleName, burnCoins)
		if err != nil {
			logs.WithError(err).WithField("amount", burnAmount.TruncateInt()).Error("BurnGetPower BurnCoins error")
			return sdk.ZeroDec(), core.ErrDaoBurn
		}
		
		err = k.AddBurnSupply(ctx, burnAmountInt)
		if err != nil {
			logs.WithError(err).WithField("amount", burnAmount.TruncateInt()).Error("AddBurnSupply error")
			return sdk.ZeroDec(), core.ErrDaoBurn
		}
		
		err = k.SetClusterDaoRewardQueue(ctx, fromAddr.String(), clusterId, burnAmount, params)
		if err != nil {
			return sdk.ZeroDec(), err
		}
		
		err = k.burnReward(ctx, fromAddr, burnAmount, params, cluster)
		if err != nil {
			logs.WithFields(
				logrus.Fields{
					"clusterId":  clusterId,
					"burnAmount": burnAmount,
				},
			).Error("BurnGetPower burnReward error")
			return sdk.ZeroDec(), err
		}
	}
	
	err = k.personUpdate(ctx, fromAddr.String(), toAddr.String(), cluster, burnGetPowerDec, sdk.NewDecFromInt(burnAmountInt), useFreezeAmount)
	if err != nil {
		return sdk.ZeroDec(), err
	}
	
	
	if !useFreezeAmount.IsZero() {
		k.clusterLevelUpdate(ctx, &cluster, useFreezeAmount)
	} else {
		k.clusterLevelUpdate(ctx, &cluster, burnAmount)
	}

	
	addActiveAmount, err := k.clusterInfoUpdate(ctx, toAddr.String(), &cluster, burnAmount, burnGetPowerDec, useFreezeAmount)
	if err != nil {
		return sdk.ZeroDec(), err
	}
	
	k.InitializeGasDelegation(ctx, cluster, toAddr.String())

	fromBalance := k.BankKeeper.GetAllBalances(ctx, fromAddr)

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeBurn,

		
		sdk.NewAttribute(types.AttributeSendeer, fromAddr.String()),

		
		sdk.NewAttribute(types.AttributeToAddr, toAddr.String()),

		
		sdk.NewAttribute(types.AttributeKeyAmount, burnAmountInt.String()),

		
		sdk.NewAttribute(types.AttributeKeyDaoModule, authType.NewModuleAddress(types.ModuleName).String()),

		
		sdk.NewAttribute(types.AttributeSenderBalances, fromBalance.String()),

		
		sdk.NewAttribute(types.AttributeKeyClusterChatId, cluster.ClusterChatId),
	))
	
	if !isCreate {
		
		
		
		err = k.sendClusterSalary(ctx, burnAmount, cluster, params)
		if err != nil {
			return sdk.ZeroDec(), err
		}

		
		err = k.SettlementCluster(ctx, cluster, burnAmount, params)
		if err != nil {
			return sdk.ZeroDec(), err
		}
	}
	
	err = k.SetDeviceCluster(ctx, cluster)
	if err != nil {
		return sdk.ZeroDec(), nil
	}

	return addActiveAmount, nil
}


func (k Keeper) BurnGetNotActivePower(ctx sdk.Context, fromAddr sdk.AccAddress, burnAmount sdkmath.Int) error {
	logs := core.BuildLog(core.GetStructFuncName(k), core.LmChainDaoKeeper)
	supply, err := k.GetGenesisIdoSupply(ctx)
	if err != nil {
		return err
	}

	logs.Info("genesis ido supply:", supply.String(), ", burnAmount:", burnAmount.String())

	
	if supply.IsZero() {
		logs.Info("genesis ido ended")
		return core.ErrGenesisIdoEnd
	}

	burnAmountDec := sdk.NewDecFromInt(burnAmount)
	newSupply := supply.Sub(burnAmountDec)

	
	if newSupply.LT(sdk.MustNewDecFromStr("33333333333333")) {
		logs.Info("newSupply LT 33333333333333")

		burnAmount = supply.TruncateInt()
		burnAmountDec = supply
		supply = sdk.ZeroDec()
	} else {
		supply = newSupply
	}

	logs.Info("supply:", supply.String())
	logs.Info("burnAmount:", burnAmount.String())
	logs.Info("burnAmountDec:", burnAmountDec.String())
	logs.Info("supply:", supply.String())

	am := sdk.NewCoins(sdk.NewCoin(core.BaseDenom, burnAmount))
	
	err = k.BankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, fromAddr, am)
	if err != nil {
		logs.WithError(err).WithField("am", am.String()).Error("BurnGetNotActivePower SendCoinsFromModuleToAccount error")
		return err
	}
	
	err = k.BankKeeper.SendCoinsFromAccountToModule(ctx, fromAddr, types.ModuleName, am)
	if err != nil {
		logs.WithError(err).Error("BurnGetNotActivePower SendCoinsFromAccountToModule error")
		return err
	}
	
	err = k.BankKeeper.BurnCoins(ctx, types.ModuleName, am)
	if err != nil {
		logs.WithError(err).Error("BurnGetNotActivePower BurnCoins error")
		return err
	}

	
	err = k.AddBurnSupply(ctx, burnAmount)
	if err != nil {
		logs.WithError(err).Error("BurnGetNotActivePower AddBurnSupply error")
		return err
	}
	pf, err := k.GetPersonClusterInfo(ctx, fromAddr.String())
	if err != nil {
		logs.WithError(err).Error("BurnGetNotActivePower GetPersonClusterInfo error")
		return err
	}
	clusterMember := 0
	cluster := new(types.DeviceCluster)
	for key, _ := range pf.Owner {
		cl, err := k.GetCluster(ctx, key)
		if err != nil {
			continue
		}
		members := len(cl.ClusterDeviceMembers)
		if members > clusterMember {
			cluster = &cl
			clusterMember = members
		}
	}
	
	if _, ok := cluster.ClusterPowerMembers[fromAddr.String()]; ok {
		
		_, err = k.calculateWithdrawRewards(ctx, *cluster, fromAddr.String())
		if err != nil {
			logs.WithError(err).WithField("clusterId:", cluster.ClusterId).Error("BurnGetNotActivePower calculateWithdrawRewards error")
			return err
		}
	} else {
		k.IncrementClusterPeriod(ctx, *cluster)
	}

	
	cluster.ClusterPowerMembers[fromAddr.String()] = types.ClusterPowerMember{
		Address:            fromAddr.String(),
		ActivePower:        cluster.ClusterPowerMembers[fromAddr.String()].ActivePower.Add(burnAmountDec),
		BurnAmount:         cluster.ClusterPowerMembers[fromAddr.String()].ActivePower.Add(burnAmountDec),
		PowerCanReceiveDao: cluster.ClusterPowerMembers[fromAddr.String()].PowerCanReceiveDao,
	}
	
	cluster.ClusterPower = cluster.ClusterPower.Add(burnAmountDec)
	
	cluster.ClusterBurnAmount = cluster.ClusterBurnAmount.Add(burnAmountDec)

	
	err = k.AddPowerAndBurn(ctx, burnAmountDec, burnAmountDec)
	if err != nil {
		logs.WithError(err).Error("BurnGetNotActivePower AddPowerAndBurn Error")
		return err
	}

	
	k.clusterLevelUpdate(ctx, cluster, burnAmountDec)

	err = k.SetDeviceCluster(ctx, *cluster)
	if err != nil {
		logs.WithError(err).Error("BurnGetNotActivePower SetDeviceCluster Error")
		return err
	}
	
	k.InitializeGasDelegation(ctx, *cluster, fromAddr.String())
	
	pf.ActivePower = pf.ActivePower.Add(burnAmountDec)
	pf.AllBurn = pf.AllBurn.Add(burnAmountDec)
	err = k.SetPersonClusterInfo(ctx, pf)
	if err != nil {
		logs.WithError(err).Error("BurnGetNotActivePower SetPersonClusterInfo error")
		return err
	}

	if supply.IsZero() || supply.IsNegative() {
		logs.Info("genesis ido supply is zero, start swap---")
		k.DeleteGenesisIdoSupply(ctx)
		k.UpdateGenesisIdoEndMark(ctx, true)
		
		k.SetGenesisIdoEndTime(ctx, ctx.BlockTime().Unix())

		k.SetGenesisIdoSupply(ctx, sdk.ZeroDec())

		err = k.StartSwap(ctx)
		if err != nil {
			return err
		}
		return nil
	}
	k.SetGenesisIdoSupply(ctx, supply)
	return nil
}

func (k Keeper) clusterInfoUpdate(ctx sdk.Context, toAddr string, cluster *types.DeviceCluster, burnAmount, burnGetPowerDec, useFreezeAmount sdk.Dec) (sdk.Dec, error) {
	logs := core.BuildLog(core.GetStructFuncName(k), core.LmChainDaoKeeper)
	
	if !useFreezeAmount.IsZero() {
		burnGetPowerDec = burnGetPowerDec.Add(useFreezeAmount)
		burnAmount = burnAmount.Add(useFreezeAmount) 
		
		err := k.SubNotActivePowerAmount(ctx, useFreezeAmount)
		if err != nil {
			return sdk.ZeroDec(), err
		}
	}
	
	addActivePower := burnGetPowerDec
	addActiveAmount := burnAmount
	
	err := k.AddPowerAndBurn(ctx, burnAmount, burnGetPowerDec)
	if err != nil {
		return sdk.ZeroDec(), err
	}
	
	powerCanReceiveDao := burnAmount

	
	if oldPInfo, ok := cluster.ClusterPowerMembers[toAddr]; ok {
		
		powerCanReceiveDao = oldPInfo.PowerCanReceiveDao.Add(burnAmount)
		burnAmount = oldPInfo.BurnAmount.Add(burnAmount)
		burnGetPowerDec = oldPInfo.ActivePower.Add(burnGetPowerDec)
		
		_, err = k.calculateWithdrawRewards(ctx, *cluster, toAddr)
		if err != nil {
			logs.WithError(err).WithField("clusterId:", cluster.ClusterId).Error("withdrawBurnRewards error")
			return sdk.ZeroDec(), err
		}
	} else {
		k.IncrementClusterPeriod(ctx, *cluster)
	}
	
	if _, ok := cluster.ClusterPowerMembers[toAddr]; ok == false {
		if len(cluster.ClusterPowerMembers) >= core.MaxClusterPowerMembersAmount {
			return sdk.ZeroDec(), core.ErrCluserMaxPowerMembers
		}
	}
	
	cluster.ClusterPower = cluster.ClusterPower.Add(addActivePower)
	
	cluster.ClusterBurnAmount = cluster.ClusterBurnAmount.Add(addActiveAmount)
	
	cluster.ClusterPowerMembers[toAddr] = types.ClusterPowerMember{
		Address:            toAddr,
		ActivePower:        burnGetPowerDec,
		BurnAmount:         burnAmount,
		PowerCanReceiveDao: powerCanReceiveDao,
	}

	return addActivePower, nil
}


func (k Keeper) clusterLevelUpdate(ctx sdk.Context, cluster *types.DeviceCluster, burnAmount sdk.Dec) {
	newLevel := k.CalculateClusterUpgrade(ctx, *cluster, burnAmount, 0)
	
	cluster.ClusterLevel = newLevel
	
	if newLevel > cluster.ClusterLevel {
		
		ctx.EventManager().EmitEvent(sdk.NewEvent(
			types.EventTypeClusterUpgrade,
			sdk.NewAttribute(types.AttributeKeyClusterChatId, cluster.ClusterChatId), 
			sdk.NewAttribute(types.AttributeKeyClusterName, cluster.ClusterName),     
			sdk.NewAttribute(types.AttributeKeyClusterOwner, cluster.ClusterOwner),   
			sdk.NewAttribute(types.AttributeKeyOldLevel, strconv.FormatInt(cluster.ClusterLevel, 10)),
			sdk.NewAttribute(types.AttributeKeyNewLevel, strconv.FormatInt(newLevel, 10)),
		))
	}
}


func (k Keeper) personUpdate(ctx sdk.Context, fromAddr, toAddr string, cluster types.DeviceCluster, burnGetPowerDec, burnAmount, useFreezeAmount sdk.Dec) error {
	logs := core.BuildLog(core.GetStructFuncName(k), core.LmChainDaoKeeper)
	
	pt, err := k.GetPersonClusterInfo(ctx, toAddr)
	if err != nil {
		logs.WithError(err).WithField("fromAddr", toAddr).Error("BurnGetPower GetPersonClusterInfo error")
		return err
	}

	
	if pt.FirstPowerCluster == "" && pt.Address != cluster.ClusterOwner && fromAddr == toAddr {

		if len(pt.Owner) == 0 { 
			
			pt.FirstPowerCluster = cluster.ClusterId

		} else { 

			
			
			if cluster.ClusterLeader != "" {

				
				var isLoop bool
				isLoop, err = k.FindLeaderClusterOwner(ctx, pt.Owner, cluster.ClusterLeader)
				if err != nil {
					logs.WithError(err).Error("BurnGetPower FindLeaderClusterOwner error")
					return err
				}

				
				if !isLoop {

					
					pt.FirstPowerCluster = cluster.ClusterId

					
					err = k.setAllClustersLeader(ctx, pt.Owner, cluster.ClusterId)
					if err != nil {
						logs.WithError(err).Error("BurnGetPower setAllClustersLeader error")
						return err
					}
				}

			} else {
				
				pt.FirstPowerCluster = cluster.ClusterId

				
				err = k.setAllClustersLeader(ctx, pt.Owner, cluster.ClusterId)
				if err != nil {
					logs.WithError(err).Error("BurnGetPower setAllClustersLeader error")
					return err
				}
			}
		}
	}

	
	pt.ActivePower = pt.ActivePower.Add(burnGetPowerDec)
	pt.AllBurn = pt.AllBurn.Add(burnAmount)
	if !useFreezeAmount.IsZero() {
		pt.FreezePower = pt.FreezePower.Sub(useFreezeAmount)
		pt.ActivePower = pt.ActivePower.Add(useFreezeAmount)
		pt.AllBurn = pt.AllBurn.Add(useFreezeAmount) 
	}
	pt.BePower[cluster.ClusterId] = struct{}{}

	
	err = k.SetPersonClusterInfo(ctx, pt)
	if err != nil {
		logs.WithError(err).WithField("toAddr:", toAddr).Error("BurnGetPower set p cluster error")
		return err
	}
	return nil
}


func (k Keeper) setAllClustersLeader(ctx sdk.Context, owners map[string]struct{}, LeaderId string) error {
	
	for clusterId, _ := range owners {
		logs := core.BuildLog(core.GetStructFuncName(k), core.LmChainDaoKeeper)
		
		clusterInfo, err := k.GetCluster(ctx, clusterId)
		if err != nil {
			logs.WithError(err).WithField("clusterId:", clusterId).Error("setAllClustersLeader GetCluster error")
			return err
		}

		clusterInfo.ClusterLeader = LeaderId
		
		err = k.SetDeviceCluster(ctx, clusterInfo)
		if err != nil {
			logs.WithError(err).WithField("clusterId:", clusterId).Error("setAllClustersLeader SetDeviceCluster error")
			return err
		}
	}

	return nil
}

func (k Keeper) FindLeaderClusterOwner(ctx sdk.Context, myClusters map[string]struct{}, burnClusterLeader string) (bool, error) {
	logs := core.BuildLog(core.GetStructFuncName(k), core.LmChainDaoKeeper)

	
	allLeaders, err := k.GetAllLeaders(ctx, burnClusterLeader)
	if err != nil {
		logs.WithError(err).WithField("burnClusterLeader:", burnClusterLeader).Error("FindLeaderClusterOwner GetAllLeaders error")
		return false, err
	}

	
	for _, leader := range allLeaders {
		if _, ok := myClusters[leader]; ok {
			
			return true, nil
		}
	}

	
	return false, nil
}

func (k Keeper) GetAllLeaders(ctx sdk.Context, burnClusterLeader string) ([]string, error) {
	logs := core.BuildLog(core.GetStructFuncName(k), core.LmChainDaoKeeper)
	leaders := make([]string, 0)
	current := burnClusterLeader

	repeatedLeaders := make(map[string]struct{})

	for current != "" {

		
		if _, ok := repeatedLeaders[current]; ok {
			logs.WithField("current:", current).Error("GetAllLeaders loop error")
			return leaders, core.ErrErrLoopCluster
		}

		leaders = append(leaders, current)
		parent, err := k.GetCluster(ctx, current)
		if err != nil {
			logs.WithError(err).WithField("clusterId:", current).Error("GetAllLeaders GetCluster error")
			return leaders, err
		}

		if parent.ClusterLeader == "" {
			break
		}

		repeatedLeaders[current] = struct{}{}

		current = parent.ClusterLeader
	}

	return leaders, nil
}

func (k Keeper) burnReward(ctx sdk.Context, from sdk.AccAddress, burnAmount sdk.Dec, params types.Params, cluster types.DeviceCluster) error {
	logs := core.BuildLog(core.GetStructFuncName(k), core.LmChainDaoKeeper)

	
	regInt := burnAmount.Mul(params.BurnRegisterGateRatio).TruncateInt()

	
	curInt := burnAmount.Mul(params.BurnCurrentGateRatio).TruncateInt()

	
	mintAmount := regInt.Add(curInt)

	
	daoInt := burnAmount.Mul(params.BurnDaoPool).TruncateInt()

	
	err := k.BankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(sdk.NewCoin(core.BaseDenom, mintAmount), sdk.NewCoin(core.BurnRewardFeeDenom, daoInt)))
	if err != nil {
		logs.WithError(err).Error("burnReward MintCoins Error")
		return core.ErrMintReward
	}
	
	err = k.AddMintSupply(ctx, mintAmount)
	if err != nil {
		logs.WithError(err).Error("AddMintSupply Error")
		return err
	}

	
	ownerDaoReward := sdk.NewDecFromInt(daoInt).Mul(cluster.ClusterDaoRatio).TruncateInt()
	ownerAddr, err := sdk.AccAddressFromBech32(cluster.ClusterOwner)
	if err != nil {
		return core.ParseAccountError
	}
	logs.Debug("owner dao reward:", core.MustLedgerInt2RealString(ownerDaoReward), " addr:", ownerAddr.String())
	
	err = k.BankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, ownerAddr, sdk.NewCoins(sdk.NewCoin(core.BurnRewardFeeDenom, ownerDaoReward)))
	if err != nil {
		logs.WithError(err).Error("SendCoinsFromModuleToAccount Error")
		return err
	}
	daoAddr, err := sdk.AccAddressFromBech32(cluster.ClusterDaoPool)
	if err != nil {
		return core.ParseAccountError
	}
	
	err = k.BankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, daoAddr, sdk.NewCoins(sdk.NewCoin(core.BurnRewardFeeDenom, daoInt.Sub(ownerDaoReward))))
	if err != nil {
		return err
	}
	
	err = k.DaoReward(ctx, cluster)
	if err != nil {
		logs.WithError(err).Error("daoReward Error")
		return err
	}
	
	err = k.BankKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, distributionTypes.ModuleName, sdk.NewCoins(sdk.NewCoin(core.BaseDenom, regInt.Add(curInt))))
	if err != nil {
		logs.WithError(err).Error("SendCoinsFromModuleToModule Error")
		return err
	}
	
	fromInfo, err := k.ChatKeeper.GetRegisterInfo(ctx, from.String())
	if err != nil {
		logs.WithError(err).Error("burnReward GetRegisterInfo Error")
		return err
	}
	
	regGate := fromInfo.RegisterNodeAddress
	regGateAddr, err := sdk.ValAddressFromBech32(regGate)
	if err != nil {
		logs.WithError(err).WithField("regGate:", regGate).Error("burnReward ValAddressFromBech32 Error")
		return core.ErrGateway
	}
	regvalI := k.stakingKeeper.Validator(ctx, regGateAddr)
	
	k.distributionKeeper.AllocateTokensToValidator(ctx, regvalI, sdk.NewDecCoins(sdk.NewDecCoin(core.BaseDenom, regInt)))

	
	curGate := fromInfo.NodeAddress
	curGateAddr, err := sdk.ValAddressFromBech32(curGate)
	if err != nil {
		logs.WithError(err).WithField("curGate:", curGate).Error("burnReward current ValAddressFromBech32 Error")
		return core.ErrGateway
	}
	curvalI := k.stakingKeeper.Validator(ctx, curGateAddr)
	
	k.distributionKeeper.AllocateTokensToValidator(ctx, curvalI, sdk.NewDecCoins(sdk.NewDecCoin(core.BaseDenom, curInt)))

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeBurnReward,
		
		sdk.NewAttribute(types.AttributeKeyDaoModule, authType.NewModuleAddress(types.ModuleName).String()),
		
		sdk.NewAttribute(types.AttributeKeyFeeCollecter, core.ContractAddressFee.String()),
		
		
		
		sdk.NewAttribute(types.AttributeKeyDistriModule, authType.NewModuleAddress(distributionTypes.ModuleName).String()),
		
		sdk.NewAttribute(types.AttributeGatewayReward, regInt.Add(curInt).String()),
		
		sdk.NewAttribute(types.AttributeOwnerDaoReward, ownerDaoReward.String()),
		
		sdk.NewAttribute(types.AttributeKeyClusterOwner, ownerAddr.String()),
	))
	return nil
}



func (k Keeper) CalculateBurnGetPower(ctx sdk.Context, burnAmount sdk.Dec) (sdk.Dec, error) {
	
	mark := k.GetGenesisIdoEndMark(ctx)
	if !mark {
		return burnAmount, nil
	}
	
	r, err := k.GetBurnRatio(ctx)
	if err != nil {
		return r, err
	}
	
	powerGet := burnAmount.Quo(r)
	return powerGet, nil
}

func (k Keeper) GetBurnRatio(ctx sdk.Context) (sdk.Dec, error) {
	
	
	dayMintAmount := k.GetParams(ctx).DayMintAmount

	
	totalPowerAmount, err := k.GetTotalPowerAmount(ctx)
	if err != nil {
		return sdk.ZeroDec(), err
	}
	totalNotActiveAmount, err := k.GetNotActivePowerAmount(ctx)
	if err != nil {
		return sdk.ZeroDec(), err
	}
	total := totalPowerAmount.Add(totalNotActiveAmount)
	if total.IsZero() {
		total = sdk.OneDec()
	}
	r := dayMintAmount.Quo(total).Mul(k.GetParams(ctx).BurnGetPowerRatio)
	return r, nil
}

func (k Keeper) DaoReward(ctx sdk.Context, cluster types.DeviceCluster) error {
	logs := core.BuildLog(core.GetStructFuncName(k), core.LmChainDaoKeeper)
	
	queue, err := k.GetClusterDaoRewardQueue(ctx, cluster.ClusterId)
	if err != nil {
		return err
	}
	if queue == nil || len(queue) == 0 {
		return nil
	}
	
	daoPoolAddr, err := sdk.AccAddressFromBech32(cluster.ClusterDaoPool)
	if err != nil {
		return err
	}
	
	daoPoolBalance := k.BankKeeper.GetBalance(ctx, daoPoolAddr, core.BurnRewardFeeDenom).Amount
	
	if daoPoolBalance.IsZero() {
		return nil
	}
	var newQueue []types.ClusterMemberDaoReward
	var inputs []bankTypes.Input
	var outputs []bankTypes.Output
	
	for _, member := range queue {
		daoMemberGet := sdk.ZeroInt()
		addr, err := sdk.AccAddressFromBech32(member.Address)
		if err != nil {
			return err
		}
		input := bankTypes.Input{Address: daoPoolAddr.String()}
		output := bankTypes.Output{Address: addr.String()}
		
		daoReward := daoPoolBalance.QuoRaw(2)
		if daoReward.IsZero() || daoPoolBalance.IsZero() {
			break
		}
		
		if member.DaoReward.LTE(daoReward) {
			
			
			
			
			input.Coins = sdk.NewCoins(sdk.NewCoin(core.BurnRewardFeeDenom, member.DaoReward))
			output.Coins = sdk.NewCoins(sdk.NewCoin(core.BurnRewardFeeDenom, member.DaoReward))
			logs.Debug("member dao reward:", core.MustLedgerInt2RealString(member.DaoReward), " addr:", addr.String())
			daoPoolBalance = daoPoolBalance.Sub(member.DaoReward)
			daoMemberGet = member.DaoReward
			member.DaoReward = sdk.ZeroInt()
			params := k.GetParams(ctx)
			ctx.EventManager().EmitEvent(sdk.NewEvent(
				types.EventTypeCompleteDao,
				sdk.NewAttribute(types.AttributeSendeer, daoPoolAddr.String()),
				sdk.NewAttribute(types.AttributeKeyAddress, member.Address),
				sdk.NewAttribute(types.AttributeKeyAmount, member.BurnAmount.String()),
				sdk.NewAttribute(types.AttributeKeyRate, params.ReceiveDaoRatio.String()),
			))
		}
		
		if member.DaoReward.GT(daoReward) && !daoReward.IsZero() {
			member.DaoReward = member.DaoReward.Sub(daoReward)
			
			
			
			
			input.Coins = sdk.NewCoins(sdk.NewCoin(core.BurnRewardFeeDenom, daoReward))
			output.Coins = sdk.NewCoins(sdk.NewCoin(core.BurnRewardFeeDenom, daoReward))
			logs.Debug("member dao reward:", core.MustLedgerInt2RealString(daoReward), " addr:", addr.String())
			daoMemberGet = daoReward
			daoPoolBalance = daoPoolBalance.Sub(daoReward)
		}
		if !member.DaoReward.IsZero() {
			newQueue = append(newQueue, member)
		}
		inputs = append(inputs, input)
		outputs = append(outputs, output)

		ctx.EventManager().EmitEvent(sdk.NewEvent(
			types.GetDaoQueue,
			sdk.NewAttribute(types.AttributeSendeer, core.ContractAddressDao.String()),
			sdk.NewAttribute(types.AttributeReceiver, member.Address),
			sdk.NewAttribute(types.AttributeKeyAmount, daoMemberGet.String()),
		))
	}
	err = k.BankKeeper.InputOutputCoins(ctx, inputs, outputs)
	if err != nil {
		return err
	}
	return k.UpdateClusterDaoRewardQueue(ctx, newQueue, cluster.ClusterId)
}
