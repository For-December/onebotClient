package onebotClient

import (
	"fmt"
	"github.com/lxzan/gws"
)

type eventHandler struct {
	be *BotEngine
}

func (c *eventHandler) OnOpen(socket *gws.Conn) {
	//_ = socket.SetDeadline(time.Now().Add(PingInterval + PingWait))
}

func (c *eventHandler) OnClose(socket *gws.Conn, err error) {
	//fmt.Printf("bot event [%v] 下线: %v", c.BotAccount, err.Error())
}

func (c *eventHandler) OnPing(socket *gws.Conn, payload []byte) {
	fmt.Println(string(payload))
}

func (c *eventHandler) OnPong(socket *gws.Conn, payload []byte) {
	fmt.Println(string(payload))

}

func (c *eventHandler) OnMessage(socket *gws.Conn, message *gws.Message) {
	//println(string(message.Bytes()))
	// 收到的消息放入 botEventRawChannel，由 dispatcher 处理
	// channel 使用缓冲区，使得能够连续接收消息而不阻塞
	c.be.rawBotEventChannel <- message.Bytes()
}
