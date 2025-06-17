package types

import (
	sdkmath "cosmossdk.io/math"
	"fmt"
	"freemasonry.cc/blockchain/core"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"math"
	"strconv"
	"time"
)

const DefaultVotingPeriodTime time.Duration = time.Hour * 24

var (
	KeyRate                  = []byte("Rate") 
	KeySalaryRewardRatio     = []byte("SalaryRewardRatio")
	KeyDvmRewardRatio        = []byte("DvmRewardRatio")
	KeyBurnGetPowerRatio     = []byte("BurnGetPowerRatio")
	KeyClusterLevels         = []byte("ClusterLevels")
	KeyMaxClusterMembers     = []byte("MaxClusterMembers")
	KeyDaoRewardPercent      = []byte("DaoRewardPercent")
	KeyDposRewardPercent     = []byte("DposRewardPercent")
	KeyBurnCurrentGateRatio  = []byte("BurnCurrentGateRatio")
	KeyBurnRegisterGateRatio = []byte("BurnRegisterGateRatio")
	KeyDayMintAmount         = []byte("DayMintAmount")
	KeyBurnLevels            = []byte("BurnLevels")
	KeyPowerGasRatio         = []byte("PowerGasRatio")
	KeyAdPrice               = []byte("AdPrice")
	KeyAdRate                = []byte("AdRate")
	KeyBurnRewardFeeRate     = []byte("BurnRewardFeeRate")
	KeyReceiveDaoRatio       = []byte("ReceiveDaoRatio")
	KeyConnectivityDaoRatio  = []byte("ConnectivityDaoRatio")
	KeyBurnDaoPool           = []byte("BurnDaoPool")
	KeyDaoRewardRatio        = []byte("DaoRewardRatio")
	KeyMaxOnlineRatio        = []byte("MaxOnlineRatio")
	KeyVotingPeriod          = []byte("VotingPeriod")
	KeyMintBlockInterval     = []byte("MintBlockInterval")
	KeyIBCTransferFee        = []byte("IBCTransferFee")
	KeyTranslateMin          = []byte("TranslateMin")
	KeyDaoIncreaseRatio      = []byte("DaoIncreaseRatio")
	KeyDaoIncreaseHeight     = []byte("DaoIncreaseHeight")
	KeyIdoMinMember          = []byte("IdoMinMember")
)

func NewParams(
	rate,
	daoRewardPercent,
	dposRewardPercent,
	burnCurrentGateRatio,
	burnRegisterGateRati,
	burnGetPowerRatio sdk.Dec,
	salaryRewardRatio,
	dvmRewardRatio,
	daoRewardRatio RatioLimit,
	maxClusterMembers int64,
	clusterLevels []ClusterLevel,
	dayMintAmount sdk.Dec,
	burnLevels []BurnLevel,
	powerGasRatio sdk.Dec,
	adPrice sdk.Dec,
	adRate sdk.Dec,
	burnRewardFeeRate sdk.Dec,
	receiveDaoRatio sdk.Dec,
	connectivityDaoRatio sdk.Dec,
	burnDaoPool sdk.Dec,
	maxOnlineRatio sdk.Dec,
	votingPeriod time.Duration,
	mintBlockInterval int64,
	crossFee CrossFee,
	translateMin, daoIncreaseRatio sdk.Dec,
	daoIncreaseHeight int64,
	idoMinMember int64,
) Params {
	return Params{
		Rate:                  rate,
		SalaryRewardRatio:     salaryRewardRatio,
		DvmRewardRatio:        dvmRewardRatio,
		BurnGetPowerRatio:     burnGetPowerRatio,
		ClusterLevels:         clusterLevels,
		MaxClusterMembers:     maxClusterMembers,
		DaoRewardPercent:      daoRewardPercent,
		DposRewardPercent:     dposRewardPercent,
		BurnCurrentGateRatio:  burnCurrentGateRatio,
		BurnRegisterGateRatio: burnRegisterGateRati,
		DayMintAmount:         dayMintAmount,
		BurnLevels:            burnLevels,
		PowerGasRatio:         powerGasRatio,
		AdPrice:               adPrice,
		AdRate:                adRate,
		BurnRewardFeeRate:     burnRewardFeeRate,
		ReceiveDaoRatio:       receiveDaoRatio,
		ConnectivityDaoRatio:  connectivityDaoRatio,
		BurnDaoPool:           burnDaoPool,
		DaoRewardRatio:        daoRewardRatio,
		MaxOnlineRatio:        maxOnlineRatio,
		VotingPeriod:          votingPeriod,
		MintBlockInterval:     mintBlockInterval,
		CrossFee:              crossFee,
		TranslateMin:          translateMin,
		DaoIncreaseRatio:      daoIncreaseRatio,
		DaoIncreaseHeight:     daoIncreaseHeight,
		IdoMinMember:          idoMinMember,
	}
}

