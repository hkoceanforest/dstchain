package app

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	v82 "github.com/evmos/evmos/v10/app/upgrades/v8_2"
	"github.com/evmos/evmos/v10/types"
)



func (app *Evmos) ScheduleForkUpgrade(ctx sdk.Context) {
	
	if !types.IsMainnet(ctx.ChainID()) {
		return
	}

	upgradePlan := upgradetypes.Plan{
		Height: ctx.BlockHeight(),
	}

	
	switch ctx.BlockHeight() {
	case v82.MainnetUpgradeHeight:
		upgradePlan.Name = v82.UpgradeName
		upgradePlan.Info = v82.UpgradeInfo
	default:
		
		return
	}

	
	
	if err := app.UpgradeKeeper.ScheduleUpgrade(ctx, upgradePlan); err != nil {
		panic(
			fmt.Errorf(
				"failed to schedule upgrade %s during BeginBlock at height %d: %w",
				upgradePlan.Name, ctx.BlockHeight(), err,
			),
		)
	}
}
