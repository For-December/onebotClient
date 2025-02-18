package onebotClient

import (
	"github.com/For-December/onebotClient/msg"
)

type BotEngine struct {
	plugins       []PluginInterface
	wsEndpoint    string
	authorization string

	rawBotEventChannel chan []byte

	botActionRequestChannel     chan BotAction
	rawBotActionResponseChannel chan []byte

	groupMessageChannels map[int64]chan *msg.GroupMessageContext
	listeningGroups      []int64

	groupTrie *RouteTrie
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

	be.botActionRequestChannel = make(chan BotAction, 1024)
	be.rawBotActionResponseChannel = make(chan []byte, 1024)

	eventClient := be.createEventClient()
	actionClient := be.createActionClient()

	go eventClient.ReadLoop()
	go actionClient.ReadLoop()

	go be.runEventDispatcher()
	go be.runActionDispatcher(actionClient)

	be.buildCustomPluginsTrie()
	be.startChannelPluginListeners()

}
