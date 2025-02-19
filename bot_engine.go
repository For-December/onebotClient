package onebotClient

import (
	"github.com/For-December/onebotClient/msg"
	"github.com/lxzan/gws"
)

type BotEngine struct {
	plugins       []PluginInterface
	wsEndpoint    string
	authorization string

	actionClient *gws.Conn
	eventClient  *gws.Conn

	rawBotEventChannel chan []byte

	//botActionRequestChannel     chan BotAction
	rawBotActionChannel chan []byte

	groupMessageChannels map[int64]chan *msg.GroupMessageContext
	listeningGroups      []int64

	groupTrie *RouteTrie

	actionRequests map[string]func(botAction BotActionResult)
}

func DefaultBotEngine() *BotEngine {
	return &BotEngine{
		wsEndpoint:    "/onebot/v11/ws",
		authorization: "Bearer test-114514",
	}
}
func NewBotEngine(wsEndpoint, authorization string) *BotEngine {
	return &BotEngine{
		wsEndpoint:    wsEndpoint,
		authorization: authorization,
	}
}

func (be *BotEngine) RegisterGroupPlugins(plugins ...PluginInterface) {
	be.plugins = append(be.plugins, plugins...)
}

func (be *BotEngine) RunLoopWithGroups(listeningGroups []int64) {
	if len(listeningGroups) == 0 {
		panic("listeningGroups must not be empty")
	}

	be.listeningGroups = listeningGroups
	be.groupMessageChannels = make(map[int64]chan *msg.GroupMessageContext)
	for _, group := range be.listeningGroups {
		be.groupMessageChannels[group] = make(chan *msg.GroupMessageContext, 1024)
	}

	be.rawBotEventChannel = make(chan []byte, 4096)

	//be.botActionRequestChannel = make(chan BotAction, 1024)
	be.rawBotActionChannel = make(chan []byte, 1024)

	be.actionRequests = make(map[string]func(botAction BotActionResult))

	go be.runEventLoop()
	go be.runActionLoop()

	go be.runDispatcherLoop()

	be.buildCustomPluginsTrie()
	be.startChannelPluginListeners()

}
