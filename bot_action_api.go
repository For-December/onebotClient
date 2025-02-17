package onebotClient

import "github.com/For-December/onebotClient/msg"

type BotActionAPIInterface interface {
	GetBotAccount() int64
	SendGroupMessage(chain *msg.MessageChain, callback ...func(messageId int64))
	RecallMessage(messageId int64)
	GetNextContext() *msg.MessageContext
}
