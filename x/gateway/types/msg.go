package types

import (
	"bytes"
	"freemasonry.cc/blockchain/core"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"regexp"
)

var (
	_ sdk.Msg = &MsgGatewayRegister{}
	_ sdk.Msg = &MsgGatewayIndexNum{}
	_ sdk.Msg = &MsgGatewayUndelegate{}
)

const (
	TypeMsgCreateSmartValidator   = "create_smart_validator"
	TypeMsgGatewayRegister        = "gateway_register"
	TypeMsgGatewayEdit            = "gateway_edit"
	TypeMsgGatewayIndexNum        = "gateway_index_num"
	TypeMsgGatewayUndelegation    = "gateway_undelegation"
	TypeMsgGatewayBeginRedelegate = "gateway_begin_redelegate"
	TypeMsgGatewayUpload          = "gateway_upload"
)

func NewMsgGatewayUpload(address string, gatewayKeyInfo []byte) *MsgGatewayUpload {
	return &MsgGatewayUpload{
		FromAddress:    address,
		GatewayKeyInfo: gatewayKeyInfo,
	}
}

func (msg MsgGatewayUpload) Route() string { return RouterKey }
func (msg MsgGatewayUpload) Type() string  { return TypeMsgGatewayUpload }
func (msg MsgGatewayUpload) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.FromAddress)
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{addr}
}
func (msg *MsgGatewayUpload) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}
func (msg MsgGatewayUpload) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.FromAddress)
	if err != nil {
		return sdkerrors.Wrap(err, "invalid send address")
	}
	return nil
}

func NewMsgCreateSmartValidator(
	valAddr sdk.ValAddress, pubKey string, 
	selfDelegation sdk.Coin, description stakingtypes.Description, commission stakingtypes.CommissionRates, minSelfDelegation sdk.Int,
) (*MsgCreateSmartValidator, error) {
	return &MsgCreateSmartValidator{
		Description:       description,
		DelegatorAddress:  sdk.AccAddress(valAddr).String(),
		ValidatorAddress:  valAddr.String(),
		PubKey:            pubKey,
		Value:             selfDelegation,
		Commission:        commission,
		MinSelfDelegation: minSelfDelegation,
	}, nil
}

func (msg MsgCreateSmartValidator) Route() string { return RouterKey }

func (msg MsgCreateSmartValidator) Type() string { return TypeMsgCreateSmartValidator }

func (msg MsgCreateSmartValidator) GetSigners() []sdk.AccAddress {
	
	delAddr, err := sdk.AccAddressFromBech32(msg.DelegatorAddress)
	if err != nil {
		panic(err)
	}
	addrs := []sdk.AccAddress{delAddr}
	addr, err := sdk.ValAddressFromBech32(msg.ValidatorAddress)
	if err != nil {
		panic(err)
	}
	if !bytes.Equal(delAddr.Bytes(), addr.Bytes()) {
		addrs = append(addrs, sdk.AccAddress(addr))
	}

	return addrs
}

func (msg MsgCreateSmartValidator) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgCreateSmartValidator) ValidateBasic() error {
	
	delAddr, err := sdk.AccAddressFromBech32(msg.DelegatorAddress)
	if err != nil {
		return err
	}
	if delAddr.Empty() {
		return stakingtypes.ErrEmptyDelegatorAddr
	}

	if msg.ValidatorAddress == "" {
		return stakingtypes.ErrEmptyValidatorAddr
	}

	valAddr, err := sdk.ValAddressFromBech32(msg.ValidatorAddress)
	if err != nil {
		return err
	}
	if !sdk.AccAddress(valAddr).Equals(delAddr) {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "validator address is invalid")
	}

	if msg.PubKey == "" {
		return stakingtypes.ErrEmptyValidatorPubKey
	}

	if !msg.Value.IsValid() || !msg.Value.Amount.IsPositive() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "invalid delegation amount")
	}

	if msg.Description == (stakingtypes.Description{}) {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "empty description")
	}

	if msg.Commission == (stakingtypes.CommissionRates{}) {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "empty commission")
	}

	if err := msg.Commission.Validate(); err != nil {
		return err
	}

	if !msg.MinSelfDelegation.IsPositive() {
		return sdkerrors.Wrap(
			sdkerrors.ErrInvalidRequest,
			"minimum self delegation must be a positive integer",
		)
	}

	if msg.Value.Amount.LT(msg.MinSelfDelegation) {
		return stakingtypes.ErrSelfDelegationBelowMinimum
	}
	return nil
}


