package core

import (
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	ibcTransferTypes "github.com/cosmos/ibc-go/v5/modules/apps/transfer/types"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
)

var (
	
	ContractAddressFee = authtypes.NewModuleAddress(authtypes.FeeCollectorName)

	
	ContractAddressBank = authtypes.NewModuleAddress(bankTypes.ModuleName)

	
	ContractAddressDistribution = authtypes.NewModuleAddress(distrtypes.ModuleName)

	
	ContractAddressStakingBonded = authtypes.NewModuleAddress(stakingtypes.BondedPoolName)

	
	ContractAddressStakingNotBonded = authtypes.NewModuleAddress(stakingtypes.NotBondedPoolName)

	ContractAddressGov = authtypes.NewModuleAddress(govtypes.ModuleName)

	
	ContractAddressIbcTransfer = authtypes.NewModuleAddress(ibcTransferTypes.ModuleName)

	
	ContractGatewayBonus = authtypes.NewModuleAddress(GatewayBonusAddress)

	
	

	
	

	
	ContractEvm = authtypes.NewModuleAddress(evmtypes.ModuleName)

	
	ContractAddressDao = authtypes.NewModuleAddress("dao")

	Erc20Reward = authtypes.NewModuleAddress(Erc20RewardAccount)

	
	ContractAddressValidatorBonded = authtypes.NewModuleAddress(stakingtypes.BondedPoolName)

	
	ContractAddressValidatorNotBonded = authtypes.NewModuleAddress(stakingtypes.NotBondedPoolName)

	
	ContractAddressModule = authtypes.NewModuleAddress("contract")

	
	ContractAddressExchangeInit = authtypes.NewModuleAddress("genesis_ido_reward")
)


const CrossChainInManageAccount string = "dst1qk5qku7vxd4f2hqxrhw7d93sx2khckylrxf9ju"


const CrossChainAccount string = "dst1qk5qku7vxd4f2hqxrhw7d93sx2khckylrxf9ju"


const CrossChainFeeAccount string = "dst1p69txg6ctfdll44jccyd3wrryrwp98eae4el4v"


const CrossChainAutoDump string = "dst1556h2zeg04qgaw2vkhcgdf8wzduv3nd2cxc4qq"
