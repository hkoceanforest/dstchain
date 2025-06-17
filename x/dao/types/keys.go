package types

import (
	"encoding/binary"
	"github.com/cosmos/cosmos-sdk/types/address"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/ethereum/go-ethereum/common"
)

const (
	
	ModuleName = "dao"
	
	StoreKey = ModuleName

	QuerierRoute = ModuleName

	
	RouterKey = ModuleName

	RedPacketTypeNormal = int64(1)

	RedPacketTypeLucky = int64(2)
)


var (
	ModuleAddress common.Address
	
	BurnStartingInfoPrefix          = []byte{0x00} 
	ClusterHistoricalRewardsPrefix  = []byte{0x01} 
	ClusterCurrentRewardsPrefix     = []byte{0x02} 
	ClusterOutstandingRewardsPrefix = []byte{0x03} 
	
	DeviceStartingInfoPrefix       = []byte{0x04} 
	DeviceHistoricalRewardsPrefix  = []byte{0x05} 
	DeviceCurrentRewardsPrefix     = []byte{0x06} 
	DeviceOutstandingRewardsPrefix = []byte{0x07} 
	
	DeviceClusterKey = []byte{0x08}
	
	PersonClusterInfoKey = []byte{0x09}
	
	TotalBurnAmount = []byte{0x10}
	
	TotalPowerAmount = []byte{0x11}
	
	ClusterIdKey = []byte{0x12}
	
	ClusterForGateway = []byte{0x13}
	
	ClusterDeductionFee = []byte{0x14}
	
	ClusterEvmDeductionFee = []byte{0x15}

	
	ClusterApprovePowerInfo = []byte{0x16}

	
	ClusterApproveInfo = []byte{0x17}

	
	NotActivePowerPrefix = []byte{0x18}

	
	SendDeductionFee = []byte{0x19}

	
	ClusterPolicyPrefix = []byte{0x20}

	
	MintMarkPrefix = []byte{0x21}

	
	ClusterTimePrefix = []byte{0x23}

	
	ClusterMemberRewardPrefix = []byte{0x24}

	
	CutProductionPrefix = []byte{0x25}

	
	ChainStartTimePrefix = []byte{0x26}

	
	GenesisIdoSupplyPrefix = []byte{0x27}

	
	ClusterMemberDaoRewardPrefix = []byte{0x28}

	
	RedPacketPrefix = []byte{0x29}

	
	ClusterDaoRewardSumPrefix = []byte{0x30}

	
	RidClusterIdKey = []byte{0x31}

	
	Erc20RewardPrefix = []byte{0x32}

	
	Erc20SwapPrefix = []byte{0x33}

	
	RemainderPoolKey = []byte{0x34}

	
	PowerRewardCutInfo = []byte{0x35}

	
	StartTimeKey = []byte{0x36}

	
	BurnSupplyKey = []byte{0x37}

	
	MintSupplyKey = []byte{0x38}

	
	DstDelegateSupplyKey = []byte{0x39}

	
	ClusterDaoRewardQueueKey = []byte{0x40}

	
	DstPerPowerDayKey = []byte{0x41}

	
	ExchangeLiquidityAddedKey = []byte{0x42}

	
	UsdtAuthKey = []byte{0x43}
)

func GetUsdtAuthKey() []byte {
	return append(UsdtAuthKey, address.MustLengthPrefix([]byte("usdt_auth"))...)
}

func GeExchangeLiquidityAddedKey() []byte {
	return append(ExchangeLiquidityAddedKey, address.MustLengthPrefix([]byte("exchange_liquidity_added"))...)
}

func GetDstPerPowerDayKey() []byte {
	return append(DstPerPowerDayKey, address.MustLengthPrefix([]byte("dst_per_power_day"))...)
}

func GetClusterDaoRewardQueuePrefixKey() []byte {
	return append(ClusterDaoRewardQueueKey, address.MustLengthPrefix([]byte("cluster_dao_reward_queue"))...)
}

func GetDstDelegateSupplyPrefixKey() []byte {
	return append(DstDelegateSupplyKey, address.MustLengthPrefix([]byte("dst_delegate_supply"))...)
}
func GetMintSupplyPrefixKey() []byte {
	return append(MintSupplyKey, address.MustLengthPrefix([]byte("mint_supply"))...)
}

func GetBurnSupplyPrefixKey() []byte {
	return append(BurnSupplyKey, address.MustLengthPrefix([]byte("burn_supply"))...)
}

func GetStartTimeKey() []byte {
	return append(StartTimeKey, address.MustLengthPrefix([]byte("start_time"))...)
}

func GetPowerRewardCutInfoKey(addr string) []byte {
	return append(PowerRewardCutInfo, address.MustLengthPrefix([]byte(addr))...)
}

func init() {
	ModuleAddress = common.BytesToAddress(authtypes.NewModuleAddress(ModuleName).Bytes())
}

func GetErc20SwapKey(contract string) []byte {
	return append(Erc20SwapPrefix, address.MustLengthPrefix([]byte(contract))...)
}

func GetErc20RewardKey() []byte {
	return append(Erc20RewardPrefix, address.MustLengthPrefix([]byte("ERC20Reward"))...)
}

