package gateway

import (
	"freemasonry.cc/blockchain/x/gateway/keeper"
	"freemasonry.cc/blockchain/x/gateway/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func InitGenesis(
	ctx sdk.Context,
	k keeper.Keeper,
	data types.GenesisState,
) {

	kvStore := k.KVHelper(ctx)

	k.SetParams(ctx, data.Params)

	
	if len(data.Gateways) > 0 {
		for _, g := range data.Gateways {
			gatewayIndexNum := make([]types.GatewayNumIndex, 0)
			for _, num := range g.GatewayNum {
				gatewayIndexNum = append(gatewayIndexNum, types.GatewayNumIndex{
					GatewayAddress: num.GatewayAddress,
					NumberIndex:    num.NumberIndex,
					NumberEnd:      num.NumberEnd,
					Status:         num.Status,
					Validity:       num.Validity,
					IsFirst:        num.IsFirst,
				})
			}

			gateway := types.Gateway{
				GatewayAddress:    g.GatewayAddress,
				GatewayName:       g.GatewayName,
				GatewayUrl:        g.GatewayUrl,
				GatewayQuota:      g.GatewayQuota,
				Status:            g.Status,
				GatewayNum:        gatewayIndexNum,
				Package:           g.Package,
				PeerId:            g.PeerId,
				MachineAddress:    g.MachineAddress,
				MachineUpdateTime: g.MachineUpdateTime,
				ValAccAddress:     g.ValAccAddress,
			}

			err := kvStore.Set(keeper.GatewayKey+gateway.GatewayAddress, gateway)
			if err != nil {
				panic(err)
			}
		}

	}

	gatewayNumIndexs := make(map[string]types.GatewayNumIndex)
	for _, g := range data.GatewayNumIndexs {
		gatewayNumIndexs[g.NumberIndex] = types.GatewayNumIndex{
			GatewayAddress: g.GatewayAddress,
			NumberIndex:    g.NumberIndex,
			NumberEnd:      g.NumberEnd,
			Status:         g.Status,
			Validity:       g.Validity,
			IsFirst:        g.IsFirst,
		}
	}
	err := kvStore.Set(keeper.GatewayNumKey, gatewayNumIndexs)
	if err != nil {
		panic(err)
	}

}

func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	
	
	

	
	gateways, err := k.GetAllGatewayInfo(ctx)
	if err != nil {
		panic(err)
	}

	gatewayExport := make([]types.GatewayExport, 0)

	for _, gateway := range gateways {
		gatewayNumExport := make([]types.GatewayNumIndexExport, 0)
		for _, gatewayNum := range gateway.GatewayNum {
			gatewayNumExport = append(gatewayNumExport, types.GatewayNumIndexExport{
				GatewayAddress: gatewayNum.GatewayAddress,
				NumberIndex:    gatewayNum.NumberIndex,
				NumberEnd:      gatewayNum.NumberEnd,
				Status:         gatewayNum.Status,
				Validity:       gatewayNum.Validity,
				IsFirst:        gatewayNum.IsFirst,
			})
		}

		gatewayExport = append(gatewayExport, types.GatewayExport{
			GatewayAddress:    gateway.GatewayAddress,
			GatewayName:       gateway.GatewayName,
			GatewayUrl:        gateway.GatewayUrl,
			GatewayQuota:      gateway.GatewayQuota,
			Status:            gateway.Status,
			GatewayNum:        gatewayNumExport,
			Package:           gateway.Package,
			PeerId:            gateway.PeerId,
			MachineAddress:    gateway.MachineAddress,
			MachineUpdateTime: gateway.MachineUpdateTime,
			ValAccAddress:     gateway.ValAccAddress,
		})
	}

	
	allNumberIndex, err := k.GetAllGatewayNum(ctx)
	if err != nil {
		panic(err)
	}

	numberInexExport := make(map[string]types.GatewayNumIndexExport)
	for _, gatewayNum := range allNumberIndex {
		numberInexExport[gatewayNum.NumberIndex] = types.GatewayNumIndexExport{
			GatewayAddress: gatewayNum.GatewayAddress,
			NumberIndex:    gatewayNum.NumberIndex,
			NumberEnd:      gatewayNum.NumberEnd,
			Status:         gatewayNum.Status,
			Validity:       gatewayNum.Validity,
			IsFirst:        gatewayNum.IsFirst,
		}
	}

	return &types.GenesisState{
		Params:           k.GetParams(ctx),
		Gateways:         gatewayExport,
		GatewayNumIndexs: numberInexExport,
	}
}
