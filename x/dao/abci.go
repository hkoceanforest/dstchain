package dao

import (
	"fmt"
	"freemasonry.cc/blockchain/core"
	"freemasonry.cc/blockchain/x/dao/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func BeginBlocker(ctx sdk.Context, k keeper.Keeper) {

}

func EndBlocker(ctx sdk.Context, k keeper.Keeper) []abci.ValidatorUpdate {
	logs := core.BuildLog(core.GetStructFuncName(k), core.LmChainDaoKeeper)
	
	k.PeriodicRewards(ctx)
	
	k.IncrementHistoricalRewards(ctx)
	
	k.CutProduct(ctx)
	
	k.PeriodicDaoReward(ctx)
	if ctx.BlockHeight() == int64(1) {
		
		k.RecordStartTime(ctx, ctx.BlockTime().Unix())
	}

	
	if k.GetGenesisIdoEndMark(ctx) && !k.GetExchangeLiquidityAdded(ctx) {
		logs.Info("genesis ido end, add liquidity, create pair+++++++++++++++++++++++++++++++++++++++++++++++++")

		
		err := k.GrantAuthorizationUsdt(ctx)
		if err != nil {
			return nil
		}

		
		err = k.ExchangeAddLiquidity(ctx)
		if err != nil {
			return nil
		}

		k.SetExchangeLiquidityAdded(ctx)
	}

	
	if ctx.BlockHeight()%14400 == 1 {
		err := k.SetRewardRatioYear(ctx)
		if err != nil {
			fmt.Println("set reward ratio year error")
			panic(err)
		}
	}

	
	supply, err := k.GetGenesisIdoSupply(ctx)
	if err != nil {
		logs.Error("--------------GetGenesisIdoSupply error" + err.Error())
		return []abci.ValidatorUpdate{}
	}

	a, err := k.GetNotActivePowerAmount(ctx)
	if err != nil {
		logs.Error("--------------GetNotActivePowerAmount error" + err.Error())
		return []abci.ValidatorUpdate{}
	}

	logs.Info("--------------GetGenesisIdoSupply", supply.String())
	logs.Info("--------------GetNotFreePowerAmount", a.String())

	return []abci.ValidatorUpdate{}
}
