package types

import (
	sdkmath "cosmossdk.io/math"
	"freemasonry.cc/blockchain/core"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	_ sdk.Msg = &MsgColonyRate{}
)

const (
	TypeMsgColonyRate                = "dao/ColonyRate"
	TypeMsgCreateCluster             = "dao/CreateCluster"
	TypeMsgClusterAddMembers         = "dao/ClusterAddMembers"
	TypeMsgBurnToPower               = "dao/BurnToPower"
	TypeMsgClusterChangeSalaryRatio  = "dao/ClusterChangeSalaryRatio"
	TypeMsgClusterChangeDvmRatio     = "dao/ClusterChangeDvmRatio"
	TypeMsgClusterChangeid           = "dao/ClusterChangeid"
	TypeMsgWithdrawOwnerReward       = "dao/WithdrawOwnerReward"
	TypeMsgWithdrawSwapDpos          = "dao/WithdrawSwapDpos"
	TypeMsgWithdrawDeviceReward      = "dao/WithdrawDeviceReward"
	TypeMsgDeleteMembers             = "dao/DeleteMembers"
	TypeMsgThawFrozenPower           = "dao/ThawFrozenPower"
	TypeMsgClusterMemberExit         = "dao/ClusterMemberExit"
	TypeMsgClusterChangeName         = "dao/ClusterChangeName"
	TypeMsgUpdateAdmin               = "dao/UpdateAdmin"
	TypeMsgClusterPowerApprove       = "dao/ClusterPowerApprove"
	TypePersonDvmApprove             = "dao/PersonDvmApprove"
	TypeMsgClusterAd                 = "dao/ClusterAd"
	TypeMsgReceiveBurnRewardFee      = "dao/ReceiveBurnRewardFee"
	TypeMsgClusterChangeDaoRatio     = "dao/ChangeDaoRatio"
	TypeMsgAgreeJoinCluster          = "dao/AgreeJoinCluster"
	TypeMsgRedPacket                 = "dao/RedPacket"
	TypeMsgOpenRedPacket             = "dao/OpenRedPacket"
	TypeMsgReturnRedPacket           = "dao/ReturnRedPacket"
	TypeMsgCreateClusterAddMembers   = "dao/CreateClusterAddMembers"
	TypeNewMsgAgreeJoinClusterApply  = "dao/MsgAgreeJoinClusterApply"
	TypeNewMsgStartPowerRewardRedeem = "dao/MsgStartPowerRewardRedeem"
	TypeNewMsgReceivePowerCutReward  = "dao/MsgReceivePowerCutReward"
)

func NewMsgStartPowerRewardRedeem(from string) *MsgStartPowerRewardRedeem {
	return &MsgStartPowerRewardRedeem{
		FromAddress: from,
	}
}

func (m MsgStartPowerRewardRedeem) Route() string { return sdk.MsgTypeURL(&m) }

func (m MsgStartPowerRewardRedeem) Type() string { return sdk.MsgTypeURL(&m) }

func (m MsgStartPowerRewardRedeem) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.FromAddress)
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{addr}
}
func (m MsgStartPowerRewardRedeem) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}
func (m MsgStartPowerRewardRedeem) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(m.FromAddress)
	if err != nil {
		return sdkerrors.Wrap(err, "invalid sender address")
	}
	return nil
}

func NewMsgReceivePowerCutReward(from string) *MsgReceivePowerCutReward {
	return &MsgReceivePowerCutReward{
		FromAddress: from,
	}
}

func (m MsgReceivePowerCutReward) Route() string { return sdk.MsgTypeURL(&m) }

func (m MsgReceivePowerCutReward) Type() string { return sdk.MsgTypeURL(&m) }

func (m MsgReceivePowerCutReward) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.FromAddress)
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{addr}
}
func (m MsgReceivePowerCutReward) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}
func (m MsgReceivePowerCutReward) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(m.FromAddress)
	if err != nil {
		return sdkerrors.Wrap(err, "invalid sender address")
	}
	return nil
}

