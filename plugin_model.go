package onebotClient

import "github.com/For-December/onebotClient/msg"

type PluginHandler func(
	api BotActionAPIInterface,
	ctx *msg.MessageContext,
) (done bool)

type PluginInterface interface {
	GetPaths() []string // ban [user] [duration]
	GetPluginInfo() string
	GetPluginHandler() PluginHandler
}
