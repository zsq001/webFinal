package web

type Msg struct {
	Status int         `json:"status"`
	Errors string      `json:"error,omitempty"`
	Data   interface{} `json:"data,omitempty"`
}

type MsgList struct {
	ID            uint   `json:"ID"`
	Name          string `json:"name"`
	LatestMsg     string `json:"latest_msg"`
	LatestMsgTime string `json:"latest_msg_time"`
}
