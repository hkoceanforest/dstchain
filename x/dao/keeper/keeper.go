package keeper

import (
	sdkmath "cosmossdk.io/math"
	"encoding/json"
	"errors"
	"fmt"
	"freemasonry.cc/blockchain/contracts"
	"freemasonry.cc/blockchain/x/contract"
	contractTypes "freemasonry.cc/blockchain/x/contract/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"math/big"
	"strconv"

	"freemasonry.cc/blockchain/core"
	"freemasonry.cc/blockchain/util"
	"freemasonry.cc/blockchain/x/dao/types"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/group"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/sirupsen/logrus"
	"github.com/tendermint/tendermint/libs/log"
)

type Keeper struct {
	storeKey           storetypes.StoreKey
	cdc                codec.BinaryCodec
	paramstore         paramtypes.Subspace
	ChatKeeper         types.ChatKeeper
	stakingKeeper      types.StakingKeeper
	accountKeeper      types.AccountKeeper
	BankKeeper         types.BankKeeper
	gatewayKeeper      types.GatewayKeeper
	distributionKeeper types.DistributionKeeper
	groupKeeper        types.GroupKeeper
	contractKeeper     types.ContractKeeper
}

func NewKeeper(
	storeKey storetypes.StoreKey,
	cdc codec.BinaryCodec,
	ps paramtypes.Subspace,
	ak types.AccountKeeper,
	bk types.BankKeeper,
	stakingKeeper types.StakingKeeper,
	gatewayKeeper types.GatewayKeeper,
	distributionKeeper types.DistributionKeeper,
	ck types.ChatKeeper,
	gk types.GroupKeeper,
	ek types.ContractKeeper,
) Keeper {
	
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}
	return Keeper{
		storeKey:           storeKey,
		cdc:                cdc,
		paramstore:         ps,
		accountKeeper:      ak,
		BankKeeper:         bk,
		stakingKeeper:      stakingKeeper,
		gatewayKeeper:      gatewayKeeper,
		distributionKeeper: distributionKeeper,
		ChatKeeper:         ck,
		groupKeeper:        gk,
		contractKeeper:     ek,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) KVHelper(ctx sdk.Context) StoreHelper {
	store := ctx.KVStore(k.storeKey)
	return StoreHelper{
		store,
	}
}


func (k Keeper) AddPowerAndBurn(ctx sdk.Context, burnAmount, power sdk.Dec) error {
	logs := core.BuildLog(core.GetStructFuncName(k), core.LmChainDaoKeeper)
	
	if !power.IsZero() {
		err := k.AddTotalPowerAmount(ctx, power)
		if err != nil {
			logs.WithError(err).WithField("newActiveAmount:", power).Error("BurnGetPower AddTotalPowerAmount error")
			return err
		}
	}
	if !burnAmount.IsZero() {
		
		err := k.AddTotalBurnAmount(ctx, burnAmount)
		if err != nil {
			logs.WithError(err).WithField("burnAmount:", burnAmount).Error("BurnGetPower AddTotalBurnAmount error")
			return err
		}
	}
	return nil
}

func (k Keeper) CreatClusterRelation(ctx sdk.Context, clusterId, leaderId string) error {
	cluster, err := k.GetCluster(ctx, clusterId)
	if err != nil {
		return err
	}
	cluster.ClusterLeader = leaderId

	err = k.SetDeviceCluster(ctx, cluster)
	if err != nil {
		return err
	}
	return nil
}

func (k Keeper) PeriodicDaoReward(ctx sdk.Context) {
	queue, err := k.GetDaoRewardQueue(ctx)
	if err != nil {
		return
	}
	if len(queue) > 0 {
		info := queue[0]
		if info.Height <= ctx.BlockHeight() {
			err = k.SettlementRouteBridgingReward(ctx, info.ClusterId, info.Address, false)
			if err != nil {
				return
			}
			
			queue = queue[1:]
			err = k.UpdateDaoRewardQueue(ctx, queue)
			if err != nil {
				return
			}
		}
	}
}

func (k Keeper) UpdateDaoRewardQueue(ctx sdk.Context, settlement []types.SettlementInfo) error {
	store := ctx.KVStore(k.storeKey)
	key := types.GetClusterDaoRewardQueuePrefixKey()
	bz, err := json.Marshal(settlement)
	if err != nil {
		return err
	}
	store.Set(key, bz)
	return nil
}

func (k Keeper) SetDaoRewardQueue(ctx sdk.Context, settlement types.SettlementInfo) error {
	store := ctx.KVStore(k.storeKey)
	key := types.GetClusterDaoRewardQueuePrefixKey()
	var settlementInfo []types.SettlementInfo
	if store.Has(key) {
		bz := store.Get(key)
		err := util.Json.Unmarshal(bz, &settlementInfo)
		if err != nil {
			return err
		}
	}
	settlementInfo = append(settlementInfo, settlement)
	bz, err := json.Marshal(settlementInfo)
	if err != nil {
		return err
	}
	store.Set(key, bz)
	return nil
}