func DefaultParams() Params {

	return Params{
		Rate:                  sdk.NewDecWithPrec(5, 2),   
		SalaryRewardRatio:     DefaultSalaryRewardRatio(), 
		DvmRewardRatio:        DefaultDvmRewardRatio(),    
		BurnGetPowerRatio:     sdk.NewDec(100),            
		ClusterLevels:         DefaultClusterLevelsInfo(),
		MaxClusterMembers:     700,                        
		DaoRewardPercent:      sdk.NewDecWithPrec(2, 2),   
		DposRewardPercent:     sdk.NewDecWithPrec(1, 1),   
		BurnCurrentGateRatio:  sdk.NewDecWithPrec(575, 4), 
		BurnRegisterGateRatio: sdk.NewDecWithPrec(375, 4), 
		DayMintAmount:         sdk.MustNewDecFromStr("360000000000000000000000"),
		BurnLevels:            DefaultBurnLevels(),
		PowerGasRatio:         sdk.NewDec(100),
		AdPrice:               sdk.NewDecWithPrec(5, 2),       
		AdRate:                sdk.NewDecWithPrec(2, 1),       
		BurnRewardFeeRate:     sdk.MustNewDecFromStr("0.3"),   
		ReceiveDaoRatio:       sdk.MustNewDecFromStr("0.075"), 
		ConnectivityDaoRatio:  sdk.MustNewDecFromStr("0.7"),   
		BurnDaoPool:           sdk.MustNewDecFromStr("0.1"),   
		DaoRewardRatio:        DefaultDaoRewardRatio(),
		MaxOnlineRatio:        sdk.NewDecWithPrec(60, 2), 
		VotingPeriod:          DefaultVotingPeriodTime,   
		MintBlockInterval:     3600,
		CrossFee:              DefaulIBCTransferFee(),
		TranslateMin:          sdk.MustNewDecFromStr("10000000000000000000"),
		DaoIncreaseRatio:      sdk.MustNewDecFromStr("0.03"), 
		DaoIncreaseHeight:     600,                           
		IdoMinMember:          10,                            
	}
}

