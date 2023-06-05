// Code generated by mockery v2.20.0. DO NOT EDIT.

package mocks

import (
	domain "github.com/superosystem/trainingsystem-backend/src/domain"
	helper "github.com/superosystem/trainingsystem-backend/src/helper"

	mock "github.com/stretchr/testify/mock"
)

// MenteeUsecase is an autogenerated mock type for the MenteeUsecase type
type MenteeUsecase struct {
	mock.Mock
}

// FindAll provides a mock function with given fields:
func (_m *MenteeUsecase) FindAll() (*[]domain.Mentee, error) {
	ret := _m.Called()

	var r0 *[]domain.Mentee
	var r1 error
	if rf, ok := ret.Get(0).(func() (*[]domain.Mentee, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() *[]domain.Mentee); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*[]domain.Mentee)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindByCourse provides a mock function with given fields: courseId, pagination
func (_m *MenteeUsecase) FindByCourse(courseId string, pagination helper.Pagination) (*helper.Pagination, error) {
	ret := _m.Called(courseId, pagination)

	var r0 *helper.Pagination
	var r1 error
	if rf, ok := ret.Get(0).(func(string, helper.Pagination) (*helper.Pagination, error)); ok {
		return rf(courseId, pagination)
	}
	if rf, ok := ret.Get(0).(func(string, helper.Pagination) *helper.Pagination); ok {
		r0 = rf(courseId, pagination)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*helper.Pagination)
		}
	}

	if rf, ok := ret.Get(1).(func(string, helper.Pagination) error); ok {
		r1 = rf(courseId, pagination)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindById provides a mock function with given fields: id
func (_m *MenteeUsecase) FindById(id string) (*domain.Mentee, error) {
	ret := _m.Called(id)

	var r0 *domain.Mentee
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*domain.Mentee, error)); ok {
		return rf(id)
	}
	if rf, ok := ret.Get(0).(func(string) *domain.Mentee); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.Mentee)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ForgotPassword provides a mock function with given fields: forgotPassword
func (_m *MenteeUsecase) ForgotPassword(forgotPassword *domain.MenteeForgotPassword) error {
	ret := _m.Called(forgotPassword)

	var r0 error
	if rf, ok := ret.Get(0).(func(*domain.MenteeForgotPassword) error); ok {
		r0 = rf(forgotPassword)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Login provides a mock function with given fields: menteeAuth
func (_m *MenteeUsecase) Login(menteeAuth *domain.MenteeAuth) (interface{}, error) {
	ret := _m.Called(menteeAuth)

	var r0 interface{}
	var r1 error
	if rf, ok := ret.Get(0).(func(*domain.MenteeAuth) (interface{}, error)); ok {
		return rf(menteeAuth)
	}
	if rf, ok := ret.Get(0).(func(*domain.MenteeAuth) interface{}); ok {
		r0 = rf(menteeAuth)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interface{})
		}
	}

	if rf, ok := ret.Get(1).(func(*domain.MenteeAuth) error); ok {
		r1 = rf(menteeAuth)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Register provides a mock function with given fields: menteeAuth
func (_m *MenteeUsecase) Register(menteeAuth *domain.MenteeAuth) error {
	ret := _m.Called(menteeAuth)

	var r0 error
	if rf, ok := ret.Get(0).(func(*domain.MenteeAuth) error); ok {
		r0 = rf(menteeAuth)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Update provides a mock function with given fields: id, menteeDomain
func (_m *MenteeUsecase) Update(id string, menteeDomain *domain.Mentee) error {
	ret := _m.Called(id, menteeDomain)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, *domain.Mentee) error); ok {
		r0 = rf(id, menteeDomain)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// VerifyRegister provides a mock function with given fields: menteeRegister
func (_m *MenteeUsecase) VerifyRegister(menteeRegister *domain.MenteeRegister) error {
	ret := _m.Called(menteeRegister)

	var r0 error
	if rf, ok := ret.Get(0).(func(*domain.MenteeRegister) error); ok {
		r0 = rf(menteeRegister)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewMenteeUsecase interface {
	mock.TestingT
	Cleanup(func())
}

// NewMenteeUsecase creates a new instance of MenteeUsecase. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMenteeUsecase(t mockConstructorTestingTNewMenteeUsecase) *MenteeUsecase {
	mock := &MenteeUsecase{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}