func NewMsgAgreeJoinClusterApply(from, clusterId, sign, indexNum, chatAddress, memberAddress, gatewayAddress, gatewaySign string, onlineAmount int64) *MsgAgreeJoinClusterApply {
	return &MsgAgreeJoinClusterApply{
		FromAddress:        from,
		ClusterId:          clusterId,
		Sign:               sign,
		IndexNum:           indexNum,
		ChatAddress:        chatAddress,
		MemberAddress:      memberAddress,
		MemberOnlineAmount: onlineAmount,
		GatewayAddress:     gatewayAddress,
		GatewaySign:        gatewaySign,
	}
}

func (m MsgAgreeJoinClusterApply) Route() string { return sdk.MsgTypeURL(&m) }

func (m MsgAgreeJoinClusterApply) Type() string { return sdk.MsgTypeURL(&m) }

func (m MsgAgreeJoinClusterApply) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.FromAddress)
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{addr}
}
func (m MsgAgreeJoinClusterApply) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}
func (m MsgAgreeJoinClusterApply) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(m.FromAddress)
	if err != nil {
		return sdkerrors.Wrap(err, "invalid sender address")
	}

	_, err = sdk.AccAddressFromBech32(m.MemberAddress)
	if err != nil {
		return sdkerrors.Wrap(err, "invalid member address")
	}
	return nil
}

func NewMsgCreateClusterAddMembers(from, gateAddress, clusterId, chatAddress, clusterName, ownerIndexNum, gatewaySign string, salaryRatio, daoRatio, burnAmount, freezeAmount sdk.Dec, memberOnlineAmount int64, members []Members) *MsgCreateClusterAddMembers {
	return &MsgCreateClusterAddMembers{
		FromAddress:        from,
		GateAddress:        gateAddress,
		ClusterId:          clusterId,
		SalaryRatio:        salaryRatio,
		BurnAmount:         burnAmount,
		ChatAddress:        chatAddress,
		ClusterName:        clusterName,
		FreezeAmount:       freezeAmount,
		Metadata:           "",
		ClusterDaoRatio:    daoRatio,
		Members:            members,
		MemberOnlineAmount: memberOnlineAmount,
		OwnerIndexNum:      ownerIndexNum,
		GatewaySign:        gatewaySign,
	}
}

func (m MsgCreateClusterAddMembers) Route() string { return sdk.MsgTypeURL(&m) }

func (m MsgCreateClusterAddMembers) Type() string { return sdk.MsgTypeURL(&m) }

func (m MsgCreateClusterAddMembers) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.FromAddress)
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{addr}
}
func (m MsgCreateClusterAddMembers) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}
func (m MsgCreateClusterAddMembers) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(m.FromAddress)
	if err != nil {
		return sdkerrors.Wrap(err, "invalid sender address")
	}

	return nil
}

func NewMsgReturnRedPacket(from, redPacketId string) *MsgReturnRedPacket {
	return &MsgReturnRedPacket{
		Fromaddress: from,
		Redpacketid: redPacketId,
	}
}

func (m MsgReturnRedPacket) Route() string { return sdk.MsgTypeURL(&m) }

func (m MsgReturnRedPacket) Type() string { return sdk.MsgTypeURL(&m) }

func (m MsgReturnRedPacket) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.GetFromaddress())
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{addr}
}
func (m MsgReturnRedPacket) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}
func (m MsgReturnRedPacket) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(m.Fromaddress)
	if err != nil {
		return sdkerrors.Wrap(err, "invalid sender address")
	}
	return nil
}

func NewMsgOpenRedPacket(from, redPacketId string) *MsgOpenRedPacket {
	return &MsgOpenRedPacket{
		Fromaddress: from,
		Redpacketid: redPacketId,
	}
}

func (m MsgOpenRedPacket) Route() string { return sdk.MsgTypeURL(&m) }

func (m MsgOpenRedPacket) Type() string { return sdk.MsgTypeURL(&m) }

func (m MsgOpenRedPacket) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.GetFromaddress())
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{addr}
}
func (m MsgOpenRedPacket) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}
func (m MsgOpenRedPacket) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(m.Fromaddress)
	if err != nil {
		return sdkerrors.Wrap(err, "invalid sender address")
	}
	return nil
}

