// Code generated by mockery v2.20.0. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	time "time"
)

// OtpRepository is an autogenerated mock type for the OtpRepository type
type OtpRepository struct {
	mock.Mock
}

// Get provides a mock function with given fields: ctx, key
func (_m *OtpRepository) Get(ctx context.Context, key string) (string, error) {
	ret := _m.Called(ctx, key)

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (string, error)); ok {
		return rf(ctx, key)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) string); ok {
		r0 = rf(ctx, key)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, key)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Save provides a mock function with given fields: ctx, key, value, ttl
func (_m *OtpRepository) Save(ctx context.Context, key string, value string, ttl time.Duration) error {
	ret := _m.Called(ctx, key, value, ttl)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string, time.Duration) error); ok {
		r0 = rf(ctx, key, value, ttl)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewOtpRepository interface {
	mock.TestingT
	Cleanup(func())
}

// NewOtpRepository creates a new instance of OtpRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewOtpRepository(t mockConstructorTestingTNewOtpRepository) *OtpRepository {
	mock := &OtpRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}