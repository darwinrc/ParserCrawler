package model

type Request struct {
	ReqId string `json:"reqId,omitempty"`
	Url   string `json:"url,omitempty"`
}

type Sitemap struct {
	Pages map[string][]string `json:"pages"`
}

type Response struct {
	Request
	Sitemap
}
