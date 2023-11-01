package service

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/streadway/amqp"
	mock_infra "server/internal/infra/mocks"
	"server/internal/model"
	"server/internal/repo"
	mock_repo "server/internal/repo/mocks"
	"testing"
)

func TestCrawlService_Crawl(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockCrawlerRepo(ctrl)
	mockAMQPClient := mock_infra.NewMockAMQPClient(ctrl)

	service := NewCrawlerService(mockRepo, mockAMQPClient)

	ctx := context.Background()
	testURL := "https://parsedigital.com/"

	t.Run("URL Found in Cache", func(t *testing.T) {
		expectedResponse := &model.Response{
			Status: "returned from cache",
		}
		data, _ := json.Marshal(expectedResponse)

		mockRepo.EXPECT().GetUrl(ctx, testURL).Return(string(data), nil)

		response, err := service.Crawl(ctx, testURL)
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		if !responseEquals(response, expectedResponse) {
			t.Errorf("Expected response: %v, got: %v", expectedResponse, response)
		}
	})

	t.Run("URL Not Found in Cache", func(t *testing.T) {
		mockRepo.EXPECT().GetUrl(ctx, testURL).Return("", errors.New(repo.KeyNotFound))

		response, err := service.Crawl(ctx, testURL)
		if err != nil && err.Error() != UrlNotFound {
			t.Errorf("Expected error: %v, got: %v", errors.New(UrlNotFound), err)
		}
		if response == nil {
			t.Error("Expected a non-nil response")
		}

		go func() {
			mockAMQPClient.EXPECT().SetupAMQExchange().Return(nil)
			mockAMQPClient.EXPECT().PublishAMQMessage(gomock.Any()).Return(nil)
		}()
	})

	t.Run("Error Getting URL from Cache", func(t *testing.T) {
		expectedError := errors.New("some error")

		mockRepo.EXPECT().GetUrl(ctx, testURL).Return("", expectedError)

		response, err := service.Crawl(ctx, testURL)

		if err.Error() != "error getting url from cache: some error" {
			t.Errorf("Expected error: %s, got: %s", "error getting url from cache: some error", err.Error())
		}
		if response != nil {
			t.Error("Expected a nil response")
		}
	})
}

func TestCrawlService_ConsumeFromResponseQueue(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repo.NewMockCrawlerRepo(ctrl)
	mockAMQPClient := mock_infra.NewMockAMQPClient(ctrl)

	service := NewCrawlerService(mockRepo, mockAMQPClient)

	ctx := context.Background()
	amqpMessages := make(chan amqp.Delivery)
	broadcast := make(chan []byte, 1)

	mockAMQPClient.EXPECT().SetupAMQExchange().Return(nil)
	mockAMQPClient.EXPECT().ConsumeAMQMessages().Return(amqpMessages, nil)
	mockRepo.EXPECT().StoreUrl(ctx, gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	go service.ConsumeFromResponseQueue(ctx, broadcast)

	testResponse := &model.Response{
		Status: "test response",
	}
	responseJSON, _ := json.Marshal(testResponse)

	amqpMessages <- amqp.Delivery{Body: responseJSON}

	close(amqpMessages)

	receivedMessage := <-broadcast

	if len(receivedMessage) == 0 {
		t.Error("Expected a non-empty response in broadcast")
	}
}

func responseEquals(a, b *model.Response) bool {
	return a.Status == b.Status && a.Request.ReqId == b.Request.ReqId && a.Request.Url == b.Request.Url
}
