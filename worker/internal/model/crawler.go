package model

type Request struct {
	ReqId string `json:"reqId,omitempty"`
	Url   string `json:"url,omitempty"`
}

type Response struct {
	Request
	Response Sitemap `json:"response"`
}

type Sitemap struct {
	Pages map[string][]string `json:"pages"`
}
