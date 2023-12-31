// Code generated by MockGen. DO NOT EDIT.
// Source: crawler.go

// Package mock_service is a generated GoMock package.
package mock_service

import (
	context "context"
	reflect "reflect"
	model "server/internal/model"

	gomock "github.com/golang/mock/gomock"
)

// MockCrawlerService is a mock of CrawlerService interface.
type MockCrawlerService struct {
	ctrl     *gomock.Controller
	recorder *MockCrawlerServiceMockRecorder
}

// MockCrawlerServiceMockRecorder is the mock recorder for MockCrawlerService.
type MockCrawlerServiceMockRecorder struct {
	mock *MockCrawlerService
}

// NewMockCrawlerService creates a new mock instance.
func NewMockCrawlerService(ctrl *gomock.Controller) *MockCrawlerService {
	mock := &MockCrawlerService{ctrl: ctrl}
	mock.recorder = &MockCrawlerServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCrawlerService) EXPECT() *MockCrawlerServiceMockRecorder {
	return m.recorder
}

// ConsumeFromResponseQueue mocks base method.
func (m *MockCrawlerService) ConsumeFromResponseQueue(ctx context.Context, broadcast chan []byte) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "ConsumeFromResponseQueue", ctx, broadcast)
}

// ConsumeFromResponseQueue indicates an expected call of ConsumeFromResponseQueue.
func (mr *MockCrawlerServiceMockRecorder) ConsumeFromResponseQueue(ctx, broadcast interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ConsumeFromResponseQueue", reflect.TypeOf((*MockCrawlerService)(nil).ConsumeFromResponseQueue), ctx, broadcast)
}

// Crawl mocks base method.
func (m *MockCrawlerService) Crawl(ctx context.Context, url string) (*model.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Crawl", ctx, url)
	ret0, _ := ret[0].(*model.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Crawl indicates an expected call of Crawl.
func (mr *MockCrawlerServiceMockRecorder) Crawl(ctx, url interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Crawl", reflect.TypeOf((*MockCrawlerService)(nil).Crawl), ctx, url)
}