func (k Keeper) GetDaoRewardQueue(ctx sdk.Context) ([]types.SettlementInfo, error) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetClusterDaoRewardQueuePrefixKey()
	var settlementInfo []types.SettlementInfo
	if store.Has(key) {
		bz := store.Get(key)
		err := util.Json.Unmarshal(bz, &settlementInfo)
		if err != nil {
			return nil, err
		}
	}
	return settlementInfo, nil
}

func (k Keeper) SettlementDaoReward(ctx sdk.Context, clusterId, address string) (bool, int64, error) {
	logs := core.BuildLog(core.GetStructFuncName(k), core.LmChainDaoKeeper)
	swapSupply, err := k.GetSwapDelegateSupply(ctx)
	if err != nil {
		return true, ctx.BlockHeight(), err
	}
	logs.WithField("swapSupply:", swapSupply.String()).Info("swapSupply")
	if swapSupply.IsZero() {
		return true, ctx.BlockHeight(), nil
	}
	daoSupply := k.BankKeeper.GetSupply(ctx, core.BurnRewardFeeDenom)

	if daoSupply.IsZero() {
		return true, ctx.BlockHeight(), nil
	}
	
	rate := sdk.NewDecFromInt(swapSupply).Quo(sdk.NewDecFromInt(daoSupply.Amount))

	if rate.LT(core.BaseRate) {
		return true, ctx.BlockHeight(), nil
	}
	
	diffRate := rate.Sub(core.BaseRate)
	params := k.GetParams(ctx)
	var height int64
	
	if diffRate.GTE(core.BaseRate) {
		height = core.BaseRate.Quo(params.DaoIncreaseRatio).TruncateInt().Int64() * params.DaoIncreaseHeight
	} else {
		height = diffRate.Quo(params.DaoIncreaseRatio).TruncateInt().Int64() * params.DaoIncreaseHeight
	}
	if height == 0 {
		return true, ctx.BlockHeight(), nil
	}
	settlementInfo := types.SettlementInfo{
		ClusterId: clusterId,
		Height:    ctx.BlockHeight() + height,
		Address:   address,
	}
	err = k.SetDaoRewardQueue(ctx, settlementInfo)
	if err != nil {
		return false, settlementInfo.Height, err
	}
	return false, settlementInfo.Height, nil
}