func NewMsgGatewayRegister(address, gatewayUrl, delegation, packageName, PeerId, machineAddress string, indexNumber []string) *MsgGatewayRegister {
	return &MsgGatewayRegister{
		Address:        address,
		GatewayUrl:     gatewayUrl,
		Delegation:     delegation,
		IndexNumber:    indexNumber,
		Package:        packageName,
		PeerId:         PeerId,
		MachineAddress: machineAddress,
	}
}

func (msg MsgGatewayRegister) Route() string { return RouterKey }
func (msg MsgGatewayRegister) Type() string  { return TypeMsgGatewayRegister }
func (msg MsgGatewayRegister) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Address)
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{addr}
}
func (msg *MsgGatewayRegister) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}
func (msg MsgGatewayRegister) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Address)
	if err != nil {
		return sdkerrors.Wrap(err, "invalid send address")
	}
	for _, val := range msg.IndexNumber {
		if _, ok := core.NumIndexAmount[len(val)]; !ok {
			return core.ErrGatewayNumLength
		}
	}
	matchedIp, err := regexp.MatchString("^http[s]?://(\\d|[1-9]\\d|1\\d{2}|2[0-4]\\d|25[0-5])\\.(\\d|[1-9]\\d|1\\d{2}|2[0-4]\\d|25[0-5])\\.(\\d|[1-9]\\d|1\\d{2}|2[0-4]\\d|25[0-5])\\.(\\d|[1-9]\\d|1\\d{2}|2[0-4]\\d|25[0-5]):50327$", msg.GatewayUrl)
	if err != nil {
		return sdkerrors.Wrap(err, "gateway url")
	}
	matchedDomain, err := regexp.MatchString("^(http[s]?://)([a-zA-Z0-9-_]+\\.)*[a-zA-Z0-9][a-zA-Z0-9-_]+\\.[a-zA-Z]{2,63}:50327$", msg.GatewayUrl)
	if err != nil {
		return sdkerrors.Wrap(err, "gateway url")
	}
	if !matchedIp && !matchedDomain {
		return sdkerrors.Wrap(core.ErrGatewayUrl, "gateway url not match")
	}
	matched, err := regexp.MatchString("^[a-z][a-z\\d]{0,31}\\.[a-z][a-z\\d]{2}$", msg.Package)
	if err != nil {
		return sdkerrors.Wrap(err, "gateway package")
	}
	if !matched {
		return sdkerrors.Wrap(core.ErrGatewayPackage, "gateway package not match")
	}
	return nil
}


func NewMsgGatewayEdit(address, gatewayUrl string) *MsgGatewayEdit {
	return &MsgGatewayEdit{
		Address:    address,
		GatewayUrl: gatewayUrl,
	}
}

func (msg MsgGatewayEdit) Route() string { return RouterKey }
func (msg MsgGatewayEdit) Type() string  { return TypeMsgGatewayEdit }
func (msg MsgGatewayEdit) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Address)
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{addr}
}
func (msg *MsgGatewayEdit) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}
func (msg MsgGatewayEdit) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Address)
	if err != nil {
		return sdkerrors.Wrap(err, "invalid send address")
	}
	matched, err := regexp.MatchString("^(\\d|[1-9]\\d|1\\d{2}|2[0-4]\\d|25[0-5])\\.(\\d|[1-9]\\d|1\\d{2}|2[0-4]\\d|25[0-5])\\.(\\d|[1-9]\\d|1\\d{2}|2[0-4]\\d|25[0-5])\\.(\\d|[1-9]\\d|1\\d{2}|2[0-4]\\d|25[0-5]):([0-9]|[1-9]\\d|[1-9]\\d{2}|[1-9]\\d{3}|[1-5]\\d{4}|6[0-4]\\d{3}|65[0-4]\\d{2}|655[0-2]\\d|6553[0-5])$", msg.GatewayUrl)
	if err != nil {
		return sdkerrors.Wrap(err, "gateway url")
	}
	if !matched {
		return sdkerrors.Wrap(err, "gateway url not match")
	}
	return nil
}


