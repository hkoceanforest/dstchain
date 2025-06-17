package keeper

import (
	"freemasonry.cc/blockchain/core"
	"freemasonry.cc/blockchain/x/dao/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

type Hooks struct {
	k Keeper
}

var _ bankTypes.BankHooks = Hooks{}

func (k Keeper) Hooks() Hooks { return Hooks{k} }

func (h Hooks) AfterSendCoins(ctx sdk.Context, fromAddr, toAddr sdk.AccAddress, amt sdk.Coins) error {
	return h.k.afterSendCoins(ctx, fromAddr, toAddr, amt)
}

func (k Keeper) afterSendCoins(ctx sdk.Context, fromAddr, toAddr sdk.AccAddress, amt sdk.Coins) error {
	logs := core.BuildLog(core.GetStructFuncName(k), core.LmChainDaoKeeper)
	daoModuleAcc := k.accountKeeper.GetModuleAddress(types.ModuleName)
	if toAddr.Equals(daoModuleAcc) {
		logs.Info("send coins to dao module account and burn coins ", "from:", fromAddr.String(), " to:", toAddr.String(), " amount", amt.String())
		err := k.BankKeeper.BurnCoins(ctx, types.ModuleName, amt)
		if err != nil {
			return err
		}
		for _, coin := range amt {
			err = k.AddBurnSupply(ctx, coin.Amount)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
