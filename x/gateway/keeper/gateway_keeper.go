package keeper

import (
	"freemasonry.cc/blockchain/core"
	"freemasonry.cc/blockchain/util"
	"freemasonry.cc/blockchain/x/gateway/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	GatewayKey = "gateway_" 

	GatewayNumKey = "gateway_num" 

	GatewayRedeemNumKey = "gateway_redeem_num" 
)

func (k Keeper) SetGateway(ctx sdk.Context, msg types.MsgGatewayRegister, coin sdk.Coin, valAddress string) error {
	log := core.BuildLog(core.GetStructFuncName(k), core.LmChainCommKeeper)
	kvStore := k.KVHelper(ctx)
	params := k.GetParams(ctx)
	num := coin.Amount.Quo(params.MinDelegate.Amount)
	gatewayInfo := types.Gateway{}
	if kvStore.Has(GatewayKey + valAddress) {
		return core.ErrGatewayExist
	}
	gatewayInfo.GatewayAddress = valAddress
	gatewayInfo.GatewayUrl = msg.GatewayUrl
	gatewayInfo.GatewayQuota = num.Int64()
	gatewayInfo.Package = msg.Package
	gatewayInfo.PeerId = msg.PeerId
	gatewayInfo.MachineAddress = msg.MachineAddress

	if msg.IndexNumber == nil || len(msg.IndexNumber) == 0 { 
		return core.ErrGatewayNumLen
	}

	if (num.Sub(sdk.NewInt(int64(len(gatewayInfo.GatewayNum))))).LT(sdk.NewInt(int64(len(msg.IndexNumber)))) {
		return core.ErrGatewayNum
	}
	gatewayNumArray, err := k.GatewayNumFilter(ctx, gatewayInfo, msg.IndexNumber)
	if err != nil {
		log.WithError(err).Error("GatewayNumFilter")
		return err
	}
	
	gatewayNumArray[0].IsFirst = true
	gatewayInfo.GatewayNum = append(gatewayInfo.GatewayNum, gatewayNumArray...)
	
	err = k.SetGatewayNum(ctx, gatewayNumArray)
	if err != nil {
		log.WithError(err).Error("SetGatewayNum")
		return err
	}
	
	err = k.GatewayRedeemNumFilter(ctx, gatewayNumArray)
	if err != nil {
		log.WithError(err).Error("GatewayRedeemNumFilter")
		return err
	}
	
	if gatewayInfo.GatewayNum == nil || len(gatewayInfo.GatewayNum) == 0 {
		gatewayInfo.Status = 1
	}
	
	return kvStore.Set(GatewayKey+valAddress, gatewayInfo)
}


func (k Keeper) GatewayEdits(ctx sdk.Context, gatewayUrl, valAddress string) error {
	kvStore := k.KVHelper(ctx)
	gateway, err := k.GetGatewayInfo(ctx, valAddress)
	if err != nil {
		return err
	}
	gateway.GatewayUrl = gatewayUrl
	
	return kvStore.Set(GatewayKey+valAddress, *gateway)
}


func (k Keeper) GatewayNumFilter(ctx sdk.Context, gateway types.Gateway, indexNum []string) ([]types.GatewayNumIndex, error) {
	log := core.BuildLog(core.GetStructFuncName(k), core.LmChainCommKeeper)
	var gatewayNumArray []types.GatewayNumIndex
	for _, val := range indexNum {
		
		gatewayNum, isRegister, err := k.GetGatewayNum(ctx, val)
		if err != nil {
			log.WithError(err).Error("GetGatewayNum")
			return nil, err
		}
		if !isRegister && gateway.GatewayAddress != gatewayNum.GatewayAddress { 
			return nil, core.ErrGatewayNumber
		}
		
		if gateway.GatewayNum != nil {
			for _, num := range gateway.GatewayNum {
				if num.NumberIndex == val {
					return nil, core.ErrGatewayNumber
				}
			}
		}
		if gatewayNum == nil {
			gatewayNum = &types.GatewayNumIndex{
				NumberIndex: val,
			}
		}
		gatewayNum.GatewayAddress = gateway.GatewayAddress
		gatewayNum.Status = 0
		gatewayNum.Validity = 0
		gatewayNumArray = append(gatewayNumArray, *gatewayNum)
	}
	return gatewayNumArray, nil
}


