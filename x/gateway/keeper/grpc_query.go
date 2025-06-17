package keeper

import (
	"context"
	"freemasonry.cc/blockchain/core"
	"freemasonry.cc/blockchain/x/gateway/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

type Querier struct {
	Keeper
}

func (k Querier) GatewayNumberUnbondCount(goCtx context.Context, params *types.QueryGatewayNumberUnbondCountParams) (*types.QueryGatewayNumberUnbondCountResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	log := core.BuildLog(core.GetPackageFuncName(), core.LmChainCommQuery)

	
	amountInt, ok := sdk.NewIntFromString(params.Amount)
	if !ok {
		log.Error("params parse int error")
		return nil, core.ErrDelegationCoin
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
		ctx, delAddr, valAddr, amountInt,
	)
	if err != nil {
		log.WithError(err).Error("ValidateUnbondAmount Err")
		return nil, err
	}
	param := k.GetParams(ctx)
	
	gateway, err := k.GetGatewayInfo(ctx, params.GatewayAddress)
	if err != nil {
		if err == core.ErrGatewayNotExist {
			return &types.QueryGatewayNumberUnbondCountResponse{
				Count: 0,
			}, nil
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

	return &types.QueryGatewayNumberUnbondCountResponse{
		Count: count,
	}, nil
}

var _ types.QueryServer = Querier{}

func (k Querier) Params(goCtx context.Context, p *types.QueryGatewayParams) (*types.QueryGatewayParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	params := k.GetParams(ctx)

	return &types.QueryGatewayParamsResponse{Params: params}, nil
}
