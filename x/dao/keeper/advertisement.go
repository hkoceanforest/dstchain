package keeper

import (
	sdkmath "cosmossdk.io/math"
	"freemasonry.cc/blockchain/contracts"
	"freemasonry.cc/blockchain/core"
	"freemasonry.cc/blockchain/util"
	"freemasonry.cc/blockchain/x/contract"
	"freemasonry.cc/blockchain/x/dao/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
	"strings"
)


func (k Keeper) ClustersFilter(ctx sdk.Context, startTime, endTime, level int64, rate sdk.Dec) (clusters []types.DeviceCluster) {
	k.IterateAllClusterCreateTime(ctx, func(cluster types.ClusterCreateTime) bool {
		if cluster.CreateTime >= startTime && cluster.CreateTime <= endTime {
			info, err := k.GetCluster(ctx, cluster.ClusterId)
			if err != nil {
				return true
			}
			clusters = append(clusters, info)
		}
		return false
	})
	
	newClusters := []types.DeviceCluster{}
	for _, cluster := range clusters {
		if cluster.OnlineRatio.LT(rate) || cluster.ClusterLevel < level {
			continue
		}
		newClusters = append(newClusters, cluster)
	}
	return newClusters
}

func (k Keeper) IterateAllClusterCreateTime(ctx sdk.Context, cb func(cluster types.ClusterCreateTime) (stop bool)) {
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, types.ClusterTimePrefix)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		cluster := new(types.ClusterCreateTime)
		err := util.Json.Unmarshal(iterator.Value(), cluster)
		if err != nil {
			panic(err)
		}
		if cb(*cluster) {
			break
		}
	}
}

func (k Keeper) GetExchangeRate(ctx sdk.Context, address sdk.AccAddress, params types.Params) (sdkmath.Int, error) {
	from := common.BytesToAddress(address.Bytes())
	priceAbi := contracts.PriceJSONContract.ABI
	args := params.AdPrice.Mul(core.UsdtRealToLedgerRateDec)
	resp, err := k.contractKeeper.CallEVM(ctx, priceAbi, from, contract.PriceContract, false, "getPrice", args.TruncateInt().BigInt())
	if err != nil {
		return sdkmath.ZeroInt(), err
	}

	unPackData, err := priceAbi.Unpack("getPrice", resp.Ret)
	if err != nil {
		return sdkmath.ZeroInt(), err
	}

	rate, ok := unPackData[0].(*big.Int)
	if !ok {
		return sdkmath.ZeroInt(), core.ErrPriceFormat
	}

	return sdkmath.NewIntFromBigInt(rate), nil
}


func (k Keeper) clusterAd(ctx sdk.Context, clusterIds []string, fromAddr sdk.AccAddress) error {
	var err error
	amount := sdkmath.ZeroInt()
	clusterAdAmount := []types.ClusterDaoPoolFee{}
	params := k.GetParams(ctx)
	rate, err := k.GetExchangeRate(ctx, fromAddr, params)
	if err != nil {
		return err
	}
	for _, clusterId := range clusterIds {
		if strings.Contains(clusterId, ".") {
			
			clusterId, err = k.GetClusterId(ctx, clusterId)
			if err != nil {
				return err
			}
		}
		cluster, err := k.GetCluster(ctx, clusterId)
		if err != nil {
			return err
		}
		adAmount := rate.MulRaw(int64(len(cluster.ClusterDeviceMembers)))
		daoPool := types.ClusterDaoPoolFee{ClusterId: clusterId, DaoPool: cluster.ClusterDaoPool, Amount: sdk.NewCoins(sdk.NewCoin(core.BaseDenom, adAmount))}
		clusterAdAmount = append(clusterAdAmount, daoPool)
		amount = amount.Add(adAmount)
	}
	
	adCoin := sdk.NewCoin(core.BaseDenom, amount)
	err = k.BankKeeper.SendCoinsFromAccountToModule(ctx, fromAddr, types.ModuleName, sdk.NewCoins(adCoin))
	if err != nil {
		return err
	}
	
	feeCoin := sdk.NewCoin(core.BaseDenom, sdk.NewDecFromInt(amount).Mul(params.AdRate).TruncateInt())
	err = k.BankKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, authtypes.FeeCollectorName, sdk.NewCoins(feeCoin))
	if err != nil {
		return err
	}
	
	for _, info := range clusterAdAmount {
		addr, err := sdk.AccAddressFromBech32(info.DaoPool)
		if err != nil {
			continue
		}
		
		coinRate := sdk.OneDec().Sub(params.AdRate)
		daoAmount := sdk.NewCoin(core.BaseDenom, sdk.NewDecFromInt(info.Amount.AmountOf(core.BaseDenom)).Mul(coinRate).TruncateInt())
		err = k.BankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, addr, sdk.NewCoins(daoAmount))
		if err != nil {
			return err
		}
	}
	return nil
}