func (p Params) Validate() error {
	if err := validateDec(p.Rate); err != nil {
		return err
	}
	if err := validateRatioLimit(p.SalaryRewardRatio); err != nil {
		return err
	}
	if err := validateRatioLimit(p.DvmRewardRatio); err != nil {
		return err
	}
	if err := validateDec(p.BurnGetPowerRatio); err != nil {
		return err
	}
	if err := validateClusterLevels(p.ClusterLevels); err != nil {
		return err
	}
	if err := validateInt64(p.MaxClusterMembers); err != nil {
		return err
	}
	if err := validateDec(p.DaoRewardPercent); err != nil {
		return err
	}
	if err := validateDec(p.DposRewardPercent); err != nil {
		return err
	}
	if err := validateDec(p.BurnCurrentGateRatio); err != nil {
		return err
	}
	if err := validateDec(p.BurnRegisterGateRatio); err != nil {
		return err
	}
	if err := validateDec(p.DayMintAmount); err != nil {
		return err
	}
	if err := validateBurnLevels(p.BurnLevels); err != nil {
		return err
	}
	if err := validateDec(p.PowerGasRatio); err != nil {
		return err
	}
	if err := validateDec(p.AdPrice); err != nil {
		return err
	}
	if err := validateDec(p.AdRate); err != nil {
		return err
	}

	if err := validateDec(p.BurnRewardFeeRate); err != nil {
		return err
	}
	if err := validateDec(p.ReceiveDaoRatio); err != nil {
		return err
	}
	if err := validateDec(p.ConnectivityDaoRatio); err != nil {
		return err
	}
	if err := validateDec(p.BurnDaoPool); err != nil {
		return err
	}
	if err := validateRatioLimit(p.DaoRewardRatio); err != nil {
		return err
	}
	if err := validateDec(p.MaxOnlineRatio); err != nil {
		return err
	}
	if err := validateDec(p.MaxOnlineRatio); err != nil {
		return err
	}
	if err := validateTime(p.VotingPeriod); err != nil {
		return err
	}
	if err := validateCrossFee(p.CrossFee); err != nil {
		return err
	}
	if err := validateDec(p.TranslateMin); err != nil {
		return err
	}
	if err := validateDec(p.DaoIncreaseRatio); err != nil {
		return err
	}
	if err := validateInt64(p.DaoIncreaseHeight); err != nil {
		return err
	}
	if err := validateInt64(p.IdoMinMember); err != nil {
		return err
	}
	return nil
}

func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyRate, &p.Rate, validateDec),
		paramtypes.NewParamSetPair(KeySalaryRewardRatio, &p.SalaryRewardRatio, validateRatioLimit),
		paramtypes.NewParamSetPair(KeyDvmRewardRatio, &p.DvmRewardRatio, validateRatioLimit),
		paramtypes.NewParamSetPair(KeyBurnGetPowerRatio, &p.BurnGetPowerRatio, validateDec),
		paramtypes.NewParamSetPair(KeyClusterLevels, &p.ClusterLevels, validateClusterLevels),
		paramtypes.NewParamSetPair(KeyMaxClusterMembers, &p.MaxClusterMembers, validateInt64),
		paramtypes.NewParamSetPair(KeyDaoRewardPercent, &p.DaoRewardPercent, validateDec),
		paramtypes.NewParamSetPair(KeyDposRewardPercent, &p.DposRewardPercent, validateDec),
		paramtypes.NewParamSetPair(KeyBurnCurrentGateRatio, &p.BurnCurrentGateRatio, validateDec),
		paramtypes.NewParamSetPair(KeyBurnRegisterGateRatio, &p.BurnRegisterGateRatio, validateDec),
		paramtypes.NewParamSetPair(KeyDayMintAmount, &p.DayMintAmount, validateDec),
		paramtypes.NewParamSetPair(KeyBurnLevels, &p.BurnLevels, validateBurnLevels),
		paramtypes.NewParamSetPair(KeyPowerGasRatio, &p.PowerGasRatio, validateDec),
		paramtypes.NewParamSetPair(KeyAdPrice, &p.AdPrice, validateDec),
		paramtypes.NewParamSetPair(KeyAdRate, &p.AdRate, validateDec),
		paramtypes.NewParamSetPair(KeyBurnRewardFeeRate, &p.BurnRewardFeeRate, validateDec),
		paramtypes.NewParamSetPair(KeyReceiveDaoRatio, &p.ReceiveDaoRatio, validateDec),
		paramtypes.NewParamSetPair(KeyConnectivityDaoRatio, &p.ConnectivityDaoRatio, validateDec),
		paramtypes.NewParamSetPair(KeyBurnDaoPool, &p.BurnDaoPool, validateDec),
		paramtypes.NewParamSetPair(KeyDaoRewardRatio, &p.DaoRewardRatio, validateRatioLimit),
		paramtypes.NewParamSetPair(KeyMaxOnlineRatio, &p.MaxOnlineRatio, validateDec),
		paramtypes.NewParamSetPair(KeyVotingPeriod, &p.VotingPeriod, validateTime),
		paramtypes.NewParamSetPair(KeyMintBlockInterval, &p.MintBlockInterval, validateMintBlockInterval),
		paramtypes.NewParamSetPair(KeyIBCTransferFee, &p.CrossFee, validateCrossFee),
		paramtypes.NewParamSetPair(KeyTranslateMin, &p.TranslateMin, validateDec),
		paramtypes.NewParamSetPair(KeyDaoIncreaseRatio, &p.DaoIncreaseRatio, validateDec),
		paramtypes.NewParamSetPair(KeyDaoIncreaseHeight, &p.DaoIncreaseHeight, validateInt64),
		paramtypes.NewParamSetPair(KeyIdoMinMember, &p.IdoMinMember, validateInt64),
	}
}

