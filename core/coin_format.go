package core

import (
	"errors"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/migrations/legacytx"
	"strings"
)

func StringToCoinWithRate(amStr string) (sdk.Coin, error) {
	decCoin, err := sdk.ParseDecCoin(amStr)
	if err != nil {
		return sdk.Coin{}, err
	}
	dec := decCoin.Amount 
	decInt := dec.TruncateInt()
	coin := sdk.NewCoin(decCoin.Denom, decInt)
	return coin, nil
}

func StringToCoinsWithRate(amStr string) (sdk.Coins, error) {
	decCoins, err := sdk.ParseDecCoins(amStr)
	if err != nil {
		return sdk.Coins{}, err
	}
	coins := sdk.Coins{}
	for _, decCoin := range decCoins {
		dec := decCoin.Amount
		decInt := dec.TruncateInt()
		coin := sdk.NewCoin(decCoin.Denom, decInt)
		coins = append(coins, coin)
	}
	return coins, nil
}

func MustRealString2LedgerInt(realString string) (ledgerInt sdk.Int) {
	realIntAmountDec := sdk.MustNewDecFromStr(realString)
	if realIntAmountDec.LT(MinRealAmountDec) && !realIntAmountDec.IsZero() {
		realIntAmountDec = MinRealAmountDec
	}
	return realIntAmountDec.Mul(RealToLedgerRateDec).TruncateInt()
}

func MustRealString2LedgerIntNoMin(realString string) (ledgerInt sdk.Int) {
	realIntAmountDec := sdk.MustNewDecFromStr(realString)
	return realIntAmountDec.Mul(RealToLedgerRateDec).TruncateInt()
}

func RealString2LedgerCoin(realString, denom string) sdk.Coin {
	return sdk.NewCoin(denom, MustRealString2LedgerInt(realString))
}

func MustLedgerInt2RealString(ledgerInt sdk.Int) string {
	realCoinAmount := sdk.NewDecFromInt(ledgerInt)
	rate, err := sdk.NewDecFromStr(LedgerToRealRate)
	if err != nil {
		panic(err)
	}
	return RemoveStringLastZero(realCoinAmount.Mul(rate).String())
}

func MustParseLedgerDec(ledgerDec sdk.Dec) (realAmount string) {
	rate, err := sdk.NewDecFromStr(LedgerToRealRate)
	if err != nil {
		panic(err)
	}
	return RemoveStringLastZero(ledgerDec.Mul(rate).String())
}

func NewLedgerInt(realAmount float64) sdk.Int {
	if realAmount < MinRealAmountFloat64 {
		return sdk.NewInt(int64(MinRealAmountFloat64 * RealToLedgerRate))
	}
	return sdk.NewInt(int64(realAmount * RealToLedgerRate))
}

func NewLedgerDec(realAmount float64) sdk.Dec {
	ledgerInt := NewLedgerInt(realAmount)
	return sdk.NewDecFromInt(ledgerInt)
}

func NewLedgerCoin(realAmount float64) sdk.Coin {
	ledgerInt := NewLedgerInt(realAmount)
	return sdk.NewCoin(BaseDenom, ledgerInt)
}

func NewLedgerDecCoin(realAmount float64) sdk.DecCoin {
	ledgerInt := NewLedgerInt(realAmount)
	return sdk.NewDecCoin(BaseDenom, ledgerInt)
}

func NewLedgerCoins(realAmount float64) sdk.Coins {
	ledgerInt := NewLedgerInt(realAmount)
	return sdk.NewCoins(sdk.NewCoin(BaseDenom, ledgerInt))
}

func NewLedgerFeeFromGas(gas uint64, amount float64) legacytx.StdFee {
	ledgerInt := NewLedgerInt(amount)
	fee := legacytx.NewStdFee(gas, sdk.NewCoins(sdk.NewCoin(BaseDenom, ledgerInt)))
	return fee
}

