package contract

import (
	"freemasonry.cc/blockchain/x/contract/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func BeginBlocker(ctx sdk.Context, k keeper.Keeper) {

}

func EndBlocker(ctx sdk.Context, k keeper.Keeper) []abci.ValidatorUpdate {
	if ctx.BlockHeight() == 1 {
		
		DeployUsdtContract(ctx, k)
		
		DeployTokenFactoryContract(ctx, k)
		
		DeployGenesisIdoContract(ctx, k)
		
		DeployAuthContract(ctx, k)
		
		DeployPriceContract(ctx, k)
		
		DeployRedPacketContract(ctx, k)
		
		DeploySwapSwitchContract(ctx, k)
		
		DeployWdstContract(ctx, k)
		
		DeployExchangeFactoryContract(ctx, k)
		
		DeployExchangeRouterContract(ctx, k)
		
		DeployMulticallContract(ctx, k)
	}
	return []abci.ValidatorUpdate{}
}
