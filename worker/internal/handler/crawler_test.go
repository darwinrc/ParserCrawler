package handler

import (
	"github.com/golang/mock/gomock"
	"testing"
	"time"
	mock_infra "worker/internal/infra/mocks"
	mock_service "worker/internal/service/mocks"
)

func TestCrawlHandler_Process(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAMQPClient := mock_infra.NewMockAMQPClient(ctrl)
	mockCrawlerService := mock_service.NewMockCrawlerService(ctrl)

	handler := NewCrawlerHandler(mockAMQPClient, mockCrawlerService)

	// Create a channel to signal when the goroutine has finished
	done := make(chan struct{})

	// Set up expectations for AMQPClient methods
	mockAMQPClient.EXPECT().SetupAMQExchange().Return(nil)
	mockAMQPClient.EXPECT().ConsumeAMQMessages().Return(nil, nil).AnyTimes()
	mockAMQPClient.EXPECT().PublishAMQMessage(gomock.Any()).Return(nil).AnyTimes()

	// Set up expectations for CrawlerService methods
	mockCrawlerService.EXPECT().Crawl(gomock.Any()).Return(nil)

	// Start the Process method in a goroutine
	go func() {
		defer close(done)
		handler.Process()
	}()

	// Wait for the goroutine to finish or for a timeout
	select {
	case <-done:
		// Goroutine finished
	case <-time.After(5 * time.Second): // Adjust the timeout as needed
		t.Error("Timed out waiting for the goroutine to finish")
	}
}
