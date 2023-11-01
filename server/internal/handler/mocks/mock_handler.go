// Code generated by MockGen. DO NOT EDIT.
// Source: crawler.go

// Package mock_handler is a generated GoMock package.
package mock_handler

import (
	http "net/http"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	mux "github.com/gorilla/mux"
)

// MockCrawlerHandler is a mock of CrawlerHandler interface.
type MockCrawlerHandler struct {
	ctrl     *gomock.Controller
	recorder *MockCrawlerHandlerMockRecorder
}

// MockCrawlerHandlerMockRecorder is the mock recorder for MockCrawlerHandler.
type MockCrawlerHandlerMockRecorder struct {
	mock *MockCrawlerHandler
}

// NewMockCrawlerHandler creates a new mock instance.
func NewMockCrawlerHandler(ctrl *gomock.Controller) *MockCrawlerHandler {
	mock := &MockCrawlerHandler{ctrl: ctrl}
	mock.recorder = &MockCrawlerHandlerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCrawlerHandler) EXPECT() *MockCrawlerHandlerMockRecorder {
	return m.recorder
}

// Attach mocks base method.
func (m *MockCrawlerHandler) Attach(r *mux.Router) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Attach", r)
}

// Attach indicates an expected call of Attach.
func (mr *MockCrawlerHandlerMockRecorder) Attach(r interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Attach", reflect.TypeOf((*MockCrawlerHandler)(nil).Attach), r)
}

// HandleCrawl mocks base method.
func (m *MockCrawlerHandler) HandleCrawl(w http.ResponseWriter, r *http.Request) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "HandleCrawl", w, r)
}

// HandleCrawl indicates an expected call of HandleCrawl.
func (mr *MockCrawlerHandlerMockRecorder) HandleCrawl(w, r interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HandleCrawl", reflect.TypeOf((*MockCrawlerHandler)(nil).HandleCrawl), w, r)
}