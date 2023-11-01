package repo

import (
	"context"
	"errors"
	mock_infra "server/internal/infra/mocks"
	"testing"

	"github.com/golang/mock/gomock"
	"server/internal/infra"
)

func TestCrawlerRepository_GetUrl(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock_infra.NewMockRedisClient(ctrl)
	repo := NewCrawlerRepository(mockClient)

	testURL := "https://parserdigital.com/"

	t.Run("Successful GetUrl", func(t *testing.T) {
		ctx := context.Background()
		expectedValue := "test-value"

		mockClient.EXPECT().Get(ctx, testURL).Return(expectedValue, nil)

		result, err := repo.GetUrl(ctx, testURL)
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}
		if result != expectedValue {
			t.Errorf("Expected value: %s, got: %s", expectedValue, result)
		}
	})

	t.Run("Key Not Found", func(t *testing.T) {
		ctx := context.Background()

		mockClient.EXPECT().Get(ctx, testURL).Return("", errors.New(infra.RedisKeyNotFound))

		_, err := repo.GetUrl(ctx, testURL)

		expectedError := errors.New(KeyNotFound)
		if expectedError.Error() != err.Error() {
			t.Errorf("Expected error: %v, got: %v", expectedError, err)
		}
	})

	t.Run("Error from Redis Client", func(t *testing.T) {
		ctx := context.Background()

		mockClient.EXPECT().Get(ctx, testURL).Return("", errors.New("some error"))

		_, err := repo.GetUrl(ctx, testURL)

		expectedError := errors.New("error getting url from repo: some error")
		if expectedError.Error() != err.Error() {
			t.Errorf("Expected error: %v, got: %v", expectedError, err)
		}
	})
}

func TestCrawlerRepository_StoreUrl(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock_infra.NewMockRedisClient(ctrl)
	repo := NewCrawlerRepository(mockClient)

	testKey := "testKey"
	testValue := "testValue"

	t.Run("Successful StoreUrl", func(t *testing.T) {
		ctx := context.Background()

		mockClient.EXPECT().Set(ctx, testKey, testValue).Return(nil)

		err := repo.StoreUrl(ctx, testKey, testValue)

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}
	})

	t.Run("Error from Redis Client", func(t *testing.T) {
		ctx := context.Background()

		mockClient.EXPECT().Set(ctx, testKey, testValue).Return(errors.New("some error"))

		err := repo.StoreUrl(ctx, testKey, testValue)

		expectedError := errors.New("some error")
		if expectedError.Error() != err.Error() {
			t.Errorf("Expected error: %v, got: %v", expectedError, err)
		}
	})
}
