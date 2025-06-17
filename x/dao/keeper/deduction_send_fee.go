package keeper

import (
	"freemasonry.cc/blockchain/core"
	"freemasonry.cc/blockchain/util"
	"freemasonry.cc/blockchain/x/dao/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)


func (k Keeper) CalculateSendFee(ctx sdk.Context, msg bankTypes.MsgSend, fee sdk.Coins) (sdk.Coins, error) {
	feeCoin := sdk.Coin{}
	for _, coin := range fee {
		if coin.Denom == core.BaseDenom {
			feeCoin = coin
		}
	}

	persionInfo, err := k.GetPersonClusterInfo(ctx, msg.FromAddress)
	if err != nil {
		return fee, nil
	}
	limitPower := persionInfo.ActivePower
	for clusterId, _ := range persionInfo.BePower {
		approve, err := k.GetClusterApproveInfo(ctx, clusterId, msg.FromAddress)
		if err != nil {
			return nil, err
		}
		
		if approve.ApproveAddress != "" && approve.EndBlock > ctx.BlockHeight() {
			cluster, err := k.GetCluster(ctx, clusterId)
			if err != nil {
				return fee, err
			}
			limitPower = limitPower.Sub(cluster.ClusterPowerMembers[msg.FromAddress].ActivePower)
		}
	}

	cycleFee, err := k.getSendCycleFee(ctx, msg.FromAddress)
	if err != nil {
		return nil, err
	}
	
	limit := limitPower.Quo(sdk.NewDec(100)).TruncateInt()
	
	balance := limit.Sub(cycleFee.Amount)
	
	if balance.LT(feeCoin.Amount) {
		
		err = k.setSendCycleFee(ctx, msg.FromAddress, sdk.NewCoin(core.BaseDenom, balance))
		if err != nil {
			return nil, err
		}
		
		payFee := sdk.NewCoin(feeCoin.Denom, feeCoin.Amount.Sub(balance))
		events := sdk.Events{
			sdk.NewEvent(
				types.EventTypeDeductionFee,
				sdk.NewAttribute(sdk.AttributeKeyFee, feeCoin.String()),
				sdk.NewAttribute(types.AttributeDeductionFee, sdk.NewCoin(core.BaseDenom, balance).String()),
				sdk.NewAttribute(sdk.AttributeKeyFeePayer, msg.FromAddress),
			),
		}
		ctx.EventManager().EmitEvents(events)
		
		return sdk.NewCoins(payFee), nil
	}
	err = k.setSendCycleFee(ctx, msg.FromAddress, feeCoin)
	if err != nil {
		return nil, err
	}
	events := sdk.Events{
		sdk.NewEvent(
			types.EventTypeDeductionFee,
			sdk.NewAttribute(sdk.AttributeKeyFee, feeCoin.String()),
			sdk.NewAttribute(types.AttributeDeductionFee, feeCoin.String()),
			sdk.NewAttribute(sdk.AttributeKeyFeePayer, msg.FromAddress),
		),
	}
	ctx.EventManager().EmitEvents(events)
	return nil, nil
}

func (k Keeper) setSendCycleFee(ctx sdk.Context, address string, fee sdk.Coin) error {
	store := ctx.KVStore(k.storeKey)
	data := make(map[int64]sdk.Coin)
	key := types.GetSendDeductionFeeKey(address)
	if store.Has(key) {
		bz := store.Get(key)
		err := util.Json.Unmarshal(bz, &data)
		if err != nil {
			return err
		}
	}
	cycle := ctx.BlockHeight() / 14400
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

func (k Keeper) getSendCycleFee(ctx sdk.Context, address string) (sdk.Coin, error) {
	store := ctx.KVStore(k.storeKey)
	zeroCoin := sdk.NewCoin(core.BaseDenom, sdk.ZeroInt())
	data := make(map[int64]sdk.Coin)
	key := types.GetSendDeductionFeeKey(address)
	cycle := ctx.BlockHeight() / 14400
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
