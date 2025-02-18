package main

import (
	"github.com/For-December/onebotClient"
)

func main() {
	engine := onebotClient.DefaultBotEngine()
	engine.RegisterGroupPlugins(EchoPlugin{})
	engine.RunLoopWithGroups([]int64{472737616})
}
