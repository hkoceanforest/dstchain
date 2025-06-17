package core

import (
	sdkerrors "cosmossdk.io/errors"
)

var (
	step = uint32(0)
)

func incr() uint32 {
	step += 1
	return step
}



var (
	BroadcastTimeOut = sdkerrors.Register(BaseDenom, incr(), "Broadcase time out")
)

var (
	
	ErrUserHasExisted       = sdkerrors.Register(BaseDenom, incr(), "user has existed")
	ErrRegister             = sdkerrors.Register(BaseDenom, incr(), "user register faild")
	ErrUserNotFound         = sdkerrors.Register(BaseDenom, incr(), "user not found")
	ErrAddressFormat        = sdkerrors.Register(BaseDenom, incr(), "address format error")
	ErrUserUpdate           = sdkerrors.Register(BaseDenom, incr(), "user info update error")
	ErrNumberOfGateWay      = sdkerrors.Register(BaseDenom, incr(), "number of gateway error")
	ErrGetMobile            = sdkerrors.Register(BaseDenom, incr(), "error get mobile")
	ErrGateway              = sdkerrors.Register(BaseDenom, incr(), "error gateway address")
	ErrMobileSetError       = sdkerrors.Register(BaseDenom, incr(), "error mobile set")
	ErrUserNotHaveMobile    = sdkerrors.Register(BaseDenom, incr(), "user does not have a mobile phone number")
	ErrUserMobileCount      = sdkerrors.Register(BaseDenom, incr(), "user mobile count error")
	ErrChatInfoSet          = sdkerrors.Register(BaseDenom, incr(), "chat info set error")
	ErrMobileexhausted      = sdkerrors.Register(BaseDenom, incr(), "The number index has been exhausted")
	ErrValidatorStatusError = sdkerrors.Register(BaseDenom, incr(), "validator status error")
	ErrGetGatewauInfo       = sdkerrors.Register(BaseDenom, incr(), "get gateway info error")
	ErrBurn                 = sdkerrors.Register(BaseDenom, incr(), "burn error")
	ErrMaxPhoneNumber       = sdkerrors.Register(BaseDenom, incr(), "max phone number must be between 2 and 100")
	ErrDestroyCoinDenom     = sdkerrors.Register(BaseDenom, incr(), "destroy coin denom error")
	ErrDefaultChatFee       = sdkerrors.Register(BaseDenom, incr(), "default chat fee error")
	ErrSetMobileOwner       = sdkerrors.Register(BaseDenom, incr(), "set mobile owner error")
	ErrMobileNotFount       = sdkerrors.Register(BaseDenom, incr(), "mobile not found")
	ErrMobileTransferTo     = sdkerrors.Register(BaseDenom, incr(), "can't transfer your phone number to yourself")
	ErrChatAddressExist     = sdkerrors.Register(BaseDenom, incr(), "chat address exist")
	ErrChatAddressNotExist  = sdkerrors.Register(BaseDenom, incr(), "chat address does not exist")

	
	ErrGatewayNumLen      = sdkerrors.Register(BaseDenom, incr(), "Number is empty")
	ErrGatewayNum         = sdkerrors.Register(BaseDenom, incr(), "Number segment overrun")
	ErrDelegationCoin     = sdkerrors.Register(BaseDenom, incr(), "Invalid amount")
	ErrGatewayNumber      = sdkerrors.Register(BaseDenom, incr(), "Number Already registered")
	ErrGatewayDelegation  = sdkerrors.Register(BaseDenom, incr(), "Insufficient mortgage amount")
	ErrGatewayNotExist    = sdkerrors.Register(BaseDenom, incr(), "gateway not exist")
	ErrGatewayExist       = sdkerrors.Register(BaseDenom, incr(), "gateway already registered")
	ErrGatewayFirstNum    = sdkerrors.Register(BaseDenom, incr(), "the first number segment cannot be redeemed")
	ErrGatewayNumNotFound = sdkerrors.Register(BaseDenom, incr(), "gateway number not found")
	ErrGatewayNumLength   = sdkerrors.Register(BaseDenom, incr(), "Illegal length of number segment")
	ErrValidatorNotFound  = sdkerrors.Register(BaseDenom, incr(), "Is not an validator node")
	ErrGatewayUrl         = sdkerrors.Register(BaseDenom, incr(), "gateway url not match")
	ErrGatewayPackage     = sdkerrors.Register(BaseDenom, incr(), "gateway package not match")
	ErrUnmarshal          = sdkerrors.Register(BaseDenom, incr(), "unmarshal error")

	
	ErrABIPack                  = sdkerrors.Register(BaseDenom, incr(), "contract ABI pack failed")
	ErrEmptyContractAddress     = sdkerrors.Register(BaseDenom, incr(), "invalid contract address")
	ErrTokenNotFound            = sdkerrors.Register(BaseDenom, incr(), "token not found")
	ErrStringNumber             = sdkerrors.Register(BaseDenom, incr(), "string number error")
	ErrTokenPairNotFound        = sdkerrors.Register(BaseDenom, incr(), "token pair not found")
	ErrCrossChainHashCheck      = sdkerrors.Register(BaseDenom, incr(), "cross chain hash check error")
	ErrCrossChainHashRepetitive = sdkerrors.Register(BaseDenom, incr(), "cross chain hash repetitive")

	
	ErrMachineAddress            = sdkerrors.Register(BaseDenom, incr(), "address mismatch")
	ErrMachineUpdateTime         = sdkerrors.Register(BaseDenom, incr(), "insufficient update interval")
	ErrSetCluster                = sdkerrors.Register(BaseDenom, incr(), "set cluster error")
	ErrGetCluster                = sdkerrors.Register(BaseDenom, incr(), "get cluster error")
	GetClusterNotFound           = sdkerrors.Register(BaseDenom, incr(), "cluster not found")
	GetClusterExisted            = sdkerrors.Register(BaseDenom, incr(), "cluster already existed")
	ErrSetPersonClusterInfo      = sdkerrors.Register(BaseDenom, incr(), "set person cluster info error")
	ErrGetPersonClusterInfo      = sdkerrors.Register(BaseDenom, incr(), "get person cluster info error")
	ErrClusterLevelsParams       = sdkerrors.Register(BaseDenom, incr(), "cluster levels params error")
	ErrEmptyBurnStartInfo        = sdkerrors.Register(BaseDenom, incr(), "no burn startInfo")
	ErrAddTotalBurnAmount        = sdkerrors.Register(BaseDenom, incr(), "add total burn amount error")
	ErrSubTotalBurnAmount        = sdkerrors.Register(BaseDenom, incr(), "sub total burn amount error")
	ErrGetTotalBurnAmount        = sdkerrors.Register(BaseDenom, incr(), "get total burn amount error")
	ErrDaoBurn                   = sdkerrors.Register(BaseDenom, incr(), "burn error")
	ErrMintReward                = sdkerrors.Register(BaseDenom, incr(), "mint reward error")
	ErrFreezeBurn                = sdkerrors.Register(BaseDenom, incr(), "frozen power cannot be given")
	ErrGetClusterId              = sdkerrors.Register(BaseDenom, incr(), "get cluster id error")
	ErrClusterIdNotFound         = sdkerrors.Register(BaseDenom, incr(), "cluster id not found")
	ErrParamsBurnLevels          = sdkerrors.Register(BaseDenom, incr(), "params burn level error")
	ErrDeviceRatio               = sdkerrors.Register(BaseDenom, incr(), "Incorrect device reward ratio")
	ErrSalaryRatio               = sdkerrors.Register(BaseDenom, incr(), "Incorrect salary reward ratio")
	ErrDaoRatio                  = sdkerrors.Register(BaseDenom, incr(), "Incorrect dao reward ratio")
	ErrAddClusterToGateway       = sdkerrors.Register(BaseDenom, incr(), "add cluster to gateway error")
	ErrGetClusterToGateway       = sdkerrors.Register(BaseDenom, incr(), "get cluster to gateway error")
	ErrFrozenInsufficient        = sdkerrors.Register(BaseDenom, incr(), "insufficient frozen power")
	ErrMemberNotExist            = sdkerrors.Register(BaseDenom, incr(), "cluster member does not exist")
	ErrMemberAlreadyExist        = sdkerrors.Register(BaseDenom, incr(), "cluster member already exist")
	ErrOwnerCannotExit           = sdkerrors.Register(BaseDenom, incr(), "group owner can not exit")
	ErrOwnerPermission           = sdkerrors.Register(BaseDenom, incr(), "group leader permission error")
	ErrMaxClusterMember          = sdkerrors.Register(BaseDenom, incr(), "number of people exceeds the upper limit")
	ErrNotInCluster              = sdkerrors.Register(BaseDenom, incr(), "not a cluster member")
	ErrClusterPermission         = sdkerrors.Register(BaseDenom, incr(), "cluster permission error")
	ErrCreateClusterBurn         = sdkerrors.Register(BaseDenom, incr(), "insufficient create cluster burn amount")
	ErrClusterOwnerPower         = sdkerrors.Register(BaseDenom, incr(), "cluster owner power not found")
	ErrGatewayNotFound           = sdkerrors.Register(BaseDenom, incr(), "gateway not found")
	ErrParamsLevel               = sdkerrors.Register(BaseDenom, incr(), "params level error")
	EmptyGetClustertoGate        = sdkerrors.Register(BaseDenom, incr(), "cluster for gateway not found")
	ErrSetDeviceCluster          = sdkerrors.Register(BaseDenom, incr(), "set device cluster error")
	ErrLevelInfo                 = sdkerrors.Register(BaseDenom, incr(), "level error")
	ErrGetApprovePowerInfo       = sdkerrors.Register(BaseDenom, incr(), "get approvr power info error")
	ErrClusterOwnerErr           = sdkerrors.Register(BaseDenom, incr(), "cluster owner not match")
	ErrNoBurn                    = sdkerrors.Register(BaseDenom, incr(), "no burn records")
	ErrUseFreeze                 = sdkerrors.Register(BaseDenom, incr(), "freeze power and burn cannot be used simultaneously")
	ErrApproveNotEnd             = sdkerrors.Register(BaseDenom, incr(), "approve is not end")
	ErrCluserMaxPowerMembers     = sdkerrors.Register(BaseDenom, incr(), "member has reached the maximum limit")
	ErrClusterConfigChange       = sdkerrors.Register(BaseDenom, incr(), "cluster configuration can only be changed once a day")
	ErrContractAddress           = sdkerrors.Register(BaseDenom, incr(), "contract address is not found")
	ErrVotePolicyAddress         = sdkerrors.Register(BaseDenom, incr(), "vote policy address and from address are inconsistent")
	ErrPriceFormat               = sdkerrors.Register(BaseDenom, incr(), "price format err")
	ErrParamsPeldgeLevel         = sdkerrors.Register(BaseDenom, incr(), "params pledge level error")
	ErrAuthorization             = sdkerrors.Register(BaseDenom, incr(), "authorization failure error")
	ErrRedPacketExist            = sdkerrors.Register(BaseDenom, incr(), "red packet exist")
	ErrRedPacketType             = sdkerrors.Register(BaseDenom, incr(), "red packet type error")
	ErrRedPacketMarshal          = sdkerrors.Register(BaseDenom, incr(), "red packet marshal error")
	ErrRedPacketUnmarshal        = sdkerrors.Register(BaseDenom, incr(), "red packet unmarshal error")
	ErrRedPacketAmountMax        = sdkerrors.Register(BaseDenom, incr(), "The amount of the red packet exceeds the limit")
	ErrRedPacketAmountMin        = sdkerrors.Register(BaseDenom, incr(), "The amount of the red packet is below the limit")
	ErrRedPacketCount            = sdkerrors.Register(BaseDenom, incr(), "RedPacket Count Err")
	ErrRedPacketNotExist         = sdkerrors.Register(BaseDenom, incr(), "RedPacket Not Exist")
	ErrRedPacketClusters         = sdkerrors.Register(BaseDenom, incr(), "Cannot receive red packets from other clusters")
	ErrRedPacketExpired          = sdkerrors.Register(BaseDenom, incr(), "The red envelope has expired")
	ErrRedPacketCollected        = sdkerrors.Register(BaseDenom, incr(), "The red envelope has been collected")
	ErrRedPacketRepeat           = sdkerrors.Register(BaseDenom, incr(), "Repeat red packet collection")
	ErrRedPacketInsufficient     = sdkerrors.Register(BaseDenom, incr(), "The amount of red packet is insufficient")
	ErrRedPacketEndBlock         = sdkerrors.Register(BaseDenom, incr(), "Red envelope not expired")
	ErrRedPacketSender           = sdkerrors.Register(BaseDenom, incr(), "This red envelope doesn't belong to you")
	ErrErrMemberAlreadyInCluster = sdkerrors.Register(BaseDenom, incr(), "The member is already in the cluster")
	ErrCutProductionNotStart     = sdkerrors.Register(BaseDenom, incr(), "The cut production is not started")
	ErrNoPowerRewardCycleInfo    = sdkerrors.Register(BaseDenom, incr(), "Power reward cycle info not found")
	ErrCurrentCycleAlreadyStart  = sdkerrors.Register(BaseDenom, incr(), "Current cycle already start")
	ErrNoPowerCutReward          = sdkerrors.Register(BaseDenom, incr(), "No power cut reward")
	ErrIBCTransferAmount         = sdkerrors.Register(BaseDenom, incr(), "IBC transfer amount is too less")
	ErrErrLoopCluster            = sdkerrors.Register(BaseDenom, incr(), "Loop cluster error")
	ErrIdoLimit                  = sdkerrors.Register(BaseDenom, incr(), "Subscription limit is exceeded")

	
	ErrRewardReceive                     = sdkerrors.Register(BaseDenom, incr(), "there is no reward to receive")
	ErrDelegationZero                    = sdkerrors.Register(BaseDenom, incr(), "delegate amount is zero")
	ErrDelegateAmountLtMinSelfDelegation = sdkerrors.Register(BaseDenom, incr(), "delegate amount less than min self delegation")
	ErrUnjailOperatorAddress             = sdkerrors.Register(BaseDenom, incr(), "delegate address error")
	ErrFreezePowerInsufficient           = sdkerrors.Register(BaseDenom, incr(), "freeze power insufficient")
	ErrQueryTokens                       = sdkerrors.Register(BaseDenom, incr(), "query tokens error")
	ErrNoPower                           = sdkerrors.Register(BaseDenom, incr(), "no power")
	ErrNoDaoCanReceive                   = sdkerrors.Register(BaseDenom, incr(), "dao insufficient")
	ErrInsufficientBalance               = sdkerrors.Register(BaseDenom, incr(), "Insufficient balance")
	ErrDvmRatio                          = sdkerrors.Register(BaseDenom, incr(), "Incorrect dvm reward ratio")
	ErrGatewaySign                       = sdkerrors.Register(BaseDenom, incr(), "gateway sign error")
	ErrGenesisIdoEnd                     = sdkerrors.Register(BaseDenom, incr(), "genesis ido end")
)