func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable(
		paramtypes.NewParamSetPair(KeyRate, DefaultParams().Rate, validateDec),
		paramtypes.NewParamSetPair(KeySalaryRewardRatio, DefaultParams().SalaryRewardRatio, validateRatioLimit),
		paramtypes.NewParamSetPair(KeyDvmRewardRatio, DefaultParams().DvmRewardRatio, validateRatioLimit),
		paramtypes.NewParamSetPair(KeyBurnGetPowerRatio, DefaultParams().BurnGetPowerRatio, validateDec),
		paramtypes.NewParamSetPair(KeyClusterLevels, DefaultParams().ClusterLevels, validateClusterLevels),
		paramtypes.NewParamSetPair(KeyMaxClusterMembers, DefaultParams().MaxClusterMembers, validateInt64),
		paramtypes.NewParamSetPair(KeyDaoRewardPercent, DefaultParams().DaoRewardPercent, validateDec),
		paramtypes.NewParamSetPair(KeyDposRewardPercent, DefaultParams().DposRewardPercent, validateDec),
		paramtypes.NewParamSetPair(KeyBurnCurrentGateRatio, DefaultParams().BurnCurrentGateRatio, validateDec),
		paramtypes.NewParamSetPair(KeyBurnRegisterGateRatio, DefaultParams().BurnRegisterGateRatio, validateDec),
		paramtypes.NewParamSetPair(KeyDayMintAmount, DefaultParams().DayMintAmount, validateDec),
		paramtypes.NewParamSetPair(KeyBurnLevels, DefaultParams().BurnLevels, validateBurnLevels),
		paramtypes.NewParamSetPair(KeyPowerGasRatio, DefaultParams().PowerGasRatio, validateDec),
		paramtypes.NewParamSetPair(KeyAdPrice, DefaultParams().AdPrice, validateDec),
		paramtypes.NewParamSetPair(KeyAdRate, DefaultParams().AdRate, validateDec),
		paramtypes.NewParamSetPair(KeyBurnRewardFeeRate, DefaultParams().BurnRewardFeeRate, validateDec),
		paramtypes.NewParamSetPair(KeyReceiveDaoRatio, DefaultParams().ReceiveDaoRatio, validateDec),
		paramtypes.NewParamSetPair(KeyConnectivityDaoRatio, DefaultParams().ConnectivityDaoRatio, validateDec),
		paramtypes.NewParamSetPair(KeyBurnDaoPool, DefaultParams().BurnDaoPool, validateDec),
		paramtypes.NewParamSetPair(KeyDaoRewardRatio, DefaultParams().DaoRewardRatio, validateRatioLimit),
		paramtypes.NewParamSetPair(KeyMaxOnlineRatio, DefaultParams().MaxOnlineRatio, validateDec),
		paramtypes.NewParamSetPair(KeyVotingPeriod, DefaultParams().VotingPeriod, validateTime),
		paramtypes.NewParamSetPair(KeyMintBlockInterval, DefaultParams().MintBlockInterval, validateMintBlockInterval),
		paramtypes.NewParamSetPair(KeyIBCTransferFee, DefaultParams().CrossFee, validateCrossFee),
		paramtypes.NewParamSetPair(KeyTranslateMin, DefaultParams().TranslateMin, validateDec),
		paramtypes.NewParamSetPair(KeyDaoIncreaseRatio, DefaultParams().DaoIncreaseRatio, validateDec),
		paramtypes.NewParamSetPair(KeyDaoIncreaseHeight, DefaultParams().DaoIncreaseHeight, validateInt64),
		paramtypes.NewParamSetPair(KeyIdoMinMember, DefaultParams().IdoMinMember, validateInt64),
	)
}

