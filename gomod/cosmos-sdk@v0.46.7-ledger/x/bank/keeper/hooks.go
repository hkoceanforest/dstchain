package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank/types"
)

var _ types.BankHooks = BaseSendKeeper{}

func (k BaseSendKeeper) AfterSendCoins(ctx sdk.Context, fromAddr, toAddr sdk.AccAddress, amt sdk.Coins) error {
	if k.hooks != nil {
		return k.hooks.AfterSendCoins(ctx, fromAddr, toAddr, amt)
	}
	return nil
}
