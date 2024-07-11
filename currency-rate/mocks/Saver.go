// Code generated by mockery v2.43.2. DO NOT EDIT.

package mocks

import (
	dto "github.com/AlwaysSayNo/genesis-currency-api/currency-rate/internal/shared/dto/user"
	mock "github.com/stretchr/testify/mock"
)

// Saver is an autogenerated mock type for the Saver type
type Saver struct {
	mock.Mock
}

// Save provides a mock function with given fields: _a0
func (_m *Saver) Save(_a0 dto.SaveRequestDTO) (*dto.ResponseDTO, error) {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for Save")
	}

	var r0 *dto.ResponseDTO
	var r1 error
	if rf, ok := ret.Get(0).(func(dto.SaveRequestDTO) (*dto.ResponseDTO, error)); ok {
		return rf(_a0)
	}
	if rf, ok := ret.Get(0).(func(dto.SaveRequestDTO) *dto.ResponseDTO); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*dto.ResponseDTO)
		}
	}

	if rf, ok := ret.Get(1).(func(dto.SaveRequestDTO) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewSaver creates a new instance of Saver. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewSaver(t interface {
	mock.TestingT
	Cleanup(func())
}) *Saver {
	mock := &Saver{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}