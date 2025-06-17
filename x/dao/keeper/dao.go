package keeper

import (
	"cosmossdk.io/math"
	"freemasonry.cc/blockchain/core"
	"freemasonry.cc/blockchain/x/dao/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) ReceiveDao(ctx sdk.Context, cluster types.DeviceCluster, daoPoolAddr, fromAddr sdk.AccAddress, burnAmount sdk.Dec, reward math.Int) error {

	
	newV := cluster.ClusterPowerMembers[fromAddr.String()]

	newV.PowerCanReceiveDao = newV.PowerCanReceiveDao.Sub(burnAmount)

	
	cluster.ClusterPowerMembers[fromAddr.String()] = newV

	err := k.SetDeviceCluster(ctx, cluster)
	if err != nil {
		return err
	}

	
	err = k.BankKeeper.SendCoins(ctx, daoPoolAddr, fromAddr, sdk.NewCoins(sdk.NewCoin(core.BurnRewardFeeDenom, reward)))
	if err != nil {
		return err
	}

	
	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.GetDao,
		sdk.NewAttribute(types.AttributeReceiver, fromAddr.String()),
		sdk.NewAttribute(types.AttributeKeyAmount, reward.String()),
		sdk.NewAttribute(types.AttributeSendeer, daoPoolAddr.String()),
	))

	return nil
}
