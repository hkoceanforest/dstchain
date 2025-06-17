package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {

	cdc.RegisterConcrete(&MsgCreateSmartValidator{}, MSG_SMART_CREATE_VALIDATOR, nil)
	cdc.RegisterConcrete(&MsgGatewayRegister{}, MSG_GATEWAY_REGISTER, nil)
	cdc.RegisterConcrete(&MsgGatewayEdit{}, MSG_GATEWAY_Edit, nil)
	cdc.RegisterConcrete(&MsgGatewayIndexNum{}, MSG_GATEWAY_INDEX_NUM, nil)
	cdc.RegisterConcrete(&MsgGatewayUndelegate{}, MSG_GATEWAY_UNDELEGATION, nil)
	cdc.RegisterConcrete(&MsgGatewayBeginRedelegate{}, MSG_GATEWAY_BEGIN_REDELEGATE, nil)
	cdc.RegisterConcrete(&MsgGatewayUpload{}, MSG_GATEWAY_UPLOAD, nil)
	cdc.RegisterConcrete(&MsgEmptyResponse{}, MSG_EMPTY_RESPONSE, nil)
}

func RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&MsgCreateSmartValidator{},
		&MsgGatewayRegister{},
		&MsgGatewayEdit{},
		&MsgGatewayIndexNum{},
		&MsgGatewayUndelegate{},
		&MsgGatewayBeginRedelegate{},
		&MsgGatewayUpload{},
	)
	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	amino = codec.NewLegacyAmino()

	ModuleCdc = codec.NewAminoCodec(amino)
)

func init() {
	RegisterLegacyAminoCodec(amino)
	cryptocodec.RegisterCrypto(amino)
	amino.Seal()
}
