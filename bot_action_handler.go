package onebotClient

import (
	"fmt"
	"github.com/lxzan/gws"
)

type actionHandler struct {
	be *BotEngine
}

func (c *actionHandler) OnOpen(socket *gws.Conn) {
	//_ = socket.SetDeadline(time.Now().Add(PingInterval + PingWait))
}

func (c *actionHandler) OnClose(socket *gws.Conn, err error) {
}

func (c *actionHandler) OnPing(socket *gws.Conn, payload []byte) {
	fmt.Println("ping", string(payload))
}

func (c *actionHandler) OnPong(socket *gws.Conn, payload []byte) {
	fmt.Println("pong", string(payload))

}

func (c *actionHandler) OnMessage(socket *gws.Conn, message *gws.Message) {

	c.be.rawBotActionResponseChannel <- message.Bytes()
}
