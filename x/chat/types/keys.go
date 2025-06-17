package types

import (
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/ethereum/go-ethereum/common"
)

const (
	
	ModuleName     = "chat"
	ModuleBurnName = "chat_burn"
	
	StoreKey = ModuleName

	
	RouterKey = ModuleName
)

var ModuleAddress common.Address

func init() {
	ModuleAddress = common.BytesToAddress(authtypes.NewModuleAddress(ModuleName).Bytes())
}

const (
	
	MobileSuffixLength = 4 

	
	MobileSuffixMax = 9999

	
	KeyPrefixRegisterInfo = "chat_register_info_"

	
	KeyPrefixChatSendGift = "chat_send_gift_"

	
	KeyPrefixMobileOwner = "chat_mobile_owner_"

	
	KeyPrefixGatewayIssueToken = "gateway_issue_token_"

	
	KeyChatAddress = "chat_address_"
)
