package handler

import (
	"encoding/json"
	"github.com/streadway/amqp"
	"testing"
	mock_infra "worker/internal/infra/mocks"
	"worker/internal/model"
	mock_service "worker/internal/service/mocks"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestCrawlHandler_Process(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Mocks for dependencies
	mockAMQPClient := mock_infra.NewMockAMQPClient(ctrl)
	mockCrawlerService := mock_service.NewMockCrawlerService(ctrl)

	handler := NewCrawlerHandler(mockAMQPClient, mockCrawlerService)

	// Define test data
	request := &model.Request{
		ReqId: "test-req-id",
		Url:   "https://example.com",
	}

	response := &model.Response{
		Request: *request,
		Sitemap: model.Sitemap{},
	}

	requestJSON, _ := json.Marshal(request)
	responseJSON, _ := json.Marshal(response)

	// Expectations for AMQPClient
	mockAMQPClient.EXPECT().SetupAMQExchange().Return(nil)
	mockAMQPClient.EXPECT().ConsumeAMQMessages().Return(make(<-chan amqp.Delivery), nil)
	mockAMQPClient.EXPECT().PublishAMQMessage(gomock.Any()).Return(nil).AnyTimes()

	// Expectations for CrawlerService
	mockCrawlerService.EXPECT().Crawl(request.Url).Return(&model.Sitemap{})

	handler.Process()

	// Add your assertions here based on what you expect the Process method to do.
	assert.NotNil(t, requestJSON)
	assert.NotNil(t, responseJSON)
}
