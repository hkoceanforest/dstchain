package types

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

const (
	AttributeKeyReturnAmount = "return_amount"
)

type LogTransfer struct {
	From   common.Address
	To     common.Address
	Tokens *big.Int
}
