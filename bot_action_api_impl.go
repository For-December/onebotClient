package onebotClient

import (
	"fmt"
	"github.com/For-December/onebotClient/msg"
	"time"
)

type botActionAPIImpl struct {
	groupId int64
	be      *BotEngine
}

var _ BotActionAPIInterface = (*botActionAPIImpl)(nil)

func NewBotActionAPI(groupId int64, be *BotEngine) BotActionAPIInterface {
	return &botActionAPIImpl{
		groupId: groupId,
		be:      be,
	}
}

func (receiver *botActionAPIImpl) SendGroupCqMessage(groupId int64, cqMessage string, callback ...func(messageId int64)) {

	type TParam struct {
		GroupId    int64  `json:"group_id"`
		Message    string `json:"message"`
		AutoEscape bool   `json:"auto_escape"`
	}

	// 微秒时间戳
	echoMsg := fmt.Sprintf("group_message_%d", time.Now().UnixMicro())
	botAction := NewBotAction("send_group_msg", TParam{
		GroupId:    groupId,
		Message:    cqMessage,
		AutoEscape: false,
	}, echoMsg)

	// 发送消息
	receiver.be.submitAction(botAction, func(result BotActionResult) {
		if result.Echo == echoMsg {
			// 执行回调
			fmt.Println(result.Data)
			for _, f := range callback {
				f(result.Data["message_id"].(int64))
			}
		}
	})

}

func (receiver *botActionAPIImpl) SendGroupMessage(
	groupId int64, chain *msg.GroupMessageChain, callback ...func(messageId int64)) {

	type TParam struct {
		GroupId    int64                 `json:"group_id"`
		Message    []msg.JsonTypeMessage `json:"message"`
		AutoEscape bool                  `json:"auto_escape"`
	}

	// 微秒时间戳
	echoMsg := fmt.Sprintf("group_message_%d", time.Now().UnixMicro())
	botAction := NewBotAction("send_group_msg", TParam{
		GroupId:    groupId,
		Message:    chain.ToJsonTypeMessage(),
		AutoEscape: false,
	}, echoMsg)

	// 发送消息
	receiver.be.submitAction(botAction, func(result BotActionResult) {
		if result.Echo == echoMsg {
			// 执行回调
			fmt.Printf("执行回调，data为：%v\n", result.Data)
			for _, f := range callback {
				f(result.Data["message_id"].(int64))
			}
		}
	})

}

func (receiver *botActionAPIImpl) RecallMessage(messageId int64) {

	echoMsg := fmt.Sprintf("recall_%d", messageId)

	receiver.be.submitAction(NewBotAction("delete_msg", map[string]int64{
		"message_id": messageId,
	}, echoMsg), func(result BotActionResult) {})
}

func (receiver *botActionAPIImpl) GetNextContext() *msg.GroupMessageContext {
	select {
	case ctx := <-receiver.be.groupMessageChannels[receiver.groupId]:
		return ctx
	}

}
