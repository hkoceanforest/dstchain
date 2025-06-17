package keeper

import (
	"freemasonry.cc/blockchain/x/chat/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) InitGenesis(ctx sdk.Context, data types.GenesisState) {
	k.SetParams(ctx, data.Params)

	for _, info := range data.RegisterInfos {
		userinfo := types.UserInfo{
			FromAddress:         info.FromAddress,
			RegisterNodeAddress: info.RegisterNodeAddress,
			NodeAddress:         info.NodeAddress,
			AddressBook:         info.AddressBook,
			ChatBlacklist:       info.ChatBlacklist,
			ChatWhitelist:       info.ChatWhitelist,
			Mobile:              info.Mobile,
			UpdateTime:          info.UpdateTime,
			ChatBlackEncList:    info.ChatBlackEncList,
			ChatWhiteEncList:    info.ChatWhiteEncList,
		}
		err := k.SetRegisterInfo(ctx, userinfo)
		if err != nil {
			panic(err)
		}
	}

	for _, didInfo := range data.PhonePrefixes {
		err := k.SetMobileOwner(ctx, didInfo.Mobile, didInfo.Address)
		if err != nil {
			panic(err)
		}
	}

	for _, chatAddrInfo := range data.AddressBooks {
		err := k.SetChatAddr(ctx, chatAddrInfo.ChatAddress, chatAddrInfo.FromAddress)
		if err != nil {
			panic(err)
		}
	}

}

func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {

	registerInfos := make([]types.RegisterInfo, 0)
	k.IterateRegisterInfoWithAddr(ctx,
		func(info types.RegisterInfo) (stop bool) {
			registerInfos = append(registerInfos, info)
			return false
		},
	)

	didInfos := make([]types.PhonePrefix, 0)
	k.IterateAddrWithDid(ctx,
		func(didInfo types.PhonePrefix) (stop bool) {
			didInfos = append(didInfos, didInfo)
			return false
		},
	)

	addressBooks := make([]types.AddressBook, 0)
	k.IterateAddressBook(ctx,
		func(book types.AddressBook) (stop bool) {
			addressBooks = append(addressBooks, book)
			return false
		},
	)

	return &types.GenesisState{
		Params:        k.GetParams(ctx),
		RegisterInfos: registerInfos,
		PhonePrefixes: didInfos,
		AddressBooks:  addressBooks,
	}
}
