package msg

import "fmt"

func NewGroupChain() *GroupMessageChain {
	return &GroupMessageChain{
		messages: []CommonMessage{},
	}
}

type GroupMessageChain struct {
	messages []CommonMessage
}

type JsonTypeMessage struct {
	Type string                 `json:"type"`
	Data map[string]interface{} `json:"data"`
}

func (receiver *GroupMessageChain) ToPath() string {
	resStr := ""
	for _, message := range receiver.messages {
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

func (receiver *GroupMessageChain) Size() int {
	return len(receiver.messages)
}

func (receiver *GroupMessageChain) Get(idx int) CommonMessage {
	if idx < 0 || idx >= len(receiver.messages) {
		fmt.Println("Get message chain out of range: ", idx)
		return CommonMessage{}
	}
	return receiver.messages[idx]
}

func (receiver *GroupMessageChain) ToCqStr() string {
	resStr := ""
	for _, message := range receiver.messages {
		switch message.MessageType {
		case "text":
			resStr += message.MessageContent["text"].(string)
		case "image":
			resStr += "[CQ:image,file=" + message.MessageContent["file"].(string) + "]"
		case "record":
			resStr += "[CQ:record,file=" + message.MessageContent["file"].(string) + "]"
		case "at":
			resStr += "[CQ:at,qq=" + message.MessageContent["qq"].(string) + "]"
		case "reply":
			resStr += "[CQ:reply,id=" + message.MessageContent["id"].(string) + "]"
		case "face":
			resStr += "[CQ:face,id=" + message.MessageContent["id"].(string) + "]"
		}
		//resStr += "\n"
	}
	return resStr
}

func (receiver *GroupMessageChain) ToString() string {
	resStr := ""
	for _, message := range receiver.messages {
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

func (receiver *GroupMessageChain) ToJsonTypeMessage() []JsonTypeMessage {
	message := make([]JsonTypeMessage, 0)

	for _, commonMessage := range receiver.messages {
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

func (receiver *GroupMessageChain) appendByType(messageType, messageKey, messageValue string) *GroupMessageChain {
	receiver.messages = append(receiver.messages, CommonMessage{
		MessageType:    messageType,
		MessageContent: map[string]interface{}{messageKey: messageValue},
	})
	return receiver
}

func (receiver *GroupMessageChain) Text(content string) *GroupMessageChain {
	return receiver.appendByType("text", "text", content)
}

func (receiver *GroupMessageChain) Image(file string) *GroupMessageChain {
	return receiver.appendByType("image", "file", file)
}

func (receiver *GroupMessageChain) Record(file string) *GroupMessageChain {
	return receiver.appendByType("record", "file", file)
}

func (receiver *GroupMessageChain) At(qq string) *GroupMessageChain {
	return receiver.appendByType("at", "qq", qq)
}

// Reply 通过消息id回复
func (receiver *GroupMessageChain) Reply(id string) *GroupMessageChain {
	return receiver.appendByType("reply", "id", id)
}

func (receiver *GroupMessageChain) Face(id string) *GroupMessageChain {

	// 关于id和表情的对应
	// https://github.com/kyubotics/coolq-http-api/wiki/%E8%A1%A8%E6%83%85-CQ-%E7%A0%81-ID-%E8%A1%A8
	return receiver.appendByType("face", "id", id)
}
