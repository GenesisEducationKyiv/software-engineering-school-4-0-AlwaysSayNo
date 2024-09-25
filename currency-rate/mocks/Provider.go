// Code generated by mockery v2.43.2. DO NOT EDIT.

package mocks

import (
	context "context"

	currency "github.com/AlwaysSayNo/genesis-currency-api/currency-rate/internal/shared/dto/currency"
	mock "github.com/stretchr/testify/mock"
)

// Provider is an autogenerated mock type for the Provider type
type Provider struct {
	mock.Mock
}

// GetCurrencyRate provides a mock function with given fields: ctx
func (_m *Provider) GetCurrencyRate(ctx context.Context) (*currency.ResponseDTO, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for GetCurrencyRate")
	}

	var r0 *currency.ResponseDTO
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) (*currency.ResponseDTO, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) *currency.ResponseDTO); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*currency.ResponseDTO)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewProvider creates a new instance of Provider. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewProvider(t interface {
	mock.TestingT
	Cleanup(func())
}) *Provider {
	mock := &Provider{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
