package client

import (
	"encoding/json"
	"errors"
	"freemasonry.cc/blockchain/util"
	chatTypes "freemasonry.cc/blockchain/x/chat/types"
	contractTypes "freemasonry.cc/blockchain/x/contract/types"
	daotypes "freemasonry.cc/blockchain/x/dao/types"
	"freemasonry.cc/blockchain/x/gateway/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authType "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distributionTypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	gov "github.com/cosmos/cosmos-sdk/x/gov/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	govTypes "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	"github.com/cosmos/cosmos-sdk/x/group"
	"github.com/cosmos/cosmos-sdk/x/params/types/proposal"
	paramproposal "github.com/cosmos/cosmos-sdk/x/params/types/proposal"
	stakeTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	upgradeTypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	"strings"
	
	ibcTypes "github.com/cosmos/ibc-go/v5/modules/apps/transfer/types"
)

func unmashalMsgSend(msgByte []byte) (sdk.Msg, error) {
	msg := bankTypes.MsgSend{}
	err := util.Json.Unmarshal(msgByte, &msg)
	return &msg, err
}

func unmashalMsgDelegate(msgByte []byte) (sdk.Msg, error) {
	msg := stakeTypes.MsgDelegate{}
	err := util.Json.Unmarshal(msgByte, &msg)
	return &msg, err
}

func unmashalDistMsgWithdrawDelegatorReward(msgByte []byte) (sdk.Msg, error) {
	msg := distributionTypes.MsgWithdrawDelegatorReward{}
	err := util.Json.Unmarshal(msgByte, &msg)
	return &msg, err
}

func unmashalMsgVote(msgByte []byte) (sdk.Msg, error) {
	msg := govTypes.MsgVote{}
	err := util.Json.Unmarshal(msgByte, &msg)
	return &msg, err
}

func unmashalMsgMobileTransfer(msgByte []byte) (sdk.Msg, error) {
	msg := chatTypes.MsgMobileTransfer{}
	err := util.Json.Unmarshal(msgByte, &msg)
	return &msg, err
}

func unmashalMsgBurnGetMobile(msgByte []byte) (sdk.Msg, error) {
	msg := chatTypes.MsgBurnGetMobile{}
	err := util.Json.Unmarshal(msgByte, &msg)
	return &msg, err
}

func unmashalMsgSetChatInfo(msgByte []byte) (sdk.Msg, error) {
	msg := SetChatInfo{}
	err := util.Json.Unmarshal(msgByte, &msg)
	if err != nil {
		return nil, err
	}

	realMas := chatTypes.MsgSetChatInfo{
		FromAddress:      msg.FromAddress,
		GatewayAddress:   msg.NodeAddress,
		AddressBook:      msg.AddressBook,
		ChatBlacklist:    msg.ChatBlacklist,
		ChatWhitelist:    msg.ChatWhitelist,
		UpdateTime:       msg.UpdateTime,
		ChatBlacklistEnc: msg.ChatBlacklistEnc,
		ChatWhitelistEnc: msg.ChatWhitelistEnc,
	}

	return &realMas, err
}

func unmashalMsgAppTokenIssue(msgByte []byte) (sdk.Msg, error) {
	msg := contractTypes.MsgAppTokenIssue{}
	err := util.Json.Unmarshal(msgByte, &msg)
	return &msg, err
}

func unmashalProposalParams(msgByte []byte) (sdk.Msg, error) {
	proposals := struct {
		Proposer    string `json:"proposer"`
		Deposit     string `json:"deposit"`
		Title       string `json:"title"`
		Description string `json:"description"`
		Change      string `json:"change"`
	}{}
	err := util.Json.Unmarshal(msgByte, &proposals)
	if err != nil {
		return nil, err
	}
	var changes []proposal.ParamChange
	err = json.Unmarshal([]byte(proposals.Change), &changes)
	if err != nil {
		return nil, err
	}

	for k, change := range changes {
		if change.Subspace == "staking" && change.Key == "MaxValidators" {
			changes[k].Value = strings.Replace(changes[k].Value, "\"", "", -1)
		}
	}

	govModuleAcc := authType.NewModuleAddress(gov.ModuleName)
	content := paramproposal.NewParameterChangeProposal(proposals.Title, proposals.Description, changes)
	legacyContent, err := govtypes.NewLegacyContent(content, govModuleAcc.String())

	deposit, err := sdk.ParseCoinsNormalized(proposals.Deposit)
	if err != nil {
		return nil, err
	}

	msg, err := govtypes.NewMsgSubmitProposal([]sdk.Msg{legacyContent}, deposit, proposals.Proposer, proposals.Description)
	if err != nil {
		return nil, err
	}
	err = msg.ValidateBasic()
	if err != nil {
		return nil, err
	}
	return msg, nil
}

