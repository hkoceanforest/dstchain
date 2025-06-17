package types

import (
	"github.com/cosmos/cosmos-sdk/types"
)

const (
	MsgTypeMobileTransfer = "chat/MsgTypeMobileTransfer"
	MsgTypeBurnGetMobile  = "chat/MsgTypeBurnGetMobile"
	MsgTypeSetChatInfo    = "chat/MsgTypeSetChatInfo"
)


type UserInfo struct {
	
	FromAddress string `json:"from_address" yaml:"from_address"`
	
	RegisterNodeAddress string `json:"register_node_address" yaml:"register_node_address"`
	
	NodeAddress string `json:"node_address" yaml:"node_address"`
	
	AddressBook string `json:"address_book" yaml:"address_book"`
	
	ChatBlacklist string `json:"chat_blacklist" yaml:"chat_blacklist"`
	
	ChatWhitelist string `json:"chat_whitelist" yaml:"chat_whitelist"`
	
	Mobile []string `json:"mobile" yaml:"mobile"`
	
	UpdateTime int64 `json:"update_time" yaml:"update_time"`
	
	ChatBlackEncList string `json:"chat_black_enc_list" yaml:"chat_black_enc_list"`
	
	ChatWhiteEncList string `json:"chat_white_enc_list" yaml:"chat_white_enc_list" `
}

type CustomInfo struct {
	Address     string `json:"address" yaml:"address"`
	CommAddress string `json:"comm_address" yaml:"comm_address"`
	Mobile      string `json:"mobile" yaml:"mobile"`
}

type AllUserInfo struct {
	UserInfo
	PledgeLevel         int64  `json:"pledge_level"`          
	GatewayProfixMobile string `json:"gateway_profix_mobile"` 
	IsExist             int64  `json:"is_exist"`              
}

type MortgageInfo struct {
	MortgageRemain     types.Coin           `json:"mortgage_remain"`      
	MortgageDevideInfo []MortgageDevideInfo `json:"mortgage_devide_info"` 
}

type MortgageDevideInfo struct {
	MortgageAddress string `json:"mortgage_address"` 
	MortgageAmount  string `json:"mortgage_amount"`  
	ShowBalance     bool   `json:"show_balance"`     
}

type LastReceiveLog struct {
	Height int64      `json:"height"`
	Value  types.Coin `json:"value"`
}

type MortgageAddLog struct {
	Height        int64      `json:"height"`
	MortgageValue types.Coin `json:"mortgage_value"`
}
