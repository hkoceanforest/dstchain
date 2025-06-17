package types

const (
	QueryUserInfo       = "user_info"
	QueryParams         = "params"
	QueryUserInfos      = "user_infos"
	QueryUserByMobile   = "user_by_mobile"
	QueryUsersChatInfo  = "users_chat_info"
	QueryAddrByChatAddr = "addr_by_chat_addr"
)

type QueryUserInfoParams struct {
	Address string
}

type QueryUserByMobileParams struct {
	Mobile string
}

type QueryChatSendGiftInfoParams struct {
	FromAddress string
	ToAddress   string
}

type QueryUserInfosParams struct {
	Addresses []string
}

type QueryPledgeLevelsParams struct {
	Addresses []string
}

type QueryChatSendGiftsInfoParams struct {
	FromAddress string
	ToAddresses []string
}
