// Code generated by mockery v2.43.2. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// EmailSender is an autogenerated mock type for the EmailSender type
type EmailSender struct {
	mock.Mock
}

// SendEmails provides a mock function with given fields: ctx
func (_m *EmailSender) SendEmails(ctx context.Context) error {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for SendEmails")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewEmailSender creates a new instance of EmailSender. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewEmailSender(t interface {
	mock.TestingT
	Cleanup(func())
}) *EmailSender {
	mock := &EmailSender{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