func GetRidClusterIdKey(rid string) []byte {
	return append(RidClusterIdKey, address.MustLengthPrefix([]byte(rid))...)
}

func GetClusterDaoRewardSumPrefixKey(clusterId string) []byte {
	return append(ClusterDaoRewardSumPrefix, address.MustLengthPrefix([]byte(clusterId))...)
}

func GetRedPacketKey(redPacketId string) []byte {
	return append(RedPacketPrefix, address.MustLengthPrefix([]byte(redPacketId))...)
}

func GetClusterMemberDaoReward(clusterId string) []byte {
	return append(ClusterMemberDaoRewardPrefix, address.MustLengthPrefix([]byte(clusterId))...)
}

func GetGenesisIdoSupplyPrefixKey() []byte {
	return append(GenesisIdoSupplyPrefix, address.MustLengthPrefix([]byte("GenesisIdoSupply"))...)
}

func GetNotActivePowerPrefixPrefixKey() []byte {
	return append(NotActivePowerPrefix, address.MustLengthPrefix([]byte("NotActivePowerPrefix"))...)
}

func GetChainStartTimePrefixKey() []byte {
	return append(ChainStartTimePrefix, address.MustLengthPrefix([]byte("start_time"))...)
}

func GetCutProductionPrefixKey() []byte {
	return append(CutProductionPrefix, address.MustLengthPrefix([]byte("cut_production"))...)
}

func GetClusterMemberRewardKey(clusterId, memberAddress string) []byte {
	return append(append(ClusterMemberRewardPrefix, address.MustLengthPrefix([]byte(clusterId))...), address.MustLengthPrefix([]byte(memberAddress))...)
}

func GetClusterTimeKey(clusterId string) []byte {
	return append(ClusterTimePrefix, address.MustLengthPrefix([]byte(clusterId))...)
}

func GetMintMarkPrefixKey(mark string) []byte {
	return append(MintMarkPrefix, address.MustLengthPrefix([]byte(mark))...)
}

func GetClusterEvmDeductionFeeKey(contractAddress, clusterId string) []byte {
	return append(append(ClusterEvmDeductionFee, address.MustLengthPrefix([]byte(contractAddress))...), address.MustLengthPrefix([]byte(clusterId))...)
}
func GetClusterApprovePowerInfoKey(contractAddress string) []byte {
	return append(ClusterApprovePowerInfo, address.MustLengthPrefix([]byte(contractAddress))...)
}
func GetClusterApproveInfoInfoKey(clusterIdAndMember string) []byte {
	return append(ClusterApproveInfo, address.MustLengthPrefix([]byte(clusterIdAndMember))...)
}
func GetClusterDeductionFeeKey(clusterId string) []byte {
	return append(ClusterDeductionFee, address.MustLengthPrefix([]byte(clusterId))...)
}
func GetSendDeductionFeeKey(addr string) []byte {
	return append(SendDeductionFee, address.MustLengthPrefix([]byte(addr))...)
}

func GetClusterForGatewayKey(gatewayAddress string) []byte {
	return append(ClusterForGateway, []byte(gatewayAddress)...)
}

func GetClusterIdKey(clusterChatId string) []byte {
	return append(ClusterIdKey, []byte(clusterChatId)...)
}

func GetDeviceClusterKey(clusterId string) []byte {
	return append(DeviceClusterKey, []byte(clusterId)...)
}
func GetClusterPolicyKey(policyAddress string) []byte {
	return append(ClusterPolicyPrefix, []byte(policyAddress)...)
}

func GetPersonClusterInfoKey(address string) []byte {
	return append(PersonClusterInfoKey, []byte(address)...)
}

func GetClusterHistoricalRewardsKey(v string, k uint64) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, k)
	return append(append(ClusterHistoricalRewardsPrefix, address.MustLengthPrefix([]byte(v))...), b...)
}

func GetClusterCurrentRewardsKey(v string) []byte {
	return append(ClusterCurrentRewardsPrefix, address.MustLengthPrefix([]byte(v))...)
}

func GetClusterOutstandingRewardsKey(clusterId string) []byte {
	return append(ClusterOutstandingRewardsPrefix, address.MustLengthPrefix([]byte(clusterId))...)
}

func GetBurnStartingInfoKey(clusterId, memberAddress string) []byte {
	return append(append(BurnStartingInfoPrefix, address.MustLengthPrefix([]byte(clusterId))...), address.MustLengthPrefix([]byte(memberAddress))...)
}

func GetDeviceHistoricalRewardsKey(v string, k uint64) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, k)
	return append(append(DeviceHistoricalRewardsPrefix, address.MustLengthPrefix([]byte(v))...), b...)
}

func GetDeviceCurrentRewardsKey(v string) []byte {
	return append(DeviceCurrentRewardsPrefix, address.MustLengthPrefix([]byte(v))...)
}

func GetDeviceOutstandingRewardsKey(clusterId string) []byte {
	return append(DeviceOutstandingRewardsPrefix, address.MustLengthPrefix([]byte(clusterId))...)
}

func GetDeviceStartingInfoKey(clusterId, memberAddress string) []byte {
	return append(append(DeviceStartingInfoPrefix, address.MustLengthPrefix([]byte(clusterId))...), address.MustLengthPrefix([]byte(memberAddress))...)
}
