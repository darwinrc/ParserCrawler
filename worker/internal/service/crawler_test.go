package service

import (
	"github.com/stretchr/testify/assert"
	"testing"
	mock_handler "worker/internal/handler/mocks"
	"worker/internal/model"
)

func TestCrawlService_Crawl_Success(t *testing.T) {
	url := "https://parserdigital.com/"
	mockHttp := mock_handler.NewHttpMock(url)
	mockHttp.RegisterResponders()

	mockSitemap := &model.Sitemap{
		Pages: make(map[string][]string),
	}

	service := &crawlService{
		visitedURLs: make(map[string]bool),
		sitemap:     mockSitemap,
	}

	//url := mockHttpServer.MockServer.URL
	sitemap := service.Crawl(url)

	expected := &model.Sitemap{
		Pages: map[string][]string{
			url: {
				url + "how-we-work",
				url + "career",
				url + "contact-us",
			},
			url + "how-we-work": {
				url + "cases",
				url + "people",
				url + "",
			},
			url + "career": {
				url + "apply",
				url + "visit",
				url + "",
			},
			url + "contact-us": {
				url + "address",
				url + "form",
				url + "",
			},
			url + "people":  []string(nil),
			url + "cases":   []string(nil),
			url + "form":    []string(nil),
			url + "address": []string(nil),
			url + "visit":   []string(nil),
			url + "apply":   []string(nil),
		},
	}

	assert.Equal(t, sitemap, expected)

	mockHttp.DeactivateAndReset()
}