func (k Keeper) SettlementRouteBridgingReward(ctx sdk.Context, clusterId, address string, settlementFlag bool) error {
	logs := core.BuildLog(core.GetStructFuncName(k), core.LmChainDaoKeeper)
	cluster, err := k.GetCluster(ctx, clusterId)
	if err != nil {
		return err
	}
	params := k.GetParams(ctx)
	
	
	
	
	
	ownerAddr, err := sdk.AccAddressFromBech32(cluster.ClusterOwner)
	if err != nil {
		return err
	}
	routePoolAddr, err := sdk.AccAddressFromBech32(cluster.ClusterRouteRewardPool)
	if err != nil {
		return err
	}
	daoPoolAddr, err := sdk.AccAddressFromBech32(cluster.ClusterDaoPool)
	if err != nil {
		return err
	}
	
	daoReward := k.BankKeeper.GetBalance(ctx, routePoolAddr, core.BurnRewardFeeDenom).Amount
	
	if daoReward.IsZero() {
		return nil
	}
	header := ctx.BlockHeader()
	daoSendReward := sdk.NewDecFromInt(daoReward).Mul(util.MakeHashRandomRange(header.GetAppHash(), []byte(address), ctx.BlockHeight(), big.NewInt(20), big.NewInt(120)).Quo(sdk.NewDec(100))).TruncateInt()
	
	if daoSendReward.IsZero() {
		return nil
	}
	
	daoMaxAmount := k.GetBurnLevelAmount(params, cluster.ClusterLevel)
	
	receiveAmount, err := k.GetClusterDaoRewardSum(ctx, clusterId)
	if err != nil {
		return err
	}
	
	if receiveAmount.GTE(daoMaxAmount) {
		return nil
	}
	
	if daoSendReward.Add(receiveAmount).GT(daoMaxAmount) {
		daoSendReward = daoMaxAmount.Sub(receiveAmount)
	}
	
	if daoSendReward.GT(daoReward) {
		daoSendReward = daoReward
	}

	if settlementFlag {
		flag, height, err := k.SettlementDaoReward(ctx, clusterId, address)
		if err != nil {
			logs.WithError(err).WithField("clusterId:", clusterId).WithField("address:", address).WithField("daoSendReward:", daoSendReward.String()).Error("dao delay settlement error")
			return err
		}
		if !flag {
			logs.WithField("clusterId:", clusterId).WithField("address:", address).WithField("daoSendReward:", daoSendReward.String()).Info("dao delay settlement")
			ctx.EventManager().EmitEvent(sdk.NewEvent(
				types.DaoDelay,
				sdk.NewAttribute(types.AttributeClusterId, cluster.ClusterChatId),
				sdk.NewAttribute(types.AttributeKeyAddress, address),
				sdk.NewAttribute(types.AttributeKeyAmount, daoSendReward.String()),
				sdk.NewAttribute(types.AttributeKeyBlockHeight, strconv.FormatInt(height, 10)),
			))
			return nil
		}
	}

	
	ownerDao := sdk.NewDecFromInt(daoSendReward).Mul(cluster.ClusterDaoRatio).TruncateInt()
	
	err = k.BankKeeper.SendCoins(ctx, routePoolAddr, ownerAddr, sdk.NewCoins(sdk.NewCoin(core.BurnRewardFeeDenom, ownerDao)))
	if err != nil {
		return err
	}
	memberDao := daoSendReward.Sub(ownerDao)
	
	err = k.BankKeeper.SendCoins(ctx, routePoolAddr, daoPoolAddr, sdk.NewCoins(sdk.NewCoin(core.BurnRewardFeeDenom, memberDao)))
	if err != nil {
		return err
	}
	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.GetDaoRoute,
		sdk.NewAttribute(types.AttributeSendeer, routePoolAddr.String()),
		sdk.NewAttribute(types.AttributeReceiver, ownerAddr.String()),
		sdk.NewAttribute(types.AttributeKeyAmount, ownerDao.String()),
	))
	
	
	
	
	
	
	err = k.DaoReward(ctx, cluster)
	if err != nil {
		return err
	}
	
	err = k.SetClusterDaoRewardSum(ctx, clusterId, daoSendReward)
	if err != nil {
		return err
	}
	
	if cluster.ClusterLeader != "" {
		leaderCluster, err := k.GetCluster(ctx, cluster.ClusterLeader)
		if err != nil {
			return err
		}
		
		
		
		
		
		
		
		
		
		
		daoInt := sdk.NewDecFromInt(daoSendReward).Mul(params.ConnectivityDaoRatio).TruncateInt()
		if daoInt.IsZero() {
			return nil
		}
		
		
		
		
		
		daoCoins := sdk.NewCoins(sdk.NewCoin(core.BurnRewardFeeDenom, daoInt))
		
		err = k.BankKeeper.MintCoins(ctx, types.ModuleName, daoCoins)
		if err != nil {
			return err
		}
		
		leaderAddr, err := sdk.AccAddressFromBech32(leaderCluster.ClusterRouteRewardPool)
		if err != nil {
			return err
		}
		err = k.BankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, leaderAddr, daoCoins)
		if err != nil {
			return err
		}
	}
	return nil
}


func (k Keeper) SettlementCluster(ctx sdk.Context, cluster types.DeviceCluster, burnAmount sdk.Dec, params types.Params) error {
	
	if cluster.ClusterLeader == "" || burnAmount.IsZero() {
		return nil
	}
	
	leaderCluster, err := k.GetCluster(ctx, cluster.ClusterLeader)
	if err != nil {
		return err
	}
	
	
	
	
	
	
	
	
	
	
	daoInt := burnAmount.Mul(params.BurnDaoPool).Mul(params.ConnectivityDaoRatio).TruncateInt()
	daoCoins := sdk.NewCoins(sdk.NewCoin(core.BurnRewardFeeDenom, daoInt))
	
	err = k.BankKeeper.MintCoins(ctx, types.ModuleName, daoCoins)
	if err != nil {
		return err
	}
	
	
	
	
	
	
	leaderAddr, err := sdk.AccAddressFromBech32(leaderCluster.ClusterRouteRewardPool)
	if err != nil {
		return err
	}
	return k.BankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, leaderAddr, daoCoins)
	
	
	
}

func (k Keeper) AddTotalBurnAmount(ctx sdk.Context, burnAmount sdk.Dec) error {
	logs := core.BuildLog(core.GetStructFuncName(k), core.LmChainDaoKeeper)

	key := types.TotalBurnAmount

	store := ctx.KVStore(k.storeKey)

	var oldAmount sdk.Dec
	if store.Has(key) {
		oldAmountByte := store.Get(key)
		err := util.Json.Unmarshal(oldAmountByte, &oldAmount)
		if err != nil {
			logs.WithError(err).WithField("oldAmountByte", oldAmountByte).Error("AddTotalBurnAmount oldAmount Json.Unmarshal error")
			return core.ErrAddTotalBurnAmount
		}
	} else {
		oldAmount = sdk.ZeroDec()
	}

	newAmount := oldAmount.Add(burnAmount)
	newAmountByte, err := util.Json.Marshal(newAmount)
	if err != nil {
		logs.WithError(err).WithField("newAmount", newAmount).Error("AddTotalBurnAmount newAmount Json.Marshal error")
		return core.ErrAddTotalBurnAmount
	}

	store.Set(key, newAmountByte)

	return nil
}