func NewMsgRedPacket(from, clusterId string, amount sdk.Coin, count int64, redPacketType int64) *MsgRedPacket {
	return &MsgRedPacket{
		Fromaddress: from,
		Clusterid:   clusterId,
		Amount:      amount,
		Count:       count,
		Redtype:     redPacketType,
	}
}

func (m MsgRedPacket) Route() string { return sdk.MsgTypeURL(&m) }

func (m MsgRedPacket) Type() string { return sdk.MsgTypeURL(&m) }

func (m MsgRedPacket) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.GetFromaddress())
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{addr}
}
func (m MsgRedPacket) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}
func (m MsgRedPacket) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(m.Fromaddress)
	if err != nil {
		return sdkerrors.Wrap(err, "invalid sender address")
	}
	return nil
}

func NewMsgAgreeJoinCluster(from, clusterId, sign, indexNum, chatAddress, gatewayAddress, gatewaySign string, onlineAmount int64) *MsgAgreeJoinCluster {
	return &MsgAgreeJoinCluster{
		FromAddress:        from,
		ClusterId:          clusterId,
		Sign:               sign,
		IndexNum:           indexNum,
		ChatAddress:        chatAddress,
		MemberOnlineAmount: onlineAmount,
		GatewayAddress:     gatewayAddress,
		GatewaySign:        gatewaySign,
	}
}

func (m MsgAgreeJoinCluster) Route() string { return sdk.MsgTypeURL(&m) }

func (m MsgAgreeJoinCluster) Type() string { return sdk.MsgTypeURL(&m) }

func (m MsgAgreeJoinCluster) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.GetFromAddress())
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{addr}
}
func (m MsgAgreeJoinCluster) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}
func (m MsgAgreeJoinCluster) ValidateBasic() error {

	_, err := sdk.AccAddressFromBech32(m.FromAddress)
	if err != nil {
		return sdkerrors.Wrap(err, "invalid sender address")
	}
	return nil
}

func NewMsgClusterChangeDaoRatio(from, clusterId string, ratio sdk.Dec) *MsgClusterChangeDaoRatio {
	return &MsgClusterChangeDaoRatio{
		FromAddress: from,
		ClusterId:   clusterId,
		DaoRatio:    ratio,
	}
}

func (m MsgClusterChangeDaoRatio) Route() string { return sdk.MsgTypeURL(&m) }

func (m MsgClusterChangeDaoRatio) Type() string { return sdk.MsgTypeURL(&m) }

func (m MsgClusterChangeDaoRatio) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.GetFromAddress())
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{addr}
}
func (m MsgClusterChangeDaoRatio) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}
func (m MsgClusterChangeDaoRatio) ValidateBasic() error {

	_, err := sdk.AccAddressFromBech32(m.FromAddress)
	if err != nil {
		return sdkerrors.Wrap(err, "invalid sender address")
	}
	return nil
}

func NewMsgReceiveBurnRewardFee(from, clusterId string, amount sdk.Dec) *MsgReceiveBurnRewardFee {
	return &MsgReceiveBurnRewardFee{
		FromAddress: from,
		ClusterId:   clusterId,
		Amount:      amount,
	}
}

func (m MsgReceiveBurnRewardFee) Route() string { return sdk.MsgTypeURL(&m) }

func (m MsgReceiveBurnRewardFee) Type() string { return sdk.MsgTypeURL(&m) }

func (m MsgReceiveBurnRewardFee) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.GetFromAddress())
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{addr}
}
func (m MsgReceiveBurnRewardFee) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}
func (m MsgReceiveBurnRewardFee) ValidateBasic() error {

	_, err := sdk.AccAddressFromBech32(m.FromAddress)
	if err != nil {
		return sdkerrors.Wrap(err, "invalid sender address")
	}
	return nil
}

func NewMsgClusterAd(from, adText string, clusterId []string) *MsgClusterAd {
	return &MsgClusterAd{
		FromAddress: from,
		ClusterId:   clusterId,
		AdText:      adText,
	}
}

