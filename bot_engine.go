package onebotClient

type BotEngine struct {
	plugins []PluginInterface
}

func NewBotEngine() *BotEngine {
	return &BotEngine{}
}

func (be *BotEngine) RegisterGroupPlugins(plugins ...PluginInterface) {
	be.plugins = append(be.plugins, plugins...)
}

func (be *BotEngine) RunLoop() {
}