func (k Keeper) GetTotalBurnAmount(ctx sdk.Context) (sdk.Dec, error) {
	log := core.BuildLog(core.GetStructFuncName(k), core.LmChainDaoKeeper)

	key := types.TotalBurnAmount

	store := ctx.KVStore(k.storeKey)

	var amount sdk.Dec
	if store.Has(key) {
		amountByte := store.Get(key)
		err := util.Json.Unmarshal(amountByte, &amount)
		if err != nil {
			log.WithError(err).WithField("amountByte", amountByte).Error("GetTotalBurnAmount Json.Unmarshal error")
			return sdk.ZeroDec(), core.ErrGetTotalBurnAmount
		}
	} else {
		amount = sdk.ZeroDec()
	}

	return amount, nil
}

func (k Keeper) SetNotActivePowerAmount(ctx sdk.Context, burnAmount sdk.Dec) {
	key := types.GetNotActivePowerPrefixPrefixKey()
	store := ctx.KVStore(k.storeKey)
	powerByte, err := util.Json.Marshal(burnAmount)
	if err != nil {
		panic(core.ErrAddTotalBurnAmount)
	}
	store.Set(key, powerByte)
}


func (k Keeper) SubNotActivePowerAmount(ctx sdk.Context, burnAmount sdk.Dec) error {
	log := core.BuildLog(core.GetStructFuncName(k), core.LmChainDaoKeeper)
	key := types.GetNotActivePowerPrefixPrefixKey()
	store := ctx.KVStore(k.storeKey)
	power := sdk.ZeroDec()
	if !store.Has(key) {
		return nil
	}
	powerByte := store.Get(key)
	err := util.Json.Unmarshal(powerByte, &power)
	if err != nil {
		log.WithError(err).WithField("power", powerByte).Error("SubNotActivePowerAmount power Json.Unmarshal error")
		return core.ErrSubTotalBurnAmount
	}
	power = power.Sub(burnAmount)
	if power.IsZero() || power.IsNegative() {
		store.Delete(key)
		return nil
	}
	powerByte, err = util.Json.Marshal(power)
	if err != nil {
		log.WithError(err).WithField("power", power).Error("SubNotActivePowerAmount power Json.Marshal error")
		return core.ErrSubTotalBurnAmount
	}
	store.Set(key, powerByte)
	return nil
}

func (k Keeper) GetNotActivePowerAmount(ctx sdk.Context) (sdk.Dec, error) {
	logs := core.BuildLog(core.GetStructFuncName(k), core.LmChainDaoKeeper)
	key := types.GetNotActivePowerPrefixPrefixKey()
	store := ctx.KVStore(k.storeKey)
	amount := sdk.ZeroDec()
	if store.Has(key) {
		amountByte := store.Get(key)
		err := util.Json.Unmarshal(amountByte, &amount)
		if err != nil {
			logs.WithError(err).WithField("amountByte", amountByte).Error("GetNotActivePowerAmount Json.Unmarshal error")
			return sdk.ZeroDec(), core.ErrGetTotalBurnAmount
		}
	}
	return amount, nil
}

func (k Keeper) AddTotalPowerAmount(ctx sdk.Context, burnAmount sdk.Dec) error {
	log := core.BuildLog(core.GetStructFuncName(k), core.LmChainDaoKeeper)

	key := types.TotalPowerAmount

	store := ctx.KVStore(k.storeKey)

	var oldAmount sdk.Dec
	if store.Has(key) {
		oldAmountByte := store.Get(key)
		err := util.Json.Unmarshal(oldAmountByte, &oldAmount)
		if err != nil {
			log.WithError(err).WithField("oldAmountByte", oldAmountByte).Error("AddTotalPowerAmount oldAmount Json.Unmarshal error")
			return core.ErrAddTotalBurnAmount
		}
	} else {
		oldAmount = sdk.ZeroDec()
	}

	newAmount := oldAmount.Add(burnAmount)
	newAmountByte, err := util.Json.Marshal(newAmount)
	if err != nil {
		log.WithError(err).WithField("newAmount", newAmount).Error("AddTotalPowerAmount newAmount Json.Marshal error")
		return core.ErrAddTotalBurnAmount
	}

	store.Set(key, newAmountByte)

	return nil
}

