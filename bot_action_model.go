package onebotClient

// BotAction 通过接口避免直接创建对象
type BotAction interface {
	GetAccount() int64
}
type botActionImpl struct {
	BotAccount int64
	Action     string      `json:"action"`
	Params     interface{} `json:"params"`
	Echo       string      `json:"echo"`
}

func (b botActionImpl) GetAccount() int64 {
	return b.BotAccount
}

func NewBotAction(botAccount int64, action string, params interface{}, echo string) BotAction {
	return botActionImpl{
		BotAccount: botAccount,
		Action:     action,
		Params:     params,
		Echo:       echo,
	}

}

type BotActionResult struct {
	Status  string `json:"status"`
	Retcode int    `json:"retcode"`
	Data    struct {
		MessageId int64 `json:"message_id"`
	} `json:"data"`
	Echo string `json:"echo"`
}

type HeartBeat struct {
	Interval int `json:"interval"`
	Status   struct {
		AppInitialized bool `json:"app_initialized"`
		AppEnabled     bool `json:"app_enabled"`
		AppGood        bool `json:"app_good"`
		Online         bool `json:"online"`
		Good           bool `json:"good"`
	} `json:"status"`
	MetaEventType string `json:"meta_event_type"`
	Time          int    `json:"time"`
	SelfId        int64  `json:"self_id"`
	PostType      string `json:"post_type"`
}