func (m MsgClusterAd) Route() string { return sdk.MsgTypeURL(&m) }

func (m MsgClusterAd) Type() string { return sdk.MsgTypeURL(&m) }

func (m MsgClusterAd) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.GetFromAddress())
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{addr}
}
func (m MsgClusterAd) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}
func (m MsgClusterAd) ValidateBasic() error {

	_, err := sdk.AccAddressFromBech32(m.FromAddress)
	if err != nil {
		return sdkerrors.Wrap(err, "invalid sender address")
	}
	return nil
}

func NewMsgUpdateAdmin(
	from, clusterId string, clusterAdminList []string,
) *MsgUpdateAdmin {
	return &MsgUpdateAdmin{
		FromAddress:      from,
		ClusterId:        clusterId,
		ClusterAdminList: clusterAdminList,
	}
}

func (m MsgUpdateAdmin) Route() string { return sdk.MsgTypeURL(&m) }

func (m MsgUpdateAdmin) Type() string { return sdk.MsgTypeURL(&m) }

func (m MsgUpdateAdmin) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.GetFromAddress())
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{addr}
}
func (m MsgUpdateAdmin) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}
func (m MsgUpdateAdmin) ValidateBasic() error {

	_, err := sdk.AccAddressFromBech32(m.FromAddress)
	if err != nil {
		return sdkerrors.Wrap(err, "invalid sender address")
	}
	return nil
}

func NewMsgClusterChangeName(
	from, clusterId, clusterName string,
) *MsgClusterChangeName {
	return &MsgClusterChangeName{
		FromAddress: from,
		ClusterId:   clusterId,
		ClusterName: clusterName,
	}
}

func (m MsgClusterChangeName) Route() string { return sdk.MsgTypeURL(&m) }

func (m MsgClusterChangeName) Type() string { return sdk.MsgTypeURL(&m) }

func (m MsgClusterChangeName) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.GetFromAddress())
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{addr}
}
func (m MsgClusterChangeName) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}
func (m MsgClusterChangeName) ValidateBasic() error {

	_, err := sdk.AccAddressFromBech32(m.FromAddress)
	if err != nil {
		return sdkerrors.Wrap(err, "invalid sender address")
	}
	return nil
}

func NewMsgClusterMemberExit(fromAddress, clusterId, gatewayAddress, gatewaySign string, onlineAmount int64) *MsgClusterMemberExit {
	return &MsgClusterMemberExit{
		FromAddress:        fromAddress,
		ClusterId:          clusterId,
		MemberOnlineAmount: onlineAmount,
		GatewayAddress:     gatewayAddress,
		GatewaySign:        gatewaySign,
	}
}

func (m MsgClusterMemberExit) Route() string { return sdk.MsgTypeURL(&m) }

func (m MsgClusterMemberExit) Type() string { return sdk.MsgTypeURL(&m) }

func (m MsgClusterMemberExit) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.FromAddress)
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{addr}
}
func (m MsgClusterMemberExit) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}
func (m MsgClusterMemberExit) ValidateBasic() error {

	_, err := sdk.AccAddressFromBech32(m.FromAddress)
	if err != nil {
		return sdkerrors.Wrap(err, "invalid sender address")
	}
	return nil
}

func NewMsgThawFrozenPower(fromAddress, clusterId, gatewayAddr, chatAddr string, thawAmount sdk.Dec) *MsgThawFrozenPower {
	return &MsgThawFrozenPower{
		FromAddress:    fromAddress,
		ClusterId:      clusterId,
		ThawAmount:     thawAmount,
		GatewayAddress: gatewayAddr,
		ChatAddress:    chatAddr,
	}
}

func (m MsgThawFrozenPower) Route() string { return sdk.MsgTypeURL(&m) }

func (m MsgThawFrozenPower) Type() string { return sdk.MsgTypeURL(&m) }

func (m MsgThawFrozenPower) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.FromAddress)
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{addr}
}
func (m MsgThawFrozenPower) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}
func (m MsgThawFrozenPower) ValidateBasic() error {

	_, err := sdk.AccAddressFromBech32(m.FromAddress)
	if err != nil {
		return sdkerrors.Wrap(err, "invalid sender address")
	}
	return nil
}

