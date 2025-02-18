package msg

type GroupMessageContext struct {
	BotAccount        int64              `json:"bot_account"`
	MessageType       MessageType        `json:"message_type"`
	MessageId         int64              `json:"message_id"`
	GroupMessageChain *GroupMessageChain `json:"message_chain"`

	UserId  int64
	GroupId int64

	// 内部路由使用
	Params map[string]string
}
