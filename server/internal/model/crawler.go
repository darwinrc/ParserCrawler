package model

type Request struct {
	ReqId string `json:"reqId,omitempty"`
	Url   string `json:"url,omitempty"`
}

type Response struct {
	Request
	Response string `json:"response"`
}
