package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)


func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	
	cdc.RegisterConcrete(&MsgMobileTransfer{}, MsgTypeMobileTransfer, nil)
	cdc.RegisterConcrete(&MsgBurnGetMobile{}, MsgTypeBurnGetMobile, nil)
	cdc.RegisterConcrete(&MsgSetChatInfo{}, MsgTypeSetChatInfo, nil)
}

func RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgMobileTransfer{},
		&MsgBurnGetMobile{},
		&MsgSetChatInfo{},
	)
	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewAminoCodec(amino)
)
