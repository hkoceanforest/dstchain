package types

import sdk "github.com/cosmos/cosmos-sdk/types"

var _ BankHooks = &MultiBankHooks{}

type MultiBankHooks []BankHooks

func NewMultiBankHooks(hooks ...BankHooks) MultiBankHooks {
	return hooks
}

func (m MultiBankHooks) AfterSendCoins(ctx sdk.Context, fromAddr, toAddr sdk.AccAddress, amt sdk.Coins) error {
	for _, hook := range m {
		if err := hook.AfterSendCoins(ctx, fromAddr, toAddr, amt); err != nil {
			return err
		}
	}
	return nil
}
