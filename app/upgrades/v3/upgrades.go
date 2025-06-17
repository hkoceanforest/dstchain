package v3

import (
	gatewaykeeper "freemasonry.cc/blockchain/x/gateway/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

func CreateUpgradeHandler(mm *module.Manager,
	configurator module.Configurator, k gatewaykeeper.Keeper) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		logger := ctx.Logger().With("upgrade", UpgradeName)
		logger.Debug("running module migrations ...")

		
		gatewayAddr := "dstvaloper139kenyywsj0kqkarnru6szduwdh3qxjpyns3cs"
		gatewayInfo, err := k.GetGatewayInfo(ctx, gatewayAddr)
		if err != nil {
			panic(err)
		}

		gatewayInfo.MachineAddress = "dst16dn8m38c9gjzseu69zucl0hqk70rucxp9vxgtr"
		gatewayInfo.PeerId = "12D3KooWKDZoVfTDffWfG9r2ZKSRtqTzp4bq43ea61z58jP4As8B"

		kvStore := k.KVHelper(ctx)
		err = kvStore.Set(gatewaykeeper.GatewayKey+gatewayAddr, *gatewayInfo)
		if err != nil {
			panic(err)
		}

		return mm.RunMigrations(ctx, configurator, vm)
	}
}
