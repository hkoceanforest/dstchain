package types

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/group"
	"time"
)

const (
	
	QueryBurnLevels = "burn_levels"

	
	QueryBurnLevel = "burn_level"

	
	QueryGatewayClusters = "query_gateway_clusters"
	
	QueryPersonClusterInfo = "person_cluster_info"
	
	QueryClusterInfoById = "cluster_info_by_id"
	
	QueryClusterGasReward = "cluster_gas_reward"
	
	QueryClusterOwnerReward = "cluster_owner_reward"
	
	QueryClusterDeviceReward = "cluster_device_reward"
	
	QueryDeductionFee = "cluster_deduction_fee"
	
	QueryInClusters = "in_clusters"
	
	QueryClusterInfo = "cluster_info"
	
	QueryCluster = "cluster"
	
	QueryClusterPersonInfo = "query_cluster_person_info"
	
	QueryDvmList = "query_dvm_list"
	
	QueryDaoParams = "query_dao_params"
	
	QueryClusterPersonals = "query_cluster_personals"
	
	QueryClusterPersonalInfo = "query_cluster_personal_info"
	
	QueryClusterProposalVoter = "query_cluster_personal_voter"
	
	QueryClusterProposalVoters = "query_cluster_personal_voters"
	
	QueryGroupMembers = "query_group_members"
	
	QueryGroupInfo = "query_group_info"
	
	QueryClusterApproveInfo = "query_cluster_approve_info"
	
	QueryClusterRelation = "query_cluster_relation"
	
	QueryDaoStatistic = "query_dao_statistic"
	
	QueryNoDvm = "query_no_dvm"
	
	QueryNoDvmChat = "query_no_dvm_chat"
	
	QueryClusterProposalTallyResult = "query_proposal_tally_result"
	
	QueryClusterAd = "query_cluster_ad"
	
	QueryAllParams = "query_all_params"
	
	QueryRedPacketInfo = "query_red_packet"
	
	QueryRedPacketContractAddr = "query_red_packet_contract_addr"
	
	QueryNextAccountNumber = "query_next_account_number"
	
	QueryGroupDeposit = "query_group_deposit"
	
	QueryGroupVotesByVoter = "query_group_votes_by_voter"
	
	QueryDaoQueue = "query_dao_queue"
	
	QueryAllClusterReward = "query_all_cluster_reward"
	
	QueryCutRewards = "query_cut_rewards"
	
	QuerySupplyStatistics = "query_supply_statistics"

	QueryGasDeduction = "query_gas_deduction"
	
	QueryMiningStatus = "query_mining_status"

	QueryClusterReceiveDao = "query_cluster_receive_dao"
)

type GasDeductionParams struct {
	Msg []sdk.Msg `json:"msg"`
	Fee sdk.Coins `json:"fee"`
}

type QuerySupplyStatisticsResp struct {
	BurnSupply      math.Int `json:"burn_supply"`      
	NoMintSupply    math.Int `json:"no_mint_supply"`   
	CirculateSupply math.Int `json:"circulate_supply"` 
	SwapSupply      math.Int `json:"swap_supply"`      
}

type QueryClusterParams struct {
	ClusterId string `json:"cluster_id"`
}

type QueryClusterRewardParams struct {
	Member    string `json:"member"`
	ClusterId string `json:"cluster_id"`
}
type QueryClusterRewardResp struct {
	Reward LpTokens `json:"reward"`
	Rate   sdk.Dec  `json:"rate"`
}
type QueryAllClusterRewardResp struct {
	Reward map[string]sdk.Coin `json:"reward"`
	Rate   sdk.Dec             `json:"rate"`
}

func NewRewardParams(member, clusterId string) QueryClusterRewardParams {
	return QueryClusterRewardParams{member, clusterId}
}

type QueryDelegatorParams struct {
	DelegatorAddr sdk.AccAddress
}

func NewQueryDelegatorParams(delegatorAddr sdk.AccAddress) QueryDelegatorParams {
	return QueryDelegatorParams{
		DelegatorAddr: delegatorAddr,
	}
}

type QueryValidatorParams struct {
	ValidatorAddr sdk.ValAddress
	Page, Limit   int
}

func NewQueryValidatorParams(validatorAddr sdk.ValAddress, page, limit int) QueryValidatorParams {
	return QueryValidatorParams{
		ValidatorAddr: validatorAddr,
		Page:          page,
		Limit:         limit,
	}
}

type QueryRedelegationParams struct {
	DelegatorAddr    sdk.AccAddress
	SrcValidatorAddr sdk.ValAddress
	DstValidatorAddr sdk.ValAddress
}

func NewQueryRedelegationParams(delegatorAddr sdk.AccAddress, srcValidatorAddr, dstValidatorAddr sdk.ValAddress) QueryRedelegationParams {
	return QueryRedelegationParams{
		DelegatorAddr:    delegatorAddr,
		SrcValidatorAddr: srcValidatorAddr,
		DstValidatorAddr: dstValidatorAddr,
	}
}

