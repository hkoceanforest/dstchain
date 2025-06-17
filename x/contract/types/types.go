package types

import (
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

type NftInfo struct {
	TokenId       int64  `json:"tokenId"`
	CreateAddress string `json:"createAddress"`
	Level         int64  `json:"level"`
	CreateTime    int64  `json:"createTime"`
}

type StakeFactoryCreateParams struct {
	GatewayAddress   common.Address `json:"gatewayAddress"`
	Name             string         `json:"name"`
	Symbol           string         `json:"symbol"`
	PreMintAmount    *big.Int       `json:"preMintAmount"`
	TokenDecimals    uint8          `json:"tokenDecimals"`
	StakeRate        *big.Int       `json:"stakeRate"`
	UnStakeRate      *big.Int       `json:"unStakeRate"`
	MinInflationRate *big.Int       `json:"minInflationRate"`
	MaxInflationRate *big.Int       `json:"maxInflationRate"`
	TargetStakeRate  *big.Int       `json:"targetStakeRate"`
	ScaleFactor      *big.Int       `json:"scaleFactor"`
	ChatThreshold    *big.Int       `json:"chatThreshold"`
}

type StakeFactoryStakeInfo struct {
	TokenAddress     common.Address `json:"tokenAddress"`
	StakeAddress     common.Address `json:"stakeAddress"`
	TokenSymbol      string         `json:"tokenSymbol"`
	Decimals         *big.Int       `json:"decimals"`
	StakeRate        *big.Int       `json:"stakeRate"`
	UnStakeRate      *big.Int       `json:"unStakeRate"`
	GatewayAddress   common.Address `json:"gatewayAddress"`
	MinInflationRate *big.Int       `json:"minInflationRate"`
	MaxInflationRate *big.Int       `json:"maxInflationRate"`
	TargetStakeRate  *big.Int       `json:"targetStakeRate"`
	ScaleFactor      *big.Int       `json:"scaleFactor"`
	ChatThreshold    *big.Int       `json:"chatThreshold"`
}

type QueryGatewayTokenInfoParams struct {
	FromAddress string `json:"from_address"`
}

type GatewayTokenInfo struct {
	TokenLogo string
	StakeFactoryCreateParams
}

type GatewayTokenStakeInfo struct {
	TokenLogo string
	StakeFactoryStakeInfo
}

type TokenFactoryCreateParams struct {
	Owner         common.Address `json:"owner"`
	Name          string         `json:"name"`
	Symbol        string         `json:"symbol"`
	PreMintAmount *big.Int       `json:"preMintAmount"`
	Decimals      uint8          `json:"decimals"`
}

type TokenNewCreateInfo struct {
	TokenAddress  common.Address `json:"tokenAddress"`
	Name          string         `json:"name"`
	Symbol        string         `json:"symbol"`
	Decimals      uint8          `json:"decimals"`
	PreMintAmount *big.Int       `json:"preMintAmount"`
}
