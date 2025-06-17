package types

import (
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/ethereum/go-ethereum/common"
)

const (
	
	ModuleName = "contract"
	
	StoreKey = ModuleName

	QuerierRoute = ModuleName

	
	RouterKey = ModuleName

	StakeContractDeploy = "stakeContractDeploy"

	TokenFactoryContractDeploy = "tokenFactoryContractDeploy"

	GenesisIdoReward = "genesis_ido_reward"
)

var ModuleAddress common.Address

func init() {
	ModuleAddress = common.BytesToAddress(authtypes.NewModuleAddress(ModuleName).Bytes())
}

const (
	KeyTokenFactoryContractAddress = "token_factory_contract_address"

	KeyRedPacketContractAddress = "red_packet_contract_address"
)