func NewLedgerFee(amount float64) legacytx.StdFee {
	ledgerInt := NewLedgerInt(amount)
	fee := legacytx.NewStdFee(flags.DefaultGasLimit, sdk.NewCoins(sdk.NewCoin(BaseDenom, ledgerInt)))
	return fee
}


func NewLedgerFeeZero() legacytx.StdFee {
	fee := legacytx.NewStdFee(flags.DefaultGasLimit, sdk.NewCoins(sdk.NewInt64Coin(BaseDenom, 0)))
	return fee
}

func MustParseLedgerCoins(ledgerCoins sdk.Coins) (realAmounts string) {
	for _, coin := range ledgerCoins {
		if realAmounts != "" {
			realAmounts += ","
		}
		realAmounts += MustParseLedgerCoin(coin) + coin.Denom
	}
	return
}

func MustParseLedgerFee(ledgerFee legacytx.StdFee) (realAmount string) {
	return MustParseLedgerCoins(ledgerFee.Amount)
}

func MustParseLedgerDec2(ledgerDec sdk.Dec) (realAmount sdk.Dec) {
	rate, err := sdk.NewDecFromStr(LedgerToRealRate)
	if err != nil {
		panic(err)
	}
	return ledgerDec.Mul(rate)
}

func MustParseLedgerCoinFromStr(ledgerCoinStr string) (realAmount string) {
	ledgerCoin, err := sdk.ParseCoinNormalized(ledgerCoinStr)
	if err != nil {
		panic(err)
	}
	ledgerAmount := sdk.NewDecFromInt(ledgerCoin.Amount)

	rate, err := sdk.NewDecFromStr(LedgerToRealRate)
	if err != nil {
		panic(err)
	}
	return RemoveStringLastZero(ledgerAmount.Mul(rate).String())
}

func MustParseLedgerCoin(ledgerCoin sdk.Coin) (realAmount string) {
	ledgerAmount := sdk.NewDecFromInt(ledgerCoin.Amount)
	rate, err := sdk.NewDecFromStr(LedgerToRealRate)
	if err != nil {
		panic(err)
	}
	return RemoveStringLastZero(ledgerAmount.Mul(rate).String())
}

func MustParseLedgerDecCoin(ledgerDecCoin sdk.DecCoin) (realAmount string) {
	ledgerAmount := ledgerDecCoin.Amount
	rate, err := sdk.NewDecFromStr(LedgerToRealRate)
	if err != nil {
		panic(err)
	}
	return RemoveStringLastZero(ledgerAmount.Mul(rate).String())
}

func MustParseLedgerDecCoins(ledgerDecCoins sdk.DecCoins) (realAmount string) {
	if len(ledgerDecCoins) > 0 {
		return MustParseLedgerDecCoin(ledgerDecCoins[0])
	}
	return ""
}

func MustLedgerCoin2RealCoin(ledgerCoin sdk.Coin) (realCoin RealCoin) {
	ledgerAmount := sdk.NewDecFromInt(ledgerCoin.Amount)
	rate, err := sdk.NewDecFromStr(LedgerToRealRate)
	if err != nil {
		panic(err)
	}
	return RealCoin{
		Denom:  ledgerCoin.Denom,
		Amount: RemoveStringLastZero(ledgerAmount.Mul(rate).String()),
	}
}

func MustLedgerDecCoin2RealCoin(ledgerDecCoin sdk.DecCoin) (realCoin RealCoin) {
	ledgerAmount := ledgerDecCoin.Amount
	rate, err := sdk.NewDecFromStr(LedgerToRealRate)
	if err != nil {
		panic(err)
	}
	return RealCoin{
		Denom:  ledgerDecCoin.Denom,
		Amount: RemoveStringLastZero(ledgerAmount.Mul(rate).String()),
	}
}

func MustLedgerDecCoins2RealCoins(ledgerDecCoins sdk.DecCoins) (realCoins RealCoins) {
	for i := 0; i < len(ledgerDecCoins); i++ {
		realCoins = append(realCoins, MustLedgerDecCoin2RealCoin(ledgerDecCoins[0]))
	}
	return
}

