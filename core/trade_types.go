package core

import (
	"freemasonry.cc/trerr"
)



var (
	TradeTypeTransfer            = RegisterTranserType("transfer", "", "transfer accounts")
	TradeTypeDelegation          = RegisterTranserType("bonded", "NXN", "NXN mortgage")
	TradeTypeDelegationFee       = RegisterTranserType("bonded-fee", "NXN", "NXN mortgage service charge")
	TradeTypeUnbondDelegation    = RegisterTranserType("unbonded", "NXN", "NXN redemption")
	TradeTypeUnbondDelegationFee = RegisterTranserType("unbonded-fee", "NXN", "NXN redemption fee")
	TradeTypeFee                 = RegisterTranserType("fee", "", "Service charge expenditure")
	TradeTypeDelegationReward    = RegisterTranserType("delegate-reward", "NXN", "NXN mortgage reward")
	TradeTypeCommissionReward    = RegisterTranserType("commission-reward", "NXN", "NXN Commission reward")
	TradeTypeCommunityReward     = RegisterTranserType("community-reward", "", "Community reward")
	TradeTypeValidatorUnjail     = RegisterTranserType("unjail", "POS", "POS release")
	TradeTypeValidatorMinerBonus = RegisterTranserType("validator-miner-bonus", "POS", "POS incentive")

	TradeTypeCrossChainOut   = RegisterTranserType("cross-chain-out", "", "Cross Chain Out")
	TradeTypeCrossChainFee   = RegisterTranserType("cross-chain-fee", "", "Cross Chain Fee")
	TradeTypeCrossChainIn    = RegisterTranserType("cross-chain-in", "", "Cross Chain In")
	TradeTypeGatewayRegister = RegisterTranserType("gateway-register", "", "Gateway register")

	TradeTypeChatBonus          = RegisterTranserType("bonus", "nxn", "fm bonus")
	TradeTypeBurnReg            = RegisterTranserType("burn_trade_type_reg", "", "Register gateway destroy mining")
	TradeTypeBurnCur            = RegisterTranserType("burn_trade_type_cur", "", "Current gateway destroy mining")
	TradeTypeChatSendGiftDevide = RegisterTranserType("event_devide_send_gift", "", "Chat send gift share")

	TradeTypeChatWithDraw  = RegisterTranserType("chat_withdraw", "", "chat withDraw")
	TradeTypeChatUnpledge  = RegisterTranserType("chat_unpledge", "", "chat unpledge")
	TradeTypeBurnGetMobile = RegisterTranserType("burn_get_mobile", "", "Destroy to get the phone number")
	TradeTypeSendGift      = RegisterTranserType("send_gift", "", "send gift")
	TradeTypeChatSendGift  = RegisterTranserType("chat_send_gift", "", "chat send gift")

	TradeTypeValidatorCreate = RegisterTranserType("gateway-create", "", "Create validator")
	TradeTypeGatewayEdit     = RegisterTranserType("gateway_edit", "", "gateway edit")
	TradeTypeValReg          = RegisterTranserType("create_validator", "", "validator create")
	TradeTypeValEdit         = RegisterTranserType("edit_validator", "", "validator edit")

	TradeTypeBurnGetMedalGet  = RegisterTranserType("burn_get_medal_get", "", "burn_get")
	TradeTypeDposReward       = RegisterTranserType("dpos_withdraw_reward", "NXN", "NXN withdraw reward")
	TradeTypeDposCommision    = RegisterTranserType("dpos_commision_reward", "NXN", "NXN withdraw commision reward")
	TradeTypeBurnGetMedalBurn = RegisterTranserType("burn_get_medal_burn", "", "burn")
	TradeTypeVote             = RegisterTranserType("vote", "", "vote")
	TradeTypeUnjail           = RegisterTranserType("unjail", "", "unjail")
	TradeTypeDeposit          = RegisterTranserType("proposal_deposit", "", "proposal_deposit")
	TradeTypeSetChatInfo      = RegisterTranserType("set_chat_info", "", "set_chat_info")
	TradeTypePropsalRefund    = RegisterTranserType("proposal_deposite_refund", "", "proposal_deposite_refund")

	TradeTypeChatRegister = RegisterTranserType("register", "", "register")

	TradeTypeCreateCluster              = RegisterTranserType("create_cluster", "", "create cluster")
	TradeTypeAddMembers                 = RegisterTranserType("add_members", "", "add cluster members")
	TradeTypeDeleteMembers              = RegisterTranserType("delete_members", "", "delete cluster members")
	TradeTypeChangeName                 = RegisterTranserType("cluster_change_name", "", "cluster change name")
	TradeTypeClusterExit                = RegisterTranserType("cluster_exit", "", "cluster exit")
	TradeTypeWithdrawDeviceRewards      = RegisterTranserType("withdraw_device_rewards", "", "withdraw devic rewards")
	TradeTypeBurn                       = RegisterTranserType("burn_to_power", "", "burn of exchange power")
	TradeTypeChangeDeviceRatio          = RegisterTranserType("change_device_ratio", "", "modify device reward ratio")
	TradeTypeChangeSalaryRatio          = RegisterTranserType("change_salary_ratio", "", "modify salary ratio")
	TradeTypeChangeDvmRatio             = RegisterTranserType("change_dvm_ratio", "dvm", "modify dvm ratio")
	TradeTypeChangeClusterId            = RegisterTranserType("change_cluster_id", "ID", "modify cluster ID")
	TradeTypeGetSalary                  = RegisterTranserType("get_salary", "", "get salary")
	TradeTypeUpdateAdmin                = RegisterTranserType("update_admin", "", "update administrator")
	TradeTypeWithdrawSwapDpos           = RegisterTranserType("withdraw_swap_dpos", "SwapPos", "swap pos pledge redemption")
	TradeTypeClusterAd                  = RegisterTranserType("cluster_ad", "", "cluster advertisement")
	TradeTypeGetDao                     = RegisterTranserType("get_dao", "dao", "get dao")
	TradeTypeGetDaoQueue                = RegisterTranserType("get_dao_queue", "dao", "get dao from queue")
	TradeTypeGetDaoRoute                = RegisterTranserType("get_dao_route", "dao", "get dao from route")
	TradeTypeChangeClusterDaoRatio      = RegisterTranserType("change_cluster_dao_ratio", "dao", "Modify the dao reward ratio")
	TradeEVM                            = RegisterTranserType("ethereum_tx", "EVM", "EVM trading")
	TradeTypeAgreeJoinCluster           = RegisterTranserType("agree_join_cluster", "", "Agree to add a cluster")
	TradeTypeRedPacket                  = RegisterTranserType("red_packet", "", "Send red packet")
	TradeTypeOpenRedPacket              = RegisterTranserType("open_red_packet", "", "Open red packet")
	TradeTypeRedReturnPacket            = RegisterTranserType("return_red_packet", "", "return red packet")
	TradeTypeAgreeJoinClusterApply      = RegisterTranserType("agree_join_cluster_apply", "", "Agree to join a cluster application")
	TradeTypeGroupProposalDeposit       = RegisterTranserType("group_proposal_deposit", "", "Group vote storage deposit")
	TradeTypeGroupProposalDepositRefund = RegisterTranserType("group_proposal_deposit_refund", "", "Group vote storage deposit refund")
	TradeTypeStartCut                   = RegisterTranserType("start_power_reward_redeem", "", "Receive four years of halved redemption of assets")
	TradeTypeStopCut                    = RegisterTranserType("receive_cut_reward", "", "Receive four years of halved redemption of assets")
)

var tradeTypeText = make(map[string]string)
var tradeTypeTextEn = make(map[string]string)


func RegisterTranserType(key, value, enValue string) TranserType {
	tradeTypeTextEn[key] = enValue
	tradeTypeText[key] = value
	return TranserType(key)
}


func GetTranserTypeConfig() map[string]string {
	if trerr.Language == "EN" {
		return tradeTypeTextEn
	} else {
		return tradeTypeText
	}
}

type TranserType string

func (this TranserType) GetValue() string {
	if text, ok := tradeTypeText[string(this[:])]; ok {
		return text
	} else {
		return ""
	}
}

func (this TranserType) GetKey() string {
	return string(this[:])
}
