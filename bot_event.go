package onebotClient

import (
	"fmt"
	"github.com/lxzan/gws"
	"net/http"
	"os"
)

func (be *BotEngine) runEventLoop() {
	var err error
	be.eventClient, _, err = gws.NewClient(&eventHandler{be: be}, &gws.ClientOption{
		Addr: be.wsEndpoint + "/event",
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

	be.eventClient.ReadLoop()
}

func (be *BotEngine) buildCustomPluginsTrie() {
	// 树形路由匹配注册
	be.groupTrie = NewRouteTrie(CallbackFunc{})
	for _, plugin := range be.plugins {
		paths := plugin.GetPaths()
		for _, path := range paths {
			// 优先级 !! > 精确匹配 > **
			be.groupTrie.Insert(path, plugin.GetPluginHandler())
		}
	}
}
