package types

import (
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/ethereum/go-ethereum/common"
)

const (
	
	ModuleName = "gateway"
	
	StoreKey = ModuleName

	QuerierRoute = ModuleName

	
	RouterKey = ModuleName
)

var ModuleAddress common.Address

func init() {
	ModuleAddress = common.BytesToAddress(authtypes.NewModuleAddress(ModuleName).Bytes())
}

const (
	KeyPrefixAddressBook = "comm_address_book"

	DelegateLastTimeKey = "delegate_last_time_"

	GatewayUploadKey = "gateway_upload_"
)
