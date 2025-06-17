package contract

import (
	"freemasonry.cc/blockchain/contracts"
	"freemasonry.cc/blockchain/core"
	"freemasonry.cc/blockchain/x/contract/keeper"
	"freemasonry.cc/blockchain/x/contract/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	erc20types "github.com/evmos/evmos/v10/x/erc20/types"
)

func InitGenesis(
	ctx sdk.Context,
	k keeper.Keeper,
	accountKeeper authkeeper.AccountKeeper,
	data types.GenesisState,
) {
	k.SetParams(ctx, data.Params)
}

func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	return &types.GenesisState{
		Params: k.GetParams(ctx),
	}
}

func RegisterCoin(ctx sdk.Context, k keeper.Keeper) {
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
}

func initCoinMetadata(ctx sdk.Context, k keeper.Keeper) {
	coinMetadata := []banktypes.Metadata{
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		{
			Description: "usdt coin",
			Base:        core.UsdtDenom,
			DenomUnits: []*banktypes.DenomUnit{
				{
					Denom:    core.UsdtDenom,
					Exponent: 18,
				},
			},
			Display: core.UsdtDenom,
			Name:    "bsc bridge usdt",
			Symbol:  "usdt",
		},
	}
	for _, metadatum := range coinMetadata {
		k.BankKeeper.SetDenomMetaData(ctx, metadatum)
	}
}

func DeployUsdtContract(ctx sdk.Context, k keeper.Keeper) {
	log := core.BuildLog(core.GetPackageFuncName(), core.LmChainXContract)
	log.Info("usdt contract address:", UsdtContract)
	usdt := contracts.UsdtContract
	_, err := k.CallEVMWithData(ctx, erc20types.ModuleAddress, nil, usdt.Bin, true)
	if err != nil {
		log.WithError(err).Error("usdt contract deploy error")
		panic(err)
	}
}


func DeployTokenFactoryContract(ctx sdk.Context, k keeper.Keeper) {
	log := core.BuildLog(core.GetPackageFuncName(), core.LmChainXContract)
	log.Info("token factory contract address:", TokenFactoryContract)
	StakeFactoryMedal := contracts.AppTokenIssueJSONContract
	_, err := k.CallEVMWithData(ctx, types.ModuleAddress, nil, StakeFactoryMedal.Bin, true)
	if err != nil {
		log.WithError(err).Error("token factory contract deploy error")
		panic(err)
	}
	err = k.SetTokenFactoryContractAddress(ctx, TokenFactoryContract.String())

	conAddr := k.GetTokenFactoryContractAddress(ctx)
	log.Info("get token factory contract addr:", conAddr)

	if err != nil {
		panic(err)
	}
}

func DeployGenesisIdoContract(ctx sdk.Context, k keeper.Keeper) {
	log := core.BuildLog(core.GetPackageFuncName(), core.LmChainXContract)
	log.Info("GenesisIdo contract address:", GenesisIdoContract)
	genesisIdo := contracts.GenesisIdoNContract
	_, err := k.CallEVMWithData(ctx, types.ModuleAddress, nil, genesisIdo.Bin, true)
	if err != nil {
		log.WithError(err).Error("GenesisIdo contract deploy error")
		panic(err)
	}
}


func DeployAuthContract(ctx sdk.Context, k keeper.Keeper) {
	log := core.BuildLog(core.GetPackageFuncName(), core.LmChainXContract)
	log.Info("auth contract address:", AuthContract)
	auth := contracts.AuthJSONContract
	_, err := k.CallEVMWithData(ctx, types.ModuleAddress, nil, auth.Bin, true)
	if err != nil {
		log.WithError(err).Error("auth contract deploy error")
		panic(err)
	}
}


func DeployPriceContract(ctx sdk.Context, k keeper.Keeper) {
	log := core.BuildLog(core.GetPackageFuncName(), core.LmChainXContract)
	log.Info("price contract address:", PriceContract)
	price := contracts.PriceJSONContract
	_, err := k.CallEVMWithData(ctx, types.ModuleAddress, nil, price.Bin, true)
	if err != nil {
		log.WithError(err).Error("price contract deploy error")
		panic(err)
	}
}


func DeployRedPacketContract(ctx sdk.Context, k keeper.Keeper) {
	log := core.BuildLog(core.GetPackageFuncName(), core.LmChainXContract)
	err := k.SetRedPacketContractAddress(ctx, RedPacketContract.String())
	if err != nil {
		panic(err)
	}
	log.Info("red packet contract address:", RedPacketContract)
	redPacket := contracts.RedPacketContract
	_, err = k.CallEVMWithData(ctx, types.ModuleAddress, nil, redPacket.Bin, true)
	if err != nil {
		log.WithError(err).Error("red packet contract deploy error")
		panic(err)
	}
}

func DeploySwapSwitchContract(ctx sdk.Context, k keeper.Keeper) {
	log := core.BuildLog(core.GetPackageFuncName(), core.LmChainXContract)
	log.Info("SwapSwitch contract address:", SwapSwitchContract)
	swapSwitchPacket := contracts.SwapSwitchJSONContract
	_, err := k.CallEVMWithData(ctx, types.ModuleAddress, nil, swapSwitchPacket.Bin, true)
	if err != nil {
		log.WithError(err).Error("SwapSwitch contract deploy error")
		panic(err)
	}
}

func DeployWdstContract(ctx sdk.Context, k keeper.Keeper) {
	log := core.BuildLog(core.GetPackageFuncName(), core.LmChainXContract)
	log.Info("wdst contract address:", WdstContract)
	wdst := contracts.WdstJSONContract
	_, err := k.CallEVMWithData(ctx, types.ModuleAddress, nil, wdst.Bin, true)
	if err != nil {
		log.WithError(err).Error("wdst contract deploy error")
		panic(err)
	}
}


func DeployExchangeFactoryContract(ctx sdk.Context, k keeper.Keeper) {
	log := core.BuildLog(core.GetPackageFuncName(), core.LmChainXContract)
	log.Info("exchange factory contract address:", ExchangeFactoryContract)
	exchangeFactory := contracts.ExchangeFactoryContract

	_, err := k.CallEVMWithData(ctx, types.ModuleAddress, nil, exchangeFactory.Bin, true)
	if err != nil {
		log.WithError(err).Error("exchange factory contract deploy error")
		panic(err)
	}
}


func DeployExchangeRouterContract(ctx sdk.Context, k keeper.Keeper) {
	log := core.BuildLog(core.GetPackageFuncName(), core.LmChainXContract)
	log.Info("exchange router contract address:", ExchangeRouterContract)
	exchangeRouter := contracts.ExchangeRouterContract

	_, err := k.CallEVMWithData(ctx, types.ModuleAddress, nil, exchangeRouter.Bin, true)
	if err != nil {
		log.WithError(err).Error("exchange router contract deploy error")
		panic(err)
	}
}

func DeployMulticallContract(ctx sdk.Context, k keeper.Keeper) {
	log := core.BuildLog(core.GetPackageFuncName(), core.LmChainXContract)
	log.Info("Multicall contract address:", MulticallContract)
	multicall := contracts.MulticallContract
	_, err := k.CallEVMWithData(ctx, types.ModuleAddress, nil, multicall.Bin, true)
	if err != nil {
		log.WithError(err).Error("Multicall contract deploy error")
		panic(err)
	}
}
