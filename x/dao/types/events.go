package types

const (
	
	AttributeSendeer           = "sender"
	AttributeReceiver          = "receiver"
	AttributeSenderBalances    = "sender_balances"
	AttributeKeyClusterChatId  = "cluster_chat_id"
	AttributeKeyClusterTrueId  = "cluster_true_id"
	AttributeKeyClusterOwner   = "cluster_owner"
	AttributeKeyClusterName    = "cluster_name"
	AttributeKeyFeeSingle      = "fee_payer"
	AttributeKeyMemo           = "memo"
	AttributeKeySenderBalances = "sender_balances"
	AttributeKeyDaoModule      = "dao_module"
	AttributeKeyDistriModule   = "distribution_module"
	AttributeKeyFeeCollecter   = "fee_collecter"
	AttributeKeyAddress        = "address"
	AttributeKeyAmount         = "amount"
	AttributeKeyDenom          = "denom"
	AttributeKeyRate           = "rate"
	AttributeKeyBlockHeight    = "block_height"

	
	EventTypeOpenRedPacket   = "open_red_packet"
	EventTypeRedPacket       = "red_packet"
	EventTypeReturnRedPacket = "return_red_packet"

	AttributeKeyRedPacketId          = "red_packet_id"
	AttributeKeyRedPacketCount       = "red_packet_count"
	AttributeKeyRedPacketType        = "red_packet_type"
	AttributeKeyRedPacketSerial      = "red_packet_serial"
	AttributeKeyRedPacketRemain      = "red_packet_remain"
	AttributeKeyRedPacketCountRemain = "red_packet_count_remain"
	AttributeKeyRedPacketEndBlock    = "red_packet_end_block"

	
	EventTypeAgreeJoinCluster = "agree_join_cluster"

	
	EventTypeAgreeJoinClusterApply = "agree_join_cluster_apply"
	AttributeKeyMemberAddr         = "member_addr"
	AttributeKeyIndexNum           = "index_num"

	
	EventTypeCreateCluster = "create_cluster"

	
	EventTypeClusterUpgrade = "cluster_upgrade"
	AttributeKeyOldLevel    = "cluster_old_level"
	AttributeKeyNewLevel    = "cluster_new_level"

	
	EventTypeDeleteMembers = "cluster_delete_members"

	
	EventTypeChangeName = "cluster_change_name"

	
	EventTypeClusterExit = "cluster_exit"

	
	EventTypeWithdrawSwapDpos           = "withdraw_swap_dpos"
	EventTypeWithdrawDeviceRewards      = "withdraw_device_rewards"
	EventTypeIncrementHistoricalRewards = "increment_historical_rewards"

	AttributeKeyCluster    = "cluster"
	AttributeKeyMember     = "member"
	AttributeKeyTime       = "time"
	AttributeValueCategory = ModuleName

	
	EventTypeAddMembers = "add_members"
	AttributeClusterId  = "cluster_id"
	AttributeDaoFee     = "dao_fee"

	EventTypeDeductionFee = "deduction"
	AttributeDeductionFee = "deduction_fee"

	
	EventTypeBurnReward     = "burn_reward"
	AttributeGatewayReward  = "burn_gateway_reward"
	AttributeOwnerDaoReward = "burn_owner_dao_reward"

	
	EventTypeBurn   = "burn_to_power"
	AttributeToAddr = "burn_to_addr"

	
	EventTypeChangeSalaryRatio = "change_salary_ratio"

	
	EventTypeChangeDvmRatio = "change_dvm_ratio"

	
	EventTypeChangeClusterId = "change_cluster_id"

	
	EventTypeWithdrawSalary = "get_salary"

	
	EventTypUpdateAdmin = "update_admin"

	
	GenesisIdoEvent = "InvestLog"

	
	WdstEventEndBurnToPower = "EndBurnToPower"

	AuthEventApproveDvmLog = "ApproveDvmLog"

	
	EventTypeClusterAdLog = "ClusterAd"
	AttributeAdText       = "ad_text"

	
	GetDao      = "get_dao"
	GetDaoQueue = "get_dao_queue"
	GetDaoRoute = "get_dao_route"
	DaoDelay    = "dao_delay"
	
	ChangeClusterDaoRatio = "change_cluster_dao_ratio"

	RedPacketLog = "SendBagLog"

	OpenBagLog = "OpenBagLog"

	
	EventTypeCompleteDao = "complete_dao"

	
	EventTypeStartPowerRewardRedeem = "start_power_reward_redeem"

	
	EventTypeReceiveCutReward = "receive_cut_reward"
)