func (k Keeper) GetTotalPowerAmount(ctx sdk.Context) (sdk.Dec, error) {
	log := core.BuildLog(core.GetStructFuncName(k), core.LmChainDaoKeeper)

	key := types.TotalPowerAmount

	store := ctx.KVStore(k.storeKey)

	var amount sdk.Dec
	if store.Has(key) {
		amountByte := store.Get(key)
		err := util.Json.Unmarshal(amountByte, &amount)
		if err != nil {
			log.WithError(err).WithField("amountByte", amountByte).Error("GetTotalPowerAmount Json.Unmarshal error")
			return sdk.ZeroDec(), core.ErrGetTotalBurnAmount
		}
	} else {
		amount = sdk.ZeroDec()
	}
	return amount, nil
}

func (k Keeper) GetCurrentBurnLevel(ctx sdk.Context, burnAmount sdk.Dec) (types.BurnLevel, error) {
	levels := k.GetParams(ctx).BurnLevels
	
	for i := len(levels); i > 0; i-- {
		if burnAmount.GTE(levels[i-1].BurnAmount) {
			return levels[i-1], nil
		}
	}
	return types.BurnLevel{
		Level:      1,
		BurnAmount: sdk.ZeroDec(),
		AddPercent: sdk.ZeroInt(),
		RoomAmount: sdk.ZeroInt(),
	}, nil
}

func (k Keeper) GetBurnLevelByAccAddress(ctx sdk.Context, accAddress sdk.AccAddress) (types.BurnLevel, error) {
	log := core.BuildLog(core.GetStructFuncName(k), core.LmChainDaoKeeper)
	
	p, err := k.GetPersonClusterInfo(ctx, accAddress.String())
	if err != nil {
		log.WithError(err).Error("QueryTotalMedalGetAmount")
		return types.BurnLevel{}, err
	}

	
	burnLevel, err := k.GetCurrentBurnLevel(ctx, p.AllBurn)
	if err != nil {
		log.WithError(err).Error("GetCurrentPledgeLevel")
		return types.BurnLevel{}, err
	}

	return burnLevel, nil
}

func (k Keeper) GetBurnLevelByAccAddresses(ctx sdk.Context, addresses []string) (burnLevelInfos map[string]int64, err error) {
	logs := core.BuildLog(core.GetStructFuncName(k), core.LmChainDaoKeeper)

	burnLevelInfos = make(map[string]int64, 0)

	for _, address := range addresses {
		accAddresses, err := sdk.AccAddressFromBech32(address)
		if err != nil {
			logs.WithError(err).Error("QueryBurnLevelByAccAddresses sdk.AccAddressFromBech32")
			return nil, err
		}

		burnLevelInfo, err := k.GetBurnLevelByAccAddress(ctx, accAddresses)
		if err != nil {
			logs.WithError(err).Error("QueryBurnLevelByAccAddresses err")
			return nil, err
		}
		burnLevelInfos[address] = burnLevelInfo.Level
	}
	return burnLevelInfos, nil
}

