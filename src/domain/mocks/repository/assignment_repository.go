// Code generated by mockery v2.20.0. DO NOT EDIT.

package mocks

import (
	domain "github.com/superosystem/trainingsystem-backend/src/domain"
	mock "github.com/stretchr/testify/mock"
)

// AssignmentRepository is an autogenerated mock type for the AssignmentRepository type
type AssignmentRepository struct {
	mock.Mock
}

// Create provides a mock function with given fields: assignmentDomain
func (_m *AssignmentRepository) Create(assignmentDomain *domain.Assignment) error {
	ret := _m.Called(assignmentDomain)

	var r0 error
	if rf, ok := ret.Get(0).(func(*domain.Assignment) error); ok {
		r0 = rf(assignmentDomain)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Delete provides a mock function with given fields: assignmentId
func (_m *AssignmentRepository) Delete(assignmentId string) error {
	ret := _m.Called(assignmentId)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(assignmentId)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// FindByCourseId provides a mock function with given fields: courseId
func (_m *AssignmentRepository) FindByCourseId(courseId string) (*domain.Assignment, error) {
	ret := _m.Called(courseId)

	var r0 *domain.Assignment
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*domain.Assignment, error)); ok {
		return rf(courseId)
	}
	if rf, ok := ret.Get(0).(func(string) *domain.Assignment); ok {
		r0 = rf(courseId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.Assignment)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(courseId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindByCourses provides a mock function with given fields: courseIds
func (_m *AssignmentRepository) FindByCourses(courseIds []string) (*[]domain.Assignment, error) {
	ret := _m.Called(courseIds)

	var r0 *[]domain.Assignment
	var r1 error
	if rf, ok := ret.Get(0).(func([]string) (*[]domain.Assignment, error)); ok {
		return rf(courseIds)
	}
	if rf, ok := ret.Get(0).(func([]string) *[]domain.Assignment); ok {
		r0 = rf(courseIds)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*[]domain.Assignment)
		}
	}

	if rf, ok := ret.Get(1).(func([]string) error); ok {
		r1 = rf(courseIds)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindById provides a mock function with given fields: assignmentId
func (_m *AssignmentRepository) FindById(assignmentId string) (*domain.Assignment, error) {
	ret := _m.Called(assignmentId)

	var r0 *domain.Assignment
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*domain.Assignment, error)); ok {
		return rf(assignmentId)
	}
	if rf, ok := ret.Get(0).(func(string) *domain.Assignment); ok {
		r0 = rf(assignmentId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.Assignment)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(assignmentId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: assignmentId, assignmentDomain
func (_m *AssignmentRepository) Update(assignmentId string, assignmentDomain *domain.Assignment) error {
	ret := _m.Called(assignmentId, assignmentDomain)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, *domain.Assignment) error); ok {
		r0 = rf(assignmentId, assignmentDomain)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewAssignmentRepository interface {
	mock.TestingT
	Cleanup(func())
}

// NewAssignmentRepository creates a new instance of AssignmentRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewAssignmentRepository(t mockConstructorTestingTNewAssignmentRepository) *AssignmentRepository {
	mock := &AssignmentRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}