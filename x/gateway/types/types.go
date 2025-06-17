package types

import sdk "github.com/cosmos/cosmos-sdk/types"

const (
	MSG_SMART_CREATE_VALIDATOR   = "gateway/MsgCreateSmartValidator"
	MSG_GATEWAY_REGISTER         = "gateway/MsgGatewayRegister"
	MSG_GATEWAY_Edit             = "gateway/MsgGatewayEdit"
	MSG_GATEWAY_INDEX_NUM        = "gateway/MsgGatewayIndexNum"
	MSG_GATEWAY_UNDELEGATION     = "gateway/MsgGatewayUndelegation"
	MSG_GATEWAY_BEGIN_REDELEGATE = "gateway/MsgGatewayBeginRedelegate"
	MSG_GATEWAY_UPLOAD           = "gateway/MsgGatewayUpload"
	MSG_EMPTY_RESPONSE           = "gateway/MsgEmptyResponse"
)


type Gateway struct {
	
	GatewayAddress string `json:"gateway_address"`
	
	GatewayName string `json:"gateway_name"`
	
	GatewayUrl string `json:"gateway_url"`
	
	GatewayQuota int64 `json:"gateway_quota"`
	
	Status int64 `json:"status"`
	
	GatewayNum []GatewayNumIndex `json:"gateway_num"`
	
	Package string `json:"package"`
	
	PeerId string `json:"peer_id"`
	
	MachineAddress string `json:"machine_address"`
	
	MachineUpdateTime int64 `json:"machine_update_time"`
	
	ValAccAddress string `json:"val_acc_address"`
}

type GatewayListResp struct {
	Gateway
	Token  sdk.Int `json:"token"`
	Online int64   `json:"online"`
}


type GatewayNumIndex struct {
	
	GatewayAddress string `json:"gateway_address"`
	
	NumberIndex string `json:"number_index"`
	
	NumberEnd []string `json:"number_end"`
	
	Status int64 `json:"status"`
	
	Validity int64 `json:"validity"`
	
	IsFirst bool `json:"is_first"`
}


type ValidatorInfor struct {
	ValidatorConsAddr string `json:"validator_consaddr"` 
	ValidatorStatus   string `json:"validator_status"`   
	ValidatorPubAddr  string `json:"validator_pubaddr"`  
	ValidatorOperAddr string `json:"validator_operaddr"` 
	AccAddr           string `json:"acc_addr"`           
}

type GatewayNumberCountReq struct {
	GatewayAddress string `json:"gateway_address"` 
	Amount         string `json:"amount"`
}

type IsValidReq struct {
	Number string `json:"number"` 
}


type ValidatorInfo struct {
	OperatorAddress   string          `json:"operator_address"`
	ConsAddress       sdk.ConsAddress `json:"cons_address"`
	Jailed            bool            `json:"jailed"`           
	Status            int             `json:"status"`           
	Tokens            sdk.Int         `json:"tokens"`           
	DelegatorShares   sdk.Dec         `json:"delegator_shares"` 
	Moniker           string          `json:"moniker"`          
	Identity          string          `json:"identity"`         
	Website           string          `json:"website"`          
	SecurityContact   string          `json:"security_contact"` 
	Details           string          `json:"details"`
	UnbondingHeight   int64           `json:"unbonding_height"`
	UnbondingTime     int64           `json:"unbonding_time"`
	Rate              sdk.Dec         `json:"rate"`                
	MaxRate           sdk.Dec         `json:"max_rate"`            
	MaxChangeRate     sdk.Dec         `json:"max_change_rate"`     
	MinSelfDelegation sdk.Int         `json:"min_self_delegation"` 
}


type GatewayUploadData struct {
	AppSign []byte `json:"app_sign"`
	ValKey  []byte `json:"val_key"`
	NodeKey []byte `json:"node_key"`
}