func (k Keeper) AddMember(ctx sdk.Context, clusterId, from string, members []types.Members, memberOnlineAmount int64, isCreateCluster bool) error {
	logs := core.BuildLog(core.GetStructFuncName(k), core.LmChainDaoKeeper)
	
	logs.Debug("AddMember Start-----------")
	cluster, err := k.GetCluster(ctx, clusterId)
	if err != nil {
		logs.WithError(err).WithField("clusterId", clusterId).Error("AddMember GetClusterByChatId Error")
		return err
	}
	logs.Debug("add1")
	
	clusterAuth := k.IsOwnerOrAdmin(cluster, from)
	if !clusterAuth {
		return core.ErrClusterPermission
	}
	logs.Debug("add2")
	var n float32

	
	groupMembers := []group.MemberRequest{}
	for _, m := range members {
		logs.Debug("add member 0:", m.MemberAddress)
		
		_, err = k.ChatKeeper.GetRegisterInfo(ctx, m.MemberAddress)

		logs.Debug("add member 1")
		if err != nil {
			if err == core.ErrUserNotFound { 
				logs.Debug("add member 2")
				accM, err := sdk.AccAddressFromBech32(m.MemberAddress)
				if err != nil {
					return core.ErrAddressFormat
				}
				logs.Debug("add member 3")
				accExists := k.accountKeeper.HasAccount(ctx, accM)
				if !accExists {
					n += 1
					k.accountKeeper.SetAccount(ctx, k.accountKeeper.NewAccountWithAddress(ctx, accM))
				}
				logs.Debug("add member 4")
				
				gatewayInfo, err := k.gatewayKeeper.GetGatewayInfoByNum(ctx, m.IndexNum)
				if err != nil {
					logs.WithError(err).WithFields(logrus.Fields{
						"clusterId": clusterId,
						"indexNum":  m.IndexNum,
					}).Error("AddMember GetClusterByChatId Error")
					return err
				}
				logs.Debug("add member 5")
				if gatewayInfo == nil {
					logs.WithError(err).WithFields(logrus.Fields{
						"clusterId": clusterId,
						"indexNum":  m.IndexNum,
					}).Error("AddMember gateway not found Error")
					return core.ErrGatewayNotFound
				}
				logs.Debug("add member 6")
				
				err = k.ChatKeeper.Register(ctx, m.MemberAddress, m.ChatAddress, gatewayInfo.GatewayAddress, nil)
				if err != nil {
					logs.WithError(err).Error("AddMember Register error")
					return err
				}
				logs.Debug("add member 7")
			} else {
				logs.Debug("add member 8")
				logs.WithError(err).Error("AddMember GetRegisterInfo Error")
				return err
			}
		}
		logs.Debug("add member 9")
		defer telemetry.IncrCounter(n, "new", "account")

		
		cluster.ClusterDeviceMembers[m.MemberAddress] = types.ClusterDeviceMember{
			Address:     m.MemberAddress,
			ActivePower: sdk.NewDec(1),
		}
		logs.Debug("add member 10")
		
		p, err := k.GetPersonClusterInfo(ctx, m.MemberAddress)
		if err != nil {
			logs.WithError(err).WithField("p", m).Error("AddMember GetPersonClusterInfo Error")
			return err
		}
		logs.Debug("add member 11")
		p.Device[cluster.ClusterId] = struct{}{}
		err = k.SetPersonClusterInfo(ctx, p)
		if err != nil {
			logs.WithError(err).WithField("p", m).Error("AddMember SetPersonClusterInfo Error")
			return err
		}
		logs.Debug("add member 12")
		
		k.initializeDeviceDelegation(ctx, cluster, m.MemberAddress, 2)
		logs.Debug("add member 13")
		
		k.IncrementDevicePeriod(ctx, cluster, false)
		logs.Debug("add member 14")
		groupMembers = append(groupMembers, group.MemberRequest{Address: m.MemberAddress, Weight: "20"})
		logs.Debug("add member 15")
	}

	logs.Debug("add3")

	
	if isCreateCluster {
		logs.Debug("add3.1")
		
		cluster.ClusterActiveDevice = cluster.ClusterActiveDevice + memberOnlineAmount - 1
	} else {
		logs.Debug("add3.2")
		cluster.ClusterActiveDevice = cluster.ClusterActiveDevice + memberOnlineAmount
	}

	logs.Debug("add4")

	if cluster.ClusterActiveDevice < 0 {
		cluster.ClusterActiveDevice = 0
	}

	cluster.OnlineRatio = sdk.NewDec(cluster.ClusterActiveDevice).Quo(sdk.NewDec(int64(len(cluster.ClusterDeviceMembers))))

	logs.Debug("add5")
	err = k.SetDeviceCluster(ctx, cluster)
	if err != nil {
		logs.WithError(err).WithField("clusterTrueId", clusterId).Error("AddMember SetDeviceCluster Error")
		return err
	}
	logs.Debug("add6")
	params := k.GetParamsIfExists(ctx)
	if params.IdoMinMember == 0 {
		params.IdoMinMember = types.DefaultParams().IdoMinMember
	}
	logs.Debug("add7")

	
	if !k.GetGenesisIdoEndMark(ctx) && int64(len(cluster.ClusterDeviceMembers)) >= params.IdoMinMember {
		logs.Debug("add7.5")
		ownerAcc, err := sdk.AccAddressFromBech32(cluster.ClusterOwner)
		if err != nil {
			logs.WithError(err).Error("AccAddressFromBech32 error")
			return core.ErrAddressFormat
		}
		ownerEth := common.BytesToAddress(ownerAcc.Bytes())
		abi := contracts.GenesisIdoNContract.ABI

		resp, err := k.contractKeeper.CallEVM(ctx, abi, contractTypes.ModuleAddress, contract.GenesisIdoContract, true, "addWhitelist", ownerEth)
		if err != nil {
			logs.WithError(err).Error("failed to callEVM")
			return err
		}
		if resp.Failed() {
			logs.Errorf("failed to callEVM: %s", resp.VmError)
			return errors.New(resp.VmError)
		}
	}

	logs.Debug("add8")

	
	if len(groupMembers) > 0 {
		err := k.UpdateGroupMembers(ctx, cluster, groupMembers)
		if err != nil {
			logs.WithError(err).WithField("group", clusterId).Error("AddMember AddGroupMembers Error")
			return err
		}
	}

	logs.Debug("add9")
	return nil
}

