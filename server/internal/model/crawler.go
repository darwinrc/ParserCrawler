package model

type Request struct {
	ReqId string `json:"reqId,omitempty"`
	Url   string `json:"url,omitempty"`
}

type Response struct {
	Request
	Status string `json:"status"`
	Sitemap
}

type Sitemap struct {
	Pages map[string][]string `json:"pages"`
}
