package keeper

import (
	"freemasonry.cc/blockchain/core"
	"freemasonry.cc/blockchain/util"
	"freemasonry.cc/blockchain/x/gateway/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func NewQuerier(k Keeper, legacyQuerierCdc *codec.LegacyAmino) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		var (
			res []byte
			err error
		)
		switch path[0] {
		case types.QueryGatewayInfo: 
			return queryGatewayInfo(ctx, req, k, legacyQuerierCdc)
		case types.QueryGatewayList: 
			return queryGatewayList(ctx, k)
		case types.QueryGatewayNum: 
			return queryGatewayNum(ctx, k)
		case types.QueryGatewayRedeemNum: 
			return queryGatewayRedeemNum(ctx, k)
		case types.QueryValidatorByConsAddress: 
			return queryValidatorByConsAddress(ctx, req, k, legacyQuerierCdc)
		case types.QueryGatewayNumberCount:
			return queryGatewayNumberCount(ctx, req, k, legacyQuerierCdc)
		case types.QueryGatewayNumberUnbondCount:
			return queryGatewayNumberUnbondCount(ctx, req, k, legacyQuerierCdc)
		case types.QueryGasPrice:
			return queryGasPrice(ctx)
		case types.QueryParams:
			return queryParams(ctx, req, k, legacyQuerierCdc)
		case types.QueryValidators:
			return queryAllValidator(ctx, k, legacyQuerierCdc)
		case types.QueryGatewayUpload:
			return queryGatewayUpload(ctx, req, k, legacyQuerierCdc)
		case types.QueryDelegateLastHeightKey:
			return queryDelegateLastHeight(ctx, req, k, legacyQuerierCdc)
		default:
			err = sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unknown %s query endpoint: %s", types.ModuleName, path[0])
		}

		return res, err
	}
}

func queryDelegateLastHeight(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	log := core.BuildLog(core.GetPackageFuncName(), core.LmChainCommQuery)
	var params types.QueryDelegateLastHeight
	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		log.WithError(err).Error("UnmarshalJSON")
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	lastHeight, err := k.GetGatewayDelegateLastHeight(ctx, params.DelegateAddress, params.ValidatorAddress)
	if err != nil {
		log.WithError(err).Error("GetGatewayDelegateLastHeight")
		return nil, err
	}
	heightByte, err := util.Json.Marshal(lastHeight)
	if err != nil {
		log.WithError(err).Error("heightByte Marshal Err")
		return nil, err
	}
	return heightByte, nil
}

func queryGatewayUpload(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	log := core.BuildLog(core.GetPackageFuncName(), core.LmChainCommQuery)
	var params string
	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		log.WithError(err).Error("UnmarshalJSON")
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}
	valAddr, err := sdk.ValAddressFromBech32(params)
	if err != nil {
		log.WithError(err).Error("ValAddressFromBech32")
		return nil, err
	}

	
	resByte := k.GetGatewayUpload(ctx, valAddr.String())
	return resByte, nil
}


