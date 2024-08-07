// Code generated by mockery v2.43.2. DO NOT EDIT.

package mocks

import (
	dto "github.com/AlwaysSayNo/genesis-currency-api/currency-rate/internal/module/user/dto"

	mock "github.com/stretchr/testify/mock"

	user "github.com/AlwaysSayNo/genesis-currency-api/currency-rate/internal/shared/dto/user"
)

// UserService is an autogenerated mock type for the UserService type
type UserService struct {
	mock.Mock
}

// GetAll provides a mock function with given fields:
func (_m *UserService) GetAll() ([]user.ResponseDTO, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetAll")
	}

	var r0 []user.ResponseDTO
	var r1 error
	if rf, ok := ret.Get(0).(func() ([]user.ResponseDTO, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() []user.ResponseDTO); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]user.ResponseDTO)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Save provides a mock function with given fields: _a0
func (_m *UserService) Save(_a0 dto.SaveRequestDTO) (*user.ResponseDTO, error) {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for Save")
	}

	var r0 *user.ResponseDTO
	var r1 error
	if rf, ok := ret.Get(0).(func(dto.SaveRequestDTO) (*user.ResponseDTO, error)); ok {
		return rf(_a0)
	}
	if rf, ok := ret.Get(0).(func(dto.SaveRequestDTO) *user.ResponseDTO); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*user.ResponseDTO)
		}
	}

	if rf, ok := ret.Get(1).(func(dto.SaveRequestDTO) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewUserService creates a new instance of UserService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewUserService(t interface {
	mock.TestingT
	Cleanup(func())
}) *UserService {
	mock := &UserService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