func (k Keeper) UpdateGatewayInfo(ctx sdk.Context, gateway types.Gateway) error {
	log := core.BuildLog(core.GetStructFuncName(k), core.LmChainCommKeeper)
	kvStore := k.KVHelper(ctx)
	keys := GatewayKey + gateway.GatewayAddress
	err := kvStore.Set(keys, gateway)
	if err != nil {
		log.WithError(err).Error("Set")
		return err
	}
	return nil
}


func (k Keeper) SetGatewayNum(ctx sdk.Context, gatewayNumArray []types.GatewayNumIndex) error {
	log := core.BuildLog(core.GetStructFuncName(k), core.LmChainCommKeeper)
	kvStore := k.KVHelper(ctx)
	gatewayMap := make(map[string]types.GatewayNumIndex)
	if kvStore.Has(GatewayNumKey) {
		err := kvStore.GetUnmarshal(GatewayNumKey, &gatewayMap)
		if err != nil {
			log.WithError(err).Error("GetUnmarshal")
			return err
		}
	}
	for _, val := range gatewayNumArray {
		gatewayMap[val.NumberIndex] = val
	}
	err := kvStore.Set(GatewayNumKey, gatewayMap)
	if err != nil {
		log.WithError(err).Error("Set")
		return err
	}
	return nil
}


func (k Keeper) SetGatewayRedeemNum(ctx sdk.Context, gatewayNumArray []types.GatewayNumIndex) error {
	log := core.BuildLog(core.GetStructFuncName(k), core.LmChainCommKeeper)
	kvStore := k.KVHelper(ctx)
	gatewayMap := make(map[string]types.GatewayNumIndex)
	if kvStore.Has(GatewayRedeemNumKey) {
		err := kvStore.GetUnmarshal(GatewayRedeemNumKey, &gatewayMap)
		if err != nil {
			log.WithError(err).Error("GetUnmarshal")
			return err
		}
	}
	for _, val := range gatewayNumArray {
		gatewayMap[val.NumberIndex] = val
	}
	err := kvStore.Set(GatewayRedeemNumKey, gatewayMap)
	if err != nil {
		log.WithError(err).Error("Set")
		return err
	}
	return nil
}


func (k Keeper) GatewayRedeemNumFilter(ctx sdk.Context, gatewayNumArray []types.GatewayNumIndex) error {
	log := core.BuildLog(core.GetStructFuncName(k), core.LmChainCommKeeper)
	kvStore := k.KVHelper(ctx)
	if !kvStore.Has(GatewayRedeemNumKey) {
		return nil
	}
	gatewayNumMap := make(map[string]types.GatewayNumIndex)
	err := kvStore.GetUnmarshal(GatewayRedeemNumKey, &gatewayNumMap)
	if err != nil {
		log.WithError(err).Error("GetUnmarshal")
		return err
	}
	for _, val := range gatewayNumArray {
		if _, ok := gatewayNumMap[val.NumberIndex]; ok {
			delete(gatewayNumMap, val.NumberIndex)
		}
	}
	return kvStore.Set(GatewayRedeemNumKey, gatewayNumMap)
}


func (k Keeper) GetGatewayRedeemNum(ctx sdk.Context) (map[string]types.GatewayNumIndex, error) {
	log := core.BuildLog(core.GetStructFuncName(k), core.LmChainCommKeeper)
	kvStore := k.KVHelper(ctx)
	if !kvStore.Has(GatewayRedeemNumKey) {
		return nil, nil
	}
	gatewayNumMap := make(map[string]types.GatewayNumIndex)
	err := kvStore.GetUnmarshal(GatewayRedeemNumKey, &gatewayNumMap)
	if err != nil {
		log.WithError(err).Error("GetUnmarshal")
		return nil, core.ErrUnmarshal
	}
	return gatewayNumMap, nil
}


func (k Keeper) GetGatewayInfo(ctx sdk.Context, gatewayAddress string) (*types.Gateway, error) {
	log := core.BuildLog(core.GetStructFuncName(k), core.LmChainCommKeeper)
	kvStore := k.KVHelper(ctx)
	keys := GatewayKey + gatewayAddress
	if !kvStore.Has(keys) {
		return nil, core.ErrGatewayNotExist
	}
	gateway := new(types.Gateway)
	err := kvStore.GetUnmarshal(keys, gateway)
	if err != nil {
		log.WithError(err).Error("GetUnmarshal")
		return nil, core.ErrUnmarshal
	}
	return gateway, nil
}


