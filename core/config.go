package core

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	AppName = "daodst"

	
	BaseDenom = "dst"

	
	GovDenom = sdk.DefaultBondDenom

	
	BurnRewardFeeDenom = "dao"

	
	UsdtDenom = "usdt"

	CosmosApiPort = "1317"

	RpcPort = "26657"

	P2pPort = "26656"

	ChainDefaultFeeStr string = "100000000000000000" + BaseDenom 

	
	MinRealAmountFloat64 float64 = 0.0000000001

	
	RealToLedgerRate float64 = float64(RealToLedgerRateInt64)

	
	RealToLedgerRateInt64 int64 = 1000000000000000000

	
	LedgerToRealRate = "0.000000000000000001"

	
	CommitTime = "6s"

	
	GatewayBonusAddress = "gatewayBonus"

	
	MaxClusterPowerMembersAmount = 10000

	
	DayBlockNum = 14400

	
	PeriodicRewardsBlockNum = 3600

	
	CutProductionSeconds = 126230400

	
	StartMintSeconds = 259200

	
	OracleUpdateInterval = 14400

	
	MaxRedPacketAmount = "999999999999999999999999999999999999"

	
	RedPacketTimeOut = 14400

	Erc20RewardAccount = "erc20Reward"

	CutPowerRewardRatio = "0.01"

	CutPowerRewardTimes = int64(333)

	TxTimeoutHeight = 5

	
	IdoEndAddLiquidityAmountDST = "2160000000000000000000000"

	
	IdoEndAddLiquidityAmountUSDT = "108000000000"
)

var (
	
	GenesisIdoSupply = sdk.MustNewDecFromStr("36000000000000000000000000")

	
	DefaultChainSeed = []string{
		"bd961b42cd9286af3d74015739255a68d2d72c0e@bebewallet.com:26656",
	}

	MinRealAmountDec = sdk.NewDecWithPrec(1, 10)

	RealToLedgerRateDec = sdk.MustNewDecFromStr("1000000000000000000")

	UsdtRealToLedgerRateDec = sdk.MustNewDecFromStr("1000000")

	
	NumIndexAmount = map[int]bool{
		5: true,
		6: true,
		7: true,
	}
	
	ClusterDeviceRate = sdk.NewDecWithPrec(1, 4)

	
	DstSupply = sdk.MustNewDecFromStr("1234560000000000000000000000")

	BaseRate = sdk.MustNewDecFromStr("3")

	BurnRepair = sdk.MustNewDecFromStr("32400000000000000000000000").Sub(sdk.NewDec(108000000)).Add(sdk.MustNewDecFromStr("19440000000000000000000000")).TruncateInt()
)


var (
	ValReplace = map[string]string{
		"dstvaloper1ll30h0xykgduvxxfnpy4h6yzl0770pgn7hn3lz": "dstvaloper1kxcnfvtp6ep042vudctgda29cywkpzp0v9xkmc",
	}
)
