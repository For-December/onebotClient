package onebotClient

import (
	"encoding/json"
	"fmt"
)

func (be *BotEngine) handleActionResult(actionResult BotActionResult) {
	if handler, ok := be.actionRequests[actionResult.Echo]; ok {
		handler(actionResult)
		delete(be.actionRequests, actionResult.Echo)
	} else {
		actionResultStr, _ := json.Marshal(actionResult)
		fmt.Println("no handler for action result: ", actionResultStr)
	}

}

func (be *BotEngine) submitAction(action BotAction, handler func(result BotActionResult)) {
	be.actionRequests[action.Echo] = handler
	actionBytes, _ := json.Marshal(action)
	if err := be.actionClient.WriteString(string(actionBytes)); err != nil {
		fmt.Println(action)
		delete(be.actionRequests, action.Echo)
		panic(err)
	}
}