func unmashalProposalCommunity(msgByte []byte) (sdk.Msg, error) {
	proposal := struct {
		Proposer    string `json:"proposer"`
		Deposit     string `json:"deposit"`
		Title       string `json:"title"`
		Description string `json:"description"`
		Amount      string `json:"amount"`
		Recipient   string `json:"recipient"`
	}{}
	err := util.Json.Unmarshal(msgByte, &proposal)
	if err != nil {
		return nil, err
	}
	recpAddr, err := sdk.AccAddressFromBech32(proposal.Recipient)
	if err != nil {
		return nil, err
	}
	am, err := sdk.ParseCoinsNormalized(proposal.Amount)
	if err != nil {
		return nil, err
	}
	govModuleAcc := authType.NewModuleAddress(gov.ModuleName)

	content := distributionTypes.NewCommunityPoolSpendProposal(proposal.Title, proposal.Description, recpAddr, am)
	legacyContent, err := govtypes.NewLegacyContent(content, govModuleAcc.String())

	deposit, err := sdk.ParseCoinsNormalized(proposal.Deposit)
	if err != nil {
		return nil, err
	}

	msg, err := govtypes.NewMsgSubmitProposal([]sdk.Msg{legacyContent}, deposit, proposal.Proposer, proposal.Description)
	if err != nil {
		return nil, err
	}
	err = msg.ValidateBasic()
	if err != nil {
		return nil, err
	}
	return msg, nil
}

func unmashalProposalUpgrade(msgByte []byte) (sdk.Msg, error) {
	proposal := struct {
		Proposer    string `json:"proposer"`
		Deposit     string `json:"deposit"`
		Title       string `json:"title"`
		Description string `json:"description"`
		Info        string `json:"info"`
		Height      int64  `json:"height"`
		Name        string `json:"name"`
	}{}
	err := util.Json.Unmarshal(msgByte, &proposal)
	if err != nil {
		return nil, err
	}
	plan := upgradeTypes.Plan{
		Name:   proposal.Name,
		Height: proposal.Height,
		Info:   proposal.Info,
	}

	err = UpgradeJsonValidateBasic(plan)
	if err != nil {
		return nil, err
	}

	govModuleAcc := authType.NewModuleAddress(gov.ModuleName)

	content := upgradeTypes.NewSoftwareUpgradeProposal(proposal.Title, proposal.Description, plan)
	legacyContent, err := govtypes.NewLegacyContent(content, govModuleAcc.String())

	deposit, err := sdk.ParseCoinsNormalized(proposal.Deposit)
	if err != nil {
		return nil, err
	}

	msg, err := govtypes.NewMsgSubmitProposal([]sdk.Msg{legacyContent}, deposit, proposal.Proposer, proposal.Description)
	if err != nil {
		return nil, err
	}
	err = msg.ValidateBasic()
	if err != nil {
		return nil, err
	}
	return msg, nil
}

func UpgradeJsonValidateBasic(plan upgradeTypes.Plan) error {
	info, err := plan.UpgradeInfo()
	if err != nil {
		return err
	}
	if info.Gateway == nil && info.App == nil && info.Blockchain == nil {
		return errors.New("The json content is illegal")
	}
	return nil
}













func unmashalMsgGatewayUndelegate(msgByte []byte) (sdk.Msg, error) {
	msg := types.MsgGatewayUndelegate{}
	err := util.Json.Unmarshal(msgByte, &msg)
	return &msg, err
}

func unmashalMsgGatewayEdit(msgByte []byte) (sdk.Msg, error) {
	msg := types.MsgGatewayEdit{}
	err := util.Json.Unmarshal(msgByte, &msg)
	return &msg, err
}

func unmashalMsgGatewayBeginRedelegate(msgByte []byte) (sdk.Msg, error) {
	msg := types.MsgGatewayBeginRedelegate{}
	err := util.Json.Unmarshal(msgByte, &msg)
	return &msg, err
}

func unmashalMsgGatewayUpload(msgByte []byte) (sdk.Msg, error) {
	msg := types.MsgGatewayUpload{}
	err := util.Json.Unmarshal(msgByte, &msg)
	return &msg, err
}

func unmashalMsgGatewayIndexNum(msgByte []byte) (sdk.Msg, error) {
	msg := types.MsgGatewayIndexNum{}
	err := util.Json.Unmarshal(msgByte, &msg)
	return &msg, err
}

