// Code generated by mockery v2.43.2. DO NOT EDIT.

package mock

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	models "github.com/captechconsulting/go-microservice-templates/lambda/internal/models"
)

// MockUserLister is an autogenerated mock type for the userLister type
type MockUserLister struct {
	mock.Mock
}

type MockUserLister_Expecter struct {
	mock *mock.Mock
}

func (_m *MockUserLister) EXPECT() *MockUserLister_Expecter {
	return &MockUserLister_Expecter{mock: &_m.Mock}
}

// ListUsers provides a mock function with given fields: ctx
func (_m *MockUserLister) ListUsers(ctx context.Context) ([]models.User, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for ListUsers")
	}

	var r0 []models.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) ([]models.User, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) []models.User); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.User)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockUserLister_ListUsers_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ListUsers'
type MockUserLister_ListUsers_Call struct {
	*mock.Call
}

// ListUsers is a helper method to define mock.On call
//   - ctx context.Context
func (_e *MockUserLister_Expecter) ListUsers(ctx interface{}) *MockUserLister_ListUsers_Call {
	return &MockUserLister_ListUsers_Call{Call: _e.mock.On("ListUsers", ctx)}
}

func (_c *MockUserLister_ListUsers_Call) Run(run func(ctx context.Context)) *MockUserLister_ListUsers_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *MockUserLister_ListUsers_Call) Return(_a0 []models.User, _a1 error) *MockUserLister_ListUsers_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockUserLister_ListUsers_Call) RunAndReturn(run func(context.Context) ([]models.User, error)) *MockUserLister_ListUsers_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockUserLister creates a new instance of MockUserLister. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockUserLister(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockUserLister {
	mock := &MockUserLister{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
