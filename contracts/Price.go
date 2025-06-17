package contracts

import (
	_ "embed" // embed compiled smart contract
	"encoding/json"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
)

var (
	//go:embed compiled_contracts/Price.json
	priceJSON []byte

	PriceJSONContract evmtypes.CompiledContract
)

func init() {
	err := json.Unmarshal(priceJSON, &PriceJSONContract)
	if err != nil {
		panic(err)
	}
}