func (k Keeper) UpdateClusterOnlineRate(ctx sdk.Context, clusterId string, rate sdk.Dec) error {
	logs := core.BuildLog(core.GetStructFuncName(k), core.LmChainDaoKeeper)

	
	cluster, err := k.GetCluster(ctx, clusterId)
	if err != nil {
		logs.WithError(err).WithField("clusterID", clusterId).Error("UpdateClusterOnlineRate GetCluster Error")
		return err
	}

	
	cluster.OnlineRatio = rate

	
	cluster.OnlineRatioUpdateTime = ctx.BlockTime().Unix()

	
	deviceAmountDec := sdk.NewDec(int64(len(cluster.ClusterDeviceMembers)))

	activeAmount := deviceAmountDec.Mul(rate)

	cluster.ClusterActiveDevice = activeAmount.RoundInt64()

	err = k.SetDeviceCluster(ctx, cluster)
	if err != nil {
		logs.WithError(err).WithField("clusterID", clusterId).Error("UpdateClusterOnlineRate SetDeviceCluster Error")
		return err
	}

	return nil
}


func (k Keeper) GetGenesisIdoEndMark(ctx sdk.Context) bool {
	store := ctx.KVStore(k.storeKey)
	if store.Has(types.GetMintMarkPrefixKey(core.BaseDenom)) {
		bz := store.Get(types.GetMintMarkPrefixKey(core.BaseDenom))
		if string(bz) == "0" {
			return false
		}
		return true
	}
	return false
}


func (k Keeper) UpdateGenesisIdoEndMark(ctx sdk.Context, mark bool) {
	store := ctx.KVStore(k.storeKey)
	if mark {
		store.Set(types.GetMintMarkPrefixKey(core.BaseDenom), []byte("1"))
		return
	}
	store.Set(types.GetMintMarkPrefixKey(core.BaseDenom), []byte("0"))
}


func (k Keeper) SetCutProductTime(ctx sdk.Context, time int64) {
	store := ctx.KVStore(k.storeKey)
	timeByte, err := util.Json.Marshal(time)
	if err != nil {
		panic(err)
	}
	store.Set(types.GetCutProductionPrefixKey(), timeByte)
}

func (k Keeper) GetCutProductTime(ctx sdk.Context) (int64, error) {
	store := ctx.KVStore(k.storeKey)
	var startTime int64
	timeByte := store.Get(types.GetCutProductionPrefixKey())
	if timeByte != nil {
		err := util.Json.Unmarshal(timeByte, &startTime)
		if err != nil {
			return 0, err
		}
	}
	return startTime, nil
}


func (k Keeper) SetChainStartTime(ctx sdk.Context, time int64) {
	store := ctx.KVStore(k.storeKey)
	timeByte, err := util.Json.Marshal(time)
	if err != nil {
		panic(err)
	}
	store.Set(types.GetChainStartTimePrefixKey(), timeByte)
}


func (k Keeper) SetGenesisIdoEndTime(ctx sdk.Context, time int64) {
	store := ctx.KVStore(k.storeKey)
	timeByte, err := util.Json.Marshal(time)
	if err != nil {
		panic(err)
	}
	store.Set(types.GetChainStartTimePrefixKey(), timeByte)
}

func (k Keeper) GetGenesisIdoEndTime(ctx sdk.Context) (int64, error) {
	store := ctx.KVStore(k.storeKey)
	var startTime int64
	timeByte := store.Get(types.GetChainStartTimePrefixKey())
	if timeByte != nil {
		err := util.Json.Unmarshal(timeByte, &startTime)
		if err != nil {
			return 0, err
		}
	}
	return startTime, nil
}


func (k Keeper) SetGenesisIdoSupply(ctx sdk.Context, supply sdk.Dec) {
	store := ctx.KVStore(k.storeKey)
	supplyByte, err := util.Json.Marshal(supply)
	if err != nil {
		panic(err)
	}
	store.Set(types.GetGenesisIdoSupplyPrefixKey(), supplyByte)
}

func (k Keeper) GetGenesisIdoSupply(ctx sdk.Context) (sdk.Dec, error) {
	store := ctx.KVStore(k.storeKey)
	supply := sdk.ZeroDec()
	supplyByte := store.Get(types.GetGenesisIdoSupplyPrefixKey())
	if supplyByte != nil {
		err := util.Json.Unmarshal(supplyByte, &supply)
		if err != nil {
			return supply, err
		}
	}
	return supply, nil
}

func (k Keeper) DeleteGenesisIdoSupply(ctx sdk.Context) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetGenesisIdoSupplyPrefixKey())
}


