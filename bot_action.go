package onebotClient

import (
	"fmt"
	"github.com/lxzan/gws"
	"net/http"
	"os"
)

func (be *BotEngine) createActionClient() *gws.Conn {
	client, _, err := gws.NewClient(&actionHandler{be: be}, &gws.ClientOption{
		Addr: be.wsEndpoint + "/api",
		RequestHeader: http.Header{
			"Authorization": []string{be.authorization},
		},
		ParallelEnabled: false, // 禁用并发(内置并发实现频繁创建协程，不太合适)
		Logger:          nil,
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return client
}

func (be *BotEngine) startChannelPluginListeners() {
	for _, grp := range be.listeningGroups {
		go func(group int64) {
			api := NewBotActionAPI(group, be)
			for {
				// 每个群对应一个goroutine
				ctx := api.GetNextContext()

				switch ctx.MessageType {
				case "group":
					be.groupTrie.SearchAndExec(api, ctx)
				case "private":

					// 如果符合过滤器条件，执行过滤器，否则执行函数
					//privateTrie.SearchAndExec(api, ctx)
				case "temp":
				default:

					panic(fmt.Sprintln("unknown message type: ", ctx.MessageType))

				}
			}
		}(grp)
	}
}