func validateCrossFee(i interface{}) error {
	crossFee, ok := i.(CrossFee)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if crossFee.FeeCollectionAccount == "" {
		return fmt.Errorf("FeeCollectionAccount is empty")
	}
	if crossFee.FeeAmount.IsNegative() {
		return fmt.Errorf("negative decimal coin amount: %v\n", crossFee.FeeAmount)
	}
	return nil
}

func validateTime(i interface{}) error {
	v, ok := i.(time.Duration)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v <= 0 {
		return fmt.Errorf("time must be positive: %d", v)
	}

	return nil
}

func validateDec(i interface{}) error {
	_, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return nil
}

func validateRatioLimit(i interface{}) error {
	ratioLimit, ok := i.(RatioLimit)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if ratioLimit.MaxRatio.IsNegative() {
		return fmt.Errorf("negative decimal coin amount: %v\n", ratioLimit.MaxRatio)
	}
	if ratioLimit.MinRatio.IsNegative() {
		return fmt.Errorf("negative decimal coin amount: %v\n", ratioLimit.MinRatio)
	}
	return nil
}

func validateBurnLevels(i interface{}) error {
	burnLevels, ok := i.([]BurnLevel)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	burnCount := len(burnLevels)

	if burnCount == 0 {
		return core.ErrParamsBurnLevels
	}

	
	if burnLevels[0].Level != 1 {
		return core.ErrParamsBurnLevels
	}

	
	for n := 1; n < len(burnLevels)+1; n++ {

		if burnCount > 1 && n != burnCount {

			
			if burnLevels[n-1].Level+1 != burnLevels[n].Level {
				return core.ErrParamsLevel
			}

			
			if burnLevels[n-1].BurnAmount.GTE(burnLevels[n].BurnAmount) {
				return core.ErrParamsPeldgeLevel
			}
		}
	}
	return nil
}

func validateClusterLevels(i interface{}) error {

	levels, ok := i.([]ClusterLevel)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if len(levels) == 0 {
		return core.ErrClusterLevelsParams
	}

	levelsCount := len(levels)

	
	if levels[0].Level != 0 {
		return core.ErrClusterLevelsParams
	}

	
	for n := 1; n < len(levels)+1; n++ {

		if levelsCount > 1 && n != levelsCount {

			
			if levels[n-1].Level+1 != levels[n].Level {
				return core.ErrClusterLevelsParams
			}

			
			if levels[n-1].BurnAmount.GT(levels[n].BurnAmount) {
				return core.ErrClusterLevelsParams
			}

			
			
			if levels[n-1].MemberAmount > levels[n].MemberAmount {
				return core.ErrClusterLevelsParams
			}
		}
	}

	return nil
}

