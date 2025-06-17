package v6

import (
	"freemasonry.cc/blockchain/core"
	contracttypes "freemasonry.cc/blockchain/x/contract/types"
	daoTypes "freemasonry.cc/blockchain/x/dao/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	"github.com/evmos/ethermint/x/evm/types"
)

func CreateUpgradeHandler(mm *module.Manager,
	configurator module.Configurator, daoKeeper types.DaoKeeper, bankKeeper types.BankKeeper) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		logger := ctx.Logger().With("upgrade", UpgradeName)
		logger.Debug("running module migrations ...")

		
		idoSupply, err := daoKeeper.GetGenesisIdoSupply(ctx)
		if err != nil {
			return nil, err
		}

		subAmount := sdk.MustNewDecFromStr("32400000000000000000000000")

		idoSupply = idoSupply.Sub(subAmount)

		
		if idoSupply.IsZero() || idoSupply.IsNegative() {
			daoKeeper.DeleteGenesisIdoSupply(ctx)
			daoKeeper.UpdateGenesisIdoEndMark(ctx, true)
			
			daoKeeper.SetGenesisIdoEndTime(ctx, ctx.BlockTime().Unix())

			daoKeeper.SetGenesisIdoSupply(ctx, sdk.ZeroDec())

			err = daoKeeper.StartSwap(ctx)
			if err != nil {
				return nil, err
			}

			daoKeeper.SetGenesisIdoSupply(ctx, sdk.ZeroDec())

		} else {
			
			daoKeeper.SetGenesisIdoSupply(ctx, idoSupply)
		}

		
		err = bankKeeper.BurnCoins(ctx, daoTypes.ModuleName, sdk.NewCoins(sdk.NewCoin(core.BaseDenom, subAmount.Sub(sdk.NewDec(108000000)).TruncateInt())))
		if err != nil {
			return nil, err
		}

		
		daoKeeper.SetNotActivePowerAmount(ctx, sdk.MustNewDecFromStr("3600000000000000000000000"))

		
		
		err = bankKeeper.SendCoinsFromModuleToModule(ctx, contracttypes.GenesisIdoReward, daoTypes.ModuleName, sdk.NewCoins(sdk.NewCoin(core.BaseDenom, sdk.MustNewDecFromStr("19440000000000000000000000").TruncateInt())))
		err = bankKeeper.BurnCoins(ctx, daoTypes.ModuleName, sdk.NewCoins(sdk.NewCoin(core.BaseDenom, sdk.MustNewDecFromStr("19440000000000000000000000").TruncateInt())))
		if err != nil {
			return nil, err
		}

		return mm.RunMigrations(ctx, configurator, vm)
	}
}
