package onebotClient

import (
	"encoding/json"
	"fmt"
	"github.com/For-December/onebotClient/msg"
)

func init() {

}

func (be *BotEngine) runEventDispatcher() {
	// 多个调度器处理对应channel数据，每个bot一个channel
	// 通道传参，实际上是指针传递，因为通道本事是指针
	for {
		select {
		case rawMsg := <-be.rawBotEventChannel: // 不同的bot用不同的协程
			// msg有个self id字段，因此可以不必传botAccount
			dispatcher(be, rawMsg)
		}
	}

}

func processGroupMsg(be *BotEngine, message []message, userId, groupId, messageId int64, botAccount int64) {

	if !IsTargetInArray(groupId, be.listeningGroups) {
		fmt.Println("unListen group msg, groupId is: ", groupId)
		return
	}

	groupChain := msg.NewGroupChain()

	for _, s := range message {

		switch s.Type {
		case "text":
			groupChain.Text(s.Data["text"].(string))
		case "image":
			groupChain.Image(s.Data["file"].(string))
		case "record":
			groupChain.Record(s.Data["file"].(string))
		case "at":
			groupChain.At(s.Data["qq"].(string))
		case "reply":
			groupChain.Reply(s.Data["id"].(string))
		case "face":
			groupChain.Face(s.Data["id"].(string))

		default:
			println(fmt.Sprintln("no such message type: ", s.Type))
			fmt.Println(s)
			return
		}
	}

	if _, ok := be.groupMessageChannels[groupId]; !ok {
		panic("group channel not found!")
	}

	be.groupMessageChannels[groupId] <- &msg.GroupMessageContext{
		BotAccount:        botAccount,
		MessageType:       msg.GroupMsg,
		MessageId:         messageId,
		GroupMessageChain: groupChain,
		UserId:            userId,
		GroupId:           groupId,
	}

}

// 将消息解析后调度到对应群的channel
func dispatcher(be *BotEngine, rawMsg []byte) {
	event := botEvent{}
	err := json.Unmarshal(rawMsg, &event)
	if err != nil {
		panic(fmt.Sprintln("解析事件json失败: ", string(rawMsg), err))
		return
	}

	switch event.PostType {
	case "message":
		switch event.MessageType {

		// 群组消息
		case "group":
			groupMessage := groupMessageEvent{}
			if err := json.Unmarshal(rawMsg, &groupMessage); err != nil {
				panic(err)
				return
			}

			fmt.Println("群消息👇")

			// 消息链

			processGroupMsg(be, groupMessage.Message,
				groupMessage.UserId,
				groupMessage.GroupId,
				groupMessage.MessageId,
				groupMessage.SelfId,
			)

		case "private":
			privateMessage := privateMessageEvent{}
			if err := json.Unmarshal(rawMsg, &privateMessage); err != nil {
				panic(err)
				return
			}

			switch privateMessage.SubType {
			case "friend":
				// 好友对话

			case "group":
				// 临时会话

			}

		default:
			panic(event.MessageType)

		}

	case "notice":
	case "request":
	case "meta_event":
	//events.HandleMateEvent(message)
	default:
		panic(event.PostType)

	}

}
