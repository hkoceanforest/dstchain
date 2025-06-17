package core

import (
	"freemasonry.cc/log"
	"github.com/sirupsen/logrus"
)

var (
	LmChainClient    = log.RegisterModule("bc-cli", logrus.DebugLevel)
	LmChainType      = log.RegisterModule("bc-ty", logrus.DebugLevel)
	LmChainKeeper    = log.RegisterModule("bc-kp", logrus.DebugLevel)
	LmChainChatQuery = log.RegisterModule("bc-qy-chat", logrus.DebugLevel)
	LmChainCommQuery = log.RegisterModule("bc-qy-comm", logrus.DebugLevel)
	LmChainDaoQuery  = log.RegisterModule("bc-qy-dao", logrus.DebugLevel)
	LmChainNet       = log.RegisterModule("bc-cn", logrus.DebugLevel)

	LmChainXContract   = log.RegisterModule("bc-xc", logrus.DebugLevel)
	LmChainCommKeeper  = log.RegisterModule("kp-comm", logrus.DebugLevel)
	LmChainChatKeeper  = log.RegisterModule("kp-chat", logrus.DebugLevel)
	LmChainMsgServer   = log.RegisterModule("bc-ms", logrus.DebugLevel)
	LmChainRest        = log.RegisterModule("bc-re", logrus.DebugLevel)
	LmChainMsgAnalysis = log.RegisterModule("bc-mas", logrus.DebugLevel)
	LmChainUtil        = log.RegisterModule("bc-ut", logrus.DebugLevel)
	LmChainBeginBlock  = log.RegisterModule("bc-bb", logrus.DebugLevel)
	LmChainDaoKeeper   = log.RegisterModule("kp-dao", logrus.DebugLevel)
	LmChainDaoServer   = log.RegisterModule("ms-dao", logrus.DebugLevel)
	LmChainDaoEvmHook  = log.RegisterModule("vm-hook-dao", logrus.DebugLevel)
)

var BuildLog = log.BuildLog
