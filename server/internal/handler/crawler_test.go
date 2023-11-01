package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"server/internal/model"
	mock_service "server/internal/service/mocks"
	"testing"

	"github.com/golang/mock/gomock"
	"server/internal/service"
)

func TestCrawlerHandler_HandleCrawl(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_service.NewMockCrawlerService(ctrl)

	handler := NewCrawlerHandler(mockService)

	url := "https://parserdigital.com/"

	t.Run("Successful Crawl", func(t *testing.T) {
		req, err := http.NewRequest("GET", fmt.Sprintf("/crawl?url=%s", url), nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()

		request := &model.Request{
			ReqId: "test-req-id",
			Url:   url,
		}

		sitemap := model.Sitemap{
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

		expectedResponse := &model.Response{
			Request: *request,
			Sitemap: sitemap,
		}

		mockService.EXPECT().Crawl(gomock.Any(), gomock.Any()).Return(expectedResponse, nil)

		handler.HandleCrawl(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, rr.Code)
		}

		expected, err := json.Marshal(expectedResponse)
		if err != nil {
			t.Errorf("Error marshaling expected response: %s", err)
		}

		assert.Equal(t, rr.Body.Bytes(), expected)
	})

	t.Run("URL Not Found In Cache", func(t *testing.T) {
		req, err := http.NewRequest("GET", fmt.Sprintf("/crawl?url=%s", url), nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()

		mockService.EXPECT().Crawl(gomock.Any(), gomock.Any()).Return(nil, errors.New(service.UrlNotFound))

		handler.HandleCrawl(rr, req)

		if rr.Code != http.StatusAccepted {
			t.Errorf("Expected status code %d, got %d", http.StatusAccepted, rr.Code)
		}
	})

	t.Run("Error in Crawl", func(t *testing.T) {
		req, err := http.NewRequest("GET", fmt.Sprintf("/crawl?url=%s", url), nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()

		res := &model.Response{}

		mockService.EXPECT().Crawl(gomock.Any(), gomock.Any()).Return(res, errors.New("test error"))

		handler.HandleCrawl(rr, req)

		if rr.Code != http.StatusInternalServerError {
			t.Errorf("Expected status code %d, got %d", http.StatusInternalServerError, rr.Code)
		}
	})
}