func (k Keeper) AddBurnSupply(ctx sdk.Context, amount sdkmath.Int) error {
	store := ctx.KVStore(k.storeKey)
	key := types.GetBurnSupplyPrefixKey()
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


func (k Keeper) GetBurnSupply(ctx sdk.Context) (sdkmath.Int, error) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetBurnSupplyPrefixKey()
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


func (k Keeper) AddMintSupply(ctx sdk.Context, amount sdkmath.Int) error {
	store := ctx.KVStore(k.storeKey)
	key := types.GetMintSupplyPrefixKey()
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


func (k Keeper) GetMintSupply(ctx sdk.Context) (sdkmath.Int, error) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetMintSupplyPrefixKey()
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

func (k Keeper) StartSwap(ctx sdk.Context) error {
	logs := core.BuildLog(core.GetStructFuncName(k), core.LmChainDaoKeeper)
	swapSwitchAbi := contracts.SwapSwitchJSONContract.ABI
	resp, err := k.contractKeeper.CallEVM(ctx, swapSwitchAbi, contractTypes.ModuleAddress, contract.SwapSwitchContract, true, "setIsOpen", true)
	if err != nil {
		logs.WithError(err).Error("failed to callEVM")
		return err
	}
	if resp.Failed() {
		logs.Errorf("failed to callEVM: %s", resp.VmError)
		return errors.New(resp.VmError)
	}
	return nil
}

func (k Keeper) ExchangeAddLiquidity(ctx sdk.Context) error {
	logs := core.BuildLog(core.GetStructFuncName(k), core.LmChainDaoKeeper)
	
	
	
	
	
	
	

	dstAmount := new(big.Int)
	dstAmount, ok := dstAmount.SetString(core.IdoEndAddLiquidityAmountDST, 10)
	if !ok {
		return errors.New("invalid dst amount")
	}
	dstAmountBig := (*hexutil.Big)(dstAmount)

	usdtAddress := common.HexToAddress("0x80b5a32E4F032B2a058b4F29EC95EEfEEB87aDcd")

	usdtAmount := new(big.Int)
	usdtAmount, ok = usdtAmount.SetString(core.IdoEndAddLiquidityAmountUSDT, 10)
	if !ok {
		return errors.New("invalid usdt amount")
	}

	logs.Info("ctx.BlockTime().Unix():", ctx.BlockTime().Unix())

	resp, err := k.contractKeeper.CallEVMWithValue(
		ctx,
		contracts.ExchangeRouterContract.ABI,
		common.BytesToAddress(authtypes.NewModuleAddress(contractTypes.GenesisIdoReward).Bytes()),
		contract.ExchangeRouterContract,
		true,
		"addLiquidityETH",
		dstAmountBig,
		usdtAddress,
		usdtAmount,
		new(big.Int).SetInt64(0),
		new(big.Int).SetInt64(0),
		common.HexToAddress("0x000000000000000000000000000000000000dEaD"),
		new(big.Int).SetInt64(ctx.BlockTime().Unix()+20),
	)
	if err != nil {
		logs.WithError(err).Error("failed to callEVM")
		return err
	}

	if resp.Failed() {
		logs.Errorf("failed to callEVM vmError: %s", resp.VmError)
		return errors.New(resp.VmError)
	}

	return nil
}

func (k Keeper) SetExchangeLiquidityAdded(ctx sdk.Context) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GeExchangeLiquidityAddedKey(), []byte("1"))
}

func (k Keeper) SetUsdtAuthorizationAdded(ctx sdk.Context) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetUsdtAuthKey(), []byte("1"))
}

func (k Keeper) GetUsdtAuthorizationAdded(ctx sdk.Context) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.GetUsdtAuthKey())
}

func (k Keeper) GetExchangeLiquidityAdded(ctx sdk.Context) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.GeExchangeLiquidityAddedKey())
}

func (k Keeper) GrantAuthorizationUsdt(ctx sdk.Context) error {
	logs := core.BuildLog(core.GetStructFuncName(k), core.LmChainDaoKeeper)
	logs.Info("GrantAuthorizationUsdt+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++")
	if k.GetUsdtAuthorizationAdded(ctx) {
		logs.Info("GrantAuthorizationUsdt has added")
		return nil
	}

	approveAmount, ok := new(big.Int).SetString(core.IdoEndAddLiquidityAmountUSDT, 10)
	if !ok {
		return errors.New("invalid approve amount")
	}

	resp, err := k.contractKeeper.CallEVM(
		ctx,
		contracts.UsdtContract.ABI,
		common.BytesToAddress(authtypes.NewModuleAddress(contractTypes.GenesisIdoReward).Bytes()),
		contract.UsdtContract,
		true,
		"approve",
		contract.ExchangeRouterContract,
		approveAmount,
	)
	if err != nil {
		logs.WithError(err).Error("failed to callEVM: GrantAuthorizationUsdt:" + err.Error())
		return err
	}

	if resp.Failed() {
		logs.Errorf("failed to callEVM vmError: %s", resp.VmError)
		return errors.New(resp.VmError)
	}
	k.SetUsdtAuthorizationAdded(ctx)
	return nil
}
