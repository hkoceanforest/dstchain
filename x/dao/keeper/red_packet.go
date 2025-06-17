package keeper

import (
	sdkmath "cosmossdk.io/math"
	"freemasonry.cc/blockchain/core"
	"freemasonry.cc/blockchain/util"
	"freemasonry.cc/blockchain/x/dao/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"math/big"
	"strconv"
)

func (k Keeper) SetRidClusterId(ctx sdk.Context, redpacketId, clusterTrueId string) error {
	store := ctx.KVStore(k.storeKey)

	
	key := types.GetRidClusterIdKey(redpacketId)

	
	if store.Has(key) {
		return core.ErrRedPacketExist
	}

	data := []byte(clusterTrueId)

	
	store.Set(key, data)

	return nil
}

func (k Keeper) GetRidClusterId(ctx sdk.Context, redpacketId string) (string, error) {
	store := ctx.KVStore(k.storeKey)

	
	key := types.GetRidClusterIdKey(redpacketId)

	
	if !store.Has(key) {
		return "", core.ErrRedPacketNotExist
	}

	
	dataByte := store.Get(key)

	return string(dataByte), nil
}

func (k Keeper) InitRedPacketInfo(ctx sdk.Context, fromAddr sdk.AccAddress, clusterTrueId, redpacketId string, amount sdk.Coin, count int64, PacketType int64) error {
	logs := core.BuildLog(core.GetStructFuncName(k), core.LmChainDaoKeeper)
	store := ctx.KVStore(k.storeKey)

	
	key := types.GetRedPacketKey(redpacketId)

	
	if store.Has(key) {
		return core.ErrRedPacketExist
	}

	
	redPacket := types.RedPacket{
		EndBlock:      ctx.BlockHeight() + core.RedPacketTimeOut,
		Id:            redpacketId,
		Sender:        fromAddr.String(),
		ClusterTrueId: clusterTrueId,
		Amount:        amount,
		Count:         count,
		Receive:       make([]types.RedPacketReceive, 0),
		PacketType:    PacketType,
	}

	
	redPacketData, err := util.Json.Marshal(redPacket)
	if err != nil {
		logs.WithError(err).Error("InitRedPacket Marshal Error")
		return core.ErrRedPacketMarshal
	}

	store.Set(key, redPacketData)

	return nil
}

func (k Keeper) GetRedPacketInfo(ctx sdk.Context, redpacketId string) (*types.RedPacket, error) {
	logs := core.BuildLog(core.GetStructFuncName(k), core.LmChainDaoKeeper)
	store := ctx.KVStore(k.storeKey)

	key := types.GetRedPacketKey(redpacketId)

	if !store.Has(key) {
		return nil, core.ErrRedPacketNotExist
	}

	redPacketBytes := store.Get(key)

	var redPacket types.RedPacket
	err := util.Json.Unmarshal(redPacketBytes, &redPacket)
	if err != nil {
		logs.WithError(err).Error("redPacket Unmarshal")
		return nil, core.ErrRedPacketUnmarshal
	}
	return &redPacket, nil
}

func (k Keeper) SetRedPacketInfo(ctx sdk.Context, redPacket types.RedPacket) error {
	logs := core.BuildLog(core.GetStructFuncName(k), core.LmChainDaoKeeper)
	store := ctx.KVStore(k.storeKey)

	key := types.GetRedPacketKey(redPacket.Id)

	
	if !store.Has(key) {
		return core.ErrRedPacketNotExist
	}

	
	redPacketData, err := util.Json.Marshal(redPacket)
	if err != nil {
		logs.WithError(err).Error("SetRedPacket Marshal Error")
		return core.ErrRedPacketMarshal
	}

	store.Set(key, redPacketData)

	return nil
}

func (k Keeper) ReceiveRedPacket(ctx sdk.Context, fromAddr sdk.AccAddress, redPacket *types.RedPacket, clusterChatId string) (sdkmath.Int, error) {
	if redPacket.PacketType == types.RedPacketTypeNormal {
		return k.ReceiveNormalRedPacket(ctx, fromAddr, redPacket, clusterChatId)
	} else if redPacket.PacketType == types.RedPacketTypeLucky {
		return k.ReceiveLuckyRedPacket(ctx, fromAddr, redPacket, clusterChatId)
	}
	return sdkmath.ZeroInt(), nil
}

