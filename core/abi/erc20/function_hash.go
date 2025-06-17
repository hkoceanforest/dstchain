package erc20

import (
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
)


var FunctionHashTransfer = hexutil.Encode(ethcrypto.Keccak256([]byte("Transfer(address,address,uint256)")))