func NewMsgGatewayIndexNum(address, validatorAddress string, indexNumber []string) *MsgGatewayIndexNum {
	return &MsgGatewayIndexNum{
		DelegatorAddress: address,
		ValidatorAddress: validatorAddress,
		IndexNumber:      indexNumber,
	}
}
func (msg MsgGatewayIndexNum) Route() string { return RouterKey }
func (msg MsgGatewayIndexNum) Type() string  { return TypeMsgGatewayIndexNum }
func (msg MsgGatewayIndexNum) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.DelegatorAddress)
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{addr}
}
func (msg *MsgGatewayIndexNum) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}
func (msg MsgGatewayIndexNum) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.DelegatorAddress)
	if err != nil {
		return sdkerrors.Wrap(err, "invalid send address")
	}
	if len(msg.IndexNumber) == 0 {
		return sdkerrors.Wrap(err, "invalid message")
	}
	for _, val := range msg.IndexNumber {
		if _, ok := core.NumIndexAmount[len(val)]; !ok {
			return core.ErrGatewayNumLength
		}
	}
	return nil
}


func NewMsgGatewayUndelegation(address, validatorAddress string, amount sdk.Coin, indexNumber []string) *MsgGatewayUndelegate {
	return &MsgGatewayUndelegate{
		DelegatorAddress: address,
		ValidatorAddress: validatorAddress,
		Amount:           amount,
		IndexNumber:      indexNumber,
	}
}

func (msg MsgGatewayUndelegate) Route() string { return RouterKey }
func (msg MsgGatewayUndelegate) Type() string  { return TypeMsgGatewayUndelegation }
func (msg MsgGatewayUndelegate) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.DelegatorAddress)
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{addr}
}
func (msg *MsgGatewayUndelegate) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}
func (msg MsgGatewayUndelegate) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.DelegatorAddress)
	if err != nil {
		return sdkerrors.Wrap(err, "invalid send address")
	}
	if len(msg.IndexNumber) == 0 && msg.Amount.IsZero() {
		return sdkerrors.Wrap(err, "invalid message")
	}
	for _, val := range msg.IndexNumber {
		if _, ok := core.NumIndexAmount[len(val)]; !ok {
			return core.ErrGatewayNumLength
		}
	}

	return nil
}

func NewMsgGatewayBeginRedelegate(
	delAddr sdk.AccAddress, valSrcAddr, valDstAddr sdk.ValAddress, amount sdk.Coin, indexNumber []string) *MsgGatewayBeginRedelegate {
	return &MsgGatewayBeginRedelegate{
		DelegatorAddress:    delAddr.String(),
		ValidatorSrcAddress: valSrcAddr.String(),
		ValidatorDstAddress: valDstAddr.String(),
		Amount:              amount,
		IndexNumber:         indexNumber,
	}
}

func (msg MsgGatewayBeginRedelegate) Route() string { return RouterKey }

func (msg MsgGatewayBeginRedelegate) Type() string { return TypeMsgGatewayBeginRedelegate }

func (msg MsgGatewayBeginRedelegate) GetSigners() []sdk.AccAddress {
	delAddr, err := sdk.AccAddressFromBech32(msg.DelegatorAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{delAddr}
}

func (msg MsgGatewayBeginRedelegate) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgGatewayBeginRedelegate) ValidateBasic() error {
	if msg.DelegatorAddress == "" {
		return stakingtypes.ErrEmptyDelegatorAddr
	}

	if msg.ValidatorSrcAddress == "" {
		return stakingtypes.ErrEmptyValidatorAddr
	}

	if msg.ValidatorDstAddress == "" {
		return stakingtypes.ErrEmptyValidatorAddr
	}

	if !msg.Amount.IsValid() || !msg.Amount.Amount.IsPositive() {
		return sdkerrors.Wrap(
			sdkerrors.ErrInvalidRequest,
			"invalid shares amount",
		)
	}
	for _, val := range msg.IndexNumber {
		if len(val) != 7 {
			return core.ErrGatewayNumLength
		}
	}
	return nil
}
