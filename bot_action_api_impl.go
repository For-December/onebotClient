package onebotClient

import (
	"encoding/json"
	"fmt"
	"github.com/For-December/onebotClient/msg"
	"os"
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

func (receiver *botActionAPIImpl) GetNextContextWithTimeout(timeout time.Duration) *msg.GroupMessageContext {
	select {
	case ctx := <-receiver.be.groupMessageChannels[receiver.groupId]:
		return ctx
	case <-time.After(timeout):
		return nil
	}
}

type GroupHistory struct {
	MessageTyp string      `json:"message_typ"`
	SubType    string      `json:"sub_type"`
	MessageId  int         `json:"message_id"`
	GroupId    int         `json:"group_id"`
	UserId     int         `json:"user_id"`
	Anonymous  interface{} `json:"anonymous"`
	Message    []struct {
		Type string `json:"type"`
		Data struct {
			UserId string `json:"user_id"`
			Name   string `json:"name"`
		} `json:"data"`
	} `json:"message"`
	RawMessage string `json:"raw_message"`
	Font       string `json:"font"`
	Sender     struct {
		UserId   int    `json:"user_id"`
		Nickname string `json:"nickname"`
		Card     string `json:"card"`
		Sex      string `json:"sex"`
		Age      int    `json:"age"`
		Area     string `json:"area"`
		Level    string `json:"level"`
		Role     string `json:"role"`
		Title    string `json:"title"`
	} `json:"sender"`
}

func (receiver *botActionAPIImpl) FetchGroupHistories(groupId int64, messageId string) {
	type TParam struct {
		GroupId   int64  `json:"group_id"`
		MessageId string `json:"message_id"`
		Count     int64  `json:"count"`
	}

	// 微秒时间戳
	echoMsg := fmt.Sprintf("get_group_msg_history_%d", time.Now().UnixMicro())
	botAction := NewBotAction("get_group_msg_history", TParam{
		GroupId:   groupId,
		MessageId: messageId,
		Count:     20,
	}, echoMsg)

	// 发送消息
	receiver.be.submitAction(botAction, func(result BotActionResult) {
		histories := result.Data["messages"].([]interface{})
		raw, _ := json.Marshal(histories)

		err := os.WriteFile(fmt.Sprintf("grps/grp_%d.txt", groupId), raw, 0655)
		if err != nil {
			panic(err)
			return
		}
		//histories := make([]GroupHistory, 0)
		//err := json.Unmarshal([]byte(raw), &histories)
		//if err != nil {
		//	fmt.Println(err)
		//	return
		//}
		for _, h := range histories {
			fmt.Println(h)
		}

	})
}