func (k Keeper) ReceiveLuckyRedPacket(ctx sdk.Context, fromAddr sdk.AccAddress, redPacket *types.RedPacket, clusterChatId string) (sdkmath.Int, error) {
	logs := core.BuildLog(core.GetStructFuncName(k), core.LmChainDaoKeeper)

	var luckyAmount sdkmath.Int
	if int64(len(redPacket.Receive)) == redPacket.Count-1 { 
		luckyAmount = redPacket.Remain()
	} else { 
		
		min := new(big.Int).Div(redPacket.Remain().BigInt(), big.NewInt(redPacket.Count-int64(len(redPacket.Receive))))
		max := new(big.Int).Mul(min, big.NewInt(2))
		luckyAmountDec := util.MakeHashRandomRange(ctx.BlockHeader().AppHash, nil, ctx.BlockHeight(), min, max)
		luckyAmount = luckyAmountDec.TruncateInt()
	}

	
	redPacket.Receive = append(redPacket.Receive, types.RedPacketReceive{
		Receiver: fromAddr.String(),
		Amount:   luckyAmount,
	})

	err := k.SetRedPacketInfo(ctx, *redPacket)
	if err != nil {
		logs.WithError(err).Error("SetRedPacketInfo Error")
		return sdkmath.ZeroInt(), err
	}

	
	err = k.BankKeeper.SendCoinsFromModuleToAccount(
		ctx,
		types.ModuleName,
		fromAddr,
		sdk.NewCoins(sdk.NewCoin(redPacket.Amount.Denom, luckyAmount)),
	)
	if err != nil {
		logs.WithError(err).Error("Receive LuckyRedPacket SendCoinsFromModuleToAccount Error")
		return sdkmath.ZeroInt(), err
	}

	
	countRemain := sdk.NewInt(redPacket.Count - int64(len(redPacket.Receive)))

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeOpenRedPacket,
		sdk.NewAttribute(types.AttributeSendeer, fromAddr.String()),                                          
		sdk.NewAttribute(types.AttributeSenderBalances, k.BankKeeper.GetAllBalances(ctx, fromAddr).String()), 
		sdk.NewAttribute(types.AttributeKeyAmount, luckyAmount.String()),                                     
		sdk.NewAttribute(types.AttributeKeyRedPacketId, redPacket.Id),                                        
		sdk.NewAttribute(types.AttributeKeyDenom, redPacket.Amount.Denom),                                    
		sdk.NewAttribute(types.AttributeKeyClusterChatId, clusterChatId),                                     
		sdk.NewAttribute(types.AttributeKeyRedPacketSerial, strconv.Itoa(len(redPacket.Receive))),            
		sdk.NewAttribute(types.AttributeKeyRedPacketRemain, redPacket.Remain().String()),                     
		sdk.NewAttribute(types.AttributeKeyRedPacketCountRemain, countRemain.String()),                       
		sdk.NewAttribute(types.AttributeKeyRedPacketType, strconv.FormatInt(redPacket.PacketType, 10)),       
	))

	return luckyAmount, nil
}

func (k Keeper) ReceiveNormalRedPacket(ctx sdk.Context, fromAddr sdk.AccAddress, redPacket *types.RedPacket, clusterChatId string) (sdkmath.Int, error) {
	logs := core.BuildLog(core.GetStructFuncName(k), core.LmChainDaoKeeper)

	
	redPacket.Receive = append(redPacket.Receive, types.RedPacketReceive{
		Receiver: fromAddr.String(),
		Amount:   redPacket.Amount.Amount,
	})

	err := k.SetRedPacketInfo(ctx, *redPacket)
	if err != nil {
		logs.WithError(err).Error("SetRedPacketInfo Error")
		return sdkmath.ZeroInt(), err
	}

	
	err = k.BankKeeper.SendCoinsFromModuleToAccount(
		ctx,
		types.ModuleName,
		fromAddr,
		sdk.NewCoins(redPacket.Amount),
	)
	if err != nil {
		logs.WithError(err).Error("Receive NormalRedPacket SendCoinsFromModuleToAccount Error")
		return sdkmath.ZeroInt(), err
	}

	fromBalances := k.BankKeeper.GetAllBalances(ctx, fromAddr)

	
	countRemain := sdk.NewInt(redPacket.Count - int64(len(redPacket.Receive)))

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeOpenRedPacket,
		sdk.NewAttribute(types.AttributeSendeer, fromAddr.String()),                                    
		sdk.NewAttribute(types.AttributeSenderBalances, fromBalances.String()),                         
		sdk.NewAttribute(types.AttributeKeyAmount, redPacket.Amount.Amount.String()),                   
		sdk.NewAttribute(types.AttributeKeyRedPacketId, redPacket.Id),                                  
		sdk.NewAttribute(types.AttributeKeyDenom, redPacket.Amount.Denom),                              
		sdk.NewAttribute(types.AttributeKeyClusterChatId, clusterChatId),                               
		sdk.NewAttribute(types.AttributeKeyRedPacketSerial, strconv.Itoa(len(redPacket.Receive))),      
		sdk.NewAttribute(types.AttributeKeyRedPacketRemain, redPacket.Remain().String()),               
		sdk.NewAttribute(types.AttributeKeyRedPacketCountRemain, countRemain.String()),                 
		sdk.NewAttribute(types.AttributeKeyRedPacketType, strconv.FormatInt(redPacket.PacketType, 10)), 
	))

	return redPacket.Amount.Amount, nil
}

func (k Keeper) ReturnRedPacketLogic(ctx sdk.Context, fromAddr sdk.AccAddress, redPacket *types.RedPacket) (sdkmath.Int, error) {
	logs := core.BuildLog(core.GetStructFuncName(k), core.LmChainDaoKeeper)

	
	returnCoin := sdk.NewCoin(redPacket.Amount.Denom, redPacket.Remain())

	err := k.BankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, fromAddr, sdk.NewCoins(returnCoin))
	if err != nil {
		logs.WithError(err).Error("ReturnRedPacketLogic SendCoinsFromModuleToAccount Error")
		return sdk.ZeroInt(), err
	}

	
	redPacket.IsReturn = true
	err = k.SetRedPacketInfo(ctx, *redPacket)
	if err != nil {
		return sdkmath.ZeroInt(), err
	}

	return redPacket.Remain(), nil
}