func unmashalMsgGatewayRegister(msgByte []byte) (sdk.Msg, error) {
	msg := types.MsgGatewayRegister{}
	err := util.Json.Unmarshal(msgByte, &msg)
	return &msg, err
}

func unmashalMsgCreateSmartValidator(msgByte []byte) (sdk.Msg, error) {
	msg := types.MsgCreateSmartValidator{}
	err := util.Json.Unmarshal(msgByte, &msg)
	return &msg, err
}



func unmashalDaoMsgCreateCluster(msgByte []byte) (sdk.Msg, error) {
	msg := daotypes.MsgCreateCluster{}
	err := util.Json.Unmarshal(msgByte, &msg)
	return &msg, err
}

func unmashalDaoMsgDeleteMembers(msgByte []byte) (sdk.Msg, error) {
	msg := daotypes.MsgDeleteMembers{}
	err := util.Json.Unmarshal(msgByte, &msg)
	return &msg, err
}

func unmashalDaoMsgClusterAddMembers(msgByte []byte) (sdk.Msg, error) {
	msg := daotypes.MsgClusterAddMembers{}
	err := util.Json.Unmarshal(msgByte, &msg)
	return &msg, err
}

func unmashalDaoMsgAgreeJoinCluster(msgByte []byte) (sdk.Msg, error) {
	msg := daotypes.MsgAgreeJoinCluster{}
	err := util.Json.Unmarshal(msgByte, &msg)
	return &msg, err
}

func unmashalDaoMsgAgreeJoinClusterApply(msgByte []byte) (sdk.Msg, error) {
	msg := daotypes.MsgAgreeJoinClusterApply{}
	err := util.Json.Unmarshal(msgByte, &msg)
	return &msg, err
}

func unmashalDaoMsgClusterChangeName(msgByte []byte) (sdk.Msg, error) {
	msg := daotypes.MsgClusterChangeName{}
	err := util.Json.Unmarshal(msgByte, &msg)
	return &msg, err
}

func unmashalDaoMsgChangeDaoRatio(msgByte []byte) (sdk.Msg, error) {
	msg := daotypes.MsgClusterChangeDaoRatio{}
	err := util.Json.Unmarshal(msgByte, &msg)
	return &msg, err
}

func unmashalDaoMsgClusterMemberExit(msgByte []byte) (sdk.Msg, error) {
	msg := daotypes.MsgClusterMemberExit{}
	err := util.Json.Unmarshal(msgByte, &msg)
	return &msg, err
}

func unmashalDaoMsgBurnToPower(msgByte []byte) (sdk.Msg, error) {
	msg := daotypes.MsgBurnToPower{}
	err := util.Json.Unmarshal(msgByte, &msg)
	return &msg, err
}

func unmashalDaoMsgClusterChangeSalaryRatio(msgByte []byte) (sdk.Msg, error) {
	msg := daotypes.MsgClusterChangeSalaryRatio{}
	err := util.Json.Unmarshal(msgByte, &msg)
	return &msg, err
}
func unmashalDaoMsgClusterChangeDvmRatio(msgByte []byte) (sdk.Msg, error) {
	msg := daotypes.MsgClusterChangeDvmRatio{}
	err := util.Json.Unmarshal(msgByte, &msg)
	return &msg, err
}

func unmashalDaoMsgClusterChangeId(msgByte []byte) (sdk.Msg, error) {
	msg := daotypes.MsgClusterChangeId{}
	err := util.Json.Unmarshal(msgByte, &msg)
	return &msg, err
}

func unmashalDaoMsgWithdrawDeviceReward(msgByte []byte) (sdk.Msg, error) {
	msg := daotypes.MsgWithdrawDeviceReward{}
	err := util.Json.Unmarshal(msgByte, &msg)
	return &msg, err
}

func unmashalDaoMsgWithdrawOwnerReward(msgByte []byte) (sdk.Msg, error) {
	msg := daotypes.MsgWithdrawOwnerReward{}
	err := util.Json.Unmarshal(msgByte, &msg)
	return &msg, err
}

func unmashalDaoMsgWithdrawBurnReward(msgByte []byte) (sdk.Msg, error) {
	msg := daotypes.MsgWithdrawSwapDpos{}
	err := util.Json.Unmarshal(msgByte, &msg)
	return &msg, err
}

