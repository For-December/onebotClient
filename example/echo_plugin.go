package main

import (
	"fmt"
	"github.com/For-December/onebotClient"
	"github.com/For-December/onebotClient/msg"
)

type EchoPlugin struct {
}

func (e EchoPlugin) GetPaths() []string {
	return []string{"!!"}
}

func (e EchoPlugin) GetPluginInfo() string {
	//TODO implement me
	return "测试插件"
}

func (e EchoPlugin) GetPluginHandler() onebotClient.PluginHandler {

	return func(api onebotClient.BotActionAPIInterface, ctx *msg.GroupMessageContext) (done bool) {

		fmt.Println(ctx.GroupMessageChain.Messages)
		api.SendGroupMessage(ctx.GroupId, ctx.GroupMessageChain)
		return true
	}
}