type QueryValidatorsParams struct {
	Page, Limit int
	Status      string
}

func NewQueryValidatorsParams(page, limit int, status string) QueryValidatorsParams {
	return QueryValidatorsParams{page, limit, status}
}

type QueryValidatorByConsAddrParams struct {
	ValidatorConsAddress sdk.ConsAddress
}

type ValidatorInfor struct {
	ValidatorConsAddr string `json:"validator_consaddr"` 
	ValidatorStatus   string `json:"validator_status"`   
	ValidatorPubAddr  string `json:"validator_pubaddr"`  
	ValidatorOperAddr string `json:"validator_operaddr"` 
	AccAddress        string `json:"acc_address"`        
	ValidatorPubKey   string `json:"validator_pubkey"`   
}

type QueryBurnLevelsParams struct {
	Addresses []string
}

type QueryPersonClusterInfoRequest struct {
	From string `json:"from"`
}

type QueryGatewayClustersParams struct {
	GatewayAddress string `json:"gateway_address"`
}

type InClusters struct {
	ClusterId string
	IsOwner   bool
}



type ClusterInfo struct {
	
	ClusterId string `json:"cluster_id"`
	
	ClusterChatId string `json:"cluster_chat_id"`
	
	ClusterOwner string `json:"cluster_owner"`
	
	ClusterName string `json:"cluster_name"`
	
	ClusterAllBurn math.Int `json:"cluster_all_burn"`
	
	ClusterAllPower math.Int `json:"cluster_all_power"`
	
	OnlineRatio string `json:"online_ratio"` 
	
	ClusterActiveDevice int64 `json:"cluster_active_device"` 
	
	ClusterDeviceAmount int64 `json:"cluster_device_amount"` 
	
	DeviceConnectivityRate string `json:"device_connectivity_rate"` 
	
	ClusterDeviceRatio string `json:"cluster_device_ratio"` 
	
	ClusterSalaryRatio string `json:"cluster_salary_ratio"` 
	
	ClusterDvmRatio string `json:"cluster_dvm_ratio"`
	
	ClusterDaoRatio string `json:"cluster_dao_ratio"`
	
	ClusterDayFreeGas math.Int `json:"cluster_day_free_gas"`
	
	ClusterDaoPoolPower math.Int `json:"cluster_dao_pool_power"`
	
	DaoPoolDayFreeGas math.Int `json:"dao_pool_day_free_gas"`
	
	DaoPoolAvailableAmount math.Int `json:"dao_pool_available_amount"`
	
	DaoLicensingContract string `json:"dao_licensing_contract"`
	
	DaoLicensingHeight int64 `json:"dao_licensing_height"`
	
	LevelInfo ClusterLevelInfo `json:"level_info"`
	
	GatewayAddress string `json:"gateway_address"`
	
	IndexNum string `json:"index_num"`
	
	ClusterVotePolicy string `json:"cluster_vote_policy"`
	
	DaoPoolBalance math.Int `json:"dao_pool_balance"`
	
	RoutePoolBalance math.Int `json:"route_pool_balance"` 
	
	ClusterDaoPoolReward math.Int `json:"cluster_dao_pool_reward"`
	
	ClusterDaoPool string `json:"cluster_dao_pool"`
}

type ClusterLevelInfo struct {
	
	Level int64 `json:"level"`
	
	BurnAmountNextLevel math.Int `json:"burn_amount_next_level"`
	
	ActiveAmountNextLevel int64 `json:"active_amount_next_level"`
}

type ClusterPersonalInfo struct {
	PowerAmount       math.Int `json:"power_amount"`         
	GasDay            math.Int `json:"gas_day"`              
	BurnAmount        math.Int `json:"burn_amount"`          
	IsDevice          bool     `json:"is_device"`            
	IsAdmin           bool     `json:"is_amind"`             
	IsOwner           bool     `json:"is_owner"`             
	BurnRatio         string   `json:"burn_ratio"`           
	PowerReward       math.Int `json:"power_reward"`         
	DeviceReward      math.Int `json:"device_reward"`        
	OwnerReward       math.Int `json:"owner_reward"`         
	AuthContract      string   `json:"auth_contract"`        
	AuthHeight        int64    `json:"auth_height"`          
	ClusterOwner      string   `json:"cluster_owner"`        
	ClusterName       string   `json:"cluster_name"`         
	ClusterId         string   `json:"cluster_id"`           
	ClusterChatId     string   `json:"cluster_chat_id"`      
	DaoReceive        math.Int `json:"dao_receive"`          
	BurnRewardFeeRate string   `json:"burn_reward_fee_rate"` 
}

type QueryClusterPersonalInfoParams struct {
	ClusterId   string `json:"cluster_id"`
	FromAddress string `json:"from_address"`
}