func unmashalDaoMsgThawFrozenPower(msgByte []byte) (sdk.Msg, error) {
	msg := daotypes.MsgThawFrozenPower{}
	err := util.Json.Unmarshal(msgByte, &msg)
	return &msg, err
}
func unmashalPersonDvmApprove(msgByte []byte) (sdk.Msg, error) {
	msg := daotypes.MsgPersonDvmApprove{}
	err := util.Json.Unmarshal(msgByte, &msg)
	return &msg, err
}

func unmashalDaoMsgColonyRate(msgByte []byte) (sdk.Msg, error) {
	msg := daotypes.MsgColonyRate{}
	err := util.Json.Unmarshal(msgByte, &msg)
	return &msg, err
}
func unmashalDaoMsgClusterAd(msgByte []byte) (sdk.Msg, error) {
	msg := daotypes.MsgClusterAd{}
	err := util.Json.Unmarshal(msgByte, &msg)
	return &msg, err
}
func unmashalDaoMsgReceiveBurnRewardFee(msgByte []byte) (sdk.Msg, error) {
	msg := daotypes.MsgReceiveBurnRewardFee{}
	err := util.Json.Unmarshal(msgByte, &msg)
	return &msg, err
}
func unmashalDaoMsgRedPacket(msgByte []byte) (sdk.Msg, error) {
	msg := daotypes.MsgRedPacket{}
	err := util.Json.Unmarshal(msgByte, &msg)
	return &msg, err
}
func unmashalDaoMsgOpenRedPacket(msgByte []byte) (sdk.Msg, error) {
	msg := daotypes.MsgOpenRedPacket{}
	err := util.Json.Unmarshal(msgByte, &msg)
	return &msg, err
}
func unmashalDaoMsgReturnRedPacket(msgByte []byte) (sdk.Msg, error) {
	msg := daotypes.MsgReturnRedPacket{}
	err := util.Json.Unmarshal(msgByte, &msg)
	return &msg, err
}

func unmashalDaoMsgCreateClusterAddMembers(msgByte []byte) (sdk.Msg, error) {
	msg := daotypes.MsgCreateClusterAddMembers{}
	err := util.Json.Unmarshal(msgByte, &msg)
	return &msg, err
}

func unmashalDaoMsgStartPowerRewardRedeem(msgByte []byte) (sdk.Msg, error) {
	msg := daotypes.MsgStartPowerRewardRedeem{}
	err := util.Json.Unmarshal(msgByte, &msg)
	return &msg, err
}

func unmashalDaoMsgReceivePowerCutReward(msgByte []byte) (sdk.Msg, error) {
	msg := daotypes.MsgReceivePowerCutReward{}
	err := util.Json.Unmarshal(msgByte, &msg)
	return &msg, err
}

type MsgSubmitProposalClone struct {
	GroupPolicyAddress string     `protobuf:"bytes,1,opt,name=group_policy_address,json=groupPolicyAddress,proto3" json:"group_policy_address,omitempty"`
	Proposers          []string   `protobuf:"bytes,2,rep,name=proposers,proto3" json:"proposers,omitempty"`
	Metadata           string     `protobuf:"bytes,3,opt,name=metadata,proto3" json:"metadata,omitempty"`
	Messages           []sdk.Msg  `protobuf:"bytes,4,rep,name=messages,proto3" json:"messages,omitempty"`
	Exec               group.Exec `protobuf:"varint,5,opt,name=exec,proto3,enum=cosmos.group.v1.Exec" json:"exec,omitempty"`
}

func unmashalMsgSubmitProposal(msgByte []byte) (sdk.Msg, error) {
	msg := MsgSubmitProposalClone{}
	err := util.Json.Unmarshal(msgByte, &msg)

	newmsg := group.MsgSubmitProposal{}
	newmsg.GroupPolicyAddress = msg.GroupPolicyAddress
	newmsg.Proposers = msg.Proposers
	newmsg.Metadata = msg.Metadata
	newmsg.Exec = msg.Exec
	newmsg.SetMsgs(msg.Messages)
	
	return &newmsg, err
}
func unmashalGroupMsgVote(msgByte []byte) (sdk.Msg, error) {
	msg := group.MsgVote{}
	err := util.Json.Unmarshal(msgByte, &msg)
	return &msg, err
}
func unmashalGroupMsgExec(msgByte []byte) (sdk.Msg, error) {
	msg := group.MsgExec{}
	err := util.Json.Unmarshal(msgByte, &msg)
	return &msg, err
}

func unmashalIbcMsgTransfer(msgByte []byte) (sdk.Msg, error) {
	msg := ibcTypes.MsgTransfer{}
	err := util.Json.Unmarshal(msgByte, &msg)
	return &msg, err
}
