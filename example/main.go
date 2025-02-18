package main

import (
	"github.com/For-December/onebotClient"
)

func main() {
	engine := onebotClient.NewBotEngine("ws://127.0.0.1:8080", "")
	engine.RegisterGroupPlugins(EchoPlugin{})
	engine.RunLoopWithGroups([]int64{472737616})
}
