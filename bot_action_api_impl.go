package onebotClient

import (
	"encoding/json"
	"fmt"
	"github.com/For-December/onebotClient/msg"
	"time"
)

type botActionAPIImpl struct {
	botAccount int64
	groupId    int64
	be         *BotEngine
}

var _ BotActionAPIInterface = (*botActionAPIImpl)(nil)

func NewBotActionAPI(groupId int64, be *BotEngine) BotActionAPIInterface {
	return &botActionAPIImpl{
		botAccount: 111,
		groupId:    groupId,
		be:         be,
	}
}

func (receiver *botActionAPIImpl) GetBotAccount() int64 {
	return receiver.botAccount
}

func (receiver *botActionAPIImpl) SendGroupCqMessage(groupId int64, cqMessage string, callback ...func(messageId int64)) {

	type TParam struct {
		GroupId    int64  `json:"group_id"`
		Message    string `json:"message"`
		AutoEscape bool   `json:"auto_escape"`
	}

	// 微秒时间戳，带上bot标识，理论上不会重复
	echoMsg := fmt.Sprintf("group_message_%d_%d", receiver.botAccount, time.Now().UnixMicro())
	botAction := NewBotAction(
		receiver.GetBotAccount(),
		"send_group_msg",
		TParam{
			GroupId:    groupId,
			Message:    cqMessage,
			AutoEscape: false,
		},
		echoMsg)

	// 发送消息
	receiver.be.botActionRequestChannel <- botAction

	// 如果设置了回调则处理结果
	if len(callback) > 0 {
		receiver.solveSentRes(echoMsg, callback)
	}

}

func (receiver *botActionAPIImpl) SendGroupMessage(
	groupId int64, chain *msg.GroupMessageChain, callback ...func(messageId int64)) {

	type TParam struct {
		GroupId    int64                 `json:"group_id"`
		Message    []msg.JsonTypeMessage `json:"message"`
		AutoEscape bool                  `json:"auto_escape"`
	}

	// 微秒时间戳，带上bot标识，理论上不会重复
	echoMsg := fmt.Sprintf("group_message_%d_%d", receiver.botAccount, time.Now().UnixMicro())
	botAction := NewBotAction(
		receiver.GetBotAccount(),
		"send_group_msg",
		TParam{
			GroupId:    groupId,
			Message:    chain.ToJsonTypeMessage(),
			AutoEscape: false,
		},
		echoMsg)

	// 发送消息
	receiver.be.botActionRequestChannel <- botAction

	// 如果设置了回调则处理结果
	if len(callback) > 0 {
		receiver.solveSentRes(echoMsg, callback)
	}

}
func (receiver *botActionAPIImpl) solveSentRes(echoMsg string, callback []func(messageId int64)) {
	// 发完消息后处理结果
	go func() {
		for {
			select {
			case actionResult := <-receiver.be.rawBotActionResponseChannel:
				event := BotActionResult{}

				if err := json.Unmarshal(actionResult, &event); err != nil {
					// 心跳包
					continue
				}

				if event.Echo == echoMsg {

					// 执行回调，结束协程
					fmt.Println(event.Data)
					for _, f := range callback {
						f(event.Data.MessageId)
					}

					return
				}

				// 不是所需要的，重新放入channel
				receiver.be.rawBotActionResponseChannel <- actionResult
				continue

			}
		}
	}()
}

func (receiver *botActionAPIImpl) RecallMessage(messageId int64) {

	echoMsg := fmt.Sprintf("recall_%d", messageId)

	receiver.be.botActionRequestChannel <- NewBotAction(
		receiver.botAccount,
		"delete_msg",
		map[string]int64{
			"message_id": messageId,
		},
		echoMsg)
}

func (receiver *botActionAPIImpl) GetNextContext() *msg.GroupMessageContext {
	select {
	case ctx := <-receiver.be.groupMessageChannels[receiver.groupId]:
		return ctx
	}

}
