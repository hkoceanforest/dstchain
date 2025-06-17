package keeper

import (
	"freemasonry.cc/blockchain/contracts"
	"freemasonry.cc/blockchain/core"
	"freemasonry.cc/blockchain/x/contract/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	erc20types "github.com/evmos/evmos/v10/x/erc20/types"
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
	DeployUsdtContract(ctx, m.keeper)
	return nil
}

func (m Migrator) Migrate2to3(ctx sdk.Context) error {
	DeployGenesisIdoContract(ctx, m.keeper)
	return nil
}

func (m Migrator) Migrate3to4(ctx sdk.Context) error {
	DeployGenesisIdoContract(ctx, m.keeper)
	return nil
}

func DeployUsdtContract(ctx sdk.Context, k Keeper) {
	log := core.BuildLog(core.GetPackageFuncName(), core.LmChainXContract)
	contractAddr := crypto.CreateAddress(erc20types.ModuleAddress, 0)
	log.Info("usdt contract address:", contractAddr)
	usdt := contracts.UsdtContract
	_, err := k.CallEVMWithData(ctx, erc20types.ModuleAddress, nil, usdt.Bin, true)
	if err != nil {
		log.WithError(err).Error("usdt contract deploy error")
		panic(err)
	}
}

func DeployGenesisIdoContract(ctx sdk.Context, k Keeper) {
	log := core.BuildLog(core.GetPackageFuncName(), core.LmChainXContract)
	contractAddress := authtypes.NewModuleAddress(types.TokenFactoryContractDeploy)
	from := common.BytesToAddress(contractAddress)
	contractAddr := crypto.CreateAddress(from, 1)
	log.Info("GenesisIdo contract address:", contractAddr)
	genesisIdo := contracts.GenesisIdoNContract
	_, err := k.CallEVMWithData(ctx, from, nil, genesisIdo.Bin, true)
	if err != nil {
		log.WithError(err).Error("GenesisIdo contract deploy error")
		panic(err)
	}
}
