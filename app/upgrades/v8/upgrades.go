package v8

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	"github.com/evmos/ethermint/x/evm/types"
)

func CreateUpgradeHandler(mm *module.Manager,
	configurator module.Configurator, feeKeeper types.FeeMarketKeeper) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		logger := ctx.Logger().With("upgrade", UpgradeName)
		logger.Debug("running module migrations ...")

		feeParams := feeKeeper.GetParams(ctx)
		feeParams.MinGasPrice = sdk.MustNewDecFromStr("3000000000")
		feeKeeper.SetParams(ctx, feeParams)

		return mm.RunMigrations(ctx, configurator, vm)
	}
}