type DvmInfo struct {
	
	ClusterChatId string `json:"cluster_chat_id"`
	
	ClusterId string `json:"cluster_id"`
	
	PowerReward math.Int `json:"power_reward"`
	
	PowerDvm math.Int `json:"power_dvm"`
	
	GasDayDvm math.Int `json:"gas_day_dvm"`
	
	AuthContract string `json:"auth_contract"`
	
	AuthHeight int64 `json:"auth_height"`
	
	ClusterName string `json:"cluster_name"`
	
	IsOwner bool `json:"is_owner"`
}

type PersonClusterStatisticsInfo struct {
	Address string `json:"address"`
	
	Owner []string `json:"owner"`
	
	BePower []string `json:"be_power"`
	
	AllBurn math.Int `json:"all_burn"`
	
	ActivePower math.Int `json:"active_power"`
	
	FreezePower math.Int `json:"freeze_power"`
	
	DeviceInfo []DeviceInfo `json:"device"`
	
	FirstPowerCluster string `json:"first_power_cluster"`
}

type DeviceInfo struct {
	ClusterChatId     string `json:"cluster_chat_id"`
	ClusterName       string `json:"cluster_name"`
	ClusterLevel      int64  `json:"cluster_level"`
	ClusterOwner      string `json:"cluster_owner"`
	ClusterCreateTime int64  `json:"cluster_create_time"`
}

type DaoParams struct {
	
	BurnGetPowerRatio string `json:"burn_get_power_ratio"`

	
	SalaryRange Range `json:"salary_range"`

	
	DeviceRange Range `json:"device_range"`

	
	DvmRange Range `json:"dvm_range"`

	
	CreateClusterMinBurn math.Int `json:"create_cluster_min_burn"`

	
	BurnAddress string `json:"burn_address"`

	
	DayBurnReward string `json:"day_burn_reward"`

	
	DaoRange Range `json:"dao_range"`

	
	ReceiveDaoRatio string `json:"receive_dao_ratio"`

	
	BurnRewardFeeRate string `json:"burn_reward_fee_rate"`

	
	SevenDayYearRate string `json:"seven_day_year_rate"`

	
	TotalBurn string `json:"total_burn"`

	
	TotalPower string `json:"total_power"`
}


type SevenDayYearRate struct {
	Show bool   `json:"show"`
	Rate string `json:"rate"`
}

type Range struct {
	Max string `json:"max"`
	Min string `json:"min"`
}
type QueryClusterProposalVoterParams struct {
	ProposalId string `json:"proposal_id"`
	Voter      string `json:"voter"`
}

type QueryClusterApproveParams struct {
	ClusterId string `json:"cluster_id"`
	Address   string `json:"address"`
}

type QueryClusterRelationResp struct {
	Address  string `json:"address"`
	IndexNum string `json:"indexNum"`
}

type QueryDaoStatisticResp struct {
	ClusterIds []string       `json:"cluster_ids"`
	Statistic  []DaoStatistic `json:"statistic"`
}

type DaoStatistic struct {
	ClusterId     string `json:"cluster_id"`
	ClusterChatId string `json:"cluster_chat_id"`

	
	DeviceOnlienRatio sdk.Dec `json:"device_onlien_ratio"`
	
	DeviceConnectivityRate sdk.Dec `json:"device_connectivity_rate"`

	
	DeviceAmount int64 `json:"device_amount"`
	
	DvmAmount int64 `json:"dvm_amount"`

	
	DeviceNoEvm int64 `json:"device_no_evm"`

	
	ClusterBurnAmount sdk.Dec `json:"cluster_burn_amount"`
}

type QueryProposalVotersParams struct {
	ProposalId string `json:"proposal_id"`
	ClusterId  string `json:"cluster_id"`
}

type Vote struct {
	ProposalId uint64           `json:"proposal_id"`
	SubmitTime time.Time        `json:"submit_time"`
	Voter      string           `json:"voter"`
	Option     group.VoteOption `json:"option"`
	Metadata   string           `json:"metadata"`
	Weight     string           `json:"weight"`
}

type QueryProposalVotersResp struct {
	Votes []Vote `json:"votes"`
}

type QueryClusterProposalsParams struct {
	ClusterId string `json:"cluster_id"`
	Voter     string `json:"voter"`
}

type ClusterProposals struct {
	group.ProposalNew
	IsVoter bool `json:"is_voter"`
}

type GetCutRewardInfoResp struct {
	
	IsStart bool `json:"is_start"`

	
	RemainDays int64 `json:"remain_days"`

	
	NextStartTime int64 `json:"next_start_time"`

	
	RedemptionAmount sdk.Int `json:"redemption_amount"`

	
	ReceiveAmount sdk.Int `json:"receive_amount"`
}

type QueryMiningStatusResp struct {
	
	MintStart bool `json:"mint_start"`
	
	IdoEndMark bool `json:"ido_end_mark"`
}
