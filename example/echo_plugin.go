package main

import (
	"github.com/For-December/onebotClient"
	"github.com/For-December/onebotClient/msg"
)

type EchoPlugin struct {
}

func (e EchoPlugin) GetPaths() []string {
	//TODO implement me
	panic("implement me")
}

func (e EchoPlugin) GetPluginInfo() string {
	//TODO implement me
	panic("implement me")
}

func (e EchoPlugin) GetPluginHandler() onebotClient.PluginHandler {

	return func(api onebotClient.BotActionAPIInterface, ctx *msg.MessageContext) (done bool) {
		api.SendGroupMessage(111, ctx.GroupMessageChain)
		return true
	}
}