func NewMsgDeleteMembersd(fromAddress, clusterId, gatewayAddress, gatewaySign string, onlineAmount int64, members []string) *MsgDeleteMembers {
	return &MsgDeleteMembers{
		FromAddress:        fromAddress,
		ClusterId:          clusterId,
		Members:            members,
		MemberOnlineAmount: onlineAmount,
		GatewayAddress:     gatewayAddress,
		GatewaySign:        gatewaySign,
	}
}

func (m MsgDeleteMembers) Route() string { return sdk.MsgTypeURL(&m) }

func (m MsgDeleteMembers) Type() string { return sdk.MsgTypeURL(&m) }

func (m MsgDeleteMembers) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.FromAddress)
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{addr}
}
func (m MsgDeleteMembers) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}
func (m MsgDeleteMembers) ValidateBasic() error {

	_, err := sdk.AccAddressFromBech32(m.FromAddress)
	if err != nil {
		return sdkerrors.Wrap(err, "invalid sender address")
	}
	return nil
}

func NewMsgWithdrawOwnerReward(address, clusterId string) *MsgWithdrawOwnerReward {
	return &MsgWithdrawOwnerReward{
		Address:   address,
		ClusterId: clusterId,
	}
}

func (m MsgWithdrawOwnerReward) Route() string { return sdk.MsgTypeURL(&m) }

func (m MsgWithdrawOwnerReward) Type() string { return sdk.MsgTypeURL(&m) }

func (m MsgWithdrawOwnerReward) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.Address)
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{addr}
}
func (m MsgWithdrawOwnerReward) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}
func (m MsgWithdrawOwnerReward) ValidateBasic() error {

	_, err := sdk.AccAddressFromBech32(m.Address)
	if err != nil {
		return sdkerrors.Wrap(err, "invalid sender address")
	}
	return nil
}

func NewMsgWithdrawSwapDpos(memberAddress, clusterId, daoNum string) *MsgWithdrawSwapDpos {
	return &MsgWithdrawSwapDpos{
		MemberAddress: memberAddress,
		ClusterId:     clusterId,
		DaoNum:        daoNum,
	}
}

func (m MsgWithdrawSwapDpos) Route() string { return sdk.MsgTypeURL(&m) }

func (m MsgWithdrawSwapDpos) Type() string { return sdk.MsgTypeURL(&m) }

func (m MsgWithdrawSwapDpos) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.MemberAddress)
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{addr}
}
func (m MsgWithdrawSwapDpos) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}
func (m MsgWithdrawSwapDpos) ValidateBasic() error {

	_, err := sdk.AccAddressFromBech32(m.MemberAddress)
	if err != nil {
		return sdkerrors.Wrap(err, "invalid sender address")
	}
	if m.DaoNum == "" {
		return core.FieldEmptyError
	}
	_, ok := sdkmath.NewIntFromString(m.DaoNum)
	if !ok {
		return core.ParseCoinError
	}
	return nil
}

func NewMsgWithdrawDeviceReward(memberAddress, clusterId string) *MsgWithdrawDeviceReward {
	return &MsgWithdrawDeviceReward{
		MemberAddress: memberAddress,
		ClusterId:     clusterId,
	}
}

func (m MsgWithdrawDeviceReward) Route() string { return sdk.MsgTypeURL(&m) }

func (m MsgWithdrawDeviceReward) Type() string { return sdk.MsgTypeURL(&m) }

func (m MsgWithdrawDeviceReward) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.MemberAddress)
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{addr}
}
func (m MsgWithdrawDeviceReward) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}
func (m MsgWithdrawDeviceReward) ValidateBasic() error {

	_, err := sdk.AccAddressFromBech32(m.MemberAddress)
	if err != nil {
		return sdkerrors.Wrap(err, "invalid sender address")
	}
	return nil
}

func NewMsgClusterChangeId(
	from, clusterId, newClusterId string,
) *MsgClusterChangeId {
	return &MsgClusterChangeId{
		FromAddress:  from,
		ClusterId:    clusterId,
		NewClusterId: newClusterId,
	}
}

