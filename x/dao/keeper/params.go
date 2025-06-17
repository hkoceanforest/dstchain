package keeper

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"freemasonry.cc/blockchain/x/dao/types"
)

func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.paramstore.GetParamSet(ctx, &params)
	return params
}

func (k Keeper) GetParamsIfExists(ctx sdk.Context) (params types.Params) {
	k.paramstore.GetParamSetIfExists(ctx, &params)
	return params
}

func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramstore.SetParamSet(ctx, &params)
}

func (k Keeper) GetBurnLevelAmount(params types.Params, level int64) sdkmath.Int {
	levels := params.GetClusterLevels()
	for _, l := range levels {
		if l.Level == level {
			return l.DaoLimit
		}
	}
	return sdkmath.ZeroInt()
}
