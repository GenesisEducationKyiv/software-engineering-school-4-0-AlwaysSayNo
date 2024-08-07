// Code generated by mockery v2.43.2. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// CurrencyUpdater is an autogenerated mock type for the CurrencyUpdater type
type CurrencyUpdater struct {
	mock.Mock
}

// UpdateCurrencyRates provides a mock function with given fields: ctx
func (_m *CurrencyUpdater) UpdateCurrencyRates(ctx context.Context) error {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for UpdateCurrencyRates")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewCurrencyUpdater creates a new instance of CurrencyUpdater. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewCurrencyUpdater(t interface {
	mock.TestingT
	Cleanup(func())
}) *CurrencyUpdater {
	mock := &CurrencyUpdater{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