func queryParams(ctx sdk.Context, _ abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	log := core.BuildLog(core.GetPackageFuncName(), core.LmChainCommQuery)
	params := k.GetParams(ctx)

	res, err := codec.MarshalJSONIndent(legacyQuerierCdc, params)
	if err != nil {
		log.WithError(err).Error("MarshalJSONIndent")
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

func queryGasPrice(ctx sdk.Context) ([]byte, error) {
	log := core.BuildLog(core.GetPackageFuncName(), core.LmChainCommQuery)
	gasPrice := ctx.MinGasPrices()
	gasByte, err := util.Json.Marshal(gasPrice)
	if err != nil {
		log.WithError(err).Error("Marshal")
		return nil, err
	}
	return gasByte, nil
}


func queryGatewayNumberUnbondCount(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	log := core.BuildLog(core.GetPackageFuncName(), core.LmChainCommQuery)
	var params types.GatewayNumberCountParams
	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		log.WithError(err).Error("UnmarshalJSON")
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}
	valAddr, err := sdk.ValAddressFromBech32(params.GatewayAddress)
	if err != nil {
		log.WithError(err).Error("ValAddressFromBech32")
		return nil, err
	}
	delAddr := sdk.AccAddress(valAddr)
	
	delegation, found := k.stakingKeeper.GetDelegation(ctx, delAddr, valAddr)
	if !found {
		log.WithError(stakingTypes.ErrNoDelegation).Error("GetDelegation Err")
		return nil, stakingTypes.ErrNoDelegation
	}
	
	shares, err := k.stakingKeeper.ValidateUnbondAmount(
		ctx, delAddr, valAddr, params.Amount.Amount,
	)
	if err != nil {
		log.WithError(err).Error("ValidateUnbondAmount Err")
		return nil, err
	}
	param := k.GetParams(ctx)
	
	gateway, err := k.GetGatewayInfo(ctx, params.GatewayAddress)
	if err != nil {
		if err == core.ErrGatewayNotExist {
			return []byte("0"), nil
		}
		log.WithError(err).Error("GetGatewayInfo Err")
		return nil, err
	}
	
	balanceShares := delegation.Shares.Sub(shares)
	
	num := balanceShares.QuoInt(param.MinDelegate.Amount)

	
	hode := gateway.GatewayQuota - int64(len(gateway.GatewayNum))

	count := gateway.GatewayQuota - num.TruncateInt64() - hode
	if count < 0 {
		count = 0
	}
	countByte, err := util.Json.Marshal(count)
	if err != nil {
		log.WithError(err).Error("CountByteJson Marshal Err")
		return nil, err
	}
	return countByte, nil
}


func queryGatewayNumberCount(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	log := core.BuildLog(core.GetPackageFuncName(), core.LmChainCommQuery)
	var params types.GatewayNumberCountParams
	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		log.WithError(err).Error("LegacyQuerierCdc UnmarshalJSON Err")
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}
	param := k.GetParams(ctx)
	
	gateway, err := k.GetGatewayInfo(ctx, params.GatewayAddress)
	if err != nil && err != core.ErrGatewayNotExist {
		log.WithError(err).Error("GetGatewayInfo Err")
		return nil, err
	}
	num := params.Amount.Amount.Quo(param.MinDelegate.Amount)
	if err == core.ErrGatewayNotExist { 
		return num.Marshal()
	}
	
	hode := gateway.GatewayQuota - int64(len(gateway.GatewayNum))

	count := num.Int64() + hode

	countByte, err := util.Json.Marshal(count)
	if err != nil {
		log.WithError(err).Error("CountByte Json Marshal Err")
		return nil, err
	}
	return countByte, nil
}


func queryValidatorByConsAddress(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	log := core.BuildLog(core.GetPackageFuncName(), core.LmChainCommQuery)
	var params types.QueryValidatorByConsAddrParams

	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		log.WithError(err).Error("UnmarshalJSON")
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}
	
	validator, found := k.stakingKeeper.GetValidatorByConsAddr(ctx, params.ValidatorConsAddress)
	if !found {
		log.WithError(stakingTypes.ErrNoValidatorFound).Error("GetValidatorByConsAddr Err")
		return nil, stakingTypes.ErrNoValidatorFound
	}
	res, err := codec.MarshalJSONIndent(legacyQuerierCdc, validator)
	if err != nil {
		log.WithError(err).Error("MarshalJSONIndent validator Err")
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return res, nil
}


func queryGatewayList(ctx sdk.Context, k Keeper) ([]byte, error) {
	log := core.BuildLog(core.GetPackageFuncName(), core.LmChainCommQuery)
	gatewayList, err := k.GetGatewayList(ctx)
	if err != nil {
		log.WithError(err).Error("GetGatewayList")
		return nil, err
	}
	for i := 0; i < len(gatewayList); i++ {
		valAddr, err := sdk.ValAddressFromBech32(gatewayList[i].GatewayAddress)
		if err != nil {
			log.WithError(err).Error("ValAddressFromBech32")
			return nil, err
		}
		validator, found := k.stakingKeeper.GetValidator(ctx, valAddr)
		if found {
			gatewayList[i].GatewayName = validator.GetMoniker()
		}
	}
	gatewayListByte, err := util.Json.Marshal(gatewayList)
	if err != nil {
		log.WithError(err).Error("Marshal")
		return nil, err
	}
	return gatewayListByte, nil
}