func validateInt64(i interface{}) error {
	_, ok := i.(int64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return nil
}

func validateSdkInt(i interface{}) error {
	_, ok := i.(sdkmath.Int)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return nil
}

func validateMintBlockInterval(i interface{}) error {
	iValue, ok := i.(int64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if iValue < 100 || iValue > 14400 {
		return fmt.Errorf("mint block interval must be between 100 and 14400: %d", iValue)
	}

	return nil
}

func DefaultSalaryRewardRatio() RatioLimit {
	rateLimit := RatioLimit{
		MinRatio: sdk.ZeroDec(),
		MaxRatio: sdk.NewDecWithPrec(49, 2),
	}
	return rateLimit
}

func DefaultDvmRewardRatio() RatioLimit {
	rateLimit := RatioLimit{
		MinRatio: sdk.ZeroDec(),
		MaxRatio: sdk.NewDecWithPrec(49, 2),
	}
	return rateLimit
}

func DefaultDaoRewardRatio() RatioLimit {
	rateLimit := RatioLimit{
		MinRatio: sdk.ZeroDec(),
		MaxRatio: sdk.OneDec(),
	}
	return rateLimit
}

func DefaulIBCTransferFee() CrossFee {
	crossFee := CrossFee{
		FeeCollectionAccount: "0x4918aD481e4512DedE7565A8a4ddbc06881D9Bf7",
		FeeAmount:            sdkmath.NewInt(150000),
	}
	return crossFee
}

func DefaultClusterLevelsInfo() []ClusterLevel {
	clusterLevels := make([]ClusterLevel, 0)
	defaultBurnAmount := DefaultClusterMinPledge()
	defaultMembersAmount := DefaultClusterMinMembers()
	for i := int64(0); i < int64(34); i++ {
		var clusterLevel ClusterLevel
		if i == 0 {
			clusterLevel = ClusterLevel{
				Level:        0,
				BurnAmount:   sdk.ZeroInt(),
				DaoLimit:     sdk.ZeroInt(),
				MemberAmount: 1,
			}
		} else if i == 1 {
			clusterLevel = ClusterLevel{
				Level:        1,
				BurnAmount:   sdk.NewInt(30000000000000000),
				DaoLimit:     defaultBurnAmount[i].Mul(sdk.NewInt(core.RealToLedgerRateInt64)),
				MemberAmount: defaultMembersAmount[i],
			}
		} else {
			clusterLevel = ClusterLevel{
				Level:        i,
				BurnAmount:   defaultBurnAmount[i].Mul(sdk.NewInt(core.RealToLedgerRateInt64)),
				DaoLimit:     defaultBurnAmount[i].Mul(sdk.NewInt(core.RealToLedgerRateInt64)),
				MemberAmount: defaultMembersAmount[i],
			}
		}

		clusterLevels = append(clusterLevels, clusterLevel)
	}

	return clusterLevels
}

func DefaultBurnLevels() []BurnLevel {
	var err error
	burnLevels := make([]BurnLevel, 0)
	defaultAddPercemt := DefaultAddPercent()
	defaultRoomAmount := DefaultRoomAmount()
	for i := int64(1); i < int64(34); i++ {
		amountInt := sdk.MustNewDecFromStr("0")
		if i > 1 {
			ifloat := float64(i)
			amountBaseFloat64 := math.Pow(2, ifloat-2) * 100
			amountFloat64 := amountBaseFloat64 * float64(core.RealToLedgerRateInt64)
			amountString := strconv.FormatFloat(amountFloat64, 'f', 0, 64)
			amountInt, err = sdk.NewDecFromStr(amountString)

			if err != nil {
				panic("Err Default Dao Params BurnLevel")
			}
		}

		burnLevel := BurnLevel{
			Level:      i,
			BurnAmount: amountInt,
			AddPercent: defaultAddPercemt[i],
			RoomAmount: defaultRoomAmount[i],
		}

		burnLevels = append(burnLevels, burnLevel)
	}

	return burnLevels
}

func DefaultRoomAmount() map[int64]sdk.Int {
	
	return map[int64]sdk.Int{
		1:  sdk.NewInt(0),
		2:  sdk.NewInt(0),
		3:  sdk.NewInt(0),
		4:  sdk.NewInt(0),
		5:  sdk.NewInt(1),
		6:  sdk.NewInt(2),
		7:  sdk.NewInt(3),
		8:  sdk.NewInt(4),
		9:  sdk.NewInt(5),
		10: sdk.NewInt(6),
		11: sdk.NewInt(7),
		12: sdk.NewInt(8),
		13: sdk.NewInt(9),
		14: sdk.NewInt(10),
		15: sdk.NewInt(11),
		16: sdk.NewInt(12),
		17: sdk.NewInt(13),
		18: sdk.NewInt(14),
		19: sdk.NewInt(15),
		20: sdk.NewInt(16),
		21: sdk.NewInt(17),
		22: sdk.NewInt(18),
		23: sdk.NewInt(19),
		24: sdk.NewInt(20),
		25: sdk.NewInt(21),
		26: sdk.NewInt(22),
		27: sdk.NewInt(23),
		28: sdk.NewInt(24),
		29: sdk.NewInt(25),
		30: sdk.NewInt(26),
		31: sdk.NewInt(27),
		32: sdk.NewInt(28),
		33: sdk.NewInt(29),
	}
}

func DefaultAddPercent() map[int64]sdk.Int {
	
	return map[int64]sdk.Int{
		1:  sdk.NewInt(1),
		2:  sdk.NewInt(2),
		3:  sdk.NewInt(3),
		4:  sdk.NewInt(4),
		5:  sdk.NewInt(5),
		6:  sdk.NewInt(6),
		7:  sdk.NewInt(7),
		8:  sdk.NewInt(8),
		9:  sdk.NewInt(9),
		10: sdk.NewInt(12),
		11: sdk.NewInt(15),
		12: sdk.NewInt(18),
		13: sdk.NewInt(21),
		14: sdk.NewInt(24),
		15: sdk.NewInt(27),
		16: sdk.NewInt(30),
		17: sdk.NewInt(33),
		18: sdk.NewInt(36),
		19: sdk.NewInt(39),
		20: sdk.NewInt(43),
		21: sdk.NewInt(47),
		22: sdk.NewInt(51),
		23: sdk.NewInt(55),
		24: sdk.NewInt(59),
		25: sdk.NewInt(63),
		26: sdk.NewInt(67),
		27: sdk.NewInt(71),
		28: sdk.NewInt(75),
		29: sdk.NewInt(79),
		30: sdk.NewInt(84),
		31: sdk.NewInt(89),
		32: sdk.NewInt(94),
		33: sdk.NewInt(100),
	}
}

func DefaultClusterMinPledge() map[int64]sdk.Int {
	return map[int64]sdk.Int{
		1:  sdk.NewInt(5000),
		2:  sdk.NewInt(10000),
		3:  sdk.NewInt(15000),
		4:  sdk.NewInt(20000),
		5:  sdk.NewInt(25000),
		6:  sdk.NewInt(30000),
		7:  sdk.NewInt(40000),
		8:  sdk.NewInt(50000),
		9:  sdk.NewInt(65000),
		10: sdk.NewInt(80000),
		11: sdk.NewInt(100000),
		12: sdk.NewInt(120000),
		13: sdk.NewInt(145000),
		14: sdk.NewInt(175000),
		15: sdk.NewInt(210000),
		16: sdk.NewInt(250000),
		17: sdk.NewInt(310000),
		18: sdk.NewInt(375000),
		19: sdk.NewInt(455000),
		20: sdk.NewInt(550000),
		21: sdk.NewInt(660000),
		22: sdk.NewInt(795000),
		23: sdk.NewInt(955000),
		24: sdk.NewInt(1150000),
		25: sdk.NewInt(1385000),
		26: sdk.NewInt(1665000),
		27: sdk.NewInt(2000000),
		28: sdk.NewInt(2405000),
		29: sdk.NewInt(2890000),
		30: sdk.NewInt(3470000),
		31: sdk.NewInt(4165000),
		32: sdk.NewInt(5000000),
		33: sdk.NewInt(6000000),
	}
}

func DefaultClusterMinMembers() map[int64]int64 {
	
	return map[int64]int64{
		1:  1,
		2:  2,
		3:  3,
		4:  4,
		5:  5,
		6:  6,
		7:  8,
		8:  10,
		9:  12,
		10: 15,
		11: 18,
		12: 21,
		13: 25,
		14: 30,
		15: 35,
		16: 41,
		17: 48,
		18: 57,
		19: 67,
		20: 79,
		21: 93,
		22: 109,
		23: 128,
		24: 150,
		25: 175,
		26: 204,
		27: 238,
		28: 278,
		29: 325,
		30: 379,
		31: 442,
		32: 515,
		33: 600,
	}
}
