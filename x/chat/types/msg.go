package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	_ sdk.Msg = &MsgBurnGetMobile{}
	_ sdk.Msg = &MsgMobileTransfer{}
	_ sdk.Msg = &MsgSetChatInfo{}
)

const (
	TypeMsgMobileTransfer = "mobile_transfer"
	TypeMsgBurnGetMobile  = "burn_get_mobile"
	TypeMsgSetChatInfo    = "set_chat_info"
)


func NewMsgSetChatInfo(fromAddress, gatewayAddress, addressBook, chatBlacklist, chatWhitelist, chatBlacklistEnc, chatWhitelistEnc string, updateTime int64) *MsgSetChatInfo {
	return &MsgSetChatInfo{
		FromAddress:      fromAddress,
		GatewayAddress:   gatewayAddress,
		AddressBook:      addressBook,
		ChatBlacklist:    chatBlacklist,
		ChatWhitelist:    chatWhitelist,
		UpdateTime:       updateTime,
		ChatBlacklistEnc: chatBlacklistEnc,
		ChatWhitelistEnc: chatWhitelistEnc,
	}
}

func (m MsgSetChatInfo) Route() string { return sdk.MsgTypeURL(&m) }

func (m MsgSetChatInfo) Type() string { return sdk.MsgTypeURL(&m) }

func (m MsgSetChatInfo) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.FromAddress)
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{addr}
}
func (m MsgSetChatInfo) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}
func (m MsgSetChatInfo) ValidateBasic() error {

	_, err := sdk.AccAddressFromBech32(m.FromAddress)
	if err != nil {
		return sdkerrors.Wrap(err, "invalid sender address")
	}
	return nil
}

func NewMsgBurnGetMobile(fromAddress, mobilePrefix, gatewayAddress, chatAddress string) *MsgBurnGetMobile {
	return &MsgBurnGetMobile{
		FromAddress:    fromAddress,
		MobilePrefix:   mobilePrefix,
		GatewayAddress: gatewayAddress,
		ChatAddress:    chatAddress,
	}
}

func (m MsgBurnGetMobile) Route() string { return sdk.MsgTypeURL(&m) }

func (m MsgBurnGetMobile) Type() string { return sdk.MsgTypeURL(&m) }

func (m MsgBurnGetMobile) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.FromAddress)
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{addr}
}
func (m MsgBurnGetMobile) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}
func (m MsgBurnGetMobile) ValidateBasic() error {

	_, err := sdk.AccAddressFromBech32(m.FromAddress)
	if err != nil {
		return sdkerrors.Wrap(err, "invalid sender address")
	}
	return nil
}

func NewMsgMobileTransfer(fromAddress, toAddress, mobile string) *MsgMobileTransfer {
	return &MsgMobileTransfer{
		FromAddress: fromAddress,
		ToAddress:   toAddress,
		Mobile:      mobile,
	}
}

func (m MsgMobileTransfer) Route() string { return sdk.MsgTypeURL(&m) }

func (m MsgMobileTransfer) Type() string { return sdk.MsgTypeURL(&m) }

func (m MsgMobileTransfer) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.FromAddress)
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{addr}
}
func (m MsgMobileTransfer) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}
func (m MsgMobileTransfer) ValidateBasic() error {

	_, err := sdk.AccAddressFromBech32(m.FromAddress)
	if err != nil {
		return sdkerrors.Wrap(err, "invalid sender address")
	}
	return nil
}
