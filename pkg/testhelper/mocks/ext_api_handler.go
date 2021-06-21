// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package mocks

import (
	bytes "bytes"
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// ExtAPIHandler is an autogenerated mock type for the ExtAPIHandler type
type ExtAPIHandler struct {
	mock.Mock
}

// SendToContentService provides a mock function with given fields: ctx, body, contentType
func (_m *ExtAPIHandler) SendToContentService(ctx context.Context, body bytes.Buffer, contentType string) error {
	ret := _m.Called(ctx, body, contentType)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, bytes.Buffer, string) error); ok {
		r0 = rf(ctx, body, contentType)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ValidateToken provides a mock function with given fields: ctx, token
func (_m *ExtAPIHandler) ValidateToken(ctx context.Context, token string) error {
	ret := _m.Called(ctx, token)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, token)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
