package v12

import (
	"freemasonry.cc/blockchain/x/dao/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

func CreateUpgradeHandler(mm *module.Manager,
	configurator module.Configurator, daoKeeper keeper.Keeper) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		logger := ctx.Logger().With("upgrade", UpgradeName)
		logger.Debug("running module migrations ...")

		daoParams := daoKeeper.GetParams(ctx)
		daoParams.ConnectivityDaoRatio = sdk.MustNewDecFromStr("0.6")
		daoKeeper.SetParams(ctx, daoParams)

		logger.Info("dao params updated")

		return mm.RunMigrations(ctx, configurator, vm)
	}
}
