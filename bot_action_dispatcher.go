package onebotClient

import (
	"encoding/json"
	"github.com/lxzan/gws"
)

func (be *BotEngine) runActionDispatcher(actionClient *gws.Conn) {
	for {
		select {
		case botAction := <-be.botActionRequestChannel: // bot 行为

			if actionBytes, err := json.Marshal(&botAction); err != nil {
				panic(err)
			} else if err = actionClient.WriteString(string(actionBytes)); err != nil {
				panic(err)
			}
			//case <-botActionResponseChannel: // bot 行为结果
			// 其他地方接收
		}
	}
}
