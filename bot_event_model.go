package onebotClient

type botEvent struct {
	MessageType string `json:"message_type,omitempty"`
	Time        int64  `json:"time"`
	SelId       int64  `json:"sel_id"`
	PostType    string `json:"post_type,omitempty"`
}

type message struct {
	Type string                 `json:"type"`
	Data map[string]interface{} `json:"data"`
	//File    string `json:"file"`
	//Url     string `json:"url"`
	//Summary string `json:"summary"`
	//QQ      string `json:"qq"`
	//Text    string `json:"text"`
	//SubType int64    `json:"subType"`
}

type groupMessageEvent struct {
	MessageType string      `json:"message_type"`
	SubType     string      `json:"sub_type"`
	MessageId   int64       `json:"message_id"`
	GroupId     int64       `json:"group_id"`
	UserId      int64       `json:"user_id"`
	Anonymous   interface{} `json:"anonymous"`
	Message     []message   `json:"message"`
	RawMessage  string      `json:"raw_message"`
	Font        int64       `json:"font"`
	Sender      struct {
		UserId   int64  `json:"user_id"`
		Nickname string `json:"nickname"`
		Card     string `json:"card"`
		Sex      string `json:"sex"`
		Age      int64  `json:"age"`
		Area     string `json:"area"`
		Level    string `json:"level"`
		Role     string `json:"role"`
		Title    string `json:"title"`
	} `json:"sender"`
	Time     int64  `json:"time"`
	SelfId   int64  `json:"self_id"`
	PostType string `json:"post_type"`
}

type privateMessageEvent struct {
	MessageType string    `json:"message_type"`
	SubType     string    `json:"sub_type"`
	MessageId   int64     `json:"message_id"`
	UserId      int64     `json:"user_id"`
	Message     []message `json:"message"`
	RawMessage  string    `json:"raw_message"`
	Font        int64     `json:"font"`
	Sender      struct {
		UserId   int64  `json:"user_id"`
		Nickname string `json:"nickname"`
		Sex      string `json:"sex"`
	} `json:"sender"`
	TargetId int64  `json:"target_id"`
	Time     int64  `json:"time"`
	SelfId   int64  `json:"self_id"`
	PostType string `json:"post_type"`
}
