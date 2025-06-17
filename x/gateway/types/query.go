package types

import sdk "github.com/cosmos/cosmos-sdk/types"

const (
	
	QueryGatewayInfo = "gateway_info"
	
	QueryGatewayList = "gateway_list"
	
	QueryGatewayNum = "gateway_num"
	
	QueryGatewayRedeemNum = "gateway_redeem_num"
	
	QueryValidatorByConsAddress = "validatorByConsAddress"
	
	QueryGatewayNumberCount = "gateway_number_count"
	
	QueryGatewayNumberUnbondCount = "gateway_unbond_number_count"
	
	QueryGasPrice = "gas_price"
	
	QueryParams = "params"
	
	QueryValidators = "validators"
	
	QueryGatewayUpload = "gateway_upload"
	
	QueryDelegateLastHeightKey = "delegate_last_height"
)


type QueryGatewayInfoParams struct {
	GatewayAddress  string `json:"gateway_address"`
	GatewayNumIndex string `json:"gateway_num_index"`
}


type GatewayNumberCountParams struct {
	GatewayAddress string   `json:"gateway_address"`
	Amount         sdk.Coin `json:"amount"`
}


type QueryValidatorByConsAddrParams struct {
	ValidatorConsAddress sdk.ConsAddress
}

type QueryDelegateLastHeight struct {
	ValidatorAddress string `json:"validator_address"`
	DelegateAddress  string `json:"delegate_address"`
}