func (k Keeper) GetGatewayList(ctx sdk.Context) ([]types.Gateway, error) {
	log := core.BuildLog(core.GetStructFuncName(k), core.LmChainCommKeeper)
	store := k.KVHelper(ctx)
	gatewayArray := []types.Gateway{}
	iterator := store.KVStorePrefixIterator(GatewayKey)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var gateway types.Gateway
		err := util.Json.Unmarshal(iterator.Value(), &gateway)
		if err != nil {
			log.WithError(err).Error("Unmarshal")
			return nil, err
		}
		if gateway.GatewayAddress == "" || gateway.Status == 1 {
			continue
		}
		gatewayArray = append(gatewayArray, gateway)
	}
	return gatewayArray, nil
}


func (k Keeper) GetGatewayInfoByNum(ctx sdk.Context, gatewayNum string) (*types.Gateway, error) {
	log := core.BuildLog(core.GetStructFuncName(k), core.LmChainCommKeeper)
	kvStore := k.KVHelper(ctx)
	if !kvStore.Has(GatewayNumKey) {
		return nil, nil
	}
	
	gatewayNumMap, err := k.GetGatewayNumMap(ctx)
	if err != nil {
		log.WithError(err).Error("GetGatewayNumMap")
		return nil, err
	}
	if gatewayNumMap == nil {
		return nil, nil
	}
	var gatewayAddress string
	if _, ok := gatewayNumMap[gatewayNum]; ok {
		gatewayAddress = gatewayNumMap[gatewayNum].GatewayAddress
	}
	if gatewayAddress == "" {
		return nil, nil
	}
	return k.GetGatewayInfo(ctx, gatewayAddress)
}


func (k Keeper) GetGatewayNumMap(ctx sdk.Context) (map[string]types.GatewayNumIndex, error) {
	log := core.BuildLog(core.GetStructFuncName(k), core.LmChainCommKeeper)
	kvStore := k.KVHelper(ctx)
	if !kvStore.Has(GatewayNumKey) {
		return nil, nil
	}
	gatewayNumMap := make(map[string]types.GatewayNumIndex)
	err := kvStore.GetUnmarshal(GatewayNumKey, &gatewayNumMap)
	if err != nil {
		log.WithError(err).Error("GetUnmarshal")
		return nil, err
	}
	return gatewayNumMap, nil
}


func (k Keeper) GetGatewayNum(ctx sdk.Context, gatewayNum string) (*types.GatewayNumIndex, bool, error) {
	log := core.BuildLog(core.GetStructFuncName(k), core.LmChainCommKeeper)
	kvStore := k.KVHelper(ctx)
	if !kvStore.Has(GatewayNumKey) {
		return nil, true, nil
	}
	gatewayNumMap := make(map[string]types.GatewayNumIndex)
	err := kvStore.GetUnmarshal(GatewayNumKey, &gatewayNumMap)
	if err != nil {
		log.WithError(err).Error("GetUnmarshal")
		return nil, false, err
	}
	if val, ok := gatewayNumMap[gatewayNum]; ok {
		if val.Status != 2 {
			return &val, false, nil
		}
		return &val, true, nil
	}
	return nil, true, nil
}


func (k Keeper) UpdateGatewayNum(ctx sdk.Context, gatewayNumArry []types.GatewayNumIndex) error {
	log := core.BuildLog(core.GetStructFuncName(k), core.LmChainCommKeeper)
	for _, gatewayNum := range gatewayNumArry {
		gateway, err := k.GetGatewayInfo(ctx, gatewayNum.GatewayAddress)
		if err != nil {
			log.WithError(err).Error("GetGatewayInfo")
			return err
		}
		for i, num := range gateway.GatewayNum {
			if num.GatewayAddress == gatewayNum.GatewayAddress && num.NumberIndex == gatewayNum.NumberIndex {
				gateway.GatewayNum[i] = gatewayNum
			}
		}
		err = k.UpdateGatewayInfo(ctx, *gateway)
		if err != nil {
			log.WithError(err).Error("UpdateGatewayInfo")
			return err
		}
	}
	return nil
}

