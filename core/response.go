package core

import sdk "github.com/cosmos/cosmos-sdk/types"

type BaseResponse struct {
	Info   string      `json:"info"`
	Status int         `json:"status"`
	Data   interface{} `json:"data"`
}

type BroadcastTxResponse struct {
	Height      int64  `json:"height"`
	TxHash      string `json:"tx_hash"`
	CodeSpace   string `json:"code_space"`    
	Code        uint32 `json:"code"`          
	SignedTxStr string `json:"signed_tx_str"` 
}

type AccountNumberSeqResponse struct {
	AccountNumber uint64 `json:"account_number"` 
	Sequence      uint64 `json:"sequence"`
	NotFound      bool   `json:"not_found"` 
}

type BalanceResponse struct {
	Balances sdk.Coins `json:"balances"`
}
