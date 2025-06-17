package keeper

import (
	"freemasonry.cc/blockchain/core"
	"freemasonry.cc/blockchain/x/dao/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) InitGenesis(ctx sdk.Context, data *types.GenesisState) {

	k.SetNotActivePowerAmount(ctx, core.GenesisIdoSupply)

	k.UpdateGenesisIdoEndMark(ctx, false)

	k.SetCutProductTime(ctx, ctx.BlockTime().Unix())

	k.SetGenesisIdoSupply(ctx, core.GenesisIdoSupply)

	k.accountKeeper.GetModuleAccountAndPermissions(ctx, types.ModuleName)
}

func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {

	return &types.GenesisState{}
}