func (k Keeper) GatewayNumUnbond(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress, indexNumber []string) error {
	log := core.BuildLog(core.GetStructFuncName(k), core.LmChainCommKeeper)
	if valAddr.String() == sdk.ValAddress(delAddr).String() {
		shares := sdk.ZeroDec()
		
		delegation, found := k.stakingKeeper.GetDelegation(ctx, delAddr, valAddr)
		if found {
			
			
			shares = delegation.Shares
		}
		
		gatewayInfo, err := k.GetGatewayInfo(ctx, valAddr.String())
		if err != nil {
			if err == core.ErrGatewayNotExist {
				return nil
			}
			log.WithError(err).Error("GetGatewayInfo")
			return err
		}
		params := k.GetParams(ctx)
		
		num := shares.QuoInt(params.MinDelegate.Amount)
		gatewayInfo.GatewayQuota = num.TruncateInt64()
		
		holdNum := int64(len(gatewayInfo.GatewayNum))
		
		if holdNum-int64(len(indexNumber)) > gatewayInfo.GatewayQuota {
			return core.ErrGatewayNum
		}
		if indexNumber != nil && len(indexNumber) > 0 {
			var indexNumArray []types.GatewayNumIndex
			for _, val := range indexNumber {
				
				indexNum, _, err := k.GetGatewayNum(ctx, val)
				if err != nil {
					log.WithError(err).Error("GetGatewayNum")
					return err
				}
				if indexNum == nil {
					return core.ErrGatewayNumNotFound
				}
				indexNum.Status = 1 
				indexNum.Validity = ctx.BlockHeight() + params.Validity
				indexNumArray = append(indexNumArray, *indexNum)
				
				for i, gatewayNum := range gatewayInfo.GatewayNum {
					if gatewayNum.NumberIndex == val {
						gatewayInfo.GatewayNum = append(gatewayInfo.GatewayNum[:i], gatewayInfo.GatewayNum[i+1:]...)
					}
				}
			}
			
			if len(gatewayInfo.GatewayNum) == 0 {
				gatewayInfo.Status = 1
			} else {
				
				if !gatewayInfo.GatewayNum[0].IsFirst {
					return core.ErrGatewayFirstNum
				}
			}
			
			err := k.UpdateGatewayInfo(ctx, *gatewayInfo)
			if err != nil {
				log.WithError(err).Error("UpdateGatewayInfo")
				return err
			}
			
			err = k.SetGatewayNum(ctx, indexNumArray)
			if err != nil {
				log.WithError(err).Error("SetGatewayNum")
				return err
			}
			
			err = k.SetGatewayRedeemNum(ctx, indexNumArray)
			if err != nil {
				log.WithError(err).Error("SetGatewayRedeemNum")
				return err
			}
		}
		
		err = k.UpdateGatewayInfo(ctx, *gatewayInfo)
		if err != nil {
			log.WithError(err).Error("UpdateGatewayInfo")
			return err
		}
	}
	return nil
}


func (k Keeper) GetAllGatewayInfo(ctx sdk.Context) ([]types.Gateway, error) {
	log := core.BuildLog(core.GetStructFuncName(k), core.LmChainCommKeeper)
	store := k.KVHelper(ctx)
	gatewayArray := []types.Gateway{}
	iterator := store.KVStorePrefixIterator(GatewayKey)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var gateway types.Gateway
		err := util.Json.Unmarshal(iterator.Value(), &gateway)
		if err != nil {
			log.WithError(err).Error("Unmarshal")
			return nil, err
		}
		if gateway.GatewayAddress == "" || gateway.Status == 1 {
			continue
		}
		gatewayArray = append(gatewayArray, gateway)
	}
	return gatewayArray, nil
}


func (k Keeper) GetAllGatewayNum(ctx sdk.Context) (map[string]types.GatewayNumIndex, error) {
	log := core.BuildLog(core.GetStructFuncName(k), core.LmChainCommKeeper)
	kvStore := k.KVHelper(ctx)
	if !kvStore.Has(GatewayNumKey) {
		return nil, nil
	}
	gatewayNumMap := make(map[string]types.GatewayNumIndex)
	err := kvStore.GetUnmarshal(GatewayNumKey, &gatewayNumMap)
	if err != nil {
		log.WithError(err).Error("GetUnmarshal")
		return nil, err
	}

	return gatewayNumMap, nil
}
