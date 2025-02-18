package msg

type GroupMessageContext struct {
	MessageContext
}

type MessageContext struct {
	BotAccount        int64              `json:"bot_account"`
	MessageType       MessageType        `json:"message_type"`
	MessageId         int64              `json:"message_id"`
	GroupMessageChain *GroupMessageChain `json:"message_chain"`

	FromId   int64 // 消息来源
	TargetId int64 // 消息去路

	// 内部路由使用
	Params map[string]string
}
