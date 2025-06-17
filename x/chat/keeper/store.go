package keeper

import (
	"freemasonry.cc/blockchain/util"
	"freemasonry.cc/blockchain/x/chat/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) IterateRegisterInfoWithAddr(ctx sdk.Context, handler func(addr types.RegisterInfo) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, []byte(types.KeyPrefixRegisterInfo))
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var userInfo types.UserInfo
		err := util.Json.Unmarshal(iter.Value(), &userInfo)
		if err != nil {
			panic(err)
		}

		info := types.RegisterInfo{
			FromAddress:         userInfo.FromAddress,
			RegisterNodeAddress: userInfo.RegisterNodeAddress,
			NodeAddress:         userInfo.NodeAddress,
			AddressBook:         userInfo.AddressBook,
			ChatBlacklist:       userInfo.ChatBlacklist,
			ChatWhitelist:       userInfo.ChatWhitelist,
			Mobile:              userInfo.Mobile,
			UpdateTime:          userInfo.UpdateTime,
			ChatBlackEncList:    userInfo.ChatBlackEncList,
			ChatWhiteEncList:    userInfo.ChatWhiteEncList,
		}

		if handler(info) {
			break
		}
	}
}

func (k Keeper) IterateAddrWithDid(ctx sdk.Context, handler func(didInfo types.PhonePrefix) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, []byte(types.KeyPrefixMobileOwner))
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var userInfo types.UserInfo
		err := util.Json.Unmarshal(iter.Value(), &userInfo)
		if err != nil {
			panic(err)
		}

		did := k.GetDidByKey(iter.Key())

		info := types.PhonePrefix{
			Mobile:  did,
			Address: string(iter.Value()),
		}

		if handler(info) {
			break
		}
	}
}

func (k Keeper) GetDidByKey(key []byte) string {
	keyStr := string(key)
	did := keyStr[len(types.KeyPrefixMobileOwner):]
	return did
}

func (k Keeper) IterateAddressBook(ctx sdk.Context, handler func(addressBook types.AddressBook) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, []byte(types.KeyChatAddress))
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		addr := string(iter.Value())
		chatAddr := k.GetChatAddrByKey(iter.Key())

		addressBook := types.AddressBook{
			ChatAddress: chatAddr,
			FromAddress: addr,
		}

		if handler(addressBook) {
			break
		}
	}
}

func (k Keeper) GetChatAddrByKey(key []byte) string {
	keyStr := string(key)
	addr := keyStr[len(types.KeyChatAddress):]
	return addr
}
