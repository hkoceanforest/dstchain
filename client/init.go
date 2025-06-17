package client

import (
	"encoding/json"
	"freemasonry.cc/blockchain/core"
	daotypes "freemasonry.cc/blockchain/x/dao/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	distributionTypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	ibcTypes "github.com/cosmos/ibc-go/v5/modules/apps/transfer/types"
	ibcclienttypes "github.com/cosmos/ibc-go/v5/modules/core/02-client/types"
)


var msgUnmashalHandles map[string]func(msgByte []byte) (sdk.Msg, error)

func registerUnmashalHandles(msgType string, callback func(msgByte []byte) (sdk.Msg, error), msgJson func() ([]byte, error)) {
	if msgJson != nil {
		
		
	}
	msgUnmashalHandles[msgType] = callback
}

func init() {
	
	
	
	
	
	
	
	
	

	msgUnmashalHandles = make(map[string]func(data []byte) (sdk.Msg, error))

	registerUnmashalHandles("cosmos-sdk/MsgWithdrawDelegationReward", unmashalDistMsgWithdrawDelegatorReward, func() ([]byte, error) {
		return json.Marshal(&distributionTypes.MsgWithdrawDelegatorReward{DelegatorAddress: "dst13g9x8juqmhhtkdrc7xpdn90259g45we8z5kvet", ValidatorAddress: "dstvaloper1emqqtns9xdrpyjfant4wzf7zsylc59ydypaej5"})
	})
	registerUnmashalHandles("cosmos-sdk/MsgSend", unmashalMsgSend, nil)
	registerUnmashalHandles("cosmos-sdk/MsgDelegate", unmashalMsgDelegate, nil)
	registerUnmashalHandles("cosmos-sdk/MsgVote", unmashalMsgVote, nil)

	registerUnmashalHandles("dao/MsgColonyRate", unmashalDaoMsgColonyRate, nil)

	registerUnmashalHandles("dao/MsgStartPowerRewardRedeem", unmashalDaoMsgStartPowerRewardRedeem, func() ([]byte, error) {
		return json.Marshal(&daotypes.MsgStartPowerRewardRedeem{
			FromAddress: "dst13g9x8juqmhhtkdrc7xpdn90259g45we8z5kvet",
		})
	})

	registerUnmashalHandles("dao/MsgReceivePowerCutReward", unmashalDaoMsgReceivePowerCutReward, func() ([]byte, error) {
		return json.Marshal(&daotypes.MsgReceivePowerCutReward{
			FromAddress: "dst13g9x8juqmhhtkdrc7xpdn90259g45we8z5kvet",
		})
	})

	
	registerUnmashalHandles("dao/MsgCreateCluster", unmashalDaoMsgCreateCluster, func() ([]byte, error) {
		return json.Marshal(&daotypes.MsgCreateCluster{
			FromAddress:     "dst13g9x8juqmhhtkdrc7xpdn90259g45we8z5kvet",
			GateAddress:     "dstvaloper1emqqtns9xdrpyjfant4wzf7zsylc59ydypaej5",
			ClusterId:       "123455",
			SalaryRatio:     sdk.NewDec(1),
			BurnAmount:      sdk.NewDec(1),
			ChatAddress:     "chataddr",
			ClusterName:     "name",
			FreezeAmount:    sdk.NewDec(1),
			Metadata:        "",
			ClusterDaoRatio: sdk.MustNewDecFromStr("0.5"),
		})
	})

	
	registerUnmashalHandles("dao/MsgClusterAddMembers", unmashalDaoMsgClusterAddMembers, func() ([]byte, error) {
		return json.Marshal(&daotypes.MsgClusterAddMembers{
			FromAddress: "dst13g9x8juqmhhtkdrc7xpdn90259g45we8z5kvet",
			ClusterId:   "123455",
			Members:     []daotypes.Members{{MemberAddress: "dst13g9x8juqmhhtkdrc7xpdn90259g45we8z5kvet", IndexNum: "IndexNum", ChatAddress: "ChatAddress"}}})
	})

	
	registerUnmashalHandles("dao/MsgAgreeJoinCluster", unmashalDaoMsgAgreeJoinCluster, func() ([]byte, error) {
		return json.Marshal(&daotypes.MsgAgreeJoinCluster{
			FromAddress:        "dstxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
			ClusterId:          "!xasdasda.nxn",
			Sign:               "abcabc",
			IndexNum:           "1111111",
			ChatAddress:        "dstxxxxxxxxxxxs",
			MemberOnlineAmount: 4,
			GatewayAddress:     "dstvaloper1ll30h0xykgduvxxfnpy4h6yzl0770pgn7hn3lz",
			GatewaySign:        "asdasdasdasda",
		})
	})

	
	registerUnmashalHandles("dao/MsgAgreeJoinClusterApply", unmashalDaoMsgAgreeJoinClusterApply, func() ([]byte, error) {
		return json.Marshal(&daotypes.MsgAgreeJoinClusterApply{
			FromAddress:        "dst1q26kg06nwrju5cc34u007v6tp54qartk2ql6rz",
			ClusterId:          "!zOnfR8GEhgtDGYP5:1111111.nxn",
			Sign:               "abcabc",
			IndexNum:           "1111111",
			ChatAddress:        "dstxxxxxxxxxxxs",
			MemberAddress:      "dst1q26kg06nwrju5cc34u007v6tp54qartk2ql6rz",
			MemberOnlineAmount: 3,
			GatewayAddress:     "dstvaloper1ll30h0xykgduvxxfnpy4h6yzl0770pgn7hn3lz",
			GatewaySign:        "asdasdasd",
		})
	})

	
	registerUnmashalHandles("dao/MsgDeleteMembers", unmashalDaoMsgDeleteMembers, func() ([]byte, error) {
		return json.Marshal(
			&daotypes.MsgDeleteMembers{
				FromAddress:        "dst13g9x8juqmhhtkdrc7xpdn90259g45we8z5kvet",
				ClusterId:          "123455",
				Members:            []string{"dst13g9x8juqmhhtkdrc7xpdn90259g45we8z5kvet", "dst13g9x8juqmhhtkdrc7xpdn90259g45we8z5kvet"},
				MemberOnlineAmount: 2,
				GatewayAddress:     "dstvaloper1ll30h0xykgduvxxfnpy4h6yzl0770pgn7hn3lz",
				GatewaySign:        "asdasdasdas",
			},
		)
	})

	
	registerUnmashalHandles("dao/MsgClusterChangeName", unmashalDaoMsgClusterChangeName, func() ([]byte, error) {
		return json.Marshal(&daotypes.MsgClusterChangeName{
			FromAddress: "dst13g9x8juqmhhtkdrc7xpdn90259g45we8z5kvet",
			ClusterId:   "123455",
			ClusterName: "ClusterName"})
	})

	
	registerUnmashalHandles("dao/MsgClusterMemberExit", unmashalDaoMsgClusterMemberExit, func() ([]byte, error) {
		return json.Marshal(&daotypes.MsgClusterMemberExit{
			FromAddress:        "dst13g9x8juqmhhtkdrc7xpdn90259g45we8z5kvet",
			ClusterId:          "123455",
			MemberOnlineAmount: 1,
			GatewayAddress:     "dstvaloper1ll30h0xykgduvxxfnpy4h6yzl0770pgn7hn3lz",
			GatewaySign:        "asdasdasdas",
		},
		)
	})

	
	registerUnmashalHandles("dao/MsgBurnToPower", unmashalDaoMsgBurnToPower, func() ([]byte, error) {
		return json.Marshal(&daotypes.MsgBurnToPower{
			FromAddress:     "dst13g9x8juqmhhtkdrc7xpdn90259g45we8z5kvet",
			ToAddress:       "dstvaloper1emqqtns9xdrpyjfant4wzf7zsylc59ydypaej5",
			BurnAmount:      sdk.NewDec(1),
			UseFreezeAmount: sdk.NewDec(1),
			GatewayAddress:  "dstvoloperxxxxxxxxxxxxxxxx",
			ChatAddress:     "ChatAddress",
			ClusterId:       "123455"})
	})

	
	
	
	
	
	
	

	
	registerUnmashalHandles("dao/MsgClusterChangeSalaryRatio", unmashalDaoMsgClusterChangeSalaryRatio, func() ([]byte, error) {
		return json.Marshal(&daotypes.MsgClusterChangeSalaryRatio{
			FromAddress: "dst13g9x8juqmhhtkdrc7xpdn90259g45we8z5kvet",
			SalaryRatio: sdk.NewDec(1),
			ClusterId:   "123455"})
	})
	
	registerUnmashalHandles("dao/MsgClusterChangeDvmRatio", unmashalDaoMsgClusterChangeDvmRatio, func() ([]byte, error) {
		return json.Marshal(&daotypes.MsgClusterChangeDvmRatio{
			FromAddress: "dst13g9x8juqmhhtkdrc7xpdn90259g45we8z5kvet",
			DvmRatio:    sdk.NewDec(1),
			ClusterId:   "123455"})
	})

	
	registerUnmashalHandles("dao/MsgClusterChangeId", unmashalDaoMsgClusterChangeId, func() ([]byte, error) {
		return json.Marshal(&daotypes.MsgClusterChangeId{
			FromAddress:  "dst13g9x8juqmhhtkdrc7xpdn90259g45we8z5kvet",
			NewClusterId: "NewClusterId",
			ClusterId:    "123455"})
	})

	
	registerUnmashalHandles("dao/MsgWithdrawOwnerReward", unmashalDaoMsgWithdrawOwnerReward, func() ([]byte, error) {
		return json.Marshal(&daotypes.MsgWithdrawOwnerReward{
			ClusterId: "!zOnfR8GEhgtDGYP5:1111111.nxn",
			Address:   "dstxxxxxxxxxxxxxxxxx",
		})
	})

	
	registerUnmashalHandles("dao/MsgChangeDaoRatio", unmashalDaoMsgChangeDaoRatio, func() ([]byte, error) {
		return json.Marshal(&daotypes.MsgClusterChangeDaoRatio{
			FromAddress: "dstxxxxxxxxxxxxxxxxx",
			ClusterId:   "!zOnfR8GEhgtDGYP5:1111111.nxn",
			DaoRatio:    sdk.MustNewDecFromStr("0.12"),
		})
	})

	
	registerUnmashalHandles("dao/MsgWithdrawSwapDpos", unmashalDaoMsgWithdrawBurnReward, func() ([]byte, error) {
		return json.Marshal(&daotypes.MsgWithdrawSwapDpos{
			ClusterId:     "!zOnfR8GEhgtDGYP5:1111111.nxn",
			MemberAddress: "dstxxxxxxxxxxxxxxxxx",
		})
	})

	
	registerUnmashalHandles("dao/MsgWithdrawDeviceReward", unmashalDaoMsgWithdrawDeviceReward, func() ([]byte, error) {
		return json.Marshal(&daotypes.MsgWithdrawDeviceReward{
			ClusterId:     "!zOnfR8GEhgtDGYP5:1111111.nxn",
			MemberAddress: "dstxxxxxxxxxxxxxxxxx",
		})
	})

	
	registerUnmashalHandles("dao/MsgThawFrozenPower", unmashalDaoMsgThawFrozenPower, func() ([]byte, error) {
		return json.Marshal(&daotypes.MsgThawFrozenPower{
			FromAddress:    "dstxxxxxxxxxxxxxxxxx",
			ClusterId:      "!zOnfR8GEhgtDGYP5:1111111.nxn",
			ThawAmount:     sdk.MustNewDecFromStr("12.34"),
			GatewayAddress: "dstveloperxxxxxxxxxxxxxxxxxxxxxxxx",
			ChatAddress:    "dstxxxxxxxxxxxxxxxxx",
		})
	})
	registerUnmashalHandles("dao/PersonDvmApprove", unmashalPersonDvmApprove, func() ([]byte, error) {
		return json.Marshal(&daotypes.MsgPersonDvmApprove{
			FromAddress:     "dstxxxxxxxxxxxxxxxxx",
			ClusterId:       "!zOnfR8GEhgtDGYP5:1111111.nxn",
			ApproveAddress:  "0x",
			ApproveEndBlock: "0",
		})
	})
	
	registerUnmashalHandles("dao/ClusterAd", unmashalDaoMsgClusterAd, nil)

	
	registerUnmashalHandles("dao/MsgReceiveBurnRewardFee", unmashalDaoMsgReceiveBurnRewardFee, func() ([]byte, error) {
		return json.Marshal(&daotypes.MsgReceiveBurnRewardFee{
			FromAddress: "dstxxxxxxxxxxxxxxxxx",
			ClusterId:   "!zOnfR8GEhgtDGYP5:1111111.nxn",
			Amount:      sdk.MustNewDecFromStr("1000000000000000000000"),
		})
	})

	
	registerUnmashalHandles("dao/MsgRedPacket", unmashalDaoMsgRedPacket, func() ([]byte, error) {
		return json.Marshal(&daotypes.MsgRedPacket{
			Fromaddress: "dst1sqyyuwuepcpx0sj97lcjlwmdapznjnh0pscc2k",
			Clusterid:   "!zOnfR8GEhgtDGYP5:1111111.nxn",
			Amount:      sdk.NewCoin(core.BaseDenom, sdk.NewInt(123000000)),
			Count:       3,
			Redtype:     1,
		})
	})

	
	registerUnmashalHandles("dao/MsgOpenRedPacket", unmashalDaoMsgOpenRedPacket, func() ([]byte, error) {
		return json.Marshal(&daotypes.MsgOpenRedPacket{
			Fromaddress: "dst1q26kg06nwrju5cc34u007v6tp54qartk2ql6rz",
			Redpacketid: "e6c2818ec06c5781051a9e1953942724557b28a75f95e04bfb9a800c27fb692e",
		})
	})

	
	registerUnmashalHandles("dao/MsgReturnRedPacket", unmashalDaoMsgReturnRedPacket, func() ([]byte, error) {
		return json.Marshal(&daotypes.MsgReturnRedPacket{
			Fromaddress: "dst1q26kg06nwrju5cc34u007v6tp54qartk2ql6rz",
			Redpacketid: "e6c2818ec06c5781051a9e1953942724557b28a75f95e04bfb9a800c27fb692e",
		})
	})

	
	registerUnmashalHandles("dao/MsgCreateClusterAddMembers", unmashalDaoMsgCreateClusterAddMembers, func() ([]byte, error) {
		m := make([]daotypes.Members, 0)

		m = append(m, daotypes.Members{MemberAddress: "dst1q26kg06nwrju5cc34u007v6tp54qartk2ql6rz", IndexNum: "1", ChatAddress: "chataddr1"})
		m = append(m, daotypes.Members{MemberAddress: "dst1q26kg06nwrju5cc34u007v6tp54qartk2ql6rz", IndexNum: "2", ChatAddress: "chataddr2"})

		return json.Marshal(&daotypes.MsgCreateClusterAddMembers{
			FromAddress:        "dst1q26kg06nwrju5cc34u007v6tp54qartk2ql6rz",
			GateAddress:        "dstoprator1emqqtns9xdrpyjfant4wzf7zsylc59ydypaej51",
			ClusterId:          "!zOnfR8GEhgtDGYP5:1111111.nxn",
			SalaryRatio:        sdk.NewDec(1),
			BurnAmount:         sdk.NewDec(1),
			ChatAddress:        "dst1q26kg06nwrju5cc34u007v6tp54qartk2ql6rz",
			ClusterName:        "foo",
			FreezeAmount:       sdk.NewDec(1),
			Metadata:           "",
			ClusterDaoRatio:    sdk.MustNewDecFromStr("0.5"),
			Members:            m,
			MemberOnlineAmount: 10,
			OwnerIndexNum:      "123456",
			GatewaySign:        "asdasdasd",
		})
	})

	registerUnmashalHandles("ibc/MsgTransfer", unmashalIbcMsgTransfer, func() ([]byte, error) {
		return json.Marshal(&ibcTypes.MsgTransfer{
			SourcePort:    "transfer",
			SourceChannel: "channel-0",
			Token:         sdk.NewCoin("dst", sdk.NewInt(10)),
			Sender:        "dstxxxxxxxx",
			Receiver:      "kavaxxxxxxx",
			TimeoutHeight: ibcclienttypes.Height{
				RevisionNumber: 2,
				RevisionHeight: 2,
			},
			TimeoutTimestamp: 1234567890123456789,
			Memo:             "asdklfjhdsl",
		})
	})
	registerUnmashalHandles("chat/MsgMobileTransfer", unmashalMsgMobileTransfer, nil)
	registerUnmashalHandles("chat/MsgBurnGetMobile", unmashalMsgBurnGetMobile, nil)
	registerUnmashalHandles("chat/MsgSetChatInfo", unmashalMsgSetChatInfo, nil)

	registerUnmashalHandles("gateway/MsgCreateSmartValidator", unmashalMsgCreateSmartValidator, nil)     
	registerUnmashalHandles("gateway/MsgGatewayRegister", unmashalMsgGatewayRegister, nil)               
	registerUnmashalHandles("gateway/MsgGatewayIndexNum", unmashalMsgGatewayIndexNum, nil)               
	registerUnmashalHandles("gateway/MsgGatewayUndelegate", unmashalMsgGatewayUndelegate, nil)           
	registerUnmashalHandles("gateway/MsgGatewayBeginRedelegate", unmashalMsgGatewayBeginRedelegate, nil) 
	registerUnmashalHandles("gateway/MsgGatewayUpload", unmashalMsgGatewayUpload, nil)                   
	registerUnmashalHandles("gateway/MsgGatewayEdit", unmashalMsgGatewayEdit, nil)                       

	registerUnmashalHandles("contract/AppTokenIssue", unmashalMsgAppTokenIssue, nil)
	registerUnmashalHandles("proposal_upgrade", unmashalProposalUpgrade, nil)
	registerUnmashalHandles("proposal_community", unmashalProposalCommunity, nil)
	registerUnmashalHandles("proposal_params", unmashalProposalParams, nil)
	registerUnmashalHandles("cosmos-sdk/group/MsgSubmitProposal", unmashalMsgSubmitProposal, nil)
	registerUnmashalHandles("cosmos-sdk/group/MsgVote", unmashalGroupMsgVote, nil)
	registerUnmashalHandles("cosmos-sdk/group/MsgExec", unmashalGroupMsgExec, nil)
}
