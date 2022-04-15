package keeper

import (
	v2 "github.com/Gravity-Bridge/Gravity-Bridge/module/x/gravity/migrations/v2"
	v3 "github.com/Gravity-Bridge/Gravity-Bridge/module/x/gravity/migrations/v3"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Migrator is a struct for handling in-place store migrations.
type Migrator struct {
	keeper Keeper
}

// NewMigrator returns a new Migrator.
func NewMigrator(keeper Keeper) Migrator {
	return Migrator{keeper: keeper}
}

// Migrate1to2 migrates from consensus version 1 to 2.
func (m Migrator) Migrate1to2(ctx sdk.Context) error {
	ctx.Logger().Info("Mercury Upgrade: Enter Migrate1to2()")
	return v2.MigrateStore(ctx, m.keeper.storeKey, m.keeper.cdc)
}

// Migrate2to3 migrates from consensus version 2 to 3.
func (m Migrator) Migrate2to3(ctx sdk.Context) error {
	ctx.Logger().Info("v3 Upgrade: Enter Migrate2to3()")
	return v3.MigrateStore(ctx, m.keeper.storeKey, m.keeper.cdc)
}
