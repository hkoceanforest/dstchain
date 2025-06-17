package contract

import (
	"freemasonry.cc/blockchain/x/contract/types"
	"github.com/ethereum/go-ethereum/crypto"
	erc20types "github.com/evmos/evmos/v10/x/erc20/types"
)

var (
	
	UsdtContract = crypto.CreateAddress(erc20types.ModuleAddress, 0)

	
	TokenFactoryContract = crypto.CreateAddress(types.ModuleAddress, 0)

	
	GenesisIdoContract = crypto.CreateAddress(types.ModuleAddress, 1)

	
	AuthContract = crypto.CreateAddress(types.ModuleAddress, 2)

	
	PriceContract = crypto.CreateAddress(types.ModuleAddress, 3)

	
	RedPacketContract = crypto.CreateAddress(types.ModuleAddress, 4)

	
	SwapSwitchContract = crypto.CreateAddress(types.ModuleAddress, 5)

	
	WdstContract = crypto.CreateAddress(types.ModuleAddress, 6)

	
	ExchangeFactoryContract = crypto.CreateAddress(types.ModuleAddress, 7)

	
	ExchangeRouterContract = crypto.CreateAddress(types.ModuleAddress, 8)

	
	MulticallContract = crypto.CreateAddress(types.ModuleAddress, 9)
)
