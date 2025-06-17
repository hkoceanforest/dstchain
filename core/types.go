package core

import (
	"errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/migrations/legacytx"
	"github.com/shopspring/decimal"
	abciTypes "github.com/tendermint/tendermint/abci/types"
)


type RealCoin struct {
	Denom  string `json:"denom"`
	Amount string `json:"amount"`
}

func (this RealCoin) Add(rcoin2 RealCoin) RealCoin {
	dec1 := sdk.MustNewDecFromStr(rcoin2.Amount)
	dec2 := sdk.MustNewDecFromStr(this.Amount)
	if this.Denom != rcoin2.Denom {
		panic(errors.New("added realcoin denom It has to be the same"))
	}
	return RealCoin{Amount: dec1.Add(dec2).String(), Denom: this.Denom}
}
func (this RealCoin) AmountDec() decimal.Decimal {
	if this.Amount == "" {
		return decimal.Decimal{}
	} else {
		ret, err := decimal.NewFromString(this.Amount)
		if err != nil {

		}
		return ret
	}
}

func (this RealCoin) String() string {
	return this.Amount + this.Denom
}


func (this RealCoin) FormatAmount() string {
	return RemoveStringLastZero(this.Amount)
}

func (this RealCoin) FormatString() string {
	return this.FormatAmount() + this.Denom
}


type RealCoins []RealCoin

func (this RealCoins) String() string {
	tmp := ""
	for _, v := range this {
		tmp += v.String() + " "
	}
	return tmp
}

func (this RealCoins) Get(index int) RealCoin {
	return this[index]
}


type BlockMessageI interface {
	Type() string 
	Unmarshal([]byte) (interface{}, error)
	Marshal() []byte
}

type TxBase struct {
	BlockTime int64           `json:"block_time"` 
	Height    int64           `json:"height"`
	TxHash    string          `json:"tx_hash"`
	Fee       legacytx.StdFee `json:"fee"`  
	Memo      string          `json:"memo"` 
}


type TradeBase struct {
	TradeType       TranserType `json:"trade_type"`       
	ContractAddress string      `json:"contract_address"` 
}

type TxPropValue struct {
	BlockTime int64 
	TxHash    string
	Height    int64
	Seq       int64
	Fee       legacytx.StdFee
	Memo      string
	Events    []abciTypes.Event
}

func (this *TxBase) UpdateTxBase(prop *TxPropValue) {
	this.BlockTime = prop.BlockTime
	this.Height = prop.Height
	this.TxHash = prop.TxHash
	this.Fee = prop.Fee
	this.Memo = prop.Memo
}

func (this *TradeBase) UpdateTradeBase(tradeType TranserType, contractAddr string) {
	this.TradeType = tradeType
	this.ContractAddress = contractAddr
}



type TransferData struct {
	TxBase
	TradeBase
	FromAddress string    `json:"from_address"` 
	ToAddress   string    `json:"to_address"`   
	Coins       sdk.Coins `json:"coins"`
	FromBalance string    `json:"from_balance"`
	ToBalance   string    `json:"to_balance"`
}
