package contracts

import (
	_ "embed" // embed compiled smart contract
	"encoding/json"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
)

var (
	//go:embed compiled_contracts/RedPacket.json
	redPacketJSON []byte

	RedPacketContract evmtypes.CompiledContract
)

// 0x9fc3338C63e0C9Df723d5a9242BdC4bc9d544aE2
func init() {
	err := json.Unmarshal(redPacketJSON, &RedPacketContract)
	if err != nil {
		panic(err)
	}
}
