// Code generated by mockery v2.20.0. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"
	domain "github.com/superosystem/trainingsystem-backend/src/domain"
)

// ReviewRepository is an autogenerated mock type for the ReviewRepository type
type ReviewRepository struct {
	mock.Mock
}

// Create provides a mock function with given fields: reviewDomain
func (_m *ReviewRepository) Create(reviewDomain *domain.Review) error {
	ret := _m.Called(reviewDomain)

	var r0 error
	if rf, ok := ret.Get(0).(func(*domain.Review) error); ok {
		r0 = rf(reviewDomain)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// FindByCourse provides a mock function with given fields: courseId
func (_m *ReviewRepository) FindByCourse(courseId string) ([]domain.Review, error) {
	ret := _m.Called(courseId)

	var r0 []domain.Review
	var r1 error
	if rf, ok := ret.Get(0).(func(string) ([]domain.Review, error)); ok {
		return rf(courseId)
	}
	if rf, ok := ret.Get(0).(func(string) []domain.Review); ok {
		r0 = rf(courseId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.Review)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(courseId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewReviewRepository interface {
	mock.TestingT
	Cleanup(func())
}

// NewReviewRepository creates a new instance of ReviewRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewReviewRepository(t mockConstructorTestingTNewReviewRepository) *ReviewRepository {
	mock := &ReviewRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
