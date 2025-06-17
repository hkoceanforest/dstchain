package keeper

import (
	v2 "freemasonry.cc/blockchain/x/dao/migrations/v2"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
)

var _ module.MigrationHandler = Migrator{}.Migrate1to2

type Migrator struct {
	keeper Keeper
}

func NewMigrator(keeper Keeper) Migrator {
	return Migrator{
		keeper: keeper,
	}
}

func (m Migrator) Migrate1to2(ctx sdk.Context) error {
	return v2.UpdateParams(ctx, &m.keeper.paramstore)
}
