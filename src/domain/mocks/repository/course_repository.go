// Code generated by mockery v2.20.0. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"
	domain "github.com/superosystem/trainingsystem-backend/src/domain"
)

// CourseRepository is an autogenerated mock type for the CourseRepository type
type CourseRepository struct {
	mock.Mock
}

// Create provides a mock function with given fields: courseDomain
func (_m *CourseRepository) Create(courseDomain *domain.Course) error {
	ret := _m.Called(courseDomain)

	var r0 error
	if rf, ok := ret.Get(0).(func(*domain.Course) error); ok {
		r0 = rf(courseDomain)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Delete provides a mock function with given fields: courseId
func (_m *CourseRepository) Delete(courseId string) error {
	ret := _m.Called(courseId)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(courseId)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// FindAll provides a mock function with given fields: keyword
func (_m *CourseRepository) FindAll(keyword string) (*[]domain.Course, error) {
	ret := _m.Called(keyword)

	var r0 *[]domain.Course
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*[]domain.Course, error)); ok {
		return rf(keyword)
	}
	if rf, ok := ret.Get(0).(func(string) *[]domain.Course); ok {
		r0 = rf(keyword)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*[]domain.Course)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(keyword)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindByCategory provides a mock function with given fields: categoryId
func (_m *CourseRepository) FindByCategory(categoryId string) (*[]domain.Course, error) {
	ret := _m.Called(categoryId)

	var r0 *[]domain.Course
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*[]domain.Course, error)); ok {
		return rf(categoryId)
	}
	if rf, ok := ret.Get(0).(func(string) *[]domain.Course); ok {
		r0 = rf(categoryId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*[]domain.Course)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(categoryId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindById provides a mock function with given fields: courseId
func (_m *CourseRepository) FindById(courseId string) (*domain.Course, error) {
	ret := _m.Called(courseId)

	var r0 *domain.Course
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*domain.Course, error)); ok {
		return rf(courseId)
	}
	if rf, ok := ret.Get(0).(func(string) *domain.Course); ok {
		r0 = rf(courseId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.Course)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(courseId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindByMentor provides a mock function with given fields: mentorId
func (_m *CourseRepository) FindByMentor(mentorId string) (*[]domain.Course, error) {
	ret := _m.Called(mentorId)

	var r0 *[]domain.Course
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*[]domain.Course, error)); ok {
		return rf(mentorId)
	}
	if rf, ok := ret.Get(0).(func(string) *[]domain.Course); ok {
		r0 = rf(mentorId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*[]domain.Course)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(mentorId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindByPopular provides a mock function with given fields:
func (_m *CourseRepository) FindByPopular() ([]domain.Course, error) {
	ret := _m.Called()

	var r0 []domain.Course
	var r1 error
	if rf, ok := ret.Get(0).(func() ([]domain.Course, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() []domain.Course); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.Course)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: courseId, courseDomain
func (_m *CourseRepository) Update(courseId string, courseDomain *domain.Course) error {
	ret := _m.Called(courseId, courseDomain)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, *domain.Course) error); ok {
		r0 = rf(courseId, courseDomain)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewCourseRepository interface {
	mock.TestingT
	Cleanup(func())
}

// NewCourseRepository creates a new instance of CourseRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewCourseRepository(t mockConstructorTestingTNewCourseRepository) *CourseRepository {
	mock := &CourseRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