func queryGatewayInfo(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	log := core.BuildLog(core.GetPackageFuncName(), core.LmChainCommQuery)
	var params types.QueryGatewayInfoParams
	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		log.WithError(err).Error("UnmarshalJSON")
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}
	var gateway *types.Gateway
	
	if params.GatewayAddress != "" {
		gateway, err = k.GetGatewayInfo(ctx, params.GatewayAddress)
		if err != nil {
			log.WithError(err).Error("GetGatewayInfo")
			return nil, err
		}
	}
	
	if params.GatewayNumIndex != "" {
		gateway, err = k.GetGatewayInfoByNum(ctx, params.GatewayNumIndex)
		if err != nil {
			log.WithError(err).Error("GetGatewayInfoByNum")
			return nil, err
		}
	}
	if gateway != nil {
		valAddr, err := sdk.ValAddressFromBech32(gateway.GatewayAddress)
		if err != nil {
			log.WithError(err).Error("ValAddressFromBech32")
			return nil, err
		}
		validator, found := k.stakingKeeper.GetValidator(ctx, valAddr)
		if found {
			gateway.GatewayName = validator.GetMoniker()
		}
		gatewayByte, err := util.Json.Marshal(gateway)
		if err != nil {
			log.WithError(err).Error("Marshal")
			return nil, err
		}
		return gatewayByte, nil
	}
	return nil, nil
}


func queryGatewayNum(ctx sdk.Context, k Keeper) ([]byte, error) {
	log := core.BuildLog(core.GetPackageFuncName(), core.LmChainCommQuery)
	gatewayNumMap, err := k.GetGatewayNumMap(ctx)
	if err != nil {
		log.WithError(err).Error("GetGatewayNumMap")
		return nil, err
	}
	gatewayMapByte, err := util.Json.Marshal(gatewayNumMap)
	if err != nil {
		log.WithError(err).Error("Marshal")
		return nil, err
	}
	return gatewayMapByte, nil
}


func queryGatewayRedeemNum(ctx sdk.Context, k Keeper) ([]byte, error) {
	log := core.BuildLog(core.GetPackageFuncName(), core.LmChainCommQuery)
	gatewayRedeemNumMap, err := k.GetGatewayRedeemNum(ctx)
	if err != nil {
		log.WithError(err).Error("GetGatewayRedeemNum")
		return nil, err
	}
	gatewayMapByte, err := util.Json.Marshal(gatewayRedeemNumMap)
	if err != nil {
		log.WithError(err).Error("Marshal")
		return nil, err
	}
	return gatewayMapByte, nil
}


func queryAllValidator(ctx sdk.Context, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	log := core.BuildLog(core.GetPackageFuncName(), core.LmChainCommQuery)
	validators := k.stakingKeeper.GetAllValidators(ctx)
	resp := []types.ValidatorInfo{}
	for _, val := range validators {
		validatorInfo := valInfo(val)
		consAddr, err := val.GetConsAddr()
		if err != nil {
			log.WithError(err).Error("GetConsAddr")
			continue
		}
		validatorInfo.ConsAddress = consAddr
		resp = append(resp, validatorInfo)
	}
	res, err := util.Json.Marshal(resp)
	if err != nil {
		log.WithError(err).Error("Marshal")
		return nil, err
	}
	return res, nil
}


func valInfo(val stakingTypes.Validator) types.ValidatorInfo {
	validatorInfo := types.ValidatorInfo{
		OperatorAddress:   val.OperatorAddress,
		Jailed:            val.Jailed,
		Status:            int(val.Status),
		Tokens:            val.Tokens,
		DelegatorShares:   val.DelegatorShares,
		Moniker:           val.Description.Moniker,
		Identity:          val.Description.Identity,
		Website:           val.Description.Website,
		SecurityContact:   val.Description.SecurityContact,
		Details:           val.Description.Details,
		UnbondingHeight:   val.UnbondingHeight,
		UnbondingTime:     val.UnbondingTime.Unix(),
		Rate:              val.Commission.Rate,
		MaxRate:           val.Commission.MaxRate,
		MaxChangeRate:     val.Commission.MaxChangeRate,
		MinSelfDelegation: val.MinSelfDelegation,
	}
	return validatorInfo
}
