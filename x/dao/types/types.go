package types

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"sort"
)

type DeviceCluster struct {
	ClusterId                      string                         `json:"cluster_id"`                         
	ClusterChatId                  string                         `json:"cluster_chat_id"`                    
	ClusterName                    string                         `json:"cluster_name"`                       
	ClusterOwner                   string                         `json:"cluster_owner"`                      
	ClusterGateway                 string                         `json:"cluster_gateway"`                    
	ClusterLeader                  string                         `json:"cluster_leader"`                     
	ClusterDeviceMembers           map[string]ClusterDeviceMember `json:"cluster_device_members"`             
	ClusterPowerMembers            map[string]ClusterPowerMember  `json:"cluster_power_members"`              
	ClusterPower                   sdk.Dec                        `json:"cluster_power"`                      
	ClusterLevel                   int64                          `json:"cluster_level"`                      
	ClusterBurnAmount              sdk.Dec                        `json:"cluster_burn_amount"`                
	ClusterActiveDevice            int64                          `json:"cluster_active_device"`              
	ClusterDaoPool                 string                         `json:"cluster_dao_pool"`                   
	ClusterRouteRewardPool         string                         `json:"cluster_route_reward_pool"`          
	ClusterDeviceRatio             sdk.Dec                        `json:"cluster_device_ratio"`               
	ClusterDeviceRatioUpdateHeight int64                          `json:"cluster_device_ratio_update_height"` 
	ClusterSalaryRatio             sdk.Dec                        `json:"cluster_salary_ratio"`               
	ClusterDvmRatio                sdk.Dec                        `json:"cluster_dvm_ratio"`                  
	ClusterDaoRatio                sdk.Dec                        `json:"cluster_dao_ratio"`                  
	ClusterSalaryRatioUpdateHeight ClusterChangeRatioHeight       `json:"cluster_salary_ratio_update_height"` 
	OnlineRatio                    sdk.Dec                        `json:"online_ratio"`                       
	OnlineRatioUpdateTime          int64                          `json:"online_ratio_update_time"`           
	ClusterAdminList               map[string]struct{}            `json:"cluster_admin_list"`                 
	ClusterVoteId                  uint64                         `json:"cluster_vote_id"`                    
	ClusterVotePolicy              string                         `json:"cluster_vote_policy"`                
}

type ClusterChangeRatioHeight struct {
	SalaryRatioUpdateHeight int64 `json:"salary_ratio_update_height"` 
	DvmRatioUpdateHeight    int64 `json:"dvm_ratio_update_height"`    
	DaoRatioUpdateHeight    int64 `json:"dao_ratio_update_height"`    
}

type PowerSupply struct {
	ActivePower sdk.Dec `json:"active_power"` 
}

type PersonalClusterInfo struct {
	Address string `json:"address"`
	
	Device map[string]struct{} `json:"device"`
	
	Owner map[string]struct{} `json:"owner"`
	
	BePower map[string]struct{} `json:"be_power"`
	
	AllBurn sdk.Dec `json:"all_burn"`
	
	ActivePower sdk.Dec `json:"active_power"`
	
	FreezePower sdk.Dec `json:"freeze_power"`
	
	FirstPowerCluster string `json:"first_power_cluster"`
}

type FreezeUsers map[string]sdk.Dec


type ApprovePower struct {
	ClusterId string `json:"cluster_id"`  
	Address   string `json:"address"`     
	IsDaoPool bool   `json:"is_dao_pool"` 
	EndBlock  int64  `json:"end_block"`
}

type ClusterCurApprove struct {
	ApproveAddress string `json:"approve_address"` 
	EndBlock       int64  `json:"end_block"`
}

type QueryNoDvmParams struct {
	ClusterId string   `json:"cluster_id"`
	Addrs     []string `json:"addrs"`
}

type ClusterDaoPoolFee struct {
	ClusterId string    `json:"cluster_id"`
	DaoPool   string    `json:"dao_pool"`
	Amount    sdk.Coins `json:"amount"`
}

type QueryClusterAdParams struct {
	StartTime string  `json:"start_time"` 
	EndTime   string  `json:"end_time"`   
	Rate      sdk.Dec `json:"rate"`       
	Level     string  `json:"level"`      
}

type QueryRedPacketInfoParams struct {
	RedPacketId string `json:"red_packet_id"` 
}

type QueryClusterAdResp struct {
	ClusterId []string `json:"cluster_id"`
	Amount    string   `json:"amount"`
	Quantity  int      `json:"quantity"` 
}

type HisClusterMemberReward struct {
	DeviceReward sdkmath.Int `json:"device_reward"` 
	HisReward    sdkmath.Int `json:"his_reward"`    
	Erc20Reward  sdk.Coins   `json:"erc20_reward"`  
}

type ClusterMemberDaoReward struct {
	Address    string      `json:"address"`
	BurnAmount sdkmath.Int `json:"burn_amount"` 
	Time       int64       `json:"time"`        
	DaoReward  sdkmath.Int `json:"dao_reward"`  
}

type RedPacketReceive struct {
	Receiver string      `json:"receiver"` 
	Amount   sdkmath.Int `json:"amount"`   
}

type RedPacket struct {
	EndBlock      int64              `json:"end_block"`       
	Id            string             `json:"id"`              
	Sender        string             `json:"sender"`          
	ClusterTrueId string             `json:"cluster_true_id"` 
	Amount        sdk.Coin           `json:"amount"`          
	Count         int64              `json:"count"`           
	Receive       []RedPacketReceive `json:"receive"`         
	PacketType    int64              `json:"packet_type"`     
	IsReturn      bool               `json:"is_return"`       
}

