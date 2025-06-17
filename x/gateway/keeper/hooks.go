package keeper

import (
	"freemasonry.cc/blockchain/x/gateway/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

func (k Keeper) AfterDelegationModified(ctx sdk.Context, addr sdk.AccAddress, valAddr sdk.ValAddress) error {
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	return nil
}

type Hooks struct {
	k Keeper
}

var _ types.StakingHooks = Hooks{}

var _ types.CommonHooks = Keeper{}

func (k Keeper) AfterCreateGateway(ctx sdk.Context, validator stakingTypes.Validator) {
	if k.hooks != nil {
		k.hooks.AfterCreateGateway(ctx, validator)
	}
}

func (k Keeper) Hooks() Hooks {
	return Hooks{k}
}


func (h Hooks) AfterDelegationModified(ctx sdk.Context, addr sdk.AccAddress, valAddr sdk.ValAddress) error {
	return h.k.AfterDelegationModified(ctx, addr, valAddr)
}

func (h Hooks) AfterValidatorBonded(_ sdk.Context, _ sdk.ConsAddress, _ sdk.ValAddress) error {
	return nil
}
func (h Hooks) AfterValidatorRemoved(_ sdk.Context, _ sdk.ConsAddress, _ sdk.ValAddress) error {
	return nil
}
func (h Hooks) AfterValidatorCreated(_ sdk.Context, _ sdk.ValAddress) error { return nil }
func (h Hooks) AfterValidatorBeginUnbonding(_ sdk.Context, _ sdk.ConsAddress, _ sdk.ValAddress) error {
	return nil
}
func (h Hooks) BeforeValidatorModified(_ sdk.Context, _ sdk.ValAddress) error { return nil }
func (h Hooks) BeforeDelegationCreated(_ sdk.Context, _ sdk.AccAddress, _ sdk.ValAddress) error {
	return nil
}
func (h Hooks) BeforeDelegationSharesModified(_ sdk.Context, _ sdk.AccAddress, _ sdk.ValAddress) error {
	return nil
}
func (h Hooks) BeforeValidatorSlashed(_ sdk.Context, _ sdk.ValAddress, _ sdk.Dec) error { return nil }
func (h Hooks) BeforeDelegationRemoved(_ sdk.Context, _ sdk.AccAddress, _ sdk.ValAddress) error {
	return nil
}