func MustLedgerCoins2RealCoins(ledgerCoins sdk.Coins) (realCoins RealCoins) {
	for i := 0; i < len(ledgerCoins); i++ {
		realCoins = append(realCoins, MustLedgerCoin2RealCoin(ledgerCoins[i]))
	}
	return
}

func MustRealCoin2LedgerDecCoin(realCoin RealCoin) (ledgerDecCoin sdk.DecCoin) {
	realCoinAmount, err := sdk.NewDecFromStr(realCoin.Amount)
	if err != nil {
		panic(err)
	}
	rate := sdk.NewDec(RealToLedgerRateInt64)
	return sdk.NewDecCoinFromDec(realCoin.Denom, realCoinAmount.Mul(rate))
}

func MustRealCoins2LedgerDecCoins(realCoins RealCoins) (ledgerDecCoins sdk.DecCoins) {
	for i := 0; i < len(realCoins); i++ {
		ledgerDecCoins = append(ledgerDecCoins, MustRealCoin2LedgerDecCoin(realCoins[0]))
	}
	return
}



func MustRealCoin2LedgerCoin(realCoin RealCoin) (ledgerCoin sdk.Coin) {
	realCoinAmount, err := sdk.NewDecFromStr(realCoin.Amount)
	if err != nil {
		panic(err)
	}
	rate := sdk.NewDec(RealToLedgerRateInt64)
	return sdk.NewCoin(realCoin.Denom, realCoinAmount.Mul(rate).TruncateInt())
}

func MustRealCoins2LedgerCoins(realCoins RealCoins) (ledgerCoins sdk.Coins) {
	for i := 0; i < len(realCoins); i++ {
		ledgerCoins = append(ledgerCoins, MustRealCoin2LedgerCoin(realCoins[i])).Sort()
	}
	return
}

func NewRealCoinFromStr(denom string, amount string) RealCoin {
	return RealCoin{Denom: denom, Amount: amount}
}

func NewRealCoinsFromStr(denom string, amount string) RealCoins {
	return RealCoins{NewRealCoinFromStr(denom, amount)}
}

func NewRealCoins(realCoin RealCoin) RealCoins {
	return RealCoins{realCoin}
}


func RemoveStringLastZero(balance string) string {
	if !strings.Contains(balance, ".") {
		return balance
	}
	dataList := strings.Split(balance, ".")
	zhengshu := dataList[0]
	xiaoshu := dataList[1]
	if len(dataList[1]) > 18 {
		xiaoshu = xiaoshu[:18]
	}
	xiaoshu2 := ""
	for i := len(xiaoshu) - 1; i >= 0; i-- {
		if xiaoshu[i] != '0' {
			xiaoshu2 = xiaoshu[:i+1]
			break
		}
	}
	if xiaoshu2 == "" {
		return zhengshu
	} else {
		return zhengshu + "." + xiaoshu2
	}
}


func RemoveDecLastZero(amount sdk.Dec) string {
	balance := amount.String()
	if !strings.Contains(balance, ".") {
		return balance
	}
	dataList := strings.Split(balance, ".")
	zhengshu := dataList[0]
	xiaoshu := dataList[1]
	if len(dataList[1]) > 18 {
		xiaoshu = xiaoshu[:18]
	}
	xiaoshu2 := ""
	
	for i := len(xiaoshu) - 1; i >= 0; i-- {
		if xiaoshu[i] != '0' {
			xiaoshu2 = xiaoshu[:i+1]
			break
		}
	}
	if xiaoshu2 == "" {
		return zhengshu
	} else {
		return zhengshu + "." + xiaoshu2
	}
}

func DealAmount(amount, symbol string) (string, error) {
	amountSlice := strings.Split(amount, ",")
	symbolSlice := strings.Split(symbol, ",")
	if len(amountSlice) != len(symbolSlice) {
		return "", errors.New("amount formatting error")
	}
	amountArray := []string{}
	for i, am := range amountSlice {
		amountArray = append(amountArray, am+symbolSlice[i])
	}
	return strings.Join(amountArray, ","), nil
}