var (
	ParseAccountError = sdkerrors.Register(BaseDenom, incr(), "parse account error")
	ParseCoinError    = sdkerrors.Register(BaseDenom, incr(), "parse coin error")
	FieldEmptyError   = sdkerrors.Register(BaseDenom, incr(), "field empty error")
	DenomError        = sdkerrors.Register(BaseDenom, incr(), "denom error")
	SignError         = sdkerrors.Register(BaseDenom, incr(), "invalid sign")
	SignAccountError  = sdkerrors.Register(BaseDenom, incr(), "invalid sign account")
	UnconfirmedTxsErr = sdkerrors.Register(BaseDenom, incr(), "You have an unconfirmed transaction on the chain, please wait for processing to complete")
	ParamsInvalidErr  = sdkerrors.Register(BaseDenom, incr(), "params invalid")
	DstChainNodeErr   = sdkerrors.Register(BaseDenom, incr(), "dst chain node connect fail")
	KavaChainNodeErr  = sdkerrors.Register(BaseDenom, incr(), "kava chain node connect fail")

	QueryChainInforError      = "query chain infor errors"
	ValidatorDescriptionError = "validator description length error"
	ValidatorInfoError        = "Verifier information can only be changed once in 24 hours"
	FeeIsTooLess              = "fee is too less"
	ErrorGasOut               = "The gas consumed exceeds the upper limit set by the client"
	ErrUnauthorized           = "signature verification failed invalid chainid or account number"
	ErrWrongSequence          = "account serial number expired the reason may be node block behind or repeatedly sent messages"
)