func (m MsgClusterChangeId) Route() string { return sdk.MsgTypeURL(&m) }

func (m MsgClusterChangeId) Type() string { return sdk.MsgTypeURL(&m) }

func (m MsgClusterChangeId) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.GetFromAddress())
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{addr}
}
func (m MsgClusterChangeId) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}
func (m MsgClusterChangeId) ValidateBasic() error {

	_, err := sdk.AccAddressFromBech32(m.FromAddress)
	if err != nil {
		return sdkerrors.Wrap(err, "invalid sender address")
	}
	return nil
}

func NewMsgClusterAddMembers(
	from, clusterId string, members []Members,
) *MsgClusterAddMembers {
	return &MsgClusterAddMembers{
		FromAddress: from,
		ClusterId:   clusterId,
		Members:     members,
	}
}

func (m MsgClusterAddMembers) Route() string { return sdk.MsgTypeURL(&m) }

func (m MsgClusterAddMembers) Type() string { return sdk.MsgTypeURL(&m) }

func (m MsgClusterAddMembers) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.GetFromAddress())
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{addr}
}
func (m MsgClusterAddMembers) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}
func (m MsgClusterAddMembers) ValidateBasic() error {

	_, err := sdk.AccAddressFromBech32(m.FromAddress)
	if err != nil {
		return sdkerrors.Wrap(err, "invalid sender address")
	}
	return nil
}

func NewMsgCreateCluster(
	from, gateAddress, clusterId, chatAddress, clusterName string, salaryRatio, daoRatio, burnAmount, freezeAmount sdk.Dec,
) *MsgCreateCluster {
	return &MsgCreateCluster{
		FromAddress:     from,
		GateAddress:     gateAddress,
		ClusterId:       clusterId,
		SalaryRatio:     salaryRatio,
		BurnAmount:      burnAmount,
		ChatAddress:     chatAddress,
		ClusterName:     clusterName,
		FreezeAmount:    freezeAmount,
		ClusterDaoRatio: daoRatio,
	}
}

func (m MsgCreateCluster) Route() string { return sdk.MsgTypeURL(&m) }

func (m MsgCreateCluster) Type() string { return sdk.MsgTypeURL(&m) }

func (m MsgCreateCluster) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.GetFromAddress())
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{addr}
}
func (m MsgCreateCluster) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}
func (m MsgCreateCluster) ValidateBasic() error {

	_, err := sdk.AccAddressFromBech32(m.FromAddress)
	if err != nil {
		return sdkerrors.Wrap(err, "invalid sender address")
	}
	return nil
}

func NewMsgClusterChangeSalaryRatio(from, clusterId string, salaryRatio sdk.Dec) *MsgClusterChangeSalaryRatio {
	return &MsgClusterChangeSalaryRatio{
		FromAddress: from,
		ClusterId:   clusterId,
		SalaryRatio: salaryRatio,
	}
}

func (m MsgClusterChangeSalaryRatio) Route() string { return sdk.MsgTypeURL(&m) }

func (m MsgClusterChangeSalaryRatio) Type() string { return sdk.MsgTypeURL(&m) }

func (m MsgClusterChangeSalaryRatio) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.GetFromAddress())
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{addr}
}
func (m MsgClusterChangeSalaryRatio) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}
func (m MsgClusterChangeSalaryRatio) ValidateBasic() error {

	_, err := sdk.AccAddressFromBech32(m.FromAddress)
	if err != nil {
		return sdkerrors.Wrap(err, "invalid sender address")
	}
	return nil
}

func NewMsgClusterChangeDvmRatio(from, clusterId string, dvmRatio sdk.Dec) *MsgClusterChangeDvmRatio {
	return &MsgClusterChangeDvmRatio{
		FromAddress: from,
		ClusterId:   clusterId,
		DvmRatio:    dvmRatio,
	}
}

func (m MsgClusterChangeDvmRatio) Route() string { return sdk.MsgTypeURL(&m) }

