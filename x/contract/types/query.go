package types

const (
	
	QueryParams = "params"
	
	QueryContractCode = "contract_code"

	QueryExchangeContractAddress = "exchange_contract_address"

	QueryTokenPair = "token_pair"

	QueryMainTokenBalances = "main_token_balance"
)

type QueryNftInfoParams struct {
	Address         string `json:"address"`          
	ContractAddress string `json:"contract_address"` 
}


type QueryTokenIsEnough struct {
	FromAddress    string `json:"from_address"`
	GatewayAddress string `json:"gateway_address"`
}

type QueryCrossChainHashParams struct {
	Address string `json:"address"`
	Hash    string `json:"hash"`
}
