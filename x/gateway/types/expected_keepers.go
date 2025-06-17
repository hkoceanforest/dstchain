package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

type StakingHooks interface {
	AfterDelegationModified(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) error
}

type CommonHooks interface {
	AfterCreateGateway(ctx sdk.Context, validator stakingTypes.Validator)
}
