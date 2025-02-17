package main

import (
	"github.com/For-December/onebotClient"
)

func main() {
	engine := onebotClient.NewBotEngine()
	engine.RegisterGroupPlugins(EchoPlugin{})
	engine.RunLoop()
}
