package onebotClient

import (
	"github.com/For-December/onebotClient/msg"
)

type BotActionAPIInterface interface {
	GetBotAccount() int64
	SendGroupMessage(groupId int64, chain *msg.GroupMessageChain, callback ...func(messageId int64))
	SendGroupCqMessage(groupId int64, cqMessage string, callback ...func(messageId int64))
	RecallMessage(messageId int64)
	GetNextContext() *msg.GroupMessageContext
}
