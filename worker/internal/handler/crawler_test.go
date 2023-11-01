package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/golang/mock/gomock"
	"log"
	"testing"
	mock_infra "worker/internal/infra/mocks"
	"worker/internal/model"
	mock_service "worker/internal/service/mocks"
)

func TestCrawlHandler_Process(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAMQPClient := mock_infra.NewMockAMQPClient(ctrl)
	mockCrawlerService := mock_service.NewMockCrawlerService(ctrl)

	handler := NewCrawlerHandler(mockAMQPClient, mockCrawlerService)
	url := "https://parserdigital.com/"
	data := &model.Sitemap{
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
	req := &model.Request{
		ReqId: "req-id",
		Url:   "https://parserdigital.com/",
	}
	res := &model.Response{
		Request: model.Request{
			ReqId: req.ReqId,
			Url:   req.Url,
		},
		Sitemap: *data,
	}

	body, _ := json.Marshal(res)

	go func() {
		mockCrawlerService.EXPECT().Crawl(gomock.Any()).Return(data)

		mockAMQPClient.EXPECT().SetupAMQExchange().Return(nil)
		mockAMQPClient.EXPECT().ConsumeAMQMessages().Return(nil, nil).AnyTimes()
		mockAMQPClient.EXPECT().PublishAMQMessage(body).Return(nil).AnyTimes()

		var logOutput bytes.Buffer
		log.SetOutput(&logOutput)

		handler.Process()

		logMsg := logOutput.String()
		expectedLog := fmt.Sprintf("published to the exchange: %v", string(body))
		if logMsg != expectedLog {
			t.Errorf("Expected log message: %s, got: %s", expectedLog, logMsg)
		}
	}()
}
