package keeper

import (
	sdkmath "cosmossdk.io/math"
	"freemasonry.cc/blockchain/core"
	"freemasonry.cc/blockchain/util"
	"freemasonry.cc/blockchain/x/dao/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
)

func (k Keeper) DeductEvmFee(ctx sdk.Context, toAddress common.Address) (sdkmath.Int, string, error) {
	
	approvePowerInfo, _ := k.GetContractApproveInfo(ctx, toAddress.String())
	if approvePowerInfo == nil {
		return sdkmath.ZeroInt(), "", nil
	}
	for _, val := range approvePowerInfo {
		limitPower := sdk.ZeroDec()
		if val.EndBlock < ctx.BlockHeight() {
			continue
		}
		cluster, err := k.GetCluster(ctx, val.ClusterId)
		if err != nil {
			return sdkmath.ZeroInt(), "", err
		}
		if _, ok := cluster.ClusterPowerMembers[val.Address]; !ok {
			continue
		}
		limitPower = limitPower.Add(cluster.ClusterPowerMembers[val.Address].ActivePower)
		cycleFee, err := k.getCycleEvmFee(ctx, toAddress.String(), val.ClusterId)
		if err != nil {
			return sdkmath.ZeroInt(), "", err
		}
		
		powerDec, err := k.CalculateBurnGetPower(ctx, sdk.OneDec())
		if err != nil {
			return sdkmath.ZeroInt(), "", err
		}
		
		limit := limitPower.Quo(k.GetParams(ctx).PowerGasRatio).Quo(powerDec).TruncateInt()
		
		balance := limit.Sub(cycleFee.Amount)
		if balance.IsNegative() || balance.IsZero() {
			continue
		}
		return balance, val.ClusterId, nil
	}
	return sdkmath.ZeroInt(), "", nil
}


func (k Keeper) CalculateEvmFee(ctx sdk.Context, tx sdk.Tx, toAddress common.Address, fee sdk.Coins) (sdk.Coins, error) {
	feeCoin := sdk.Coin{}
	for _, coin := range fee {
		if coin.Denom == core.BaseDenom {
			feeCoin = coin
		}
	}
	balance, clusterId, err := k.DeductEvmFee(ctx, toAddress)
	if err != nil {
		return nil, err
	}
	if balance.IsZero() {
		return fee, nil
	}
	
	if balance.LT(feeCoin.Amount) {
		
		err = k.setCycleEvmFee(ctx, toAddress.String(), clusterId, sdk.NewCoin(core.BaseDenom, balance))
		if err != nil {
			return nil, err
		}
		
		payFee := sdk.NewCoin(feeCoin.Denom, feeCoin.Amount.Sub(balance))
		events := sdk.Events{
			sdk.NewEvent(
				types.EventTypeDeductionFee,
				sdk.NewAttribute(sdk.AttributeKeyFee, feeCoin.String()),
				sdk.NewAttribute(types.AttributeDeductionFee, sdk.NewCoin(core.BaseDenom, balance).String()),
				sdk.NewAttribute(sdk.AttributeKeyFeePayer, toAddress.String()),
			),
		}
		ctx.EventManager().EmitEvents(events)
		
		return sdk.NewCoins(payFee), nil
	}
	err = k.setCycleEvmFee(ctx, toAddress.String(), clusterId, feeCoin)
	if err != nil {
		return nil, err
	}
	events := sdk.Events{
		sdk.NewEvent(
			types.EventTypeDeductionFee,
			sdk.NewAttribute(sdk.AttributeKeyFee, feeCoin.String()),
			sdk.NewAttribute(types.AttributeDeductionFee, feeCoin.String()),
			sdk.NewAttribute(sdk.AttributeKeyFeePayer, toAddress.String()),
		),
	}
	ctx.EventManager().EmitEvents(events)
	return nil, nil
}

func (k Keeper) setCycleEvmFee(ctx sdk.Context, contractAddress, clusterId string, fee sdk.Coin) error {
	store := ctx.KVStore(k.storeKey)
	data := make(map[int64]sdk.Coin)
	key := types.GetClusterEvmDeductionFeeKey(contractAddress, clusterId)
	if store.Has(key) {
		bz := store.Get(key)
		err := util.Json.Unmarshal(bz, &data)
		if err != nil {
			return err
		}
	}
	cycle := ctx.BlockHeight() / core.DayBlockNum
	if _, ok := data[cycle]; ok {
		data[cycle] = data[cycle].Add(fee)
	} else {
		data[cycle] = fee
	}
	dataByte, err := util.Json.Marshal(data)
	if err != nil {
		return err
	}
	store.Set(key, dataByte)
	return nil
}

func (k Keeper) getCycleEvmFee(ctx sdk.Context, contractAddress, clusterId string) (sdk.Coin, error) {
	store := ctx.KVStore(k.storeKey)
	zeroCoin := sdk.NewCoin(core.BaseDenom, sdk.ZeroInt())
	data := make(map[int64]sdk.Coin)
	key := types.GetClusterEvmDeductionFeeKey(contractAddress, clusterId)
	cycle := ctx.BlockHeight() / core.DayBlockNum
	if store.Has(key) {
		bz := store.Get(key)
		err := util.Json.Unmarshal(bz, &data)
		if err != nil {
			return zeroCoin, err
		}
		if _, ok := data[cycle]; ok {
			return data[cycle], nil
		}
		return zeroCoin, nil
	}
	return zeroCoin, nil
}
