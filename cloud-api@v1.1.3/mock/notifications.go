// Code generated by MockGen. DO NOT EDIT.
// Source: /go_workspace/src/github.com/videocoin/cloud-api/notifications/v1/notifications_service.pb.go

// Package mock is a generated GoMock package.
package mock

import (
	gomock "github.com/golang/mock/gomock"
)

// MockNotificationServiceClient is a mock of NotificationServiceClient interface
type MockNotificationServiceClient struct {
	ctrl     *gomock.Controller
	recorder *MockNotificationServiceClientMockRecorder
}

// MockNotificationServiceClientMockRecorder is the mock recorder for MockNotificationServiceClient
type MockNotificationServiceClientMockRecorder struct {
	mock *MockNotificationServiceClient
}

// NewMockNotificationServiceClient creates a new mock instance
func NewMockNotificationServiceClient(ctrl *gomock.Controller) *MockNotificationServiceClient {
	mock := &MockNotificationServiceClient{ctrl: ctrl}
	mock.recorder = &MockNotificationServiceClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockNotificationServiceClient) EXPECT() *MockNotificationServiceClientMockRecorder {
	return m.recorder
}

// MockNotificationServiceServer is a mock of NotificationServiceServer interface
type MockNotificationServiceServer struct {
	ctrl     *gomock.Controller
	recorder *MockNotificationServiceServerMockRecorder
}

// MockNotificationServiceServerMockRecorder is the mock recorder for MockNotificationServiceServer
type MockNotificationServiceServerMockRecorder struct {
	mock *MockNotificationServiceServer
}

// NewMockNotificationServiceServer creates a new mock instance
func NewMockNotificationServiceServer(ctrl *gomock.Controller) *MockNotificationServiceServer {
	mock := &MockNotificationServiceServer{ctrl: ctrl}
	mock.recorder = &MockNotificationServiceServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockNotificationServiceServer) EXPECT() *MockNotificationServiceServerMockRecorder {
	return m.recorder
}