func (m MsgClusterChangeDvmRatio) Type() string { return sdk.MsgTypeURL(&m) }

func (m MsgClusterChangeDvmRatio) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.GetFromAddress())
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{addr}
}
func (m MsgClusterChangeDvmRatio) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}
func (m MsgClusterChangeDvmRatio) ValidateBasic() error {

	_, err := sdk.AccAddressFromBech32(m.FromAddress)
	if err != nil {
		return sdkerrors.Wrap(err, "invalid sender address")
	}
	return nil
}







func NewMsgBurnToPower(
	from, to, clusterId, gatewayAddress, chatAddress string, burnAmoumt, useFreezeAmount sdk.Dec,
) *MsgBurnToPower {

	return &MsgBurnToPower{
		FromAddress:     from,
		ToAddress:       to,
		ClusterId:       clusterId,
		BurnAmount:      burnAmoumt,
		UseFreezeAmount: useFreezeAmount,
		GatewayAddress:  gatewayAddress,
		ChatAddress:     chatAddress,
	}
}

func (m MsgBurnToPower) Route() string { return sdk.MsgTypeURL(&m) }

func (m MsgBurnToPower) Type() string { return sdk.MsgTypeURL(&m) }

func (m MsgBurnToPower) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.GetFromAddress())
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{addr}
}
func (m MsgBurnToPower) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}
func (m MsgBurnToPower) ValidateBasic() error {

	_, err := sdk.AccAddressFromBech32(m.FromAddress)
	if err != nil {
		return sdkerrors.Wrap(err, "invalid sender address")
	}
	return nil
}

func NewMsgColonyRate(
	address, gatewayAddress string,
	rate []ColonyRate,
) *MsgColonyRate {

	return &MsgColonyRate{
		Address:        address,
		GatewayAddress: gatewayAddress,
		OnlineRate:     rate,
	}
}

func (m MsgColonyRate) Route() string { return sdk.MsgTypeURL(&m) }

func (m MsgColonyRate) Type() string { return sdk.MsgTypeURL(&m) }

func (m MsgColonyRate) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.Address)
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{addr}
}
func (m MsgColonyRate) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}
func (m MsgColonyRate) ValidateBasic() error {

	_, err := sdk.AccAddressFromBech32(m.Address)
	if err != nil {
		return sdkerrors.Wrap(err, "invalid sender address")
	}
	return nil
}

func NewMsgClusterPowerApprove(approveAddress, clusterId, approveEndBlock, fromAddress string) *MsgClusterPowerApprove {
	return &MsgClusterPowerApprove{
		ApproveAddress:  approveAddress,
		ClusterId:       clusterId,
		ApproveEndBlock: approveEndBlock,
		FromAddress:     fromAddress,
	}
}

func (m MsgClusterPowerApprove) Route() string { return sdk.MsgTypeURL(&m) }

func (m MsgClusterPowerApprove) Type() string { return sdk.MsgTypeURL(&m) }

func (m MsgClusterPowerApprove) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.FromAddress)
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{addr}
}
func (m MsgClusterPowerApprove) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}
func (m MsgClusterPowerApprove) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(m.FromAddress)
	if err != nil {
		return sdkerrors.Wrap(err, "invalid from address")
	}
	return nil
}

func NewMsgPersonDvmApprove(approveAddress, clusterId, approveEndBlock, fromAddress string) *MsgPersonDvmApprove {
	return &MsgPersonDvmApprove{
		ApproveAddress:  approveAddress,
		ClusterId:       clusterId,
		ApproveEndBlock: approveEndBlock,
		FromAddress:     fromAddress,
	}
}

func (m MsgPersonDvmApprove) Route() string { return sdk.MsgTypeURL(&m) }

func (m MsgPersonDvmApprove) Type() string { return sdk.MsgTypeURL(&m) }

func (m MsgPersonDvmApprove) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.FromAddress)
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{addr}
}
func (m MsgPersonDvmApprove) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}
func (m MsgPersonDvmApprove) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(m.FromAddress)
	if err != nil {
		return sdkerrors.Wrap(err, "invalid from address")
	}
	return nil
}
