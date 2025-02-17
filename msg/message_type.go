package msg

// MessageType 消息类型
type MessageType = string

const (
	GroupMsg   MessageType = "group"
	PrivateMsg MessageType = "private"
	TempMsg    MessageType = "temp"
)

const (
	GroupScope   uint32 = 1 << 0
	PrivateScope uint32 = 1 << 1
)
