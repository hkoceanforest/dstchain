package types

import (
	"freemasonry.cc/blockchain/core"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	_ sdk.Msg = &MsgAppTokenIssue{}
)

const (
	TypeMsgAppTokenIssuePath = "contract/AppTokenIssue"
	TypeMsgRegisterErc20     = "contract/RegisterErc20"
)

func NewMsgRegisterErc20(address, contractAddress, denom string, owner int32) *MsgRegisterErc20 {
	return &MsgRegisterErc20{
		Address:         address,
		ContractAddress: contractAddress,
		Denom:           denom,
		Owner:           owner,
	}
}

func (m MsgRegisterErc20) Route() string { return sdk.MsgTypeURL(&m) }
func (m MsgRegisterErc20) Type() string  { return sdk.MsgTypeURL(&m) }
func (m MsgRegisterErc20) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.Address)
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{addr}
}
func (m *MsgRegisterErc20) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}
func (m MsgRegisterErc20) ValidateBasic() error {
	if m.ContractAddress == "" {
		return core.ErrEmptyContractAddress
	}
	if m.Address != core.ContractAddressGov.String() {
		return core.SignAccountError
	}
	return nil
}


func NewMsgAppTokenIssue(
	fromAddress, name, symbol string,
	preMintAmount string,
	tokenDecimals string,
	logo string,
) *MsgAppTokenIssue {

	return &MsgAppTokenIssue{
		FromAddress:   fromAddress,
		Name:          name,
		Symbol:        symbol,
		PreMintAmount: preMintAmount,
		LogoUrl:       logo,
		Decimals:      tokenDecimals,
	}
}

func (m MsgAppTokenIssue) Route() string { return sdk.MsgTypeURL(&m) }
func (m MsgAppTokenIssue) Type() string  { return sdk.MsgTypeURL(&m) }
func (m MsgAppTokenIssue) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.FromAddress)
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{addr}
}
func (m *MsgAppTokenIssue) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}
func (m MsgAppTokenIssue) ValidateBasic() error {

	_, err := sdk.AccAddressFromBech32(m.FromAddress)
	if err != nil {
		return sdkerrors.Wrap(err, "invalid sender address")
	}
	if len(m.Name) > 255 || len(m.Symbol) > 255 || len(m.Decimals) > 255 || len(m.PreMintAmount) > 255 {
		return core.ParamsInvalidErr
	}
	if len(m.LogoUrl) > 500 {
		return core.ParamsInvalidErr
	}
	return nil
}
