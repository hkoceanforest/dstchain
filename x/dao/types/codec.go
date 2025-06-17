package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {

	cdc.RegisterConcrete(&MsgColonyRate{}, TypeMsgColonyRate, nil)
	cdc.RegisterConcrete(&MsgBurnToPower{}, TypeMsgBurnToPower, nil)
	cdc.RegisterConcrete(&MsgClusterChangeSalaryRatio{}, TypeMsgClusterChangeSalaryRatio, nil)
	cdc.RegisterConcrete(&MsgClusterChangeDvmRatio{}, TypeMsgClusterChangeDvmRatio, nil)
	cdc.RegisterConcrete(&MsgCreateCluster{}, TypeMsgCreateCluster, nil)
	cdc.RegisterConcrete(&MsgClusterAddMembers{}, TypeMsgClusterAddMembers, nil)
	cdc.RegisterConcrete(&MsgClusterChangeId{}, TypeMsgClusterChangeid, nil)
	cdc.RegisterConcrete(&MsgWithdrawOwnerReward{}, TypeMsgWithdrawOwnerReward, nil)
	cdc.RegisterConcrete(&MsgWithdrawSwapDpos{}, TypeMsgWithdrawSwapDpos, nil)
	cdc.RegisterConcrete(&MsgWithdrawDeviceReward{}, TypeMsgWithdrawDeviceReward, nil)
	cdc.RegisterConcrete(&MsgDeleteMembers{}, TypeMsgDeleteMembers, nil)
	cdc.RegisterConcrete(&MsgThawFrozenPower{}, TypeMsgThawFrozenPower, nil)
	cdc.RegisterConcrete(&MsgClusterMemberExit{}, TypeMsgClusterMemberExit, nil)
	cdc.RegisterConcrete(&MsgClusterChangeName{}, TypeMsgClusterChangeName, nil)
	cdc.RegisterConcrete(&MsgUpdateAdmin{}, TypeMsgUpdateAdmin, nil)
	cdc.RegisterConcrete(&MsgClusterPowerApprove{}, TypeMsgClusterPowerApprove, nil)
	cdc.RegisterConcrete(&MsgPersonDvmApprove{}, TypePersonDvmApprove, nil)
	cdc.RegisterConcrete(&MsgClusterAd{}, TypeMsgClusterAd, nil)
	cdc.RegisterConcrete(&MsgReceiveBurnRewardFee{}, TypeMsgReceiveBurnRewardFee, nil)
	cdc.RegisterConcrete(&MsgClusterChangeDaoRatio{}, TypeMsgClusterChangeDaoRatio, nil)
	cdc.RegisterConcrete(&MsgAgreeJoinCluster{}, TypeMsgAgreeJoinCluster, nil)
	cdc.RegisterConcrete(&MsgRedPacket{}, TypeMsgRedPacket, nil)
	cdc.RegisterConcrete(&MsgOpenRedPacket{}, TypeMsgOpenRedPacket, nil)
	cdc.RegisterConcrete(&MsgReturnRedPacket{}, TypeMsgReturnRedPacket, nil)
	cdc.RegisterConcrete(&MsgCreateClusterAddMembers{}, TypeMsgCreateClusterAddMembers, nil)
	cdc.RegisterConcrete(&MsgAgreeJoinClusterApply{}, TypeNewMsgAgreeJoinClusterApply, nil)
	cdc.RegisterConcrete(&MsgStartPowerRewardRedeem{}, TypeNewMsgStartPowerRewardRedeem, nil)
	cdc.RegisterConcrete(&MsgReceivePowerCutReward{}, TypeNewMsgReceivePowerCutReward, nil)

}

func RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&MsgColonyRate{},
		&MsgBurnToPower{},
		&MsgClusterChangeSalaryRatio{},
		&MsgClusterChangeDvmRatio{},
		&MsgCreateCluster{},
		&MsgClusterAddMembers{},
		&MsgClusterChangeId{},
		&MsgWithdrawOwnerReward{},
		&MsgWithdrawSwapDpos{},
		&MsgWithdrawDeviceReward{},
		&MsgDeleteMembers{},
		&MsgThawFrozenPower{},
		&MsgClusterMemberExit{},
		&MsgClusterChangeName{},
		&MsgUpdateAdmin{},
		&MsgClusterPowerApprove{},
		&MsgPersonDvmApprove{},
		&MsgClusterAd{},
		&MsgReceiveBurnRewardFee{},
		&MsgClusterChangeDaoRatio{},
		&MsgAgreeJoinCluster{},
		&MsgRedPacket{},
		&MsgOpenRedPacket{},
		&MsgReturnRedPacket{},
		&MsgCreateClusterAddMembers{},
		&MsgAgreeJoinClusterApply{},
		&MsgStartPowerRewardRedeem{},
		&MsgReceivePowerCutReward{},
	)
	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewAminoCodec(amino)
)

func init() {
	RegisterLegacyAminoCodec(amino)
	cryptocodec.RegisterCrypto(amino)
	amino.Seal()
}