func (p RedPacket) Remain() sdkmath.Int {

	total := sdk.ZeroInt()
	received := sdk.ZeroInt()

	if p.PacketType == RedPacketTypeNormal {
		total = p.Amount.Amount.Mul(sdk.NewInt(p.Count))
	} else if p.PacketType == RedPacketTypeLucky {
		total = p.Amount.Amount
	} else {
		return sdk.ZeroInt()
	}

	for _, receive := range p.Receive {
		received = received.Add(receive.Amount)
	}

	if total.LTE(received) {
		return sdk.ZeroInt()
	} else {
		return total.Sub(received)
	}

}

func (p RedPacket) IsReceived(addr string) bool {
	for _, receive := range p.Receive {
		if receive.Receiver == addr {
			return true
		}
	}

	return false
}

type TokenInfo struct {
	Symbol   string `json:"symbol"`
	Contract string `json:"contract"`
}

type Erc20Swap struct {
	Contract string    `json:"contract"` 
	Symbol   string    `json:"symbol"`   
	Token0   TokenInfo `json:"token0"`
	Token1   TokenInfo `json:"token1"`
}

type FreezeInitInfo struct {
	Address    string `json:"address"`
	BaseAmount string `json:"base_amount"`
	USDTAmount string `json:"usdt_amount"`
}

type QueryDaoQueueParams struct {
	Page      int    `json:"page"`
	PageSize  int    `json:"page_size"`
	ClusterId string `json:"cluster_id"`
	Member    string `json:"member"`
}

type QueryDaoQueueResp struct {
	Total  int                      `json:"total"`
	Queues []ClusterMemberDaoReward `json:"queues"`
	Rate   sdk.Dec                  `json:"rate"`
}

type GatewayClusterExport struct {
	GatewayAddress string `json:"gateway_address"`
	Clusters       map[string]struct{}
}

type LpTokens []LpToken

type LpToken struct {
	LpToken    sdk.Coin  `json:"lp_token"`
	LpContract string    `json:"lp_contract"`
	Token0     TokenInfo `json:"token0"`
	Token1     TokenInfo `json:"token1"`
}

func (coins LpTokens) Add(lpTokenB ...LpToken) (coalesced LpTokens) {

	uniqCoins := make(map[string]LpTokens, len(coins)+len(lpTokenB))
	
	for _, cL := range []LpTokens{coins, lpTokenB} {
		for _, c := range cL {
			uniqCoins[c.LpToken.Denom] = append(uniqCoins[c.LpToken.Denom], c)
		}
	}

	for denom, cL := range uniqCoins {
		comboCoin := LpToken{LpToken: sdk.NewCoin(denom, sdkmath.NewInt(0))}
		for _, c := range cL {
			comboCoin.LpToken = comboCoin.LpToken.Add(c.LpToken)
			comboCoin.LpContract = c.LpContract
			comboCoin.Token0 = c.Token0
			comboCoin.Token1 = c.Token1
		}
		if !comboCoin.LpToken.IsZero() {
			coalesced = append(coalesced, comboCoin)
		}
	}
	return coalesced
}

type ClusterDeviceMember struct {
	Address     string  `json:"address"`
	ActivePower sdk.Dec `json:"active_power"` 
}

type ClusterPowerMember struct {
	Address            string  `json:"address"`
	ActivePower        sdk.Dec `json:"active_power"`          
	BurnAmount         sdk.Dec `json:"burn_amount"`           
	PowerCanReceiveDao sdk.Dec `json:"power_can_receive_dao"` 
}

type ClusterCreateTime struct {
	ClusterId  string `json:"cluster_id"`
	CreateTime int64  `json:"create_time"`
}

type PowerRewardCycleInfo struct {
	Address   string              `json:"address"`
	CycleInfo map[int64]CycleInfo `json:"cycle_info"`
}

type CycleInfo struct {
	Cycle             int64                       `json:"cycle"`               
	AllReward         sdk.Dec                     `json:"all_reward"`          
	CutPerReward      sdk.Int                     `json:"cut_per_reward"`      
	RemainReward      sdk.Int                     `json:"remain_reward"`       
	ReceiveTimes      int64                       `json:"receive_times"`       
	AllCutReward      sdk.Int                     `json:"all_cut_reward"`      
	StartTime         int64                       `json:"start_time"`          
	ClusterRewardList map[string]ClusterCutReward `json:"cluster_reward_list"` 
	Status            int64                       `json:"status"`              
}

type ClusterCutReward struct {
	ClusterId string  `json:"cluster_id"`
	AllReward sdk.Dec `json:"all_reward"` 
	CutReward sdk.Int `json:"cut_reward"` 
	IsReceive bool    `json:"is_receive"` 
}

type GatewaySign struct {
	Members      string `json:"members"`
	OnlineAmount int64  `json:"online_amount"`
	Seq          int64  `json:"seq"`
	MemberAdd    bool   `json:"member_add"`
}

func SortSliceMembers(members []string) {
	sort.Slice(members, func(i, j int) bool {
		return members[i] < members[j]
	})
}

type SettlementInfo struct {
	ClusterId string `json:"cluster_id"`
	Height    int64  `json:"height"`
	Address   string `json:"address"`
}
