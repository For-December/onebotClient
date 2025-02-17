package msg

type GroupMessageChain struct {
	*MessageChain
}

func NewGroupChain(targetId int64) *MessageChain {
	return &MessageChain{
		Messages: []CommonMessage{},
		fromId:   0,
		targetId: targetId,
	}
}

func NewPrivateChain(targetId int64) *MessageChain {
	return &MessageChain{
		Messages: []CommonMessage{},
		fromId:   0,
		targetId: targetId,
	}
}

// NewReceivedChain 用于构建接收到的消息链
// fromId 为发送者的QQ号
// targetId 为接收者的QQ号
func NewReceivedChain(fromId int64, targetId int64) *MessageChain {
	return &MessageChain{
		Messages: []CommonMessage{},
		fromId:   fromId,
		targetId: targetId,
	}
}

type MessageChain struct {
	Messages []CommonMessage `json:"messages"`

	fromId   int64 // 消息来源
	targetId int64 // 消息去路
}

type JsonTypeMessage struct {
	Type string                 `json:"type"`
	Data map[string]interface{} `json:"data"`
}

func (receiver *MessageChain) GetFromId() int64 {
	return receiver.fromId
}

func (receiver *MessageChain) GetTargetId() int64 {
	return receiver.targetId
}

func (receiver *MessageChain) ToPath() string {
	resStr := ""
	for _, message := range receiver.Messages {
		switch message.MessageType {
		case "text":
			resStr += " " + message.MessageContent["text"].(string)
		case "image":
			resStr += " # " + message.MessageContent["file"].(string) + ""
		case "record":
			resStr += " <-record->[" + message.MessageContent["file"].(string) + "]"
		case "at":
			resStr += " @ " + message.MessageContent["qq"].(string) + ""
		case "reply":
			resStr += " & " + message.MessageContent["id"].(string) + ""
		case "face":
			resStr += " face(" + message.MessageContent["id"].(string) + ")"
		}
		//resStr += "\n"
	}
	return resStr
}

func (receiver *MessageChain) ToString() string {
	resStr := ""
	for _, message := range receiver.Messages {
		switch message.MessageType {
		case "text":
			resStr += message.MessageContent["text"].(string)
		case "image":
			resStr += "<-image->[" + message.MessageContent["file"].(string) + "]"
		case "record":
			resStr += "<-record->[" + message.MessageContent["file"].(string) + "]"
		case "at":
			resStr += "@(" + message.MessageContent["qq"].(string) + ")"
		case "reply":
			resStr += "reply(" + message.MessageContent["id"].(string) + ")"
		case "face":
			resStr += "face(" + message.MessageContent["id"].(string) + ")"
		}
		//resStr += "\n"
	}
	return resStr
}

func (receiver *MessageChain) ToJsonTypeMessage() []JsonTypeMessage {
	message := make([]JsonTypeMessage, 0)

	for _, commonMessage := range receiver.Messages {
		switch commonMessage.MessageType {
		case "text":
			message = append(message, JsonTypeMessage{
				Type: "text",
				Data: map[string]interface{}{"text": commonMessage.MessageContent["text"].(string)},
			})
		case "image":
			message = append(message, JsonTypeMessage{
				Type: "image",
				Data: map[string]interface{}{"file": commonMessage.MessageContent["file"].(string)},
			})
		case "record":
			message = append(message, JsonTypeMessage{
				Type: "record",
				Data: map[string]interface{}{"file": commonMessage.MessageContent["file"].(string)},
			})

		case "at":
			message = append(message, JsonTypeMessage{
				Type: "at",
				Data: map[string]interface{}{"qq": commonMessage.MessageContent["qq"].(string)},
			})
		case "reply":
			message = append(message, JsonTypeMessage{
				Type: "reply",
				Data: map[string]interface{}{"id": commonMessage.MessageContent["id"].(string)},
			})
		case "face":
			message = append(message, JsonTypeMessage{
				Type: "face",
				Data: map[string]interface{}{"id": commonMessage.MessageContent["id"].(string)},
			})

		}
	}

	return message
}

func (receiver *MessageChain) appendByType(messageType, messageKey, messageValue string) *MessageChain {
	receiver.Messages = append(receiver.Messages, CommonMessage{
		MessageType:    messageType,
		MessageContent: map[string]interface{}{messageKey: messageValue},
	})
	return receiver
}

func (receiver *MessageChain) Text(content string) *MessageChain {
	return receiver.appendByType("text", "text", content)
}

func (receiver *MessageChain) Image(file string) *MessageChain {
	return receiver.appendByType("image", "file", file)
}

func (receiver *MessageChain) Record(file string) *MessageChain {
	return receiver.appendByType("record", "file", file)
}

func (receiver *MessageChain) At(qq string) *MessageChain {
	return receiver.appendByType("at", "qq", qq)
}

// Reply 通过消息id回复
func (receiver *MessageChain) Reply(id string) *MessageChain {
	return receiver.appendByType("reply", "id", id)
}

func (receiver *MessageChain) Face(id string) *MessageChain {

	// 关于id和表情的对应
	// https://github.com/kyubotics/coolq-http-api/wiki/%E8%A1%A8%E6%83%85-CQ-%E7%A0%81-ID-%E8%A1%A8
	return receiver.appendByType("face", "id", id)
}
